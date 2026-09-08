---
id: E-012-release-v0.6.0
parent_goal: GOAL-001-nav-group-collapsible
doc: execution-entry
status: recorded
date: 2026-09-08
created: 2026-09-08
updated: 2026-09-08
version: 0.1.0
---

# E-012 · apps/api/v0.6.0 PR、合并、tag 与 Release 资产

## 已发生事实

1. 发布 PR [#11](https://github.com/magicvr/schema-ui-core/pull/11) `feat: ship grouped Admin navigation and apps/api v0.6.0` 已创建并合并。
   - head：`049f541f70a1ef2dd8802afa482d725425672cf9`
   - merge commit：`56ff05321202568a9b19593d3924ed0b7c3ae65a`
   - 合并时间：2026-09-08 09:09:40Z
2. PR 级 Hosted CI run [34207998495](https://github.com/magicvr/schema-ui-core/actions/runs/34207998495) 的 9 个 job 全部 `success`：web、api、api-postgres、4 个 browser E2E 方言/Profile 组合、mvp/admin container smoke。
3. 合并后的 `main` push CI run [34208455485](https://github.com/magicvr/schema-ui-core/actions/runs/34208455485) 的 9 个 job 也全部 `success`。Actions 的 Node 20 deprecation/cache annotations 为 warning，不影响 job 结论。
4. 在 merge commit `56ff05321202568a9b19593d3924ed0b7c3ae65a` 上创建并推送 annotated tag `apps/api/v0.6.0`；tag object 为 `5fe57d79c19db29c89d3f754048af2e3272e2432`，peeled target 与 merge commit 一致。
5. GitHub Release [schema-ui-core apps/api v0.6.0](https://github.com/magicvr/schema-ui-core/releases/tag/apps/api/v0.6.0) 已于 2026-09-08 09:37:59Z 发布，草稿已转正式 Release。
6. Release 资产共 7 个，均为 `state: uploaded`：
   - `schema-ui-v0.6.0-darwin-amd64.zip`
   - `schema-ui-v0.6.0-darwin-arm64.zip`
   - `schema-ui-v0.6.0-linux-amd64.zip`
   - `schema-ui-v0.6.0-linux-arm64.zip`
   - `schema-ui-v0.6.0-windows-amd64.zip`
   - `schema-ui-v0.6.0-windows-arm64.zip`
   - `SHA256SUMS`
7. 资产完整性已独立核对：本地构建后逐文件 SHA-256 校验通过；再从已发布 Release 下载全部 zip，并按 `SHA256SUMS` 复核，6 个二进制压缩包均 `REMOTE_OK`。

## 证据边界

- 本地 `scripts/pre-release-smoke.sh` 未执行成功，原因是当前 checkout 未配置 `ADMIN_INITIAL_PASSWORD` / `AUTH_JWT_SECRET`，未猜测或输出秘密；Hosted container smoke 已在 PR 和合并后 CI 中通过，作为本次发布门禁证据。
- 本条记录 Hosted CI、merge、tag、Release 的已发生事实，不改变 Goal/VP 的 `status` 或 `progress`；VP-034 愿景层仍按 `/vision` 另行处理。
