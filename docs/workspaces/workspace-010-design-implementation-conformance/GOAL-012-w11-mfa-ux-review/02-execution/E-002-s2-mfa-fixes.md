---
id: E-002-s2-mfa-fixes
doc: execution-entry
goal: GOAL-012-w11-mfa-ux-review
date: 2026-08-15
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-002 · S2 MFA 三缺陷修复（M-01～M-03）

## 事实（全部已落地并验证）

1. **M-01 二维码**：新增 qrcode-generator 依赖 + QrCode 组件（SVG）；mfa-manager 绑定区显示二维码（otpauthURL 编码）与手动密钥兜底；新增 i18n 键 schema.account.mfa.qr / qrHint（zh+en）；qr-code.test.tsx 3 用例。
2. **M-02/M-03 401→400**：apps/api/internal/handler/mfa.go 新增 writeSelfServiceMFAError（ErrMFAInvalid→400 MFA_INVALID；NotEnrolled/PendingOnly/Active→400；其余 500），confirm/disable/rotate 三处切换；登录二步验证保持 401 不变。mfa_test.go：fake Confirm 校验动态码；错码 confirm/disable 断言 400；TestMFALoginTwoStep 401 断言保留。
3. **M-02 解绑成功 UX**：mfa-manager disable 成功后本地置 disabled 状态 + sessionStorage["mfa.disabledNotice"]=1 + logout()（AuthContext 置 unauthenticated）；LoginPage 消费标记显示一次性提示横幅（可关闭，i18n login.mfaDisabledNotice）。
4. **M-03 重填**：confirm 错码 400 → 错误文案（error.mfaInvalid）内联显示、enrollPayload 保留、不登出（回归测试覆盖）。
5. **验证**：go test ./... 全绿；web 受影响测试（mfa-manager 5 用例含 2 个新回归、qr-code 3、LoginPage 12 含 notice 1）全绿；tsc 0。
