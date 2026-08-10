---
id: GOAL-001-localization-and-system-settings
doc: audit-entry
record_id: A-001
source: independent
auditor: grok-4.5
auditor_host: grok-build-cli
thinking: high
audit_type: design-plan
scope: S0 契约冻结 · D-002（I-L10N-001～005）+ F-V029 + E-002 + 00-meta/01-decision/02-execution 索引
status: recorded
verdict: conditional
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# A-001 · S0 契约冻结独立交叉审计（design-plan）

| 字段 | 值 |
|------|-----|
| **source** | `independent` |
| **auditor** | grok-4.5（grok build CLI，`--effort high`） |
| **日期** | 2026-08-09 |
| **类型** | design-plan（方案冻结复审） |
| **scope** | S0：D-002 关闭 `I-L10N-001`～`005` + 契约冻结；`attachments/F-V029-coverage-table-s0-freeze.md`；E-002 盘点事实；`00-meta` / `01-decision` / `02-execution` 索引同步 |
| **verdict** | **conditional** |
| **声明** | 本意见**不修改**任何 `status` / `progress` / `goal-tree` 状态列 / 方案正文；响应、finding 闭合与阶段放行归 `/govern` 与用户裁决。 |

## 1. 范围与材料

### 1.1 工作区边界（已校验）

| 字段 | 值 | 核对 |
|------|-----|------|
| workspace_id | `workspace-007-localization-and-system-settings` | `workspace.md` |
| root_goal | `GOAL-001-localization-and-system-settings` | 与目标 `id` 一致 |
| canonical_scope | `docs/workspaces/workspace-007-localization-and-system-settings/` | 本审计仅读该范围 + 愿景链（VP/Charter 对齐） |
| shared_materials_catalog | `none` | 无共享资料引用；未把跨区资料当关闭证据 |
| primary_plan | `VP-007-localization-and-system-settings` | 与 Root `plan_refs` 一致 |

**未**读取其他 workspace 的过程台账（仅对本区事实与 VP-007 / 仓库内 schema/代码做只读抽查）。

### 1.2 已读材料

| 材料 | 路径 |
|------|------|
| 工作区 | `docs/workspaces/workspace-007-localization-and-system-settings/workspace.md` |
| 目标树 | `…/goal-tree.md` |
| Root meta | `…/GOAL-001-…/00-meta.md` |
| 决策索引 + D-001/D-002 | `…/01-decision.md`、`01-decision/D-001-*.md`、`01-decision/D-002-*.md` |
| 执行索引 + E-002 | `…/02-execution.md`、`02-execution/E-002-*.md` |
| F-V029 冻结表 | `…/attachments/F-V029-coverage-table-s0-freeze.md` |
| 审计索引（本轮前） | `…/03-audit.md`（空台账） |
| VP-007 | `docs/vision/plans/VP-007-localization-and-system-settings.md` |
| 协议/schema 抽查 | `docs/schemas/app-manifest.schema.json`、`node.schema.json`、`component-registry.json` |
| 代码抽查 | `apps/api/internal/kernel/profile.go`、`handler/settings.go`、`apps/web/src/theme/theme.ts`、`protocol/app-manifest.ts` |

### 1.3 本轮关注点对照

| # | 关注点 | 本意见落点 |
|---|--------|------------|
| (1) | I-L10N-001～005 关闭是否合规（P-004/P-005） | §3 |
| (2) | F-V029 分母可枚举与覆盖完备性 | §4 |
| (3) | 冻结契约 vs `schema-ui-docs@v2.7.0` / VP-007 | §5 |
| (4) | 审计模式 `independent` 是否满足 | §6 |

## 2. 成果（有证据）

| # | 主张 | 证据 | 独立复核 |
|---|------|------|----------|
| 1 | S0 决策与执行 ledger 已落盘 | `D-002`、`E-002`、`F-V029-*.md`；`01-decision.md`/`02-execution.md` 索引有行 | **成立** |
| 2 | 五条 I-L10N 均有用户选定方案 + 冻结语义写入 D-002 | D-002 §I-L10N-001～005；`01-decision.md` 信息表状态 `verified` | **成立**（见 §3 细节与 F-003） |
| 3 | F-V029 固定 UI 7 面 + pageId 并集 12 + M1～M4 + 缺失翻译纪律 | `attachments/F-V029-coverage-table-s0-freeze.md` | **成立**（见 §4） |
| 4 | pageId 并集与 Profile 模块边界可对代码 | `kernel/profile.go` mvp=8 模块、admin=+settings/activity；core pages 8 + users/roles/settings/activity = 12 | **成立** |
| 5 | `titleKey`/`labelKey` 属既有 app-manifest 兼容面 | `docs/schemas/app-manifest.schema.json` PageEntry/NavLink/NavGroup；前端 parser 已识别字段 | **成立** |
| 6 | `node.props` 开放对象（禁 CSS 名）→ 本地 registry 的 `*Key` 不改写上游 page/node 语义 | `node.schema.json` props；D-002 边界声明 | **成立** |
| 7 | theme localStorage 单通道可作 I-L10N-002 参照 | `apps/web/src/theme/theme.ts` | **成立** |
| 8 | 公开 branding 现状可扩展 | `GET /api/branding` 现仅 `siteTitle`/`logoUrl`；D-002 定为 additive 扩展 | **成立** |
| 9 | 错误 envelope 现状为 `{error,message}`；稳定码可机读 | 多处 `writeError`；本轮枚举约 **32** 个不同码字符串 | **基本成立**（见 F-005） |
| 10 | 审计模式声明为 S0=`independent` | D-002「审计模式」；本 A-001 即为该节点意见 | **模式满足**（响应仍待 `/govern`） |

## 3. 关注点 (1) · I-L10N-001～005 关闭合规（P-004 / P-005）

| ID | 最晚阶段（台账） | 状态主张 | 用户书面留痕 | 证据/冻结 | 合规判定 |
|----|------------------|----------|--------------|-----------|----------|
| I-L10N-001 | 多语种方案冻结前（S0） | verified | D-002：前端 key 解析 | 兼容盘点 + 回退链 + 边界 | **合规** |
| I-L10N-002 | 用户控制实施前（S1） | verified（提前关闭） | D-002：localStorage 单通道 + 优先级 | theme 参照 + 匿名/登录同通道 | **合规**（提前关闭合法） |
| I-L10N-003 | Settings API 方案冻结前（S3） | verified（提前关闭） | D-002：兼容扩展 `/api/branding` | 字段表 + 缓存/事件边界 | **合规** |
| I-L10N-004 | 后端提示本地化实施前（S4） | verified（路径 a） | D-002：路径 (a) 有界服务端协商 | envelope 扩展 + 编目范围 + **实施证据明示 S4** | **信息门禁可接受**；不得当作 exit 5 已交付（F-003） |
| I-L10N-005 | Localization 设置实施前（S3） | verified | D-002：UTC 存储 + 显示转换 | `siteTimezone` / 无效时区语义 | **合规** |

**P-004/P-005 核对摘要**

- 五条均为 `required`；关闭路径均为**用户选定方案 + D-002 书面冻结**，不是 silent residual。
- 无 `deferred` 逾期；无「仅口头」关闭。
- I-L10N-002/003/005 在最晚阶段前**提前** verified，允许，且有冻结语义可实施。
- I-L10N-004：信息问题是「选路径 a 还是 residual b + 兼容边界」；D-002 选定 (a) 并冻结 envelope，**满足信息门禁关闭**。VP-007 exit 5 仍要求 S4 实施证据——状态字 `verified` 容易被编排器/读者误当成 exit 5 关闭（→ F-003 recommended）。
- 未发现把未知伪装成已完成实施事实的写法；E-002 盘点与 D-002 边界大体一致。

## 4. 关注点 (2) · F-V029 分母

| 分母块 | 冻结内容 | 可枚举？ | 与 VP-007 F-V029 要求 | 结论 |
|--------|----------|----------|------------------------|------|
| 固定 UI | U1～U7（登录、Shell header、导航、语种切换、通用反馈、通用组件、Admin Settings 四类） | **是**（7 行） | 覆盖 VP 列出的固定 UI 面 | **满足** |
| Manifest/Schema | 12 pageId + schemaUrl；mvp/admin 可达列；settings/activity mvp=`N/A` 写模块边界 | **是**（12 行） | 要求双 Profile Runtime Manifest pageId/schema 并集，禁止「代表性」 | **满足**；与 `profile.go` + core 8 pages 一致 |
| 主流程 | M1～M4 | **是**（4 行） | 与 VP 文案一致 | **满足** |
| 缺失翻译纪律 | 事件名 + 回退链 + residual 排除规则 | **是** | 与 VP 纪律一致 | **满足** |
| 矩阵语义 | 行语义声明 `zh-CN`/`en-US` × `mvp`/`admin` × 匿名/已认证；证据路径 S2/S5 填 | 分母行可枚举；单元格证据**尚未填**（S0 预期） | exit 6 复用同一分母 | **S0 冻结可接受**；不得用空证据路径宣称已覆盖 |

**结论**：F-V029 作为 S0 **分母冻结**充分可枚举，覆盖固定 UI + 双 Profile pageId/schema 并集 + M1～M4 + 缺失 key 纪律。无「代表性收缩」伪装。

## 5. 关注点 (3) · 兼容面与 VP-007

### 5.1 `schema-ui-docs@v2.7.0` / 本地兼容面

| 冻结点 | 与兼容面关系 | 结论 |
|--------|--------------|------|
| 前端解析 `titleKey`/`labelKey`（及 registry `*Key`） | app-manifest 已声明；不改 page/node 协议语义；字面 `title`/`label`/`content` 保留 | **一致** |
| 缺失 key 可观察 + 安全回退 | 运行时行为，不扩展上游 schema 必填面 | **一致** |
| 本地 registry 未来 additive `placeholderKey` 等 | D-002 限本地 registry + props 开放；声明不冒充上游 | **可接受**（须在 S2 决策再登记） |
| 错误 envelope additive `messageKey`/`params` | 兼容扩展，保留 `message` 英文 | **方向一致**（实施在 S4） |
| `/api/branding` additive 字段 | 旧客户端只读 siteTitle/logoUrl 不受影响 | **方向一致** |

未发现「用私有协议字段冒充 v2.7.0 上游语义」的冻结写法。

### 5.2 与 VP-007 冲突 / 收缩

| VP-007 承诺 | D-002 / F-V029 冻结 | 判定 |
|-------------|---------------------|------|
| 语种 `zh-CN`/`en-US`/`auto` 与优先级 | 一致 | **对齐** |
| 四类设置 General/Branding/Localization/Appearance | 四类名一致 | **名称对齐** |
| Branding：**Logo、favicon、浅色/深色 Logo** | 仅冻结 `logoUrl`；深/浅色「首版仅保留单 logoUrl」；**favicon 未列入且无 residual** | **范围收缩 · 必改**（F-001） |
| exit 5 路径 a/b | 选定路径 a | **对齐**（实施待 S4） |
| F-V029 分母结构 | 已落盘 | **对齐** |
| Settings 仅 admin Profile | F-V029 N/A + D-002 | **对齐** |
| Non-goals：不重定义 v2.7.0 | D-002 边界 | **对齐** |

成功边界 3（`00-meta` 镜像 VP exit 3）仍写「品牌资产」等产品面可用；若 S0 冻结字段永久缩为单 `logoUrl` 且无 residual，则与 VP 方向级范围冲突。

## 6. 关注点 (4) · 审计模式 independent

| 检查 | 结果 |
|------|------|
| D-002 是否声明 S0 冻结 = independent | **是** |
| 本轮是否由独立会话 / 独立角色写入 `source: independent` | **是**（本 A-001） |
| 是否越权改 status/progress/方案 | **否** |
| 编排器是否已在审计响应前把 S0 标 done、open required=0 | **是（过程瑕疵）** → F-002 |

**模式本身满足**（独立意见已落盘）。**放行前提**仍是：`/govern` 汇总本意见、闭合 required findings 后，再确认 S0→S1。

## 7. Findings

### F-001 · Branding 首版字段相对 VP-007 静默收缩（favicon / 深浅色 Logo）

| 字段 | 值 |
|------|-----|
| 严重度 | **high** |
| 建议 | **required** |
| 状态 | open |
| 关联 | VP-007 交付范围 Branding；exit 3；D-002「Settings 四类字段」；`00-meta` 成功边界 3 |

**描述**：VP-007 明确 Branding 含 **Logo、favicon、浅色/深色 Logo**（URL 形态，非上传）。D-002 仅冻结 `logoUrl`，并写「深/浅色 Logo 首版仅保留单 logoUrl」；**favicon 完全未出现**，亦无 `accepted-residual`（范围、缓解、责任人、复审触发）。此为对 VP 方向范围的收缩，不能仅凭实现便利静默成立。

**证据**：
- `docs/vision/plans/VP-007-localization-and-system-settings.md` 交付范围表 Branding 行
- `01-decision/D-002-s0-contract-freeze-info-gates.md`「冻结的其他契约」§3
- 现状代码仅 `siteTitle`/`logoUrl`（`handler/settings.go`）——现状可解释起点，不能单独构成对 VP 的合法 residual

**建议闭合路径**（三选一，须书面）：
1. **fixed**：把 `faviconUrl` 与（若仍要）`logoUrlLight`/`logoUrlDark` 写入 S0 冻结字段表与后续 S3 方案；或  
2. **accepted-residual**：用户书面接受本波次仅 `logoUrl`、排除 favicon/分主题 Logo，含范围/缓解/责任人/复审触发；或  
3. **VP 修订**（`/vision`）：正式改写 VP-007 Branding 首版字段后再冻结 Goal 契约。

### F-002 · S0 在 independent 意见响应前已标完成，且「开放 required=0」不一致

| 字段 | 值 |
|------|-----|
| 严重度 | **med** |
| 建议 | **required** |
| 状态 | open |
| 关联 | P-002/P-003；D-002 审计模式；`00-meta` S0=done；`goal-tree.md`；本轮前 `03-audit.md` 空台账 |

**描述**：D-002 规定 S0 契约冻结审计模式为 `independent`。本轮审计前：`00-meta` 已将 S0 标 **done**、`progress: 1/6`；`goal-tree.md` 写「S0 已完成」且「开放 required：Goal 审计 0」，同时又写「A-001 意见待编排器响应」。索引 `03-audit.md` 信息就绪区仍写 I-L10N `open`（与 `01-decision`/`00-meta` 的 verified 不同步）。在独立意见落盘并响应前宣称阶段完成与 open required=0，违反「未合法闭合的 required 不放行」的操作纪律。

**证据**：`00-meta.md` 路线图表；`goal-tree.md` 维护说明；本轮前 `03-audit.md` 空索引 + 信息表仍 open。

**建议**：`/govern` 响应本 A-001 后，再确认是否维持 S0 done；在 F-001 等 required 未闭合前，不得以「审计 0 开放」放行 S1 实施。

### F-003 · I-L10N-004 `verified` 易被误读为 VP exit 5 已关闭

| 字段 | 值 |
|------|-----|
| 严重度 | med |
| 建议 | **recommended** |
| 状态 | open |
| 关联 | I-L10N-004；VP-007 exit 5；D-002 §I-L10N-004 |

**描述**：信息门禁关闭（选定路径 a + envelope 冻结）**可以**标 verified；但 exit 5 路径 (a) 仍要求 S4「协商 / Content-Language / 失败回退」实施证据。建议在信息表或 D-002 加醒目标注「exit 5 证据 = S4，非本 verified」。

### F-004 · E-002 里程碑 checkpoint hash 仍为占位

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| 状态 | open |
| 关联 | E-002「里程碑 checkpoint」；P-002 checkpoint 纪律 |

**描述**：`E-002` 写 `commit：<S0 checkpoint hash 待填>`。S0 文档产物已列路径，但 checkpoint 事实不完整。建议在 S0 响应/提交后回填真实 hash（仅 owned paths）。

### F-005 · 错误码「约 35」未钉死可回归枚举

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| 状态 | open |
| 关联 | I-L10N-004；S4 错误码契约测试 |

**描述**：E-002/D-002 称「约 35 码」。独立抽查 `writeError(..., "<CODE>", ...)` 字面量约 **32** 个稳定码（另有 kernel/migration 等非 HTTP writeError 族）。S0 盘点方向正确，但 S4 契约测试需要**显式枚举表**（纳入/排除、编目 vs INTERNAL 回退），否则「码集合不变」不可回归。

## 8. 必改项汇总（required）

1. **F-001（high）**：调和 Branding 字段与 VP-007（恢复 favicon/分主题 Logo 冻结 **或** 用户书面 residual **或** 修订 VP）。未闭合前，不得将 S0 契约视为与 VP-007 无冲突的最终冻结面。  
2. **F-002（med）**：`/govern` 正式响应本 A-001；在 required 闭合前不得宣称 Goal 审计 open required=0 并放行 S1；修正索引/表述一致性（信息表与审计台账）。

## 9. 与既有意见的异同

- 本目标 `03-audit` 此前无 A 条目（scaffold 模式 `none`，D-001）。  
- 本 A-001 为 S0 首条 **independent** 意见；无 self 历史可对照。

## 10. 结论与给编排器/用户的下一步

**verdict: conditional** — I-L10N-001～005 的用户书面关闭与 F-V029 分母冻结整体扎实，v2.7.0 兼容策略方向正确，independent 模式已由本意见落实；但 **Branding 相对 VP-007 的字段收缩未 residual/未改 VP（F-001）**，且 **S0 完成宣称早于独立审计响应（F-002）**，不可无条件放行 S0→S1。

**建议下一步（`/govern`）**

1. 展示 F-001/F-002，请用户裁决闭合路径并留痕。  
2. 闭合 required 后，再确认 S0 检查点与是否进入 S1。  
3. recommended：F-003 标注 exit 5 边界；F-004 回填 checkpoint；F-005 钉死错误码枚举（可并入 S4 方案，不必阻塞 S1 核心，除非用户要求）。

### 声明

本意见 `source: independent`，**不修改** status / progress / 方案正文 / goal-tree 状态列。响应、修正与放行由 **`/govern`** 处理。
