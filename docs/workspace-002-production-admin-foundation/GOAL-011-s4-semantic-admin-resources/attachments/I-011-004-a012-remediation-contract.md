---
title: I-011-004 · A-012 生产边界整改契约与验收矩阵
status: active
doc_type: contract
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-011-s4-semantic-admin-resources
version: 0.1.0
related_info: I-011-004
related_decision: D-006
supersedes: null
---

# I-011-004 · A-012 生产边界整改契约（冻结）

> **性质**：回答 A-012 F-001～F-005 如何走 `fixed`，并冻结用户对 F-003“补齐角色授权/grant 管理路径”的产品边界。D-006 已将本信息项置为 `verified`；这只解除方案选择门禁，不构成实现、验收或 finding closure。

## 1. 范围与兼容边界

- 本契约是 I-011-001 v0.2.0 与 I-011-003 v0.2.0 的 A-012 限定补充；历史 S1～S5 收据和旧版本契约保持可追溯，不重写为当时已具备以下能力。
- 在本整改范围内，I-011-001 §3.4/§8 中“自定义角色无 grants / grant 管理界面非目标”由本契约 §4 取代；交付的是 users/roles 所需的最小授权闭环，不是完整 IAM。
- I-011-003 的 Renderer 零 diff 是 2026-08-03 S4 历史基线事实。本次允许增加通用、结构化的密码字段传输与行级 action 禁用能力，并进行资源中性重命名；不得为 users/roles action id 写特判。
- F-006（legacy roles JSON 双写）继续为 recommended/non-blocking，不纳入本次 required closure。

## 2. F-001 · 角色委派边界

1. 新增显式权限 `roles.assign`，默认只授予 `admin`；`users.write` 不再隐含角色委派权。
2. create 携带非空 `roles` 或 patch 提交 `roles` 字段时，操作者必须具有 `roles.assign`；缺失时返回 403 `ROLE_ASSIGNMENT_FORBIDDEN`。
3. 目标角色有效权限的并集必须是操作者当前有效权限的子集；操作者不能借用户管理授予自己没有的权限。
4. 分配或保留 `admin` 角色只允许当前 `admin` 操作者；self-demote、最后管理员保护继续生效。
5. 角色引用仍必须先存在；API 不得隐式创建角色。正向与负向测试至少覆盖 users.write-only actor、超集委派、admin 委派与合法子集委派。

## 3. F-002 · 密码与会话生命周期

1. Schema 使用专用 `password` 控件；通用表单提交与资源解码均按字段类型保留密码原始字符串，不做 trim、大小写或 Unicode 归一化。
2. create/patch 密码必须是 string，UTF-8 字节长度为 8～72 bytes，且不能全为空白；失败返回 400 `INVALID_PASSWORD`。72 bytes 与 bcrypt 输入上限对齐。
3. 密码仍为仅写字段，不进入任何响应或 operation log。
4. 修改密码时，在更新用户的同一数据库事务内撤销该用户全部 refresh token；旧 refresh token 随后必须失败，新密码登录与新 refresh token 必须成功。
5. 已签发 access token 不建立服务端撤销表，仍在其既有短 TTL 到期前有效；该边界必须在契约与测试中显式记录，不得把 refresh 撤销表述为 access token 即时失效。

## 4. F-003 · 用户角色与 grant 管理闭环

1. users Schema 提供三个独立写路径：编辑资料、管理角色、修改密码。角色输入必须进入真实 `roles` 请求体；空密码不得因普通编辑被提交。
2. roles list/detail/create/update 响应增加：`permissions`（permission key 数组）、`menuItems`（menu id 数组）、`assignedUsers`（整数）、`editable`/`deletable`（boolean）。字段均来自持久化关系或可验证派生，不使用前端猜测。
3. 自定义角色 create/update 可在同一事务写入 `role_permissions` 与 `role_menu_items`；引用不存在分别返回 400 `INVALID_PERMISSION_REF` / `INVALID_MENU_ITEM_REF`。更新采用集合替换语义，省略字段表示保持不变，create 省略表示空集合。
4. system 角色仍不可修改/删除；in-use 自定义角色仍不可删除。Schema 通过通用结构化行条件读取 `editable`/`deletable` 来禁用动作，后端 409 不变量继续作为最终保护。
5. 自定义角色 grants 必须进入既有有效权限投影：分配该角色后，重新登录/刷新得到的身份具有相应 permission/menu；取消 grant 后新身份不再具有。系统角色种子 grants 仍由 seed 权威维护。
6. grant 管理只覆盖仓库内既有 permission/menu catalog 与自定义角色集合替换；权限目录 CRUD、条件策略、角色继承、SSO/SCIM、多租户、批量授权仍为非目标。

## 5. F-004 · records 活动面洁净度

1. 活动 CI 的 persistence/smoke 路径改为 `/api/users`，使用与真实 seed 一致的用户 id/总数变量；不得再引用 `/api/records` 或 `SMOKE_RECORD_ID`。
2. Web 通用 transport 文件和导出符号改为资源中性命名；生产 import 不保留 `records` 专名。
3. 活动 CI 增加静态洁净度检查，阻止 `/api/records`、`SMOKE_RECORD_ID` 与 Web 生产 transport 专名回流。
4. 历史 migration、migration compatibility 测试、历史治理文档和 operation-log legacy 值保留，不因静态检查被改写或删除。

## 6. F-005 · 最后管理员事务原子性

1. 删除路径的目标读取、角色判定、admin 计数、refresh token 清理与 users 删除必须在同一数据库事务/串行化边界内执行；不使用事务外角色快照决定删除。
2. 删除必须检查受影响行数，零行按 `USER_NOT_FOUND` 处理；任一失败整体回滚。
3. 可控并发回归必须证明：两个管理员被并发删除时最多一个成功，最终至少保留一个 admin；self-delete 与单管理员删除仍分别失败。

## 7. 验收矩阵

| Finding | 必须通过的证据 |
|---------|----------------|
| F-001 | handler/store 权限委派正负向测试；users.write-only 不能授予角色；超集/admin 越权失败；合法子集成功 |
| F-002 | API 密码边界/原字节登录测试；Web password 控件与不 trim 测试；改密撤销 refresh、新密码登录成功、access-token TTL 边界测试 |
| F-003 | roles grants store/handler CRUD 与有效权限投影测试；users/roles 真实 fixture 请求断言；system/in-use 行禁用；真实浏览器 users 角色/改密与 roles grant 路径 |
| F-004 | CI persistence 使用 users；smoke 新默认与变量一致；静态洁净度 gate；Web 通用 transport 重命名后测试/build 通过 |
| F-005 | 同事务删除实现；RowsAffected 断言；并发双管理员测试；原 self/last-admin 回归通过 |
| 全量 | `go vet ./...`、`go test ./... -count=1`、Web tests、TypeScript/build、Playwright、`git diff --check`；不可用的环境证据必须单独报告，不能写成通过 |

## 8. 关闭规则

- D-006 只决定 `fixed` 路径。完成实现与本地验证后，`/govern` 可记录候选修复收据，但不得自行把 A-012 independent finding 写成已独立验证。
- 必须请求仅覆盖 A-012 F-001～F-005 的 `/audit` finding-closure 复审。只有复审逐项确认且无新增 required，才可把对应 finding 标记为 `fixed` 并评估 GOAL-011 重新关门。
- F-006 不阻断上述复审或关门；若本轮顺带处理，仍须单独记录，不得扩大 required scope。

## 9. 修订记录

| 版本 | 日期 | 变更 |
|------|------|------|
| 0.1.0 | 2026-08-04 | D-006 冻结：F-001～F-005 全部 fixed 路径；F-003 选择完整角色授权/grant 管理；定义实施与验收边界 |
