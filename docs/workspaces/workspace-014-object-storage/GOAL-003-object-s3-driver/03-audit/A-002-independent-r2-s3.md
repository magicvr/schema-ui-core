---
id: A-002-independent-r2-s3
title: 独立交叉审计 · R2 S3 兼容接入（R2→R3 门禁）
source: independent
date: 2026-08-21
scope: GOAL-003 全部交付（S3 适配器 / 端口零改动 / 凭证面 / readyz 扩依赖 / API 子集 / 测试）
verdict: pass
status: recorded
parent: GOAL-003-object-s3-driver
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# A-002 · 独立交叉审计：R2 S3 兼容接入（verdict: pass）

## 范围与区间

- auditor: grok-build (grok-4.6 · reasoning high)
- type: stage
- covered: Root D-004/D-005 与 GOAL-003 D-001 对照 commit `1545134`；`kernel.ObjectStore` 零改动；S3 适配器与 LocalStore 端口语义；static-credentials-only；凭证泄露面；readyz 仅 `driver=s3` 扩依赖；API 子集（无 multipart/presigned/建桶）；stub/live 测试覆盖；go.mod 依赖引入面。
- excluded: R3 三类落盘改线、R4 公共面去 `os.File`、R5 双路径 live 验收、目标 `status`/`progress`（本意见不改）。
- 工作区：`workspace-014-object-storage`（`root_goal: GOAL-001-object-storage`，`shared_materials_catalog: none`）。
- 对照提交：实现 `1545134`。方案冻结 GOAL-003 D-001；被审裁定 Root D-004 / D-005。
- P-005：I-001 / I-003 已由 Root D-004 / D-005 `verified`；GOAL-003 无自有到期 required 信息项。

## 成果与证据

| 主张 | 证据 |
|------|------|
| R1 端口合同零改动：`1545134` 未触碰 `kernel/objectstore.go`；Ping 不进接口 | `git diff 1545134^..1545134 --name-only`（无 kernel）；`kernel/objectstore.go:103-110`；`s3.go:33,70-72`；`composition.go:316` 以方法值作探针 |
| 适配器语义与 LocalStore 对等：upsert / 幂等删 / 哨兵 / 升序 List / 缺失命名空间空切片 / 校验先于 IO | `s3.go:83-91,141-159,161-176,178-192,194-203,205-218,223-255`；对照 `local.go:108-154,156-186,188-199,218-249`；`s3_test.go:124-231` |
| API 子集未越界：仅 Put/Get/Head/Delete/ListV2 + HeadBucket | `s3.go:37-44`；`objectstore/` 无 CreateBucket/UploadPart/CreateMultipart/Presign/PutBucket |
| 凭证 static-only：显式 `StaticCredentialsProvider`，源码不调用 `LoadDefaultConfig` | `s3.go:50-67`；`go list` 的 objectstore 导入无 `config` |
| 适配器无日志；user metadata 仅 name/type/kind/owner；readyz 响应不回显探针错误 | `objectstore/*.go` 无 slog/log；`s3.go:93-113`；`health.go:113-119` |
| readyz 仅 `ObjectsDriver=="s3"` 扩依赖；缺省 local 路径仍走无探针注册链 | `composition.go:309-318`；`config.go:281,509-511`；`health.go:45-47,109-111`；`testhelpers_test.go:122` + `health_test.go:32-48` 仍 200 |
| 分页聚合有 stub 覆盖；miss→哨兵有 stub 覆盖；live 未设 env 则 skip | `s3_test.go:166-178,202-231,233-243`；`s3_live_test.go:15-22` |
| 本审复跑：`go test ./internal/objectstore/` 与 `TestReadyz*` 绿 | 2026-08-21 独立执行，exit 0 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| (a) 端口零改动且与 LocalStore 逐点等价 | 达成 | 上表；Put 走单对象 user metadata，比本地边车更原子，符合 D-004 |
| (b) 凭证无泄露面（日志 / 错误串 / user metadata） | 达成 | 适配器不打日志；`notFound` 不包装 SDK 原文；readyz JSON 无 err；meta 四键 |
| (c) static-credentials-only 与 D-005 一致 | 达成（实现） | 运行时禁用默认链；go.mod 残留未使用链模块见 R-002 |
| (d) readyz 仅显式 s3 生效、local 缺省零变化 | 达成 | composition 分支 + 既有 `RegisterWithReadiness` 测试仍绿 |
| (e) API 子集无 multipart/presigned/建桶 | 达成 | `s3API` 白名单；适配器不调其余 `*s3.Client` 方法 |
| (f) 失败路径与分页覆盖 | 大部分 | 分页有测；miss/校验/Ping 失败有测；传输错误与 HTTP 404 分支无测（R-003） |

## Findings

### R-001 · main.go 仍声称 S3 接线未落地，与 R2 readyz 事实矛盾（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `apps/api/cmd/server/main.go:41-48` 启动 Warn：「S3 wiring lands in workspace-014 R2; file storage still uses the local disk adapter」；注释仍写「readyz does not cover the backend」。`composition.go:305-317` 已在 `driver=s3` 时构造适配器并把 `Ping` 挂进 readyz。 |
| closure | — |

调用方改线确属 R3，警告后半句仍真。前半句与注释在 R2 后为假：显式 s3 时 readyz **已经** HeadBucket。运维若按该 Warn 判断「S3 未接入」会误诊 503。建议改为只声明三类落盘仍走本地，readyz 已覆盖后端。

### R-002 · go.mod 将直接依赖标为 indirect，并残留未使用的默认链模块（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | open |
| evidence | `apps/api/go.mod:16-33`：`service/s3`、`credentials`、`aws-sdk-go-v2`、`smithy-go` 均 `// indirect`，但 `s3.go:12-16` 直接导入它们。`go mod why -m`：`config` / `feature/ec2/imds` / `service/sso` / `service/sts` / `service/signin` 均为「main module does not need」。`service/internal/presigned-url` 是 `service/s3` 真传递（SigV4），不是公开 Presign API。 |
| closure | — |

运行时未走默认链（见成果表），故不构成 D-005 违约。`go mod tidy` 会提升直接依赖并删掉未使用链模块；当前清单与 import 图不一致，后续升级易把 IMDS/SSO 重新当成「需要的」。建议 tidy 后只保留实际图。

### R-003 · 失败路径与 composition 接线缺少自动化锁（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | D-001 测试策略要求覆盖「全部端口行为 + 错误映射」。`s3_test.go` 覆盖 miss 哨兵、校验零调用、List 分页、Ping 失败、空桶拒绝；**未**覆盖 Put/Delete/List 后端错误传递、Exists 非 404（403/500）与 `mapS3Error` HTTP 404 分支（`s3.go:134-137`）。`composition.go:310-317` 的 `ObjectsDriver=="s3"` 接线无 composition 测试；readyz 机制测在 `health_probes_test.go:17-51`，与 driver 分支脱节。 |
| closure | — |

分页覆盖充分（`s3_test.go:202-218`）。R3 把适配器接到写路径前，建议补传输错误与 driver=s3 探针接线各一例。

### N-001 · D-001 写 `LoadDefaultConfig`，实现用显式 `aws.Config`

`D-001-s3-adapter-freeze.md` 第 16 行要求 `config.LoadDefaultConfig` 注入 static credentials。`s3.go:57-66` 改为手写 `aws.Config{Credentials: NewStaticCredentialsProvider(...)}` + `s3.NewFromConfig`。此偏差**更贴** D-005「禁用默认链」，不构成违约。方案正文未同步。

### N-002 · stub 分页 token 形不等于 S3 ContinuationToken

与 A-001 N-101 同结论。`s3_test.go:87-108` 用上一页最后 key 当 token；适配器只透传 `NextContinuationToken`（`s3.go:247-249`），不解析内容。无行为影响。

### N-003 · live 集成默认 skip；全量 `go test ./...` 台账未闭合

`s3_live_test.go:20-22` 未设 `S3_TEST_*` 则 skip。E-001 正文把全量测试指到不存在的 E-002；索引却写「全量测试绿」。本审只复跑 objectstore 与 `TestReadyz*`。真实 MinIO/R2 round-trip 属 R5。

### N-004 · Put/Delete/List/Ping 对非 miss 错误 `%w` 包装 SDK 原文

`s3.go:75,156,200,234`。当前 readyz 丢弃探针 err，适配器无日志，故 R2 无泄露。R3 若 `log.Error("%v", err)` 可能打出 Access Key Id 类 SDK 句。不必改映射；调用方勿把 SDK 原文打进日志。

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同

- A-001 self `pass`、开放 required = 0。本意见同意：端口零改动、语义对等、API 子集、static-only 运行时、readyz 条件扩依赖、N-101/N-102。
- 增补：self 未写 main.go 过期 Warn（R-001）、go.mod 未 tidy / 残留默认链模块（R-002）、传输错误与 composition 接线测试缺口（R-003）、D-001 与实现客户端构造偏差（N-001）。
- 无冲突：不把上述升为 required；不反对 R2→R3。

## 结论 + 建议给编排器/用户的下一步

R2 主张可核对，I-001/I-003 门禁仍关闭，无开放 required。可用 `/govern` 响应本意见：R-001～R-003 可在 R3 接线前顺手处理或书面跟踪；不必为这些项阻断 R2→R3。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。

---

## 编排器响应与闭合（/govern · 2026-08-21）

| finding | 处置 | 证据 |
|---------|------|------|
| R-001 | fixed：main.go 警告文案改为「readyz 已覆盖后端探针，三类落盘仍本地至 R3」 | cmd/server/main.go（A-002 R-001 注释块） |
| R-002 | fixed：go mod tidy，直接依赖归位、未用链模块移除 | apps/api/go.mod diff |
| R-003 | fixed：传输错误注入 + TestS3TransportErrorsPropagate；newObjectProbe 提取 + TestNewObjectProbe（driver=s3 接线锁） | objectstore/s3_test.go；composition/objectprobe_test.go |
| N-001 | fixed：D-001 §1 文本同步实际实现 | GOAL-003 D-001 |
| N-002 / N-004 | 留痕不动作（N-002 无影响；N-004 转为 R3 接线约束） | 本节记录；R3 立项时引用 |
| N-003 | fixed：全量套件证据补记于 E-002；live 证据归 R5 | GOAL-003 02-execution/E-002 |

**验证**：go build exit 0；objectstore/composition 测试全绿；全量 go test ./... exit 0。

**门禁判定**：verdict pass、开放 required 0 → R2→R3 门禁放行；GOAL-003 三检查点满足，结项（Root progress 2/5）。
