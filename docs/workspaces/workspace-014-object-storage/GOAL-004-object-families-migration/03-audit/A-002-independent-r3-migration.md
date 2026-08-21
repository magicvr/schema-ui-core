---
id: A-002-independent-r3-migration
title: 独立交叉审计 · R3 三类落盘收口走端口（R3→R4 门禁）
source: independent
date: 2026-08-21
scope: GOAL-004 R3→R4 门禁（三类落盘经端口 / 行为保持与 D-001 §5 / ObjectInfo.ModTime / 单实例装配 / DELETE 幽灵边车与 import owner / 测试覆盖）
verdict: conditional
status: recorded
parent: GOAL-004-object-families-migration
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# A-002 · 独立交叉审计：R3 三类落盘收口（verdict: conditional）

## 范围与区间

- auditor: grok-build (grok-4.6 · reasoning high)
- type: stage
- covered: 被审方案 D-001；实现提交 `d99221f`；三类 handler 持久化是否经端口；GC/配额/owner/遗留兼容与 §5 声明；`ObjectInfo.ModTime` 相对 R1 冻结面；`newObjectStore` local/s3 两分支；file-library DELETE 幽灵边车与 import owner；测试覆盖；`git diff HEAD~1 --stat` 越界面。
- excluded: R4 公共面收口、R5 双路径 live 验收、目标 `status`/`progress`（本意见不改）；未复跑全量 `go test ./...`（以 E-001 自述 + 本轮定向复跑为准）。
- 工作区：`workspace-014-object-storage`（`root_goal: GOAL-001-object-storage`，`canonical_scope` 匹配，`shared_materials_catalog: none`）。
- 对照提交：`d99221f`。方案：GOAL-004 D-001。R1 对照：GOAL-002 D-001。P-005：GOAL-004 无自有到期 required 信息项；Root I-004 `recorded`（non-blocking，R5 叙事）。
- 既有意见：A-001 self `pass`（编排器已落盘）。本意见独立形成，不采纳其 verdict。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 三类生产持久化已改端口；handler 生产文件无 `os.ReadDir`/`WriteFile`/`Remove` | `raster_assets.go:271-280,287-301,304-311,347-369,386-408`；`upload.go:128-144,155-167,175-209`；`filelibrary.go:111-129,205-213`；`import.go:388-401`；对本四文件 grep 无直接 IO |
| 模块合同改 `kernel.ObjectStore`；composition 单次构造后注入五处消费方 | `filelibrary/provider.go:32-34,64`；`datatransfer/provider.go:32-33,60`；`composition.go:311-321,327-328,346-347,403,409` |
| local root 派生与 s3 探针符合 D-001 §1；未知 driver 由配置面 fail-closed | `composition.go:590-605`；`config.go:509-526`；`main.go` 无「s3 未接线」警告 |
| `ObjectInfo.ModTime` 为结构体加性字段；接口方法集未改；两适配器同提交补齐 | `kernel/objectstore.go:88-96,108-115`；`local.go:185`；`s3.go:191-193`；GOAL-002 D-001 §1 方法集对照 |
| owner 门禁保留：upload GET、library confirm、import（空 owner fail-closed） | `upload.go:357-363`；`filelibrary.go:294-298`；`import.go:168-172`；`data_transfer_test.go:197-211`；`filelibrary_test.go:244-267` |
| 幽灵边车 DELETE：Exists(body)→Delete→!exists 则 404；local Exists 只看 body | `filelibrary.go:205-217`；`local.go:202-216,188-199`；D-001 §5 第一点 |
| 无越界：`d99221f` 18 文件均在 api 三类落盘/适配器/装配 + GOAL-004 E-001 | `git show d99221f --stat` |
| 本轮定向复跑绿 | `go test ./internal/composition -run TestNewObjectStoreWiring`；`go test ./internal/objectstore`（exit 0） |

## 对照成功标准

| 焦点 | 状态 | 证据 |
|------|------|------|
| (a) 三类持久化全部经端口 | 达成 | 上表；漏网 IO 未发现 |
| (b) 行为保持 + §5 完整诚实 | 部分 | GC/owner/布局兼容成立；§5 未声明幽灵边车 **list 不可见 + 配额仍计入**（F-001） |
| (c) ModTime 不破坏 R1 冻结面 | 达成 | 方法集零变化；非发布 API；无无键字面量依赖 |
| (d) 单实例装配 local/s3 正确 | 代码达成、测试弱 | `composition.go:590-605,311-409` 正确；测试未锁 root/类型/指针同一（R-001） |
| (e) DELETE 幽灵边车与 import owner 安全 | 达成（实现） | 见成果表；幽灵 DELETE 无专项测试（R-002） |
| (f) 测试覆盖 | 有明显缺口 | 见 R-001～R-003；检查点 2 全量绿本轮未复跑（N-004） |

## Findings

| ID | 级别 | 主张 | 证据 | 建议处置 |
|----|------|------|------|----------|
| F-001 | required / med | D-001 §5 对已知差异声明不完整：幽灵边车（仅 `.meta.json`、无 body）旧 scan 入列，新 scan 因 Stat miss 被跳过；同一 id 仍被 List 出来并在配额路径保守计入。残损数据会出现「库列表看不见、配额仍占、UI 无法点删」——超出 §5 仅记录的 204→404。 | 旧：`d99221f^` `filelibrary.go` scan 只迭代 `.meta.json`、body Stat 失败仍 `append`。新：`filelibrary.go:118-120` Stat err→`continue`；`local.go:241-246` List 把 sidecar 名映射为 id；`upload.go:187-195` Stat err→`files++` 且 `bytes += maxUploadBytes`。§5：`D-001-r3-migration-freeze.md:36-40`。 | 补记 §5（范围=前端口残损；清理=按 id DELETE 或盘面）或恢复 sidecar-only 入列。文档补记即可闭合。 |
| R-001 | recommended / med | D-001 §6 / 检查点 3 要求的「单实例断言（local 与 s3）」与 root 派生未被测试锁住。`TestNewObjectStoreWiring` 的 s3 分支丢弃 store、不断言 `*S3Store`、不断言消费者指针同一、不设 `ObjectsLocalRoot`/`DBPath`。代码审查显示单次 `newObjectStore` 注入五处，实现正确。 | `objectprobe_test.go:13-47`；`composition.go:311-314,319,327-328,346-347,403,409,600-604`；D-001 §1、§6；`00-meta.md` 检查点 3。 | 补：local root 派生；s3 返回类型；同一指针注入（或明确把检查点 3 降为代码审查门禁）。 |
| R-002 | recommended / low | 幽灵边车 DELETE（Exists=false → best-effort Delete → 404）无专项测试；现测只覆盖有 body 的 204 与再删 404。实现与 §5 第一点一致，安全面依赖 `files.delete`。 | `filelibrary.go:205-217`；`filelibrary_test.go:167-178`。 | 加 sidecar-only fixture：404 + 边车被清。 |
| R-003 | recommended / low | `ModTime` 是 file-library `created` 的唯一来源，适配器与 library 测试均不断言非零；回归成零值会得到 `0001-01-01T00:00:00.000Z`。 | `kernel/objectstore.go:92-95`；`local.go:185`；`s3.go:191-193`；`filelibrary.go:85,126`；`rfc3339.go:7-8`；`local_test.go:46-52` / `s3_test.go:154-157` / `filelibrary_test.go:117-126` 均不查 `created`/`ModTime`。 | Stat 测断言 `ModTime` 非零；list/detail 断言 `created` 形如 RFC3339。 |
| R-004 | recommended / low | 坏 JSON 边车 + 仍在的 body：旧 CountOwner/quota 走「读失败则保守 +1」；新 `parseMeta` 吞错返回空 meta，Stat 成功后按 owner 不匹配跳过，不再保守计入。远程不可利用（需改盘），但 A-001「坏 meta 保守计」过满。 | 旧 CountOwner：`d99221f^` unmarshal 失败 `count++`。新：`local.go:94-105,185`；`raster_assets.go:356-367`；`upload.go:187-202`；A-001 N-002 行。 | 留痕到 §5 或 CountOwner/quota 对空 owner + Stat 成功走保守分支。 |
| N-001 | note | CountOwner 对 List 出的 ghost id（Stat NotFound）无条件 +1，旧实现按 meta.owner 匹配才 +1。两方向都保守，归属变宽。同意 A-001 N-201。 | `raster_assets.go:356-361`；对照旧 CountOwner。 | 留痕即可。 |
| N-002 | note | 空 Name/Type/Owner 的 scan 跳过 = 无 meta 字段的 body 不可见。S3 无 sidecar 形态。同意 A-001 N-202。§5「依旧不可见」对 body-only 为真；对可解析的 `{}` 边车旧实现会入列（现跳过）——被 F-001 吸收。 | `filelibrary.go:123-125`；D-001 §2、§5。 | 无需动作。 |
| N-003 | note | main.go「s3 未接线」整段已删，前提消失。同意 A-001 N-203。 | `git show d99221f -- apps/api/cmd/server/main.go`；现 `main.go:47-51`。 | R4 若再出现「已配置未生效」窗口再重建提示。 |
| N-004 | note | 检查点 2 全量绿以 E-001 自述为准。本轮独立复跑 composition 接线测与 `objectstore` 包 exit 0；未复跑 handler 全套。 | E-001「验证证据」；本轮 `go test` 输出。 | 编排器若要以检查点 2 放行，应确认全量套件仍绿。 |
| N-005 | note | `newObjectStore` 对非 `s3` 一律走 local。生产经 `config.Load` 拒绝未知 driver；测试若绕过 Load 传入 `"S3"` 会静默 local。 | `composition.go:591-605`；`config.go:509-526`。 | 非门禁；可选在 `newObjectStore` 再校验一次。 |

## 必改项汇总

1. **F-001**：把「幽灵边车不再出现在 library 列表、但仍被配额/CountOwner 保守计入」写入 D-001 §5（或恢复入列）。未闭合前不得把「§5 已穷尽已知差异」当作 R3→R4 无条件放行依据。

开放 required = 1（F-001）。无到期 required 信息项。

## 与既有意见的异同

- 同意 A-001：生产路径无直接 IO；owner 三处保留；ModTime 加性合法；装配代码正确；N-201/N-202/N-203。
- 不同意 A-001「开放 required = 0 / verdict pass」：§5 对幽灵边车 **list×配额** 组合未声明（F-001）。
- A-001「坏 meta 保守计」过满，降为 R-004。
- A-001「接线锁测试 / 单实例装配回归」：测试锁的是探针而非单实例/S3 类型/root（R-001），不升 required（代码审查已覆盖正确性）。

## 结论

实现把三类落盘收口到同一 `kernel.ObjectStore` 实例，端口方法集未被破坏，owner 门禁与 DELETE 幽灵边车的 HTTP 差异与方案一致。R3→R4 不能无条件放行：§5 漏记幽灵边车在列表面消失但仍占配额这一残损数据组合（F-001）。建议 `/govern` 先补记 §5（或恢复入列）闭合 F-001，再决定是否带 R-001～R-004 进入 R4。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。

---

## 编排器响应与闭合（/govern · 2026-08-21）

| finding | 处置 | 证据 |
|---------|------|------|
| F-001 (required) | **fixed**（采纳"补记 §5"路径）：D-001 §5 补记幽灵边车组合差异——"列表不可见、配额仍占、UI 无法点删"，含范围（前端口残损）与清理路径；同节补记 R-004 的坏 JSON 边车口径变化与 CountOwner ghost +1 | GOAL-004 D-001 §5（A-002 F-001/R-004 标注的两条补记） |
| R-001 | **fixed**：TestNewObjectStoreWiring 强化——local 分支断言 *objectstore.LocalStore 类型 + root 派生（DBPath fallback 与 ObjectsLocalRoot override 两例）；s3 分支断言 *objectstore.S3Store 类型。检查点 3 恢复为测试锁 | composition/objectprobe_test.go |
| R-002 | **fixed**：新增 TestFileLibraryGhostSidecarDelete——sidecar-only fixture 断言 404 且边车被 best-effort 清理 | handler/filelibrary_test.go |
| R-003 | **fixed**：TestLocalStatModTimeNonZero 断言 local Stat 的 ModTime 非零；S3 stub 未设 LastModified 不适用（真实 LastModified 由 live 测试在 R5 覆盖） | objectstore/local_test.go |
| R-004 | **fixed**（文档路径，见 F-001 同批 §5 补记）；如需恢复保守口径留 R4 评估 | D-001 §5 |
| N-001/N-002/N-003/N-005 | 留痕不动作（与 A-001 一致；构造函数二次校验列为可选加固） | 本节记录 |
| N-004 | fixed：响应修复后全量 go test ./... 复跑 exit 0（本提交前完成验证） | E-002 |

**门禁判定**：F-001 以 fixed 合法闭合，开放 required = 0 → R3→R4 门禁放行；GOAL-004 三检查点满足，结项（Root progress 3/5）。
