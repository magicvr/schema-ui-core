package mfa

import (
	"strings"
	"testing"
	"time"
)

// RFC 6238 Appendix B test vectors (SHA1). The 8-digit official values are
// truncated to 6 digits by taking the same HOTP value modulo 1000000.
func TestTotpRFC6238Vectors(t *testing.T) {
	const secret = "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ" // base32 of ASCII "12345678901234567890"
	cases := []struct {
		unix  int64
		want6 string
	}{
		{59, "287082"},
		{1111111109, "081804"},
		{1111111111, "050471"},
		{1234567890, "005924"},
		{2000000000, "279037"},
		{20000000000, "353130"},
	}
	for _, c := range cases {
		got, err := totpCode(secret, c.unix/30)
		if err != nil {
			t.Fatalf("totpCode(%d): %v", c.unix, err)
		}
		if got != c.want6 {
			t.Fatalf("totpCode(%d) = %s, want %s", c.unix, got, c.want6)
		}
	}
}

func TestValidateTotpWindowAndReplay(t *testing.T) {
	const secret = "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"
	now := time.Unix(1_000_000_000, 0)
	step := now.Unix() / totpPeriodSeconds
	code, err := totpCode(secret, step)
	if err != nil {
		t.Fatal(err)
	}

	// Same step validates.
	matched, ok := ValidateTotp(secret, code, now, 1, 0)
	if !ok || matched != step {
		t.Fatalf("same-step = %d, %v; want %d, true", matched, ok, step)
	}
	// Same-step replay rejected once lastUsedStep == matched.
	if _, ok := ValidateTotp(secret, code, now, 1, step); ok {
		t.Fatalf("same-step replay accepted")
	}
	// A neighbouring step (30s earlier) still validates within the window.
	prevCode, _ := totpCode(secret, step-1)
	if _, ok := ValidateTotp(secret, prevCode, now.Add(-30*time.Second), 1, 0); !ok {
		t.Fatalf("windowed step rejected")
	}
	// Wrong code rejected; garbage secret rejected.
	if _, ok := ValidateTotp(secret, "000000", now, 1, 0); ok {
		t.Fatalf("wrong code accepted")
	}
	if _, ok := ValidateTotp("!!!", code, now, 1, 0); ok {
		t.Fatalf("garbage secret accepted")
	}
}

func TestGenerateSecretAndURL(t *testing.T) {
	s1, err := GenerateSecret()
	if err != nil || len(s1) != 32 { // 20 bytes → 32 base32 chars
		t.Fatalf("secret = %q, %v", s1, err)
	}
	s2, _ := GenerateSecret()
	if s1 == s2 {
		t.Fatalf("secrets must be random")
	}
	url := otpauthURL("Schema UI Core", "admin", s1)
	if !strings.HasPrefix(url, "otpauth://totp/") || !strings.Contains(url, "secret="+s1) {
		t.Fatalf("otpauth url = %q", url)
	}
}
