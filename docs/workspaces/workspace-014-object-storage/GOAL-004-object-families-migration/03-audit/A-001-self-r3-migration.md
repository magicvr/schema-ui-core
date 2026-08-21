---
id: A-001-self-r3-migration
title: 自审 · R3 三类落盘收口
source: self
date: 2026-08-21
scope: GOAL-004 全部交付（三类落盘端口化 / 单实例装配 / 行为保持 / 测试）
verdict: pass
parent: GOAL-004-object-families-migration
version: 0.1.0
---

# A-001 · 自审：R3 三类落盘收口（verdict: pass）

## 合同与证据指回

- **持久化全经端口**：grep 复核 handler 生产代码无 os.ReadDir/os.WriteFile 于三类路径（仅 *_test.go 直盘断言）；RasterAssetStore/uploadStore/fileLibraryEntity/importHandler 全部持有 kernel.ObjectStore。
- **行为保持**：
  - GC：List+Delete 未引用集 = 原 ReadDir+Remove 语义（Root D-003 授权）。
  - 配额：List+Stat，ghost/坏 meta 保守计 maxUploadBytes（A-002 N-002 处置落实）。
  - owner 门禁：upload GET / library confirm / import 三处 owner 校验逐点保留。
  - 遗留兼容：本地布局字节不变；legacy 直写用例（upload_test legacy fixture）原样通过。
- **加性微演化**：ObjectInfo.ModTime 为结构体字段新增，方法集零变化；两适配器同仓同改（D-001 §4）。
- **装配**：newObjectStore local 分支 root=ObjectsLocalRoot∥filepath.Dir(db.path)（与现行派生一致）；s3 分支复用 GOAL-003 适配器 + HeadBucket 探针；接线锁测试 TestNewObjectStoreWiring。

## Findings

| 编号 | 级别 | 内容 | 处置 |
|------|------|------|------|
| N-201 | note | CountOwner 对 ghost id 无条件 +1（旧实现读 meta 判 owner 后 +1）——两方向都保守，方向略变宽。 | 留痕；随独立审复核 |
| N-202 | note | file-library scan 对空 meta 行跳过 = 遗留无 meta body 不可见（D-001 §2 声明）；S3 上不存在该形态（meta 内嵌）。 | 无需动作 |
| N-203 | note | main.go 警告整段移除——若未来出现新的"配置了但未生效"窗口需重建同类提示。 | R4 复核 |
| — | required | 无 | 开放 required = 0 |

## 过程留痕

实施中一处测试修正（stray namespace 子目录）在提交前完成并全量复跑绿（吸取 R2 a4c68ef 教训：全量验证先于提交）。

## 结论

R3 自审 pass，开放必改 0。已并行发起独立审计（grok build · grok-4.6 · high），意见落盘后由编排器响应。
