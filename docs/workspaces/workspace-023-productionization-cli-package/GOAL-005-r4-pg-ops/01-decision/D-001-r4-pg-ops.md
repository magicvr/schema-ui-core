# D-001 · R4 PG 与运维定案（2026-08-29）

## 裁决

1. **PG 实例**（用户指令）：docker postgres:16 容器 gf-pg（host 15432 → 5432；库 golden）；本地测试环境动作、可逆。
2. **下游组合根双方言**：golden-field cmd/server/main.go 参数化 -dialect sqlite|postgres + -dsn（默认 sqlite 保持内嵌默认——契约平等）；CLI 模板同步（双轨一致）。
3. **运维文档范围**：ops-playbook（启动/升级/迁移/备份/停机）+ compose 样例（postgres + 下游应用容器）。
4. **团队化**：golden-field CI 消费回归 workflow（workflow_dispatch + repository_dispatch(published)）。