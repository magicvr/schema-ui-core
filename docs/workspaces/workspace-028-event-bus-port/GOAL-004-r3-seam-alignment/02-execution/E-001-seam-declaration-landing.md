---
id: E-001
date: 2026-09-01
scope: R3 C1/C2/C3（接缝声明 + 对齐登记 + 命名约定）
summary: D-001 决策落盘 + I-028-004 用户确认闭合
status: completed
---

# E-001 · R3 接缝声明与对齐决策落盘

**日期**：2026-09-01  
**执行者**：实现团队

## 实施内容

### 1. 接缝声明文档（C1 判据 #3）

**产物**：`01-decision/D-001-seam-declaration-and-alignment.md`

**内容**：
- **三层架构边界**：应用契约层 ↔ 进程内运输层（本VP已交付）↔ 持久化运输层（gated）
- **接缝约定**：
  - 进程内 ↔ outbox：transactional outbox writer 插入点 + JSON serializable 约束
  - outbox ↔ broker：poller + re-publish 路径
- **红线遵守验证**：
  - ✅ 不引入 broker 客户端依赖
  - ✅ 不预裁 RT-Q06 outbox 表结构
  - ✅ 不消耗 RT-Q02 trigger
  - ✅ Profile 默认集 / 模块矩阵 / Manifest 未变更

**验证**：grep 代码库未发现 Kafka/RabbitMQ/NATS 依赖；outbox 表未创建

### 2. Admin gated 对齐登记（C2 判据 #4 + I-028-004）

**产物**：D-001 §2 + Root Goal I-028-004 闭合

**决策内容**（用户 2026-09-01 确认）：

| event 类别 | 注册责任 | trigger 状态 | 本 VP 角色 |
|-----------|----------|-------------|-----------|
| 系统级事件 | 本 VP 或基础设施 VP | 无 gated | 可预注册（当前未注册） |
| 业务域事件 | 各业务域 VP | 视业务域 VP 而定 | 不负责；仅提供 Register API |
| Admin 功能事件 | Admin 功能分支 VP | **trigger-gated** | 不负责；**不解除 gated** |

**对齐声明**：
- 本 VP 不解除 Admin typed domain event 扩展接缝的 trigger-gated
- 事件 schema 定义 / Register 调用时机 / topic 命名归 Admin 功能 VP 负责
- 不预置任何业务域/Admin event schema（R2 Memory `topics map` 初始为空）

**I-028-004 闭合路径**：
1. D-001 §2.2 提出 lead 建议（注册权属划分 + 不解除 gated）
2. 用户确认（2026-09-01）
3. Root Goal 00-meta.md I-028-004 状态：待确认 → **verified**

**验证**：
- grep `kernel/eventbus.go`：无预注册 topic
- grep `composition.go` `newEventBus`：无业务 event schema
- D-001 明确声明不解除 Admin gated

### 3. topic/订阅命名与契约测试 harness（C3 判据 #5）

**产物**：D-001 §3

**命名约定**：
- **topic 格式**：`<domain>.<aggregate>.<event>`（小写，点分隔，event 过去式）
- **正则约束**：`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`（继承 R1 D-002）
- **示例**：
  - 系统级：`system.startup`, `system.config.reloaded`
  - 业务域：`profile.created`, `schema.published`
  - Admin（gated）：`user.invited`, `role.assigned`

**订阅生命周期约定**：
- 创建：fx lifecycle OnStart 中调用 Subscribe
- 销毁：OnStop 中调用 Unsubscribe（可选，Stop 自动 close）
- 并发：同 topic 可多订阅者，payload 独立 copy，无顺序保证

**契约测试 harness**（D-001 §3.3）：
- event schema JSON serializability 测试模板
- Register 幂等性测试模板
- Publish/Subscribe 集成测试模板

**登记位置**：D-001（不纳入 Redis key 轨道）

**验证**：D-001 §3 完整定义命名约定与测试模式

## 产物清单

| 文件 | 类型 | 行数 | 说明 |
|------|------|------|------|
| `01-decision/D-001-seam-declaration-and-alignment.md` | 决策文档 | ~400 | 三层架构边界 + 对齐声明 + 命名约定 |
| Root Goal `00-meta.md` I-028-004 | 信息项更新 | 1 行修改 | 待确认 → verified |

## 执行验证

```bash
# 验证 D-001 文件存在
ls docs/workspaces/workspace-028-event-bus-port/GOAL-004-r3-seam-alignment/01-decision/D-001-seam-declaration-and-alignment.md

# 验证无 broker 依赖
grep -r "kafka\|rabbitmq\|nats\|amqp" apps/api/go.mod  # 无匹配

# 验证无预注册 topic
grep -A 20 "func newEventBus" apps/api/internal/composition/composition.go  # 无 Register 调用

# 验证 I-028-004 闭合
grep "I-028-004.*verified" docs/workspaces/workspace-028-event-bus-port/GOAL-001-event-bus-port/00-meta.md
```

**结果**：全部验证通过

## 进度评估

| 检查点 | 状态 | 证据 |
|--------|------|------|
| C1（判据 #3 接缝声明） | ✅ 完成 | D-001 §1 三层架构 + 接缝约定 + 红线验证 |
| C2（判据 #4 + I-028-004 对齐） | ✅ 完成 | D-001 §2 注册权属 + 用户确认 + I-028-004 verified |
| C3（判据 #5 命名约定） | ✅ 完成 | D-001 §3 topic 格式 + 订阅生命周期 + 测试 harness |

R3 决策落盘阶段已完成，所有必需信息项（I-028-004）已闭合。进入审计阶段。

---

## 备注

- R3 为纯声明阶段，无代码变更
- I-028-004 从 non-blocking 升 required（因 R1 选注册表），本轮用户确认闭合
- 对齐声明明确：本 VP 不解除 Admin typed domain event gated
- 命名约定与测试 harness 由 D-001 提供，业务域 VP 按需采纳
