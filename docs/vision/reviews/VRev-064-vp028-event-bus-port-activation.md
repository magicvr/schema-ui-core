---
id: VRev-064-vp028-event-bus-port-activation
doc_type: vision-review
title: VP-028 激活就绪 · 进程内事件总线端口（架构分支）
source: self
date: 2026-09-01
scope: VP-028-event-bus-port 意图 / 退出判据 / 非目标 / P-005 / 架构类 freshness（`29727510`）
verdict: pass
open_required: 0
status: active
created: 2026-09-01
updated: 2026-09-01
parent: null
version: 0.1.0
---

# VRev-064 · VP-028 激活就绪（进程内事件总线端口）

## 背景与触发

用户 2026-09-01 指令：「/vision 激活 vp-028，然后交 /govern 开设新工作区」。VP-028（进程内事件总线运输端口 · 架构分支 · 承接 RT-Q02 运输端口前置）经 VRev-058（self `pass`）与 VRev-059（grok build independent `conditional` → V-F101/102/103/104 响应 **fixed** 2026-08-31 · 开放 required=0）审视修订至 v0.1.1，本次为激活就绪审视。

## 审视要点

### 1. 意图与退出判据可判定性

**pass**。VP-028 v0.1.1 八条退出判据（端口契约冻结 / 进程内实现可用 / 接缝声明落盘 / 对齐登记不解除 Admin gated / 共享约定登记 / 停机与边界语义 / 边界保持 / 审计闭合）均可核验；VRev-059 响应已落入正文（V-F101 定位统一为运输端口 + 不解除 Admin typed domain event gated + EventBus ≠ Job、V-F102「不消耗 RT-Q02 trigger」解释规则、V-F103 序列化取舍改为 R1 显式取舍 + I-028-002 缓冲满最小语义、V-F104 异步投递须声明 SIGTERM 取消订阅/排空否则同步投递、V-F100 共享约定收窄为 topic/订阅命名 + 契约测试 harness、不纳入 Redis key 轨道）。

### 2. 非目标与红线

**pass**。不预制 outbox / 外部 broker（不引入客户端依赖、不预裁 RT-Q06 表结构）；不重开历史 VP；不改 Profile 默认集 / Manifest（VP-008 `go` 红线）；不解除 Admin 功能分支 typed domain event 扩展接缝的 trigger-gated；缓存/限流语义归 VP-026/027 独立交付；业务域事件产品语义归业务域 VP。

### 3. 信息需求（P-005）

**pass**。I-028-001（类型化机制：接口断言 vs 注册表；含可序列化约束取舍）/ I-028-002（投递语义默认 + 缓冲满最小语义）/ I-028-003（handler 错误语义）required → R1 裁决；I-028-004（事件类型注册权属）non-blocking，若 I-028-001 选注册表则升 required。均标注最晚阶段，无伪装已知。

### 4. 架构类轻量 freshness（`5744868d` → `29727510`）

**PASS**，不暂挂 `go`：

| 域 | 变更 | 判定 |
|----|------|------|
| 协议 pin / provenance（`apps/web/src/protocol/upstream`） | 零变更 | ✅ |
| 依赖锁（go.mod / go.sum / package.json / lockfiles） | 零变更 | ✅ |
| 迁移台账（`migrate.go` / `modules/*/migration`） | 零变更 | ✅ |
| Profile 装配（`kernel/profile.go`） | 零变更 | ✅ |
| 区间代码变更 | `kernel/ratelimit.go` / `internal/ratelimit` / `handler/*` 使用点迁移 / `composition` / `cache-redis-seam-and-track.md`（全部为 VP-027 已审结目交付——R1～R4 落地 + 关门双审 + VRev-063 `pass` · 2026-09-01） | ✅ 不涉及内核端口面以外契约 / Store / 模块契约 |

消费候选 HEAD `29727510`（workspace-027 R4 关门 commit；freshness 三字段落 Root D-001）。本 VP 属架构分支、非业务域 VP，H-002「业务域 VP 激活前确认同进程主要形态」发现机制不适用；H-002 仍为 Charter 冻结假设。

### 5. 组合对齐

**pass**。VP-028 `vision_ref` = `schema-ui-core-admin-foundation@0.4.0` 精确匹配现行 Charter；roadmap 组合表第 28 行 / RT-Q02 承接注记同步；lead `workspace-028-event-bus-port` 命名沿用 VP slug 惯例（VP-013～027 一致形态）；三端口第三个——不属 Redis 轨道（owner 仍为 `docs/architecture/cache-redis-seam-and-track.md`）；共享约定为本 VP 的 topic/订阅命名 + 契约测试 harness。

## Verdict

**pass**（0 required）。

VP-028 意图/判据/非目标/信息需求已就绪，架构类 freshness PASS（不暂挂 `go`），可激活并交 `/govern` 开区。

## Findings

### 必改（required）

无。

### 建议（recommended）

无新增（VRev-058/059 全部 findings 已闭合）。
