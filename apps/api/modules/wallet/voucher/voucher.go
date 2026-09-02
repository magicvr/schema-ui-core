// Package voucher owns the prepaid instrument (vouchers/gift cards) domain
// (VP-029 · GOAL-002 / GOAL-003): batch generation, code hashing, voiding, and
// atomic single-transaction redemption into wallet ledgers.
package voucher

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// Domain sentinels.
var (
	ErrNotFound               = errors.New("voucher not found")
	ErrVoucherAlreadyRedeemed = errors.New("voucher already redeemed")
	ErrVoucherVoid            = errors.New("voucher is void")
	ErrVoucherExpired         = errors.New("voucher has expired")
	ErrVoucherInvalid         = errors.New("invalid voucher status")
	ErrVoucherConflict        = errors.New("voucher redemption conflict (race lost)")
	ErrSubjectNotFound        = errors.New("subject not found")
	ErrCurrencyMismatch       = errors.New("voucher currency mismatch")
	ErrInvalidInput           = errors.New("invalid voucher input")
)

// Status represents the lifecycle status of a prepaid voucher.
type Status string

const (
	StatusUnused   Status = "unused"
	StatusRedeemed Status = "redeemed"
	StatusVoid     Status = "void"
)

// Voucher is a stored prepaid voucher row.
type Voucher struct {
	ID          string     `json:"id"`
	BatchID     string     `json:"batchId"`
	CodeHash    string     `json:"-"`
	CodePrefix  string     `json:"codePrefix"`
	Amount      int64      `json:"amount"`
	Currency    string     `json:"currency"`
	Status      Status     `json:"status"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	RedeemedBy  *string    `json:"redeemedBy,omitempty"`
	RedeemedAt  *time.Time `json:"redeemedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// GeneratedVoucher pairs the created voucher metadata with its one-time plaintext code.
// The plaintext code is NEVER stored in database or audit logs.
type GeneratedVoucher struct {
	Voucher Voucher `json:"voucher"`
	Code    string  `json:"code"`
}

// RedeemResult contains the receipt of a successful voucher redemption.
type RedeemResult struct {
	VoucherID string `json:"voucherId"`
	BatchID   string `json:"batchId"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	AccountID string `json:"accountId"`
	EntryID   string `json:"entryId"`
	Balance   int64  `json:"balanceAfter"`
}

// alphabet Crockford Base32-like (unambiguous uppercase characters: no 0/O, 1/I/L).
const codeAlphabet = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"

// GenerateCode creates a high-entropy 24-character plaintext voucher code,
// its 6-character non-secret prefix, and its hex SHA-256 digest.
// 24 characters with 32-char alphabet = 24 * 5 = 120 bits of entropy (>80 bit baseline).
func GenerateCode() (code string, prefix string, hash string, err error) {
	var sb strings.Builder
	sb.Grow(24)
	alphaLen := big.NewInt(int64(len(codeAlphabet)))
	for i := 0; i < 24; i++ {
		n, err := rand.Int(rand.Reader, alphaLen)
		if err != nil {
			return "", "", "", fmt.Errorf("generate voucher code: %w", err)
		}
		sb.WriteByte(codeAlphabet[n.Int64()])
	}
	code = sb.String()
	prefix = code[:6]
	hash = HashCode(code)
	return code, prefix, hash, nil
}

// HashCode computes the canonical SHA-256 hex digest of a voucher code.
func HashCode(code string) string {
	clean := strings.ToUpper(strings.TrimSpace(code))
	h := sha256.Sum256([]byte(clean))
	return hex.EncodeToString(h[:])
}

// ConstantTimeCompare compares two code hashes in constant time to resist timing attacks.
func ConstantTimeCompare(hash1, hash2 string) bool {
	return subtle.ConstantTimeCompare([]byte(hash1), []byte(hash2)) == 1
}

// newID generates a time-ordered millisecond prefix + random hex ID.
func newID(now time.Time) (string, error) {
	randBytes := make([]byte, 12)
	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("voucher: random bytes: %w", err)
	}
	return fmt.Sprintf("%016x%s", now.UnixMilli(), hex.EncodeToString(randBytes)), nil
}
