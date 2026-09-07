// Package store owns the biz.digital-offer domain queries (VP-031 · GOAL-002
// D-002 v1.0.0 §2/§3/§4): offers with a frozen entitlement form, append-only
// purchase vouchers and entitlement rows. Purchase rows are the money-path
// receipt — insert-only, no update or delete path exists.
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pagination"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TxRunner is the platform persistence boundary consumed by the repository.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// Domain sentinels mapped by the handler to frozen error codes (D-002 §9).
var (
	ErrNotFound                = errors.New("digital offer not found")
	ErrPurchaseNotFound        = errors.New("digital purchase not found")
	ErrEntitlementNotFound     = errors.New("digital entitlement not found")
	ErrOfferNotOnSale          = errors.New("digital offer is not on sale")
	ErrCurrencyMismatch        = errors.New("offer currency does not match the subject wallet account")
	ErrSubjectNotFound         = errors.New("subject not found")
	ErrRequestIdConflict       = errors.New("request id already used with a different offer")
	ErrEntitlementInsufficient = errors.New("insufficient entitlement count")
	ErrInvalidOffer            = errors.New("invalid digital offer")
	ErrFormConflict            = errors.New("entitlement form fields are immutable")
	ErrVersionConflict         = errors.New("digital offer version conflict")
	errPurchaseUniqueRace      = errors.New("digital purchase request race")
)

// Offer forms and statuses (D-002 §2).
const (
	FormDuration = "duration"
	FormCount    = "count"

	StatusDraft   = "draft"
	StatusOnSale  = "on_sale"
	StatusOffSale = "off_sale"
)

// Entitlement stored statuses (D-002 §3): expired/exhausted are derived
// predicates, never stored transitions.
const (
	EntitlementActive = "active"
	EntitlementVoided = "voided"
)

// Offer is one sellable digital service item (D-002 §2).
type Offer struct {
	ID               string
	Name             string
	Description      string
	PriceAmount      int64
	Currency         string
	EntitlementForm  string
	DurationSeconds  int64 // duration form only
	CountPerPurchase int64 // count form only
	Status           string
	Version          int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Purchase is one fulfilled money-path receipt (D-002 §4.1). Append-only.
type Purchase struct {
	ID            string
	SubjectID     string
	OfferID       string
	OfferName     string
	Amount        int64
	Currency      string
	FreezeEntryID string
	DeductEntryID string
	RequestID     string
	Status        string
	CreatedAt     time.Time
}

// Entitlement is one granted right (D-002 §3). Validity is derived lazily.
type Entitlement struct {
	ID             string
	SubjectID      string
	OfferID        string
	PurchaseID     string
	Form           string
	ExpiresAt      *time.Time // duration form
	RemainingCount *int64     // count form
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// OfferFilter carries admin offer list parameters.
type OfferFilter struct {
	Q        string
	Status   string
	Page     int
	PageSize int
}

// PurchaseFilter carries admin purchase list parameters.
type PurchaseFilter struct {
	Q         string
	SubjectID string
	OfferID   string
	Page      int
	PageSize  int
}

// EntitlementFilter carries admin entitlement list parameters.
type EntitlementFilter struct {
	Q         string
	SubjectID string
	OfferID   string
	Status    string
	Page      int
	PageSize  int
}

// Repository owns the digital-offer domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the repository over a platform transaction runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// Runner exposes the persistence boundary so the service can open
// caller-owned transactions spanning store and wallet operations (D-002 §4.2).
func (r *Repository) Runner() TxRunner { return r.runner }

func isUniqueViolation(err error) bool { return kernel.IsUniqueViolation(err) }

func pageBounds(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return page, pageSize
}

// InsertOfferInTx creates one offer row inside a caller-owned transaction
// (fail-closed audit pairing, D-002 §7).
func (r *Repository) InsertOfferInTx(tx kernel.Tx, o Offer) error {
	_, err := tx.Exec(context.Background(),
		`INSERT INTO digital_offers (id, name, description, price_amount, currency, entitlement_form, duration_seconds, count_per_purchase, status, version, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.ID, o.Name, o.Description, o.PriceAmount, o.Currency, o.EntitlementForm,
		nullableInt(o.EntitlementForm == FormDuration, o.DurationSeconds),
		nullableInt(o.EntitlementForm == FormCount, o.CountPerPurchase),
		o.Status, o.Version, o.CreatedAt.Unix(), o.UpdatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("insert digital offer: %w", err)
	}
	return nil
}

// GetOffer loads one offer by id.
func (r *Repository) GetOffer(id string) (*Offer, error) {
	var offer Offer
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		return scanOffer(tx, id, &offer)
	})
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

// GetOfferInTx loads one offer inside a caller-owned transaction.
func (r *Repository) GetOfferInTx(tx kernel.Tx, id string) (*Offer, error) {
	var offer Offer
	if err := scanOffer(tx, id, &offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

func scanOffer(tx kernel.Tx, id string, offer *Offer) error {
	var price, created, updated int64
	var duration, count any
	err := tx.QueryRow(context.Background(),
		`SELECT id, name, description, price_amount, currency, entitlement_form, duration_seconds, count_per_purchase, status, version, created_at, updated_at
		 FROM digital_offers WHERE id = ?`, id,
	).Scan(&offer.ID, &offer.Name, &offer.Description, &price, &offer.Currency,
		&offer.EntitlementForm, &duration, &count, &offer.Status, &offer.Version, &created, &updated)
	if errors.Is(err, kernel.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("query digital offer: %w", err)
	}
	offer.PriceAmount = price
	if v, ok := duration.(int64); ok {
		offer.DurationSeconds = v
	}
	if v, ok := count.(int64); ok {
		offer.CountPerPurchase = v
	}
	offer.CreatedAt = time.Unix(created, 0)
	offer.UpdatedAt = time.Unix(updated, 0)
	return nil
}

// ListOffers returns the admin offer page (newest first, wallet precedent).
func (r *Repository) ListOffers(filter OfferFilter) ([]Offer, int, error) {
	offers := []Offer{}
	total := 0
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		where, args, err := offerWhere(filter)
		if err != nil {
			return err
		}
		if err := tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM digital_offers "+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count digital offers: %w", err)
		}
		page, pageSize := pageBounds(filter.Page, filter.PageSize)
		rows, err := tx.Query(context.Background(),
			`SELECT id, name, description, price_amount, currency, entitlement_form, COALESCE(duration_seconds, 0), COALESCE(count_per_purchase, 0), status, version, created_at, updated_at
			 FROM digital_offers `+where+" ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?",
			append(args, pageSize, pagination.Offset(page, pageSize, total))...)
		if err != nil {
			return fmt.Errorf("list digital offers: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var o Offer
			if err := scanOfferRow(rows, &o); err != nil {
				return err
			}
			offers = append(offers, o)
		}
		return rows.Err()
	})
	return offers, total, err
}

func offerWhere(filter OfferFilter) (string, []any, error) {
	where := "WHERE 1=1"
	args := []any{}
	if q := searchQ(filter.Q); q != "" {
		where += " AND (LOWER(name) LIKE LOWER(?) ESCAPE '\\' OR LOWER(id) LIKE LOWER(?) ESCAPE '\\')"
		like := "%" + escapeLike(q) + "%"
		args = append(args, like, like)
	}
	switch filter.Status {
	case "":
	case StatusDraft, StatusOnSale, StatusOffSale:
		where += " AND status = ?"
		args = append(args, filter.Status)
	default:
		return "", nil, ErrInvalidOffer
	}
	return where, args, nil
}

type offerRows interface {
	Scan(dest ...any) error
}

func scanOfferRow(rows offerRows, o *Offer) error {
	var price, created, updated int64
	if err := rows.Scan(&o.ID, &o.Name, &o.Description, &price, &o.Currency, &o.EntitlementForm, &o.DurationSeconds, &o.CountPerPurchase, &o.Status, &o.Version, &created, &updated); err != nil {
		return fmt.Errorf("scan digital offer: %w", err)
	}
	o.PriceAmount = price
	o.CreatedAt = time.Unix(created, 0)
	o.UpdatedAt = time.Unix(updated, 0)
	return nil
}

// ListOnSale returns the C-end catalog (D-002 §5.3): on_sale offers only,
// oldest first.
func (r *Repository) ListOnSale() ([]Offer, error) {
	offers := []Offer{}
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		rows, err := tx.Query(context.Background(),
			`SELECT id, name, description, price_amount, currency, entitlement_form, COALESCE(duration_seconds, 0), COALESCE(count_per_purchase, 0), status, version, created_at, updated_at
			 FROM digital_offers WHERE status = ? ORDER BY created_at ASC, id ASC`, StatusOnSale)
		if err != nil {
			return fmt.Errorf("list on-sale offers: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var o Offer
			if err := scanOfferRow(rows, &o); err != nil {
				return err
			}
			offers = append(offers, o)
		}
		return rows.Err()
	})
	return offers, err
}

// UpdateOfferInTx applies mutable fields (name/description/price/status) with
// the optimistic-lock version guard (D-002 §2). Immutable form fields are not
// part of the UPDATE set by construction.
func (r *Repository) UpdateOfferInTx(tx kernel.Tx, id string, name, description *string, price *int64, status *string, expectedVersion int64, now time.Time) (*Offer, error) {
	sets := []string{"updated_at = ?", "version = version + 1"}
	args := []any{now.Unix()}
	if name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *name)
	}
	if description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *description)
	}
	if price != nil {
		sets = append(sets, "price_amount = ?")
		args = append(args, *price)
	}
	if status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *status)
	}
	args = append(args, id, expectedVersion)
	res, err := tx.Exec(context.Background(),
		"UPDATE digital_offers SET "+strings.Join(sets, ", ")+" WHERE id = ? AND version = ?", args...)
	if err != nil {
		return nil, fmt.Errorf("update digital offer: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("digital offer update result: %w", err)
	}
	if affected == 0 {
		var probe Offer
		if scanErr := scanOffer(tx, id, &probe); scanErr != nil {
			return nil, ErrNotFound
		}
		return nil, ErrVersionConflict
	}
	return r.GetOfferInTx(tx, id)
}

// GetPurchaseByRequestInTx resolves the idempotency guard row (D-002 §4.4
// attempt step 1). ErrPurchaseNotFound = no prior purchase for the pair.
func (r *Repository) GetPurchaseByRequestInTx(tx kernel.Tx, subjectID, requestID string) (*Purchase, error) {
	var p Purchase
	var amount, created int64
	err := tx.QueryRow(context.Background(),
		`SELECT id, subject_id, offer_id, offer_name, amount, currency, freeze_entry_id, deduct_entry_id, request_id, status, created_at
		 FROM digital_purchases WHERE subject_id = ? AND request_id = ?`, subjectID, requestID,
	).Scan(&p.ID, &p.SubjectID, &p.OfferID, &p.OfferName, &amount, &p.Currency,
		&p.FreezeEntryID, &p.DeductEntryID, &p.RequestID, &p.Status, &created)
	if errors.Is(err, kernel.ErrNoRows) {
		return nil, ErrPurchaseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query digital purchase by request: %w", err)
	}
	p.Amount = amount
	p.CreatedAt = time.Unix(created, 0)
	return &p, nil
}

// InsertPurchaseInTx appends one fulfilled purchase voucher inside the
// caller-owned transaction. A UNIQUE(subject_id, request_id) violation is the
// D-002 §4.4 purchase race; it surfaces as errPurchaseUniqueRace for the
// service retry loop (classified via kernel.IsUniqueViolation on our OWN table
// insert — no dependency on wallet-internal sentinels).
func (r *Repository) InsertPurchaseInTx(tx kernel.Tx, p Purchase) error {
	_, err := tx.Exec(context.Background(),
		`INSERT INTO digital_purchases (id, subject_id, offer_id, offer_name, amount, currency, freeze_entry_id, deduct_entry_id, request_id, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.SubjectID, p.OfferID, p.OfferName, p.Amount, p.Currency,
		p.FreezeEntryID, p.DeductEntryID, p.RequestID, p.Status, p.CreatedAt.Unix(),
	)
	if isUniqueViolation(err) {
		return errPurchaseUniqueRace
	}
	if err != nil {
		return fmt.Errorf("insert digital purchase: %w", err)
	}
	return nil
}

// ListPurchases returns the admin purchase page (read-only audit view).
func (r *Repository) ListPurchases(filter PurchaseFilter) ([]Purchase, int, error) {
	purchases := []Purchase{}
	total := 0
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		where := "WHERE 1=1"
		args := []any{}
		if q := searchQ(filter.Q); q != "" {
			where += " AND (LOWER(id) LIKE LOWER(?) ESCAPE '\\' OR LOWER(offer_name) LIKE LOWER(?) ESCAPE '\\')"
			like := "%" + escapeLike(q) + "%"
			args = append(args, like, like)
		}
		if filter.SubjectID != "" {
			where += " AND subject_id = ?"
			args = append(args, filter.SubjectID)
		}
		if filter.OfferID != "" {
			where += " AND offer_id = ?"
			args = append(args, filter.OfferID)
		}
		if err := tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM digital_purchases "+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count digital purchases: %w", err)
		}
		page, pageSize := pageBounds(filter.Page, filter.PageSize)
		rows, err := tx.Query(context.Background(),
			`SELECT id, subject_id, offer_id, offer_name, amount, currency, freeze_entry_id, deduct_entry_id, request_id, status, created_at
			 FROM digital_purchases `+where+" ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?",
			append(args, pageSize, pagination.Offset(page, pageSize, total))...)
		if err != nil {
			return fmt.Errorf("list digital purchases: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var p Purchase
			var amount, created int64
			if err := rows.Scan(&p.ID, &p.SubjectID, &p.OfferID, &p.OfferName, &amount, &p.Currency, &p.FreezeEntryID, &p.DeductEntryID, &p.RequestID, &p.Status, &created); err != nil {
				return fmt.Errorf("scan digital purchase: %w", err)
			}
			p.Amount = amount
			p.CreatedAt = time.Unix(created, 0)
			purchases = append(purchases, p)
		}
		return rows.Err()
	})
	return purchases, total, err
}

// InsertEntitlementInTx appends one entitlement row inside the caller-owned
// transaction (form mirrors the offer at purchase time, D-002 §3).
func (r *Repository) InsertEntitlementInTx(tx kernel.Tx, e Entitlement) error {
	var expires any
	if e.ExpiresAt != nil {
		expires = e.ExpiresAt.Unix()
	}
	var remaining any
	if e.RemainingCount != nil {
		remaining = *e.RemainingCount
	}
	_, err := tx.Exec(context.Background(),
		`INSERT INTO digital_entitlements (id, subject_id, offer_id, purchase_id, form, expires_at, remaining_count, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.SubjectID, e.OfferID, e.PurchaseID, e.Form, expires, remaining, e.Status, e.CreatedAt.Unix(), e.UpdatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("insert digital entitlement: %w", err)
	}
	return nil
}

// GetEntitlement loads one entitlement row.
func (r *Repository) GetEntitlement(id string) (*Entitlement, error) {
	var e Entitlement
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		return scanEntitlement(tx, id, &e)
	})
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// ListEntitlementsBySubjectOfferInTx returns every entitlement row for one
// (subject, offer) pair — the §5.1 Check aggregate reads them all and derives
// the reason (no status filter here).
func (r *Repository) ListEntitlementsBySubjectOfferInTx(tx kernel.Tx, subjectID, offerID string) ([]Entitlement, error) {
	rows, err := tx.Query(context.Background(),
		`SELECT id, subject_id, offer_id, purchase_id, form, expires_at, remaining_count, status, created_at, updated_at
		 FROM digital_entitlements WHERE subject_id = ? AND offer_id = ? ORDER BY created_at ASC, id ASC`,
		subjectID, offerID)
	if err != nil {
		return nil, fmt.Errorf("list entitlements by subject/offer: %w", err)
	}
	defer rows.Close()
	var out []Entitlement
	for rows.Next() {
		var e Entitlement
		var created, updated int64
		var expires, remaining any
		if err := rows.Scan(&e.ID, &e.SubjectID, &e.OfferID, &e.PurchaseID, &e.Form, &expires, &remaining, &e.Status, &created, &updated); err != nil {
			return nil, fmt.Errorf("scan digital entitlement: %w", err)
		}
		decodeEntitlementTimes(&e, expires, remaining, created, updated)
		out = append(out, e)
	}
	return out, rows.Err()
}

// GetEntitlementInTx loads one entitlement inside a caller-owned transaction
// (nested Run is forbidden — D-002 §7 void must stay on one tx).
func (r *Repository) GetEntitlementInTx(tx kernel.Tx, id string) (*Entitlement, error) {
	var e Entitlement
	if err := scanEntitlement(tx, id, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func scanEntitlement(tx kernel.Tx, id string, e *Entitlement) error {
	var created, updated int64
	var expires, remaining any
	err := tx.QueryRow(context.Background(),
		`SELECT id, subject_id, offer_id, purchase_id, form, expires_at, remaining_count, status, created_at, updated_at
		 FROM digital_entitlements WHERE id = ?`, id,
	).Scan(&e.ID, &e.SubjectID, &e.OfferID, &e.PurchaseID, &e.Form, &expires, &remaining, &e.Status, &created, &updated)
	if errors.Is(err, kernel.ErrNoRows) {
		return ErrEntitlementNotFound
	}
	if err != nil {
		return fmt.Errorf("query digital entitlement: %w", err)
	}
	decodeEntitlementTimes(e, expires, remaining, created, updated)
	return nil
}

func decodeEntitlementTimes(e *Entitlement, expires, remaining any, created, updated int64) {
	if v, ok := expires.(int64); ok {
		t := time.Unix(v, 0)
		e.ExpiresAt = &t
	}
	if v, ok := remaining.(int64); ok {
		e.RemainingCount = &v
	}
	e.CreatedAt = time.Unix(created, 0)
	e.UpdatedAt = time.Unix(updated, 0)
}

// ListEntitlements returns the admin entitlement page.
func (r *Repository) ListEntitlements(filter EntitlementFilter) ([]Entitlement, int, error) {
	entitlements := []Entitlement{}
	total := 0
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		where := "WHERE 1=1"
		args := []any{}
		if q := searchQ(filter.Q); q != "" {
			where += " AND (LOWER(id) LIKE LOWER(?) ESCAPE '\\' OR LOWER(subject_id) LIKE LOWER(?) ESCAPE '\\')"
			like := "%" + escapeLike(q) + "%"
			args = append(args, like, like)
		}
		if filter.SubjectID != "" {
			where += " AND subject_id = ?"
			args = append(args, filter.SubjectID)
		}
		if filter.OfferID != "" {
			where += " AND offer_id = ?"
			args = append(args, filter.OfferID)
		}
		switch filter.Status {
		case "":
		case EntitlementActive, EntitlementVoided:
			where += " AND status = ?"
			args = append(args, filter.Status)
		default:
			return ErrInvalidOffer
		}
		if err := tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM digital_entitlements "+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count digital entitlements: %w", err)
		}
		page, pageSize := pageBounds(filter.Page, filter.PageSize)
		rows, err := tx.Query(context.Background(),
			`SELECT id, subject_id, offer_id, purchase_id, form, expires_at, remaining_count, status, created_at, updated_at
			 FROM digital_entitlements `+where+" ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?",
			append(args, pageSize, pagination.Offset(page, pageSize, total))...)
		if err != nil {
			return fmt.Errorf("list digital entitlements: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var e Entitlement
			var created, updated int64
			var expires, remaining any
			if err := rows.Scan(&e.ID, &e.SubjectID, &e.OfferID, &e.PurchaseID, &e.Form, &expires, &remaining, &e.Status, &created, &updated); err != nil {
				return fmt.Errorf("scan digital entitlement: %w", err)
			}
			decodeEntitlementTimes(&e, expires, remaining, created, updated)
			entitlements = append(entitlements, e)
		}
		return rows.Err()
	})
	return entitlements, total, err
}

// GetEntitlementByPurchaseInTx loads the entitlement granted by one purchase
// (D-002 §4.4 replay path: the read-back result carries its entitlement).
func (r *Repository) GetEntitlementByPurchaseInTx(tx kernel.Tx, purchaseID string) (*Entitlement, error) {
	var e Entitlement
	var created, updated int64
	var expires, remaining any
	err := tx.QueryRow(context.Background(),
		`SELECT id, subject_id, offer_id, purchase_id, form, expires_at, remaining_count, status, created_at, updated_at
		 FROM digital_entitlements WHERE purchase_id = ?`, purchaseID,
	).Scan(&e.ID, &e.SubjectID, &e.OfferID, &e.PurchaseID, &e.Form, &expires, &remaining, &e.Status, &created, &updated)
	if errors.Is(err, kernel.ErrNoRows) {
		return nil, ErrEntitlementNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query digital entitlement by purchase: %w", err)
	}
	decodeEntitlementTimes(&e, expires, remaining, created, updated)
	return &e, nil
}

// VoidEntitlementInTx transitions active → voided inside the caller-owned
// transaction (D-002 §7). RowsAffected=0 with an existing row means the row
// was already voided (idempotent success); a missing row is ErrEntitlementNotFound.
func (r *Repository) VoidEntitlementInTx(tx kernel.Tx, id string, now time.Time) (alreadyVoided bool, err error) {
	res, err := tx.Exec(context.Background(),
		`UPDATE digital_entitlements SET status = ?, updated_at = ? WHERE id = ? AND status = ?`,
		EntitlementVoided, now.Unix(), id, EntitlementActive)
	if err != nil {
		return false, fmt.Errorf("void digital entitlement: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("void digital entitlement result: %w", err)
	}
	if affected == 1 {
		return false, nil
	}
	var e Entitlement
	if scanErr := scanEntitlement(tx, id, &e); scanErr != nil {
		return false, ErrEntitlementNotFound
	}
	return true, nil
}

// ConsumeCandidate is one active count row in consumption order (D-002 §5.2:
// created_at ASC, id ASC).
type ConsumeCandidate struct {
	ID             string
	RemainingCount int64
}

// ListConsumeCandidatesInTx returns the active count rows for (subject, offer)
// in consumption order.
func (r *Repository) ListConsumeCandidatesInTx(tx kernel.Tx, subjectID, offerID string) ([]ConsumeCandidate, error) {
	rows, err := tx.Query(context.Background(),
		`SELECT id, remaining_count FROM digital_entitlements
		 WHERE subject_id = ? AND offer_id = ? AND form = ? AND status = ? AND remaining_count > 0
		 ORDER BY created_at ASC, id ASC`,
		subjectID, offerID, FormCount, EntitlementActive)
	if err != nil {
		return nil, fmt.Errorf("list consume candidates: %w", err)
	}
	defer rows.Close()
	var out []ConsumeCandidate
	for rows.Next() {
		var c ConsumeCandidate
		if err := rows.Scan(&c.ID, &c.RemainingCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// DecrementEntitlementInTx atomically takes exactly `take` from one row. The
// five-condition predicate (id/subject/offer/form/status/balance) is the
// concurrency guard (D-002 §5.2); RowsAffected=0 means the row lost a race
// and must not be counted.
func (r *Repository) DecrementEntitlementInTx(tx kernel.Tx, id, subjectID, offerID string, take int64, now time.Time) (bool, error) {
	res, err := tx.Exec(context.Background(),
		`UPDATE digital_entitlements SET remaining_count = remaining_count - ?, updated_at = ?
		 WHERE id = ? AND subject_id = ? AND offer_id = ? AND form = ? AND status = ? AND remaining_count >= ?`,
		take, now.Unix(), id, subjectID, offerID, FormCount, EntitlementActive, take)
	if err != nil {
		return false, fmt.Errorf("consume digital entitlement: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("consume digital entitlement result: %w", err)
	}
	return affected == 1, nil
}

// nullableInt renders the form-conditional column: the active form stores its
// parameter, the other form stores SQL NULL (D-002 §2/§3 CHECK mutex).
func nullableInt(active bool, v int64) any {
	if !active {
		return nil
	}
	return v
}

// searchQ normalizes a user-supplied search term: trimmed and hard-capped to
// 100 runes (byte slicing could split a multi-byte UTF-8 rune and make
// PostgreSQL reject the parameter).
func searchQ(q string) string {
	q = strings.TrimSpace(q)
	runes := []rune(q)
	if len(runes) > 100 {
		runes = runes[:100]
	}
	return string(runes)
}

// escapeLike escapes SQL LIKE metacharacters (wallet W13 F-011 principle).
func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
