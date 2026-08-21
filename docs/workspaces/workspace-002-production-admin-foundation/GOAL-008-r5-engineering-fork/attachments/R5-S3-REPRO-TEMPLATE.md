---
title: R5-S3 · 独立复现记录模板（I-008-002 v0.1.3，含 W16-F01 首登改密）
status: template
doc_type: reproduction-record-template
created: 2026-08-20
parent: GOAL-008-r5-engineering-fork
version: 0.1.0
---

# R5-S3 · 独立复现记录模板（I-008-002 v0.1.3）

> 复制模板为 `R5-S3-REPRO-<NNN>.md` 后填写。字段要求见
> [`I-008-002-fork-reproduction-protocol.md`](I-008-002-fork-reproduction-protocol.md) v0.1.3 §4；
> 记录必须以**真实执行事实**为准，禁止把计划/假设写成已通过。

## 1. protocol

- 协议：`I-008-002` **v0.1.3**（本模板据此）。

## 2. attempt

- 编号：`R5-S3-REPRO-<NNN>`｜日期：`{{YYYY-MM-DD}}`（时区 `{{TZ}}`）
- 操作者/执行器标识：`{{operator}}`
- 独立性声明：`same-operator-clean-session`（隔离 shell/checkout/DB；未复用已启动服务）

## 3. source

- 仓库 URL / ref / commit：`{{repository-url}}` `{{ref}}`
- 工作树状态：`{{clean / diff-hash}}`（预 T0 `git status --short` 附件留痕）

## 4. path

- 路径：`local-dual-process` 或 `compose`
- API base URL：`{{API_BASE_URL}}`（compose 默认不发布宿主端口；smoke 时经 `scripts/pre-release-smoke.sh` loopback override `127.0.0.1:25080`；否则写实测 URL）
- Web base URL：`{{WEB_BASE_URL}}`（默认 `http://localhost:25081`；local `http://localhost:${WEB_PORT:-25173}`）

## 5. platform

- OS/架构、Git、Go、Node、npm、Docker/Compose、bash/curl 版本（照实填写）

## 6. cache precondition

- 已完成的依赖/镜像准备命令（`go mod download` / `npm ci` / 镜像拉取）与是否命中缓存；哪些耗时被排除。

## 7. W16-F01 · 首登强制改密（必填）

> v0.1.3：fresh seed 的 admin `must_change_password=1`，smoke/复现**必须**走真实改密流程，禁止清标志 / 改库 / 加 dev 跳过开关。

| 项 | 记录 |
|----|------|
| 初始登录 | username `admin` / 初始密码来源（`ADMIN_INITIAL_PASSWORD`，只写来源类别，不写值） |
| 强制改密页面/断言 | 首登后出现强制改密界面（en `Change your password` / zh `修改初始密码`），或 API 路径 `POST /api/account/password` |
| 新密码 | `SMOKE_PASSWORD_NEW`（默认 `<SMOKE_PASSWORD>-changed`；只写占位描述，不写真实值） |
| 改密方式 | 真实 UI 表单或 `POST /api/account/password`；改密后以返回的新 token/新密码继续 |
| 改密后验证 | 旧密码登录 `401`；新密码登录 `200` 且 `user.mustChangePassword=false`；业务 API 可用 |

## 8. timing

- 起点/终点 UTC 时间与单调计时秒数；四个终点各自时间戳（`<= 900s`）。

## 9. checks

| # | 检查 | 结果 | 证据 |
|---|------|------|------|
| 0 | **W16-F01 首登强制改密**（§7） | PASS/FAIL | 页面/接口输出路径 |
| 1 | `GET ${API_BASE_URL}/healthz` + `${API_BASE_URL}/readyz` → 200 `status=ok` | PASS/FAIL | run log |
| 2 | `${WEB_BASE_URL}/api/auth/login`（初始密码后改密，用新密码）→ 200 非空 token | PASS/FAIL | run log |
| 3 | `${WEB_BASE_URL}/api/accounts/me`（Bearer）→ 200 user+features | PASS/FAIL | run log |
| 4 | 浏览器 `${WEB_BASE_URL}/list-edit-lifecycle` → `List + edit lifecycle` + 列表 `Acme Console` 已加载 | PASS/FAIL | 截图/日志 |
| 5 | smoke 输出：`SM-001～SM-006`（含 disposable `SM-006=PASS`）；可选 `SM-008`（真实浏览器/CSP） | PASS/FAIL | logs |

## 10. secrets

- 配置来源类别与是否脱敏；禁止记录 token/password/secret 值。

## 11. result

- pass/fail、失败原因、重试编号、日志/截图/命令输出路径。
