---
id: E-001-sweep-evidence
title: R4 核查——公共面扫描证据 + driver 二次校验加固
status: recorded
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-005-public-surface-sweep
version: 0.1.0
---

# E-001 · R4 扫描证据与加固

## 扫描证据（2026-08-21，commit d99221f 之后的工作树）

1. **导出函数路径参数**：正则 `func [A-Z]\w*\(...(dir|path)... string` 于 internal 全域仅命中 `store.OpenWithCatalog(path, ...)`——SQL 数据库打开（VP-013 Store 方言），**非对象存储合同**。
2. **`*os.File`**：internal 非测试代码**零引用**。
3. **`uploadDir` 残留**：仅 *_test.go 直盘断言（布局字节兼容承诺的合法使用面）。
4. **模块构造器清点**：datatransfer.New / filelibrary.New 已取 kernel.ObjectStore；account/settings 取 handler 的 RasterAssetStore 包装（无路径暴露）；systemmonitoring.New 的 dbPath 为其自有 SQLite 监控库——Store 方言边界，同第 1 条。
5. **加固落地（GOAL-003 A-002 N-005）**：composition.newObjectStore 对非 local/s3 的未知 driver 显式报错，杜绝手搭 Config 静默落到本地适配器。

## 边界声明

本核查的对象存储合同 = 三类第一方落盘（avatars / brand-assets / uploads）的 Handler 与模块公共契约。数据库文件路径（db.path、monitoring.db）属持久化方言（VP-013 已 closed），不在本 VP 扫描分母。

## 验证

go build exit 0；TestNewObjectStoreWiring 绿；全量 go test ./... 于提交前复跑（结果记于提交 hash 对应 checkpoint）。
