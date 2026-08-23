---
source: self
date: 2026-08-23
scope: GOAL-036 范围升级（全盘修复 + 防复发，S5）
verdict: pass
---

# D-002 · 范围升级：全盘修复此类问题并防复发

## 用户裁决（2026-08-23，书面）

对「其他页面是否也存在类似问题」的排查结论之上，用户决定：**升级 GOAL-036（W25）为全盘修复此类问题，并确保以后不会出现此类问题**。配套拍板：

1. **纳入**：系统监控页自动刷新定向化（只刷 `/api/system-monitoring/status`，6 张 statCard 的共享端点）；
2. **纳入**：schema 渲染层自定义组件注册校验测试（防新组件脱管）；
3. **出局**：大表 `COUNT(*)` 计数优化（activity / operationlog / file-library 等列表页）——理由：与本类页面性能反模式正交，属**容量课题**（索引健康时当前无碍），不并入本波；如需处理另行立项；
4. **文件夹改名**：`GOAL-036-w25-wallet-page-performance` → `GOAL-036-w25-page-performance-guardrails`（引用同步）。

## 全盘扫描结论（触发裁决的证据）

26 页 schema + 10 个自定义组件扫描（E-001）：

| 反模式 | 命中页 | 处置 |
|--------|--------|------|
| 同 URL 重复展示节点 | my-wallet 3×、system-monitoring 6×、data-display 3× | 渲染层 in-flight 合并全局覆盖（无需逐页改） |
| 挂载即写 + 整页重拉 | 仅 wallet-ensure | 探活后写契约（已改） |
| schema 每次导航重取 | 全部页面 | App 级文档缓存全局覆盖 |
| SQLite 单连接串行 + fsync | 全部 API | store 层全局修复 |

另确认：file-library/users/roles 中的 `{type:"custom", handler:"…"}` 属**动作层自定义动作**（非渲染层组件），不在注册校验范围；monitoring-auto-refresh / account-session-toolbar 的 `reloadList` 调用均为用户触发或显式轮询，语义正确。

## 防复发机制设计（三支柱）

1. **行为回归测试（把机制钉死在测试里）**：
   - `store_wal_test.go`：`sqliteDSN` pragma 单测 + 文件库白盒断言（池 = 4、`journal_mode=wal`、`busy_timeout=5000`、`synchronous=1`）+ 内存库单连接 + `DB_POOL_MAX_OPEN` 覆盖——任何回退 `MaxOpenConns=1` 的改动立即红；
   - `render.test.tsx`：三 statCard 合并（既有）+ 新增 statCard+chart 合并、`refreshList` 定向刷新（只重拉目标 URL）回归；
   - `custom-components.schema.test.ts`：模块 schema 声明的全部渲染层自定义组件必须已注册（镜像 main.tsx 的自注册导入）。
2. **开发规范**：`docs/architecture/module-contribution-playbook.md` §6「页面数据面性能规范」——同源多展示节点由合并机制兜底（允许）；自定义组件**禁止挂载即写 + 整页 `reloadList()`**，须探活后写；定时刷新用 `refreshList` 定向；新组件必须注册 + 行为回归测试。
3. **语义取舍留痕**：
   - `refreshList` 只作用于标准展示形态（`DISPLAY_LIST_QUERY`，无路由绑定的 statCard/chart 场景）；
   - monitoring 定向刷新后**事件表不再随 tick 刷新**（手动刷新或等待整页刷新时更新）——符合「监控状态为主」的页面语义，用户已确认范围。

## 未选方案

- 对重复 dataSource 施加 schema 校验/告警：**不采用**——合并机制已使其无害，校验只增噪声；由规范文档说明「允许但无需规避」。
- `refreshList` 支持多 URL / 通配：**不采用**——当前唯一消费方（monitoring）单端点即可，多 URL 可后续按需扩展。
- 大表 COUNT 物化统计：**出局**（见上），另行立项时再评估。