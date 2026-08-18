// Package concurrency defines the shared optimistic-concurrency wire contract.
package concurrency

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrPreconditionRequired = errors.New("concurrency: precondition required")
	ErrInvalidPrecondition  = errors.New("concurrency: invalid precondition")
)

// ETag returns the strong entity tag for a non-negative resource version.
func ETag(version int64) string {
	return fmt.Sprintf(`"v%d"`, version)
}

// ResolveExpectedVersion resolves If-Match, expectedVersion, and the legacy
// version field. Every provided source must name the same non-negative value.
func ResolveExpectedVersion(ifMatch []string, expectedVersion, legacyVersion *int64) (int64, error) {
	versions := make([]int64, 0, 3)
	if len(ifMatch) > 0 {
		if len(ifMatch) != 1 {
			return 0, ErrInvalidPrecondition
		}
		version, err := parseETag(ifMatch[0])
		if err != nil {
			return 0, err
		}
		versions = append(versions, version)
	}
	for _, version := range []*int64{expectedVersion, legacyVersion} {
		if version == nil {
			continue
		}
		if *version < 0 {
			return 0, ErrInvalidPrecondition
		}
		versions = append(versions, *version)
	}
	if len(versions) == 0 {
		return 0, ErrPreconditionRequired
	}
	for _, version := range versions[1:] {
		if version != versions[0] {
			return 0, ErrInvalidPrecondition
		}
	}
	return versions[0], nil
}

func parseETag(raw string) (int64, error) {
	tag := strings.TrimSpace(raw)
	if len(tag) < 4 || tag[0] != '"' || tag[len(tag)-1] != '"' || !strings.HasPrefix(tag[1:], "v") {
		return 0, ErrInvalidPrecondition
	}
	digits := tag[2 : len(tag)-1]
	if digits == "" || strings.ContainsAny(digits, " \t\r\n,+-") {
		return 0, ErrInvalidPrecondition
	}
	version, err := strconv.ParseInt(digits, 10, 64)
	if err != nil || version < 0 {
		return 0, ErrInvalidPrecondition
	}
	return version, nil
}
