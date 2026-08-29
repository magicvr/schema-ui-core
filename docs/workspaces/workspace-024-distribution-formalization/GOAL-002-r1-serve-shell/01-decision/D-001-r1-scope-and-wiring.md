---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-serve-shell
version: 0.1.0
---

# D-001 · serve 面构成与模板形态（2026-08-29 · 用户 P-004 裁决）

## 决策

1. **serve 面构成 = 方案 A · 标准下游组合闭环**（用户裁决 · I-001 关闭）：
   - 新建公开包 `apps/api/server`：config 装载（RT-K01 语义子集 + fail-closed）+
     标准下游组合装配（`core.server-registration` / `core.auth-session` /
     `core.manifest-route` / `core.navigation-capability` / `core.schema-render` /
     `core.operationlog` / `admin.users`）+ 中央面接线（healthz/readyz、登录、
     schema 文档、manifest/bootstrap、模块贡献路由）+ RT-D02 §1 全序停机。
   - **不采用**方案 B（全量对等：主仓 composition 公开化，面过大、与薄装配边界冲突）
     与方案 C（最小壳：无中央面，探针/登录不可用，不算可用基架）。
2. **模板形态 = 方案 A · 薄封装单一形态**（用户裁决）：
   - `cmd/server` 改为 thin wrapper（`-config` / `-dialect` / `-dsn` / `-addr`
     flag 解析 → `server.Serve`）；组合装配代码移入公开 serve 面。
   - `-dialect/-dsn` 语义与旧模板一致（postgres → DSN；sqlite → 文件路径），
     保持 `schema-ui upgrade` 冒烟与 golden-field 既有调用兼容；golden-field
     主密升级在 R2/R3 实证时执行（其现行 main 已注册历史，不回改）。
3. **有界下游基线**（明确不装配，主仓形态保留）：jobs / metrics / tracing /
   objects / mail-admin / settings / mfa / captcha / wallet 等中央面与模块；
   manifest 按标准组合模块集发布（与「kernel + ≥1 标准模块」试点基线一致）。
4. **config 公开面**：`server.LoadConfig(path)`——代码默认 → YAML
   （`${VAR}`/`${VAR:-default}` 插值 fail-closed，KnownFields 严格）+ 进程 env
   定向覆盖（`HTTP_ADDR` / `HTTP_SHUTDOWN_TIMEOUT` / `DB_*` / `AUTH_JWT_SECRET` 等）；
   非法 `http.shutdown_timeout`、方言/DSN 配对错误、非 dev 缺密钥 → 拒绝启动。
5. **RT-D02 接线**：`Serve`（信号版，CLI/模板用）与 `Run(ctx, opts, signals)`
   （可测版）共享同一生命周期：signal/ctx → `shutdown.starting` → `http.Server.Shutdown`
   （预算 = `Config.ShutdownTimeout`，默认 10s）→ kernel runtime 逆序 Stop →
   Store Close → `shutdown.complete`（nil）/ `shutdown.timeout|error`（错误）。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| serve 面构成 | B 全量对等 / C 最小壳 | B：把 internal composition 全公开，等同 fork 替代、破坏薄装配试点边界；C：无登录/schema/manifest 中央面，探针不绿、不算可用基架 |
| 模板 | B 双入口并存 | 两套组合代码膨胀、双轨漂移风险（GOAL-003 F-001「双轨同构」教训） |
| config | 完整镜像主仓 config（.env 文件层/全键） | 下游基线只需可核对子集；秘密走进程 env，与主仓 env 键对齐 |

## 信息门禁

- I-001（serve 面构成 + 模板形态）→ **verified**（用户 P-004 裁决，2026-08-29）。
- I-002（config 装载默认形态）→ **verified**（内嵌默认 + 显式文件 + env 定向覆盖，本决策 §4）。