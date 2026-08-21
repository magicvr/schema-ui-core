---
id: E-005-w9-recommended-hardening
goal: GOAL-009-w9-api-web-security-audit
status: done
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-009-w9-api-web-security-audit
version: 0.1.0
---

# E-005 · W9 A-005 recommended 三项加固实施 + 回归（2026-08-21）

## 事实

用户指令"先修正 A-005 的三条 recommended"已执行；三项全部实施并带针对性回归锁。

| A-005 recommended | 实施 |
|-------------------|------|
| R-F-001 L2 校验器仅测试路径 | apps/web/src/renderer/render.tsx：生产渲染路径接入 validatePermissions —— 结构非法时 console.error（codes+paths）并对全部已注册权限目标 fail-closed deny；未标记目标保持协议默认放行。回归锁：render-l2-permissions.test.tsx（非法结构 → 按钮禁用 + console.error；合法结构 → 可用且无噪音） |
| R-F-002 恢复码 CAS 秒级令牌同秒窗口 | mfa/store/repository.go UpdateRecoveryCodesIfUnchanged 的 OCC 令牌从 updated_at（Unix 秒）改为前值 recovery_codes_hash 本身（compare-and-swap on value）——同秒并发双通过窗口在时序上不可能；service.go 消费方改传前值。回归锁：mfa/store/repository_test.go（前进/同步拒绝/旧令牌拒绝/新令牌落地） |
| R-F-003 缺原缺陷形状回归锁 | 新增 6 组锁：kernel/unique_violation_test.go（双方言文案+wrap+反例）；authsession/accounts_lock_test.go（F-004 计数逐次递增/阈值开锁清零/重计/未知用户）；mfa/store/repository_test.go（F-005 水位 CAS 四态）；jobs/runner_panic_test.go（F-007 panic → durable JOB_HANDLER_FAILED + panic 详情）；scheduledtasks/scheduler_panic_test.go（F-007 Execute panic → failed run 且循环存活）；scheduledtasks/store/cron_posix_test.go（F-025 POSIX OR/AND 十用例，含 1-31 全集等价性） |

## 回归证据（全绿）

- API：go vet ./... exit 0；go test ./... **exit 0 全部包通过**（含新增 6 组回归锁）。
- Web：npm test **75 文件 / 1077 测试全部通过**（含新增 render-l2-permissions 2 项）；npm run build exit 0。
- 过程记录：首轮全量曾因新测试自身缺陷失败（Go const 冒号等号语法、调度断言未计 panic 任务重试、Web 控制组夹具缺 permissions.inheritance capability、cron 用例星期事实错误），均已修正后复跑全绿——夹具修正过程本身验证了 L2 接线与 CAS 语义的判别力。

## 边界

- 本条为 recommended 加固，不改变 A-005/A-006 的 required 闭合结论。
- checkpoint：关门前最终验证后提交，hash 补记于 02-execution.md 追记节。