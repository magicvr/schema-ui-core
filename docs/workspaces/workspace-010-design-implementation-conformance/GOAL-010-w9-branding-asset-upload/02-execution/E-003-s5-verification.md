---
id: GOAL-010-w9-branding-asset-upload
doc: execution
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-003 · S5 验证（单测 + 全量回归 + 活栈点验）

2026-08-15 完成 S5 验证：

## 单元与集成测试

- `apps/api/internal/handler/branding_assets_test.go`（8 例）：上传+公开 GET（1024→512 JPEG、nosniff/immutable 头、可解码）、favicon 64px + 透明 PNG、SVG/空文件/坏 kind/匿名 401/editor 403/文本 415/超限 413（1 KiB 策略）、GET 404 与路径穿越、替换/清空/重置清理、启动 GC。
- `TestBrandingConfig`（4 例）：默认值 / YAML 层 / 越界质量回退 / env 覆盖。
- `TestErrorCodeContractPinnedSet` + `TestErrorCatalogCoversFrozenCodesExceptInternal`：新错误码已登记契约与双语目录。
- provider 测试：路由数 5→7（含 POST/GET assets）；settings schema 键完整性（schema-keys.structural）、s5 分母渲染、startup-config 渲染（设置页四分类含上传字段）通过。
- **Go 全量 `go test ./...`：PASS（exit 0）**；**Web vitest 全量 967/967 PASS**；gofmt/vet 干净。
- 注：`tsc -b` 在 HEAD 上即存在 6 处既有类型错误（git stash 验证与本波无关），vitest 构建链不受影响。

## 活栈点验（HTTP 级）

- 本地启动真实 API（admin profile + dev session）：
  - 上传 1024×1024 PNG → 200，返回 `/api/branding/assets/{id}`，重编码为 512×512 image/jpeg（7609 B）；
  - 公开 GET → 200 + `Cache-Control: public, max-age=31536000, immutable` + nosniff + CSP sandbox；
  - PATCH 设置 logoUrl → 持久化；`/api/branding` 投影即时反映；
  - 替换（A→B）：旧资产 A 即时删除；清空字段：B 删除；恢复默认：目录清零，旧 URL 404。
