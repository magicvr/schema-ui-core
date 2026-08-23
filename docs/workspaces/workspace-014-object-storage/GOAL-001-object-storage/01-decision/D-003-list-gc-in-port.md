---
id: D-003-list-gc-in-port
title: 枚举能力进端口：List + Stat（I-005 闭合）
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-001-object-storage
version: 0.1.0
---

# D-003 · List/GC 归属：枚举进端口，GC/配额留应用层（闭合 I-005）

## 决定

内核对象存储端口在 VP 下限（put/get/delete/exists）之上**增加两个只读枚举原语**：

- `List(ctx, ns) → []id`：枚举一个命名空间内的全部对象 id。
- `Stat(ctx, ns, id) → {Size, Meta}`：取元数据与字节数，不读 body。

启动 GC（孤儿清理）、每用户配额扫描等**应用层语义保留在调用方**，经上述原语实现；不把 GC 本身做进端口。

## 理由（证据）

1. 现行 GC 与配额都依赖目录列举：composition.go:326-348 启动时 `GC(referenced)` 清理未引用 brand/avatar 资产；upload.go `quotaReached` 扫描 meta 文件统计 per-owner 文件数/字节。若端口无枚举能力，S3 兼容实现无法对等交付这两个行为——直接违反 VP-014 "两实现合同平等，不得残缺"（意图 3 / 首波冻结表）。
2. uploads 目录**没有 DB 引用集**（owner 记录在 meta 边车文件里，不在数据库），"仅靠 DB 引用集做 GC、List 留私有"的替代方案在 uploads 命名空间根本不成立。
3. 最小公约数可行：本地适配器 = `os.ReadDir`（现行做法）；S3 兼容 = ListObjectsV2 分页（MinIO/R2/AWS 全支持，属 I-001 公约数内的 API）。两者都能以 O(keys) 实现同语义。
4. 性能非本 VP 退出分母：现行 quotaReached 已是每次上传 O(files) 扫描（注释明示 admin 工具可接受）；S3 上同复杂度只是更慢，语义不变。

## 未选方案

- **List 不进端口、GC/S3 跳过**：S3 上孤儿永久堆积、配额失效 = 缩水实现，违反 VP 意图。
- **GC 进端口（如 `GC(referenced)` 方法）**：把应用层引用语义（site_settings/users 的 URL 字段解析）下沉到存储层，端口被迫理解 URL 形态；且 GC 策略可能按 namespace 演化。保持"端口给原语、应用定政策"。

## 影响

- R1 冻结的端口方法集 = Put/Get/Stat/Delete/Exists/List（见 workspace-014 GOAL-002 D-001）。
- R3 收口时，三类落盘的 GC/配额改写为经端口原语的实现，行为与今日逐点对应。
