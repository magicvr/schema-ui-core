---
id: VRev-059-vp026-028-ports-independent
doc_type: vision-review
title: VP-026/027/028 端口 VP 计划独立复审 · 同进程基座基础设施早期化
source: independent
date: 2026-08-31
scope: VP-026-cache-port / VP-027-rate-limiter-port / VP-028-event-bus-port（planned · 三端口意图、退出判据、非目标、P-005、RT-Q02/Q03/Q05 衔接与跨 VP 约定）
verdict: conditional
open_required: 1
status: active
created: 2026-08-31
updated: 2026-08-31
parent: null
version: 0.1.0
---

# VRev-059 · VP-026/027/028 端口 VP 计划独立复审（2026-08-31）

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | grok-4.6（思考强度 xhigh · `/vision-audit`） |
| scope | VP-026 / VP-027 / VP-028 planned 意图与退出判据；Charter 0.4.0 成功边界 #6 / H-002；roadmap RT-Q02/Q03/Q05 |
| audit_type | vision-plan |
| verdict | **conditional** |
| 建议 class | editorial |
| open required | **1**（V-F099） |

## 背景

用户 2026-08-31 裁决：缓存 / 限流 / 事件总线按 **3 个独立 VP** 立项，分组标准 = **触发条件独立 × 关门能力独立**（纠正「演化轨道相似即合并」）。三 VP 均 `planned`、0 区、`vision_ref` = `schema-ui-core-admin-foundation@0.4.0`，承接 Charter 0.4.0 同进程基座（H-002 / 成功边界 #6）。VR-053 记为 **editorial**（不改 Charter 目的/边界/非目标）。

本报告为独立 Vision Review（`source: independent`）。对照 [VRev-058](VRev-058-vp026-028-ports-planned-self.md)（self · `pass` · 0 required）只作参考，**不盲从**。本意见不修改 Charter / VP / Goal status。

**独立核验材料**：

| 材料 | 核验用途 |
|------|----------|
| `docs/vision/charter.md` 0.4.0 | 成功边界 #6、非目标、H-002 |
| `docs/vision/plans/VP-026-cache-port.md` v0.1.0 | 缓存端口意图/判据/信息项 |
| `docs/vision/plans/VP-027-rate-limiter-port.md` v0.1.0 | 限流端口意图/判据/信息项 |
| `docs/vision/plans/VP-028-event-bus-port.md` v0.1.0 | 事件总线意图/判据/信息项 |
| `docs/vision/roadmap.md` v0.61.0 | 组合表第 26–28 行；并行规则 3；RT-Q02/Q03/Q05/Q06；Admin typed domain event |
| `docs/vision/revisions.md` VR-053 | editorial 立项主张 |
| `docs/vision/alignment.md` §5 / §9 | planned 0 区合法；Vision Review 台账 |
| `apps/api/internal/handler/rate_limit.go` 及全部 `newLoginRateLimiter` 构造点 | VP-027 既有使用点分母 |
| `apps/api/internal/mail/runtime.go` `cachedAdapter` | I-026-004 不强制迁移是否成立 |
| `apps/api/modules/scheduledtasks/scheduler.go` `lastRun` | 「调度去重 ≠ 限流」 |
| workspace-009 GOAL-012 D-002 | W12 单实例边界 + Redis 预登记（VP-027 相邻冻结） |

---

## 逐条审视

### 1. 3 个独立 VP 的分组是否成立

**判断：pass**（附 recommended V-F100：关门交叉依赖削弱独立性，但不推翻分组）

独立核对三组触发与可单独关门的交付物：

| VP | 触发（独立） | 可单独关门的本波交付 | 不可被另一 VP 拖累的理由 |
|----|--------------|----------------------|--------------------------|
| VP-026 | RT-Q03：多实例 **或** C 端业务域接入（评估义务）；本波无真实缓存消费者 | Cache 端口 + 内存供应商 + Redis **接缝声明** | 无消费者也可凭端口测试关门；不必等限流迁移或 EventBus 订阅者 |
| VP-027 | RT-Q05：多实例 **或** C 端限流无法共用进程内 limiter；**已有**认证面消费者 | 端口 + 滑动窗口内存供应商 + **既有使用点迁移不回归** + Redis 接缝 | 迁移关门与缓存有无消费者无关 |
| VP-028 | RT-Q02：多实例 / 跨机长任务 / 领域事件 fan-out（broker 仍 gated）；进程内解耦需求最弱 | EventBus 端口 + channel 实现 + outbox/MQ **接缝声明** | 不实现 broker；不必等 Redis 轨道 |

「缓存+限流因同 Redis 演化轨道而合并」会把 **无消费者的缓存端口** 与 **有安全回归面的限流迁移** 绑在同一关门分母上；「EventBus 当共享基建并入」会把最弱触发的运输端口绑进 Redis 轨道。用户采用的「触发独立 × 关门独立」比演化轨道合并更贴本仓库 VP-013～VP-025 单能力族模式。

独立性被削弱的一点：三 VP 退出判据 5 都要求「三端口供应商约定在（首个实现 VP 的）D-001 登记」。这是跨 VP 关门交叉依赖，见 V-F100；**不**使分组本身不成立。

---

### 2. 各 VP 意图是否落在 Charter 0.4.0 边界内

**判断：pass**（附 recommended V-F102：须把「不预制」的解释规则写进 VP，以免激活时误耗 trigger）

Charter 成功边界 #6（`charter.md` L43）分两句：

1. 基座提供业务域可直接消费的基础设施端口（已列持久化 / 认证 / Job / 可观测 / 对象存储 / 出站邮件）。
2. 「缓存、限流扩展、队列等有状态横切能力在业务域模块触发时由架构分支按需交付，**不预制**。」

H-002（L67）冻结的是：上述横切能力须在业务域触发时有**路线图登记**；复核触发与 freshness 发现机制已由 V-F093/V-F095 闭合。

三 VP 的实际分母是 **端口契约 + 进程内默认 + 外部实现接缝声明，且明确不引入 Redis/broker 客户端**。这与 #6 第二句的「不预制」可调和，条件是把「不预制」解释为 **不预制外部有状态实现（Redis / outbox / broker）**，而不是「连端口都不得在触发前存在」。用户 2026-08-31 的早期化裁决 + VR-053 editorial 构成该解释的书面授权；三 VP 均写「不改 Charter」，机读 `vision_ref` 精确匹配 `@0.4.0`。

未构成 P-006 §6.3 冲突（下级并未肯定 Charter 非目标所禁止之事，也未收缩仍生效的成功边界）。张力在于：**本波在业务域 VP 尚未激活、RT-Q0N trigger 尚未成立时交付进程内能力**。该张力必须在 VP 正文写明「本波不消耗 RT-Q02/Q03/Q05 的 trigger；Redis/outbox/broker 仍 gated」，否则后续读者会把「已有通用缓存端口」读成 Charter 允许预制 Redis。见 V-F102。

非目标「不预制 C 端 API 业务逻辑」（Charter L52）未被触碰：三 VP 都不定义订单/支付等产品语义。

---

### 3. 各 VP 方向级退出判据是否可判定、无歧义、可核验

**判断：conditional**（VP-026 / VP-028 可判定；VP-027 退出判据 3 **不可核验** → V-F099 required）

| VP | 判据 | 独立判定 |
|----|------|----------|
| VP-026 1–4、6–7 | 端口快测、双策略+可插拔、内存有界/TTL/驱逐、接缝「不引入 Redis 客户端」、未改 Charter/Profile、required=0 | **可核验**。负向证据形态明确（`go.mod` 无 Redis 客户端）。 |
| VP-026 5 | 三端口约定 D-001 | 可核验「有无登记」，但所有权与 EventBus 范畴不清（V-F100） |
| VP-027 1–2、4、6–7 | 端口快测、滑动窗口内存供应商、Redis 接缝、边界保持 | **可核验**。演进 `loginRateLimiter`（`rate_limit.go` L12–27：滑动窗口 + 容量驱逐 + `allow` 不注册 key）与判据 2 对齐。 |
| **VP-027 3** | 「登录 / captcha / recovery / MFA **全部**接入端口；D-001 P1 测试通过」 | **不可核验**。意图、首波冻结表、退出判据 3 三份清单不一致；相对代码构造点不完整。见下表。 |
| VP-028 1–3、6–7 | 端口、channel（发布/订阅/退订/并发/顺序/panic 隔离）、接缝不引入 broker、边界保持 | **可核验** |
| VP-028 4 | 「与 Admin typed domain event 扩展接缝登记对齐（端口已交付，运输仍 gated）」 | 证据形态偏弱（一行 roadmap 注记即可勾选），须防被读成解除 Admin gated（V-F101） |
| VP-028 5 | 同 VP-026/027 的 D-001 | 同 V-F100 |

**VP-027 三份清单 vs 代码构造点**（独立扫描 `newLoginRateLimiter(`）：

| 代码构造点 | 文件 | 意图（VP-027 L33） | 冻结表（L45） | 退出判据 3（L74） |
|------------|------|-------------------|---------------|-------------------|
| 登录 15min/20 | `auth.go:60` | 有 | 有 | 有 |
| 验证码生成 1min/10 | `captcha.go:36` | 有 | 有（captcha） | 有 |
| 密码修改 15min/5 | `account_self.go:51` | **有** | **无** | **无** |
| 自助恢复 15min/20 | `recovery.go:58` | 有 | 有 | 有 |
| MFA verify 独立桶 15min/10 | `mfa.go:121–124`（A6，注释写明**不与登录桶共用**） | **无**（只写 MFA step-up） | 被含糊的「MFA」覆盖？ | 被含糊的「MFA」覆盖？ |
| MFA step-up 15min/5 | `mfa.go:129–132`（enroll/disable/recovery-rotate） | 有 | 被含糊的「MFA」覆盖？ | 被含糊的「MFA」覆盖？ |
| **邀请接受** 15min/10 | `invites.go:308`（W13 F-001 预认证面 CPU DoS 刹车） | **无** | **无** | **无** |

「全部接入」在分母不完整时可以把未列点当作已覆盖或当作不存在。邀请接受与 MFA verify 均来自安全审计 required 修复，不是可静默遗漏的边角。账号分层锁定（GOAL-014，DB 行锁）**不是** `loginRateLimiter`，三 VP 未将其吞进限流——此项卫生成立，但须在分母里显式排除。

VP-028 I-028-002 把「缓冲容量与背压处理」标为 R1 required，同时非目标写「背压……消费者触发后评估」（L51）。R1 只应冻结「缓冲满时的最小语义」（阻塞 / 丢弃 / 返回错误），不得把完整背压产品拉进本波。不升级为 required。

---

### 4. 非目标边界是否卫生

**判断：pass**（限流使用点分母缺口记在第 3 点 / V-F099，不在本条重复升格）

| 边界主张 | 独立核验 |
|----------|----------|
| 缓存 ≠ 限流 ≠ 消息 | 三 VP 互指非目标；端口动词正交（Get/Set/Delete vs Allow/Record vs Publish/Subscribe） |
| 不重开 VP-012 | VP-028 声明新运输端口、不重开 correlation/审计/Job。Job 六态仍是 VP-012 已交付面。缺一句「EventBus 不是 Job 替代、持久化工作仍走 Job 端口」——建议级，见 V-F101 |
| 调度去重不是限流 | `scheduler.go` L34–39：`lastRun` / `unscheduled` 是分钟槽与失败记录去重，领域调度状态。VP-027 L55 排除成立 |
| mail `cachedAdapter` 不强制迁移 | `runtime.go` L130–132：`updatedAt` + `sender` 版本戳失效，不是通用 TTL 缓存。I-026-004 non-blocking +「不强制」成立 |
| 不替代 VP-009 / VP-010 | 三 VP 将安全/符合性 gap 归持续程序 |
| 分布式锁仍 gated | VP-026 排除 RT-Q04，成立 |
| 业务配额 / 产品事件语义 | 分别归业务域；成立 |

---

### 5. 信息需求（P-005）是否充分且分级正确

**判断：pass**（附 recommended：两处分级在 R1 选择后可能必须升级）

对照 P-005 最小列：三 VP 均有编号、问题、required/non-blocking、影响门禁、最晚阶段、状态（待裁决/待确认）。与本仓库 planned VP 惯例（如 VP-021 初创表）一致；「验证/收集动作」在 VP 层常推迟到 lead Root 台账，不把缺列升为 required。

分级独立评估：

| 项 | 级别 | 独立意见 |
|----|------|----------|
| I-026-001 API 形态 | required · R1 | 正确。零值/未命中是契约冻结问题 |
| I-026-002 TTL 清理 | required · R2 | 正确。后台协程 vs 惰性影响判据 3 |
| I-026-003 命名空间形态 | non-blocking · R1 | 可接受：退出判据 1 要求「有命名空间」，不要求在模块 ID vs 参数之间先裁。若 R1 把命名空间写进公共签名，建议升 required |
| I-026-004 mail 迁移 | non-blocking | 正确（见第 4 点） |
| I-027-001 Allow/Record 形态 | required · R1 | 正确。现状即为拆分（`allow`/`record`/`clear`/`retryAfterSeconds`） |
| I-027-002 演进 vs 双轨 | required · R2 | 正确；**与 V-F099 耦合**：双轨时未列入分母的构造点会留在旧类型上 |
| I-027-003 窗口策略 / 是否与缓存共用策略接口 | non-blocking | 默认应 **不共用**（否则破坏端口独立）。级别正确 |
| I-027-004 路由+用户复合 key | non-blocking | 本波既有点已有自己的 key 维度；可留给业务域。正确 |
| I-028-001 类型化机制 | required · R1 | 正确 |
| I-028-002 同步/异步与缓冲 | required · R1 | 正确；「背压」须收窄（第 3 点） |
| I-028-003 handler 错误语义 | required · R1 | 正确 |
| I-028-004 事件类型注册归谁 / 是否保持 Admin gated | **non-blocking** · R3 | 若 I-028-001 选**注册表**，注册权属变成方案冻结问题，须升为 required。见 V-F101 |

未伪装为已决定。关键 API / 清理 / 投递 / 迁移策略均 required。充分性在 planned 阶段成立，但 **VP-027 使用点清单本身不是信息项、却是退出判据 3 的分母**——这不是「再加一条 I-027-00N」就能代替的，必须改判据。见 V-F099。

---

### 6. 与 RT-Q02/Q03/Q05 trigger-gated 语义衔接是否有矛盾

**判断：pass**（无硬矛盾；有两处必须靠解释规则/措辞才能保持无矛盾 → V-F101 / V-F102）

| 登记项 | roadmap 状态 | VP 承接 | 独立核对 |
|--------|--------------|---------|----------|
| RT-Q03「缓存（Redis 等）」L143 | 仍 **trigger-gated** | VP-026 planned = 端口+内存+接缝；**Redis 实现仍 gated** | 行状态继续描述 Redis 项，与端口 VP 并存。禁止「先上 Redis 再找场景」未被 VP-026 违反（VP-026 明确不引入客户端，且 JWT/会话/热配置无消费者不预制） |
| RT-Q05「登录/API 限流跨实例」L145 | 仍 **trigger-gated** | VP-027 planned = 端口化**已有**进程内 limiter + 接缝 | 与 W12 D-002（单实例官方边界、Redis 预登记不实施、login/recovery 预算保持）相容，只要本波不改窗口常量、不引入 Redis 客户端 |
| RT-Q02「外部消息队列 / Job broker」L142 | 仍 **trigger-gated** | VP-028 = 「应用契约前置」 | **措辞问题**：RT-Q02 本身是运输项；「应用契约」按并行规则 3（roadmap L75）归属 **Admin 功能**。VP-028 正文又说只交运输端口。见 V-F101 |
| RT-Q06 事务 outbox L146 | 仍 **trigger-gated** | VP-028 声明接缝、**不实现** outbox、不进 outbox 表设计 | 无矛盾。接缝不得预裁 RT-Q06 表结构（VP 已排除） |

H-002 触发句「业务域 VP 激活即视为触发条件成立」**尚未发生**（无业务域 VP）。因此三 VP 是 trigger 之前的端口早期化，不是 RT-Q0N 的 Redis/broker 交付。roadmap 组合表第 26–28 行与 RT-Q 注记已原子同步，机读状态一致。

---

### 7. 跨 VP 共享「供应商约定」是否合理

**判断：pass with recommended**（约定本身合理；把 EventBus 塞进 Redis key 约定 + 用 Goal D-001 当跨 VP 权威不合理 → V-F100）

合理部分：VP-026 与 VP-027 **同 Redis 演化轨道**，key 前缀 / 命名空间 / 测试 harness（内存假供应商、契约测试）应避免两套方言。共享的是**约定**不是交付物，与「关门独立」兼容——前提是约定有**单一所有者**，且不是三条退出判据互相等待。

不合理部分：

1. VP-026 L34 把「**Redis 连接管理** / key 前缀 / 命名空间 / 测试 harness」写成三端口（含 VP-028）统一约定。VP-028 的演化是 outbox/MQ，不是 Redis key。这是范畴错配。
2. 三 VP 退出判据 5 都把「D-001 已登记」当作本 VP 关门条件。并行开区时「首个实现 VP」未定义（VRev-058 V-F096 方向正确，但只谈时序，未指出 EventBus 错配与 D-001 跨区权威问题）。
3. D-001 是**工作区 Goal 决策**。workspace-protocol 不允许把一区 Goal 决定直接当成另一区/另一 VP 的状态真相。三 VP 若分三区，应升到架构文档或指定单一 owner VP，其余「继承或本地默认」，而不是三条判据绑同一份 Goal 记录。

---

### 8. 是否有遗漏的触发条件、边界冲突或范围滑移风险

**判断：conditional**（主风险 = V-F099；其余为 recommended）

已识别、须处理：

1. **VP-027 使用点分母遗漏**（required，V-F099）：邀请接受、MFA verify 独立桶；密码修改仅在意图出现；「MFA」一词掩盖两个独立 limiter。
2. **VP-028「应用契约前置」**（recommended，V-F101）：可能把 Admin typed domain event（roadmap L333，仍全部 trigger-gated）拖进架构 VP，或被读成解除 gated。
3. **「不预制」解释未落盘**（recommended，V-F102）：激活时可能误把端口 VP 当作 RT-Q03/Q05 trigger 已消耗。
4. **相邻已 closed 合同未引用**（recommended，V-F104）：
   - VP-009 W12 D-002：单实例官方边界、Redis 仅预登记、login/recovery 15min/20/`IP|identifier`/`Retry-After` 保持。VP-027「行为不回归」应显式继承，避免端口化变成语义重写。
   - VP-021 停机合同：若 I-026-002 选后台清理协程、或 I-028-002 选异步 channel，新内核端口必须声明 SIGTERM 下取消订阅 / 排空 / 停清理，否则与已 closed 停机合同冲突。
5. **EventBus 序列化张力**（recommended，V-F103；独立确认 VRev-058 V-F098）：I-028-001「负载限定可序列化（为 outbox 预留）」与进程内 channel 可不序列化的轻量初衷冲突；R1 必须写未选方案。

未发现：多愿景、`vision_ref` 不匹配、planned 0 区不合法、重开 VP-012 记录、把调度去重或 mail 版本戳缓存强行纳入、或 RT-Q04/Q06 被本波实现。

**机读对齐（附）**：三 VP `doc_type: vision-plan`、id 与文件名一致、`status: planned`、`lead_workspace` 空、`vision_ref` = 现行 Charter `@0.4.0`。alignment §5：planned 允许 0 工作区。`workspaces.md` 无 026/027/028 行，正确。单愿景不变量成立。

---

## Verdict

**conditional**

分组标准成立，Charter 0.4.0 对齐在「端口 + 进程内默认、外部实现仍 gated」解释下成立，VP-026 / VP-028 方向级退出判据可判定，非目标与 RT-Q 衔接无硬矛盾，P-005 分级大体正确。

**不可宣称三 VP 作为一组「方向已稳」**，也 **不得激活 VP-027**，直到 V-F099 合法闭合：VP-027 退出判据 3 的既有使用点分母相对代码不完整，且意图 / 冻结表 / 判据三处不一致，无法核验「全部接入」「行为不回归」。

VP-026 / VP-028 无本报告 required 项；是否单独激活由 `/vision` 按「关门独立」决定，但不得用二者的就绪掩盖 VP-027 分母缺口。建议 class = **editorial**（改 VP-027 分母与跨 VP 措辞，不改 Charter）。

对照 VRev-058：self 将退出判据一律判 pass，并将使用点问题写成 recommended（V-F097，只谈测试基线）。独立核验后 **升级为 required**——缺口是分母本身，不是测试措辞。V-F096/V-F098 方向独立确认，改写为 V-F100/V-F103，不关闭 self 原条。

---

## Findings

### 必改（required）

| id | finding | 严重度 | 证据 | 影响门禁 | 关闭要求 |
|----|---------|--------|------|----------|----------|
| V-F099 [required] | VP-027 既有 `loginRateLimiter` 使用点分母不完整，且意图 / 首波冻结 / 退出判据 3 三份清单不一致。「全部接入端口」因此不可核验。代码 7 处构造：登录、验证码生成、密码修改、自助恢复、MFA verify 独立桶、MFA step-up、邀请接受。意图列 5 处（有密码修改与 MFA step-up，无 invite / MFA verify）；冻结表与判据 3 列 4 类且「MFA」有歧义；邀请接受（W13 F-001）与 MFA verify（A6，注释写明独立于登录桶）完全未出现。若 I-027-002 选双轨，未列点会留在旧类型；即便演进类型，回归套件仍可能不覆盖未列点。 | high（安全面迁移分母） | VP-027 L33 / L45 / L74；`auth.go:60`；`captcha.go:36`；`account_self.go:51`；`recovery.go:58`；`mfa.go:121–132`；`invites.go:282–308` | **VP-027** 激活、R1 分母冻结、退出判据 3、宣称三 VP 方向已稳。不阻断 VP-026/028 单独激活（须在 `/vision` 响应中写明 scope） | 以代码扫描冻结**完整使用点清单**并对齐意图、首波冻结、退出判据 3 三处。每个构造点：迁入 / 显式移出（移出须写明不回归责任与复审触发）。显式排除 GOAL-014 分层锁定（非限流器）。回归证据形态写明：各迁入点既有 handler 测试全量通过 + `rate_limit.go` 单元语义（allow 不注册 key、容量驱逐、Retry-After、trusted-proxy/`loginClientIP`）。闭合路径：`fixed`（改 VP-027 正文）或用户书面 `accepted-residual` / `user-overruled` 并留痕 |

### 建议（recommended）

| id | finding | 建议 |
|----|---------|------|
| V-F100 [recommended] | 三 VP 退出判据 5 都要求「三端口供应商约定」在首个实现 VP 的 D-001 登记；VP-026 L34 把 **Redis 连接管理** 纳入含 VP-028 的统一约定。VP-028 不走 Redis。结果：(1) 关门交叉依赖，削弱「关门独立」；(2) EventBus 范畴错配；(3) Goal D-001 不能当跨工作区权威。独立确认并扩展 VRev-058 V-F096。 | 指定单一 owner：建议 VP-026 只覆盖 Redis 轨道（026+027）的 key 前缀/命名空间/连接管理/harness。VP-028 改为 topic/订阅命名 + 契约测试 harness，**删除** Redis key 前缀作为其关门条件。约定升到架构短文或 owner VP 决策，其余 VP 写「继承或本地默认」。不要三条判据互等同一份 D-001。 |
| V-F101 [recommended] | VP-028 同时自称 RT-Q02「**应用契约**前置」（VP L23/L66；roadmap L142）与「只交付**运输**端口」（意图/非目标）。roadmap 并行规则 3（L75）：领域事件**应用契约** → Admin 功能；outbox/broker → 架构。Admin typed domain event 仍全部 trigger-gated（L333）。I-028-004 把「类型注册归谁 / 是否保持 gated」标 non-blocking。存在把应用契约拖进架构 VP、或误解除 Admin gated 的滑移风险。亦未写明 EventBus ≠ Job 端口。 | 统一措辞为「运输端口 + 进程内实现 + outbox/MQ 接缝声明」，删除或限定「应用契约前置」。写明**不**解除 Admin typed domain event 的 trigger-gated。若 I-028-001 选择注册表，将 I-028-004 升为 R1 required。补一句：持久化/重试工作仍走 VP-012 Job，本 VP 不替代。 |
| V-F102 [recommended] | Charter #6「不预制」+ RT-Q02/Q03/Q05 仍 trigger-gated，与无业务域触发下交付端口+进程内默认的早期化之间，解释规则未写入 VP。当前「不引入 Redis/broker 客户端」可核验，但激活时可能被读成已消耗 RT-Q0N trigger。 | 各 VP 增加一句解释规则：本波 = 基座可消费面早期化（端口+进程内默认+接缝）；**不消耗** RT-Q02/Q03/Q05 trigger；Redis/outbox/broker **实现**仍须等触发后评估。不改 Charter。 |
| V-F103 [recommended] | I-028-001 将「负载是否限定可序列化（为 outbox 预留）」与类型化机制绑在同一 required 项。进程内 channel 可传非序列化负载；为 RT-Q06 预留则可序列化约束。该取舍未在 VP 内作显式未选方案。独立确认 VRev-058 V-F098。另：I-028-002「背压处理」与非目标「背压不进本波」需在 R1 收窄为缓冲满时的最小语义。 | R1 契约冻结把「进程内允许非序列化 vs 预留 outbox 的可序列化约束」写成显式取舍（含未选方案）。I-028-002 的背压范围限定为满缓冲阻塞/丢弃/返回错误，完整背压产品仍 gated。 |
| V-F104 [recommended] | 相邻已 closed 合同未登记：(a) VP-009 W12 D-002——单实例官方边界、Redis 仅预登记、login/recovery 预算/键空间/`Retry-After` 保持现状；VP-027 未引用，端口化可能重写窗口常量。(b) VP-021 停机合同——后台 TTL 清理或异步 EventBus handler 会引入 goroutine/channel，须声明 SIGTERM 语义。 | VP-027 相邻 VP 表增加 W12 D-002：本波不改官方单实例边界、不实施 Redis、不改 login/recovery 窗口常量（除非用户另裁）。VP-026/028：若 R1/R2 选择后台协程或异步投递，须声明停机时停止清理/取消订阅/排空；否则选惰性清理与同步投递，避开新生命周期。 |

---

## 声明

本意见不修改 Charter / VP / Goal status，不写 Goal `03-audit`。required finding 的响应由 `/vision` 追加在本报告中；原 verdict 与 finding 原文不得改写。实施工作交 `/govern`。

建议 `/vision` 响应输入：

1. **先闭合 V-F099**（改 VP-027 使用点分母，或用户书面 residual/overruled）后，才可宣称方向已稳或激活 VP-027。
2. V-F100～V-F104 可在激活事务内 editorial 处理；其中 V-F102 的「不消耗 trigger」一句建议三 VP 同时补，避免激活叙事分叉。
3. 不改 Charter；不把本报告当作 Redis/outbox/broker 已立项。

---

## `/vision` 响应（2026-08-31 · 编排器 · append-only）

| finding | 响应 | 状态 |
|---------|------|------|
| V-F099 [required] | **fixed**：VP-027 意图 / 首波冻结表 / 退出判据 3 三处对齐为代码扫描的 **7 处构造点**（登录 `auth.go:60` / 验证码 `captcha.go:36` / 密码修改 `account_self.go:51` / 恢复 `recovery.go:58` / MFA verify 独立桶 `mfa.go:121` / MFA step-up `mfa.go:129` / 邀请接受 `invites.go:308`），每处标注参数；显式排除 GOAL-014 分层锁定（DB 行锁，非限流器）；回归证据形态写明（既有 handler 测试全量 + `rate_limit.go` 单元语义 + W12 D-002 窗口常量保持）。VP-027 v0.1.1。原 verdict（conditional）与原文未改写。 |
| V-F100 [recommended] | **fixed**：Redis 轨道约定收窄为 VP-026/027（key 前缀 / 命名空间 / 连接管理 / 测试 harness），登记于架构短文或 owner VP 决策（单一所有者，不跨区绑同一份 Goal D-001）；VP-028 共享约定改为 topic/订阅命名 + 契约测试 harness，删除 Redis key 前缀作为其关门条件。三 VP 退出判据 5 同步修订。 |
| V-F101 [recommended] | **fixed**：VP-028 定位统一为「运输端口 + 进程内实现 + outbox/MQ 接缝声明」，删除「应用契约前置」含糊措辞；明确不解除 Admin typed domain event 的 trigger-gated（应用契约归 Admin 功能分支，并行规则 3）；补 EventBus ≠ Job 端口声明；I-028-004 注明若选注册表则升 required。roadmap RT-Q02 注记与组合表 VP-028 行同步。 |
| V-F102 [recommended] | **fixed**：三 VP 均补解释规则句「本波 = 基座可消费面早期化（端口 + 进程内默认 + 接缝声明），**不消耗** RT-Q02/Q03/Q05 trigger；Redis/outbox/broker **实现**仍须等触发后评估」。 |
| V-F103 [recommended] | **fixed**：VP-028 I-028-001 改为 R1 显式取舍（进程内非序列化 vs outbox 可序列化预留，记录未选方案）；I-028-002 背压收窄为「缓冲满时最小语义 = 阻塞 / 丢弃 / 返回错误，完整背压仍 gated」。 |
| V-F104 [recommended] | **fixed**：VP-027 相邻边界表补 W12 D-002（单实例边界 / 窗口常量保持）；VP-026/028 补停机语义声明（后台清理协程 / 异步投递须声明 SIGTERM 排空，否则选惰性清理 / 同步投递避开新生命周期）。 |

**闭合结果**：V-F099 required → **fixed**；V-F100～V-F104 recommended → **fixed**。开放 required = **0**。VRev-059 原 verdict（conditional）保留为历史结论；二轮复审后可重新审视方向已稳状态。
