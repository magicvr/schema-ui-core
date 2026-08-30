package migration

import (
	"context"
	"fmt"
	"regexp"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
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

func migrateOperationLogWallet(tx kernel.Tx) error {
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

func migrateOperationLogWalletDeduct(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogWalletDeductDDL, "wallet-deduct-events-expanded")
}

func migrateOperationLogAvatarEvents(tx kernel.Tx) error {
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

// ---- postgres-flavored operation_log DDL (R3 dual-dialect ledger; R1 v1.4
// §3/§4). The canonical sqlite literals above stay untouched (checksums
// stable). The postgres bodies derive from them so the event CHECK lists stay
// in lockstep; the only dialect difference is the Unix time columns
// (created_at / archived_at) becoming BIGINT.
var pgTimeColRe = regexp.MustCompile(`(?m)^(\s*)(created_at|archived_at)\s+INTEGER NOT NULL`)

func pgTimeDDL(stmts []string) []string {
	out := make([]string, len(stmts))
	for i, s := range stmts {
		out[i] = pgTimeColRe.ReplaceAllString(s, "$1$2 BIGINT NOT NULL")
	}
	return out
}

var (
	operationLogPGDDL                   = pgTimeDDL(operationLogDDL)
	operationLogExpandPGDDL             = pgTimeDDL(operationLogExpandDDL)
	operationLogSettingsPGDDL           = pgTimeDDL(operationLogSettingsDDL)
	operationLogAccountEventsPGDDL      = pgTimeDDL(operationLogAccountEventsDDL)
	operationLogDataTransferPGDDL       = pgTimeDDL(operationLogDataTransferDDL)
	operationLogFileEventsPGDDL         = pgTimeDDL(operationLogFileEventsDDL)
	operationLogDictionaryPGDDL         = pgTimeDDL(operationLogDictionaryDDL)
	operationLogTasksPGDDL              = pgTimeDDL(operationLogTasksDDL)
	operationLogCaptchaPGDDL            = pgTimeDDL(operationLogCaptchaDDL)
	operationLogRecyclePGDDL            = pgTimeDDL(operationLogRecycleDDL)
	operationLogDataPermissionPGDDL     = pgTimeDDL(operationLogDataPermissionDDL)
	operationLogMFAPGDDL                = pgTimeDDL(operationLogMFADDL)
	operationLogWalletPGDDL             = pgTimeDDL(operationLogWalletDDL)
	operationLogWalletDeductPGDDL       = pgTimeDDL(operationLogWalletDeductDDL)
	operationLogAvatarEventsPGDDL       = pgTimeDDL(operationLogAvatarEventsDDL)
	operationLogWalletJobsPGDDL         = pgTimeDDL(operationLogWalletJobsDDL)
	operationLogServiceCredentialsPGDDL = pgTimeDDL(operationLogServiceCredentialsDDL)
	operationLogMailEventsPGDDL         = pgTimeDDL(operationLogMailEventsDDL)
	operationLogArchivePGDDL            = pgTimeDDL(operationLogArchiveDDL)
)

// pgExecDDL builds a postgres apply that runs the given (postgres-flavored)
// DDL statements in order.
func pgExecDDL(ddl []string) func(kernel.Tx) error {
	return func(tx kernel.Tx) error {
		for _, stmt := range ddl {
			if _, err := tx.Exec(context.Background(), stmt); err != nil {
				return fmt.Errorf("apply postgres DDL: %w", err)
			}
		}
		return nil
	}
}

// pgRebuild builds a postgres apply that runs the operation_log rebuild
// (rename → create → copy → drop → index) with postgres-flavored DDL.
func pgRebuild(ddl []string, label string) func(kernel.Tx) error {
	return func(tx kernel.Tx) error {
		return rebuildOperationLog(tx, ddl, label)
	}
}

// pgRebuildWithCorrelation is pgRebuild for the 0043/0045 migrations whose
// postgres apply must also preserve the operation_log_correlation side table
// (its FK depends on operation_log, which the rebuild rename breaks).
func pgRebuildWithCorrelation(ddl []string, label string) func(kernel.Tx) error {
	return func(tx kernel.Tx) error {
		return rebuildOperationLogWithCorrelation(tx, ddl, label)
	}
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
			ApplyPostgres:        pgExecDDL(operationLogPGDDL),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_expand"},
			Version:              5,
			Name:                 "operation_log_expand",
			Checksum:             kernel.MigrationChecksum(operationLogExpandDDL, "0005:operation-log-expand:v1"),
			Apply:                migrateOperationLogExpand,
			ApplyPostgres:        pgRebuild(operationLogExpandPGDDL, `expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_settings"},
			Version:              8,
			Name:                 "operation_log_settings",
			Checksum:             kernel.MigrationChecksum(operationLogSettingsDDL, "0008:operation-log-settings:v1"),
			Apply:                migrateOperationLogSettings,
			ApplyPostgres:        pgRebuild(operationLogSettingsPGDDL, `settings-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_account_events"},
			Version:              14,
			Name:                 "operation_log_account_events",
			Checksum:             kernel.MigrationChecksum(operationLogAccountEventsDDL, "0014:operation-log-account-events:v1"),
			Apply:                migrateOperationLogAccountEvents,
			ApplyPostgres:        pgRebuild(operationLogAccountEventsPGDDL, `account-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_data_transfer"},
			Version:              15,
			Name:                 "operation_log_data_transfer",
			Checksum:             kernel.MigrationChecksum(operationLogDataTransferDDL, "0015:operation-log-data-transfer:v1"),
			Apply:                migrateOperationLogDataTransfer,
			ApplyPostgres:        pgRebuild(operationLogDataTransferPGDDL, `data-transfer-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_file_events"},
			Version:              18,
			Name:                 "operation_log_file_events",
			Checksum:             kernel.MigrationChecksum(operationLogFileEventsDDL, "0018:operation-log-file-events:v1"),
			Apply:                migrateOperationLogFileEvents,
			ApplyPostgres:        pgRebuild(operationLogFileEventsPGDDL, `file-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_dictionary"},
			Version:              20,
			Name:                 "operation_log_dictionary",
			Checksum:             kernel.MigrationChecksum(operationLogDictionaryDDL, "0020:operation-log-dictionary:v1"),
			Apply:                migrateOperationLogDictionary,
			ApplyPostgres:        pgRebuild(operationLogDictionaryPGDDL, `dictionary-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_tasks"},
			Version:              22,
			Name:                 "operation_log_tasks",
			Checksum:             kernel.MigrationChecksum(operationLogTasksDDL, "0022:operation-log-tasks:v1"),
			Apply:                migrateOperationLogTasks,
			ApplyPostgres:        pgRebuild(operationLogTasksPGDDL, `tasks-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_captcha"},
			Version:              24,
			Name:                 "operation_log_captcha",
			Checksum:             kernel.MigrationChecksum(operationLogCaptchaDDL, "0024:operation-log-captcha:v1"),
			Apply:                migrateOperationLogCaptcha,
			ApplyPostgres:        pgRebuild(operationLogCaptchaPGDDL, `captcha-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_data_permission"},
			Version:              28,
			Name:                 "operation_log_data_permission",
			Checksum:             kernel.MigrationChecksum(operationLogDataPermissionDDL, "0028:operation-log-data-permission:v1"),
			Apply:                migrateOperationLogDataPermission,
			ApplyPostgres:        pgRebuild(operationLogDataPermissionPGDDL, `data-permission-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_mfa"},
			Version:              30,
			Name:                 "operation_log_mfa",
			Checksum:             kernel.MigrationChecksum(operationLogMFADDL, "0030:operation-log-mfa:v1"),
			Apply:                migrateOperationLogMFA,
			ApplyPostgres:        pgRebuild(operationLogMFAPGDDL, `mfa-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_wallet"},
			Version:              32,
			Name:                 "operation_log_wallet",
			Checksum:             kernel.MigrationChecksum(operationLogWalletDDL, "0032:operation-log-wallet:v1"),
			Apply:                migrateOperationLogWallet,
			ApplyPostgres:        pgRebuild(operationLogWalletPGDDL, `wallet-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_wallet_deduct"},
			Version:              34,
			Name:                 "operation_log_wallet_deduct",
			Checksum:             kernel.MigrationChecksum(operationLogWalletDeductDDL, "0034:operation-log-wallet-deduct:v1"),
			Apply:                migrateOperationLogWalletDeduct,
			ApplyPostgres:        pgRebuild(operationLogWalletDeductPGDDL, `wallet-deduct-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_avatar_events"},
			Version:              36,
			Name:                 "operation_log_avatar_events",
			Checksum:             kernel.MigrationChecksum(operationLogAvatarEventsDDL, "0036:operation-log-avatar-events:v1"),
			Apply:                migrateOperationLogAvatarEvents,
			ApplyPostgres:        pgRebuild(operationLogAvatarEventsPGDDL, `avatar-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_recycle"},
			Version:              26,
			Name:                 "operation_log_recycle",
			Checksum:             kernel.MigrationChecksum(operationLogRecycleDDL, "0026:operation-log-recycle:v1"),
			Apply:                migrateOperationLogRecycle,
			ApplyPostgres:        pgRebuild(operationLogRecyclePGDDL, `recycle-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_correlation"},
			Version:              41,
			Name:                 "operation_log_correlation",
			Checksum:             kernel.MigrationChecksum(operationLogCorrelationDDL, "0041:operation-log-correlation:v1"),
			Apply:                migrateOperationLogCorrelation,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_wallet_jobs"},
			Version:              43,
			Name:                 "operation_log_wallet_jobs",
			Checksum:             kernel.MigrationChecksum(operationLogWalletJobsDDL, "0043:operation-log-wallet-jobs:v1"),
			Apply:                migrateOperationLogWalletJobs,
			ApplyPostgres:        pgRebuildWithCorrelation(operationLogWalletJobsPGDDL, `wallet-job-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_service_credentials"},
			Version:              45,
			Name:                 "operation_log_service_credentials",
			Checksum:             kernel.MigrationChecksum(operationLogServiceCredentialsDDL, "0045:operation-log-service-credentials:v1"),
			Apply:                migrateOperationLogServiceCredentials,
			ApplyPostgres:        pgRebuildWithCorrelation(operationLogServiceCredentialsPGDDL, `service-credential-events-expanded`),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_archive"},
			Version:              47,
			Name:                 "operation_log_archive",
			Checksum:             kernel.MigrationChecksum(operationLogArchiveDDL, "0047:operation-log-archive:v1"),
			Apply:                migrateOperationLogArchive,
			ApplyPostgres:        pgExecDDL(operationLogArchivePGDDL),
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_session"},
			Version:              48,
			Name:                 "operation_log_session",
			Checksum:             kernel.MigrationChecksum(operationLogSessionDDL, "0048:operation-log-session:v1"),
			Apply:                migrateOperationLogSession,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operation_log_mail_events"},
			Version:              53,
			Name:                 "operation_log_mail_events",
			Checksum:             kernel.MigrationChecksum(operationLogMailEventsDDL, "0053:operation-log-mail-events:v1"),
			Apply:                migrateOperationLogMailEvents,
			ApplyPostgres:        MailEventsPGApply(),
		},
	}
}

// operationLogMailEventsDDL (0053 · VP-017 R7): adds the outbound-mail admin
// surface audit events (channel config update + test send) to the frozen
// CHECK enumeration.
var operationLogMailEventsDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login','wallet.account-create','wallet.account-update','wallet.adjust','wallet.freeze','wallet.unfreeze','wallet.reconcile','wallet.deduct-frozen','account.avatar-change','wallet.reconcile.queued','wallet.reconcile.failed','wallet.reconcile.cancelled','service-credentials.create','service-credentials.use','service-credentials.revoke','mail.channel-update','mail.test-send')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogArchiveDDL (0047): cold store for expired audit rows.
var operationLogArchiveDDL = []string{
	`CREATE TABLE operation_log_archive (
  id          TEXT PRIMARY KEY,
  event       TEXT NOT NULL,
  actor_id    TEXT NOT NULL,
  actor_name  TEXT NOT NULL,
  record_id   TEXT,
  detail      TEXT,
  created_at  INTEGER NOT NULL,
  archived_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_archive_created_at ON operation_log_archive(created_at DESC)`,
	`CREATE TABLE operation_log_archive_correlation (
  operation_id   TEXT PRIMARY KEY,
  correlation_id TEXT NOT NULL
)`,
}

func migrateOperationLogArchive(tx kernel.Tx) error {
	for _, stmt := range operationLogArchiveDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create operation_log_archive: %w", err)
		}
	}
	return nil
}

var operationLogSessionDDL = []string{
	`CREATE TABLE operation_log_session (
  operation_id TEXT PRIMARY KEY REFERENCES operation_log(id) ON DELETE CASCADE,
  session_id   TEXT NOT NULL
)`,
	`CREATE INDEX idx_operation_log_session_id ON operation_log_session(session_id)`,
	`CREATE TABLE operation_log_archive_session (
  operation_id TEXT PRIMARY KEY,
  session_id   TEXT NOT NULL
)`,
}

func migrateOperationLogSession(tx kernel.Tx) error {
	for _, stmt := range operationLogSessionDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create operation_log_session: %w", err)
		}
	}
	return nil
}

func migrateOperationLog(tx kernel.Tx) error {
	for _, stmt := range operationLogDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create operation_log: %w", err)
		}
	}
	return nil
}

func migrateOperationLogExpand(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogExpandDDL, "expanded")
}

func migrateOperationLogSettings(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogSettingsDDL, "settings-expanded")
}

func migrateOperationLogAccountEvents(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogAccountEventsDDL, "account-events-expanded")
}

func migrateOperationLogDataTransfer(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogDataTransferDDL, "data-transfer-expanded")
}

func migrateOperationLogFileEvents(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogFileEventsDDL, "file-events-expanded")
}

func migrateOperationLogDictionary(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogDictionaryDDL, "dictionary-events-expanded")
}

func migrateOperationLogTasks(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogTasksDDL, "tasks-events-expanded")
}

func migrateOperationLogCaptcha(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogCaptchaDDL, "captcha-events-expanded")
}

func migrateOperationLogRecycle(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogRecycleDDL, "recycle-events-expanded")
}

func migrateOperationLogCorrelation(tx kernel.Tx) error {
	for _, stmt := range operationLogCorrelationDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create operation_log_correlation: %w", err)
		}
	}
	return nil
}

func migrateOperationLogWalletJobs(tx kernel.Tx) error {
	return rebuildOperationLogWithCorrelation(tx, operationLogWalletJobsDDL, "wallet-job-events-expanded")
}

func migrateOperationLogServiceCredentials(tx kernel.Tx) error {
	return rebuildOperationLogWithCorrelation(tx, operationLogServiceCredentialsDDL, "service-credential-events-expanded")
}

// rebuildOperationLogWithCorrelation is the shared dialect-neutral dance used
// by 0043/0045: back up the correlation side table (operation_log_correlation
// FK depends on operation_log, which the rebuild rename breaks), drop it,
// rebuild operation_log, recreate it, restore the rows. Works on sqlite and
// postgres (dialect difference lives in the DDL slice).
func rebuildOperationLogWithCorrelation(tx kernel.Tx, ddl []string, label string) error {
	if _, err := tx.Exec(context.Background(), `CREATE TEMP TABLE operation_log_correlation_backup AS
SELECT operation_id, correlation_id FROM operation_log_correlation`); err != nil {
		return fmt.Errorf("backup operation log correlations: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE operation_log_correlation`); err != nil {
		return fmt.Errorf("drop operation log correlations before rebuild: %w", err)
	}
	if err := rebuildOperationLog(tx, ddl, label); err != nil {
		return err
	}
	for _, statement := range operationLogCorrelationDDL {
		if _, err := tx.Exec(context.Background(), statement); err != nil {
			return fmt.Errorf("recreate operation log correlations: %w", err)
		}
	}
	if _, err := tx.Exec(context.Background(), `INSERT INTO operation_log_correlation (operation_id, correlation_id)
SELECT operation_id, correlation_id FROM operation_log_correlation_backup`); err != nil {
		return fmt.Errorf("restore operation log correlations: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE operation_log_correlation_backup`); err != nil {
		return fmt.Errorf("drop operation log correlation backup: %w", err)
	}
	return nil
}

func migrateOperationLogDataPermission(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogDataPermissionDDL, "data-permission-events-expanded")
}

func migrateOperationLogMFA(tx kernel.Tx) error {
	return rebuildOperationLog(tx, operationLogMFADDL, "mfa-events-expanded")
}

// MailEventsDDL exposes the 0053 DDL for checksum tooling and tests.
func MailEventsDDL() []string { return operationLogMailEventsDDL }

func migrateOperationLogMailEvents(tx kernel.Tx) error {
	return rebuildOperationLogMailEvents(tx, operationLogMailEventsDDL)
}

// MailEventsPGApply returns the postgres-flavored 0053 apply.
func MailEventsPGApply() func(kernel.Tx) error {
	return func(tx kernel.Tx) error {
		return rebuildOperationLogMailEvents(tx, pgTimeDDL(operationLogMailEventsDDL))
	}
}

// rebuildOperationLogMailEvents is the shared dialect-neutral dance used by
// 0053: like rebuildOperationLogWithCorrelation, plus the 0048
// operation_log_session side table — BOTH carry an FK that the rebuild rename
// would rewrite onto operation_log_old (leaving dangling references after the
// drop), so each is backed up, dropped before the rename, recreated after,
// and refilled.
func rebuildOperationLogMailEvents(tx kernel.Tx, ddl []string) error {
	if _, err := tx.Exec(context.Background(), `CREATE TEMP TABLE operation_log_correlation_backup AS
SELECT operation_id, correlation_id FROM operation_log_correlation`); err != nil {
		return fmt.Errorf("backup operation log correlations: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `CREATE TEMP TABLE operation_log_session_backup AS
SELECT operation_id, session_id FROM operation_log_session`); err != nil {
		return fmt.Errorf("backup operation log sessions: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE operation_log_session`); err != nil {
		return fmt.Errorf("drop operation log sessions before rebuild: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE operation_log_correlation`); err != nil {
		return fmt.Errorf("drop operation log correlations before rebuild: %w", err)
	}
	if err := rebuildOperationLog(tx, ddl, "mail-events-expanded"); err != nil {
		return err
	}
	for _, statement := range operationLogCorrelationDDL {
		if _, err := tx.Exec(context.Background(), statement); err != nil {
			return fmt.Errorf("recreate operation log correlations: %w", err)
		}
	}
	// Recreate ONLY the session table + index (operationLogSessionDDL[2:] is
	// the 0048-era archive_session table, which still exists untouched).
	for _, statement := range operationLogSessionDDL[:2] {
		if _, err := tx.Exec(context.Background(), statement); err != nil {
			return fmt.Errorf("recreate operation log sessions: %w", err)
		}
	}
	if _, err := tx.Exec(context.Background(), `INSERT INTO operation_log_correlation (operation_id, correlation_id)
SELECT operation_id, correlation_id FROM operation_log_correlation_backup`); err != nil {
		return fmt.Errorf("restore operation log correlations: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `INSERT INTO operation_log_session (operation_id, session_id)
SELECT operation_id, session_id FROM operation_log_session_backup`); err != nil {
		return fmt.Errorf("restore operation log sessions: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE operation_log_correlation_backup`); err != nil {
		return fmt.Errorf("drop operation log correlation backup: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE operation_log_session_backup`); err != nil {
		return fmt.Errorf("drop operation log session backup: %w", err)
	}
	return nil
}

func rebuildOperationLog(tx kernel.Tx, ddl []string, label string) error {
	if _, err := tx.Exec(context.Background(), `ALTER TABLE operation_log RENAME TO operation_log_old`); err != nil {
		return fmt.Errorf("rename operation_log: %w", err)
	}
	if _, err := tx.Exec(context.Background(), ddl[0]); err != nil {
		return fmt.Errorf("create operation_log %s: %w", label, err)
	}
	if _, err := tx.Exec(context.Background(),
		`INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, detail, created_at)
		 SELECT id, event, actor_id, actor_name, record_id, detail, created_at FROM operation_log_old`,
	); err != nil {
		return fmt.Errorf("migrate operation_log rows: %w", err)
	}
	if _, err := tx.Exec(context.Background(), `DROP TABLE operation_log_old`); err != nil {
		return fmt.Errorf("drop operation_log_old: %w", err)
	}
	if _, err := tx.Exec(context.Background(), ddl[1]); err != nil {
		return fmt.Errorf("create operation_log index: %w", err)
	}
	return nil
}

// operationLogCorrelationDDL (0041 · VP-012 R1): keeps request correlation
// separate from the stable business-event detail JSON, preserving existing
// auth/audit detail contracts while allowing operations to be traced.
var operationLogCorrelationDDL = []string{
	`CREATE TABLE operation_log_correlation (
  operation_id   TEXT PRIMARY KEY REFERENCES operation_log(id) ON DELETE CASCADE,
  correlation_id TEXT NOT NULL
)`,
	`CREATE INDEX idx_operation_log_correlation_id ON operation_log_correlation(correlation_id)`,
}

// operationLogWalletJobsDDL (0043 · VP-012 R4): separates enqueue and
// terminal failure/cancellation from the existing successful reconcile event.
var operationLogWalletJobsDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login','wallet.account-create','wallet.account-update','wallet.adjust','wallet.freeze','wallet.unfreeze','wallet.reconcile','wallet.deduct-frozen','account.avatar-change','wallet.reconcile.queued','wallet.reconcile.failed','wallet.reconcile.cancelled')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}

// operationLogServiceCredentialsDDL (0045 · VP-012 R6): adds the service
// credential lifecycle events while preserving the correlation side table.
var operationLogServiceCredentialsDDL = []string{
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login','wallet.account-create','wallet.account-update','wallet.adjust','wallet.freeze','wallet.unfreeze','wallet.reconcile','wallet.deduct-frozen','account.avatar-change','wallet.reconcile.queued','wallet.reconcile.failed','wallet.reconcile.cancelled','service-credentials.create','service-credentials.use','service-credentials.revoke')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at INTEGER NOT NULL
)`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
}
