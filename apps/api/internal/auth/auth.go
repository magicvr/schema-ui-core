// Package auth implements the R2 authentication core (GOAL-005 D-003 / D-004):
// short-lived JWT access tokens, opaque revocable refresh tokens stored hashed
// in SQLite, bcrypt password verification, and request-level identity wiring so
// business handlers authenticate from the request instead of a process-injected
// static session.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
)

// Sentinel errors surfaced to handlers for mapping to HTTP status codes.
var (
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	ErrAccountLocked      = errors.New("auth: account locked")
	ErrAccountDisabled    = errors.New("auth: account disabled")
	ErrInvalidToken       = errors.New("auth: invalid token")
	ErrExpiredToken       = errors.New("auth: expired token")
	ErrTokenRevoked       = errors.New("auth: token revoked")
)

// MFARequiredError is returned by Login when the account's second factor must
// be completed before tokens are issued (S-10 · GOAL-017 D-002 §3). It carries
// the verified user id so the login handler can begin a second-factor proof.
type MFARequiredError struct {
	UserID string
}

func (e *MFARequiredError) Error() string { return "auth: second factor required" }

// Account-lock policy (GOAL-004 S4-6): 5 consecutive password failures open a
// 15-minute lock window; the lock expires automatically once now passes
// locked_until and a successful login resets the counter.
const (
	LockThresholdFailures = 5
	LockWindow            = 15 * time.Minute
)

// Repository is the account/session persistence surface used by auth.
type Repository interface {
	UserByUsername(string) (*authsession.User, error)
	UserByID(string) (*authsession.User, error)
	CreateRefreshToken(authsession.RefreshToken) error
	RefreshTokenByHash(string) (*authsession.RefreshToken, error)
	RevokeRefreshToken(string, time.Time) error
	PermissionsForUser(string) ([]string, error)
	FeaturesForUser(string) (map[string]bool, error)
	// Account-lock surface (GOAL-004 S4-6).
	RecordLoginFailure(string, int, time.Time, time.Time) (bool, error)
	ResetLoginFailures(string, time.Time) error
	RevokeAllRefreshTokensForUser(string, time.Time) error
}

// MFAEnforcer is the optional second-factor login gate (S-10 · GOAL-017
// D-002 §3): Required reports whether the user must complete a second factor
// before tokens are issued. A nil enforcer keeps the login contract
// byte-identical. Satisfied structurally by the handler.MFAVerifier.
type MFAEnforcer interface {
	Required(userID string) bool
}

// Authenticator holds the signing secret, token TTLs and the auth-session
// repository.
type Authenticator struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	repository Repository
	devSession bool // explicit opt-in: static dev session fallback (M9)
	// mfa is the optional second-factor gate (S-10 · GOAL-017 D-002 §3);
	// nil = the login contract is byte-identical to the pre-MFA behavior.
	mfa MFAEnforcer
	// OnLockOpened is an optional best-effort hook fired when a login failure
	// opens the account-lock window (F-04 · GOAL-006 D-002 §3: system
	// notification source). Nil = no hook; failures never block Login.
	OnLockOpened func(userID string)
}

// New builds an Authenticator. The caller is responsible for a non-empty secret
// in non-development environments (fail-closed at startup). devSession gates
// the explicit local-development fallback to StaticDevSession; it must be false
// in production (acceptance M9).
func New(secret []byte, accessTTL, refreshTTL time.Duration, runner authsession.TxRunner, devSession bool) *Authenticator {
	return NewWithRepository(secret, accessTTL, refreshTTL, authsession.NewRepository(runner), devSession)
}

// NewWithRepository builds an Authenticator over the shared module repository
// constructed by the composition root.
func NewWithRepository(secret []byte, accessTTL, refreshTTL time.Duration, repository Repository, devSession bool) *Authenticator {
	return &Authenticator{secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL, repository: repository, devSession: devSession}
}

// timingDummyHash is compared against when the requested user does not exist,
// so a missing user burns the same bcrypt time as a wrong password and login
// responses cannot be used to enumerate usernames (D2 timing side channel).
var timingDummyHash = func() string {
	h, err := HashPassword("dummy-timing-password", 10)
	if err != nil {
		panic("auth: hash timing dummy: " + err.Error())
	}
	return h
}()

// Login verifies username/password against the store and issues a fresh
// access/refresh token pair. Fail-closed: a missing user and a bad password
// both yield ErrInvalidCredentials (no user enumeration). A locked account
// (locked_until in the future) yields ErrAccountLocked before password work.
func (a *Authenticator) Login(username, password string, now time.Time) (accessToken, refreshToken string, user account.User, err error) {
	u, err := a.repository.UserByUsername(username)
	if errors.Is(err, authsession.ErrNotFound) {
		VerifyPassword(timingDummyHash, password)
		return "", "", account.User{}, ErrInvalidCredentials
	}
	if err != nil {
		return "", "", account.User{}, err
	}
	if u.LockedUntil > now.Unix() {
		return "", "", account.User{}, ErrAccountLocked
	}
	// F-03 (GOAL-005 D-002 §3): a disabled account fails closed before any
	// password work; the admin-facing operation is visible, so the state is
	// surfaced as a distinct 403 rather than a generic credential failure.
	if !u.Enabled {
		return "", "", account.User{}, ErrAccountDisabled
	}
	if !VerifyPassword(u.PasswordHash, password) {
		// Consecutive-failure accounting: reaching the threshold opens the lock
		// window. When the lock just opened, live sessions are revoked too — a
		// locked account must not keep rotating (access tokens still expire
		// normally; refresh rotation is the durable session channel).
		locked, err := a.repository.RecordLoginFailure(u.ID, LockThresholdFailures, now.Add(LockWindow), now)
		if err != nil {
			return "", "", account.User{}, err
		}
		if locked {
			_ = a.repository.RevokeAllRefreshTokensForUser(u.ID, now)
			if a.OnLockOpened != nil {
				a.OnLockOpened(u.ID)
			}
		}
		return "", "", account.User{}, ErrInvalidCredentials
	}
	if u.FailedLoginCount != 0 || u.LockedUntil != 0 {
		if err := a.repository.ResetLoginFailures(u.ID, now); err != nil {
			return "", "", account.User{}, err
		}
	}
	// S-10 (GOAL-017 D-002 §3): the second factor gates token issuance after
	// the password factor succeeded. nil enforcer → original behavior.
	if a.mfa != nil && a.mfa.Required(u.ID) {
		return "", "", account.User{}, &MFARequiredError{UserID: u.ID}
	}
	return a.issue(u, now)
}

// SetMFAEnforcer installs the optional second-factor gate (S-10 · GOAL-017
// D-002 §3). nil restores the pre-MFA login contract.
func (a *Authenticator) SetMFAEnforcer(mfa MFAEnforcer) {
	a.mfa = mfa
}

// IssueTokensFor issues a fresh access/refresh pair for an already
// second-factor-verified user (S-10 · GOAL-017 D-002 §3). Fail-closed on the
// same terminal states as Login (locked / disabled).
func (a *Authenticator) IssueTokensFor(userID string, now time.Time) (accessToken, refreshToken string, user account.User, err error) {
	u, err := a.repository.UserByID(userID)
	if err != nil {
		return "", "", account.User{}, err
	}
	if u.LockedUntil > now.Unix() {
		return "", "", account.User{}, ErrAccountLocked
	}
	if !u.Enabled {
		return "", "", account.User{}, ErrAccountDisabled
	}
	return a.issue(u, now)
}

// Refresh validates a raw refresh token and rotates it: the old token is
// revoked and a new access/refresh pair is issued. Unknown tokens fail closed.
func (a *Authenticator) Refresh(rawRefresh string, now time.Time) (accessToken, newRefresh string, user account.User, err error) {
	rt, err := a.repository.RefreshTokenByHash(HashToken(rawRefresh))
	if errors.Is(err, authsession.ErrNotFound) {
		return "", "", account.User{}, ErrInvalidToken
	}
	if err != nil {
		return "", "", account.User{}, err
	}
	if rt.RevokedAt != nil {
		return "", "", account.User{}, ErrTokenRevoked
	}
	if now.After(rt.ExpiresAt) {
		return "", "", account.User{}, ErrExpiredToken
	}
	// Rotation is atomic: the guarded revoke is a single UPDATE that exactly one
	// concurrent caller can win. Losing the race (already revoked) means the
	// token was already used or revoked — fail closed instead of issuing a
	// second live pair (double-rotation hardening).
	if err := a.repository.RevokeRefreshToken(rt.ID, now); err != nil {
		if errors.Is(err, authsession.ErrAlreadyRevoked) {
			return "", "", account.User{}, ErrTokenRevoked
		}
		return "", "", account.User{}, err
	}
	u, err := a.repository.UserByID(rt.UserID)
	if err != nil {
		return "", "", account.User{}, err
	}
	// Fail-closed if the account is locked or disabled: revoke-on-lock is
	// best-effort, so Refresh must not mint a new pair from a leftover live
	// refresh token. Same 401 envelope as an invalid token (no extra lock or
	// disable oracle on this path).
	if u.LockedUntil > now.Unix() || !u.Enabled {
		return "", "", account.User{}, ErrInvalidToken
	}
	return a.issue(u, now)
}

// Logout revokes the presented refresh token and returns the revoked token's
// user id (for the operation log, I-008-003 §5). It is idempotent: unknown or
// already-revoked tokens are treated as success (user id empty) so a logout
// cannot be replayed.
func (a *Authenticator) Logout(rawRefresh string, now time.Time) (string, error) {
	rt, err := a.repository.RefreshTokenByHash(HashToken(rawRefresh))
	if errors.Is(err, authsession.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if rt.RevokedAt != nil {
		return rt.UserID, nil
	}
	err = a.repository.RevokeRefreshToken(rt.ID, now)
	if err != nil && !errors.Is(err, authsession.ErrAlreadyRevoked) {
		return rt.UserID, err
	}
	// ErrAlreadyRevoked under concurrency: another logout won the race, which is
	// the same end state (token revoked) — treat as success for idempotency.
	return rt.UserID, nil
}

// issue creates a fresh access/refresh pair for an authenticated user and
// persists the (hashed) opaque refresh token. The access token carries the
// user's current token_version so a later password change revokes it.
func (a *Authenticator) issue(u *authsession.User, now time.Time) (string, string, account.User, error) {
	access, err := SignAccessToken(a.secret, u.ID, u.TokenVersion, a.accessTTL, now)
	if err != nil {
		return "", "", account.User{}, err
	}
	raw, hash, err := NewOpaqueToken()
	if err != nil {
		return "", "", account.User{}, err
	}
	rt := authsession.RefreshToken{
		ID:        newID(),
		UserID:    u.ID,
		TokenHash: hash,
		ExpiresAt: now.Add(a.refreshTTL),
		CreatedAt: now,
	}
	if err := a.repository.CreateRefreshToken(rt); err != nil {
		return "", "", account.User{}, err
	}
	acct, err := a.accountFromUser(u)
	if err != nil {
		return "", "", account.User{}, err
	}
	return access, raw, acct, nil
}

// accessClaims extends the registered claims with the user's token_version
// (W4 P0-3). The middleware rejects a token whose version is older than the
// persisted value, so a password change (which bumps the version) revokes
// every already-signed access token immediately instead of leaving a
// ~accessTTL window where a stolen token still works.
type accessClaims struct {
	TokenVersion int `json:"tv"`
	jwt.RegisteredClaims
}

// SignAccessToken mints a short-lived HMAC-SHA256 access token whose subject is
// the user id and which carries the user's current token_version.
func SignAccessToken(secret []byte, userID string, tokenVersion int, ttl time.Duration, now time.Time) (string, error) {
	claims := accessClaims{
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParsedAccessToken is the verified access-token identity: subject user id and
// the token_version claim carried at issue time.
type ParsedAccessToken struct {
	UserID       string
	TokenVersion int
}

// ParseAccessToken verifies signature and expiry, returning the subject user id
// and token_version. Fail-closed: any parse, method or expiry failure is an
// error.
func ParseAccessToken(secret []byte, raw string) (ParsedAccessToken, error) {
	token, err := jwt.ParseWithClaims(raw, &accessClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("auth: unexpected signing method %v", t.Header["alg"])
		}
		return secret, nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return ParsedAccessToken{}, err
	}
	claims, ok := token.Claims.(*accessClaims)
	if !ok || !token.Valid {
		return ParsedAccessToken{}, ErrInvalidToken
	}
	return ParsedAccessToken{UserID: claims.Subject, TokenVersion: claims.TokenVersion}, nil
}

// HashPassword hashes a plaintext password with bcrypt at the given cost.
func HashPassword(plain string, cost int) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	return string(h), err
}

// VerifyPassword reports whether plaintext matches the stored bcrypt hash.
func VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// NewOpaqueToken generates a 256-bit URL-safe random token and returns both the
// raw value (returned to the client once) and its SHA-256 hex hash (persisted).
func NewOpaqueToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	return raw, HashToken(raw), nil
}

// HashToken hex-encodes the SHA-256 of a raw token.
func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand failure is effectively fatal; fall back to a non-cryptographic
		// timestamp+seq id so a broken RNG does not wedge the server silently.
		return fmt.Sprintf("rt-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// accountFromUser builds the identity snapshot, resolving the user's persisted
// permission keys (GOAL-006 S4: resource gates check keys, not role strings).
func (a *Authenticator) accountFromUser(u *authsession.User) (account.User, error) {
	perms, err := a.repository.PermissionsForUser(u.ID)
	if err != nil {
		return account.User{}, fmt.Errorf("resolve permissions for %s: %w", u.ID, err)
	}
	return account.User{ID: u.ID, Name: u.Name, Roles: u.Roles, Permissions: perms}, nil
}

// Features returns the boolean menu projection for an authenticated identity
// (GOAL-006 S5). The explicit dev-session fallback resolves to its own static
// features for parity; real identities resolve from the persisted menu grants.
func (a *Authenticator) Features(user account.User) (map[string]bool, error) {
	if a.devSession && user.ID == account.StaticDevSession().User.ID {
		return account.StaticDevSession().Features, nil
	}
	return a.repository.FeaturesForUser(user.ID)
}

// UserByID exposes the auth-session identity row to core auth-event logging.
func (a *Authenticator) UserByID(id string) (*authsession.User, error) {
	return a.repository.UserByID(id)
}

// --- request-level identity ---

type ctxKey struct{}

// WithIdentity stores the authenticated account.User in the request context.
func WithIdentity(ctx context.Context, u account.User) context.Context {
	return context.WithValue(ctx, ctxKey{}, u)
}

// IdentityFrom returns the authenticated identity and true, or the zero value
// and false when the request was not authenticated.
func IdentityFrom(ctx context.Context) (account.User, bool) {
	u, ok := ctx.Value(ctxKey{}).(account.User)
	return u, ok
}

func bearer(r *http.Request) (string, bool) {
	const prefix = "Bearer "
	h := r.Header.Get("Authorization")
	if len(h) > len(prefix) && strings.EqualFold(h[:len(prefix)], prefix) {
		return strings.TrimSpace(h[len(prefix):]), true
	}
	return "", false
}

// Middleware verifies the Bearer access token and loads the user, storing the
// identity in context. Any failure writes 401 UNAUTHENTICATED (fail-closed),
// unless devSession is enabled — the explicit local-development fallback that
// substitutes StaticDevSession instead of rejecting.
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := bearer(r)
		if !ok {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no access token")
			return
		}
		parsed, err := ParseAccessToken(a.secret, raw)
		if err != nil {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "invalid or expired access token")
			return
		}
		u, err := a.repository.UserByID(parsed.UserID)
		if err != nil {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "unknown access token subject")
			return
		}
		// W4 P0-3: a password change bumps the user's token_version; an access
		// token signed at the older version is rejected immediately. Stale
		// tokens and unknown subjects both fail closed as UNAUTHENTICATED.
		if parsed.TokenVersion != u.TokenVersion {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "access token superseded")
			return
		}
		// F-03 (GOAL-005 D-002 §3): disabling bumps token_version too, but the
		// persisted state is authoritative — a disabled identity is rejected
		// here even if a token somehow carried the current version (fail
		// closed, same envelope as superseded; no state oracle).
		if !u.Enabled {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "account disabled")
			return
		}
		acct, err := a.accountFromUser(u)
		if err != nil {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "could not resolve identity")
			return
		}
		next.ServeHTTP(w, r.WithContext(WithIdentity(r.Context(), acct)))
	})
}

// injectDevSession is the explicit opt-in local-development fallback (M9). It
// is only reachable when devSession is true.
func (a *Authenticator) injectDevSession(w http.ResponseWriter, r *http.Request, next http.Handler) {
	dev := account.StaticDevSession().User
	next.ServeHTTP(w, r.WithContext(WithIdentity(r.Context(), dev)))
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": code, "message": message})
}

// writeLocalizedError writes the VP-007 S4 envelope: cataloged codes get a
// locale-negotiated message + messageKey + Content-Language; uncataloged codes
// stay English with no key (I-L10N-004 path a).
func writeLocalizedError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	locale := errorcatalog.Negotiate(r)
	body, contentLanguage, cataloged := errorcatalog.Body(code, message, locale)
	if !cataloged {
		writeError(w, status, code, message)
		return
	}
	w.Header().Set("Content-Language", contentLanguage)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}