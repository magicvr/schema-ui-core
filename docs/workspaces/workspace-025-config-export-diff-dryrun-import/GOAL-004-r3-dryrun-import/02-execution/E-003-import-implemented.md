---
status: recorded
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-004-r3-dryrun-import
version: 0.1.0
---

# E-003 · C2 import 实现与测试（2026-08-30）

1. **裁决**：2026-08-30 用户 GUI 裁决 **I-025-004 = 方案 A**（原子替换 + 应用前备份）→ `verified`；语义冻结于 `01-decision/D-002-import-write-semantics.md`。
2. **代码**：`apps/api/cmd/schema-ui/configpkg.go` 追加 `cmdConfigImport`——① 预检前置（复用 `dryRun`：结构 + env fail-closed，任一 fail → 拒绝、目标零触碰 exit 1）；② 生成目标文本（包 config 树 marshal + 生成头注释，敏感键不入文件）；③ 应用前备份 `<file>.pre-import.bak`（存在时；备份失败 fail-closed）；④ 写 `.tmp` → `server.LoadConfig(tmp)` 装载校验（fail-closed 同纪律）→ `os.Rename` 原子替换；⑤ 任一失败清 tmp、原文件未被触碰；⑥ 报告 checks/applied/backup + yaml/json + exit 0/1/2。`main.go` usage 更新（import 转正）。
3. **测试**：+4 用例（TestImportRoundtrip：export→import→re-export→diff 无差；TestImportBackup：`.bak` 内容 = 导入前旧值；TestImportRejectsAndKeepsUntouched：预检拒绝/坏包拒绝/装载校验失败 → 目标原样 + 无 tmp 泄漏；TestImportDefaultFile：缺省 `config.yaml`）——`go test ./cmd/schema-ui/` PASS（18 用例）。
4. **回归**：`go test ./...`（apps/api 全量 49 包）PASS。
5. **CLI 冒烟（实证）**：export → import（exit 0 · applied 报告）→ 再 export → diff `[]`（**跨实例往返一致 = 判据 #1 全闭环**）；二次导入生成 `.pre-import.bak`；生成文件含 `${APP_ENV:-development}` 引用与生成头。
6. **验证覆盖判据**：VP-025 判据 #3（dry-run 无副作用）/ #4（导入不破坏：tmp+rename 原子 · 失败路径目标原样）交付面完成。