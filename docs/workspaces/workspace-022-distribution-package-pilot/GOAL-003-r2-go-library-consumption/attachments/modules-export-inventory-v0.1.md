# B 层模块包导出符号盘点 v0.1（2026-08-29 · F-002 回填）

> 来源：pps/api/modules/* 顶层包非测试文件导出扫描（^(func|type|const|var) [A-Z]）。
> 目的：契约冻结面 §2 的 B 层全量符号清单（构造签名面）；构造参数含 C 层类型者标注 ⚠️。

## account

- `const ModuleID = "admin.account"`
- `type Provider struct`
- `func New(a *auth.Authenticator, repository *authsession.Repository, operations operationlog.Recorder…`

## activity

- `const ModuleID = "admin.activity"`
- `type Provider struct`
- `func New(a *auth.Authenticator, operations operationlog.Reader) *Provider`

## authsession

- `var roleKeyRe = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)`
- `func userWithRoles(tx kernel.Tx, row kernel.Row) (*User, error)`
- `func rolesForUser(tx kernel.Tx, userID string) ([]string, error)`
- `const globalFailureWindow = 24 * time.Hour`
- `func scanUser(row interface{ Scan(...any) error }) (*User, error)`
- `func ensureRole(tx kernel.Tx, key string, now int64) error`
- `func linkUserRole(tx kernel.Tx, userID, key string, now int64) error`
- `const legacyLoginSource = "-"`
- `const lockCounterWindow = 15 * time.Minute`
- `var errSourceInsertRace = errors.New("authsession: login failure insert race")`
- `func normalizeLoginSource(ip string) string`
- `var ErrSessionNotFound = errors.New("authsession: session not found")`
- `func userHasRoleKey(tx kernel.Tx, userID, key string) (bool, error)`
- `func boolInt(value bool) int`
- `func countEnabledAdminUsersExcluding(tx kernel.Tx, id string) (int, error)`
- `func countEnabledAdminUsers(tx kernel.Tx) (int, error)`
- `func normalizeEmailInput(raw string) (string, error)`
- `func stringPtr(s string) *string { return &s }`
- `func nullIfNil(s *string) any`
- `func hashCode(emailCode string) string`
- `func generateEmailCode() (string, error)`
- `const emailVerificationMailSubject = "账号邮箱验证码 · Email verification code"`
- `func emailVerificationMailBody(code string, expires time.Time) string`
- `func sendVerificationMail(ctx context.Context, sender kernel.MailSender, to, code string, expires ti…`
- `type verificationOutcome int`
- `var errNotSentinel = errors.New("authsession: controlled verification outcome")`
- `type Invite struct`
- `func newInviteToken() (raw, hash string, err error)`
- `func hashInviteToken(raw string) string`
- `func scanInvite(row kernel.Row) (*Invite, error)`
- `const inviteColumns = `id, roles, invited_by, email, expires_at, consumed_at, revoked_at, last_sent_…`
- `func getInviteRow(tx kernel.Tx, id string) (kernel.Row, error)`
- `type InviteStatusFilter string`
- `func ParseInviteStatus(raw string) InviteStatusFilter`
- `func inviteSortSQL(sort, order string) string`
- `func inviteWhereQ(q string) (string, []any)`
- `const maxNotificationsPerUser = 500`
- `type Notification struct`
- `type NotificationFilter struct`
- `func scanNotification(row interface{ Scan(...any) error }) (*Notification, error)`
- `var ErrPasswordPolicyViolation = fmt.Errorf("authsession: password violates the active policy")`
- `var ErrPasswordPolicyNotSeeded = errors.New("authsession: password policy row not seeded")`
- `type PasswordPolicy struct`
- `func randomHexID(prefix string) string`
- `func countCategories(plain string) int`
- `const recoveryMailSubject = "密码重置验证码 · Password reset code"`
- `func recoveryMailBody(code string, expires time.Time) string`
- `func sendRecoveryMail(ctx context.Context, sender kernel.MailSender, to, code string, expires time.T…`
- `type RecoveryTarget struct`
- `type RecoveryOutcome int`
- `func trimIdentifier(raw string) string`
- `func trimCode(raw string) string { return strings.TrimSpace(raw) }`
- `type TxRunner interface`
- `type Repository struct`
- `func NewRepository(runner TxRunner) *Repository`
- `type User struct`
- `type RefreshToken struct`
- `type UserFilter struct`
- `type UserPatch struct`
- `type Role struct`
- `type RolePatch struct`
- `type RoleFilter struct`
- `type PermissionCatalogEntry struct`
- `type MenuItemCatalogEntry struct`
- `func dedupeKeys(keys []string) []string`
- `func sameRoleSet(a, b []string) bool`
- `func hydrateRole(tx kernel.Tx, role *Role) error`
- `func replaceRolePermissions(tx kernel.Tx, roleID string, keys []string) error`
- `func replaceRoleMenuItems(tx kernel.Tx, roleID string, ids []string) error`
- `func scanRoleRow(row interface{ Scan(...any) error }, role *Role) error`
- `func rolesWhere(query string, system *bool) (string, []any)`
- `func rolesSortSQL(sort, order string) string`
- `func menuItemLabel(pageRef string) string`
- `type ServiceCredential struct`
- `type ServiceCredentialAudit func(kernel.Tx) error`
- `type ServiceCredentialRevokeAudit func(kernel.Tx, ServiceCredential) error`
- `var ErrCredentialNameTaken = errors.New("authsession: service credential name already taken")`
- `func sortedScopes(scopes []string) []string`
- `func scanServiceCredential(row interface{ Scan(...any) error }) (*ServiceCredential, error)`
- `func countAdminUsersExcludingBatch(tx kernel.Tx, ids []string) (int, error)`
- `func countAdminUsersExcluding(tx kernel.Tx, id string) (int, error)`
- `func scanUserListRow(row interface{ Scan(...any) error }) (*User, error)`
- `func usersWhere(query string, enabled, locked *bool) (string, []any)`
- `func usersSortSQL(sort, order string) string`

## compiled

- `func PersistenceProviders() []kernel.Provider`
- `func PersistenceCatalog() ([]kernel.MigrationContribution, error)`

## dashboard

- `const ModuleID = "admin.dashboard"`
- `type Provider struct{}`
- `func New() *Provider`

## datadictionary

- `const ModuleID = "admin.data-dictionary"`
- `type Provider struct`
- `func New(a *auth.Authenticator, repository *datadictionarystore.Repository, operations operationlog.…`

## datapermission

- `const ModuleID = "admin.data-permission"`
- `type Service struct`
- `func NewService(repo *datapermissionstore.Repository, enforceable []string) *Service`
- `type Provider struct`
- `func New(a *auth.Authenticator, service *Service, operations operationlog.Recorder) *Provider`

## datatransfer

- `const ModuleID = "admin.data-transfer"`
- `type Provider struct`
- `func New(a *auth.Authenticator, repository *authsession.Repository, operations operationlog.Recorder…`

## filelibrary

- `const ModuleID = "admin.file-library"`
- `type Provider struct`
- `func New(a *auth.Authenticator, operations operationlog.Recorder, objects kernel.ObjectStore) *Provi…`

## logincaptcha

- `const challengeTTL = 5 * time.Minute`
- `var ErrInvalidCaptcha = errors.New("captcha verification failed")`
- `type Service struct`
- `func NewService(repository *store.Repository) *Service`
- `func answerHash(id, answer string) string`
- `func randInt(min, max int) int`
- `func newChallengeID() string`
- `const ModuleID = "admin.login-captcha"`
- `type Provider struct`
- `func New(a *auth.Authenticator, service *Service, operations operationlog.Recorder) *Provider`

## mfa

- `const ModuleID = "admin.mfa"`
- `type Provider struct`
- `func New(a *auth.Authenticator, service *Service, operations operationlog.Recorder, revoker handler.…`
- `type Service struct`
- `func NewService(repo *store.Repository, serverSecret, previousSecret []byte) *Service`
- `func randomRecoveryCode() string`
- `func encryptSecret(key, plain []byte) string`
- `func decryptWithKey(key []byte, encoded string) string`
- `func GenerateSecret() (string, error)`
- `func totpCode(secretBase32 string, step int64) (string, error)`
- `func ValidateTotp(secretBase32, code string, now time.Time, window int, lastUsedStep int64) (int64, …`
- `func otpauthURL(issuer, account, secret string) string`
- `func urlEscape(s string) string`

## notifications

- `const ModuleID = "admin.notifications"`
- `type Provider struct`
- `func New(a *auth.Authenticator, repository *authsession.Repository) *Provider`

## operationlog

- `const DetailSchemaVersion = 1`
- `const RedactedValue = "[REDACTED]"`
- `type DetailChange struct`
- `type DetailEnvelope struct`
- `func NewDetail(action string, before, after map[string]any) (string, error)`
- `func ParseDetail(raw string) (DetailEnvelope, error)`
- `func redactObject(input map[string]any) (map[string]any, error)`
- `func redactValue(value any, key string) (any, error)`
- `func isSensitiveKey(key string) bool`
- `type Operation struct`
- `type OperationFilter struct`
- `type Recorder interface`
- `type TransactionalRecorder interface`
- `type Reader interface`
- `type TxRunner interface`
- `var ErrNotFound = errors.New("operationlog: not found")`
- `type Repository struct`
- `func NewRepository(runner TxRunner) *Repository`
- `func scanOperation(row interface{ Scan(...any) error }) (Operation, error)`
- `func operationsWhere(filter OperationFilter) (string, []any)`
- `func operationsSortSQL(sort, order string) string`
- `type RetentionPolicy struct`
- `func StartRetentionSweep(repo *Repository, loadPolicy func() (RetentionPolicy, error), interval time…`

## recyclebin

- `const ModuleID = "admin.recycle-bin"`
- `type Provider struct`
- `func New(a *auth.Authenticator, service *Service, operations operationlog.Recorder) *Provider`
- `var ErrRestoreConflict = errors.New("recycle restore conflict")`
- `type Service struct`
- `func NewService(repository *recyclestore.Repository, dictionary *store.Repository, tasks *tasksstore…`
- `func toHandlerItem(item recyclestore.Item) handler.RecycleItem`
- `func isConflict(err error) bool`
- `func dictTypeFromPayload(payload map[string]any, now time.Time) store.DictType`
- `func dictEntryFromPayload(payload map[string]any, now time.Time) store.DictEntry`
- `func taskFromPayload(payload map[string]any, now time.Time) tasksstore.Task`
- `func stringField(m map[string]any, key string) string`
- `func boolField(m map[string]any, key string) bool`
- `func intField(m map[string]any, key string) int`
- `func timeField(m map[string]any, key string, fallback time.Time) time.Time`
- `func newSnapshotID() (string, error)`

## roles

- `const ModuleID = "admin.roles"`
- `type Provider struct`
- `func New(a *auth.Authenticator, repository *authsession.Repository, operations operationlog.Recorder…`

## scheduledtasks

- `const ModuleID = "admin.scheduled-tasks"`
- `type Provider struct`
- `func New(a *auth.Authenticator, repository *store.Repository, operations operationlog.Recorder, tras…`
- `type TaskHandler func(ctx context.Context, task store.Task, now time.Time) error`
- `const tickInterval = 30 * time.Second`
- `type Scheduler struct`
- `func NewScheduler(repository *store.Repository) *Scheduler`
- `var randReader = rand.Read`
- `var runIDSeq atomic.Uint64`
- `func newRunID() string`

## schemarender

- `const ModuleID = "core.schema-render"`
- `type Provider struct{}`
- `func New() *Provider { return &Provider{} }`

## settings

- `const ModuleID = migration.ModuleID`
- `type Provider struct`
- `func New(a *auth.Authenticator, repository *settingsrepository.Repository, operations operationlog.R…`

## systemmonitoring

- `const ModuleID = "admin.system-monitoring"`
- `type Provider struct`
- `func New(a *auth.Authenticator, st kernel.Store, plan kernel.Plan, ready func() bool, dbPath string,…`

## users

- `const ModuleID = "admin.users"`
- `type Provider struct`
- `func New(a *auth.Authenticator, repository *authsession.Repository, operations operationlog.Recorder…`

## wallet

- `const ReconcileJobKind = "wallet.reconcile"`
- `type reconcileJobPayload struct`
- `type JobService struct`
- `func NewJobService(service *Service, repository *jobs.Repository, runner *jobs.Runner, operations op…`
- `func reconciliationResult(run walletstore.ReconciliationRun) map[string]any`
- `const ModuleID = "admin.wallet"`
- `type Service struct`
- `func NewService(repo *walletstore.Repository) *Service`
- `var entryIDSeq atomic.Uint64`
- `func newID(now time.Time) (string, error)`
- `type Provider struct`
- `func New(a *auth.Authenticator, service *Service, jobs *JobService, operations operationlog.Recorder…`

