---
id: E-008-r4-migration-and-matrix
doc: execution-entry
parent_goal: GOAL-001-nav-group-collapsible
status: recorded
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# E-008 · R4 当前 sidebar 全量迁移与 Profile/route 矩阵

## 已发生事实

1. 当前已注册的 17 个 sidebar NodeID 已全部在对应模块 Provider 声明五组 Group 元数据；Dashboard 顶层单例、top/user slot、通知铃面和 Examples authored group 保持边界。
2. 新增 API runtime Profile matrix：`mvp`、`admin`、`demo`、`custom=admin+channel.telegram`、`custom=admin+channel.telegram+biz.digital-offer`；验证五组成员、不产生空组、optional 共组、top/user pageRef 不丢失。
3. 新增 Web route projection matrix：覆盖 17 个 grouped/sidebar NodeID 的直接 route，并覆盖 `users-invites`、`wallet-entries`、`dictionary-entries`、`task-runs`、`telegram-operator` 等已登记内页/动态深链父级自动展开。
4. 更新 `docs/architecture/module-contribution-playbook.md` 至 v1.2.0，落盘 group key/order/metadata、跨模块共组、冲突 fail-closed、slot 边界和迁移规范。
5. 证据矩阵已写入 `attachments/r4-navigation-route-profile-matrix.md`。

## 验证事实

- API：`go test ./internal/composition -run TestR4NavigationGroupProfileMatrix -count=1` 通过。
- Web：`nav-groups-r4.test.ts` 2/2 通过；R3 `nav-groups.test.tsx` 3/3、`navigation.test.ts` 7/7 继续通过。
- Web 全量回归（包含 R3）：98 files / 1337 tests 通过；`tsc -b` 与 `vite build` 通过。

## 当前门禁

- `I-034-004`：route/deep-link projection matrix 已有证据；需在 R4 self audit 中核对所有页面与父级映射。
- `I-034-005`：optional/custom/demo runtime matrix 已有证据；需在 R4 self audit 中核对无丢失/无空组/权限过滤。
- R4 尚未标记 completed；待 self audit、全量 API 回归和关门前 independent audit。
