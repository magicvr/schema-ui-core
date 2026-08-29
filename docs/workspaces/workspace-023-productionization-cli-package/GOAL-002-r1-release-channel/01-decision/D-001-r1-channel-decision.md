---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-release-channel
version: 0.1.0
---

# D-001 · R1 通道定案（用户裁决）

## 裁决（2026-08-29 用户）

| 项 | 定案 |
|----|------|
| Go 通道 | **公共 proxy + push tag**（授权推远端）；首版 tag = ~~v0.1.0~~ → **`apps/api/v0.1.0`**（实证中发现子目录 tag 约定，见 E-001） |
| npm 通道 | **GitHub Packages**（`https://npm.pkg.github.com` · scope 须 = owner → 包名 `@magicvr/schema-ui-*`）；**发布凭据待用户提供**（GH token · `write:packages`） |
| sumdb 时延 | 知识项：新 tag 收录前可用 `GOSUMDB=off` 临时绕行；收录后重跑默认校验补齐 go.sum（E-001 §知识项 2） |

## 未选方案

- 私有 GOPROXY / Verdaccio / npmjs scope：未选（用户裁决 ①/①）。
- 「不推远端、file 实证到底」：未选（用户选真实推送）。

## 待用户提供（阻塞 S2）

GitHub Packages 发布凭据：`GH_TOKEN`（scope `write:packages`）或等价 npm 认证配置（`//npm.pkg.github.com/:_authToken=...`）。配置完成后告知变量名即可。