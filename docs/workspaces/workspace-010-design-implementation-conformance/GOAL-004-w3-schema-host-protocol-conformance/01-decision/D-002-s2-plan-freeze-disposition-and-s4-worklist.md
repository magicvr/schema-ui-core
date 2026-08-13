---
id: D-002
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: S2 方案冻结 — 上游权威处置映射与 S4 工作清单
status: accepted
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# D-002 · S2 方案冻结 — 上游权威处置映射与 S4 工作清单

## 背景

上游 ADR-0034～0037 已于 2026-08-13 accepted，v2.8.0 正式发布（tag `521cff8`）并固定入仓
（E-004/E-005）。S2 出口的实质输入已齐备；本决定把上游逐项处置（附件 §1c，权威 = ADR-0034
D10/D6）映射为本仓 S2 冻结结论与 S4 工作清单，冻结 S2 方案范围。

## 决定

### 1. 处置权威与本地映射（I-001/I-002 闭环方式）

1. **逐项处置以附件 §1c 为本地投影、以上游 ADR-0034 D10/D6 为裁定权威**（A-002 F-4 注记
   维持）：95/95 行均已有 adopt-now / reserve-extension / explicitly-out 处置与理由 → 附件
   §6 门禁 (a) 满足。
2. **adopt-now 行的 shape/state/security/fixtures 齐备性**：由已 accepted 的 ADR-0035
   （bootstrap 生命周期与归一化 auth state）、ADR-0036（Host failure result、return intent、
   A11y）、ADR-0037（claim/registry/evidence）及 v2.8.0 正式机器制品（三个 host schema +
   capability-registry + 99 host fixtures + app-manifest 41）覆盖；本仓生产模块逐字段消费
   （E-004，零排除）→ 附件 §6 门禁 (b) 满足。
3. **reserve-extension / explicitly-out 行**：一律为「上游 deferred / 协议不负责」记录，不创建
   空 capability、不冒充已保留扩展点（§1b 规则维持）→ 无冒充。
4. **IMP-001～004 裁定（ADR-0034 D6）与本地动作**：
   - IMP-001（Settings PATCH 未知字段）`explicitly-out`：消费者业务 API 非协议 request
     object → **S4 无动作**；本仓 settings.go 已是业务 API 语义（D6 权威）。
   - IMP-002（Users 导航 label 双源）`adopt-now`：消费侧单源投影（manifest `labelKey`
     命中优先、`label` 字面 fallback）已由 `apps/web/src/app/navigation.ts` 实现；服务端
     provider `Label` 仅落 DB `menu_items` 元数据，不经 UI 渲染消费 → S4 动作 = **核验 +
     证据固化**（served manifest 导航来源断言测试），不新增 provider 消费权威。
   - IMP-003（Host 侧结构/语义验证与 fixtures 消费）`adopt-now`：已由 E-004 生产入口
     （host 三模块 + manifest 验证 + 99/41 fixtures 零排除）满足 → S4 无新动作，S5 复核。
   - IMP-004（row selection → drawer/detail）`reserve-extension`：上游保留独立 overlay
     ADR（D6），本波不实现、不冒充 → 无动作。
   → 附件 §6 门禁 (c) 满足。

### 2. S4 工作清单（本决定冻结；未列项目不动）

| # | 项目 | 来源 | 动作 |
|---|------|------|------|
| S4-1 | **return intent 登录流接入** | E-004 residual #3；ADR-0036 D6 | 捕获（boot reauth-required 终态 / 会话中期 auth-lost；协议 allowlist 收窄；sessionStorage + nonce + 过期 + `/login` 自循环拒绝）→ 登录成功后校验并恢复 path?query |
| S4-2 | **session adapter reauth-required 映射** | E-004 residual #5；ADR-0035 D4 | refresh token 存在但轮换失败（revoked/expired）→ `reauth-required`（终态）；logout 清除后 → `anonymous`；`locked` 无服务端状态源，映射层保留、生产源记为 residual（见 §3） |
| S4-3 | **hostOwnedPaths 显式集合** | E-004 residual #6 | 声明 `["/login"]`：认证态下命中 host-owned path 不再产生 `HOST_ROUTE_NOT_FOUND`，导航回 home |
| S4-4 | **multi-round $deps residual 纠错** | E-004 residual #2（**过时**） | 引擎已于 `e18edce`（2026-08-08）落地且 stage3 reactions 套件零排除——**纠错台账**并更新 claim 生成脚本的 residuals 文本，重生成 claim |
| S4-5 | **304/ETag** | E-004 residual #4 | 维持 200-only 装载（ADR-0035 D6 conditional GET 为「可用于」可选优化；200-only 是合规路径）→ residual 关闭为**无动作**，注明依据 |
| S4-6 | **account-locked 生产源** | E-004 residual #5 后半 | 无服务端锁定状态可消费；映射 + fixtures 已 pin 行为 → 拟议 residual（见 §3，S6 用户 P-004 决策） |
| S4-7 | **IMP-002 导航 label 单源固化** | ADR-0034 D6；F-4 | 消费侧 labelKey 命中优先/字面 fallback 已实现（`navigation.ts`）；served manifest 导航来源仅 fragment（provider `NavigationContribution.Label` 只落 DB `menu_items` RBAC 元数据，不经 UI 渲染）→ 固化证据断言测试（served manifest 导航 labelKey/label 一致 + DB label 不参与 UI 投影），不新增 provider 消费权威 |

### 3. 已登记 residual（S2 冻结时点；S6 对照，最终需用户 P-004 书面决策）

1. **account-locked 生产源缺位**（S4-6）：**拟议 `accepted-residual`**（范围：本波不新增账号锁定安全特性；复审触发：认证迭代引入锁定状态时）。—— 该 residual 在 S6 关门时点须经用户 P-004 书面接受或驳回，本决定不预写接受。
2. **304/ETag conditional GET 复用未实现**：合规路径（200-only）已过 fixtures；属可选优化（ADR-0035 D6「可用于」），不作为本波工作。
3. **return intent 的 registered 扩展 allowlist**（manifest `returnIntentQueryKeys`）：boot auth 终态发生在 manifest-load 之前，捕获时仅可用协议 allowlist（收窄合法）；登录后恢复不扩张 allowlist。会话中期 auth-lost 同此。
4. **页面协议 2.7 mandatory residual（R5）**：multi-round `$deps` 已实现（S4-4 纠错），该 residual 随纠错关闭；claim `pageVersions: ["2.7"]` 维持（2.8 未改页面 schema 字段集）。

### 4. 门禁联动

- 附件 §6 四复选框：本决定落盘后 (a)(b)(c) 勾选；(d) 在 cross 方案审视（self + grok
  independent）落盘且 required findings 闭合后勾选。
- S2 冻结后 S4 仅执行 §2 清单；实施中发现协议仍不足 → 回流 S2/I-002（D-001 §影响）。

## 影响

- I-001：初始基线（§1c 95 行）已为逐项对照结果；`collecting` → 以本决定 §1 + 附件 §1c 为
  冻结证据闭环（随 E-006 更新为 verified）。
- I-002：处置与理由已在上游 accepted 权威齐备；本决定为本地冻结记录。
- I-005：本仓消费证据（迁移矩阵零动作项、registry 弃用机制、正反 fixtures、版本协商）
  随 E-006 录入 → verified。
- I-006：上游裁定维持（D6 IMP-004 独立 overlay ADR）。
