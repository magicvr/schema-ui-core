package concurrency

import (
	"errors"
	"testing"
)

func ptr(value int64) *int64 { return &value }

func TestETag(t *testing.T) {
	if got := ETag(12); got != `"v12"` {
		t.Fatalf("ETag(12) = %q", got)
	}
}

func TestResolveExpectedVersion(t *testing.T) {
	tests := []struct {
		name     string
		header   []string
		expected *int64
		legacy   *int64
		want     int64
		wantErr  error
	}{
		{name: "header", header: []string{` "v2" `}, want: 2},
		{name: "expected zero", expected: ptr(0), want: 0},
		{name: "legacy", legacy: ptr(3), want: 3},
		{name: "all agree", header: []string{`"v4"`}, expected: ptr(4), legacy: ptr(4), want: 4},
		{name: "missing", wantErr: ErrPreconditionRequired},
		{name: "mismatch", header: []string{`"v1"`}, expected: ptr(2), wantErr: ErrInvalidPrecondition},
		{name: "negative", expected: ptr(-1), wantErr: ErrInvalidPrecondition},
		{name: "weak", header: []string{`W/"v1"`}, wantErr: ErrInvalidPrecondition},
		{name: "wildcard", header: []string{"*"}, wantErr: ErrInvalidPrecondition},
		{name: "list", header: []string{`"v1", "v2"`}, wantErr: ErrInvalidPrecondition},
		{name: "unquoted", header: []string{"v1"}, wantErr: ErrInvalidPrecondition},
		{name: "internal whitespace", header: []string{`"v 1"`}, wantErr: ErrInvalidPrecondition},
		{name: "multiple headers", header: []string{`"v1"`, `"v1"`}, wantErr: ErrInvalidPrecondition},
		{name: "overflow", header: []string{`"v9223372036854775808"`}, wantErr: ErrInvalidPrecondition},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ResolveExpectedVersion(tc.header, tc.expected, tc.legacy)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}
			if err == nil && got != tc.want {
				t.Fatalf("version = %d, want %d", got, tc.want)
			}
		})
	}
}
