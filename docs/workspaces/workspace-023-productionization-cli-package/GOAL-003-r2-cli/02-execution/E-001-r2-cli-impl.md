---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-cli
version: 0.1.0
---

# E-001 · CLI 实现与 create 验证（2026-08-29）

## 实现（apps/api/cmd/schema-ui · 零新依赖）

- `main.go`：`create`/`add`/`upgrade` 子命令（标准库手写解析——flag 包与位置参数交错解析不可靠，已规避）；`add` = 可用模块清单（kernel.BuiltinModules）+ registry 装配；`upgrade` = `go get @latest` / `pnpm add @latest` + 探针回归（--dry-run 支持）。
- `modules.go`：同模块读 kernel 模块清单。
- `templates/`（go:embed）：go.mod/main.go/README/package.json/npmrc 插值模板 + 静态资产（探针/css/gitignore——**embed 排除点文件**，.gitignore/.npmrc 改名 `.tmpl` 规避）。
- 分发 = `go install github.com/magicvr/schema-ui-core/apps/api/cmd/schema-ui@vX`（随模块 tag）。

## create 验证（双端全绿）

```
schema-ui create demo-admin --module github.com/acme/demo-admin --dir …/demo-admin
→ 11 文件（Go 组合根 + web 骨架 + 探针 + README + .npmrc + .gitignore）

Go:   go mod tidy（proxy 拉 apps/api v0.1.0）→ go run ./cmd/server
      demo-admin kernel=2.0.0 profile=admin dialect=sqlite fresh=true contrib{r=10 p=2 perm=3 nav=1 frag=1}
Web:  pnpm install（registry）→ probe/render/token 三探针 PASS
```

与 golden-field 手工骨架**同构**（模板源自其实证内容）——双轨对照成立（S2）。