---
id: E-001-dual-path-evidence
title: R5 双路径证据采集（本地回归 + MinIO live + readyz 阴性对照）
status: recorded
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-006-dual-path-evidence
version: 0.1.0
---

# E-001 · R5 双路径证据

## 路径一：本地盘默认回归

全量 `go test ./... ` exit 0（含 handler 直盘兼容用例、file-library/data-transfer 端到端、composition 单实例装配）——本地缺省路径持续可开发/快测（VP 判据 3）。

## 路径二：S3 兼容生产向验收（真实 MinIO 容器）

环境：docker minio/minio:latest，端点 http://127.0.0.1:9000，桶 schema-ui-r5（mc 建桶；适配器按合同不建桶）。

1. **配置接入**：真实 config.yaml `driver: s3` + endpoint/bucket/凭证经 Load fail-closed 校验通过并启动服务进程（首次因 db.path 反斜杠 YAML 转义失败——属测试配置书写问题非产品缺陷，修正后通过）。
2. **读写删除**：`TestS3LiveRoundTrip -count=1 -v` → **PASS** (0.04s)——Ping(HeadBucket)/Put/Get/Stat/List/Delete 全链路对真实 S3 兼容后端。
3. **就绪探针**：driver=s3 启动真实进程（go run ./cmd/server），GET /readyz：
   - 后端在线 → **200 {"status":"ok"}**
   - docker stop MinIO → **503**（阴性对照：探针真实在链路）
   - docker start MinIO → **200 恢复**

VP-014 判据 4 的"配置接入 / 读写删除 / 就绪探针"三项全部可核对（要求为至少其一）。容器已清理，不留常驻资源。

## 结论

R5 双路径证据齐备。待关门审计后结项 Root（判据逐条核对见 A-001 self 关门审计）。
