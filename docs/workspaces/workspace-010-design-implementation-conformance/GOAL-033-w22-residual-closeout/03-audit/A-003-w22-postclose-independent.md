---
id: A-003-w22-postclose-independent
doc: audit-entry
goal: GOAL-033-w22-residual-closeout
source: independent
date: 2026-08-23
scope: W22 修复结果复核（A 组 ×6 + B 组 ×6 + H ×3 的代码/证据与源台账回写全链）+ 关门叙事与台账卫生
audit_type: close-out
verdict: conditional
---

# A-003 · W22 修复结果独立复核（post-close）

## 范围与区间

- 核对提交 `85a733f`（fix(api,web): 修复 W22 残余问题并完成 GOAL-033 结项，工作树干净）的代码面：A2 迁移 0049、A3 导航边界测试、A4 密码可见切换、A5 上传嗅探、A6 MFA verify 限流、H3 dogfood 改名、e2e 选择器修复；
- 复核 A1 补跑日志（`attachments/e2e-admin-m3-rerun-final.log`）与本机端口排除区间现状（`netsh` 复跑）；
- 定向复跑测试：`internal/kernel`（Navigation）= ok 0.709s；`internal/store`（Migrate0049）= ok 2.451s；`internal/handler`（MFAVerify|ContainsActiveContent|UploadA5）= ok 3.764s；
- 逐项核对 A/B/H 各残余的**源台账回写**（Q2 路径 grep + 精读），与 D-001 用户书面批准范围对照；
- 核对本目标五件套台账结构（01-decision / 02-execution / 03-audit 索引完整性、P-005 信息表、goal-tree、workspace.md 波次表）。

## 成果（有证据）

| 项 | 核对结果 | 证据 |
|----|----------|------|
| A2 迁移 0049 | 真实落地：`migration.go:366,437` `migrateSeedAdminMustChangePassword`；`migrate_0049_test.go` 双测试；黄金断言五处同步（identity/migrate/operations/restart） | 提交 85a733f；store 定向复跑 ok |
| A3 导航边界测试 | 纯测试新增（`navigation_order_test.go` ×3），生产行为零改动 | 提交；kernel 定向复跑 ok |
| A4 密码可见切换 | `LoginPage.tsx:285,292-304` type 切换 + Eye/EyeOff + `[data-password-toggle]` + i18n 双语 key + 组件测试 +3；自引入 e2e strict-mode 选择器回归已波内修复（5 处 `exact:true`） | 提交；登录面代码精读 |
| A5 上传嗅探 | `upload.go:96-142` 8 KiB 头窗口 + `<svg`/`<script`/`<?xml` + `(?i)\bon[a-z]+=` 启发式，拒绝在 MIME 之上；下载头不变；端到端/边界测试齐 | 提交；handler 定向复跑 ok |
| A6 verify 限流 | `mfa.go:81-107` 独立桶 `newLoginRateLimiter(15m, 10, 1<<16)`，键 = client IP，429 + Retry-After + RATE_LIMITED；3 条测试 | 提交；handler 定向复跑 ok |
| H3 dogfood 改名 | `app-manifest.(admin\|mvp).json` → `-dogfood.json`（纯 rename，0 字节变更）；9 处引用覆盖 8 个测试文件 | 提交 |
| A1 补跑 | 环境阻塞解除属实：本机 `netsh` 现无 8011–8110 区间；`rerun-final.log` 显示 API 25080 / Vite 25173 启动成功、M1 断言失败（`localization.spec.ts:62` 停留 `/`）与 E-006 叙事一致；基线实验（stash 后 HEAD 同败）结论置信 | netsh 复跑；两份附件日志 |
| B 组回写 | 全部落痕：W3 A-011 F-003「fixed」、GOAL-005 R4-I004 与 GOAL-006 C1-I003 到期复核、GOAL-010 A-003 C4-004 与 GOAL-012 E-004 续期、GOAL-011 A-004 F-IND-C5-002「fixed」、W8 S1 ledger F-007 closure、W13 Root I-001/I-004 回写 + D-002 追认节 | 各源文件 grep/精读 |
| D-002 裁决 | B1 续期 2027-02-01 / B3 续期 2026-12-01 / W13 追认三项均含 owner、复审日期、触发；用户书面（ask_user_question） | `01-decision/D-002-w22-p004-adjudications.md` |
| H1/H2 零写入 | 成立：GOAL-018 `done 5/5`、D-005 决策在；W17 N-001 已 closed（note）且独立 A-002 认可 | W11 GOAL-018 meta、W17 A-001 |
| 回归证据 | E-006 记录的 go build 0 / store 45.044s / vitest 76 文件 1088 / tsc 0 / build 0 与提交一致；我另定向复跑三项绿 | E-006 + 复跑输出 |
| goal-tree | W22 树节点与状态表 done 18/18 已同步 | goal-tree.md:53,58,135 |

## Findings

| 编号 | 级别 | 描述 | 证据路径 |
|------|------|------|----------|
| F-001 | **required · med** | **A 组源台账回写缺失 ×6（用户书面批准范围未交付）**。D-001 用户裁决明确「A 组全部 6 项…含代码 + 测试 + **源目标台账回写**」；E-002 亦写明 A3「待 S4 回写源台账转 fixed」。实际六个源台账 residual 行均未追加任何闭合注记（无 2026-08-23 标注、无 GOAL-033 引用），仍为 accepted-residual：A2 `GOAL-025/A-002-self-response.md:24`（F-003）、A3 `workspace-011 GOAL-013/E-006-s5-closeout.md:25` 与 `A-003-s5-independent-closeout.md:218`（F-004）、A4 `workspace-006 Root/A-012-response-a011-and-closeout.md:41`（F-VUI-011）、A5 `workspace-009 GOAL-002/03-audit.md:33`（N-002）、A6 `workspace-009 GOAL-011/A-004-w11-closure-response.md:62`（R-001）、H3 `GOAL-002-w1/A-006-w1-closeout-response.md:35`（F-006）。对照：B 组全部落痕、A1 源台账已于 2026-08-09 discharge（无需补）。后果：残余真相源仍显示「开放残余」，未来 accepted-residual 全库扫描将再次命中同一批项目——正是本波旨在消除的失效模式。 | D-001 `01-decision/D-001-scope-and-p004.md:18`；E-002:16；六个源文件行 |
| F-002 | recommended | **02-execution.md 索引缺 E-006 行**（索引仅 E-001～E-005，而 E-006 是 C17 全量回归 + C2 补跑结论的唯一载体）；连带 A-001 自审「执行索引 E-001～E-006 与事实一致」与文件实际不符。 | `02-execution.md:16-20`；A-001:24 |
| F-003 | recommended | **01-decision.md（索引文件）整篇复刻 D-002 正文**，与 ledger 权威版 `01-decision/D-002-w22-p004-adjudications.md` 存在实质文字分歧：索引版缓解句多「下载/审计面既有控制」半句 → accepted-residual 范围表述双版本；且缺 D-001/D-002 索引行。 | `01-decision.md` vs `01-decision/D-002-…md`（Compare-Object 证实 DIFFERENT） |
| F-004 | recommended | **00-meta P-005 信息表 I-003/I-004 仍标 open**，但 E-003/E-004 已给出结论与证据（I-003：B5 freshness trigger 为事件型「每个后续业务 VP 激活前」、未触发、维持有效；I-004：W13 Root 双账回写 verified）。目标已 done，required 信息项留 open 与事实不符。 | `00-meta.md:73-74`；E-003/E-004 |
| F-005 | recommended · low | **workspace-010 `workspace.md` 波次表未登记 W22/GOAL-033 行**（updated 2026-08-22 止于 W21），与 goal-tree 波次档案不一致；另 N-001（`/dashboard` 断言漂移）移交仅存于 A-001，未在波次档案/roadmap 登记追踪槽。 | `workspace.md:46-78`；A-001 |

## 必改项汇总

- F-001（required）：按 B 组先例，在六个源台账 residual 行追加「2026-08-23 兑现复核：fixed（GOAL-033 E-005 证据）」注记；不改历史原文。
- F-002～F-005（recommended）：补 E-006 索引行并更正 A-001 表述；01-decision.md 改索引 + 统一 D-002 缓解句（若「下载/审计面既有控制」为意图范围应同步进 ledger 版并留痕，否则删分歧句）；刷新 I-003/I-004 状态；补 workspace.md W22 行（可选：登记 N-001 移交槽）。

## 与既有意见的异同

- **与 A-001（self · pass）**：同意代码/测试/回归证据扎实、B 组回写齐全、关闭叙事大体合法；**不同意**「执行索引 E-001～E-006 与事实一致」一条（F-002）；且 A-001 未核出 A 组源台账回写缺失（F-001）。
- **与 A-002（independent · pass，安全面）**：A5/A6 的 diff 级结论一致（含 R-A5-1/R-A6-1 recommended 记录在案）；本意见未发现 A5/A6 实现缺陷。
- A1 处置叙事（履约关闭而非 fixed）+ N-001 移交：认可，非必改。

## 结论 + 建议给编排器/用户

修复结果**主体真实可复核**：六项代码修复全部在位且定向复跑通过，B 组复核与 D-002 裁决落痕完整，关门流程经用户书面确认。但用户书面批准范围中的「源目标台账回写」有 6 项未交付，残余真相源仍为开放状态（F-001 required），另有 4 项台账卫生缺口。**verdict：conditional**——不推翻已关闭结论与代码修复，但 F-001 须在响应中补齐（或由用户书面裁决其他闭合路径），并处理 F-002～F-005。

建议 `/govern` 下一步：① 响应本 A-003；② 补 A 组 6 处源台账回写注记（F-001 fixed）＋ 四个台账小项；③ 刷新本目标台账后复审 A-003 关闭；④ 若拒绝任一 recommended，按 P-004 留痕。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。