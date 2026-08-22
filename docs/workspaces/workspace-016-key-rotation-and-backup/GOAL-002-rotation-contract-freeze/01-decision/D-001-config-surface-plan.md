---
id: D-001
doc: decision-entry
goal: GOAL-002-rotation-contract-freeze
status: accepted
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# D-001 · R1 配置面实施方案

## 决策

按 Root D-002 合同，配置面实现切片如下（全部限于 `internal/config` + 配置样例 + 单测）：

1. **字段**：`Config.AuthJWTSecretPrevious string`（空 = 单密钥模式）。
2. **解析双通道**：YAML `auth.jwt_secret_previous`（指针可选语义，沿 `strPtrOr` 先例）；env `AUTH_JWT_SECRET_PREVIOUS`（env 永远覆盖 YAML，沿既有次序）。
3. **ValidateProd 规则**（仅非开发环境块内，与现行结构一致；development 保持低门槛不受影响）：
   - previous 非空时：长度 ≥ `minJWTSecretLen`（32）且 `containsLettersAndDigits` 通过——与 current 同规则同措辞风格；
   - previous == current（非空相等）→ 启动失败（同值守卫）；
   - 错误信息只点名键名（`AUTH_JWT_SECRET_PREVIOUS`），不携带值。
4. **样例文档**：`config.default.yaml` 与 `configs/config.yaml` 的 auth 节增加注释键（默认不启用），注明轮换操作顺序（先加 previous → 重启 → 移除旧 key → 再重启）。
5. **测试矩阵**：
   - production + 合规不同 previous → 通过；
   - production + 短/全字母/全数字 previous → fail closed 且错误点名 PREVIOUS 键；
   - production + previous == current → fail closed（同值守卫）;
   - development + 弱 previous → 通过（低门槛不变式）;
   - 缺省无 previous → 行为零变化（既有 production 用例继续通过）;
   - YAML 层解析 + env 覆盖 AUTH_JWT_SECRET_PREVIOUS 各一条。

## 为什么

- 校验放 ValidateProd 而非 Load：与 current secret 的熵规则同一位置、同一测试面，zero-value/test Config 也被覆盖。
- dev 不设防：保持"轮换不是 mvp/dev 硬依赖"的 VP 内嵌默认。
- 不在本切片改 `internal/auth` / composition：R1 只冻结配置面；验签消费是 R2（GOAL-003），避免半接线状态。

## 未选方案

- 在 Load 里做同值守卫：Load 无 AppEnv 语境且现有先例把生产规则集中在 ValidateProd。
- previous 用独立强制度（如 ≥48）：无依据；统一 32 与 current 对齐，审计面最小。
