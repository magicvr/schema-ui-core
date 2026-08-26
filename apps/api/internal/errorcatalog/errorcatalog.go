// Package errorcatalog owns the frozen user-visible error catalog and locale
// negotiation (VP-007 S4 · I-L10N-004 path a, user-confirmed 2026-08-09).
//
// Bounded server-side negotiation: cataloged codes return a localized message
// for the request locale (zh-CN/en-US), plus a stable machine-readable
// messageKey. Uncataloged codes (INTERNAL etc.) stay English with no key and
// never leak diagnostics. The `error` code itself is never translated.
package errorcatalog

import (
	"net/http"
	"strings"
)

// Entry is one cataloged error code's localizable surface.
type Entry struct {
	MessageKey string
	En         string
	Zh         string
}

// Catalog covers every user-visible cataloged code (D-002 appendix A
// enumeration; the handler contract test pins the full set). `INTERNAL` is
// deliberately absent: it must never be localized or carry a messageKey.
var Catalog = map[string]Entry{
	"UNAUTHENTICATED":       {"error.unauth", "no active session", "未登录或会话已失效"},
	"UNAUTHORIZED":          {"error.unauthorized", "invalid username or password", "用户名或密码错误"},
	"REFRESH_TOKEN_EXPIRED": {"error.refreshTokenExpired", "sign-in expired; please sign in again", "登录已过期，请重新登录"},
	"NOT_FOUND":             {"error.notFound", "not found", "未找到"},
	"METHOD_NOT_ALLOWED":    {"error.methodNotAllowed", "method not allowed", "不允许的请求方法"},
	"FORBIDDEN":             {"error.forbidden", "you do not have permission for this action", "您没有执行此操作的权限"},
	"MUST_CHANGE_PASSWORD":  {"error.mustChangePassword", "password change required before continuing", "请先修改初始密码后再继续"},
	"SERVICE_MAINTENANCE":   {"error.serviceMaintenance", "service is under maintenance", "服务正在维护中"},
	"SERVICE_DEGRADED":      {"error.serviceDegraded", "service is operating in degraded mode", "服务当前处于降级模式"},
	"SERVICE_READ_ONLY":     {"error.serviceReadOnly", "service is read-only", "服务当前为只读模式"},

	"INVALID_LOGIN_BODY":     {"error.invalidLoginBody", "body must be JSON with username and password", "请求体必须是包含用户名和密码的 JSON"},
	"LOGIN_FAILED":           {"error.loginFailed", "authentication unavailable", "认证服务暂不可用"},
	"INVALID_REFRESH_BODY":   {"error.invalidRefreshBody", "body must be JSON with refreshToken", "请求体必须是包含 refreshToken 的 JSON"},
	"REFRESH_FAILED":         {"error.refreshFailed", "refresh unavailable", "刷新服务暂不可用"},
	"INVALID_LOGOUT_BODY":    {"error.invalidLogoutBody", "body must be JSON with refreshToken", "请求体必须是包含 refreshToken 的 JSON"},
	"LOGOUT_FAILED":          {"error.logoutFailed", "logout unavailable", "退出服务暂不可用"},
	"INVALID_PATCH_BODY":     {"error.invalidPatchBody", "body must be JSON", "请求体必须是 JSON"},
	"SCHEMA_NOT_FOUND":       {"error.schemaNotFound", "no page document for that pageId", "该 pageId 没有对应的页面文档"},
	"SETTINGS_NOT_FOUND":     {"error.settingsNotFound", "no settings with that id", "该 id 没有对应的设置"},
	"INVALID_SITE_TITLE":     {"error.invalidSiteTitle", "siteTitle must not be empty", "站点标题不能为空"},
	"INVALID_LOGO_URL":       {"error.invalidLogoUrl", "logoUrl fields must be empty, a same-origin path, or an http(s) URL", "Logo URL 必须为空、同源路径或 http(s) URL"},
	"INVALID_DEFAULT_LOCALE": {"error.invalidDefaultLocale", "defaultLocale must be auto, zh-CN or en-US", "默认语种必须是 auto、zh-CN 或 en-US"},
	"INVALID_DEFAULT_THEME":  {"error.invalidDefaultTheme", "defaultTheme must be auto, light or dark", "默认主题必须是 auto、light 或 dark"},
	"INVALID_TIMEZONE":          {"error.invalidTimezone", "siteTimezone must be auto or a valid IANA timezone", "默认时区必须是 auto 或有效的 IANA 时区"},
	"INVALID_DEFAULT_CURRENCY":  {"error.invalidDefaultCurrency", "defaultCurrency must be empty or a valid ISO 4217 code", "默认货币必须是空或有效的 ISO 4217 代码"},
	"INVALID_RETENTION_DAYS":    {"error.invalidRetentionDays", "operationLogRetentionDays must be between 1 and 3650", "审计日志保留天数必须是 1 到 3650 之间的整数"},
	"INVALID_EXPIRATION_ACTION": {"error.invalidExpirationAction", "operationLogExpirationAction must be archive or delete", "过期处理必须是归档或删除"},

	"INVALID_SORT_FIELD":    {"error.invalidSortField", "unsupported sort field", "不支持的排序字段"},
	"INVALID_SORT_ORDER":    {"error.invalidSortOrder", "order must be asc or desc", "排序方向必须是 asc 或 desc"},
	"INVALID_PAGE":          {"error.invalidPage", "page must be a positive integer", "页码必须是正整数"},
	"INVALID_PAGE_SIZE":     {"error.invalidPageSize", "pageSize must be a positive integer not exceeding 100", "每页条数必须是 1–100 的整数"},
	"INVALID_DATE_FILTER":   {"error.invalidDateFilter", "from/to must be YYYY-MM-DD or RFC3339", "from/to 必须是 YYYY-MM-DD 或 RFC3339"},
	"INVALID_CREATE_BODY":   {"error.invalidCreateBody", "body must be JSON", "请求体必须是 JSON"},
	"INVALID_CREATE_FIELD":  {"error.invalidCreateField", "invalid create field", "创建字段无效"},
	"INVALID_PATCH_FIELD":   {"error.invalidPatchField", "invalid patch field", "更新字段无效"},
	"INVALID_BODY":          {"error.invalidBody", "expected a JSON object with an ids array", "请求体应为包含 ids 数组的 JSON 对象"},
	"EMPTY_SELECTION":       {"error.emptySelection", "ids must contain at least one key", "至少选择一个条目"},
	"INVALID_SELECTION_KEY": {"error.invalidSelectionKey", "ids entries must be scalar keys", "ids 条目必须是标量键"},

	"USERNAME_TAKEN":         {"error.usernameTaken", "username already exists", "用户名已存在"},
	"ROLE_KEY_TAKEN":         {"error.roleKeyTaken", "role key already exists", "角色键已存在"},
	"ROLE_IN_USE":            {"error.roleInUse", "role is assigned to users", "该角色仍被用户使用"},
	"ROLE_SYSTEM":            {"error.roleSystem", "system roles cannot be modified", "系统角色不可修改"},
	"INVALID_ROLE_KEY":       {"error.invalidRoleKey", "invalid role key format", "角色键格式无效"},
	"INVALID_PERMISSION_REF": {"error.invalidPermissionRef", "permissions contain an unknown key", "权限包含未知键"},
	"ROLE_GRANT_FORBIDDEN":   {"error.roleGrantForbidden", "you may not grant roles you do not hold", "不能授予您本身不持有的角色"},
	"ROLE_NOT_FOUND":         {"error.roleNotFound", "no role with that key", "没有该键对应的角色"},
	"USER_NOT_FOUND":         {"error.userNotFound", "no user with that id", "没有该 id 对应的用户"},
	"CATALOG_NOT_FOUND":      {"error.catalogNotFound", "no catalog with that id", "没有该 id 对应的目录"},

	// F-006 补齐实际发出的 domain 错误码（S1 台账）：此前缺失导致 zh-CN 下英文回退。
	"LAST_ADMIN":                {"error.lastAdmin", "cannot remove the last admin user", "不能移除最后一个管理员"},
	"SELF_OPERATION":            {"error.selfOperation", "self operation is not allowed", "不允许对自身账号执行该操作"},
	"INVALID_ROLE_REF":          {"error.invalidRoleRef", "roles contain an unknown role key", "角色引用未知或无效"},
	"ROLE_ASSIGNMENT_FORBIDDEN": {"error.roleAssignmentForbidden", "you may not assign roles you do not hold", "不能分配您本身不持有的角色"},
	"ADMIN_ACCOUNT_FORBIDDEN":   {"error.adminAccountForbidden", "only an admin may manage admin accounts", "只有管理员可以管理管理员账号"},
	"INVALID_MENU_ITEM_REF":     {"error.invalidMenuItemRef", "menuItems contain an unknown id", "导航引用了未知的菜单项"},

	"INVALID_UPLOAD":  {"error.invalidUpload", "expected a multipart file part named file", "请求应为包含名为 file 的 multipart 文件"},
	"FILE_TOO_LARGE":  {"error.fileTooLarge", "file exceeds the size limit", "文件超过大小限制"},
	"FILE_NOT_FOUND":  {"error.fileNotFound", "file not found", "文件不存在"},
	"ASSET_NOT_FOUND": {"error.assetNotFound", "no brand asset with that id", "没有该 id 对应的品牌图片"},
	"INVALID_KIND":    {"error.invalidKind", "kind must be logo or favicon", "kind 必须是 logo 或 favicon"},

	// S-02 (GOAL-007 D-002 §4): file-library codes.
	"INVALID_UPLOAD_BODY": {"error.invalidUploadBody", "body must be JSON with a file field", "请求体必须是包含 file 字段的 JSON"},
	"INVALID_FILE_ID":     {"error.invalidFileId", "invalid file id", "文件 ID 无效"},

	// S-04 (GOAL-010 D-002 §4): scheduled-task codes.
	"TASK_NOT_FOUND":     {"error.taskNotFound", "no scheduled task with that id", "没有该 id 对应的定时任务"},
	"TASK_RUN_NOT_FOUND": {"error.taskRunNotFound", "no task run with that id", "没有该 id 对应的运行记录"},
	"TASK_KEY_TAKEN":     {"error.taskKeyTaken", "a scheduled task with that key already exists", "该键已存在对应的定时任务"},
	"INVALID_CRON":       {"error.invalidCron", "invalid cron expression", "cron 表达式无效"},
	"INVALID_HANDLER":    {"error.invalidHandler", "unknown task handler", "未知的任务处理器"},

	// S-01 (GOAL-008 D-002 §3): dictionary codes.
	"DICT_TYPE_NOT_FOUND":   {"error.dictTypeNotFound", "no dict type with that id", "没有该 id 对应的字典类型"},
	"DICT_ENTRY_NOT_FOUND":  {"error.dictEntryNotFound", "no dict entry with that id", "没有该 id 对应的字典条目"},
	"DICT_TYPE_KEY_TAKEN":   {"error.dictTypeKeyTaken", "a dict type with that key already exists", "该键已存在对应的字典类型"},
	"DICT_ENTRY_KEY_TAKEN":  {"error.dictEntryKeyTaken", "an entry with that key already exists in the dict type", "该字典类型下已存在相同键的条目"},
	"DICT_KEY_NOT_FOUND":    {"error.dictKeyNotFound", "no dict type with that key", "没有该键对应的字典类型"},
	"INVALID_FILE":          {"error.invalidFile", "file part is invalid", "文件内容无效"},
	"UNSUPPORTED_FILE_TYPE": {"error.unsupportedFileType", "file type is not allowed", "不允许的文件类型"},
	"STORAGE_UNAVAILABLE":   {"error.storageUnavailable", "storage is temporarily unavailable", "存储服务暂不可用"},
	"RATE_LIMITED":          {"error.rateLimited", "too many failed login attempts; try again later", "登录失败次数过多，请稍后重试"},
	// GOAL-004 S4-6: account lock terminal (423) — distinct from the per-IP
	// rate limit: this is a per-account state with a lock window.
	"ACCOUNT_LOCKED":        {"error.accountLocked", "account is temporarily locked; try again later", "账号已暂时锁定，请稍后重试"},
	"UPLOAD_QUOTA_EXCEEDED": {"error.uploadQuotaExceeded", "upload rejected: per-user quota exceeded", "上传被拒绝：超出每用户配额"},
	// W7 F-004: account avatar per-user upload quota.
	"AVATAR_QUOTA_EXCEEDED": {"error.avatarQuotaExceeded", "avatar upload quota exceeded; assign or clear your current avatar first", "头像上传超过配额；请先设置或清除当前头像"},

	// F-03 (GOAL-005 D-002 §3/§4): disabled terminal + self-service account codes.
	"ACCOUNT_DISABLED":      {"error.accountDisabled", "account is disabled; contact an administrator", "账号已被停用，请联系管理员"},
	"INVALID_PASSWORD":      {"error.invalidPassword", "current password is incorrect or the new password is invalid", "当前密码错误或新密码无效"},
	"INVALID_PASSWORD_BODY": {"error.invalidPasswordBody", "body must be JSON with currentPassword and newPassword", "请求体必须是包含 currentPassword 和 newPassword 的 JSON"},
	"SESSION_NOT_FOUND":     {"error.sessionNotFound", "no session with that id", "没有该 id 对应的会话"},
	"INVALID_STATUS_FILTER": {"error.invalidStatusFilter", "status must be active, revoked, or empty", "status 只能是 active、revoked 或空"},

	// F-02 (GOAL-004 D-002 §3/§4): data-transfer codes.
	"RESOURCE_NOT_FOUND":   {"error.resourceNotFound", "no transfer surface for that resource", "该资源没有对应的传输面"},
	"INVALID_CSV":          {"error.invalidCsv", "could not parse CSV", "CSV 解析失败"},
	"INVALID_IMPORT_BODY":  {"error.invalidImportBody", "body must be JSON with fileId", "请求体必须是包含 fileId 的 JSON"},
	"INVALID_EXPORT_LIMIT": {"error.invalidExportLimit", "pageSize must not exceed 10000", "导出每页条数不能超过 10000"},

	// F-04 (GOAL-006 D-002 §4): notification codes.
	"INVALID_SETTINGS_BODY":  {"error.invalidSettingsBody", "body must be JSON with enabled", "请求体必须是包含 enabled 的 JSON"},
	"NOTIFICATION_NOT_FOUND": {"error.notificationNotFound", "no notification with that id", "没有该 id 对应的通知"},

	// S-11 (GOAL-011 D-002 §2): login captcha code.
	"INVALID_CAPTCHA": {"error.invalidCaptcha", "captcha verification failed", "验证码校验失败"},

	// S-10 (GOAL-017 D-002 §3/§4): MFA codes.
	"INVALID_MFA_BODY":    {"error.invalidMfaBody", "body must be JSON with proof and code", "请求体必须是包含 proof 和 code 的 JSON"},
	// W7 F-007: MFA enrollment step-up requires the current password.
	"MFA_CURRENT_PASSWORD_REQUIRED": {"error.mfaCurrentPasswordRequired", "currentPassword is required to start MFA enrollment", "启用 MFA 需要验证当前密码"},
	"MFA_INVALID":         {"error.mfaInvalid", "invalid second-factor code", "第二因素验证码无效"},
	"MFA_PROOF_EXPIRED":   {"error.mfaProofExpired", "second-factor proof expired; sign in again", "第二因素验证已过期，请重新登录"},
	"MFA_PROOF_EXHAUSTED": {"error.mfaProofExhausted", "too many failed attempts; sign in again", "失败次数过多，请重新登录"},
	"MFA_NOT_ENROLLED":    {"error.mfaNotEnrolled", "no MFA enrollment for this account", "该账号未启用 MFA"},
	"MFA_PENDING_ONLY":    {"error.mfaPendingOnly", "MFA is not activated yet", "MFA 尚未激活"},
	"MFA_ALREADY_ACTIVE":  {"error.mfaAlreadyActive", "MFA is already active", "MFA 已激活"},

	// S-09 (GOAL-016 D-002 §3): data-permission codes.
	"INVALID_SCOPE":         {"error.invalidScope", "scope must be all or self", "数据范围必须是 all 或 self"},
	"INVALID_SCOPE_BODY":    {"error.invalidScopeBody", "body must be JSON with userId and scopes", "请求体必须是包含 userId 和 scopes 的 JSON"},
	"SCOPE_NOT_ENFORCEABLE": {"error.scopeNotEnforceable", "resource is not wired for row-level scoping", "该资源未接入行级数据范围"},

	// S-14 (GOAL-019 D-002 §3): wallet codes.
	"INVALID_WALLET_BODY":         {"error.invalidWalletBody", "body must be JSON with ownerType and ownerId", "请求体必须是包含 ownerType 和 ownerId 的 JSON"},
	"INVALID_WALLET_OWNER":        {"error.invalidWalletOwner", "ownerId is required", "缺少 ownerId"},
	"INVALID_WALLET_ACCOUNT":      {"error.invalidWalletAccount", "accountId is required", "缺少 accountId"},
	"INVALID_WALLET_STATUS":       {"error.invalidWalletStatus", "status must be active or disabled", "状态必须是 active 或 disabled"},
	"WALLET_NOT_FOUND":            {"error.walletNotFound", "wallet account not found", "钱包账户不存在"},
	"JOB_NOT_FOUND":               {"error.jobNotFound", "job not found", "任务不存在"},
	"JOB_NOT_CANCELLABLE":         {"error.jobNotCancellable", "job cannot be cancelled", "任务不可取消"},
	"JOB_NOT_RETRYABLE":           {"error.jobNotRetryable", "job cannot be retried", "任务不可重试"},
	"JOB_RESULT_NOT_READY":        {"error.jobResultNotReady", "job result is not ready", "任务结果尚未就绪"},
	"JOB_RESULT_EXPIRED":          {"error.jobResultExpired", "job result has expired", "任务结果已过期"},
	"JOB_ATTEMPTS_EXHAUSTED":      {"error.jobAttemptsExhausted", "job attempts exhausted", "任务重试次数已耗尽"},
	"JOB_HANDLER_FAILED":          {"error.jobHandlerFailed", "job handler failed", "任务执行失败"},
	"WALLET_OWNER_TAKEN":          {"error.walletOwnerTaken", "an account for that owner already exists", "该持有方已存在钱包账户"},
	"WALLET_DISABLED":             {"error.walletDisabled", "wallet account is disabled", "钱包账户已停用"},
	"INSUFFICIENT_BALANCE":        {"error.insufficientBalance", "insufficient balance for this mutation", "余额不足，无法执行该变动"},
	"LEDGER_VERSION_CONFLICT":     {"error.ledgerVersionConflict", "the account changed concurrently; reload and retry", "账户已被并发修改，请刷新后重试"},
	"LEDGER_IDEMPOTENCY_CONFLICT": {"error.ledgerIdempotencyConflict", "idempotency key was already used with a different payload", "幂等键已被不同载荷使用"},
	"INVALID_LEDGER_ENTRY":        {"error.invalidLedgerEntry", "invalid ledger entry", "非法账本流水"},
	"PRECONDITION_REQUIRED":       {"error.preconditionRequired", "provide If-Match or expectedVersion", "请提供 If-Match 或 expectedVersion"},
	"INVALID_PRECONDITION":        {"error.invalidPrecondition", "version preconditions must be valid and agree", "版本前置条件必须有效且一致"},
	// GOAL-020 D-001 §2: user wallet accounts are auto-created (get-or-create).
	"WALLET_USER_AUTO_ONLY": {"error.walletUserAutoOnly", "user wallet accounts are created automatically", "用户钱包账户由系统自动创建"},

	// W14 F-06 (GOAL-019): operations detail not found (previously uncataloged).
	"OPERATION_NOT_FOUND": {"error.operationNotFound", "no operation with that id", "没有该 id 对应的操作日志"},

	// S-12 (GOAL-012 D-002 §5): recycle-bin codes.
	"RECYCLE_ITEM_NOT_FOUND":        {"error.recycleItemNotFound", "no recycle item with that id", "没有该 id 对应的回收站记录"},
	"RECYCLE_RESTORE_CONFLICT":      {"error.recycleRestoreConflict", "a row with that key already exists; resolve the conflict and retry", "存在相同键的行，解决冲突后重试"},
	"RECYCLE_ITEM_ALREADY_RESTORED": {"error.recycleItemAlreadyRestored", "recycle item is already restored", "回收站记录已恢复"},

	// VP-017 R7 (GOAL-008 D-001): outbound-mail admin surface codes.
	"INVALID_MAIL_CONFIG":   {"error.invalidMailConfig", "invalid outbound-mail configuration", "出站邮件配置无效"},
	"MAIL_SWITCH_REJECTED":  {"error.mailSwitchRejected", "new channel configuration failed validation; the previous channel keeps serving", "新渠道配置校验未通过，继续沿用原渠道"},
	"MAIL_SEND_FAILED":      {"error.mailSendFailed", "the test message could not be sent", "测试邮件发送失败"},

	// workspace-018 R3 (GOAL-004 D-001 §3): account email identity codes.
	"EMAIL_INVALID":         {"error.emailInvalid", "invalid email address", "邮箱地址无效"},
	"EMAIL_TAKEN":           {"error.emailTaken", "email already bound or pending on another account", "该邮箱已被其他账号绑定或待校验"},
	"EMAIL_NOT_PENDING":     {"error.emailNotPending", "no pending email verification for this account", "该账号没有待校验的邮箱"},
	"EMAIL_CODE_INVALID":    {"error.emailCodeInvalid", "verification code is invalid", "验证码无效"},
	"EMAIL_CODE_EXPIRED":    {"error.emailCodeExpired", "verification code expired; request a new one", "验证码已过期，请重新获取"},
	"EMAIL_RESEND_COOLDOWN": {"error.emailResendCooldown", "please wait before requesting another code", "请求过于频繁，请稍后再试"},
	"EMAIL_SEND_FAILED":     {"error.emailSendFailed", "the verification email could not be sent", "验证邮件发送失败"},

	// workspace-019 R2 (GOAL-003 D-001 §6): self-recovery codes. Unknown
	// account / no challenge / wrong code deliberately share ONE code so the
	// pre-auth surface stays enumeration-neutral.
	"INVALID_RECOVERY_BODY":            {"error.invalidRecoveryBody", "body must be JSON with account, code and newPassword", "请求体必须是包含 account、code 和 newPassword 的 JSON"},
	"RECOVERY_CODE_INVALID":            {"error.recoveryCodeInvalid", "recovery code is invalid", "恢复码无效"},
	"RECOVERY_CODE_EXPIRED":            {"error.recoveryCodeExpired", "recovery code expired; request a new one", "恢复码已过期，请重新获取"},
	"RECOVERY_SECOND_FACTOR_REQUIRED":  {"error.recoverySecondFactorRequired", "your second factor is required to finish recovery", "完成恢复需要第二因素验证码"},

	// workspace-019 R3 (GOAL-004 D-001 §3): invitation codes. Unknown /
	// expired / consumed / revoked share ONE code on the pre-auth surface.
	"INVALID_INVITE_BODY": {"error.invalidInviteBody", "body must be JSON with token, username and password", "请求体必须是包含 token、username 和 password 的 JSON"},
	"INVITE_INVALID":      {"error.inviteInvalid", "invitation is unknown, expired, already used or revoked", "邀请无效：不存在、已过期、已使用或已撤销"},
	"INVITE_ROLE_GONE":    {"error.inviteRoleGone", "invited roles changed; ask for a new invitation", "邀请中的角色已变更，请索取新邀请"},
}

// SupportedLocales are the negotiation targets in preference order.
var SupportedLocales = []string{"zh-CN", "en-US"}

// NegotiateLocale picks the first supported locale from an Accept-Language
// header value; falls back to en-US when nothing matches.
func NegotiateLocale(header string) string {
	if strings.TrimSpace(header) == "" {
		return "en-US"
	}
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		tag := part
		if idx := strings.IndexAny(part, ";"); idx >= 0 {
			tag = strings.TrimSpace(part[:idx])
		}
		lower := strings.ToLower(tag)
		if lower == "zh-cn" || lower == "zh" {
			return "zh-CN"
		}
		if lower == "en-us" || lower == "en" || strings.HasPrefix(lower, "en-") {
			return "en-US"
		}
	}
	return "en-US"
}

// Body builds the wire envelope for a code in the given locale. Returns
// (body, contentLanguage, ok); ok=false when the code is uncataloged (the
// caller keeps the English generic message and omits messageKey).
//
// GOAL-014 D-002 §2: for field-validation codes (INVALID_CREATE_FIELD /
// INVALID_PATCH_FIELD) the caller's message carries the concrete field
// reason (e.g. "name must not be empty"); it is appended after the cataloged
// text so the user sees both the localized code surface and the specific
// field problem (previously the catalog text replaced it entirely).
func Body(code, message, locale string) (map[string]any, string, bool) {
	entry, found := Catalog[code]
	if !found {
		return nil, locale, false
	}
	text := entry.En
	if locale == "zh-CN" {
		text = entry.Zh
	}
	body := map[string]any{
		"error":      code,
		"message":    text,
		"messageKey": entry.MessageKey,
	}
	// Field-validation codes append the concrete field reason (GOAL-014
	// D-002 §2); every other cataloged code keeps the pure localized text so
	// caller-supplied diagnostics never leak into the response.
	if isFieldValidationCode(code) {
		if msg := strings.TrimSpace(message); msg != "" {
			body["message"] = text + ": " + msg
		}
	}
	return body, locale, true
}

// isFieldValidationCode reports codes whose message carries a field-level
// reason that should be surfaced next to the localized catalog text.
func isFieldValidationCode(code string) bool {
	return code == "INVALID_CREATE_FIELD" || code == "INVALID_PATCH_FIELD"
}

// FieldError is one field-level validation failure (GOAL-014 D-002 §2.1).
type FieldError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// BodyWithFields is Body plus an optional fieldErrors array. It is the same
// envelope with an extra optional key, so old consumers that ignore unknown
// keys keep working unchanged.
func BodyWithFields(code, message, locale string, fields []FieldError) (map[string]any, string, bool) {
	body, lang, ok := Body(code, message, locale)
	if !ok || len(fields) == 0 {
		return body, lang, ok
	}
	body["fieldErrors"] = fields
	return body, lang, true
}

// Negotiate is a convenience for http handlers: header from the request.
func Negotiate(r *http.Request) string {
	if r == nil {
		return "en-US"
	}
	return NegotiateLocale(r.Header.Get("Accept-Language"))
}
