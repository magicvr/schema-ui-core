package handler

// Error code contract (VP-007 S4 · C1): the stable machine-readable code set
// is pinned here and may not drift. The frozen enumeration lives in Root
// D-002 appendix A (31 literals + 8 domain families); S3 added the three
// settings validation codes; the auth middleware codes are included.

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
)

var frozenLiteralCodes = []string{
	"EMPTY_SELECTION", "FILE_NOT_FOUND", "FILE_TOO_LARGE", "FORBIDDEN", "INTERNAL",
	"INVALID_BODY", "INVALID_CREATE_BODY", "INVALID_CREATE_FIELD", "INVALID_FILE",
	"INVALID_LOGIN_BODY", "INVALID_LOGO_URL", "INVALID_LOGOUT_BODY", "INVALID_PAGE",
	"INVALID_PAGE_SIZE", "INVALID_DATE_FILTER", "INVALID_PATCH_BODY", "INVALID_PATCH_FIELD", "INVALID_REFRESH_BODY",
	"INVALID_SELECTION_KEY", "INVALID_SITE_TITLE", "INVALID_SORT_FIELD", "INVALID_SORT_ORDER",
	"INVALID_UPLOAD", "LOGIN_FAILED", "LOGOUT_FAILED", "REFRESH_FAILED", "SCHEMA_NOT_FOUND",
	"SETTINGS_NOT_FOUND", "STORAGE_UNAVAILABLE", "UNAUTHENTICATED", "UNAUTHORIZED",
	"REFRESH_TOKEN_EXPIRED", "NOT_FOUND", "METHOD_NOT_ALLOWED",
	"UNSUPPORTED_FILE_TYPE",
	// W4 P0-2 per-user upload quota.
	"UPLOAD_QUOTA_EXCEEDED",
	// S3 settings validation additions (D-002 appendix A family).
	"INVALID_DEFAULT_LOCALE", "INVALID_DEFAULT_THEME", "INVALID_TIMEZONE",
	"INVALID_RETENTION_DAYS", "INVALID_EXPIRATION_ACTION",
	// D2 login rate limiting.
	"RATE_LIMITED",
	// GOAL-004 S4-6 account lock terminal (423).
	// W7 F-009: ACCOUNT_LOCKED/ACCOUNT_DISABLED are retired from the login
	// surface — locked/disabled accounts return the generic UNAUTHORIZED so the
	// login endpoint no longer discloses account state. They remain cataloged
	// for backwards compatibility (see frozenRetiredCodes).
	// F-03 (GOAL-005): self-service account codes.
	"INVALID_PASSWORD", "INVALID_PASSWORD_BODY", "SESSION_NOT_FOUND",
	// F-03 sessions list status filter.
	"INVALID_STATUS_FILTER",
	// W16-F01 (GOAL-025): forced initial-password-change gate.
	"MUST_CHANGE_PASSWORD",
	// F-02 (GOAL-004): data-transfer codes.
	"RESOURCE_NOT_FOUND", "INVALID_CSV", "INVALID_IMPORT_BODY", "INVALID_EXPORT_LIMIT",
	// F-04 (GOAL-006): notification codes.
	"INVALID_SETTINGS_BODY", "NOTIFICATION_NOT_FOUND",
	// S-02 (GOAL-007): file-library codes.
	"INVALID_UPLOAD_BODY", "INVALID_FILE_ID",
	// S-11 (GOAL-011 D-002 §2): login captcha code.
	"INVALID_CAPTCHA",
	// S-09 (GOAL-016 D-002 §3): data-permission codes.
	"INVALID_SCOPE", "INVALID_SCOPE_BODY", "SCOPE_NOT_ENFORCEABLE",
	// S-10 (GOAL-017 D-002 §3/§4): MFA codes.
	"INVALID_MFA_BODY", "MFA_INVALID", "MFA_PROOF_EXPIRED", "MFA_PROOF_EXHAUSTED", "MFA_NOT_ENROLLED", "MFA_PENDING_ONLY", "MFA_ALREADY_ACTIVE",
	// W7 F-007: MFA enrollment step-up requires the current password.
	"MFA_CURRENT_PASSWORD_REQUIRED",
	// W7 F-004: account avatar per-user upload quota.
	"AVATAR_QUOTA_EXCEEDED",
	// W9 (GOAL-010 D-001): brand asset upload codes.
	"ASSET_NOT_FOUND", "INVALID_KIND",
	// S-14 (GOAL-019 D-002 §3): wallet codes.
	"INVALID_WALLET_BODY", "INVALID_WALLET_OWNER", "INVALID_WALLET_ACCOUNT", "INVALID_WALLET_STATUS", "WALLET_NOT_FOUND", "WALLET_OWNER_TAKEN", "WALLET_DISABLED",
	"INSUFFICIENT_BALANCE", "LEDGER_VERSION_CONFLICT", "LEDGER_IDEMPOTENCY_CONFLICT", "INVALID_LEDGER_ENTRY",
	"PRECONDITION_REQUIRED", "INVALID_PRECONDITION",
	// GOAL-020 (D-001 §2): user accounts are auto-created.
	"WALLET_USER_AUTO_ONLY",
	// VP-012 R4 (GOAL-005 D-002 §7): actor-scoped async Job HTTP codes.
	"JOB_NOT_FOUND", "JOB_NOT_CANCELLABLE", "JOB_NOT_RETRYABLE", "JOB_RESULT_NOT_READY", "JOB_RESULT_EXPIRED",
	// VP-017 R7 (workspace-017 GOAL-008): outbound-mail admin surface codes.
	"INVALID_MAIL_CONFIG", "MAIL_SWITCH_REJECTED", "MAIL_SEND_FAILED",
}

// frozenStoredCodes are stable Job terminal codes persisted in Job.error.
// They are cataloged for clients reading a failed Job representation, but are
// not emitted as top-level handler errors.
var frozenStoredCodes = []string{"JOB_ATTEMPTS_EXHAUSTED", "JOB_HANDLER_FAILED"}

// frozenRetiredCodes remain in the error catalog for backward compatibility
// but are no longer emitted by current handlers (W7 F-009: the login endpoint
// no longer discloses lock/disable account state via distinct codes).
var frozenRetiredCodes = []string{"ACCOUNT_LOCKED", "ACCOUNT_DISABLED"}

// frozenOperationalCodes are selected through the R5 mode switch rather than
// emitted as direct string literals at every call site.
var frozenOperationalCodes = []string{"SERVICE_MAINTENANCE", "SERVICE_DEGRADED", "SERVICE_READ_ONLY"}

// frozenDomainCodes are the verbatim domain rejection codes (resources/roles/
// users), plus the dynamic {RESOURCE}_NOT_FOUND family and the F-006 additions
// (LAST_ADMIN / SELF_OPERATION / INVALID_ROLE_REF / ROLE_ASSIGNMENT_FORBIDDEN /
// INVALID_MENU_ITEM_REF) emitted via DomainError in the users/roles handlers.
// ADMIN_ACCOUNT_FORBIDDEN (D-001 P1) guards the target-side delegation
// boundary: non-admin actors cannot reset an admin's password or demote one.
var frozenDomainCodes = []string{
	"USERNAME_TAKEN", "ROLE_KEY_TAKEN", "ROLE_IN_USE", "ROLE_SYSTEM",
	"INVALID_ROLE_KEY", "INVALID_PERMISSION_REF", "ROLE_GRANT_FORBIDDEN",
	"ROLE_NOT_FOUND", "USER_NOT_FOUND", "CATALOG_NOT_FOUND",
	"LAST_ADMIN", "SELF_OPERATION", "INVALID_ROLE_REF", "ROLE_ASSIGNMENT_FORBIDDEN",
	"INVALID_MENU_ITEM_REF", "ADMIN_ACCOUNT_FORBIDDEN",
	// S-01 (GOAL-008 D-002 §3): dictionary domain codes (DomainError / notFound).
	"DICT_TYPE_NOT_FOUND", "DICT_ENTRY_NOT_FOUND", "DICT_TYPE_KEY_TAKEN",
	"DICT_ENTRY_KEY_TAKEN", "DICT_KEY_NOT_FOUND",
	// S-04 (GOAL-010 D-002 §4): scheduled-task codes.
	"TASK_NOT_FOUND", "TASK_RUN_NOT_FOUND", "TASK_KEY_TAKEN", "INVALID_CRON", "INVALID_HANDLER",
	// S-12 (GOAL-012 D-002 §5): recycle-bin codes.
	"RECYCLE_ITEM_NOT_FOUND", "RECYCLE_RESTORE_CONFLICT", "RECYCLE_ITEM_ALREADY_RESTORED",
	// W14 F-06 (GOAL-019): operations detail not found.
	"OPERATION_NOT_FOUND",
	// workspace-018 R3 (GOAL-004 D-001 §3): account email identity codes.
	"EMAIL_INVALID", "EMAIL_TAKEN", "EMAIL_NOT_PENDING",
	"EMAIL_CODE_INVALID", "EMAIL_CODE_EXPIRED", "EMAIL_RESEND_COOLDOWN",
	"EMAIL_SEND_FAILED",
}

var codeLiteralPattern = regexp.MustCompile(`(?:writeError|writeLocalizedError|writeLocalizedFieldError)\(w, [^,]+, [^,]+, "([A-Z_]+)"`)
var notFoundCodePattern = regexp.MustCompile(`NotFoundCode:\s*"([A-Z_]+)"`)
var domainErrorCodePattern = regexp.MustCompile(`Code:\s*"([A-Z_]+)"`)

func collectCodeLiterals(t *testing.T) map[string]bool {
	t.Helper()
	found := map[string]bool{}
	for _, dir := range []string{".", "../auth"} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read dir %s: %v", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
				continue
			}
			if entry.Name() == "localize_test.go" || entry.Name() == "error_contract_test.go" {
				continue
			}
			src, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				t.Fatalf("read %s: %v", entry.Name(), err)
			}
			for _, match := range codeLiteralPattern.FindAllStringSubmatch(string(src), -1) {
				found[match[1]] = true
			}
			for _, match := range notFoundCodePattern.FindAllStringSubmatch(string(src), -1) {
				found[match[1]] = true
			}
			for _, match := range domainErrorCodePattern.FindAllStringSubmatch(string(src), -1) {
				found[match[1]] = true
			}
		}
	}
	return found
}

func TestErrorCodeContractPinnedSet(t *testing.T) {
	found := collectCodeLiterals(t)

	// Every frozen literal must still exist in the shipped code.
	for _, code := range frozenLiteralCodes {
		if !found[code] {
			t.Errorf("frozen literal code %q is no longer emitted by any handler", code)
		}
	}
	// No unexpected literal codes (drift guard).
	allowed := map[string]bool{}
	for _, code := range append(append(append(append(append([]string{}, frozenLiteralCodes...), frozenDomainCodes...), frozenStoredCodes...), frozenOperationalCodes...), frozenRetiredCodes...) {
		allowed[code] = true
	}
	for code := range found {
		if !allowed[code] {
			t.Errorf("new error code %q emitted without a contract update (D-002 appendix A)", code)
		}
	}
}

func TestErrorCatalogCoversFrozenCodesExceptInternal(t *testing.T) {
	// INTERNAL must never be cataloged (no localization, no messageKey).
	if _, ok := errorcatalog.Catalog["INTERNAL"]; ok {
		t.Fatal("INTERNAL must not be cataloged")
	}
	// Every other frozen literal + domain code must have a bilingual entry
	// with a stable messageKey.
	for _, code := range append(append(append(append(append([]string{}, frozenLiteralCodes...), frozenDomainCodes...), frozenStoredCodes...), frozenOperationalCodes...), frozenRetiredCodes...) {
		if code == "INTERNAL" {
			continue
		}
		entry, ok := errorcatalog.Catalog[code]
		if !ok {
			t.Errorf("code %q is not cataloged", code)
			continue
		}
		if entry.En == "" || entry.Zh == "" {
			t.Errorf("code %q missing en/zh message", code)
		}
		if entry.MessageKey == "" {
			t.Errorf("code %q missing messageKey", code)
		}
	}
	// The catalog must not contain codes outside the frozen set.
	for code := range errorcatalog.Catalog {
		known := false
		for _, frozen := range append(append(append(append(append([]string{}, frozenLiteralCodes...), frozenDomainCodes...), frozenStoredCodes...), frozenOperationalCodes...), frozenRetiredCodes...) {
			if code == frozen {
				known = true
				break
			}
		}
		if !known {
			t.Errorf("catalog contains code %q outside the frozen set", code)
		}
	}
}

func TestAuthMiddlewareCodesAreCataloged(t *testing.T) {
	// The auth package shares the same catalog; its emitted codes must be
	// covered (the contract test above scans internal/auth sources too).
	for _, code := range []string{"UNAUTHENTICATED"} {
		if _, ok := errorcatalog.Catalog[code]; !ok {
			t.Errorf("auth middleware code %q is not cataloged", code)
		}
	}
}

var _ = auth.Authenticator{}
var _ = sort.Strings
