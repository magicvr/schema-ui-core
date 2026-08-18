package authsession

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"sync"
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

func TestServiceCredentialConcurrentDuplicateName(t *testing.T) {
	repository, _ := openRepository(t, "service-credential-concurrent.db", true)
	now := time.Now().UTC()
	credentials := []ServiceCredential{
		{ID: "11111111111111111111111111111111", Name: "Deploy Agent", TokenPrefix: "sui_sc_11111111", TokenHash: "1111111111111111111111111111111111111111111111111111111111111111"},
		{ID: "22222222222222222222222222222222", Name: "deploy agent", TokenPrefix: "sui_sc_22222222", TokenHash: "2222222222222222222222222222222222222222222222222222222222222222"},
	}
	for index := range credentials {
		credentials[index].Scopes = []string{"users.read"}
		credentials[index].ExpiresAt = now.Add(time.Hour)
		credentials[index].CreatedBy = "user-admin"
		credentials[index].CreatedAt = now
		credentials[index].UpdatedAt = now
	}
	start := make(chan struct{})
	errorsByIndex := make([]error, len(credentials))
	var wait sync.WaitGroup
	for index := range credentials {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			<-start
			errorsByIndex[index] = repository.CreateServiceCredential(credentials[index], nil)
		}(index)
	}
	close(start)
	wait.Wait()
	successes, duplicates := 0, 0
	for _, err := range errorsByIndex {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrCredentialNameTaken):
			duplicates++
		default:
			t.Fatalf("unexpected concurrent create error: %v", err)
		}
	}
	if successes != 1 || duplicates != 1 {
		t.Fatalf("concurrent results successes=%d duplicates=%d errors=%v", successes, duplicates, errorsByIndex)
	}
}

func TestServiceCredentialConcurrentRevokeTransitionsOnce(t *testing.T) {
	repository, _ := openRepository(t, "service-credential-concurrent-revoke.db", true)
	now := time.Now().UTC()
	credential := ServiceCredential{
		ID: "33333333333333333333333333333333", Name: "Revoke Agent",
		TokenPrefix: "sui_sc_33333333", TokenHash: "3333333333333333333333333333333333333333333333333333333333333333",
		Scopes: []string{"users.read"}, ExpiresAt: now.Add(time.Hour), CreatedBy: "user-admin",
		CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.CreateServiceCredential(credential, nil); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make([]bool, 2)
	errorsByIndex := make([]error, 2)
	audits := 0
	var auditMu sync.Mutex
	var wait sync.WaitGroup
	for index := range results {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			<-start
			_, results[index], errorsByIndex[index] = repository.RevokeServiceCredential(credential.ID, now.Add(time.Minute), func(*sql.Tx, ServiceCredential) error {
				auditMu.Lock()
				audits++
				auditMu.Unlock()
				return nil
			})
		}(index)
	}
	close(start)
	wait.Wait()
	changes := 0
	for index, err := range errorsByIndex {
		if err != nil {
			t.Fatalf("revoke %d: %v", index, err)
		}
		if results[index] {
			changes++
		}
	}
	if changes != 1 || audits != 1 {
		t.Fatalf("concurrent revoke changes=%d audits=%d results=%v", changes, audits, results)
	}
}
