---
id: E-019
doc: execution-entry
goal: GOAL-001-store-dialects
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-019 · PR #4 CI 全绿治理：container-smoke(admin) 矩阵腿修复

## 2026-08-20 · 推送 PR 后的一处 CI 红 → 定位与修复

### 已发生事实（事实层）

- 推送 `dev`、建 PR #4（dev→main，**不合并**）后 CI 结果：`api` ✅、`api + postgres` ✅（新增双方言 job 首个绿 run）、web ✅、browser E2E(mvp/admin) ✅、**container smoke(mvp) ✅、container smoke(admin) ❌**。
- 失败详情：`SM-007=FAIL: admin Manifest 缺少 settings 页面`（exit 5）。
- **根因（已定位，非 workspace-013 逻辑改动引入）**：
  - API 的 profile/启用模块集按 W7/T-06 **只认 operator 配置 `apps/api/configs/config.yaml`**（`APP_PROFILE` env 刻意不读）。
  - 该 operator 配置文件**在 origin/main 上不存在**，是 dev 分支（441 提交范围内，T-06/W7 阶段）引入；其 `app.profile: mvp` → mvp 下 settings 关闭。
  - `container-smoke` 矩阵 `[mvp, admin]` 用 `PRERELEASE_PROFILE` 只改变 smoke 的**期望**，未改变 API 实际 profile → admin 腿期望 settings 而 API 无 → 红。
- **修复**（`scripts/pre-release-smoke.sh`）：显式给定 `PRERELEASE_PROFILE` 时，生成一份 temp operator 配置（复制 `configs/config.yaml` 并把 `app.profile:` 盖写成该 profile），并在 compose override 给 api 设 `CONFIG_FILE` 指向它（bind mount 只读）；独立未指定 profile 的路径保持逐字节不变。
- `CONFIG_FILE` 受 config.Load 支持（显式路径必须存在）；sed/pwsh 模拟 `profile: mvp→admin` 缩进正确；本地无 bash（无 WSL）→ 全量验证由 CI 复跑承担。

### 证据

| 主张 | 路径 / 证据 |
|------|-------------|
| PR + CI run | github.com/magicvr/schema-ui-core/pull/4（run 32432087363） |
| 根因 | `apps/api/configs/config.yaml` `app.profile: mvp`；`scripts/smoke.sh` SM-007 admin required_pages 含 settings；compose 注释 W7 不读 APP_PROFILE |
| 修复 | `scripts/pre-release-smoke.sh`（PRERELEASE_PROFILE 分支盖写 + CONFIG_FILE override） |
