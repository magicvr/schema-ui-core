# Independent Audit Prompt for GOAL-003 R2

你是本项目的**独立审计员**（independent auditor），负责对 workspace-028 GOAL-003 R2 进程内 EventBus 实现进行交叉审计。

## 审计范围

GOAL-003 R2 全部检查点：
- **C1**: Memory EventBus 实现（`apps/api/internal/eventbus/memory.go` + 测试）
- **C2**: 配置层（`apps/api/internal/config/config.go` / YAML / env）
- **C3**: composition 注入（`apps/api/internal/composition/composition.go`）

## 审计依据

1. **契约**：workspace-028 GOAL-002 `01-decision/D-002-contract-v0.1.0.md`（R1 冻结契约）
2. **决策**：GOAL-003 `01-decision.md`（D-001 架构 / D-002 配置 / D-003 注入）
3. **执行记录**：GOAL-003 `02-execution.md`（E-001 / E-002）
4. **自审意见**：GOAL-003 `03-audit/A-001-self-audit.md`（F-001～F-004）

## 审计任务

请执行以下步骤：

1. **读取文档**：
   - `docs/workspaces/workspace-028-event-bus-port/GOAL-002-r1-contract-freeze/01-decision/D-002-contract-v0.1.0.md`
   - `docs/workspaces/workspace-028-event-bus-port/GOAL-003-r2-memory-impl/01-decision.md`
   - `docs/workspaces/workspace-028-event-bus-port/GOAL-003-r2-memory-impl/02-execution.md`
   - `docs/workspaces/workspace-028-event-bus-port/GOAL-003-r2-memory-impl/03-audit/A-001-self-audit.md`

2. **读取实现代码**：
   - `apps/api/internal/eventbus/memory.go`
   - `apps/api/internal/eventbus/memory_test.go`
   - `apps/api/internal/config/config.go`（EventBus 相关部分）
   - `apps/api/internal/composition/composition.go`（newEventBus / registerLifecycle OnStop）

3. **核心审计点**：
   - **契约符合性**：Memory 实现是否满足 D-002 的五节语义（Register / Publish / Subscribe / Unsubscribe / Stop）？
   - **线程安全**：mutex / atomic / WaitGroup / select 是否正确使用？是否存在数据竞态或死锁风险？
   - **Stop 语义**：是否正确 drain 已入队事件、等待 handler、拒绝新操作、尊重 ctx timeout？
   - **payload 隔离**：每订阅者是否真正收到独立 copy（不受其他 handler 修改影响）？
   - **panic 隔离**：handler panic 是否被 recover，且不影响其他订阅者？
   - **配置边界**：buffer_size 的 YAML / env / default / fail-closed 逻辑是否正确？
   - **composition 注入**：fx.Provide / OnStop 是否正确连接？
   - **测试覆盖**：测试是否覆盖关键路径（含 -race）？

4. **验证自审意见**：
   - F-001（select 优先级修正）是否真正解决了 Stop 后发送的问题？
   - F-002（测试竞态修正）是否消除了数据竞态？
   - F-003（双重验证）与 F-004（eventBusPort 未使用）的 accepted-as-is 是否合理？

5. **输出审计报告**（Markdown 格式）：
   ```markdown
   # A-002 · R2 进程内实现独立审计（independent · grok-4.6 high）
   
   ## 审计方法
   （简述审计步骤）
   
   ## 审计发现
   
   ### F-NNN · 标题
   **严重性**：required / recommended / informational  
   **范围**：文件名或模块  
   **发现**：（具体问题）  
   **建议**：（修正建议）  
   **状态**：open
   
   （如无新发现，说明"无额外发现"）
   
   ## 审计结论
   
   **verdict**: pass / conditional / fail  
   **理由**：（简述）
   
   ### 对自审意见的确认
   - F-001: ✅ 确认修正有效 / ❌ 仍存在问题
   - F-002: ✅ 确认修正有效 / ❌ 仍存在问题
   - F-003: ✅ 同意 accepted-as-is / ❌ 需修改
   - F-004: ✅ 同意 accepted-as-is / ❌ 需修改
   
   ### 开放 required findings
   - （列出所有 required 级别的未闭合 findings）
   - （若无，写"0 required findings"）
   ```

## 约束

- **本轮用户硬约束禁止落盘**：你的审计报告输出为文本，由调用者写入 `03-audit/A-002-independent-audit-grok.md`
- **独立性**：你的审计独立于自审（A-001），可以同意或不同意自审结论
- **verdict 判定**：
  - `pass`：无 required findings，可放行
  - `conditional`：有 recommended findings 但无 required，建议修正后放行
  - `fail`：有 required findings，必须修正
- **严重性分级**：
  - `required`（必改）：契约违反、安全漏洞、数据竞态、死锁风险
  - `recommended`（建议）：性能问题、代码重复、测试覆盖不足
  - `informational`（提示）：代码风格、文档完善

开始审计。
