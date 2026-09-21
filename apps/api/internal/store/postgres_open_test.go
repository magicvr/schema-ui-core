package store

import (
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestMissingDatabaseHintIsActionableAndPreservesClassification(t *testing.T) {
	missing := &pgconn.PgError{Code: "3D000", Message: `database "schema_ui_dev" does not exist`}
	err := missingDatabaseHint(missing)
	if !errors.Is(err, missing) {
		t.Fatalf("hint error does not preserve the original pg error: %v", err)
	}
	for _, want := range []string{"schema_ui_dev", "dev.cmd init-db", "go run ./cmd/dbsetup", "SQLSTATE"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("hint error %q does not contain %q", err, want)
		}
	}

	other := &pgconn.PgError{Code: "28P01", Message: "password authentication failed"}
	if got := missingDatabaseHint(other); got != other {
		t.Fatalf("non-3D000 error was rewritten: got %v, want original %v", got, other)
	}
}
