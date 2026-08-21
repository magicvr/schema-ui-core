---
id: D-001-r3-migration-freeze
title: R3 收口迁移方案冻结
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-004-object-families-migration
version: 0.1.0
---

# D-001 · R3 收口迁移方案冻结

依据：Root D-002/D-003、GOAL-002 D-001（端口合同）、A-002 N-002/N-004 约束。端口方法集不改；两处**加性微演化**见 §4。

## 1. 单实例装配

composition 构造唯一 `kernel.ObjectStore`：driver ∈ {local, ""} → `objectstore.NewLocal(root)`，root = `cfg.ObjectsLocalRoot`（空则 `filepath.Dir(cfg.DBPath)`）；driver=s3 → `S3Store`（复用 newObjectProbe 改造为 newObjectStore 返回 实例+探针）。全部消费方注入同一实例。

## 2. 各落盘面迁移

- **RasterAssetStore**：构造器 `dir string` → `kernel.ObjectStore` + 固定 ns（brand-assets / avatars）。save→Put、load→Get、Delete/DeleteOrphan/DeleteOrphanOwnedBy→Delete+Get(owner)、CountOwner/GC/DeleteAll→List+Stat(+Delete)。HTTP 面（重编码管线、URL 前缀、AssetIDFromURL、权限注释）不动——settings/account 模块继续拿 `*handler.RasterAssetStore`，模块合同零变化。
- **uploadStore**：dir→objects。save→Put(meta Name/Type/Owner)；load→Get 返回类型化 ObjectMeta（调用方同步改字段访问）；quotaReached→List+Stat，保守计数保留（Stat miss/错误的 id 计 files+maxUploadBytes，沿 A-002 N-002 处置）。RegisterUpload 第三参 dir→store。
- **file-library**：`FileLibraryRoutes(a, uploadDir, ...)`→`(a, objects kernel.ObjectStore, ...)`；entity.scan→List+Stat（Created 取自 §4 的 ModTime）；DELETE 路由先 Exists 再 Delete，!Exists 时仍做一次 best-effort Delete 清理孤儿边车后回 404（孤儿边车清理获 A-002 N-002 授权；204→404 的幽灵边车差异记录于 §5）。无任何 meta 字段的行跳过（遗留无 meta body 保持不可见，行为对齐旧 scan 只看 .meta.json）。
- **data-transfer import**：`ImportRoutes(a, repo, operations, uploadDir, moduleID)`→`(a, repo, operations, objects, moduleID)`；loadUploadedFile→objects.Get(uploads ns, id)；LoadUploadedFile 导出包装删除（唯一调用方在本包内改造）。
- **modules**：filelibrarymodule.New / datatransfermodule.New 的 uploadDir 参数→kernel.ObjectStore（模块只消费端口——VP 意图 1）。

## 3. main.go 警告退场

R3 后 driver=s3 即全链路生效，"s3 未接线"警告失去事实基础：整段移除（readyz 探针保留）。

## 4. 端口加性微演化（不破坏 R1 冻结面）

1. `ObjectInfo` 增加 `ModTime time.Time` 字段：file-library 列表/详情的 created 排序与展示需要；两个适配器同仓同改，方法集零变化，非发布 API。S3 用 HeadObject LastModified。
2. 无其他改动。Put/Get/Stat/Delete/Exists/List 签名原样。

## 5. 已知行为差异（有意、留痕）

- 幽灵边车 DELETE：旧实现回 204 并清理；新实现回 404 且 best-effort 清理（§2）。仅影响前端口时代的残损数据。
- **幽灵边车列表/配额组合（A-002 F-001 补记）**：仅 sidecar、无 body 的残损 id——旧 scan 会入列（可点删），新 scan 因 Stat 失败而**不入列**；但同一 id 仍被 List 出并在配额扫描与 CountOwner 中**保守计入**。组合效果："列表看不见、配额仍占、UI 无法点删"。范围=前端口时代残损数据（端口自身不会产生该形态）；清理路径=按 id 直接 DELETE（会做 best-effort 边车清理）或运维盘面清理。
- 坏 JSON 边车 + body 在（A-002 R-004 补记）：旧配额/CountOwner 对读失败的 meta 保守 +1；新 parseMeta 容忍损坏返回空 meta，owner 为空被跳过、body 正常计入字节。方向由"保守多计"变"按实计"，远程不可利用；如需恢复保守口径在 R4 一并评估。
- 无 meta body 在库列表中依旧不可见（跳过空 meta 行），配额扫描则保守计入——两处口径不同系沿用现状。
- CountOwner 对 ghost id 无条件 +1（旧按 owner 匹配后 +1）——保守方向略变宽（A-001 N-201 / A-002 N-001）。
- driver=s3 上线即新写走 S3；既有本地文件不迁移（Root 信息表 I-004 用户裁决：运维自备）。

## 6. 测试策略

- 既有 handler/module 测试改注入 `NewLocal(t.TempDir())`；直接磁盘操作的用例（legacy 兼容、quota）经 \<root\>/uploads 布局保持有效（布局字节兼容是 R1 承诺）。
- 新增：composition 单实例断言（local 与 s3 两分支）、RasterAssetStore 经端口的 GC/CountOwner 行为测试。
