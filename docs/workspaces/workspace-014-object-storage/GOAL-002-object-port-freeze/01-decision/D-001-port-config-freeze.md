---
id: D-001-port-config-freeze
title: R1 端口与配置面冻结
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-002-object-port-freeze
version: 0.1.0
---

# D-001 · R1 端口与配置面冻结

依据 Root D-002（单桶 + 命名空间前缀）与 D-003（List+Stat 进端口）。本决策冻结 R1 交付的合同面。

## 1. 端口类型（internal/kernel/objectstore.go）

- `ObjectNamespace` 枚举：`avatars` / `brand-assets` / `uploads`（Root D-002；其他值 fail-closed 拒绝）。
- `ObjectMeta{Name, Type, Kind, Owner string}`：统一两类现行 meta 形态（upload = name/type/owner；raster = type/kind/owner）。全字段可选字符串，由调用方决定必填性（如 uploads owner-only 语义留在 HTTP 面）。
- `ObjectInfo{ID, Size, Meta}`：Stat 返回体（配额扫描需要字节数，不读 body）。
- `ObjectStore` 接口：`Put / Get / Stat / Delete / Exists / List`。
  - Put = 幂等 upsert（现行写入方总是生成新随机 id，覆写语义不与现状冲突）。
  - Delete = 幂等（缺失返回 nil，与现行 RasterAssetStore.Delete 一致）。
  - Get/Stat 缺失 → `kernel.ErrObjectNotFound` 哨兵（errors.Is 可判；与 kernel.ErrNoRows 同风格）。
  - List 返回**升序 id 切片**（确定性：本地 ReadDir 有序；S3 ListObjectsV2 字典序）。
- id 规则冻结：`^[0-9a-f]{32}$`（与现行 uploadFileIDPattern 完全一致），端口层统一校验，杜绝路径/key 注入（fail-closed）。
- 公共类型**零本地路径、零 `os.File`、零 io.Reader 流式**：R1 全部 []byte 整体读写，与三类现行调用方一致（最大 8 MiB 上限已在 HTTP 面）；流式化不在本 VP 分母。

## 2. 本地盘适配器（internal/objectstore，缺省实现）

- 布局：`<root>/<namespace>/<id>` + `<root>/<namespace>/<id>.meta.json` 边车——**与现行磁盘布局逐点兼容，零迁移**。
- 边车序列化只写非空键（新写入与两类遗留形态字节兼容）；读取容忍缺失/损坏边车（与现行 load() 一致）。
- Put 原子性：body 临时文件 + rename；边车写失败则回滚删除 body（沿袭 W7 F-013 fail-closed）。
- List 对缺失目录返回空切片 + nil（与现行 GC/quota 的 IsNotExist 语义一致）。

## 3. 配置面（storage.objects）

```yaml
storage:
  objects:
    driver: local          # local（缺省）| s3
    local:
      root: ""             # 缺省 = filepath.Dir(db.path)（composition 派生，R3 接线）
    s3:
      endpoint: ""         # 必填（driver=s3）
      region: ""           # 可选；空 = 驱动缺省
      bucket: ""           # 必填
      access_key_id: ""    # 必填
      secret_access_key: "" # 必填；仅经 env 插值注入，不入库不入仓
      use_path_style: true # MinIO/R2 需要；AWS 可关
```

- env 覆盖：`STORAGE_OBJECTS_DRIVER / STORAGE_OBJECTS_LOCAL_ROOT / STORAGE_OBJECTS_S3_ENDPOINT / STORAGE_OBJECTS_S3_REGION / STORAGE_OBJECTS_S3_BUCKET / STORAGE_OBJECTS_S3_ACCESS_KEY_ID / STORAGE_OBJECTS_S3_SECRET_ACCESS_KEY / STORAGE_OBJECTS_S3_USE_PATH_STYLE`（沿用现行 envOr 层）。
- fail-closed 校验（Load + ValidateProd 每环境执行，沿 validateDB 先例）：
  - driver 属于 {local, s3}（空 = local）；
  - driver=local 时设置任何 s3.* 非空键 → LoadError（误配拦截）；
  - driver=s3 时 endpoint / bucket / access_key_id / secret_access_key 缺一 → LoadError（I-003 的"凭证 fail-closed"提前到配置面落地，键名即本节冻结值）。

## 4. 未选方案

- **io.Reader 流式端口**：三类调用方全部整块 []byte，流式是假需求；S3 SDK PutObject 支持 bytes.NewReader。R5 后如需大文件再议（不在分母）。
- **端口带 Content-Type 参数独立于 meta**：Type 已在 ObjectMeta，避免双通道漂移。
- **driver=s3 时容忍缺凭证、运行时报错**：违反 fail-closed 配置先例（db.postgres 密钥缺一即拒）。

## 5. 影响与边界

- readyz 本阶段**不变**（仅 driver=s3 显式配置后才扩依赖，R2 落地）。
- I-001（S3 驱动公约数）与 I-003（凭证注入机制）仍 open，最晚 R2 实施前关闭——本决策只冻结配置**键名**，不选 SDK。
- R3 接线时 composition 用 `cfg.ObjectsLocalRoot`（空则 `filepath.Dir(cfg.DBPath)`）构造唯一 LocalStore，三类调用方按命名空间消费同一实例。
