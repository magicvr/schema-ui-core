---
id: GOAL-001-design-system-and-ui-experience
doc: audit-entry
record_id: A-010
source: self
scope: 用户驳回 closeout · 主区宽度自适应缺口
verdict: conditional
status: recorded
parent: GOAL-001-design-system-and-ui-experience
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# A-010 · 用户驳回关门 · 流体宽度缺口

## 背景

closeout-ready 提案后用户**未**确认 D-008，并书面指出：除顶栏外页面主内容不随浏览器宽度自适应。

## Finding

### F-VUI-008 · 壳 body 使用 max-w 岛，主区不随视口伸缩

| 字段 | 值 |
|------|-----|
| level | **required**（阻断 Root done，因用户列为关门缺口） |
| status | **fixed**（E-008） |
| evidence | 原 `max-w-[1440px] mx-auto` 包住 sidenav+main；修复后 `data-shell-width=fluid` + `w-full`；shell.test 断言禁止 1440 岛 |
| impact | 用户可观察布局与顶栏不一致 |
| closure | fixed via E-008 code |

## 结论

**verdict: conditional** 直至 E-008 回归绿后视为 F-VUI-008 fixed。Root 仍 **active / closeout-ready**；开放 required 在 F-VUI-008 fixed 后回到 0，仍须用户**显式** D-008 确认才可 done。
