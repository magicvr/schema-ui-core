---
id: A-002-independent-r4-sweep
title: 独立交叉审计 · R4 公共面收尾核查（R4→R5 门禁）
source: independent
date: 2026-08-21
scope: GOAL-005 R4→R5 门禁（E-001 扫描声明可复现性 / 边界声明 / newObjectStore 未知 driver 二次校验 / 公共契约路径残留）
verdict: pass
status: recorded
parent: GOAL-005-public-surface-sweep
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# A-002 · 独立交叉审计：R4 公共面收尾核查（verdict: pass）

## 范围与区间

- auditor: grok-build (grok-4.6 · reasoning high)
- type: stage
- covered: HEAD `8aa0abc`；E-001 三维扫描声明的独立复跑；`apps/api/internal/modules/*/provider.go` 构造器；`composition.newObjectStore` 未知 driver 二次校验；`store.OpenWithCatalog` 与 systemmonitoring `dbPath` 是否属 Store 方言；pkg/web/脚本漏扫维度；Handler/模块公共契约是否仍把本地路径或 `os.File` 当存储合同。
- excluded: R5 双路径 live 验收；目标 `status`/`progress`（本意见不改）；未复跑 `apps/api` 全量 `go test ./...`（以 E-001 自述 + 本轮定向复跑为准）。
- 工作区：`workspace-014-object-storage`（`root_goal: GOAL-001-object-storage`，`canonical_scope` 匹配，`shared_materials_catalog: none`）。
- 对照提交：`8aa0abc`。被审核查声明：E-001。对照意图：VP-014 意图 1 / 方向级退出判据 1；GOAL-002 D-001 端口公共类型零路径零 `os.File`。
- P-005：GOAL-005 无自有信息项。Root I-001～I-003 / I-005 已 verified；I-004 `recorded`（non-blocking，R5 叙事）。无到期 required 信息项阻断本门禁。
- 既有意见：A-001 self `pass`（编排器已落盘，开放 required = 0）。本意见独立形成，不采纳其 verdict。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| `*os.File` 仓库源码零引用 | 独立 `rg '*os.File'` 于 `*.go/*.ts/*.tsx/*.js/*.mjs/*.py/*.sh`：无命中 |
| `uploadDir` 仅 *_test.go | 独立 `rg uploadDir`：仅 `handler/upload_test.go`、`testhelpers_test.go`、`filelibrary_test.go` |
| 三类 Handler 公共签名已是端口，无 `dir string` | `RegisterUpload` `upload.go:219`；`FileLibraryRoutes` `filelibrary.go:168`；`ImportRoutes` `import.go:38`；`NewBrandingAssetStore`/`NewAvatarAssetStore` `raster_assets.go:85,94`；`handler/*.go` 生产文件 `dir string` 零命中 |
| 模块构造器：filelibrary/datatransfer 取 `kernel.ObjectStore`；account/settings 取 raster 包装 | `filelibrary/provider.go:33`；`datatransfer/provider.go:32`；`account/provider.go:33`；`settings/provider.go:34` |
| 装配单次 `newObjectStore` 注入上传 / brand / avatar / 两模块 | `composition.go:311-321,327-328,346-347,403,409` |
| 未知 driver 二次校验：非 `s3`/`local`/`""` 显式报错 | `composition.go:590-610`（`8aa0abc` 相对前一提交 +6 行） |
| 配置面仍 fail-closed 拒绝未知 driver | `config.go:509-526,913-925`；`config_objects_test.go:53-60` |
| 边界：`OpenWithCatalog` 打开 SQLite 文件 | `store.go:29-36,39-45`（`sql.Open("sqlite", path)`） |
| 边界：monitoring `dbPath` 仅 `os.Stat` 取 `DBSizeBytes` | `systemmonitoring/provider.go:36,69`；`systemmonitoring.go:84,120-133`；JSON 不回显路径 |
| VP-014 不含数据库端口；公共契约改对象端口 | VP-014 意图 1、与 VP-013 边界行；GOAL-002 D-001 §1 |
| pkg/web/脚本无对象存储路径合同 | `apps/api/pkg/version/version.go` 仅版本变量；web 走 `/api/account/avatars/`、`/api/upload`；`scripts/` 无 `uploadDir`/`os.File` |
| 本轮定向复跑绿 | `go test ./internal/composition -run TestNewObjectStoreWiring`；`./internal/config -run TestObjects`；`./internal/objectstore`；`./internal/handler`（162s）；`go build ./...` — 均 exit 0 |

## 对照成功标准

| 焦点 | 状态 | 证据 |
|------|------|------|
| (a) E-001 扫描可独立复现；漏扫维度 | 结论可复现；命中集合声明过满 | 三维结论成立；正则独立复现 3 处而非 1 处（R-001）；pkg/web/脚本独立干净（N-001） |
| (b) 边界声明成立 | 达成 | OpenWithCatalog / monitoring dbPath / OpenOptions.Path 均属 Store 方言或配置派生，不是三类落盘 Handler/模块合同 |
| (c) driver 二次校验正确 | 代码达成；测试未锁新分支 | `composition.go:600-605` 与 Load 白名单一致；`TestNewObjectStoreWiring` 未覆盖未知 driver（R-002） |
| (d) 其他把本地路径当存储合同的公共契约残留 | 未发现阻断项 | Handler/模块签名干净；`NewLocal(root)` 为适配器缝（N-002）；`os.ErrNotExist` 映射为遗留错误方言（N-003）；`FormFile` 为 HTTP 入口（N-006） |
| 检查点 4 测试绿 | 本轮定向绿；全量以 E-001 自述为准 | 见成果表末行；N-004 |

## verdict

**pass**

## findings

| 编号 | 级别 | 主张 | 证据 | 建议处置 |
|------|------|------|------|----------|
| R-001 | recommended / low | E-001「导出函数路径参数于 internal 全域仅命中 `store.OpenWithCatalog`」按所给正则独立复现不成立：同形还命中测试助手 `OpenStore` 与测试构建 `OpenSeeded`。两者仍是 SQLite 打开，不是对象存储合同，故扫描**结论**成立，**命中集合**过满。 | 命令：`rg -n --glob '*.go' 'func [A-Z]\w*\([^)]*(dir\|path)[^)]*string' apps/api/internal` → `testsupport/store.go:15`、`store/migration_catalog_test.go:30`、`store/store.go:31`。E-001 第 1 条。 | 补记 E-001 命中集合（含「测试助手 / *_test.go」）或收窄分母为「非测试、非 testsupport」。不阻断 R4→R5。 |
| R-002 | recommended / low | R4 唯一代码改动（未知 driver 拒绝）没有测试锁。`TestNewObjectStoreWiring` 覆盖 local / 空 driver / local-root override / s3，不覆盖 `"gcs"` / `"S3"` 手搭 Config。实现与 Load 白名单一致，审查正确。 | `objectprobe_test.go:15-70`；`composition.go:591-605`；对照 `config.go:509-526`。`8aa0abc` 只改 `composition.go`，未改测试。 | 加 subtest：未知 driver 返回 error 且 store/probe 为 nil。 |
| N-001 | note | E-001 扫描分母是 `internal`。独立补扫 `apps/api/pkg/`、`apps/web/src/`、仓库 `scripts/`：无对象存储路径/`os.File`/`uploadDir` 公共契约。同意 A-001 N-301。 | `pkg/version/version.go`；web `auth-context.test.tsx:82` 等为 API URL；`scripts/` `uploadDir` 零命中。 | 留痕。 |
| N-002 | note | `objectstore.NewLocal(root string)` 与 `LocalStore.Root()` 携带文件系统根，属本地适配器实现缝，不是 Handler/模块公共契约。`store.OpenOptions.Path` 是 SQLite 路径（postgres 时仅作文件根派生输入），属 VP-013 Store 方言 + VP-014 配置约定。 | `local.go:41-46`；`store/options.go:11-18`；`composition.go:606-610`；VP-014「默认仍由 filepath.Dir(db.path) 派生本地根」。 | 无需动作。 |
| N-003 | note | 三类 load 路径仍把 `kernel.ErrObjectNotFound` 映射为 `os.ErrNotExist` 再判 404。HTTP 面返回 `FILE_NOT_FOUND`/`ASSET_NOT_FOUND`，不把 `*os.File` 或本地路径暴露为合同。 | `upload.go:155-163`；`raster_assets.go:287-297`；`filelibrary.go:233-237,285-288`。 | 可选后续改为直接 `errors.Is(..., kernel.ErrObjectNotFound)`。非本门禁。 |
| N-004 | note | 检查点 4 全量绿以 E-001 自述为准。本轮独立复跑 wiring/config/objectstore/handler 全包与 `go build ./...` exit 0；未跑 `apps/api` 全量 `go test ./...`。 | E-001「验证」段；本轮命令输出。 | 编排器若要以检查点 4 放行，确认提交前全量套件仍绿即可。 |
| N-005 | note | HEAD 执行索引仍写「待登记」，而同提交已有 E-001；工作区 `goal-tree.md` 尚未列入 GOAL-005（仍写「届时创建」）。属台账包装，不改变扫描结论。本意见不改 tree。 | `git show HEAD:.../02-execution.md` 索引空行；`goal-tree.md:20-27`；工作树 `02-execution.md` 已补 E-001 行。 | `/govern` 响应时登记索引并同步 goal-tree。 |
| N-006 | note | `POST /api/upload` 仍 `r.FormFile("file")`（stdlib `multipart.File`，底层或为 `*os.File` 临时文件）。随后 `io.ReadAll` + 端口 `Put`，存储合同是 `[]byte` + `kernel.ObjectStore`。 | `upload.go:262-276,137`。 | 无需动作。 |
| N-007 | note | 同意 A-001 N-302：`RasterAssetStore` / `uploadStore` 类型名含 Store，已是端口包装。重命名会扩大 diff，不进 R4 分母。 | `raster_assets.go:43-50`；`upload.go:117-121`。 | 留痕。 |

## 必改项汇总

无。开放 required = 0。无到期 required 信息项。

## 与既有意见的异同

- 同意 A-001：三类公共签名已走端口；`*os.File` 零引用；`uploadDir` 仅测试；OpenWithCatalog / monitoring dbPath 属 Store 方言；`newObjectStore` 二次校验代码正确；N-301/N-302。
- 不同意 A-001「三维扫描全部符合声明、可独立复跑」的字面：命中集合少记两处测试/助手 SQLite 打开（R-001）。不升 required，因为结论与边界仍成立。
- A-001 未记录未知 driver 分支无测试锁（R-002）与 HEAD 台账包装缺口（N-005）。

## 结论

R4 公共面收尾核查的实质门禁成立：Handler / 模块公共契约不再把本地路径或 `os.File` 当作三类落盘的存储合同；SQL 路径与 monitoring `dbPath` 确属 Store 方言；未知 driver 二次校验实现正确。E-001 命中集合需补记、未知 driver 缺测试锁，均为 recommended。开放 required = 0，R4→R5 可放行。

建议 `/govern`：响应本意见（补记 E-001 命中集合、可选补未知 driver 测试、同步执行索引与 goal-tree），闭合 recommended 或留痕后推进 R5。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。

---

## 编排器响应与闭合（/govern · 2026-08-21）

| finding | 处置 | 证据 |
|---------|------|------|
| R-001 | fixed：E-001 第 1 条命中集合补记（OpenStore / migration_catalog_test 两处测试侧 SQL 打开），分母说明修正；结论不变 | GOAL-004 E-001（A-002 R-001 补记标注） |
| R-002 | fixed：TestNewObjectStoreWiring 新增 unknown-driver subtest（"gcs" 与大小写敏感 "S3" 均断言 nil store/probe + error） | composition/objectprobe_test.go |
| N-001 | fixed：E-001 补"补充扫描"节采纳 pkg/web/scripts 独立补扫结论 | GOAL-004 E-001 |
| N-002/N-006/N-007 | 留痕不动作 | 本节记录 |
| N-003 | 留痕：load 处 os.ErrNotExist 中转属 HTTP 404 惯用桥，非合同暴露；如后续清理属纯美学 | 本节记录 |
| N-004 | fixed：提交前全量 go test ./... FULL_TEST_EXIT=0 已完成（8aa0abc 前）；本响应测试亦单跑绿 | E-001 验证段 + 本次复跑 |
| N-005 | fixed：E-001 索引已登记、goal-tree 已列 GOAL-005（随本响应提交） | 02-execution.md；goal-tree.md |

**门禁判定**：verdict pass、开放 required = 0 → R4→R5 门禁放行；GOAL-005 四检查点满足，结项（Root progress 4/5）。
