---
id: A-002-independent-r1-freeze
title: 独立交叉审计 · R1 端口与配置面冻结（R1→R2 门禁）
source: independent
date: 2026-08-21
scope: GOAL-002 全部交付（D-001 合同 / kernel 端口 / 本地适配器 / 配置面 / 三类调用方零迁移 / R2 可接入性）
verdict: conditional
status: recorded
parent: GOAL-002-object-port-freeze
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# A-002 · 独立交叉审计：R1 端口与配置面冻结（verdict: conditional）

## 范围与区间

- auditor: grok-build (grok-4.6 · reasoning high)
- type: stage
- covered: VP-014 意图对齐；D-001 冻结合同 vs `34db126` 实现；本地适配器与现行磁盘/边车兼容；`storage.objects` fail-closed 与 secret 面；端口是否足以支撑 R2 S3 接入；测试对冻结面的覆盖；对照 `raster_assets.go` / `upload.go` / `composition.go` 零迁移。
- excluded: R2 SDK 选型（I-001）、R3 接线、readyz 扩依赖、目标 status/progress（本意见不改）。
- 工作区：`workspace-014-object-storage`（`root_goal: GOAL-001-object-storage`，`shared_materials_catalog: none`）。
- 对照提交：方案冻结 `d403832`；实现 `34db126`。

## 成果与证据

| 主张 | 证据 |
|------|------|
| D-001 冻结面与 VP-014 意图一致：三命名空间、公共面无路径/`os.File`、本地缺省、S3 显式配置；List+Stat 由 Root D-003 补齐以保两实现平等 | VP-014 意图 1–4；D-002；D-003；`kernel/objectstore.go:27-41,93-110` |
| 本地布局 `<root>/<ns>/<id>` + `<id>.meta.json` 与现行三目录一致；`34db126` 未改三类调用方 | `composition.go:308-337`；`upload.go:135-147`；`raster_assets.go:277-283`；`git show --stat 34db126`（无 handler/composition） |
| 遗留双形态边车可读；损坏边车容忍 | `local.go:70-106`；`local_test.go:130-175` |
| 未知 driver / local+s3 键 / s3 缺凭证在 Load 失败 | `config.go:509-524,913-924`；`config_objects_test.go:53-77` |
| 端口六方法 + 升序 List + `ErrObjectNotFound` + 32hex id 足以让 R2 在适配器内实现，不必改合同 | `kernel/objectstore.go:59-110`；`local.go:218-249`；S3 PutObject/ListObjectsV2/NoSuchKey 可对等 |
| 相关包测试绿 | `go test ./internal/objectstore/ ./internal/config/ ./internal/kernel/` exit 0 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| (a) D-001 无越界/无缺失 | 达成 | 未引入签名 URL/流式/GC-in-port/多桶；VP 下限 + D-003 枚举均在端口 |
| (b) 本地兼容声明真实 | 达成 | 路径/id/边车键与现行写入方同构；零调用方迁移由 commit 范围证明 |
| (c) 配置 fail-closed 无漏洞 | 部分 | 配对规则成立，但误配错误串会打出 secret/凭证值 |
| (d) 端口可支撑 R2 无破坏性改动 | 达成 | 见上；I-001/I-003 仍属 R2 实施前信息项，不构成本冻结缺口 |
| (e) 测试覆盖冻结关键行为 | 大部分 | round-trip/upsert/幂等删/哨兵/校验/遗留 meta/List；缺缺失边车与 sidecar 回滚 |

## Findings

### F-001 · driver=local 误配错误把 s3 凭证值写入 LoadError（required）

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | closed-fixed（2026-08-21，见文末编排器响应） |
| evidence | `apps/api/internal/config/config.go:512-514` 将 `firstNonEmpty(...)` 的**值**（含 `ObjectsS3SecretAccessKey`）插入 `%q`；`config.go:884-892` 注释声称「name the offending s3 key」但实现返回的是值；`apps/api/cmd/server/main.go:30-32` `logger.Error("startup failed", "err", err)` 会把该串打进启动日志。`config_objects_test.go:62-68` 只覆盖 endpoint 误配，不覆盖「仅 secret 非空」路径。 |
| closure | fixed：`firstNonEmpty` 删除，改为 `firstSetS3Key`（只返回键名）；Load 与 `localS3KeyMisconfig` 统一只报 `storage.objects.s3.<key>`；补 `TestObjectsMisconfigErrorDoesNotLeakSecret`（secret-only 路径断言错误串不含值）。 |

仅设置 `STORAGE_OBJECTS_S3_SECRET_ACCESS_KEY`（或 YAML 字面 secret）而 `driver` 仍为缺省 `local` 时，错误信息与 slog 会带出 secret。D-001 §3 / VP-014 配置面要求 secret 不入库、不入仓；启动失败日志是另一条泄露面。建议：错误只报键名（`secret_access_key` 等），禁止插值；补「仅 secret/access_key 非空」测试。

### R-001 · ValidateProd 未复检 local+s3 误配（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | D-001:53-56 写明三条规则由 Load **+ ValidateProd** 执行；`config.go:508` 亦称 pairing「re-checked by ValidateProd」。`config.go:898-906` 对 `""`/`local` 直接 `return nil`，不调用 `firstNonEmpty`。生产路径因 `ValidateProd` 先包装 `LoadError`（`config.go:838-839`）而不构成绕过。 |
| closure | — |

建议 `validateObjects` 与 `validateDB` 同构，把 local 误配复查收进 ValidateProd，避免注释/合同与实现分叉。

### R-002 · R1 接受 driver=s3 但无消费者（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `config.go:517-521` 在凭证齐全时 Load 成功；`34db126` 未改 `composition.go:308-337`，三类落盘仍写本地目录；readyz 未扩。D-001 §5 已声明 R3 才接线。 |
| closure | — |

在 R2 适配器落地前，显式 `driver=s3` 会被当成已配置但实际仍落本地。建议 R2 接线前至少打明确日志，或 R1 对尚无适配器的 `s3` 失败（若选择后者需改冻结合同，须用户裁决）。

### R-003 · 冻结行为缺两则测试（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | D-001:31-32 要求「读取容忍缺失边车」与「边车写失败则回滚删除 body」。`local_test.go:164-174` 覆盖损坏边车，未覆盖缺失边车；无 sidecar `WriteFile` 失败回滚用例。 |
| closure | — |

### N-001 · Put 覆写后边车失败会删掉旧 body

与 A-001 N-001 一致。`local.go:147-151` 回滚 `os.Remove(bodyPath)`。现行写入方均新随机 id（`upload.go:127-131`、`raster_assets.go:269-273`）。S3 双对象同样无跨 key 事务；R3 接线评审确认调用方不复用 id。属合同内部分失败窗口。

### N-002 · List 可能列出无 body 的 sidecar-only id

`local.go:241-247` 对 `.meta.json` `TrimSuffix` 后若 id 合法即计入。`Exists`/`Get`/`Stat` 只看 body（`local.go:202-215,156-169,172-185`）。R3 配额若 `List` 后 `Stat`，须按 `ErrObjectNotFound` 跳过或 Delete 清理，不必改端口。

### N-003 · I-001 / I-003 仍 open，不构成本冻结缺口

Root 信息表：I-001/I-003 最晚 R2 实施前。D-001:66-67 只冻关键名、不选 SDK。空 `region`、YAML 字面 secret（与 `DB_PASSWORD` 同构，`config.go:909-912`）、`env.example` 未列 `STORAGE_OBJECTS_*` 均归 I-003/I-001，不在本门禁放行条件内。

## 必改项汇总

- **F-001**（required, ~~open~~ → **closed-fixed** 2026-08-21）：误配 LoadError 不再包含 secret/access key 值；详见文末编排器响应。

## 与既有意见的异同

- A-001 self `pass`、开放 required = 0。本意见同意其兼容性、端口完备性与 N-001/N-002 边界。
- 分歧：self 未审 `firstNonEmpty` 插值泄露面；本意见据此 **conditional**，F-001 未闭合前不得把本合同当作无条件 R2 输入。

## 结论 + 建议给编排器/用户的下一步

冻结合同与本地实现大体可核对，R2 不必破坏性改端口。先用 `/govern` 响应 F-001（fixed：改错误串只报键名并补测；或用户书面 residual/overruled）。R-001～R-003 可与 F-001 同批或顺延 R2。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。

---

## 编排器响应与闭合（/govern · 2026-08-21）

> 编排器对本意见全部 findings 的响应；required 走三路径之一闭合，recommended 同批处理并留证据。

| finding | 处置 | 闭合路径 | 证据 |
|---------|------|----------|------|
| F-001 (required) | **fixed**：删除返回值的 `firstNonEmpty`，新增 `firstSetS3Key`（只返回键名）与共享错误构造 `localS3KeyMisconfig`；Load 与 validateObjects 均只报 `storage.objects.s3.<key>`，零值插值 | fixed | config.go（firstSetS3Key/localS3KeyMisconfig）；新增用例 TestObjectsMisconfigErrorDoesNotLeakSecret（secret-only 断言错误串含键名且**不含**值） |
| R-001 | **fixed**：validateObjects local 分支补 s3 键复查，与 Load 同措辞、与 validateDB 同构；补 ValidateProd 复查用例 | fixed（recommended 提前实施） | config.go validateObjects；TestObjectsMisconfigErrorDoesNotLeakSecret/subtest 2 |
| R-002 | **fixed**：main.go 在 ValidateProd 通过后对 `driver=s3` 打 Warn——明示 S3 接线在 workspace-014 R2 落地前文件存储仍走本地适配器、readyz 未扩 | fixed（选日志方案；"driver=s3 直接失败"会改冻结合同，不做） | cmd/server/main.go（A-002 R-002 注释块 + logger.Warn） |
| R-003 | **fixed**：补缺失边车 Get/Stat 容忍用例（TestLocalMissingSidecarTolerated）与边车写失败回滚用例（TestLocalPutRollsBackWhenMetaWriteFails，sidecar 预建为目录使 WriteFile 确定性失败） | fixed（recommended 提前实施） | internal/objectstore/local_test.go 新增两用例 |
| N-001 | 认可（合同内部分失败窗口）；R3 接线评审时核对调用方不复用 id | 留痕不动作 | 本节记录；R3 子目标立项时引用 |
| N-002 | 认可；R3 配额/GC 经 List+Stat 时按 ErrObjectNotFound 跳过或清理 sidecar-only id | 留痕不动作 | 本节记录；R3 子目标立项时引用 |
| N-003 | 认可；env.example 已顺带补 STORAGE_OBJECTS_* 注释模板（I-003 范围内的文档项），I-001/I-003 保持 open 至 R2 实施前关闭 | 留痕+顺手项已做 | configs/env.example object storage 段 |

**验证**：go build ./... exit 0；go vet（config/objectstore/cmd/server）exit 0；go test ./internal/{config,objectstore,kernel}/ 全绿。修复提交见 GOAL-002 执行台账 E-002（commit b6ac23e）。

**门禁判定**：F-001 以 fixed 合法闭合后，本意见开放 required = 0；R1→R2 门禁放行。verdict 由 conditional 对应的条件已满足，编排器据此推进 R1 收尾（GOAL-002 done、Root progress 1/5）。
