// Package migration owns the biz.digital-offer schema (VP-031 · GOAL-002
// D-002 v1.0.0 §2/§3/§4/§10): offers with a frozen entitlement form,
// append-only purchase vouchers and entitlement rows. Placeholder style is
// dialect-neutral '?' (R1 v1.4 §2); timestamps are Unix seconds like wallet.
package migration

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// ModuleID is the biz.digital-offer module owner.
const ModuleID = "biz.digital-offer"

// offerDDL (0070): offers carry the D-002 §2 form mutex — a duration offer
// must have duration_seconds >= 1 and a NULL count, a count offer the mirror
// image. entitlement_form, currency and the form parameters are immutable
// after create; price/name/description/status go through the version column.
var offerDDL = []string{
	`CREATE TABLE digital_offers (
  id                 TEXT PRIMARY KEY,
  name               TEXT NOT NULL,
  description        TEXT NOT NULL DEFAULT '',
  price_amount       INTEGER NOT NULL CHECK (price_amount > 0),
  currency           TEXT NOT NULL,
  entitlement_form   TEXT NOT NULL CHECK (entitlement_form IN ('duration','count')),
  duration_seconds   INTEGER,
  count_per_purchase INTEGER,
  status             TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','on_sale','off_sale')),
  version            INTEGER NOT NULL DEFAULT 0,
  created_at         INTEGER NOT NULL,
  updated_at         INTEGER NOT NULL,
  CHECK (
    (entitlement_form = 'duration' AND duration_seconds IS NOT NULL AND duration_seconds >= 1 AND count_per_purchase IS NULL)
    OR
    (entitlement_form = 'count' AND count_per_purchase IS NOT NULL AND count_per_purchase >= 1 AND duration_seconds IS NULL)
  )
)`,
	`CREATE INDEX idx_digital_offers_status ON digital_offers(status, created_at DESC)`,
	`CREATE TABLE digital_purchases (
  id              TEXT PRIMARY KEY,
  subject_id      TEXT NOT NULL,
  offer_id        TEXT NOT NULL,
  offer_name      TEXT NOT NULL,
  amount          INTEGER NOT NULL CHECK (amount > 0),
  currency        TEXT NOT NULL,
  freeze_entry_id TEXT NOT NULL,
  deduct_entry_id TEXT NOT NULL,
  request_id      TEXT NOT NULL,
  status          TEXT NOT NULL DEFAULT 'fulfilled' CHECK (status = 'fulfilled'),
  created_at      INTEGER NOT NULL,
  UNIQUE (subject_id, request_id)
)`,
	`CREATE INDEX idx_digital_purchases_subject ON digital_purchases(subject_id, created_at DESC)`,
	`CREATE INDEX idx_digital_purchases_offer ON digital_purchases(offer_id)`,
	`CREATE TABLE digital_entitlements (
  id              TEXT PRIMARY KEY,
  subject_id      TEXT NOT NULL,
  offer_id        TEXT NOT NULL,
  purchase_id     TEXT NOT NULL,
  form            TEXT NOT NULL CHECK (form IN ('duration','count')),
  expires_at      INTEGER,
  remaining_count INTEGER,
  status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','voided')),
  created_at      INTEGER NOT NULL,
  updated_at      INTEGER NOT NULL,
  CHECK (
    (form = 'duration' AND expires_at IS NOT NULL AND remaining_count IS NULL)
    OR
    (form = 'count' AND remaining_count IS NOT NULL AND remaining_count >= 0 AND expires_at IS NULL)
  )
)`,
	`CREATE INDEX idx_digital_entitlements_subject ON digital_entitlements(subject_id, status)`,
	`CREATE INDEX idx_digital_entitlements_purchase ON digital_entitlements(purchase_id)`,
}

// offerPGDDL is the postgres variant: money columns (price_amount / amount),
// counts (duration_seconds / count_per_purchase / remaining_count) and Unix
// time columns (created_at / updated_at / expires_at) are BIGINT (R1 v1.4 §3,
// wallet precedent).
var offerPGDDL = []string{
	`CREATE TABLE digital_offers (
  id                 TEXT PRIMARY KEY,
  name               TEXT NOT NULL,
  description        TEXT NOT NULL DEFAULT '',
  price_amount       BIGINT NOT NULL CHECK (price_amount > 0),
  currency           TEXT NOT NULL,
  entitlement_form   TEXT NOT NULL CHECK (entitlement_form IN ('duration','count')),
  duration_seconds   BIGINT,
  count_per_purchase BIGINT,
  status             TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','on_sale','off_sale')),
  version            INTEGER NOT NULL DEFAULT 0,
  created_at         BIGINT NOT NULL,
  updated_at         BIGINT NOT NULL,
  CHECK (
    (entitlement_form = 'duration' AND duration_seconds IS NOT NULL AND duration_seconds >= 1 AND count_per_purchase IS NULL)
    OR
    (entitlement_form = 'count' AND count_per_purchase IS NOT NULL AND count_per_purchase >= 1 AND duration_seconds IS NULL)
  )
)`,
	`CREATE INDEX idx_digital_offers_status ON digital_offers(status, created_at DESC)`,
	`CREATE TABLE digital_purchases (
  id              TEXT PRIMARY KEY,
  subject_id      TEXT NOT NULL,
  offer_id        TEXT NOT NULL,
  offer_name      TEXT NOT NULL,
  amount          BIGINT NOT NULL CHECK (amount > 0),
  currency        TEXT NOT NULL,
  freeze_entry_id TEXT NOT NULL,
  deduct_entry_id TEXT NOT NULL,
  request_id      TEXT NOT NULL,
  status          TEXT NOT NULL DEFAULT 'fulfilled' CHECK (status = 'fulfilled'),
  created_at      BIGINT NOT NULL,
  UNIQUE (subject_id, request_id)
)`,
	`CREATE INDEX idx_digital_purchases_subject ON digital_purchases(subject_id, created_at DESC)`,
	`CREATE INDEX idx_digital_purchases_offer ON digital_purchases(offer_id)`,
	`CREATE TABLE digital_entitlements (
  id              TEXT PRIMARY KEY,
  subject_id      TEXT NOT NULL,
  offer_id        TEXT NOT NULL,
  purchase_id     TEXT NOT NULL,
  form            TEXT NOT NULL CHECK (form IN ('duration','count')),
  expires_at      BIGINT,
  remaining_count BIGINT,
  status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','voided')),
  created_at      BIGINT NOT NULL,
  updated_at      BIGINT NOT NULL,
  CHECK (
    (form = 'duration' AND expires_at IS NOT NULL AND remaining_count IS NULL)
    OR
    (form = 'count' AND remaining_count IS NOT NULL AND remaining_count >= 0 AND expires_at IS NULL)
  )
)`,
	`CREATE INDEX idx_digital_entitlements_subject ON digital_entitlements(subject_id, status)`,
	`CREATE INDEX idx_digital_entitlements_purchase ON digital_entitlements(purchase_id)`,
}

func migrateOffers(tx kernel.Tx, ddl []string) error {
	for _, stmt := range ddl {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return err
		}
	}
	return nil
}

func migrateOffersSQLite(tx kernel.Tx) error { return migrateOffers(tx, offerDDL) }

func migrateOffersPostgres(tx kernel.Tx) error { return migrateOffers(tx, offerPGDDL) }

// Descriptors returns the immutable 0070 digital-offer history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "digital_offers"},
			Version:              70,
			Name:                 "digital_offers",
			Checksum:             kernel.MigrationChecksum(offerDDL, "0070:digital-offers:v1"),
			Apply:                migrateOffersSQLite,
			ApplyPostgres:        migrateOffersPostgres,
		},
	}
}
