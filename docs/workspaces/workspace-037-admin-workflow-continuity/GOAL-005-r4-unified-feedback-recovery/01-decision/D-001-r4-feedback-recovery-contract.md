---
id: D-001-r4-feedback-recovery-contract
doc: decision
status: accepted
goal_id: GOAL-005-r4-unified-feedback-recovery
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-005-r4-unified-feedback-recovery
version: 1.0.0
---

# D-001 · R4 反馈分类、重试与无障碍合同

承接 Root R1 D-005，R4 采用以下合同：

- 用户可见文案优先 catalog `messageKey/params`；缺少已知 key 时使用安全的既有 fallback，不把 token、密码、原始请求体或敏感服务端载荷写入反馈表面；诊断 code/correlation 只能作为可定位信息。
- 成功操作只产生一次短暂 `role=status` 反馈；错误保持 `role=alert` 直到用户 dismiss 或页面状态被新结果替换；同一结果不会因重渲染重复宣告。
- 只有幂等读取（列表、statCard/chart、页面 Schema/manifest 等读取路径）提供显式 retry；retry 发起一次新的读取并清理旧错误，不自动循环，不附着到写入成功/失败结果。
- 写入失败不自动 retry；表单值、fieldErrors 与 dirty 状态保留，用户再次点击提交才产生新的写请求。409/403/维护/不可用/离线/超时不伪装为成功。
- 普通 resource 的 operational failure 与 HostFailureScreen 的认证、协议、渲染终态保持边界；Host 终态继续由 Host recovery action 处理，不由普通 Toast 代替。
- dismiss、retry 与错误/成功呈现必须键盘可达，使用既有 `status/alert` 角色合同；不因 retry 或倒计时造成焦点反复抢占。

本决定不引入新的 API、错误 envelope、第二套全局状态管理或 R5 组合验收范围。
