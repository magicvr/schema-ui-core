---
id: GOAL-005-r4-evidence-closeout
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## D-001 · R4 实施方案（lead 方案冻结）

### 触发

R1～R3 已关门（GOAL-002/003/004 done）；Root `progress 3/4`。R4 = 证据与关门：把合同六节收成可核对证据并闭合 Root 成功标准。

### 决定（方案）

1. **证据矩阵（C1）**：落盘 `attachments/r4-evidence-matrix.md`——按合同 §1～§6 逐条映射证据（服务端/前端工具/快测文件+用例数+commit），并给 `zh-CN`/`en-US` 双 locale 范例表（时区展示、货币展示、输入解析、round-trip）。
2. **无越界核账（C2）**：逐项核对 §5 越界清单 + Root 边界（汇率/计费、RT-T03、Profile 默认集、模块矩阵、Manifest、`docs/contracts/` stage 门禁、热加载、VP-007/VP-012 重开、钱包演示面重开）→ `attachments/r4-evidence-matrix.md` 边界节；API 机器合同不变量（§3.3：RFC3339 UTC / int64 / JSON number）复证。
3. **核账项处置（C3）**：
   - GOAL-002：F-001（设置面字段解释留痕）→ 本书已留痕（R1 D-001 §4.1 无独立数字字段）+ F-002（映射表）→ 已履约（GOAL-004 C2）→ 双双 **closed**。
   - GOAL-003：F-001（epoch 输入控件按 §2.3）→ R4 核对 renderer 无新增含时间控件 → **closed（未触发条款）**；F-002（TIMEZONE_OPTIONS 扩展留痕）→ 常用集保持可核对（switcher 常量）→ **closed（留痕于本账）**。
   - GOAL-004：F-002/F-005（grouping 位序）→ **final residual**（延续用户 2026-08-26 接受范围；文档已留痕，R4 不再加严，理由：Admin 输入容差可接受、Intl 解析器一致性优先）；F-006（币种目录）→ **final residual**（句法近似 + §4.3 映射表为权威；全 ISO 4217 目录属后续工作区可选项）；F-007（安全整数）→ **评估：low 成本加严**——`parseLocalizedMoney`/`formatMoney` 对超出 `Number.MAX_SAFE_INTEGER` 的中间值拒绝（parse → null；format → ""），补快测；若实施顺利即 **fixed**。
4. **关门路径（C4）**：Root `GOAL-001` 03-audit 自审 A-001（self · 成功标准 1～4 逐条）→ 本地 grok build（grok-4.6 · high）independent（Root scope）→ 意见合并响应 → 用户书面确认 → Root done 4/4 → goal-tree 收官、workspace.md 结项记录、VP-020 关门记录填写（vision 层）。

### 为什么

- 证据矩阵复用 R2/R3 快测与双全量回归，避免重复实现；双 locale 范例与快测同源（合同 §6 核对方式）。
- 越界核账逐条对应 §5 清单，防「名义声明无证据」。
- F-007 安全整数为低风险高确定性加严（一行守卫 + 快测），消除 JS number 与 int64 主张的缺口；其余项维持已被用户接受的 residual 范围，不扩大 R4 范围。
- Root 关门走项目级 independent 路径（R3 已证明能抓真实缺陷：bodyMapping/site 通道）。

### 未选方案

- F-005 grouping 位序严格校验（正则位序检查）：复杂度/收益比低，Admin 输入场景可容差；维持 residual。
- F-006 全 ISO 4217 目录（API 侧枚举 + 前端目录）：范围扩张超出本波；若后续工作区需要可立项。
- 每个页面手写范例 UI：R4 只要求「快测 + 范例」可核对，范例表 + 快测同源即可，不造演示页。