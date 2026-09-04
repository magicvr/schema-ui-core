package store

import (
	"path/filepath"
	"testing"
)

func TestMigrateV66TelegramRowPreservesExistingConfigOnV67(t *testing.T) {
	path := filepath.Join(t.TempDir(), "telegram-v66.db")
	db := rawOpen(t, path)
	preUpgrade := &Store{db: db, path: path}
	for _, migration := range compiledMigrations[:66] {
		if err := preUpgrade.applyMigration(migration); err != nil {
			_ = db.Close()
			t.Fatalf("apply v%d (%s): %v", migration.Version, migration.Name, err)
		}
	}
	if _, err := db.Exec(`INSERT INTO telegram_config (id, bot_token_enc, webhook_secret_enc, updated_at) VALUES (1, ?, ?, ?)
`, "legacy-token-ciphertext", "legacy-secret-ciphertext", 123); err != nil {
		_ = db.Close()
		t.Fatalf("insert pre-v67 telegram row: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	upgraded, err := OpenWithCatalog(path, compiledMigrations)
	if err != nil {
		t.Fatalf("upgrade v66 to v67: %v", err)
	}
	defer upgraded.Close()

	var tokenCiphertext, secretCiphertext, mode, origin string
	var updatedAt int64
	if err := upgraded.db.QueryRow(`SELECT bot_token_enc, webhook_secret_enc, mode, webhook_public_base_url, updated_at FROM telegram_config WHERE id = 1`).Scan(
		&tokenCiphertext, &secretCiphertext, &mode, &origin, &updatedAt,
	); err != nil {
		t.Fatalf("read upgraded telegram row: %v", err)
	}
	if tokenCiphertext != "legacy-token-ciphertext" || secretCiphertext != "legacy-secret-ciphertext" ||
		mode != "polling" || origin != "" || updatedAt != 123 {
		t.Fatalf("upgraded telegram row = token %q secret %q mode %q origin %q updated_at %d", tokenCiphertext, secretCiphertext, mode, origin, updatedAt)
	}
	applied, err := upgraded.appliedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) < 68 || applied[65].version != 66 || applied[66].version != 67 || applied[66].name != "telegram_config_connection" || applied[67].version != 68 || applied[67].name != "telegram_ingress" {
		tail := applied
		if len(tail) > 3 {
			tail = tail[len(tail)-3:]
		}
		t.Fatalf("upgraded migration tail = %+v", tail)
	}
}
