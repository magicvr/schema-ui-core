package authsession

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"testing"
	"time"
)

func TestServiceCredentialPersistenceContract(t *testing.T) {
	repository, st := openRepository(t, "service-credentials.db", true)
	now := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)
	raw := "sui_sc_abcdefghijklmnopqrstuvwxyzABCDEFGH123456789"
	sum := sha256.Sum256([]byte(raw))
	credential := ServiceCredential{
		ID: "0123456789abcdef0123456789abcdef", Name: "Build Agent",
		TokenPrefix: raw[:15], TokenHash: hex.EncodeToString(sum[:]),
		Scopes:    []string{"files.write", "records.read", "files.write"},
		ExpiresAt: now.Add(24 * time.Hour), CreatedBy: "user-admin",
		CreatedAt: now, UpdatedAt: now,
	}
	audited := false
	if err := repository.CreateServiceCredential(credential, func(tx *sql.Tx) error {
		audited = true
		_, err := tx.Exec(`INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, created_at)
VALUES ('op-credential-create', 'service-credentials.create', 'user-admin', 'Admin', ?, ?)`, credential.ID, now.UnixMilli())
		return err
	}); err != nil {
		t.Fatalf("CreateServiceCredential: %v", err)
	}
	if !audited {
		t.Fatal("transactional audit callback was not called")
	}

	got, err := repository.ServiceCredentialByHash(credential.TokenHash)
	if err != nil {
		t.Fatalf("ServiceCredentialByHash: %v", err)
	}
	if got.Name != credential.Name || got.TokenPrefix != raw[:15] || len(got.Scopes) != 2 || got.Scopes[0] != "files.write" || got.Scopes[1] != "records.read" {
		t.Fatalf("credential = %+v", got)
	}
	var rawMatches int
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		return tx.QueryRow(`SELECT COUNT(*) FROM service_credentials WHERE token_hash = ? OR token_prefix = ?`, raw, raw).Scan(&rawMatches)
	}); err != nil {
		t.Fatal(err)
	}
	if rawMatches != 0 {
		t.Fatal("raw credential was persisted")
	}

	duplicate := credential
	duplicate.ID = "fedcba9876543210fedcba9876543210"
	duplicate.Name = "build agent"
	duplicate.TokenHash = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := repository.CreateServiceCredential(duplicate, nil); !errors.Is(err, ErrCredentialNameTaken) {
		t.Fatalf("case-insensitive duplicate = %v, want ErrCredentialNameTaken", err)
	}

	if err := repository.MarkServiceCredentialUsed(credential.ID, now.Add(time.Hour)); err != nil {
		t.Fatalf("MarkServiceCredentialUsed: %v", err)
	}
	got, err = repository.ServiceCredentialByID(credential.ID)
	if err != nil || got.LastUsedAt == nil {
		t.Fatalf("used credential = %+v, err=%v", got, err)
	}
	revoked, changed, err := repository.RevokeServiceCredential(credential.ID, now.Add(2*time.Hour), nil)
	if err != nil || !changed || revoked.RevokedAt == nil {
		t.Fatalf("revoke = %+v changed=%v err=%v", revoked, changed, err)
	}
	_, changed, err = repository.RevokeServiceCredential(credential.ID, now.Add(3*time.Hour), nil)
	if err != nil || changed {
		t.Fatalf("idempotent revoke changed=%v err=%v", changed, err)
	}
}

func TestServiceCredentialAuditFailureRollsBackMutation(t *testing.T) {
	repository, st := openRepository(t, "service-credential-rollback.db", true)
	now := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)
	forced := errors.New("forced audit failure")
	err := repository.CreateServiceCredential(ServiceCredential{
		ID: "0123456789abcdef0123456789abcdef", Name: "Rollback",
		TokenPrefix: "sui_sc_abcdefgh", TokenHash: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Scopes: []string{"records.read"}, ExpiresAt: now.Add(time.Hour), CreatedBy: "user-admin",
		CreatedAt: now, UpdatedAt: now,
	}, func(*sql.Tx) error { return forced })
	if !errors.Is(err, forced) {
		t.Fatalf("CreateServiceCredential error = %v, want forced audit failure", err)
	}
	if count := repositoryQueryInt(t, st, `SELECT COUNT(*) FROM service_credentials`); count != 0 {
		t.Fatalf("credentials after rollback = %d, want 0", count)
	}
}
