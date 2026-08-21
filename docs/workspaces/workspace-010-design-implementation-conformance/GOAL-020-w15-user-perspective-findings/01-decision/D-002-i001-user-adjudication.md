---
id: GOAL-020-w15-user-perspective-findings
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-002 · W15 I-001 用户书面裁决：F01～F14 全部 in-scope（父目标等子目标）

## 背景

I-001（required · 最晚 S5）要求书面裁定 W15-F01～W15-F14 的 in-scope / defer、批次与 GOAL-020 是否等整改完成再关门。「推进直到闭门」本身不是该清单。2026-08-17 用户经 GUI 问询选择推荐项（见本文件）。

## 裁决（用户书面，2026-08-17）

1. **范围**：**F01～F14 全部 in-scope**，无 defer。
2. **批次 / 顺序**：按 D-001 建议 **A → B → C**，作为 **GOAL-020 下级子目标**渐进添加：
   - **批 A**：W15-F01（会话网络容灾）、W15-F02（表格重试）、W15-F04（404/405 JSON 信封）、W15-F05（CORS/安全头，**留本区**）、W15-F07（refresh 错误码）。
   - **批 B**：W15-F03（时间格式）、W15-F11（GET 幂等）、W15-F10（429 / 配额）、W15-F12（分页默认）。
   - **批 C**：W15-F06（改密后提示）、W15-F08（校验 reason 本地化）、W15-F09（error Toast 常驻）、W15-F13（当前会话标记）、W15-F14（细节缺口）。
3. **闭门等待**：GOAL-020 **不得**在子目标完成前标 `done`。全部整改子目标 `done` 后才可关门。
4. **方案冻结（随本裁决）**：
   - **F03**：统一输出 RFC3339 毫秒串；**不改字段名**（`created` 保持 `created`），避免契约破坏。
   - **F05**：本区实施，不移交 workspace-009。
   - **F11**：GET 只读；缺失返回 404（`WALLET_NOT_FOUND`）；创建仍由 POST 触发。记录 go-impact：改变 workspace-011 钱包自动开户读路径。

## 影响

- I-001 → **closed**（P-004 ✓）。
- S5 检查点 = 本裁决 + 子目标结构落盘；整改实施在 GOAL-021 起的下级目标。
- 父目标进度分母扩为 S1～S5 + R1～R3（8 个等权检查点）。
