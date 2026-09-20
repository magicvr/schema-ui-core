package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalcontract"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// requireDenominatorColumn asserts a column really belongs to the frozen VP-040
// denominator with the claimed legacy unit. Reading the unit back from
// internal/temporalcontract (the frozen 90-column denominator) makes each
// sub-test's "this is the seconds/milliseconds family" claim checkable instead
// of asserted by prose.
func requireDenominatorColumn(t *testing.T, table, column, wantUnit string) {
	t.Helper()
	for _, c := range temporalcontract.Columns() {
		if c.Table == table && c.Column == column {
			if c.Unit != wantUnit {
				t.Fatalf("denominator says %s.%s has unit %q, want %q", table, column, c.Unit, wantUnit)
			}
			return
		}
	}
	t.Fatalf("%s.%s is not in the frozen VP-040 denominator", table, column)
}

// TestR3BUnitFamilyWireMatrix is the R3-B unit-family matrix: at least one
// endpoint per persisted unit family of the frozen denominator, asserting the
// wire value of the converted column (Root D-018 §4 R3-B).
//
//	family      columns                                endpoint
//	seconds     digital_offers.created_at/updated_at    GET /api/digitaloffer/offers
//	milliseconds jobs.created_at/updated_at/finished_at  GET /api/jobs
//	nullable    service_credentials.revoked_at/last_used_at  GET /api/service-credentials/{id}
//	sentinel    mail_config.updated_at (legacy 0 → NULL) GET /api/mail/config
//
// Every sub-test pins a trailing-zero microsecond instant straight into the
// column, so a width-losing or rounding path cannot pass by accident.
func TestR3BUnitFamilyWireMatrix(t *testing.T) {
	t.Run("seconds family: digital_offers via the admin offer list", func(t *testing.T) {
		requireDenominatorColumn(t, "digital_offers", "created_at", "sec")
		requireDenominatorColumn(t, "digital_offers", "updated_at", "sec")

		env, _, _ := newDigitalOfferEnv(t)
		token := env.login(t, testSeedUsername, testSeedPassword)

		created := serveDigitalOffer(t, env, token, http.MethodPost, "/api/digitaloffer/offers",
			`{"name":"R3B seconds","priceAmount":100,"currency":"CNY","entitlementForm":"count","countPerPurchase":1}`)
		if created.Code != http.StatusCreated {
			t.Fatalf("create offer = %d %s", created.Code, created.Body.String())
		}
		offerID := between(created.Body.String(), `"id":"`, `"`)
		if offerID == "" {
			t.Fatalf("create offer returned no id: %s", created.Body.String())
		}

		// Pin both seconds-family instants so the assertion is decisive.
		if err := env.st.Run(context.Background(), func(tx kernel.Tx) error {
			_, err := tx.Exec(context.Background(),
				`UPDATE digital_offers SET created_at = ?, updated_at = ? WHERE id = ?`,
				trailingZeroInstant, wholeSecondInstant, offerID)
			return err
		}); err != nil {
			t.Fatalf("pin digital_offers instants: %v", err)
		}

		list := serveDigitalOffer(t, env, token, http.MethodGet, "/api/digitaloffer/offers?pageSize=50", "")
		if list.Code != http.StatusOK {
			t.Fatalf("list offers = %d %s", list.Code, list.Body.String())
		}
		var envelope struct {
			Items []map[string]any `json:"items"`
		}
		if err := json.Unmarshal(list.Body.Bytes(), &envelope); err != nil {
			t.Fatalf("decode offers: %v", err)
		}
		for _, row := range envelope.Items {
			if row["id"] != offerID {
				continue
			}
			if row["createdAt"] != trailingZeroInstant {
				t.Fatalf("createdAt = %v, want %q", row["createdAt"], trailingZeroInstant)
			}
			if row["updatedAt"] != wholeSecondInstant {
				t.Fatalf("updatedAt = %v, want %q", row["updatedAt"], wholeSecondInstant)
			}
			return
		}
		t.Fatalf("offer %s missing from the admin list: %s", offerID, list.Body.String())
	})

	t.Run("milliseconds family: jobs via the management list", func(t *testing.T) {
		requireDenominatorColumn(t, "jobs", "created_at", "ms")
		requireDenominatorColumn(t, "jobs", "updated_at", "ms")
		requireDenominatorColumn(t, "jobs", "finished_at", "ms")

		env := newAuthTestEnv(t)
		repository := newJobsTestRepository(t, env)
		mountJobsRoutes(t, env, repository)

		at, err := time.Parse("2006-01-02T15:04:05.000000Z", trailingZeroInstant)
		if err != nil {
			t.Fatalf("parse seed instant: %v", err)
		}
		seedJob(t, repository, "job-r3b-ms", "wallet.reconcile", "actor-r3b", at)

		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, adminToken(t, env), http.MethodGet, "/api/jobs", ""))
		if rr.Code != http.StatusOK {
			t.Fatalf("GET /api/jobs = %d %s", rr.Code, rr.Body.String())
		}
		var envelope struct {
			Items []map[string]any `json:"items"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &envelope); err != nil {
			t.Fatalf("decode jobs: %v", err)
		}
		if len(envelope.Items) != 1 {
			t.Fatalf("jobs items = %v, want the seeded row", envelope.Items)
		}
		row := envelope.Items[0]
		if row["createdAt"] != trailingZeroInstant || row["updatedAt"] != trailingZeroInstant {
			t.Fatalf("job instants = createdAt %v / updatedAt %v, want %q",
				row["createdAt"], row["updatedAt"], trailingZeroInstant)
		}
		// jobs.finished_at is a nullable milliseconds column: a queued job has no
		// instant, so the field must be absent rather than fabricated.
		if got, present := row["finishedAt"]; present {
			t.Fatalf("queued job carries finishedAt=%v, want the field absent for a NULL column", got)
		}
	})

	t.Run("nullable family: service_credentials via the metadata detail", func(t *testing.T) {
		requireDenominatorColumn(t, "service_credentials", "revoked_at", "sec")
		requireDenominatorColumn(t, "service_credentials", "last_used_at", "sec")

		env := newAuthTestEnv(t)
		token := adminToken(t, env)
		expiresAt := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02T15:04:05.000000Z")
		created := httptest.NewRecorder()
		env.mux.ServeHTTP(created, bearer(t, token, http.MethodPost, "/api/service-credentials",
			`{"name":"R3B nullable","scopes":["users.read"],"expiresAt":"`+expiresAt+`"}`))
		if created.Code != http.StatusCreated {
			t.Fatalf("create credential = %d %s", created.Code, created.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(created.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode create: %v", err)
		}
		id, _ := body["id"].(string)
		if id == "" {
			t.Fatalf("create credential returned no id: %s", created.Body.String())
		}

		detail := func() map[string]json.RawMessage {
			t.Helper()
			rr := httptest.NewRecorder()
			env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/service-credentials/"+id, ""))
			if rr.Code != http.StatusOK {
				t.Fatalf("GET credential = %d %s", rr.Code, rr.Body.String())
			}
			return decodeObject(t, rr.Body.Bytes())
		}

		// A fresh credential has both nullable columns NULL: the wire must say
		// null, never a fabricated epoch instant.
		nulls := detail()
		for _, field := range []string{"revokedAt", "lastUsedAt"} {
			if got := string(nulls[field]); got != "null" {
				t.Fatalf("%s = %s on a fresh credential, want JSON null", field, got)
			}
		}

		// The nullable columns are genuinely NULL in storage (so the JSON above
		// is a projection of absence, not of an epoch value).
		if err := env.st.Run(context.Background(), func(tx kernel.Tx) error {
			var revokedAt, lastUsedAt *string
			if err := tx.QueryRow(context.Background(),
				`SELECT revoked_at, last_used_at FROM service_credentials WHERE id = ?`, id).
				Scan(&revokedAt, &lastUsedAt); err != nil {
				return err
			}
			if revokedAt != nil || lastUsedAt != nil {
				t.Fatalf("stored nullable columns = (%v, %v), want SQL NULL", revokedAt, lastUsedAt)
			}
			return nil
		}); err != nil {
			t.Fatalf("read stored nullable columns: %v", err)
		}

		// Setting them must surface the canonical fixed-6 form.
		if err := env.st.Run(context.Background(), func(tx kernel.Tx) error {
			_, err := tx.Exec(context.Background(),
				`UPDATE service_credentials SET revoked_at = ?, last_used_at = ? WHERE id = ?`,
				trailingZeroInstant, wholeSecondInstant, id)
			return err
		}); err != nil {
			t.Fatalf("set nullable columns: %v", err)
		}
		set := detail()
		var revokedAt, lastUsedAt string
		if err := json.Unmarshal(set["revokedAt"], &revokedAt); err != nil {
			t.Fatalf("revokedAt is not a string: %s", set["revokedAt"])
		}
		if err := json.Unmarshal(set["lastUsedAt"], &lastUsedAt); err != nil {
			t.Fatalf("lastUsedAt is not a string: %s", set["lastUsedAt"])
		}
		if revokedAt != trailingZeroInstant || lastUsedAt != wholeSecondInstant {
			t.Fatalf("nullable instants = (%q, %q), want (%q, %q)",
				revokedAt, lastUsedAt, trailingZeroInstant, wholeSecondInstant)
		}
	})

	t.Run("sentinel-converted family: mail_config via the config read", func(t *testing.T) {
		requireDenominatorColumn(t, "mail_config", "updated_at", "ms")

		env := newAuthTestEnv(t)
		key := []byte(strings.Repeat("m", 32))
		sw, err := mail.NewSwitcher(env.st, key, mail.SeedConfig{Channel: mail.RuntimeChannelMock}, nil)
		if err != nil {
			t.Fatalf("NewSwitcher: %v", err)
		}
		RegisterMailAdmin(env.mux, env.a, sw, env.operations)
		token := adminToken(t, env)

		config := func() map[string]json.RawMessage {
			t.Helper()
			rr := httptest.NewRecorder()
			env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/mail/config", ""))
			if rr.Code != http.StatusOK {
				t.Fatalf("GET /api/mail/config = %d %s", rr.Code, rr.Body.String())
			}
			return decodeObject(t, rr.Body.Bytes())
		}
		setUpdatedAt := func(sql string, args ...any) {
			t.Helper()
			if err := env.st.Run(context.Background(), func(tx kernel.Tx) error {
				_, err := tx.Exec(context.Background(), sql, args...)
				return err
			}); err != nil {
				t.Fatalf("set mail_config.updated_at: %v", err)
			}
		}

		// The legacy INTEGER 0 sentinel became SQL NULL in v73. Storing NULL is
		// therefore the post-conversion state the wire must reflect as JSON null
		// — never as 1970/zero.
		setUpdatedAt(`UPDATE mail_config SET updated_at = NULL WHERE id = 1`)
		if err := env.st.Run(context.Background(), func(tx kernel.Tx) error {
			var stored *string
			if err := tx.QueryRow(context.Background(),
				`SELECT updated_at FROM mail_config WHERE id = 1`).Scan(&stored); err != nil {
				return err
			}
			if stored != nil {
				t.Fatalf("stored mail_config.updated_at = %q, want SQL NULL (the legacy 0 sentinel is gone)", *stored)
			}
			return nil
		}); err != nil {
			t.Fatalf("read stored mail_config.updated_at: %v", err)
		}
		if got := string(config()["updated_at"]); got != "null" {
			t.Fatalf("unset sentinel column = %s, want JSON null", got)
		}

		// Once configured, the same field is canonical fixed-6.
		setUpdatedAt(`UPDATE mail_config SET updated_at = ? WHERE id = 1`, trailingZeroInstant)
		var got string
		if err := json.Unmarshal(config()["updated_at"], &got); err != nil {
			t.Fatalf("updated_at is not a string: %s", config()["updated_at"])
		}
		if got != trailingZeroInstant {
			t.Fatalf("configured updated_at = %q, want %q", got, trailingZeroInstant)
		}
	})
}
