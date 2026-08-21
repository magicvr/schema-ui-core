package objectstore

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ---- stub: minimal s3API over an in-memory map ----

type stubAPIError struct{ code string }

func (e *stubAPIError) ErrorCode() string            { return e.code }
func (e *stubAPIError) ErrorMessage() string         { return e.code }
func (e *stubAPIError) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }
func (e *stubAPIError) Error() string                { return "s3 stub: " + e.code }

type s3Object struct {
	body []byte
	meta map[string]string
}

type fakeS3 struct {
	objects   map[string]s3Object
	bucketErr error
	calls     int
	pageSize  int // list page size for pagination tests (0 = single page)
}

func newFakeS3() *fakeS3 { return &fakeS3{objects: map[string]s3Object{}} }

func (f *fakeS3) PutObject(_ context.Context, params *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	f.calls++
	body, err := io.ReadAll(params.Body)
	if err != nil {
		return nil, err
	}
	f.objects[*params.Key] = s3Object{body: body, meta: params.Metadata}
	return &s3.PutObjectOutput{}, nil
}

func (f *fakeS3) GetObject(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	f.calls++
	obj, ok := f.objects[*params.Key]
	if !ok {
		return nil, &stubAPIError{code: "NoSuchKey"}
	}
	return &s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader(obj.body)), Metadata: obj.meta}, nil
}

func (f *fakeS3) HeadObject(_ context.Context, params *s3.HeadObjectInput, _ ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	f.calls++
	obj, ok := f.objects[*params.Key]
	if !ok {
		return nil, &stubAPIError{code: "NotFound"}
	}
	size := int64(len(obj.body))
	return &s3.HeadObjectOutput{ContentLength: &size, Metadata: obj.meta}, nil
}

func (f *fakeS3) DeleteObject(_ context.Context, params *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	f.calls++
	delete(f.objects, *params.Key)
	return &s3.DeleteObjectOutput{}, nil
}

func (f *fakeS3) ListObjectsV2(_ context.Context, params *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	f.calls++
	var keys []*string
	for k := range f.objects {
		if strings.HasPrefix(k, *params.Prefix) {
			key := k
			keys = append(keys, &key)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return *keys[i] < *keys[j] })
	if params.ContinuationToken != nil {
		tok := *params.ContinuationToken
		for i, k := range keys {
			if *k == tok {
				keys = keys[i+1:]
				break
			}
		}
	}
	page := len(keys)
	if f.pageSize > 0 && f.pageSize < page {
		page = f.pageSize
	}
	out := &s3.ListObjectsV2Output{}
	for _, k := range keys[:page] {
		out.Contents = append(out.Contents, s3types.Object{Key: k})
	}
	if f.pageSize > 0 && page < len(keys) {
		truncated := true
		out.IsTruncated = &truncated
		out.NextContinuationToken = keys[page-1]
	}
	return out, nil
}

func (f *fakeS3) HeadBucket(_ context.Context, _ *s3.HeadBucketInput, _ ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	f.calls++
	if f.bucketErr != nil {
		return nil, f.bucketErr
	}
	return &s3.HeadBucketOutput{}, nil
}

// ---- tests ----

func newTestS3() *S3Store { return &S3Store{client: newFakeS3(), bucket: "test-bucket"} }

func TestS3RoundTrip(t *testing.T) {
	s := newTestS3()
	ctx := context.Background()

	if err := s.Put(ctx, kernel.ObjectNamespaceUploads, idA, []byte("hello"), kernel.ObjectMeta{
		Name: "notes.txt", Type: "text/plain", Owner: "user-1",
	}); err != nil {
		t.Fatalf("put: %v", err)
	}
	body, meta, err := s.Get(ctx, kernel.ObjectNamespaceUploads, idA)
	if err != nil || string(body) != "hello" || meta.Name != "notes.txt" || meta.Owner != "user-1" {
		t.Fatalf("get = %q/%+v/%v", body, meta, err)
	}
	info, err := s.Stat(ctx, kernel.ObjectNamespaceUploads, idA)
	if err != nil || info.Size != 5 || info.Meta.Owner != "user-1" {
		t.Fatalf("stat = %+v/%v", info, err)
	}
	if ok, _ := s.Exists(ctx, kernel.ObjectNamespaceUploads, idA); !ok {
		t.Fatal("exists = false, want true")
	}
	// upsert replaces body and metadata
	if err := s.Put(ctx, kernel.ObjectNamespaceUploads, idA, []byte("v2"), kernel.ObjectMeta{Type: "application/json"}); err != nil {
		t.Fatal(err)
	}
	body, meta, _ = s.Get(ctx, kernel.ObjectNamespaceUploads, idA)
	if string(body) != "v2" || meta.Type != "application/json" || meta.Name != "" {
		t.Fatalf("after upsert = %q/%+v", body, meta)
	}
	// delete idempotent + sentinel
	for i := 0; i < 2; i++ {
		if err := s.Delete(ctx, kernel.ObjectNamespaceUploads, idA); err != nil {
			t.Fatalf("delete #%d: %v", i+1, err)
		}
	}
	if _, _, err := s.Get(ctx, kernel.ObjectNamespaceUploads, idA); !errors.Is(err, kernel.ErrObjectNotFound) {
		t.Fatalf("get after delete err = %v", err)
	}
	if ok, err := s.Exists(ctx, kernel.ObjectNamespaceUploads, idA); ok || err != nil {
		t.Fatalf("exists after delete = %v/%v", ok, err)
	}
}

func TestS3NotFoundSentinel(t *testing.T) {
	s := newTestS3()
	ctx := context.Background()
	if _, _, err := s.Get(ctx, kernel.ObjectNamespaceAvatars, idB); !errors.Is(err, kernel.ErrObjectNotFound) {
		t.Fatalf("get miss = %v", err)
	}
	if _, err := s.Stat(ctx, kernel.ObjectNamespaceAvatars, idB); !errors.Is(err, kernel.ErrObjectNotFound) {
		t.Fatalf("stat miss = %v", err)
	}
	if ok, err := s.Exists(ctx, kernel.ObjectNamespaceAvatars, idB); ok || err != nil {
		t.Fatalf("exists miss = %v/%v", ok, err)
	}
}

// Validation happens before any network call: the fake counts invocations.
func TestS3ValidationFailClosed(t *testing.T) {
	f := newFakeS3()
	s := &S3Store{client: f, bucket: "b"}
	ctx := context.Background()
	badNS := kernel.ObjectNamespace("../escape")
	if err := s.Put(ctx, badNS, idA, []byte("x"), kernel.ObjectMeta{}); err == nil {
		t.Fatal("unknown namespace must fail")
	}
	for _, bad := range []string{"", "../" + idA, strings.ToUpper(idA), idA[:31]} {
		if err := s.Put(ctx, kernel.ObjectNamespaceUploads, bad, []byte("x"), kernel.ObjectMeta{}); err == nil {
			t.Fatalf("invalid id %q accepted", bad)
		}
	}
	if f.calls != 0 {
		t.Fatalf("rejected calls reached the backend %d times", f.calls)
	}
	if _, err := s.List(ctx, badNS); err == nil {
		t.Fatal("list with unknown namespace must fail")
	}
}

func TestS3ListPaginationAndIsolation(t *testing.T) {
	f := newFakeS3()
	f.pageSize = 1
	s := &S3Store{client: f, bucket: "b"}
	ctx := context.Background()
	for _, id := range []string{idB, idA} {
		if err := s.Put(ctx, kernel.ObjectNamespaceBrandAssets, id, []byte("x"), kernel.ObjectMeta{}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Put(ctx, kernel.ObjectNamespaceUploads, idA, []byte("other"), kernel.ObjectMeta{}); err != nil {
		t.Fatal(err)
	}
	got, err := s.List(ctx, kernel.ObjectNamespaceBrandAssets)
	if err != nil || len(got) != 2 || got[0] != idA || got[1] != idB {
		t.Fatalf("paginated list = %v err %v", got, err)
	}
	// cross-namespace isolation
	if err := s.Delete(ctx, kernel.ObjectNamespaceBrandAssets, idA); err != nil {
		t.Fatal(err)
	}
	if ok, _ := s.Exists(ctx, kernel.ObjectNamespaceUploads, idA); !ok {
		t.Fatal("cross-namespace delete leaked")
	}
	// missing namespace = empty list
	got, err = s.List(ctx, kernel.ObjectNamespaceAvatars)
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("missing ns list = %v err %v", got, err)
	}
}

func TestS3Ping(t *testing.T) {
	f := newFakeS3()
	s := &S3Store{client: f, bucket: "b"}
	if err := s.Ping(context.Background()); err != nil {
		t.Fatalf("ping = %v", err)
	}
	f.bucketErr = &stubAPIError{code: "AccessDenied"}
	if err := s.Ping(context.Background()); err == nil {
		t.Fatal("ping with backend error must fail")
	}
}

func TestNewS3RejectsEmptyBucket(t *testing.T) {
	if _, err := NewS3("http://x", "us-east-1", "  ", "k", "s", true); err == nil {
		t.Fatal("empty bucket must be rejected")
	}
}