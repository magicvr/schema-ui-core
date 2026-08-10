---
id: GOAL-001-full-protocol-contract-v2-7-0
doc: execution-entry
record_id: E-004
status: recorded
parent: null
created: 2026-08-08
updated: 2026-08-08
version: 0.1.0
---

# E-004 · S5 回归、文档与验证计划执行

## 2026-08-08 · S5 回归 + 验证计划（Verification plan 全部观测）

### 已发生事实

1. **发现路径文档更新**（覆盖权威 = `I-PROTO-FULL-001`）：
   - `docs/architecture/overview.md`：新增 VP-006 active + workspace-005 + 覆盖权威指针；移除「无 active 交付 VP」过时表述。
   - `QUICKSTART.md` §5：新增「协议覆盖权威」段（I-PROTO-FULL-001 链接 + 声明背书纪律）。
   - `docs/README.md`：推荐阅读后新增覆盖权威注。
   - `docs/vision/protocol-inventory-v2.7.0.md` §4：「尚未冻结」→ 已冻结（`I-PROTO-FULL-001`）。
   - grep 全文：**无**未背书的「已完整支持 v2.7.0」声明（仅覆盖表/决策/审计中的纪律性表述）。
2. **回归重跑**（捕获到 `{SCRATCH}`）：
   - `go test ./...`（apps/api）：**全包 ok**（go-test.log，exit 0，无 FAIL）。
   - `npm run test`（apps/web vitest）：**25 文件 / 569 测试全绿**（vitest.log）。
   - `npm run build`（tsc -b && vite build）：通过（dist 489.56 kB JS）。
   - `npm run test:e2e`（playwright，WEB_PORT=9999，mvp profile）：**2/2 通过**（schema-crud 8.6s + shell 0.9s）。
3. **验证计划 step 3 · API 真实入口启动（两次）**：`go run ./cmd/server`（admin profile，临时 DB）启动两次（api-launch-1.log / api-launch-2.log）：
   - `/healthz`、`/readyz` → 200 且 body 含 `"status":"ok"`（两次一致）。
   - 登录 `/api/auth/login` → accessToken；`/api/accounts/me` → `user.name=Admin`、`roles=admin,editor`、features 投影（两次一致）。
   - `/.well-known/schema-ui/app-manifest.json` → 12 页含 admin-list-batch/data-display/form-with-upload。
   - `/api/schema/admin-list-batch` → 200；`POST /api/users/batch-delete`（未知 id）→ 404（权限门 + 实体边界真实生效）。
4. **验证计划 step 4 · Web 构建 + 真实页面加载（headless）**：`vite preview`（构建产物 + preview 代理）+ 真实 API；Chromium 1280×800：
   - 登录 → 导航 `/data-display`（statCard「Total roles」+ chart SVG 448×160 绘制）→ `/admin-list-batch`（勾选 → 批量按钮禁用解除 → Confirm → 「Items deleted」→ 清选）→ `/form-with-upload`（真实文件选择 → multipart 上传 → 字段值 `Value: /api/files/…`）。
   - **零 console/page 错误**（CONSOLE_ERRORS []）；bodyScrollWidth=1280=viewport（无横向溢出）；截图 `{SCRATCH}/web-launch.png`；日志 `{SCRATCH}/web-launch.log`。
   - 过程中发现并修复两个真实缺陷：① 上传真实传输未携带文件字节（`UploadableFile.blob` 修复）；② `vite preview` 无代理（`preview.proxy` 对齐 dev）。均为产品代码修复，有对应测试。
5. **环境限制记录**：`scripts/smoke.sh` 为 bash 脚本（Windows 不可运行）——其等价断言（healthz/readyz、登录、身份、SPA 挂载）已由本 E-004 step 3/4 的真实端点 + playwright e2e 覆盖；close-out 审计记录此限制。
6. 上传缺陷修复后 vitest 569 全绿 + headless 复验通过（UPLOAD_OK）。

### 证据

| 主张 | 路径 |
|------|------|
| go 回归 | `{SCRATCH}/go-test.log`（exit 0） |
| vitest 回归 | `{SCRATCH}/vitest.log`（569/569） |
| API 启动×2 | `{SCRATCH}/api-launch-1.log`、`api-launch-2.log`（healthz/readyz/login/me/manifest/batch 全观测） |
| build + e2e | `npm run build` 通过；`npm run test:e2e` 2/2 |
| headless 加载 | `{SCRATCH}/web-launch.log`（BATCH_OK/UPLOAD_OK/METRICS/CONSOLE_ERRORS []/HEADLESS_VERIFIED）+ `{SCRATCH}/web-launch.png` |
| 文档指针 | overview.md / QUICKSTART.md / docs-README / protocol-inventory（grep 复核） |
| 无未背书声明 | 全仓 grep「已完整支持/完整支持 v2.7.0」复核（仅纪律性表述） |
