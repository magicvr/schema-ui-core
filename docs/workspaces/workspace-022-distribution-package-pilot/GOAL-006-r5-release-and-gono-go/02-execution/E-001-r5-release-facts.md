---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-006-r5-release-and-gono-go
version: 0.1.0
---

# E-001 · 发布流水线与消费回归事实（2026-08-29）

## S1 · 发布流水线

- `scripts/pack-npm-packages.mjs`：扫描 `dist-lib/@schema-ui/*` → `pnpm pack` → `dist-lib/artifacts/*.tgz`（一键可复现）。
- 产物：`schema-ui-protocol-0.2.0.tgz`（73 KB）· `schema-ui-renderer-0.1.0.tgz`（99 KB）。
- Go tag：`git tag v0.0.2`（本地；proxy 发布为 go 后外部动作）。

## S2 · golden 消费回归（tarball = registry 安装语义）

- golden-web `package.json` 依赖改为 **tarball 路径**（`file:*.tgz`）→ `pnpm install`（快照复制进 store，与 `pnpm add @schema-ui/protocol@0.2.0` 语义一致）。
- 结果：protocol probe PASS · render probe PASS（1573 B）· token override PASS · V2 能力 `normalizePageID(' Roles ') → 'roles'` ✅。

## S3 · go/no-go 报告

`attachments/gono-go-report-v1.md`：判据 #1–5 全达成 · 触发框架判向 = **倾向推进** · Charter 修订草案 + pin 建议。用户裁决 = **GO**（VR-050 已执行）。

## 闭环

- freeze-face v1.2.0（Web 六包边界 + peer 矩阵 + 发布形态注记）随本目标定稿。
- 残余：npm registry 上传 / G2 细化 / F-006 d.ts 等 → go 后清单（报告 §5）。