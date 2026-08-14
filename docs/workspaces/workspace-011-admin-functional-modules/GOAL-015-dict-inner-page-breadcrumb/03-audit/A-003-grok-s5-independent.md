---
id: A-003
goal: GOAL-015-dict-inner-page-breadcrumb
source: independent
date: 2026-08-14
scope: S1～S4 交叉审计（execution-facts + go 影响）
verdict: fail → 处置后 required 全部 fixed
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-003 · grok 4.6 独立审计（S5）与处置

## 审计主体

- auditor：Grok 4.6（high think，`--no-subagents --disable-web-search`），用户书面指定。
- 原始意见：**verdict fail**，5 项 required（F-001～F-005）+ 3 项 recommended（F-006～F-008）+ 2 项 recommended（F-009）+ 2 项 note（F-010/F-011）。

## 处置（编排器）

### required · 全部 fixed

- **F-001（high）条目导航重置 visitStack → 主路径无面包屑/返回** → fixed：render.tsx navigate 分支改走 Host onNavigate（pushState + visitStack），App.tsx SchemaPageSurface→RenderPage 接线；App.integration 新增行导航→内页 trail+返回集成测试（9/9）。
- **F-002（high）JOIN 后 ORDER BY 歧义（sort/updatedAt 500）** → fixed：repository.go 排序列限定 de. 前缀；dictionary_test.go 增 dictKey+sort=sort、sort=updatedAt、dictKey+sort+page 回归（handler 146s 全绿）。
- **F-003（medium）dictTypeName 列 sortable 服务端白名单缺失 → 400** → fixed：SortFields 增 dictTypeName，store 映射 dt.name；sort=dictTypeName 回归入测试。
- **F-004（medium）progress 4/5 非检查点派生** → fixed：00-meta 检查点按事实重算——S1～S4 已勾（S4 在 F-002/F-003 修复与 I-001～I-004 闭合后勾选），progress 4/5；派生说明文字更正。
- **F-005（medium）I-001～I-004 未闭合 + D-003/A-002 损坏（NaN）** → fixed：I-001～I-004 全部 closed（挂证据列）；D-003/A-002 重写（NaN 清除，version 1.1/1.2）。

### recommended / note · 处置

- **F-006（create 表单显示类型键而非类型名）** → fixed：行导航 navigateMapping 增 dictTypeName=$row.name；create 表单增 dictTypeName readOnly 显示（路由种子）；T-DE-03 断言补类型名。
- **F-007（tombstone 全量列表 + 只读必填创建）** → 保留为产品取舍，待用户裁决（见 A-003 结论）：(a) 无 dictKey 内页 fail-closed；(b) 接受无键=全局列表。当前为 (b)+创建表单必填 fail-safe。
- **F-008（运行时未门禁 data.route-binding；claim 停 2.8.0）** → fixed：dispatchParsedNode 对带 $context.route.* 绑定的 data 节点按页面 meta 门禁（>=2.9 + capability，否则 fail-closed 报错）；generate-claim.mjs 升 2.9.0（content c87c22ad…/fixture 89baddbc…，pageVersions [2.7,2.9]、manifestVersions [2.7,2.8,2.9]、capabilities 增两项、suites 增 request-construction/component-format），claim 重新生成；boot.ts supportedCapabilities 增两项。
- **F-009（测试与台账数字不严）** → fixed：readOnly 门禁 4 例落 form-controls.test.ts；extraQuery 合并 + F-001 防绕过 2 例落 resource.test.ts；E-008 计数更正（946/946 + 修正单测描述）。
- **F-010（readOnly 可伪造；schema 非安全边界）** → note 接受：写入 D-003/A-002（过滤与只读是 UX/契约，服务端仍独立鉴权；与 dictionary.write 同权非提权）。
- **F-011（路由绑定读 location.search，path params 恒空）** → note 接受：与 A-002 F-003 同源；后续路径参数页需接入 Host 路由上下文。

## 复核证据

- web 全量回归（53 文件 / 946+ 测试）绿；protocol+host 489/489；api go test ./... 绿。
- 独立复核点（审计员复核）：vendor pin 30 项 LF sha256 对齐；docs/schemas 与上游 81aa1d8 字节一致；T-DE-01/02/03/04 复跑通过。

## 结论

- 5 项 required findings 已按 fixed 闭合（可核对：测试/代码/记录引用如上）。
- F-007 为产品取舍，**待用户裁决**（accepted-residual 或 fail-closed 改造）；不影响技术关门，但在关门汇报中显式列出。
- S5 关门放行条件：required findings 全闭合（达成）+ goal-tree 5/5 + 用户裁决 F-007。
