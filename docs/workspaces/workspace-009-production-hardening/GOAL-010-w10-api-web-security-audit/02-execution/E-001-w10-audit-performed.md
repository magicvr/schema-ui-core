---
id: E-001-w10-audit-performed
goal: GOAL-010-w10-api-web-security-audit
status: done
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-010-w10-api-web-security-audit
version: 0.1.0
---

# E-001 · W10 独立审计执行与报告落盘（2026-08-21）

## 事实

1. 用户发起独立审计指令："独立审计 apps/api 和 apps/web 的代码实现是否存在 bug 或安全漏洞"。
2. 审计执行方式（同日完成）：
   - 主线（DSH 会话模型）深读安全关键路径：auth（auth.go:665 行）、config（config.go:1005 行）、session（session.go:73 行）、composition（composition.go:694 行）、handler/auth.go、handler/account.go、handler/account_avatar.go、handler/account_self.go、handler/captcha.go、handler/mfa.go、handler/export.go、handler/import.go、handler/filelibrary.go、handler/schema.go、handler/bootstrap.go、handler/rate_limit.go、handler/route_envelope.go、handler/operational.go、handler/health.go、handler/notifications.go、handler/service_credentials.go、handler/branding_assets.go、handler/raster_assets.go、handler/datapermission.go、handler/manifest.go、store/open.go、store/store.go、errorcatalog/errorcatalog.go、cmd/server/main.go、configs/config.yaml、configs/env.example、Dockerfile、compose.yaml。
   - 2 个并行子代理广度审计：① api（handlers + store + modules + auth + config + composition 全覆盖）；② web（account/auth-client.ts、tokens.ts、App.tsx、host/boot.ts、renderer/render.tsx、renderer/form-controls.tsx、renderer/reaction-engine.ts、protocol/load-page.ts、nginx.conf、vite.config.ts、package.json、Dockerfile）。
   - 交叉验证：所有 P1/P2 结论由主线逐条重读源码确认。
3. 审计结论：P0=0；P1=1（env.example 硬编码真实数据库凭据）；P2=6（CSRF 架构姿势、window.open 无 noopener、fetch 无超时、文件下载消毒边界、刷新令牌并发原子性 PG、服务凭据作用域无上限）；P3=5（informational/已接受取舍）。
4. 报告落盘：A-001 摘要+findings（`03-audit/A-001-w10-independent.md`）+ 全文附件（`attachments/audit-A-001-w10-full-report.md`）。
5. 目标五件套 + 三个 ledger 目录 + attachments/ 一次建齐；goal-tree 与 workspace.md 波次表已同步。

## 产物

- [03-audit/A-001-w10-independent.md](../03-audit/A-001-w10-independent.md)
- [attachments/audit-A-001-w10-full-report.md](../attachments/audit-A-001-w10-full-report.md)
- [goal-tree.md](../../goal-tree.md)（树+表新增 GOAL-010）
- [workspace.md](../../workspace.md)（波次表新增 W10 行）

## 后续（计划，非事实）

- S2：用户裁决 required 范围与 go 宣称影响（I-002）。
- S3/S4：修复实施与复核（未开始）。