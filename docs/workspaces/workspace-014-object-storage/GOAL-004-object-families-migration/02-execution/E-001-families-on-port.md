---
id: E-001-families-on-port
title: R3 实施——三类落盘全部改走端口
status: recorded
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-004-object-families-migration
version: 0.1.0
---

# E-001 · R3 实施：三类落盘收口

按 [D-001](../01-decision/D-001-r3-migration-freeze.md) 完成：

## 迁移事实

1. **端口加性演化**：`ObjectInfo` 增加 `ModTime`（file-library created 面）；local 用 os.Stat、S3 用 HeadObject LastModified。方法集零变化。
2. **RasterAssetStore**：构造器 `dir`→`(objects, ns)`；save/load/Delete/DeleteOrphanOwnedBy/CountOwner/DeleteAll/GC 全部经端口；HTTP 面与 settings/account 模块合同零变化。
3. **uploadStore**：dir→objects；Put 原子写（W7 F-013 由适配器承接）；quotaReached→List+Stat 保守计数。
4. **file-library**：FileLibraryRoutes 参数 uploadDir→kernel.ObjectStore；scan→List+Stat（空 meta 行跳过=遗留无 meta body 不可见）；DELETE 先 Exists 再 Delete，幽灵 id best-effort 清理后回 404（D-001 §5 差异）。
5. **data-transfer import**：ImportRoutes 参数同改；loadUploaded 经端口 owner 校验；LoadUploadedFile 导出包装删除（无外部使用方）。
6. **modules**：filelibrary.New / datatransfer.New 的 uploadDir→kernel.ObjectStore（模块只消费端口）。
7. **composition**：`newObjectStore(cfg)` 构造唯一实例（local root = cfg.ObjectsLocalRoot ∥ filepath.Dir(db.path)；s3 显式）+ readyz 探针；main.go 的"s3 未接线"警告移除（接线完成，前提消失）。

## 验证证据

- go build ./... exit 0；go test ./... -run XXXNONE 全测试二进制编译通过。
- grep 复核：handler 生产代码三类路径无 os.ReadDir/os.WriteFile 残留（仅 *_test.go 直盘断言合法保留）。
- 全量 `go test ./...` exit 0（handler 140s / composition 32s 含单实例装配回归；两处测试直盘路径补 namespace 子目录修正）。

## 边界

R4（公共面去 os.File 收尾核查）未开始——RasterAssetStore/uploadStore 内部已无 os.File 合同，模块公共契约的 dir 参数已随 R3 清零。
