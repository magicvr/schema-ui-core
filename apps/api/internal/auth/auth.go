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
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// Sentinel errors surfaced to handlers for mapping to HTTP status codes.
var (
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	ErrInvalidToken       = errors.New("auth: invalid token")
	ErrExpiredToken       = errors.New("auth: expired token")
	ErrTokenRevoked       = errors.New("auth: token revoked")
)

// Authenticator holds the signing secret, token TTLs and the auth store.
type Authenticator struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	store      *store.Store
	devSession bool // explicit opt-in: static dev session fallback (M9)
}

// New builds an Authenticator. The caller is responsible for a non-empty secret
// in non-development environments (fail-closed at startup). devSession gates
// the explicit local-development fallback to StaticDevSession; it must be false
// in production (acceptance M9).
func New(secret []byte, accessTTL, refreshTTL time.Duration, st *store.Store, devSession bool) *Authenticator {
	return &Authenticator{secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL, store: st, devSession: devSession}
}

// Login verifies username/password against the store and issues a fresh
// access/refresh token pair. Fail-closed: a missing user and a bad password
// both yield ErrInvalidCredentials (no user enumeration).
func (a *Authenticator) Login(username, password string, now time.Time) (accessToken, refreshToken string, user account.User, err error) {
	u, err := a.store.UserByUsername(username)
	if errors.Is(err, store.ErrNotFound) {
		return "", "", account.User{}, ErrInvalidCredentials
	}
	if err != nil {
		return "", "", account.User{}, err
	}
	if !VerifyPassword(u.PasswordHash, password) {
		return "", "", account.User{}, ErrInvalidCredentials
	}
	return a.issue(u, now)
}

// Refresh validates a raw refresh token and rotates it: the old token is
// revoked and a new access/refresh pair is issued. Unknown tokens fail closed.
func (a *Authenticator) Refresh(rawRefresh string, now time.Time) (accessToken, newRefresh string, user account.User, err error) {
	rt, err := a.store.RefreshTokenByHash(HashToken(rawRefresh))
	if errors.Is(err, store.ErrNotFound) {
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
	if err := a.store.RevokeRefreshToken(rt.ID, now); err != nil && !errors.Is(err, store.ErrAlreadyRevoked) {
		return "", "", account.User{}, err
	}
	u, err := a.store.UserByID(rt.UserID)
	if err != nil {
		return "", "", account.User{}, err
	}
	return a.issue(u, now)
}

// Logout revokes the presented refresh token and returns the revoked token's
// user id (for the operation log, I-008-003 §5). It is idempotent: unknown or
// already-revoked tokens are treated as success (user id empty) so a logout
// cannot be replayed.
func (a *Authenticator) Logout(rawRefresh string, now time.Time) (string, error) {
	rt, err := a.store.RefreshTokenByHash(HashToken(rawRefresh))
	if errors.Is(err, store.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if rt.RevokedAt != nil {
		return rt.UserID, nil
	}
	return rt.UserID, a.store.RevokeRefreshToken(rt.ID, now)
}

// issue creates a fresh access/refresh pair for an authenticated user and
// persists the (hashed) opaque refresh token.
func (a *Authenticator) issue(u *store.User, now time.Time) (string, string, account.User, error) {
	access, err := SignAccessToken(a.secret, u.ID, a.accessTTL, now)
	if err != nil {
		return "", "", account.User{}, err
	}
	raw, hash, err := NewOpaqueToken()
	if err != nil {
		return "", "", account.User{}, err
	}
	rt := store.RefreshToken{
		ID:        newID(),
		UserID:    u.ID,
		TokenHash: hash,
		ExpiresAt: now.Add(a.refreshTTL),
		CreatedAt: now,
	}
	if err := a.store.CreateRefreshToken(rt); err != nil {
		return "", "", account.User{}, err
	}
	acct, err := a.accountFromUser(u)
	if err != nil {
		return "", "", account.User{}, err
	}
	return access, raw, acct, nil
}

// SignAccessToken mints a short-lived HMAC-SHA256 access token whose subject is
// the user id.
func SignAccessToken(secret []byte, userID string, ttl time.Duration, now time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseAccessToken verifies signature and expiry, returning the subject user id.
// Fail-closed: any parse, method or expiry failure is an error.
func ParseAccessToken(secret []byte, raw string) (string, error) {
	token, err := jwt.ParseWithClaims(raw, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("auth: unexpected signing method %v", t.Header["alg"])
		}
		return secret, nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}
	return claims.Subject, nil
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
// permission keys (GOAL-006 S4: records gates check keys, not role strings).
func (a *Authenticator) accountFromUser(u *store.User) (account.User, error) {
	perms, err := a.store.PermissionsForUser(u.ID)
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
	return a.store.FeaturesForUser(user.ID)
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
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "no access token")
			return
		}
		userID, err := ParseAccessToken(a.secret, raw)
		if err != nil {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "invalid or expired access token")
			return
		}
		u, err := a.store.UserByID(userID)
		if err != nil {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "unknown access token subject")
			return
		}
		acct, err := a.accountFromUser(u)
		if err != nil {
			if a.devSession {
				a.injectDevSession(w, r, next)
				return
			}
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "could not resolve identity")
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
