---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-002-r1-serve-shell
version: 0.1.0
---

# E-001 · 目标建立（2026-08-29）

1. **立项**：承接 Root 纲领 R1（serve 壳闭环 · VP-024 判据 #1 · go 后清单 ①）；goal-tree 同步（Root 0/7，新增本目标 active 0/5）。
2. **现场核验（设计依据）**：
   - 生成骨架 `main.go.tmpl` 现状 = 装配冒烟（OpenStore → auth/repo/ops/mailer → users provider → RegisterContributions → 打印贡献计数退出），**无 HTTP 服务、无 config 装载、无停机接线**；
   - 主仓服务面 = `internal/composition`（Fx 全量模块装配 + `newMux` 中央面 + `registerLifecycle` 全序停机）+ `internal/server.New`（timeouts + requestid + nosniff/CORS）+ `internal/config.Load`（YAML `${VAR}`/`${VAR:-def}` 插值 fail-closed + CONFIG_FILE/env 层）；
   - RT-D02 合同（workspace-021 D-002 v0.1.1）：停机全序 9 步、drain 预算 `http.shutdown_timeout`（默认 10s）、退出码 0/1、迁移无运行时窗口、双方言同 Close 路径；
   - 中央面注册清单：`RegisterWithMFAProbes`（auth+health/ready+探针）、`RegisterSchemas`、`RegisterManifest`/`RegisterBootstrapWithAvailability`、`RegisterUpload` 等。
3. **信息门禁**：I-001（serve 面构成 + 模板形态）required · 最晚 S1 → 用户裁决中（P-004）。