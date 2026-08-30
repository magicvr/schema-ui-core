---
title: E-001 · W15 S3 API 修正实施与回归
status: active
created: 2026-08-30
updated: 2026-08-30
parent: null
version: 0.1.0
---

# E-001 · W15 S3 API 修正实施与回归（F-001～F-004）

日期：2026-08-30 · checkpoint：`609cd6d6`

## 实施事实

1. **F-001 · 默认暴露面收紧**（`server/config.go` + `server/config.default.yaml` + `cmd/schema-ui/templates/config.yaml.tmpl`）：
   - 默认监听 `:25080` → `127.0.0.1:25080`（代码默认、内嵌默认 YAML、create 骨架模板三处同步；对外暴露需显式 `-addr`/`http.addr`/`HTTP_ADDR`）。
   - `validate()` 增加空 `APP_ENV` fail-closed（内嵌默认显式 `development`，仅自定义 YAML 省略时触发；与主仓「refusing to guess」语义一致）。
2. **F-002 · JWT secret 强度门禁**（`internal/config/config.go` + `server/config.go`）：
   - 导出 `ValidateJWTSecretStrength(keyName, secret)`（≥32 字符 + 字母数字混合，错误不携带密钥值），`ValidateProd` 当前/previous 密钥分支复用（单一来源）。
   - `server.Config.validate()` 非 development 调用同一校验。
3. **F-003 · bootstrap 密码策略**（`modules/authsession/password_policy.go` + `cmd/server/main.go` + `server/serve.go`）：
   - 导出 `ValidateSeedPassword`（8–72 字节、非空，复用策略边界常量；fresh DB 的 0057 播种行默认即此）。
   - `cmd/server` `resolveSeedHash` 与 public serve `bootstrapAdmin` 均在非 development 先过策略再 bcrypt，不满足 fail-closed；development 保留 `admin` 回退。
4. **F-004 · MFA step-up 一次性语义**（`modules/mfa/service.go`）：
   - `requireActiveSecondFactor` TOTP 校验成功后以匹配 step 执行 `repo.AdvanceLastUsedStep`（CAS）；CAS 失败（同窗重放/并发双提交）→ `ErrMFAInvalid`；`maybeRewrap` 仅赢者执行。disable / recovery rotate / VerifySecondFactor 共用路径一并收紧。

## 新增测试

- `cmd/server/main_test.go`：resolveSeedHash 策略负例（production 缺省/短/7 字节/73 字节 → fail；compliant → pass；dev 回退保持）。
- `serve_test.go` `TestRunRejectsWeakSeedPasswordNonDev`：production + 4 字节种子 → bootstrap 启动失败（监听前）。
- `config_test.go`：默认回环断言、`TestLoadConfigRequiresExplicitAppEnv`（空 env fail-closed）、`TestLoadConfigJWTSecretStrengthNonDev`（短/纯字母/纯数字 secret 负例 + compliant 正例）；既有用例按新语义显式声明 `env: development`。
- `password_policy_test.go` `TestValidateSeedPassword`：7/8/72/73 字节与空白边界。
- `service_test.go` `TestServiceStepUpTotpReplayRejected`：同码二次 rotate 拒绝、VerifySecondFactor 重放拒绝、下一 step 新码仍可用。

## 回归验证

- `go build ./...` exit 0；`go vet ./...` exit 0。
- 定向包：`internal/config` / `server` / `cmd/server` / `modules/authsession` / `modules/mfa` 全绿。
- 全量 `go test ./...`（apps/api）exit 0（含 F-007 后复跑）。

## 状态

F-001～F-004 实现 + 负例 + 回归完成；闭合证据待独立审计后在 03-audit 响应节正式标记 fixed。