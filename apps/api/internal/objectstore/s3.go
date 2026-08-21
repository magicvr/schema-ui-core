package objectstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// S3Store implements kernel.ObjectStore on any S3-compatible backend
// (workspace-014 GOAL-003 D-001, per Root D-004/D-005): one configured
// bucket, keys "<namespace>/<id>", object metadata in S3 user metadata so a
// put stays atomic (no sidecar second object like the local adapter). The
// AWS SDK v2 client is the MinIO/R2/AWS common denominator; no third
// storage dialect is introduced.
type S3Store struct {
	client s3API
	bucket string
}

// compile-time proof the adapter satisfies the frozen port.
var _ kernel.ObjectStore = (*S3Store)(nil)

// s3API is the frozen API subset the adapter uses (Root D-004). *s3.Client
// satisfies it structurally; tests substitute fakes.
type s3API interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	HeadBucket(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
}

// NewS3 builds the adapter from the frozen storage.objects.s3 config face.
// Credentials are static-only (Root D-005): the SDK default chain (~/.aws,
// IMDS, roles) is never consulted, matching the fail-closed load contract.
// An empty region falls back to us-east-1 (ignored by most compatible backends).
func NewS3(endpoint, region, bucket, accessKeyID, secretAccessKey string, usePathStyle bool) (*S3Store, error) {
	if strings.TrimSpace(bucket) == "" {
		return nil, errors.New("objectstore: s3 bucket must not be empty")
	}
	if region == "" {
		region = "us-east-1"
	}
	awscfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
	}
	client := s3.NewFromConfig(awscfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		o.UsePathStyle = usePathStyle
	})
	return &S3Store{client: client, bucket: bucket}, nil
}

// Ping reports backend reachability for readyz. It is deliberately NOT part
// of the frozen port: composition consumes it via an optional capability
// assertion, keeping the R1 contract untouched.
func (s *S3Store) Ping(ctx context.Context) error {
	if _, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(s.bucket)}); err != nil {
		return fmt.Errorf("objectstore: head bucket %q: %w", s.bucket, err)
	}
	return nil
}

// objectKey validates ns and id, then returns the single-bucket key.
// Validation happens before any network call (fail-closed, same rules as
// the local adapter - shared port semantics).
func objectKey(ns kernel.ObjectNamespace, id string) (string, error) {
	if !kernel.ValidObjectNamespace(ns) {
		return "", fmt.Errorf("objectstore: unknown namespace %q", string(ns))
	}
	if !kernel.ValidObjectID(id) {
		return "", fmt.Errorf("objectstore: invalid object id %q", id)
	}
	return string(ns) + "/" + id, nil
}

// s3Metadata renders user metadata with empty values omitted (mirrors the
// local sidecar writer; empty-valued user metadata is provider-unreliable).
func s3Metadata(meta kernel.ObjectMeta) map[string]string {
	m := map[string]string{}
	for k, v := range map[string]string{
		"name":  meta.Name,
		"type":  meta.Type,
		"kind":  meta.Kind,
		"owner": meta.Owner,
	} {
		if v != "" {
			m[k] = v
		}
	}
	return m
}

// metaFromS3 reads user metadata back into the port record.
func metaFromS3(md map[string]string) kernel.ObjectMeta {
	return kernel.ObjectMeta{Name: md["name"], Type: md["type"], Kind: md["kind"], Owner: md["owner"]}
}

// notFound wraps backend misses in the port sentinel (errors.Is-able),
// recognizing both smithy error codes and bare 404 responses.
func notFound(op, key string, err error) error {
	return fmt.Errorf("objectstore: %s %s: %w", op, key, kernel.ErrObjectNotFound)
}

// mapS3Error maps backend miss signals to the sentinel and passes everything
// else through unchanged.
func mapS3Error(op, key string, err error) error {
	if err == nil {
		return nil
	}
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch ae.ErrorCode() {
		case "NoSuchKey", "NotFound":
			return notFound(op, key, err)
		}
	}
	var re *awshttp.ResponseError
	if errors.As(err, &re) && re.HTTPStatusCode() == 404 {
		return notFound(op, key, err)
	}
	return err
}

func (s *S3Store) Put(ctx context.Context, ns kernel.ObjectNamespace, id string, body []byte, meta kernel.ObjectMeta) error {
	key, err := objectKey(ns, id)
	if err != nil {
		return err
	}
	in := &s3.PutObjectInput{
		Bucket:   aws.String(s.bucket),
		Key:      aws.String(key),
		Body:     bytes.NewReader(body),
		Metadata: s3Metadata(meta),
	}
	if meta.Type != "" {
		in.ContentType = aws.String(meta.Type)
	}
	if _, err := s.client.PutObject(ctx, in); err != nil {
		return fmt.Errorf("objectstore: put %s: %w", key, err)
	}
	return nil
}

func (s *S3Store) Get(ctx context.Context, ns kernel.ObjectNamespace, id string) ([]byte, kernel.ObjectMeta, error) {
	key, err := objectKey(ns, id)
	if err != nil {
		return nil, kernel.ObjectMeta{}, err
	}
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	if err != nil {
		return nil, kernel.ObjectMeta{}, mapS3Error("get", key, err)
	}
	defer out.Body.Close()
	body, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, kernel.ObjectMeta{}, fmt.Errorf("objectstore: read %s: %w", key, err)
	}
	return body, metaFromS3(out.Metadata), nil
}

func (s *S3Store) Stat(ctx context.Context, ns kernel.ObjectNamespace, id string) (kernel.ObjectInfo, error) {
	key, err := objectKey(ns, id)
	if err != nil {
		return kernel.ObjectInfo{}, err
	}
	out, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	if err != nil {
		return kernel.ObjectInfo{}, mapS3Error("stat", key, err)
	}
	info := kernel.ObjectInfo{ID: id, Meta: metaFromS3(out.Metadata)}
	if out.ContentLength != nil {
		info.Size = *out.ContentLength
	}
	if out.LastModified != nil {
		info.ModTime = *out.LastModified
	}
	return info, nil
}

func (s *S3Store) Delete(ctx context.Context, ns kernel.ObjectNamespace, id string) error {
	key, err := objectKey(ns, id)
	if err != nil {
		return err
	}
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)}); err != nil {
		return fmt.Errorf("objectstore: delete %s: %w", key, err)
	}
	return nil
}

func (s *S3Store) Exists(ctx context.Context, ns kernel.ObjectNamespace, id string) (bool, error) {
	key, err := objectKey(ns, id)
	if err != nil {
		return false, err
	}
	if _, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)}); err != nil {
		if mapped := mapS3Error("exists", key, err); errors.Is(mapped, kernel.ErrObjectNotFound) {
			return false, nil
		} else {
			return false, mapped
		}
	}
	return true, nil
}

// List aggregates every page of ListObjectsV2 under the namespace prefix
// and returns ids ascending (S3 lists lexicographically; the explicit sort
// keeps the port deterministic regardless of backend ordering quirks).
func (s *S3Store) List(ctx context.Context, ns kernel.ObjectNamespace) ([]string, error) {
	if !kernel.ValidObjectNamespace(ns) {
		return nil, fmt.Errorf("objectstore: unknown namespace %q", string(ns))
	}
	prefix := string(ns) + "/"
	seen := map[string]bool{}
	ids := []string{}
	in := &s3.ListObjectsV2Input{Bucket: aws.String(s.bucket), Prefix: aws.String(prefix)}
	for {
		out, err := s.client.ListObjectsV2(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("objectstore: list %s: %w", string(ns), err)
		}
		for _, obj := range out.Contents {
			if obj.Key == nil {
				continue
			}
			id := strings.TrimPrefix(*obj.Key, prefix)
			if !kernel.ValidObjectID(id) || seen[id] {
				continue
			}
			seen[id] = true
			ids = append(ids, id)
		}
		if out.IsTruncated != nil && *out.IsTruncated && out.NextContinuationToken != nil {
			in.ContinuationToken = out.NextContinuationToken
			continue
		}
		break
	}
	sort.Strings(ids)
	return ids, nil
}
