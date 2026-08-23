---
title: D-001 · S1 根因冻结：N-001 非路由回归，而是 e2e 挂具 store 隔离失效
status: frozen
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-034-w23-admin-login-home-redirect
version: 0.1.0
---

# D-001 · S1 根因冻结

## 结论（一句话）

N-001（admin 登录后停留 `/` 未跳 `/dashboard`）**不是 home 推导/路由回归**：根因是本机 gitignored 的 `apps/api/configs/.env`（2026-08-21 建立，`DB_DIALECT=postgres`）被 `config.Load` 的 env-file 层读入，把 e2e 挂具本应隔离的临时 SQLite 静默改写为开发者 Postgres —— 全新种子库从未被创建，`admin/admin` 对开发者库 401，登录链在第一步就断开，界面自然停在 `/`。

## 证据链（代码级）

1. **全量 admin e2e 复现（本波 2026-08-23）**：`APP_PROFILE=admin npx playwright test` → 6/10 失败，**全部**为登录门禁失败；首个用例（`force-password-change.spec.ts:21`）在**全新种子**上以 `admin/admin` 提交即得 `alert: invalid username or password`（401 信封 `error.unauthorized`）。单纯「路由不回跳」不可能产生登录 401。

2. **挂具隔离失效**：`playwright.config.ts` 为 API 传入 `CONFIG_FILE=<临时 overlay>`、`DB_PATH=<临时 e2e.db>`、`ADMIN_INITIAL_PASSWORD=admin`，但**未传 `DB_DIALECT`**。`config.Load` 层面顺序（config.go:314-318）：`CONFIG_ENV_FILE`（默认 `configs/.env`，相对 API CWD `apps/api`）→ YAML → 内嵌默认；`.env`「never overrides an already-set process env」——挂具没设 `DB_DIALECT`，`.env` 的 `DB_DIALECT=postgres` 生效。

3. **产物佐证**：两次 e2e 运行后，`%TEMP%\schema-ui-e2e-*` 临时目录**没有任何 SQLite 文件**（`e2e.db` 不存在）；直接复现时同样现象。`DB_PATH` 对 postgres 方言无意义。

4. **凭证级直接复现**：
   - 复现环境 A（同挂具语义，未钉方言）：临时 config overlay + `ADMIN_INITIAL_PASSWORD=admin` → `POST /api/auth/login {admin,admin}` → **401**；
   - 复现环境 B（仅加 `DB_DIALECT=sqlite`）：同参数 → **200**，返回 access/refresh token，`user.mustChangePassword=true`（强制改密屏幕可进入，链路完整）。
   - 修复方向由此唯一确定：让 e2e 挂具显式钉死方言，恢复隔离语义。

5. **W22 基线实验结论失效**：`git stash` 只回滚已跟踪文件，**无法移除 gitignored `configs/.env`**，基线实验全程仍在 postgres 方言下运行，故「stash 后 HEAD 同败」只能证明「同一环境配置下同败」，不能证明「先于 W22 的代码回归」。仓库根 `e2e-baseline.log` 的尾部甚至显示 stash 后 API 因 postgres 迁移台账 49 不匹配**直接启动失败**（`Process from config.webServer was not able to start. Exit code: 1`）——该日志无法佐证「在 :62 处以 `/` 停留失败」的描述。W22 E-006 的「疑似 W14–W21 home 推导/路由漂移」属**未按凭证校验的推测**，本波更正。

## 为何此前各波全量绿未拦截

- Go 单测与 vitest 均为**密闭环境**：不读取 API 进程的 env-file，与 `configs/.env` 无交集 → 全绿不具甄别力。
- 浏览器 e2e 是唯一读取该 env-file 的消费面；`.env` 建于 2026-08-21（W21 窗口），此后本地 e2e 全部在 postgres 方言下运行。W21 关门证据不含 e2e；W22 的补跑（含基线实验）已在该污染下执行，其失败被误读为「登录成功但未跳转」。

## 修复决策

| # | 决策 | 理由 |
|---|------|------|
| D1 | `apps/web/playwright.config.ts` API webServer env 显式钉死 `DB_DIALECT: "sqlite"`（含注释说明 W23 事故） | 恢复挂具隔离承诺：无论本机 `.env`/进程 env 如何，e2e 永远使用临时 SQLite 与种子密码 |
| D2 | `apps/web/e2e/localization.spec.ts` `signInZh` 竞态修复：一次性 `isVisible()` 改为等待式三段流程（先 `waitForURL` home → 再等强制改密屏 → fallback 共享密码），与 `sign-in.ts` 既有模式对齐 | 首登 submit 到强制改密屏渲染是异步往返，`isVisible()` 不等待 → 即使种子正确也可能误入 fallback 分支提交错误密码（当前代码固有竞态，独立于 D1 根因）；等待式消除同类抖动 |
| D3 | 不改任何产品代码（Go/Web src）与协议/manifest 契约语义 | 根因在测试/挂具隔离层；产品登录链实测正确（证据 4 环境 B 全链 200） |
| D4 | 关门审计模式 `self` | 按 00-meta 边界声明：未触及协议/manifest 契约语义，不升级 independent |
| D5 | 残余风险记录：开发者在自有环境跑 e2e 前若修改挂具 env 覆盖（或未来新增方言类 env 键），隔离依赖 `DB_DIALECT` 钉值的存在——已用注释 + 本 D-001 留痕；不做额外机制 | 有界、可复审（触发条件：e2e 再次出现「SQLite 未生成」类现象） |

## 未选方案

- **给 config.Load 增加「e2e 模式」开关**：引入全局行为特异面，超出本波修复范围，且掩盖挂具自身责任。
- **删除/改名本机 `configs/.env`**：那是开发者本地 postgres 工作配置（gitignored、用户所有），与挂具无关；删它治标不治本。
- **在 overlay config.yaml 中写 `db.dialect: sqlite`**：等效于 D1 但语义较弱（YAML 仍可被 env 覆盖），且挂具环境变量已是既有模式（`DB_PATH`/`ADMIN_INITIAL_PASSWORD` 同层），D1 一致性强。