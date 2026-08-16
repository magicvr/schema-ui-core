---
id: A-004
doc: audit
source: self
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# A-004 · T-06 通知中心交互修正 自审（source: self）

- **日期**：2026-08-16
- **scope**：铃铛下拉条目可点击、列表点击即读+展开、移除行内标为已读、未读数即时刷新
- **verdict**：**pass**

## 核对清单

| 项 | 结论 | 证据 |
|----|------|------|
| 问题 1（铃铛条目不可点击） | ✅ | 条目为 menuitem 按钮：点击 → read POST + onOpenItem 深链跳转；单测断言 onOpenItem 与 POST |
| 问题 2（列表点击不展开不标记） | ✅ | 行点击展开详情面板 + 未读自动标记；单测（展开/aria-expanded/read POST/已读不重复/深链） |
| 问题 3（行内标为已读多余） | ✅ | notifications.json 移除 markRead 行 action；保留全部标为已读与搜索/筛选 |
| 问题 4（未读数不即时刷新） | ✅ | read/read-all 响应头 → 铃铛订阅命名空间即时重查徽标（含下拉条目）；Go 头断言 + 铃铛事件单测 |
| 方案卫生 | ✅ | 未改协议 schema（GOAL-018 本地 custom 扩展合法，D-VAL 绿）；通用 Renderer 无产品端点感知；无误导性 toast |
| 回归 | ✅ | Go 0 FAIL；vitest 1037/1037；tsc 0；e2e admin/mvp 8/8（含通知页冒烟） |
| go 判定 | ✅ | 仅响应头 + 前端组件 → 无装配语义变化 → 无影响、不暂挂 |

## Findings

无 required/必改 findings。

## 附带说明（不阻断）

- 通知「有数据」e2e 缺位（见 E-011 说明）：系统无管理端造数端点；交互已由组件级单测覆盖。
