---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-006-r5-release-and-gono-go
version: 0.1.0
---

# D-001 · 发布通道定案与 go/no-go 裁决（用户）

## 裁决

用户 2026-08-29 书面裁决（两项）：

1. **go/no-go = GO · 推进 Charter strategic 修订**（VR-050 已执行：Charter `@0.3.0` + pin `v2.9.0` + 22 VP re-align；VRev-050 pass）。
2. **I-003 发布通道定案**：npm = 本地 tarball 产物（`scripts/pack-npm-packages.mjs` → dist-lib/artifacts），上传公开 registry 属 go 后外部动作（凭据/网络）；Go = 单模块粗粒度 `git tag`（本地试点 `v0.0.2`）+ proxy 发布流程文档。
3. **I-007 闭合**：pin `v2.8.0 → v2.9.0`（`81aa1d8`；Charter/roadmap/provenance 同步）。

## 未选方案

- 不选 GitHub Packages/npmjs 上传为本期发行动作（无凭据、需外部网络；发布脚本已就绪可切换 target）。
- 不选 multi-module 细 tag（G2）——G1 保持。