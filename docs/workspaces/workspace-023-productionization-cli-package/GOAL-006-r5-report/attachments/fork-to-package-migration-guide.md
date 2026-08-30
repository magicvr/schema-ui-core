# fork → 包 迁移指南（VP-023 R5 · 2026-08-29）

> 面向既有 fork 下游的过渡路径（**指南先行**；工具化迁移 = go 后）。
> 战略语境：fork 与包消费并存（Charter 0.3.0）——本指南给「想换到包消费」的 fork 仓一张地图，不强制迁移。

## 1. 判定：你的 fork 属于哪类

| 类型 | 特征 | 建议 |
|------|------|------|
| A · 纯装配型（不改内核/渲染器） | 只加业务模块 + schema + 配置；未动 kernel/modules/assembly | **直接迁移**（本指南 2–3 步） |
| B · 轻度定制型 | 改了少量 shell/主题/UI（覆盖而非重写） | 迁移 + Token/覆盖面复核 |
| C · 深度定制型（改内核/渲染器主路径） | 修改了 kernel 契约面或 renderer 核心 | **保持 fork**（包形态暂不覆盖；登记录入 go 后服务面） |

## 2. 迁移步骤（A/B 型）

1. **包面打平**：把 fork 仓的 `apps/api`/`apps/web` 源码删除，替换为包依赖：`go.mod` require apps/api @tag（当前 `v0.2.0`）+ `web/package.json` 依赖 `@magicvr/schema-ui-*`（六包）。
2. **组合根重建**：用 `schema-ui create` 生成骨架 → 迁入业务模块（fork 的业务模块 = 包依赖 + 组合根注册）与 schema 文档（协议同 pin，schema 平移即可）。
3. **主题/定制保留**：把 fork 的 `brand.css`/Token 覆盖复制进新骨架的 `web/`（覆盖面纪律：只覆盖 index.css 已声明 Token——`token-check.mjs` 校验）。
4. **数据与迁移**：数据库文件直接可复用（同全局台账/checksum）；PG 路径 = ops-playbook 备份/恢复命令。
5. **升级路径切换**：`git pull upstream` → `schema-ui upgrade`（bump + changelog 迁移说明；零冲突）。

## 3. 验证（迁移后必跑）

```bash
schema-ui upgrade --dry-run   # 版本面预检
go run ./cmd/server -dialect sqlite -dsn ./data.db   # 迁移台账核对（旧库升级）
cd web && node probe.mjs && node probe-render.mjs && node probe-six.mjs && node token-check.mjs
```

## 4. 残余（go 后）

- C 类深度定制 fork 的包化承载面（kernel 契约扩展通道）——建议面 = assembly 扩展 + 六包 external 化后评估。
- schema 文档目录化迁移工具（当前 = 手工平移 + 探针验证）。