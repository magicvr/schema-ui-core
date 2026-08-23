---
id: GOAL-033-w22-residual-closeout
title: W22 · accepted-residual 残余全库清点收口（A 组修复 ×6 + B 组触发复核 ×6 + 台账卫生 ×3）
status: done
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-001-design-implementation-conformance
version: 1.0.0
progress: 18/18
---

# GOAL-033-w22-residual-closeout · W22 · accepted-residual 残余全库清点收口

## 概述

2026-08-23 用户指令：扫描全部治理文档中标记为 `accepted-residual` 的项目并评估可一次性处理项。全库扫描命中约 420 行（17 工作区 + vision 层），剔除规则引用/闭合路径选项/否定句后，真实存续 residual 共 23 项，分三组：

- **A 组（一次性收口，转 fixed）**：6 项代码/证据级修复；
- **B 组（复审触发已到期/已发生）**：6 项按触发条件复核并落盘结论；
- **C 组（设计裁决类，保留不动）**：11 项不在本目标范围（D3 匿名可读、manifest-only 权威序、双探针防御深度等，均为有效设计裁决或有未到期的触发）。

用户本轮书面裁决（ask_user_question）：**执行 A 组全部 + B 组全部 + 台账卫生小项**，并指定在本工作区新设子目标承载治理上下文。

## 范围清单

### A 组 · 一次性收口（residual → fixed）

| # | 来源（Q2 引用） | 内容 | 动作 |
|---|---|---|---|
| A1 | [W7·GOAL-007 E-002](../../../workspace-007-localization-and-system-settings/GOAL-007-s6-settings-form-page/02-execution/E-002-respond-a002-findings.md) F-002 | e2e admin M3 浏览器证据被端口排除区间阻塞 | 换端口/宿主补跑一次并附日志 → 转 fixed |
| A2 | [W10·GOAL-025 A-002](../GOAL-025-w16-rectification-batch-a/03-audit/A-002-self-response.md) F-003 | 老库种子 admin 未置 `must_change_password=1` | 补迁移/回填 + 测试 → 转 fixed |
| A3 | [W11·GOAL-013 E-006](../../../workspace-011-admin-functional-modules/GOAL-013-nav-order-config/02-execution/E-006-s5-closeout.md) F-004 | 导航重复/大小写边界无测（行为符合预期） | 补两条边界单测 → 转 fixed |
| A4 | [W6·Root A-012](../../../workspace-006-design-system-and-ui-experience/GOAL-001-design-system-and-ui-experience/03-audit/A-012-response-a011-and-closeout.md) F-VUI-011 | 登录密码可见切换未做 | 实现 toggle + 测试 → 转 fixed |
| A5 | [W9·GOAL-002 03-audit](../../../workspace-009-production-hardening/GOAL-002-audit-findings-remediation/03-audit.md) N-002 | 上传扩展拒绝 best-effort（无 `<script`/`<svg` 标记形态可入库） | 加内容启发式拒绝 + SVG/多格式测试 → 转 fixed |
| A6 | [W9·GOAL-011 A-004](../../../workspace-009-production-hardening/GOAL-011-w11-api-web-security-audit/03-audit/A-004-w11-closure-response.md) R-001 | `/api/auth/mfa/verify` 无独立 HTTP 限流 | verify 层加限流中间件 + 测试 → 转 fixed |

### B 组 · 触发到期复核（结论落盘：续期 residual / 升级整改 / 确认已兑现）

| # | 来源 | 到期原因 |
|---|---|---|
| B1 | [W3·GOAL-006 C1-I003](../../../workspace-003-modular-admin-architecture/GOAL-006-r4-c1-freeze-decision/01-decision.md) R4-I004 | owner `magicvr` review date **2026-08-05 已过期** |
| B2 | [W3·Root A-011](../../../workspace-003-modular-admin-architecture/GOAL-001-modular-admin-architecture/03-audit/A-011-a010-cohesion-response.md) F-003b | 触发="VP 退出 #4 取证前"，VP-003 已于 2026-08-06 关门取证 |
| B3 | [W3·GOAL-010 A-003](../../../workspace-003-modular-admin-architecture/GOAL-010-r4-c4-schema-other-migration/03-audit/A-003-r4-c4-schema-migration-response.md) C4-004 | allowlist 深化挂 R5/R6，两者均已完成 |
| B4 | [W3·GOAL-011 E-002](../../../workspace-003-modular-admin-architecture/GOAL-011-r4-c5-acceptance/02-execution/E-002-r4-c5-verification.md) C5-002 | "完整矩阵补入 R5 数据门禁"，R5 已 done，需核实 |
| B5 | [W8·GOAL-007 D-001](../../../workspace-008-admin-module-readiness/GOAL-007-s5-admission-audit-and-verdict/01-decision/D-001-s5-go-decision.md) F-007 | go 裁决时书面确认的 deferred residual，按 freshness trigger 复核 |
| B6 | [W13·GOAL-006-r5 A-001](../../../workspace-013-store-dialects/GOAL-006-r5-dual-path-acceptance/03-audit/A-001-independent-r5-execution-closeout.md) F-002 | Root I-001/I-004 的 residual 回写状态待核实 |

### 台账卫生（零风险文档项）

| # | 来源 | 动作 |
|---|---|---|
| H1 | W11·GOAL-017 F-004（个人中心 MFA UI residual 提案） | GOAL-018-mfa-manager-ui 已交付关门 5/5，补一条「已由后续目标兑现」确认性回写 |
| H2 | W17·[GOAL-001 A-002](../../../workspace-017-outbound-mail/GOAL-001-outbound-mail/03-audit/A-002-independent-closeout.md) 对 N-001 的纠偏 | self 将分母外 note 误标 accepted-residual，加更正注记（不改历史原文） |
| H3 | W10·[GOAL-002-w1 A-006](../GOAL-002-w1-examples-optional-module/03-audit/A-006-w1-closeout-response.md) F-006 | web 夹具「范例启用」改名 dogfood 语义 |

## 路线图与检查点（进度分母 = 18）

```text
S1 立项与治理容器          → C1 五件套 + goal-tree 同步
S2 A 组代码收口            → C2=A1 · C3=A2 · C4=A3 · C5=A4 · C6=A5 · C7=A6
S3 B 组触发复核            → C8=B1 · C9=B2 · C10=B3 · C11=B4 · C12=B5 · C13=B6
S4 台账卫生                → C14=H1 · C15=H2 · C16=H3
S5 回归与关门              → C17 全量回归（go test + vitest + tsc + build）
                            → C18 关门审计（self + independent 安全面）+ 台账/goal-tree 终态同步
```

## 信息需求登记（P-005）

| 编号 | 问题 | 级别 | 影响门禁 | 最晚需要阶段 | 收集动作 | 状态 | 备注 |
|------|------|------|----------|--------------|----------|------|------|
| I-001 | 本机 8080 端口排除区间是否仍在？能否换端口绑定跑通 e2e？ | required | A1 | C2 执行前 | `netsh interface ipv4 show excludedportrange protocol=tcp` + 试跑 | **verified** | 2026-08-23 复核：8011–8110 区间已不在排除表，8080 无监听占用，可绑定 |
| I-002 | B1–B4 各触发条件是否确已发生（review date 过期 / VP-003 exit#4 取证完成 / R5、R6 完成）？ | required | B1–B4 复核 | C8 前 | 本轮扫描 + goal-tree 交叉核对 | verified | goal-tree：VP-003 closed 2026-08-06；W3 区 R5/R6 均 done；review date 2026-08-05 < 今日 2026-08-23 |
| I-003 | B5 F-007 的 freshness trigger 具体内容与当前状态？ | required | B5 | C12 前 | 读 W8 D-001 全文 | **closed（E-003）** | 触发为事件型「后续业务 VP 激活前」，未发生，F-007 维持有效 |
| I-004 | B6 Root I-001/I-004 是否已按 residual 回写？ | required | B6 | C13 前 | 读 W13 Root 信息表 | **closed（E-003/E-004）** | 双账已回写、D-002 追认节落盘 |

## 边界

- 不触碰 C 组设计裁决类 residual（含 W9·D3 匿名可读、W11·manifest-only 权威序、W10·双探针等）。
- 不修改历史审计意见原文；所有响应以响应节/决策追加方式留痕。
- 跨区回写仅限上表点名的目标台账；引用遵守 workspace-protocol §2.6（文档 Q2 路径）。

## 审计模式声明

实施批次默认 `self`；因 A5/A6 涉及安全面（上传内容过滤、认证端点限流），关门审计按项目默认走 `independent`（本地 grok build · grok-4.6 · high，`/audit`），意见落盘本目标 `03-audit/` 后由编排器合并响应。
