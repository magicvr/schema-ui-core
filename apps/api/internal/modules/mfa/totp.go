// RFC 6238 TOTP self-implementation (S-10 · GOAL-017 D-002 §1): HMAC-SHA1,
// 30-second period, 6 digits, ±window validation with same-step replay
// rejection. Self-implemented per repository discipline (GOAL-010 D-002 §7:
// no new third-party dependency; D-001 §5) and locked by the RFC 6238
// Appendix B test vectors.
package mfa

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

const (
	totpPeriodSeconds = 30
	totpDigits        = 6
	totpWindow        = 1 // ±1 step = 3 candidates
)

// GenerateSecret returns a 20-byte (160-bit) random TOTP secret, Base32
// encoded without padding (RFC 4226 key-length recommendation).
func GenerateSecret() (string, error) {
	raw := make([]byte, 20)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("mfa: generate totp secret: %w", err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw), nil
}

// totpCode computes the 6-digit TOTP code for a secret at a given time step
// (RFC 6238 HOTP truncation, 6 digits).
func totpCode(secretBase32 string, step int64) (string, error) {
	secret, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(strings.TrimSpace(secretBase32)))
	if err != nil {
		return "", fmt.Errorf("mfa: decode totp secret: %w", err)
	}
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(step))
	mac := hmac.New(sha1.New, secret)
	_, _ = mac.Write(msg[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	code := (uint32(sum[offset])&0x7f)<<24 |
		(uint32(sum[offset+1]) << 16) |
		(uint32(sum[offset+2]) << 8) |
		uint32(sum[offset+3])
	return fmt.Sprintf("%0*d", totpDigits, code%1000000), nil
}

// ValidateTotp checks a 6-digit code within ±window steps around now. It
// returns the matched time step (for same-step replay rejection) and true on
// success. A matched step <= lastUsedStep is rejected as a replay.
func ValidateTotp(secretBase32, code string, now time.Time, window int, lastUsedStep int64) (int64, bool) {
	step := now.Unix() / totpPeriodSeconds
	for offset := -window; offset <= window; offset++ {
		candidate := step + int64(offset)
		want, err := totpCode(secretBase32, candidate)
		if err != nil {
			return 0, false
		}
		if want == strings.TrimSpace(code) {
			if candidate <= lastUsedStep {
				return 0, false // same-window replay rejected (D-002 §6)
			}
			return candidate, true
		}
	}
	return 0, false
}

// otpauthURL builds the standard otpauth:// URI for authenticator apps.
func otpauthURL(issuer, account, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&period=%d&digits=%d",
		urlEscape(issuer), urlEscape(account), secret, urlEscape(issuer), totpPeriodSeconds, totpDigits)
}

func urlEscape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", "%20"), ":", "%3A")
}
