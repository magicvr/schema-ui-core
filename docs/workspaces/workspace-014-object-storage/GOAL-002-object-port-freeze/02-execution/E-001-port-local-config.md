---
id: E-001-port-local-config
title: R1 实施——端口类型 / 本地适配器 / 配置面（全量测试绿）
status: recorded
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-002-object-port-freeze
version: 0.1.0
---

# E-001 · R1 实施：端口 + 本地适配器 + 配置面

## 事实

按 [D-001](../01-decision/D-001-port-config-freeze.md) 实施并验证：

1. **端口**：`apps/api/internal/kernel/objectstore.go` —— `ObjectStore` 接口（Put/Get/Stat/Delete/Exists/List）、`ObjectNamespace` 三值枚举、`ObjectMeta{Name,Type,Kind,Owner}`、`ObjectInfo{ID,Size,Meta}`、`ErrObjectNotFound` 哨兵、`ValidObjectID`（32hex，同现行 uploadFileIDPattern）。公共类型零本地路径 / 零 `os.File`。
2. **本地适配器**：`apps/api/internal/objectstore/local.go` —— 布局 `<root>/<ns>/<id>`+`<id>.meta.json`（与现行磁盘逐点兼容）；边车只写非空键；Put 临时文件+rename、边车失败回滚 body；List 升序、缺失命名空间=空切片；编译期 `var _ kernel.ObjectStore = (*LocalStore)(nil)`。
3. **配置面**：`apps/api/internal/config/config.go` + 两份 YAML —— `storage.objects.{driver,local.root,s3.*}`；env 覆盖 `STORAGE_OBJECTS_*`；fail-closed：未知 driver / driver=local 带 s3 键 / driver=s3 缺 endpoint/bucket/access_key_id/secret_access_key 均 LoadError；`ValidateProd.validateObjects` 每环境复查。readyz 未动。
4. **operator YAML**：`apps/api/configs/config.yaml` 与内嵌 default 增加 storage 段（driver: local 缺省）。

## 验证证据

- `go build ./...` → exit 0。
- `go vet ./internal/{kernel,objectstore,config}/` → exit 0。
- `go test ./internal/kernel/ ./internal/objectstore/ ./internal/config/` → ok（含新增 local_test.go 四组用例：round-trip/upsert/幂等删除、NotFound 哨兵、ns+id fail-closed 且零磁盘副作用、遗留 meta 双形态兼容+损坏容忍、List 有序/跨命名空间隔离/缺失目录语义；config_objects_test.go 七组用例：缺省 local、s3 解析、env 覆盖、未知 driver、误配拦截、缺凭证、ValidateProd 复查）。
- 全量 `go test ./...` → exit 0（handler 154s / composition 31s 包含在内，无 FAIL）。

## Git checkpoint

- 方案冻结 docs：`d403832`（D-002/D-003 + GOAL-002 立项 + goal-tree）。
- 本切片实施提交：见 E-001 后续登记（commit 后补 hash 于索引表备注）。

## 边界声明

未接线 composition（R3）；未引入任何第三方依赖；未改 readyz；未触碰三类落盘调用方行为。
