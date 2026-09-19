---
title: R3 前端批量触发点侦察（自定义组件契约 · 选择可达性 · 按钮落点）
status: draft
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.1.0
---

# R3 前端批量触发点侦察

> **性质**：只读侦察（READ-ONLY）。本文件是本次任务**唯一**的写入；未修改任何源码 / schema / 测试。
>
> **落点说明**：任务书指定的 `GOAL-004-r3-async-batch-承接/attachments/` 在本仓库**不存在**（workspace-038 下实际目标文件夹名为 `GOAL-004-r3-async-batch-operation`，其 `attachments/` 已存在且只有 `.gitkeep`）。按任务书的回退规则，本报告写入工作区根。若需归档进目标目录，移动即可（内容不含相对路径依赖）。
>
> **冻结前提（本次不复核、不再讨论）**：异步契约 = 新增 `POST /api/jobs/batch-export` → 202 + job 投影；ADR-0022 的同步 `batchMapping`/`batch-delete` 路径保持字节级不变且**不得复用**于异步提交；页面 schema **不得**声明 `actions.batch.request`（守卫 marker 为 `/batch-delete/`）；触发器 = 经 `registerCustomComponent` 注册的自定义组件（先例 `monitoring-auto-refresh`）。

## 0. 结论速览

| # | 问题 | 结论 | 关键证据 |
|---|---|---|---|
| 1 | 当前表格选择是否在自定义组件的 `context` 上 | **否**。`context` 只承载 host/导航数据（`user`/`features`/`route`），页面 CRUD 状态走 React Context，不在该 record 里 | `apps/web/src/app/AuthGate.tsx:83-86`、`apps/web/src/app/App.tsx:786-793`、`apps/web/src/renderer/render.tsx:122`、`:3240` |
| 2 | 选择是否**可达** | **是**——通过 `useSchemaCrud()?.selection("<tableId>")`（`SchemaCrudValue.selection`，公开导出） | `apps/web/src/renderer/render.tsx:271`、`:927-936`、`:1501`、`apps/web/src/renderer/index.ts:10,15` |
| 3 | 该 seam 今天是否被任何自定义组件用过 | **没有**。全仓库只有 `schema-table.tsx` 调 `crud.selection(...)`/`setSelection(...)` | `apps/web/src/renderer/schema-table.tsx:993,1020,1050`（grep 无其它命中） |
| 4 | 选择生效的前置条件 | users 表节点必须先声明 `props.selection.mode === "multiple"`，否则 `crud.setSelection` 永不被调用，`crud.selection()` 恒为 `undefined` | `apps/web/src/renderer/schema-table.tsx:932-935`、`:993` |
| 5 | 自定义节点能否作为 table 的 toolbar 项 | **不能**。toolbar 渲染只有 `<button>` 分支，无 custom 分支 | `apps/web/src/renderer/schema-table.tsx:1326-1367` |
| 6 | 推荐落点 | users 页 `body`（section `users`）里、table `users-table` 的**兄弟节点**（table 之前或之后皆可），`props.targetTable="users-table"`，先例 `activity-export` | `apps/api/modules/activity/schema/activity.json:124-131`、`apps/api/modules/datapermission/schema/data-permission.json:149-154`、`apps/web/src/components/activity-export.tsx:19-29` |
| 7 | 冻结设计会否触发 capability 守卫 | **否**（既不声明 `actions.batch.request` 也不声明 `table.selection`）。**反例警告**：若顺手给 users.json 加 `table.selection`，守卫会**失败**（其 marker 是 `"requiresSelection"`，users 页没有该键） | `apps/web/src/protocol/capability-declaration.guard.test.ts:51,57,106-126` |
| 8 | 同步 batch 路径的回归锚点 | 三组：端到端集成测试 2 条 + schema 能力守卫（含 `/batch-delete/` marker）+ 上游构造用例 11 条 batchRequest 全量执行 | §6 |

---

## 1. 自定义组件契约（`CustomComponentProps.context` 究竟有什么）

### 1.1 注册表与 props 形状

`apps/web/src/renderer/custom-components.ts:11-31`：

```ts
export interface CustomComponentProps {
  node: RenderCustomNode;
  context: Record<string, unknown>;
  /** Children declared under the custom node (rarely used). */
  children?: RenderNode[];
}

const registry = new Map<string, ComponentType<CustomComponentProps>>();

/** Registers a custom node component (idempotent; later registrations win). */
export function registerCustomComponent(
  key: string,
  component: ComponentType<CustomComponentProps>,
): void {
  registry.set(key, component);
}
```

节点类型允许 `props` 与 `children`：`apps/web/src/renderer/render.types.ts:243-250`

```ts
export interface RenderCustomNode {
  type: "custom";
  id?: string;
  /** Registered component key (GOAL-018: renderer custom-component registry). */
  component: string;
  props?: Record<string, unknown>;
  children?: RenderNode[];
}
```

### 1.2 renderer 如何调用（含未注册兜底）

`apps/web/src/renderer/render.tsx:3219-3241`：

```tsx
    case "custom": {
      // GOAL-018: custom nodes dispatch to the module-level registry.
      // C-010 (GOAL-041 S2): an unregistered component renders an obvious
      // placeholder ... never a silent blank ...
      const Custom = getCustomComponent(node.component);
      if (Custom === null) {
        console.error(
          `[schema-ui] unknown custom component "${node.component}"` + ...
        );
        return ( <div role="alert" ...>unknown custom component: <code>{node.component}</code></div> );
      }
      return <Custom node={node} context={context} children={node.children} />;
    }
```

**注意**：`children={node.children}` 只是把原始节点数组**传**给组件；renderer **不提供**公开的 `dispatchNode` 导出（`apps/web/src/renderer/index.ts:8-35` 只导出 `RenderPage`/`useSchemaCrud`/`SchemaCrudContext`/`registerCustomComponent`/类型/i18n），所以自定义组件**无法**自行渲染 schema 子节点。

### 1.3 `context` 的实际 key（穷举）

`context` 是宿主传入的**导航/host context**，一路透传：`App.tsx:786-793` → `RenderPage` → `render.tsx:3334`（body 节点）/`:3351`（modal 内容）。

```tsx
// apps/web/src/app/App.tsx:786-793
    <RenderPage
      document={state.document as RenderPageDocument}
      context={
        {
          ...(context as Record<string, unknown>),
          route: { params, query },
        } as Record<string, unknown>
      }
```

```ts
// apps/web/src/app/AuthGate.tsx:83-86 —— 生产宿主的 context 全部内容
  const context: NavigationContext = {
    user: user === null ? undefined : (user as unknown as Record<string, unknown>),
    features: session?.features,
  };
```

| key | 来源 | 何时存在 |
|---|---|---|
| `user` | `AuthGate.tsx:84`（登录用户；`$context.user.permissions/roles` 表达式即读它） | 生产恒定 |
| `features` | `AuthGate.tsx:85`（功能开关 map） | 生产恒定 |
| `route.params` / `route.query` | `App.tsx:791`（`render.tsx:866-881` 会做类型清洗） | 生产恒定 |
| `pageId` | 仅当宿主显式传入；`render.tsx:841-846` 用它作为 `meta.pageId` 缺失时的回退 | 生产**不**传（生产用 `meta.pageId`） |
| `modalRow` | `render.tsx:3351`：`context: { ...context, modalRow: crud.activeModal?.row ?? null }` —— **只**对 modal 的 `content` 生效 | 仅 modal 内 |
| 其它 | 测试宿主可任传（`denominator-render.test.tsx:119` 传 `context={{}}`） | 不确定 |

**`context` 里没有**：`selection`、`selections`、`tableId`、`selectedRow`、`crud`、`pageId`（生产）。

### 1.4 选择状态实际住在哪里（关键追踪）

```ts
// apps/web/src/renderer/render.tsx:818  —— provider 内部 state
  const [selections, setSelections] = useState<Record<string, unknown[]>>({});
```

```ts
// apps/web/src/renderer/render.tsx:927-957
  const selection = useCallback(
    (id: string): TableSelection | undefined => {
      const keys = selections[id];
      if (keys === undefined) {
        return undefined;
      }
      return { keys, count: keys.length };
    },
    [selections],
  );
  const setSelection = useCallback((id: string, keys: unknown[]) => {
    setSelections((prev) => ({ ...prev, [id]: normalizeSelection(keys).keys }));
  }, []);
  const clearSelection = useCallback((id: string) => { ... }, []);
  // ADR-0022 D2: any data reload success clears every table selection. ...
  const reloadList = useCallback(() => {
    setSelections({});
    listInFlight.current = new Map();
    setReloadToken((token) => token + 1);
  }, []);
```

这段 state **只**通过 React Context 暴露，并**不在** `context` record 上：

```ts
// apps/web/src/renderer/render.tsx:295-301
/** Page CRUD context consumed by custom components (rendered nodes). */
export const SchemaCrudContext = createContext<SchemaCrudValue | null>(null);

/** Reads the page-level Schema CRUD provider (null when rendered bare). */
export function useSchemaCrud(): SchemaCrudValue | null {
  return useContext(SchemaCrudContext);
}
```

```ts
// apps/web/src/renderer/render.tsx:270-273（interface SchemaCrudValue 片段）
  /** Per-table selection state (keys + count; normalized by the table). */
  selection: (tableId: string) => TableSelection | undefined;
  setSelection: (tableId: string, keys: unknown[]) => void;
  clearSelection: (tableId: string) => void;
```

`selection` 被放进 provider value 的依赖列表，因此**选择变化会让所有 context 消费者重渲染**（自定义组件能拿到最新值）：

```ts
// apps/web/src/renderer/render.tsx:1490-1527（节选）+ :1528-1562（deps 数组含 selection）
  const value = useMemo<SchemaCrudValue>(
    () => ({
      ...
      selection,
      setSelection,
      clearSelection,
      ...
    }),
    [ ..., selection, setSelection, clearSelection, ... ],
  );
  return <SchemaCrudContext.Provider value={value}>{children}</SchemaCrudContext.Provider>;  // :1564
```

Provider 包住整页（自定义节点必然在其内，故 `useSchemaCrud()` 非 null）：

```tsx
// apps/web/src/renderer/render.tsx:3380-3396
    <SchemaCrudProvider document={document} context={context} initialFetcher={dataFetcher} onNavigate={onNavigate}>
      <RenderPageSurface ... />
    </SchemaCrudProvider>
```

**`SchemaCrudValue` 可用 key 全表**（`render.tsx:204-285`）：`userId`、`pageId`、`tableFilterFields`、`notifyFeedback`、`registerDirtySource`、`selectedRow`、`selectRow`、`tableQuery`、`setTableQuery`、`selection`、`setSelection`、`clearSelection`、`reloadToken`、`reloadList`、`fetchList`、`refreshList`、`listRefreshToken`、`activeModal`、`modalRow`、`openModal`、`closeModal`、`pendingConfirm`、`requestConfirm`、`resolveConfirm`、`feedback`、`registerFetcher`、`fetcher`、`runRowAction`、`invokeAction`、`invokeBatchAction`、`uploadFiles`、`submitForm`、`searchFormSubmit`、`effectivePermission`、`route`。

`TableSelection` 形状与归一化规则：`render.tsx:184-188`（`{keys: unknown[]; count: number}`）+ `apps/web/src/protocol/conformance/request-construction.ts:623-640`（**标量键、按序去重、count = keys.length**，字符串与数字不互转）。

### 1.5 一句话回答

> **选择不在 `context` 上；但它在页面级 `SchemaCrudContext` 上，自定义组件可通过 `useSchemaCrud()?.selection("<tableId>")` 读取（这就是唯一可用的 seam），并且选择变化会触发该组件重渲染。**

`context` 无法提供 table id，**组件必须自带**（先例：`activity-export` 用 `node.props.targetTable`，见 §2.4）。

---

## 2. `monitoring-auto-refresh` 先例全解（我们要照抄的模板）

文件：`apps/web/src/components/monitoring-auto-refresh.tsx`（全文 62 行）。

### 2.1 如何注册

- 模块作用域**副作用注册**，最后一行：`monitoring-auto-refresh.tsx:62` → `registerCustomComponent("monitoring-auto-refresh", MonitoringAutoRefresh);`
- 注册表幂等、后注册覆盖：`custom-components.ts:20-26`。
- 生产侧必须在入口做副作用 import：`apps/web/src/main.tsx:17` → `import "@/components/monitoring-auto-refresh";`
- 测试侧也必须 import（否则守卫转红）：`apps/web/src/renderer/custom-components.schema.test.ts:21`、`apps/web/src/renderer/denominator-render.test.tsx:43`、`apps/web/src/i18n/s5-denominator-render.test.tsx:39`。

### 2.2 消费的 `context` key

**一个都没有。** 签名是 `export function MonitoringAutoRefresh(_props: CustomComponentProps)`（`:27`），函数体只用 `useTranslate()` 与 `useSchemaCrud()`（`:28-29`）。→ 先例证明：**自定义组件不需要 `context`**；需要页面状态时走 `useSchemaCrud()`。

### 2.3 如何调用 crud 与渲染控件

```tsx
// apps/web/src/components/monitoring-auto-refresh.tsx:17-59（节选）
const OPTIONS = [
  { value: 0, labelKey: "monitoringRefresh.off" },
  { value: 5000, labelKey: "monitoringRefresh.5s" },
  { value: 10000, labelKey: "monitoringRefresh.10s" },
  { value: 30000, labelKey: "monitoringRefresh.30s" },
];
/** The display dataSource refreshed by the polling tick (system-monitoring.json). */
const STATUS_SOURCE = "/api/system-monitoring/status";

export function MonitoringAutoRefresh(_props: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const [intervalMs, setIntervalMs] = useState(0);

  useEffect(() => {
    if (intervalMs <= 0) { return; }
    const id = window.setInterval(() => {
      crud?.refreshList(STATUS_SOURCE);      // ← 定向刷新（不是 reloadList）
    }, intervalMs);
    return () => window.clearInterval(id);
  }, [intervalMs, crud]);

  return (
    <div className="flex items-center gap-2 text-sm">
      <span className="text-muted-foreground">{t("monitoringRefresh.label")}</span>
      <select aria-label={t("monitoringRefresh.label")} value={String(intervalMs)}
        onChange={(event) => setIntervalMs(Number(event.target.value))}
        className="h-8 rounded-md border border-input bg-background px-2 text-sm"
        data-monitoring-refresh>
        {OPTIONS.map((option) => (
          <option key={option.value} value={String(option.value)}>{t(option.labelKey)}</option>
        ))}
      </select>
    </div>
  );
}
```

要点：
- **间隔选项**：`0 / 5000 / 10000 / 30000` ms，`0 = 关闭`（默认值 0，`:30`）；i18n key `monitoringRefresh.label|off|5s|10s|30s`（`apps/web/src/i18n/messages/en-US.json:83-87`、`zh-CN.json:83-87`）。
- **调用的是 `crud?.refreshList(url)`**（`render.tsx:245`、`:970-977`）：按 dataSource 的 display 查询做**定向**刷新，`listInFlight.delete(key)` 后 bump token；语义上是"刷新某个展示面"，**不会**清空选择。
- **渲染**：一个 flex 容器 + `<span>` 标签 + 原生 `<select>`，带 `data-monitoring-refresh` 测试钩子（无 shadcn 依赖，样式用 Tailwind 原子类）。
- 组件是**自包含**的：不读 `node.props`（对照 `activity-export` 读 `node.props.targetTable`）。

### 2.4 在页面 schema 里的声明

```json
// apps/api/modules/systemmonitoring/schema/system-monitoring.json:12-27
  "body": {
    "type": "section",
    "id": "system-monitoring",
    "children": [
      {
        "type": "custom",
        "id": "monitoring-refresh-control",
        "component": "monitoring-auto-refresh"
      },
      {
        "type": "text",
        "props": { "text": "Read-only system status and recent events.", "textKey": "schema.systemMonitoring.text.intro" }
      },
```

即：**`body`（根 section）的 `children[0]`，无 `props`**，与后续 `statCard`/`table` 平级。该页 `requiredCapabilities` 只有 `app.manifest`/`app.navigation`/`table.sort`（`:6-10`）——**自定义节点不引入任何能力声明**（守卫只扫 `component` 与注册表的一致性，见 §3.3）。

### 2.5 这个先例的测试形态（照抄测试时用）

- 测试替身组件直接调 crud：`apps/web/src/renderer/render.test.tsx:67-81`（`PerfRefreshController` + `registerCustomComponent("perf-refresh-controller-test", ...)`），其用法断言在 `:455-500`（定向刷新只打一个 URL）与 `:505+`（不与 in-flight 请求合流）。
- 未注册组件的兜底断言：`render.test.tsx:1576-1600`（`role="alert"` + 含 component key + `console.error`）。

---

## 3. 页面如何声明自定义节点（真实样例盘点）

### 3.1 两类 `"type":"custom"` 必须区分（易混）

| 类别 | 形状 | 分发路径 | 例子 |
|---|---|---|---|
| **渲染层自定义节点**（本任务相关） | `{"type":"custom","component":"<key>"}` | `render.tsx:3219-3241` → 注册表 | 见 §3.2 |
| 动作层自定义 handler | `{"type":"custom","handler":"<name>"}`（在 `actions.*` 里） | `render.tsx:347-363` `CUSTOM_HANDLER_URLS` 白名单 | `users.json:115-118`（`export.users`）、`roles.json:53-56`、`filelibrary/schema/file-library.json:62-73` |

`custom-components.schema.test.ts:5-7` 明确写了这一区分（"作用域仅为渲染层 custom 节点"）。**用 `"type":"custom"` 做 grep 会把两者混在一起**（我第一遍就踩到）。

### 3.2 全部渲染层自定义节点（14 处）

| # | 模块 / schema | pageId | component key | 落点 |
|---|---|---|---|---|
| 1 | `systemmonitoring/schema/system-monitoring.json:16-20` | system-monitoring | `monitoring-auto-refresh` | body 根 section 的 `children[0]`，无 props |
| 2 | `activity/schema/activity.json:124-131` | activity | `activity-export` | 根 section 子节点：search form 之后、table `operations-table` **之前**；`props.targetTable` |
| 3 | `datapermission/schema/data-permission.json:149-154` | data-permission | `data-permission-scopes` | 根 section 子节点，**在 table 之后**（table 的 toolbar 收尾于 `:147`） |
| 4 | `notifications/schema/notifications.json:61-68` | notifications | `notification-center` | 根 section **最后一个**子节点；`props.targetTable` |
| 5 | `users/schema/users-invites.json:39-44` | users-invites | `invite-issue-card` | 根 section `children[0]`，无 props |
| 6 | `users/schema/users-invites.json:26-33` | users-invites | `invite-resend-dialog` | **modal 动作的 `content`**（`resendInvite`）——modal 内 `context.modalRow` 可用 |
| 7 | `account/schema/account.json:115-119` | account | `email-identity` | `section#tab-profile` 子节点（tabs 内 section） |
| 8 | `account/schema/account.json:191-195` | account | `mfa-manager` | `section#tab-security` 子节点 |
| 9 | `account/schema/account.json:206-210` | account | `account-session-toolbar` | `section#tab-sessions` 子节点，**在 `sessions-table` 之前**（`:211-214`） |
| 10 | `settings/schema/settings.json:466-471` | settings | `password-policy-tab` | `section#tab-password-policy` 的**唯一**子节点（tabs → section → custom） |
| 11 | `settings/schema/mail.json:15-20` | mail | `mail-admin-tab` | 根 section 的**唯一**子节点（整页就一个组件） |
| 12 | `channel/telegram/schema/telegram-settings.json:22-27` | telegram-settings | `telegram-admin-tab` | 根 section 的 `children[0]` |
| 13 | `channel/telegram/schema/telegram-operator.json:11-18` | telegram-operator | `telegram-admin-tab`（`props.surface="operator"`） | **`body` 直接就是 custom 节点**（无 section 包裹） |
| 14 | `wallet/schema/my-wallet.json:59-63` | my-wallet | `wallet-ensure` | 根 section `children[0]`，无 props |

**位置规律**：custom 节点**总是** section（或 tabs 内 section、或 body、或 modal content）的**子节点**，与其它节点平级；**从未**出现在 `table` 的 `toolbar` 或 `props` 里。

### 3.3 守卫与注册表一致性（新增组件必改点）

`apps/web/src/renderer/custom-components.schema.test.ts:97-112`：递归扫 `apps/api/modules/**/schema/*.json`，把 `{type:"custom", component:"…"}` 全部收集，断言 `getCustomComponent(key) !== null`。
**副作用 import 清单是硬编码的**（`:18-37`）——新增组件必须同时加进 `main.tsx` 与 `custom-components.schema.test.ts`（历史教训写在 `:28-30`：workspace-018 R3 漏了 import，守卫在 HEAD 直接红）。
另有 3 个"全量页面渲染"测试会因未注册组件而红：`denominator-render.test.tsx:156`（`expect(container.textContent).not.toContain("unknown custom component")`，36 页全量）、`s5-denominator-render.test.tsx`、`render.test.tsx:1576+`。

### 3.4 toolbar 能不能放自定义节点？——**不能**

`apps/web/src/renderer/schema-table.tsx:1316-1369`（toolbar 全量渲染逻辑）：`toolbar.map((trigger) => { ... return (<button ... />) })`，**只有按钮**；`trigger` 支持 `key/label/labelKey/actionRef/permissions/permissionIntent/batchMapping/requiresSelection/confirm/confirmKey`（见 `:1327-1338` 与 `users.json:429-464` 的实际用法），**没有任何 `component`/custom 分支**。

同时：
- **custom 节点不能嵌进 table 节点内部**：`schema-table.tsx` 全文**没有任何 `children` 引用**（grep `children` 零命中），`SchemaTable` 只读 `props.columns/actions/toolbar/filters/...`（`:530-542`）。
- table 的 `props` 也没有"扩展插槽"。

→ 结论：想让触发器出现在 toolbar 那一行，**必须改 `schema-table.tsx`（新增 custom 分支）与协议 toolbar item 形状**，属于突破冻结范围的改动。

---

## 4. 选择可达性：选项逐条核实（诚实版）

### (a) `context` 已经带选择？——**否**
`context` 的所有 key 见 §1.3；`render.tsx:3351` 唯一注入的页面态是 `modalRow`，且只对 modal 内容。`selections` 从未进入 `context`。

### (b) `crud` / 页面上下文有读取选择的方法？——**是，且这是今天唯一可用 seam**

```ts
// apps/web/src/renderer/render.tsx:271
  selection: (tableId: string) => TableSelection | undefined;
```

- 读取：`const sel = useSchemaCrud()?.selection("users-table")` → `{keys, count} | undefined`（`render.tsx:927-936`）。
- 反转依赖：`selection` 的 `useCallback` deps 是 `[selections]`（`:935`），value memo deps 含 `selection`（`:1537`）→ **勾选/清空会重渲染组件**。
- 写入（本任务不需要）：`setSelection`/`clearSelection`（`:272-273`、`:937-949`）。
- 公开面：`apps/web/src/renderer/index.ts:10`（`useSchemaCrud`）、`:15`（`type TableSelection`）。
- 现网用法只有 `schema-table.tsx:993`（`const currentSelection = selectionEnabled ? crud?.selection(tableId) : undefined;`）——**因此这条 seam 已在生产路径被验证，只是还没有自定义组件用过**。
- **前置条件（硬）**：`schema-table.tsx:932-935` 才决定是否渲染多选列并回调 `setSelection`：

```ts
  const selectionEnabled =
    typeof node.props?.selection === "object" &&
    node.props?.selection !== null &&
    (node.props.selection as Record<string, unknown>).mode === "multiple";
```

  users 表今天**没有** `props.selection`（§7）→ 必须先补 `props.selection: {"mode":"multiple"}`，否则 `crud.selection("users-table")` 永远 `undefined`。
- **选择会被清空的时机**（设计触发器时必须知道）：任何 reload 成功（`render.tsx:950-957`，ADR-0022 D2，含 `reloadList()`）、以及筛选/分页/排序/搜索变化（`schema-table.tsx:937-943`，`clearSelection`）。→ 异步 202 提交**不能**顺手 `reloadList()`（否则用户刚选的行立刻被清空，且与"提交后仍在跑"的体验不符）；这与冻结设计"不复用同步路径"一致。

### (c) 组件自己渲染一张表？——**理论可行，实践不可取**
`SchemaTable` 是具名导出（`schema-table.tsx:530` `export function SchemaTable({ node, fetcher, pageTitle })`），组件确实可以自造一个 table node 并渲染它。但那会是**第二张表**：它的选择与页面 `users-table` 的选择是两套 state（虽然都落在同一个 `selections[id]` map 里，只要 id 不同就互不相干），并且会重复请求 `/api/users`。→ **不解决"读取另一张表的选择"**，只是把问题变成"再造一张表"。**不推荐**。

### (d) 其它 seam？

| 候选 | 今天是否可用 | 证据 |
|---|---|---|
| `crud.fetcher` / `crud.notifyFeedback` | ✅ 可用（异步提交要用） | `render.tsx:256-258`、`:211-212`；`activity-export.tsx:21-22`（`crud?.fetcher ?? globalThis.fetch`） |
| `crud.invokeAction(item, row)` 派发页面动作 | ⚠️ 可用但**不满足**异步需求：`type:"request"` 动作走 `runRowAction` → 成功即 `reloadList()` 且**丢弃响应体**；`type:"custom"` 动作只支持 `CUSTOM_HANDLER_URLS` 的 **GET 下载**（无 POST/202/body） | `render.tsx:347-363`、`:352-353`；同步 batch 的等价"丢弃"证据见 §6 |
| 页面 header actions portal（`data-page-list-actions-host`） | ❌ 单所有者：`claim(ownerId)` 在已被占用时返回 false；`SchemaTable` 已用 `tableId` claim | `apps/web/src/renderer/list-surface.tsx:23-30`；`schema-table.tsx:560-575` |
| `children` 插槽（自定义组件渲染 schema 子节点） | ❌ 无公开 dispatch 导出（§1.2） | `renderer/index.ts:8-35` |
| 模块级共享 store（"总选择表"） | ❌ 不存在；注册表只按 component key 存组件 | `custom-components.ts:18` |
| `window`/URL 中转选择 | ❌ 无任何实现痕迹 | 全仓 grep `selection` 仅命中 renderer/protocol |

**结论**：**(b) 是唯一现成 seam**；其余要么不存在，要么需要改 renderer/schema-table（超出冻结范围）。

---

## 5. 按钮落点：逐选项的技术判定与建议

### 5.1 候选与支持状态

| 选项 | renderer/table 今天是否支持 | 需要改什么 |
|---|---|---|
| **A. custom 节点作为 section 内 table 的兄弟节点** | ✅ 支持（先例 §3.2 #2 `activity-export`、#3 `data-permission-scopes`、#9 `account-session-toolbar`） | schema 加节点 +（必须）给 users-table 加 `props.selection.mode=multiple` + 新组件文件 + 注册/import/守卫/i18n |
| B. custom 节点放进 `section` 但更深一层（tabs→section→custom） | ✅ 支持（先例 #10） | 同 A；users 页没有 tabs，无必要 |
| C. custom 节点作为 `body` 根（无 section） | ✅ 技术上支持（先例 #13） | 会把整页替换掉 —— **不适用** |
| D. custom 节点作为 table 的 **toolbar 项** | ❌ **不支持**：`schema-table.tsx:1326-1367` 只有 `<button>` 分支 | 需改 `schema-table.tsx` + 协议 toolbar item 形状 → 突破冻结范围 |
| E. custom 节点嵌在 table 节点内部（`children`/插槽） | ❌ 不支持：`schema-table.tsx` 无 `children` 读取；`RenderTableNode`/`props` 无插槽 | 同上 |
| F. 组件 portal 进页面 header actions 区 | ❌ 不可用（单 owner，已被 table claim） | 需改 `list-surface.tsx` 的所有权模型 |
| G. 复用现有 `export` toolbar 项（`exportUsers`，action-level custom handler） | ❌ 语义不符：GET CSV 全量下载，无选择、无 job | —— |

### 5.2 建议（推荐 A）

在 `apps/api/modules/users/schema/users.json` 的 `body`（section `users`，`:308-311`）内，**紧邻 `users-table` 之后**（或之前）插入：

```json
{
  "type": "custom",
  "id": "users-batch-export",
  "component": "batch-export-trigger",
  "props": { "targetTable": "users-table" }
}
```

并在 `users-table.props` 增加 `"selection": { "mode": "multiple" }`（对齐 `admin-list-batch.json:46-50`）。

理由（全部有码证据）：
1. **现状唯一受支持的声明位置**（§3.2 / §5.1A），零 renderer 改动 → 不触碰 ADR-0022 的同步路径。
2. **table id 必须靠 `props.targetTable` 传入**（`context` 没有 table id）：照抄 `activity-export.tsx:26-29`
   ```ts
   const targetTable =
     isRecord(node.props) && typeof node.props.targetTable === "string" && node.props.targetTable !== ""
       ? node.props.targetTable
       : "operations-table";
   ```
   `notification-center` 更贴近"贴着表放的控件"（`notifications.json:65-67` 同样用 `targetTable`）。
3. **选择读取**：`useSchemaCrud()?.selection(tableId)`（§4b）。建议对 `undefined`/`count===0` 走 `disabled`（对齐 toolbar 的 `feedback.selectRowFirst` 文案：`schema-table.tsx:1336-1348`）。
4. **提交**：`await (crud?.fetcher ?? fetch)("/api/jobs/batch-export", {method:"POST", headers:{"Content-Type":"application/json"}, body: JSON.stringify({ids: sel.keys})})`；`crud.fetcher` 是 host 注入的**带鉴权**传输（`render.tsx:256-258`、`activity-export.tsx:22`）。**不要**调 `reloadList()`（会清空选择，§4b）。
5. 反馈：`crud.notifyFeedback({kind:"success"|"error", ...})`（`render.tsx:211-212`，value `:1495`）。
6. 新组件必须同步（§3.3）：`apps/web/src/components/<new>.tsx`（自注册）+ `main.tsx` 副作用 import + `custom-components.schema.test.ts` import + `denominator-render.test.tsx`/`s5-denominator-render.test.tsx` import + i18n（`en-US.json`/`zh-CN.json` 各加 key，否则 `schema-keys.structural.test.ts` 的结构守卫可能红）。
7. **不要**给 users.json 加 `table.selection`：守卫 marker 是 `"requiresSelection"`（`capability-declaration.guard.test.ts:51`），users 页没有该键 → 加了必红；而 renderer **没有**"用了 `props.selection` 就必须声明 `table.selection`"的运行时闸门（唯一同类闸门是 `table.sort`：`conformance/table-sort.ts:79-105`；另有 `data.route-binding` `render.tsx:3140-3163`、`form.record.load` `:1637`）。`table.selection` 本身在 host 支持集内（`apps/web/src/protocol/host-support.json:18`）、since 2.2（`conformance-claim.json:38`），所以**不声明也不违约**。

### 5.3 若最终坚持放在 toolbar 那一行

必须改 3 处（属未冻结范围，需另立决策）：`schema-table.tsx:1326-1367`（toolbar 映射加 custom 分支）、协议 toolbar item 形状（`docs/schemas/.../node.schema.json` 上游 pin + 本地校验）、以及 capability 守卫语义（toolbar 里的 custom 与 `actions.batch.request` 的关系）。**不建议在本轮做。**

---

## 6. 同步 batch 路径的回归面（证明"零回归"的锚点）

改动集合（新增 custom 节点 + users 表加 `selection`）**不触碰**下列任何文件；下列断言因此可原样复跑作为 no-regression 证据。

### 6.1 端到端集成测试（唯一跑真 UI 的真同步链路）

`apps/web/src/app/representative-pages.integration.test.tsx`

| 位置 | 测试名 | 断言要点 |
|---|---|---|
| `:360-403` | `fails closed when a batch toolbar trigger fires with an empty selection` | 「Batch delete」按钮 `disabled === true`；点击后**没有**确认文案 `Delete the selected users?`；`fetchSpy` 中 `POST /api/users/batch-delete` 命中数 **0** |
| `:491-575` | `runs the ADR-0022 batch flow end-to-end (select → confirm → request → reload clears selection)` | 2 个 `input[aria-label='Select row']`；按钮初始 disabled（`:542`）；勾两行 → enabled（`:549`）→ 点击 → 出现确认 `Delete the selected users?`（`:553`）→ Confirm → **恰好 1 次** POST `/api/users/batch-delete`（`:563-564`）且 body `{"ids":["usr-1","usr-2"]}`（`:565`）→ 成功反馈 `Items deleted`（`:568`）→ 选择被清空、按钮回 disabled（`:569-574`） |

被测 fixture 页：`apps/api/modules/dev/examples/schema/admin-list-batch.json`（capabilities `:6-14`；`selection.mode=multiple` `:46-50`；toolbar `batchMapping.body.ids = $selection.keys` `:70-85`）；页面登记见 `representative-pages.integration.test.tsx:50-63`（`MIGRATED_PAGE_IDS` 含 `admin-list-batch`）与 `apps/api/modules/dev/examples/manifest/fragment.json:14-18`。

### 6.2 schema / capability 守卫

`apps/web/src/protocol/capability-declaration.guard.test.ts`

- `:38-58` marker 表：**`:57` `"actions.batch.request": (text) => /\/batch-delete/.test(text)`**；`:51` `"table.selection": (text) => /"requiresSelection"/.test(text)`。
- `:93-129` **动态**遍历 `apps/api/modules/**/schema/*.json`，对每个文件断言"声明的能力必须在本文件文本上有 marker（或 EXEMPT / `INTENT_OVERRIDES`）"（`:106-126`，`expect(violations).toEqual([])`）。
- 风险提示：`admin-list-batch.json` 是**唯一**同时声明 `actions.batch.request` + `table.selection` 的页面；若有人把它的 `url` 从 `/api/users/batch-delete` 改掉（例如改走异步入口），`:57` marker 立即失效 → 守卫红。**冻结要求"byte-identical"，因此该文件不得动。**

### 6.3 上游构造用例（构造层不可回归）

`apps/web/src/protocol/conformance/stage3-fixtures.test.ts:397-407`：

```ts
describe("stage 3 · request-construction fixtures (incl. batch, I-PROTO-FULL-001)", () => {
  const suite = loadSuite("request-construction");
  // Full suite executes: 64 non-batch + 11 batchRequest (ADR-0022 include).
  assertCoverage(suite.value, cases.map((c) => c.id), {}, "request-construction");
  for (const fixtureCase of cases) {
    it(fixtureCase.id, () => { expect(constructRequest(fixtureCase.input)).toEqual(fixtureCase.expected); });
  }
```

即上游 `apps/web/src/protocol/upstream/request-construction.cases.json`（11 条 `kind:"batchRequest"`，首个在 `:1026`，含 `/api/orders/batch-delete` `:1029`）**逐条执行**；另有 sha256 pin（`stage3-fixtures.test.ts:114-149`，pin 源 `upstream/provenance-v2.9.json:92`）。

### 6.4 其余"改到就会红"的清单（本次**不**涉及，但列全）

- `apps/web/src/renderer/schema-table.tsx:1326-1367`（toolbar 批量分支：`isBatch` 判定 `:1329-1338`、`invokeBatchAction` 调用 `:1351-1355`）——**无专属单测**，只由 6.1 覆盖。
- `apps/web/src/renderer/render.tsx:685-779`（`runBatchRequest`：`:765-778` `if (!response.ok) {...} return { ok: true };` —— **成功路径只判 `response.ok`，202 body 被丢弃**）、`:1212-1254`（`invokeBatchAction`：成功 → `setFeedback` + `reloadList()`）、`:1287-1310`（`resolveConfirm` 的 batch 分支）、`:328-335`（`batchSuccessMessageFor`，按 URL 后缀 `/batch-delete` 取复数文案）。
  → 这三处都是**冻结"字节级不变"的对象**；异步入口必须绕开（这正是"不复用"的技术理由：202 + job 投影会被整段丢弃）。
- `apps/web/src/app/searchable.ts:470`（命令面板对 `requiresSelection`/`batchMapping` 触发器的过滤）。

---

## 7. `apps/api/modules/users/schema/users.json` 形状小结

- **meta**（`:2-17`）：`pageId: "users"`、`title: "Users"`、`protocolVersion: "2.7"`；
  `requiredCapabilities`（`:6-16`，9 项）= `app.manifest`、`app.navigation`、`permissions.inheritance`、`actions.row.request`、`actions.page.trigger`、`table.sort`、`form.controls.advanced`、`actions.upload`、`record.view.load`。
  → **无 `table.selection`，无 `actions.batch.request`**。
- **actions**（`:18-307`，共 18 个）：`request` ×10（createUser/updateUser/updateUserRoles/changeUserPassword/deleteUser/enableUser/disableUser/unlockUser/resetUserMfa/submitImport）、`upload` ×1（`uploadCsv` `:110-114`）、`custom` handler ×1（`exportUsers` → `"export.users"` `:115-118`；URL 白名单 `render.tsx:352-353` = `GET /api/export/users`，**整表 CSV blob 下载**，无选择、无 job）、`navigate` ×1（`openInvites` `:19-22`）、`modal` ×5（openImport/openCreate/openEdit/openRoles/openPassword）。
- **body**（`:308-637`）：`section` id `users`，3 个子节点，顺序为
  1. `form` id **`users-search`**（`:312-366`，`mode:"search"`，`targetTable:"users-table"`，字段 `q`/`enabled`/`locked`）；
  2. `table` id **`users-table`**（`:367-585`）；
  3. `recordView` id `user-detail`（`:586-635`）。
- **`users-table` 详情**（`:367-585`）：`permissionCascade {edit,delete}` + permissions（`:370-379`）；`props.columns` 8 列（id/username/name/email/roles/enabled/mfaEnabled/updatedAt，`:381-427`）；`props.dataSource = "/api/users"`（`:428`）；`props.toolbar` 4 项（`:429-464`）= `create→openCreate`、`invites→openInvites`、`export→exportUsers`、`import→openImport`；`props.actions` 8 个行操作（`:465-583`）。
- **`props.selection`：不存在**（已逐行确认 `props` 块 `:380-584` 无 `selection` 键）。→ 与任务书判断一致。
- 影响面提示：其它消费 users.json 的测试（改 schema 会一并被牵动，需保持绿）
  `apps/web/src/i18n/schema-keys.structural.test.ts:30`（结构键完备性，分母含 users.json）、
  `apps/web/src/i18n/ui-bilingual.test.tsx:30`、
  `apps/web/src/i18n/s5-denominator-render.test.tsx:59`、
  `apps/web/src/renderer/error-localization.test.tsx:30`、
  `apps/web/src/renderer/denominator-render.test.tsx:129-156`（36 页全量渲染 + 无 unknown custom 占位）。

---

## 8. 与本次冻结设计的接口对照（便于直接落地）

| 冻结项 | 侦察结论 | 落点 |
|---|---|---|
| 触发机制 = `registerCustomComponent` 自定义组件 | ✅ 唯一可行（toolbar 不支持 custom，header portal 单 owner） | §3.4 / §5.1 |
| 不复用同步 batch 路径 | ✅ 技术上必须绕开：同步路径成功只判 `response.ok` 并 `reloadList()`，202 body 与选择都会被丢弃/清空 | `render.tsx:765-778`、`:1240-1248`、`:950-957` |
| schema 不得声明 `actions.batch.request` | ✅ 不需要声明；且 `INTENT_OVERRIDES` 也未登记 users 页，声明了就会红 | `capability-declaration.guard.test.ts:57,106-126` |
| 异步提交 `POST /api/jobs/batch-export` → 202 + job 投影 | ⚠️ 后端路由**尚不存在**（见 §待确认 U-1）；前端可用 `crud.fetcher` 直接提交并自行解析 202 body | `render.tsx:256-258`；`activity-export.tsx:22` |
| 选择来源 | `useSchemaCrud().selection("users-table")`，前置 `props.selection.mode=multiple` | §1.4 / §4b |

---

## 待确认 / 未知

| # | 未知项 | 现状证据 | 影响 / 需要谁裁决 |
|---|---|---|---|
| U-1 | **`POST /api/jobs/batch-export` 后端路由是否已存在**（本次只读前端侦察，未展开后端） | `apps/api/internal/handler/jobs.go:34` 定义 `JobsBasePath = "/api/jobs"`；`provider.go:49-51` 只声明 `GET /api/jobs`、`GET /api/jobs/{id}`、`GET /api/jobs/{id}/result`；全仓 `*.go` grep `batch-export` 只命中测试里的 **kind 字符串** `jobs.batch-export`（`internal/jobs/list_test.go:26-70`、`internal/handler/jobs_test.go:54-131`），**未见 POST 提交路由** | 前端可先按冻结契约实现（POST + 202）；联调前需确认路由/权限键（疑似 `jobs.read` 之外的新键？） |
| U-2 | 202 响应体的 **job 投影字段名**（`id`/`kind`/`status`/`resultUrl`…） | 现有读面投影见 `internal/handler/jobs.go`、`internal/jobs/list.go:132`（`ResultURL`）；测试断言字段 `resultUrl`（`jobs_test.go:261`） | 决定前端如何展示"已提交"与后续跳转（job center 页 `jobs/schema/jobs.json`） |
| U-3 | 异步入口是否需要**新权限键**（如 `users.export` / `jobs.submit`）与按钮的权限表达式 | users 页现有导出用 `data.export`（`users.json:450-452`）；job 模块权限为 `jobs.read`（`profile.go:217`） | 需要产品/治理裁决；影响 schema 的 `permissions` 与 §5.2 的 disabled 逻辑 |
| U-4 | 提交成功后**是否清空选择** | `reloadList()` 会清空（`render.tsx:950-957`），但异步提交不应触发 reload | 需要 UX 裁决；组件可显式 `clearSelection`（`render.tsx:272`）或保留选择 |
| U-5 | 提交**上限**（选中行数）与 0 选中的交互（disabled vs 提示） | toolbar 先例：0 选中 → disabled + `feedback.selectRowFirst`（`schema-table.tsx:1336-1348`） | 需与后端 `batch-export` 的入参上限一致 |
| U-6 | 新组件的 i18n key 命名与文案（中英） | 先例 `monitoringRefresh.*`（`en-US.json:83-87`） | 影响 `schema-keys.structural.test.ts` 是否触发（新增 schema 文本必须成对给 key） |
| U-7 | 本报告归档路径 | 任务书给的 `GOAL-004-r3-async-batch-承接/`（含中文"承接"）不存在；实际目标目录是 `GOAL-004-r3-async-batch-operation/`（有 `attachments/` + `.gitkeep`） | 是否移入 `GOAL-004-r3-async-batch-operation/attachments/` 由编排器决定（未擅自改动目标目录） |
| U-8 | `crud.selection()` 这条 seam 是否被治理视为"稳定契约" | 它是 `SchemaCrudValue` 的公开成员并被 `renderer/index.ts` 导出，但**尚无任何自定义组件使用**，也无专属测试 | 若被视为不稳定，需在本次实施中补一条组件级测试把它钉住（建议：新组件单测 + 集成断言） |
