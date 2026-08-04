package migration

import (
	"errors"
	"reflect"
	"testing"
)

const testChecksum = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestCollectSortsAndValidatesContiguousPlan(t *testing.T) {
	plan, err := Collect([]Entry{
		{Version: 2, Name: "second", ModuleID: "module.second", Checksum: testChecksum},
		{Version: 1, Name: "first", ModuleID: "module.first", Checksum: testChecksum},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := plan.Versions(), []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("versions = %v, want %v", got, want)
	}
}

func TestCollectRejectsInvalidIdentityChecksumDuplicatesAndGaps(t *testing.T) {
	cases := []struct {
		name string
		code ErrorCode
		data []Entry
	}{
		{name: "invalid checksum", code: CodeInvalid, data: []Entry{{Version: 1, Name: "one", ModuleID: "one", Checksum: "bad"}}},
		{name: "duplicate version", code: CodeDuplicateVersion, data: []Entry{{Version: 1, Name: "one", ModuleID: "one", Checksum: testChecksum}, {Version: 1, Name: "other", ModuleID: "other", Checksum: testChecksum}}},
		{name: "duplicate name", code: CodeDuplicateName, data: []Entry{{Version: 1, Name: "one", ModuleID: "one", Checksum: testChecksum}, {Version: 2, Name: "one", ModuleID: "two", Checksum: testChecksum}}},
		{name: "gap", code: CodeOutOfOrder, data: []Entry{{Version: 1, Name: "one", ModuleID: "one", Checksum: testChecksum}, {Version: 3, Name: "three", ModuleID: "three", Checksum: testChecksum}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Collect(tc.data)
			var migrationErr *Error
			if !errors.As(err, &migrationErr) || migrationErr.Code != tc.code {
				t.Fatalf("error = %v, want code %s", err, tc.code)
			}
		})
	}
}
