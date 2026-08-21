---
id: E-006
doc: execution
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-006 · S3 追加回归

- **Go**：`go test ./...` 全量 **0 FAIL**（含 avatar 7 例、branding 重构回归、迁移台账 35/36、authsession/profile 用例）。
- **Web**：vitest **1029/1029**（63 文件；schema-keys 结构测试自动覆盖新 `schema.account.field.avatarUrl` 键）；`tsc -b` 0。
- **e2e（Playwright + 真实 Go API + SQLite 新鲜库）**：admin **8/8** + mvp **8/8**（各 1 跳过跨 profile 用例）。shell.spec 头像全链路通过：真实上传（multipart PNG）→ 200 + `/api/account/avatars/{id}` URL → PATCH profile 提交 → reload 后用户菜单触发器显示头像 img；移动端（390px）汉堡在最左（首元素）与品牌条断言通过。
- **e2e 首跑发现并修复**：`parseAuthUser` 未解析 `avatarUrl`（/me 快照字段被丢弃导致用户菜单不显示头像）——已修复并复跑通过。
