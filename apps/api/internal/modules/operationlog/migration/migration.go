package migration

import (
	"database/sql"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

const ModuleID = "core.operationlog"

var operationLogDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

var operationLogExpandDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

var operationLogSettingsDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogDataTransferDDL (0015 · F-02 GOAL-004 D-002 §3/§4): adds the two
// data-transfer events to the event CHECK (rebuild like 0005/0008/0014).
var operationLogDataTransferDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogFileEventsDDL (0018 · S-02 GOAL-007 D-002 §4): adds the three
// file-library events to the event CHECK (rebuild like 0005/0008/0014/0015).
var operationLogFileEventsDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogDictionaryDDL (0020 · S-01 GOAL-008 D-002 §5): adds the three
// dictionary events to the event CHECK (rebuild like 0005/0008/0014/0015/0018).
var operationLogDictionaryDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogTasksDDL (0022 · S-04 GOAL-010 D-002 §4): adds the three
// scheduled-task events to the event CHECK (rebuild like 0005/0008/0014/0015/0018/0020).
var operationLogTasksDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogCaptchaDDL (0024 · S-11 GOAL-011 D-002 §3): adds the captcha
// settings event to the event CHECK (rebuild like 0005/0008/0014/0015/0018/0020/0022).
var operationLogCaptchaDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}


// operationLogDataPermissionDDL (0028 · S-09 GOAL-016 D-002 §3): adds the two
// data-permission events to the event CHECK (rebuild like 0005/0008/0014/0015/0018/0020/0022/0024/0026).
var operationLogDataPermissionDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogMFADDL (0030 · S-10 GOAL-017 D-002 §2): adds the six MFA events
// to the event CHECK (rebuild like 0005/0008/0014/0015/0018/0020/0022/0024/0026/0028).
var operationLogMFADDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogWalletDDL (0032 · S-14 GOAL-019 D-002 §2): adds the six wallet
// events to the event CHECK (rebuild like 0005/0008/0014/0015/0018/0020/0022/0024/0026/0028/0030).
var operationLogWalletDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login','wallet.account-create','wallet.account-update','wallet.adjust','wallet.freeze','wallet.unfreeze','wallet.reconcile')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

func migrateOperationLogWallet(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogWalletDDL, "wallet-events-expanded")
}

// operationLogAvatarEventsDDL (0036 · W13 T-05 GOAL-014): adds the
// account.avatar-change event to the CHECK (rebuild like 0032/0034).
var operationLogAvatarEventsDDL = []string{
`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login','wallet.account-create','wallet.account-update','wallet.adjust','wallet.freeze','wallet.unfreeze','wallet.reconcile','wallet.deduct-frozen','account.avatar-change')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogWalletDeductDDL (0034 · GOAL-021 D-001 §3): adds the
// wallet.deduct-frozen event to the CHECK (rebuild like 0032).
var operationLogWalletDeductDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login','wallet.account-create','wallet.account-update','wallet.adjust','wallet.freeze','wallet.unfreeze','wallet.reconcile','wallet.deduct-frozen')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

func migrateOperationLogWalletDeduct(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogWalletDeductDDL, "wallet-deduct-events-expanded")
}

func migrateOperationLogAvatarEvents(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogAvatarEventsDDL, "avatar-events-expanded")
}

// operationLogRecycleDDL (0026 · S-12 GOAL-012 D-002 §5): adds the two
// recycle events to the event CHECK (rebuild like 0005/0008/0014/0015/0018/0020/0022/0024).
var operationLogRecycleDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogAccountEventsDDL (0014 · F-03 GOAL-005 D-002 §3): adds the five
// account-lifecycle events to the event CHECK. SQLite cannot ALTER a CHECK, so
// the table is rebuilt like 0005/0008.
var operationLogAccountEventsDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// Descriptors returns the immutable 0004, 0005 and 0008 operation-log history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log"},
			Version:              4,
			Name:                 "operation_log",
			Checksum:             kernel.MigrationChecksum(operationLogDDL, "0004:operation-log:v1"),
			Apply:                migrateOperationLog,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_expand"},
			Version:              5,
			Name:                 "operation_log_expand",
			Checksum:             kernel.MigrationChecksum(operationLogExpandDDL, "0005:operation-log-expand:v1"),
			Apply:                migrateOperationLogExpand,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_settings"},
			Version:              8,
			Name:                 "operation_log_settings",
			Checksum:             kernel.MigrationChecksum(operationLogSettingsDDL, "0008:operation-log-settings:v1"),
			Apply:                migrateOperationLogSettings,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_account_events"},
			Version:              14,
			Name:                 "operation_log_account_events",
			Checksum:             kernel.MigrationChecksum(operationLogAccountEventsDDL, "0014:operation-log-account-events:v1"),
			Apply:                migrateOperationLogAccountEvents,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_data_transfer"},
			Version:              15,
			Name:                 "operation_log_data_transfer",
			Checksum:             kernel.MigrationChecksum(operationLogDataTransferDDL, "0015:operation-log-data-transfer:v1"),
			Apply:                migrateOperationLogDataTransfer,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_file_events"},
			Version:              18,
			Name:                 "operation_log_file_events",
			Checksum:             kernel.MigrationChecksum(operationLogFileEventsDDL, "0018:operation-log-file-events:v1"),
			Apply:                migrateOperationLogFileEvents,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_dictionary"},
			Version:              20,
			Name:                 "operation_log_dictionary",
			Checksum:             kernel.MigrationChecksum(operationLogDictionaryDDL, "0020:operation-log-dictionary:v1"),
			Apply:                migrateOperationLogDictionary,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_tasks"},
			Version:              22,
			Name:                 "operation_log_tasks",
			Checksum:             kernel.MigrationChecksum(operationLogTasksDDL, "0022:operation-log-tasks:v1"),
			Apply:                migrateOperationLogTasks,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_captcha"},
			Version:              24,
			Name:                 "operation_log_captcha",
			Checksum:             kernel.MigrationChecksum(operationLogCaptchaDDL, "0024:operation-log-captcha:v1"),
			Apply:                migrateOperationLogCaptcha,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_data_permission"},
			Version:              28,
			Name:                 "operation_log_data_permission",
			Checksum:             kernel.MigrationChecksum(operationLogDataPermissionDDL, "0028:operation-log-data-permission:v1"),
			Apply:                migrateOperationLogDataPermission,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_mfa"},
			Version:              30,
			Name:                 "operation_log_mfa",
			Checksum:             kernel.MigrationChecksum(operationLogMFADDL, "0030:operation-log-mfa:v1"),
			Apply:                migrateOperationLogMFA,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_wallet"},
			Version:              32,
			Name:                 "operation_log_wallet",
			Checksum:             kernel.MigrationChecksum(operationLogWalletDDL, "0032:operation-log-wallet:v1"),
			Apply:                migrateOperationLogWallet,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_wallet_deduct"},
			Version:              34,
			Name:                 "operation_log_wallet_deduct",
			Checksum:             kernel.MigrationChecksum(operationLogWalletDeductDDL, "0034:operation-log-wallet-deduct:v1"),
			Apply:                migrateOperationLogWalletDeduct,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_avatar_events"},
			Version:              36,
			Name:                 "operation_log_avatar_events",
			Checksum:             kernel.MigrationChecksum(operationLogAvatarEventsDDL, "0036:operation-log-avatar-events:v1"),
			Apply:                migrateOperationLogAvatarEvents,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_recycle"},
			Version:              26,
			Name:                 "operation_log_recycle",
			Checksum:             kernel.MigrationChecksum(operationLogRecycleDDL, "0026:operation-log-recycle:v1"),
			Apply:                migrateOperationLogRecycle,
		},
	}
}

func migrateOperationLog(tx *sql.Tx) error {
	for _, stmt := range operationLogDDL {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("create operation_log: %w", err)
		}
	}
	return nil
}

func migrateOperationLogExpand(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogExpandDDL, "expanded")
}

func migrateOperationLogSettings(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogSettingsDDL, "settings-expanded")
}

func migrateOperationLogAccountEvents(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogAccountEventsDDL, "account-events-expanded")
}

func migrateOperationLogDataTransfer(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogDataTransferDDL, "data-transfer-expanded")
}

func migrateOperationLogFileEvents(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogFileEventsDDL, "file-events-expanded")
}

func migrateOperationLogDictionary(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogDictionaryDDL, "dictionary-events-expanded")
}

func migrateOperationLogTasks(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogTasksDDL, "tasks-events-expanded")
}

func migrateOperationLogCaptcha(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogCaptchaDDL, "captcha-events-expanded")
}

func migrateOperationLogRecycle(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogRecycleDDL, "recycle-events-expanded")
}

func migrateOperationLogDataPermission(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogDataPermissionDDL, "data-permission-events-expanded")
}

func migrateOperationLogMFA(tx *sql.Tx) error {
	return rebuildOperationLog(tx, operationLogMFADDL, "mfa-events-expanded")
}

func rebuildOperationLog(tx *sql.Tx, ddl []string, label string) error {
	if _, err := tx.Exec(`ALTER TABLE operation_log RENAME TO operation_log_old`); err != nil {
		return fmt.Errorf("rename operation_log: %w", err)
	}
	if _, err := tx.Exec(ddl[0]); err != nil {
		return fmt.Errorf("create operation_log %s: %w", label, err)
	}
	if _, err := tx.Exec(
		`INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, detail, created_at)
		 SELECT id, event, actor_id, actor_name, record_id, detail, created_at FROM operation_log_old`,
	); err != nil {
		return fmt.Errorf("migrate operation_log rows: %w", err)
	}
	if _, err := tx.Exec(`DROP TABLE operation_log_old`); err != nil {
		return fmt.Errorf("drop operation_log_old: %w", err)
	}
	if _, err := tx.Exec(ddl[1]); err != nil {
		return fmt.Errorf("create operation_log index: %w", err)
	}
	return nil
}