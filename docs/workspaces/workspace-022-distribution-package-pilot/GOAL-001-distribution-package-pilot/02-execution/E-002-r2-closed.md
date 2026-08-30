---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-001-distribution-package-pilot
version: 0.1.0
---

# E-002 · R2 完成事实（2026-08-29）

1. **internal 外移重构**（方案 A · 用户裁决）：`internal/kernel` → `apps/api/kernel`、`internal/modules` → `apps/api/modules`；223 文件 import 改写 + playbook 同步；build 0 + 全量测试回归（1 个 PG drain 全量并发时序敏感单跑 PASS，A-001 F-003 登记）；freeze-face v1.0.1 勘误。
2. **C 层泄漏收敛**（方案 β · 用户裁决）：新增 `apps/api/assembly` 公开装配工厂（OpenStore / NewAuthenticator / NewMailSender）；机制 = Go 类型推断不命名消费 + `kernel.Store.Run` ⊇ `authsession.TxRunner` 结构同构。
3. **黄金下游仓闭环**（GOAL-003 A-002）：`golden.local/consumer` 装配 users 全链——迁移台账收集 → SQLite Open（fresh=true 迁移从零 apply）→ Authenticator/Repository/Recorder/MailSender → `RegisterContributions` 贡献 = Descriptor 声明（routes=10 pages=2 perms=3 nav=1 frag=1）；**零 internal 命名**。
4. **契约收口**：freeze-face v1.1.0（B+ 层增列 + B 层盘点引用）；F-001/F-002 fixed（GOAL-002 A-001 响应回填）。
5. **判据 #1 满足声明**；Root progress 1/5 → 2/5；GOAL-003 done 4/4、GOAL-002 done 4/4。
6. 残余：PG external 消费（F-005，R4/R5 复审）；assembly experimental 签名 (v0.1.0)；A-001 F-003（PG drain 全量并发时序）随 R5/VP-009。

下一步：R3 立项（Web 包闭环）——npm 包组拆分（protocol/renderer/shell/ui + 主题覆盖）+ 空下游 app 渲染 + Token 覆盖（I-002）。