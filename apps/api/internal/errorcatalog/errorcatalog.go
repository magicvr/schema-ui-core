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
	"UNAUTHENTICATED": {"error.unauth", "no active session", "未登录或会话已失效"},
	"UNAUTHORIZED":    {"error.unauthorized", "invalid username or password", "用户名或密码错误"},
	"FORBIDDEN":       {"error.forbidden", "you do not have permission for this action", "您没有执行此操作的权限"},

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
	"INVALID_TIMEZONE":       {"error.invalidTimezone", "siteTimezone must be auto or a valid IANA timezone", "默认时区必须是 auto 或有效的 IANA 时区"},

	"INVALID_SORT_FIELD":     {"error.invalidSortField", "unsupported sort field", "不支持的排序字段"},
	"INVALID_SORT_ORDER":     {"error.invalidSortOrder", "order must be asc or desc", "排序方向必须是 asc 或 desc"},
	"INVALID_PAGE":           {"error.invalidPage", "page must be a positive integer", "页码必须是正整数"},
	"INVALID_PAGE_SIZE":      {"error.invalidPageSize", "pageSize must be a positive integer not exceeding 100", "每页条数必须是 1–100 的整数"},
	"INVALID_CREATE_BODY":    {"error.invalidCreateBody", "body must be JSON", "请求体必须是 JSON"},
	"INVALID_CREATE_FIELD":   {"error.invalidCreateField", "invalid create field", "创建字段无效"},
	"INVALID_PATCH_FIELD":    {"error.invalidPatchField", "invalid patch field", "更新字段无效"},
	"INVALID_BODY":           {"error.invalidBody", "expected a JSON object with an ids array", "请求体应为包含 ids 数组的 JSON 对象"},
	"EMPTY_SELECTION":        {"error.emptySelection", "ids must contain at least one key", "至少选择一个条目"},
	"INVALID_SELECTION_KEY":  {"error.invalidSelectionKey", "ids entries must be scalar keys", "ids 条目必须是标量键"},

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
	"LAST_ADMIN":             {"error.lastAdmin", "cannot remove the last admin user", "不能移除最后一个管理员"},
	"SELF_OPERATION":         {"error.selfOperation", "self operation is not allowed", "不允许对自身账号执行该操作"},
	"INVALID_ROLE_REF":       {"error.invalidRoleRef", "roles contain an unknown role key", "角色引用未知或无效"},
	"ROLE_ASSIGNMENT_FORBIDDEN": {"error.roleAssignmentForbidden", "you may not assign roles you do not hold", "不能分配您本身不持有的角色"},
	"ADMIN_ACCOUNT_FORBIDDEN":   {"error.adminAccountForbidden", "only an admin may manage admin accounts", "只有管理员可以管理管理员账号"},
	"INVALID_MENU_ITEM_REF":  {"error.invalidMenuItemRef", "menuItems contain an unknown id", "导航引用了未知的菜单项"},

	"INVALID_UPLOAD":        {"error.invalidUpload", "expected a multipart file part named file", "请求应为包含名为 file 的 multipart 文件"},
	"FILE_TOO_LARGE":        {"error.fileTooLarge", "file exceeds the size limit", "文件超过大小限制"},
	"FILE_NOT_FOUND":        {"error.fileNotFound", "file not found", "文件不存在"},

	// S-02 (GOAL-007 D-002 §4): file-library codes.
	"INVALID_UPLOAD_BODY":   {"error.invalidUploadBody", "body must be JSON with a file field", "请求体必须是包含 file 字段的 JSON"},
	"INVALID_FILE_ID":       {"error.invalidFileId", "invalid file id", "文件 ID 无效"},

	// S-04 (GOAL-010 D-002 §4): scheduled-task codes.
	"TASK_NOT_FOUND":      {"error.taskNotFound", "no scheduled task with that id", "没有该 id 对应的定时任务"},
	"TASK_RUN_NOT_FOUND":  {"error.taskRunNotFound", "no task run with that id", "没有该 id 对应的运行记录"},
	"TASK_KEY_TAKEN":      {"error.taskKeyTaken", "a scheduled task with that key already exists", "该键已存在对应的定时任务"},
	"INVALID_CRON":        {"error.invalidCron", "invalid cron expression", "cron 表达式无效"},

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

	// F-03 (GOAL-005 D-002 §3/§4): disabled terminal + self-service account codes.
	"ACCOUNT_DISABLED":      {"error.accountDisabled", "account is disabled; contact an administrator", "账号已被停用，请联系管理员"},
	"INVALID_PASSWORD":      {"error.invalidPassword", "current password is incorrect or the new password is invalid", "当前密码错误或新密码无效"},
	"INVALID_PASSWORD_BODY": {"error.invalidPasswordBody", "body must be JSON with currentPassword and newPassword", "请求体必须是包含 currentPassword 和 newPassword 的 JSON"},
	"SESSION_NOT_FOUND":     {"error.sessionNotFound", "no session with that id", "没有该 id 对应的会话"},

	// F-02 (GOAL-004 D-002 §3/§4): data-transfer codes.
	"RESOURCE_NOT_FOUND":  {"error.resourceNotFound", "no transfer surface for that resource", "该资源没有对应的传输面"},
	"INVALID_CSV":         {"error.invalidCsv", "could not parse CSV", "CSV 解析失败"},
	"INVALID_IMPORT_BODY": {"error.invalidImportBody", "body must be JSON with fileId", "请求体必须是包含 fileId 的 JSON"},
	"INVALID_EXPORT_LIMIT": {"error.invalidExportLimit", "pageSize must not exceed 10000", "导出每页条数不能超过 10000"},

	// F-04 (GOAL-006 D-002 §4): notification codes.
	"INVALID_SETTINGS_BODY":  {"error.invalidSettingsBody", "body must be JSON with enabled", "请求体必须是包含 enabled 的 JSON"},
	"NOTIFICATION_NOT_FOUND": {"error.notificationNotFound", "no notification with that id", "没有该 id 对应的通知"},
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
func Body(code, message, locale string) (map[string]any, string, bool) {
	entry, found := Catalog[code]
	if !found {
		return nil, locale, false
	}
	text := entry.En
	if locale == "zh-CN" {
		text = entry.Zh
	}
	return map[string]any{
		"error":      code,
		"message":    text,
		"messageKey": entry.MessageKey,
	}, locale, true
}

// Negotiate is a convenience for http handlers: header from the request.
func Negotiate(r *http.Request) string {
	if r == nil {
		return "en-US"
	}
	return NegotiateLocale(r.Header.Get("Accept-Language"))
}