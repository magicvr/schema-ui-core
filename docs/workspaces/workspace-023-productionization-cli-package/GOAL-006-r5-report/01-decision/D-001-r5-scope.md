# D-001 · R5 范围定案（2026-08-29）

## 决策

1. **breaking 演练**：不向 registry 实发 breaking 版本（影响真实消费面）——以流程演练交付（semver-breaking-policy §3 + changelog 模板 + R2 升级演练对照）；**实演** = go 后首个 major 发布时执行（复审触发登记）。
2. **CLI 上手计时**：引 R2 实测（demo-admin create→双端绿 分钟级）+ 本轮复测（golden-field 从零走查计时）。
3. **默认主路径建议**：报告给出（cli+包 为默认主路径候选）；**不改 Charter 措辞**（用户既定：fork 并存表述维持）。