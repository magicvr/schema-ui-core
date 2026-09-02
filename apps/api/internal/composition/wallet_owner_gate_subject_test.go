// VP-029 A-005 F-001 regression (workspace-029 · GOAL-001 · A-006 response):
// the composition OwnerExistsFunc behind the by-owner auto-create surface must
// check the live USER table ONLY. POST /api/wallet/by-owner/{id} and
// .../adjust must reject a registered external subject id (404 USER_NOT_FOUND)
// and never insert an owner_type=user wallet_accounts row — subject existence
// stays inside CreateAccount(owner_type=subject) / voucher Redeem.
package composition

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
)

func TestWalletByOwnerGateRejectsRegisteredSubject(t *testing.T) {
	hash, err := auth.HashPassword("pw", 4)
	if err != nil {
		t.Fatal(err)
	}
	// No seed: testsupport pre-reconciles with its own (partial) contribution
	// list, which would checksum-drift against the composition's production
	// contributions inside testMux. Replicate the production boot order
	// instead: open store → bootstrap admin → let newMux reconcile the real
	// contribution set once.
	st, err := testsupport.OpenStore(":memory:", "admin", hash, false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if err := authsessiondata.Bootstrap(context.Background(), st, "admin", hash); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := st.WithTx(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec(`UPDATE users SET must_change_password = 0`)
		return err
	}); err != nil {
		t.Fatal(err)
	}

	plan, err := ResolvePlan(&config.Config{ProfileName: "admin"})
	if err != nil {
		t.Fatal(err)
	}
	secret := []byte("test-secret")
	a := auth.New(secret, 0, 0, st, false)
	mux, err := testMux(a, st, plan, &readinessGate{})
	if err != nil {
		t.Fatal(err)
	}

	repo := authsession.NewRepository(st)
	admin, err := repo.UserByID("user-admin")
	if err != nil {
		t.Fatalf("seeded admin: %v", err)
	}
	token, err := auth.SignAccessToken(secret, admin.ID, admin.TokenVersion, "session-a005", 15*time.Minute, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}

	// A real registered external subject on the same store the composed
	// wallet service uses (composition wires walletService over this store).
	now := time.Now().UTC()
	sub, _, err := subject.NewStore(st).GetOrCreateSubject(context.Background(), "issuer-a005", "ext-1", now)
	if err != nil {
		t.Fatal(err)
	}
	escaped := url.PathEscape(sub.ID)

	userRows := func(ownerID string) int {
		return compositionCount(t, st, `SELECT COUNT(*) FROM wallet_accounts WHERE owner_type = ? AND owner_id = ?`, "user", ownerID)
	}
	do := func(method, path, body string) *httptest.ResponseRecorder {
		t.Helper()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		mux.ServeHTTP(rr, req)
		return rr
	}

	// POST by-owner create: a subject id must fail the user-only gate.
	rr := do(http.MethodPost, "/api/wallet/by-owner/"+escaped, "")
	if rr.Code != http.StatusNotFound || !strings.Contains(rr.Body.String(), "USER_NOT_FOUND") {
		t.Fatalf("by-owner create with subject id = %d %s, want 404 USER_NOT_FOUND", rr.Code, rr.Body.String())
	}
	if rows := userRows(sub.ID); rows != 0 {
		t.Fatalf("subject id minted %d owner_type=user account(s), want 0", rows)
	}

	// POST by-owner adjust: the money path shares the same user-only gate.
	rr = do(http.MethodPost, "/api/wallet/by-owner/"+escaped+"/adjust", `{"amountDelta":100,"memo":"must not land"}`)
	if rr.Code != http.StatusNotFound || !strings.Contains(rr.Body.String(), "USER_NOT_FOUND") {
		t.Fatalf("by-owner adjust with subject id = %d %s, want 404 USER_NOT_FOUND", rr.Code, rr.Body.String())
	}
	if rows := userRows(sub.ID); rows != 0 {
		t.Fatalf("subject id adjust minted %d owner_type=user account(s), want 0", rows)
	}

	// Positive control: a real user id still auto-opens its user account.
	rr = do(http.MethodPost, "/api/wallet/by-owner/"+admin.ID, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("by-owner create with real user = %d %s, want 200", rr.Code, rr.Body.String())
	}
	if rows := userRows(admin.ID); rows != 1 {
		t.Fatalf("real user account rows = %d, want 1", rows)
	}
}
