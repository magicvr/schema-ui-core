# attachments · GOAL-036（W25）

本目录存放诊断与验证证据：

| 文件 | 内容 |
|------|------|
| `I-001-evidence.md` | 浏览器 e2e 回归证据：修复前失败点、根因链（含 A-001 独立审补充的 CASCADE/FK 机制归因）、修复与修复后 9/9×2 |
| `I-002-evidence.md` | 活栈计时复核：双栈（基线 `0878d7f` vs 当前）请求数与呈现耗时原始数据 |

测量脚本（`measure.js` / `run-stack.ps1`）位于会话临时目录 `%TEMP%\w25-metrics\`，未入库（A-001 F-004 recommended 记录在案）；如需可复现基准，建议后续纳入 dev 脚本。