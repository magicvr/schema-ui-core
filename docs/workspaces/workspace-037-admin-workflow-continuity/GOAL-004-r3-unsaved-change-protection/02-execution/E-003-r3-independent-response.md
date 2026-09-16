---
id: E-003-r3-independent-response
doc: execution
goal_id: GOAL-004-r3-unsaved-change-protection
status: recorded
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-004-r3-unsaved-change-protection
version: 0.1.0
---

# E-003 · 响应 A-002 required/recommended

2026-09-17，响应本目标 A-002 independent 的 `conditional` 意见：

- F-001（required）按 `fixed` 路径处理：内部导航测试在确认前后都保持 `dirty=true`，并断言第二次真正调用 `window.confirm` 后才切换到 `/catalog`。
- F-002（recommended）已补真实 Schema default form 改值后点击 App 菜单的集成回归，确认取消后仍停在当前页面。
- F-003（recommended）已补客户端 required 校验不发请求仍保留 dirty、transport throw 保留 dirty，以及 modal 成功提交后卸载并清理 dirty source。
- F-004（recommended）已把矩阵的 search 直接回归范围收窄为 q + “search form 不注册”的结构证据，避免把未直接操作的筛选/排序/分页写成测试事实。
- F-005（recommended）已在 App `onNavigate` 对 exact committed href 做 no-op，并补 dirty 同路径不确认、不新增历史的回归；该取舍已补入 D-001。

修正后受影响测试 8 个文件、140 项通过，`npx tsc -p tsconfig.app.json --noEmit` 通过。F-001 尚需 independent recheck 复核后，才能作为 C4 关门依据；在此之前不修改目标状态。
