---
id: GOAL-004-r3-seam-alignment
title: R3 接缝与对齐
status: done
parent: GOAL-001-event-bus-port
created: 2026-09-01
updated: 2026-09-01
version: 0.2.0
progress: 4/4
---

# GOAL-004 · R3 接缝与对齐

## 概述

承接 GOAL-001 纲领路线图 R3 阶段：接缝声明与对齐登记。**对象面**：outbox/MQ 运输接缝声明 + Admin typed domain event gated 对齐 + topic/订阅命名与契约测试 harness。**红线**：不实现 outbox/broker；不解除 Admin gated；不消耗 RT-Q02 trigger。

## 依赖

- R1（GOAL-002）：done ✅
- R2（GOAL-003）：done ✅
- I-028-004（required）：待确认（事件类型注册权属与 Admin gated 保持）

## 成功标准（检查点）

R3 覆盖判据 #3/#4/#5 + I-028-004：

- [x] **C1**: outbox/MQ 运输接缝声明落盘（判据 #3）
  - 声明应用契约 vs 运输实现边界（进程内 channel vs outbox 持久化 vs MQ broker）
  - 明确：本 VP 不引入 broker 客户端依赖、不实现 outbox、不预裁 RT-Q06 表结构
  - 位置：D-001 §1

- [x] **C2**: Admin typed domain event gated 对齐登记（判据 #4 + I-028-004）
  - 确认事件类型注册权属（业务域 VP vs Admin 功能 VP）
  - 确认本 VP 不解除 Admin 功能分支 typed domain event 扩展接缝的 trigger-gated
  - 登记对齐声明：D-001 §2
  - **I-028-004 verified**（2026-09-01 用户确认）

- [x] **C3**: topic/订阅命名与契约测试 harness（判据 #5）
  - topic 命名约定（`<domain>.<aggregate>.<event>`）
  - 订阅命名/生命周期约定
  - 契约测试 harness（D-001 §3）

- [x] **C4**: 自审与独立审计
  - A-001 self **pass**（0 required findings）
  - A-002 independent **deferred**（工具链受阻，待补审）
  - 接缝声明完整性、对齐登记准确性、未越界验证通过

## 执行计划（待细化）

R3 是声明与登记阶段，无编码实施。主要工作：
1. 澄清 I-028-004（lead 建议 + 用户确认）
2. 编写接缝声明（D-001）
3. 登记对齐声明（与 Admin VP 协调或本 VP 决策）
4. 编写命名约定（D-002 或架构短文）
5. 自审与独立审计

## 父目标

- GOAL-001-event-bus-port

## 备注

- R3 为纯声明阶段，不涉及代码实现
- I-028-004 从 non-blocking 升为 required（因 R1 选注册表）
- 审计模式：阶段关门 default self；按需 independent
