package migration

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Entry is the framework-independent migration metadata consumed by the
// composition root and migration audit tooling.
type Entry struct {
	Version  int
	Name     string
	ModuleID string
	Checksum string
}

type Plan struct {
	Entries []Entry
}

type ErrorCode string

const (
	CodeInvalid          ErrorCode = "MIGRATION_INVALID"
	CodeDuplicateVersion ErrorCode = "MIGRATION_DUPLICATE_VERSION"
	CodeDuplicateName    ErrorCode = "MIGRATION_DUPLICATE_NAME"
	CodeOutOfOrder       ErrorCode = "MIGRATION_OUT_OF_ORDER"
)

type Error struct {
	Code    ErrorCode
	Version int
	Name    string
	Detail  string
}

func (e *Error) Error() string {
	identity := e.Name
	if e.Version > 0 {
		identity = fmt.Sprintf("%04d %s", e.Version, e.Name)
	}
	if identity == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Detail)
	}
	return fmt.Sprintf("%s [%s]: %s", e.Code, identity, e.Detail)
}

func Collect(entries []Entry) (Plan, error) {
	if len(entries) == 0 {
		return Plan{}, nil
	}

	seenVersions := make(map[int]Entry, len(entries))
	seenNames := make(map[string]Entry, len(entries))
	collected := make([]Entry, 0, len(entries))
	for _, raw := range entries {
		entry := raw
		entry.Name = strings.TrimSpace(entry.Name)
		entry.ModuleID = strings.TrimSpace(entry.ModuleID)
		entry.Checksum = strings.ToLower(strings.TrimSpace(entry.Checksum))
		if entry.Version <= 0 || entry.Name == "" || entry.ModuleID == "" {
			return Plan{}, &Error{Code: CodeInvalid, Version: entry.Version, Name: entry.Name, Detail: "version, name and module id are required"}
		}
		if len(entry.Checksum) != 64 {
			return Plan{}, &Error{Code: CodeInvalid, Version: entry.Version, Name: entry.Name, Detail: "checksum must be 64 hexadecimal characters"}
		}
		if _, err := hex.DecodeString(entry.Checksum); err != nil {
			return Plan{}, &Error{Code: CodeInvalid, Version: entry.Version, Name: entry.Name, Detail: "checksum must be hexadecimal"}
		}
		if previous, exists := seenVersions[entry.Version]; exists {
			return Plan{}, &Error{Code: CodeDuplicateVersion, Version: entry.Version, Name: entry.Name, Detail: fmt.Sprintf("version already belongs to %s", previous.Name)}
		}
		if previous, exists := seenNames[entry.Name]; exists {
			return Plan{}, &Error{Code: CodeDuplicateName, Version: entry.Version, Name: entry.Name, Detail: fmt.Sprintf("name already belongs to version %d", previous.Version)}
		}
		seenVersions[entry.Version] = entry
		seenNames[entry.Name] = entry
		collected = append(collected, entry)
	}

	sort.Slice(collected, func(i, j int) bool { return collected[i].Version < collected[j].Version })
	for i, entry := range collected {
		if i == 0 {
			if entry.Version != 1 {
				return Plan{}, &Error{Code: CodeOutOfOrder, Version: entry.Version, Name: entry.Name, Detail: "migration plan must start at version 1"}
			}
			continue
		}
		if entry.Version != collected[i-1].Version+1 {
			return Plan{}, &Error{Code: CodeOutOfOrder, Version: entry.Version, Name: entry.Name, Detail: fmt.Sprintf("expected version %d after %d", collected[i-1].Version+1, collected[i-1].Version)}
		}
	}

	return Plan{Entries: collected}, nil
}

func (p Plan) Versions() []int {
	versions := make([]int, 0, len(p.Entries))
	for _, entry := range p.Entries {
		versions = append(versions, entry.Version)
	}
	return versions
}
