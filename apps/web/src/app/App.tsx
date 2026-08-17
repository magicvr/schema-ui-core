import {
  Activity,
  Boxes,
  ChevronDown,
  CircleHelp,
  FolderKanban,
  FormInput,
  Home,
  LayoutDashboard,
  LogOut,
  Menu,
  PanelLeft,
  Pencil,
  Search,
  Settings,
  Shield,
  Table2,
  UserRound,
  Wallet,
  X,
  Zap,
  type LucideIcon,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState, type MouseEvent } from "react";

import {
  applyDocumentBranding,
  defaultBranding,
  DEFAULT_SITE_TITLE,
  fetchBranding,
  subscribeToBrandingChanges,
  type Branding,
} from "@/app/branding";
import { projectNavigation, type ProjectedItem } from "@/app/navigation";
import { LocaleSwitcher } from "@/components/locale-switcher";
import { ThemeToggle } from "@/components/theme-toggle";
import { Breadcrumbs, resolveBreadcrumbTrail } from "@/components/ui/breadcrumbs";
import { resolveTextProp } from "@/i18n/catalog";
import { useTranslate } from "@/i18n/runtime";
import { applySystemDefaultTheme } from "@/theme/theme";
import {
  type AppManifest,
  type NavigationContext,
  type PageEntry,
  matchRoute,
  resolveInitialRoute,
  stripPathQuery,
} from "@/protocol/app-manifest";
import { PageSchemaError, loadPageDocument } from "@/protocol/load-page";
import type { RenderPageDocument } from "@/renderer/render";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table.tsx";
import { HostFailureScreen } from "@/app/HostFailureScreen";
import { NotificationBell } from "@/app/notification-bell";
import { nextFailureId, type HostFailure } from "@/host/failure";

const iconRegistry: Record<string, LucideIcon> = {
  activity: Activity,
  boxes: Boxes,
  dashboard: LayoutDashboard,
  folder: FolderKanban,
  form: FormInput,
  help: CircleHelp,
  home: Home,
  logout: LogOut,
  menu: Menu,
  pen: Pencil,
  reaction: Zap,
  search: Search,
  settings: Settings,
  shield: Shield,
  table: Table2,
  user: UserRound,
  wallet: Wallet,
};

export interface AppProps {
  manifest: AppManifest;
  navigationContext?: NavigationContext;
  /** Set when the boot /me session failed; surfaces a non-blocking notice. */
  accountError?: unknown;
  /** Injectable fetch for page-schema documents (defaults to `globalThis.fetch`). */
  schemaFetcher?: typeof fetch;
  /** Injectable fetch for table data sources such as `/api/users` (GOAL-011). */
  resourceFetcher?: typeof fetch;
  /** Authenticated user rendered in the header; present → show a sign-out button. */
  currentUser?: { id: string; name?: string; avatarUrl?: string } | null;
  /** Revokes the session (AuthProvider flips to the login page). */
  onLogout?: () => void;
  /** Optional branding override (tests); defaults to live GET /api/branding. */
  branding?: Branding;
}

function currentLocationPath() {
  return `${window.location.pathname}${window.location.search}`;
}

/** Route-not-found failure occurrence (new ID per distinct unmatched path). */
function routeNotFoundFailure(): HostFailure {
  return {
    failureVersion: "1.0",
    failureId: nextFailureId(),
    scope: "route",
    kind: "not-found",
    hostCode: "HOST_ROUTE_NOT_FOUND",
    message: { messageKey: "hostFailure.notFound" },
    recoveryActions: [{ type: "home" }, { type: "back" }],
  };
}

/**
 * Declared host-owned paths (ADR-0036 D3a): host surfaces outside the manifest
 * page registry. An authenticated user landing on one is taken home rather
 * than shown HOST_ROUTE_NOT_FOUND — the path belongs to the Host, not the app.
 */
const HOST_OWNED_PATHS = ["/login"];

/**
 * W13 T-02: shell brand link (logo + site title + subtitle). Rendered twice —
 * once in the mobile-only brand bar, once in the desktop single-row header —
 * with the caller's className deciding visibility. The logo/title block is a
 * single semantic unit either way.
 */
function BrandLink({
  href,
  onClick,
  branding,
  appName,
  t,
  className,
}: {
  href: string;
  onClick: (event: MouseEvent<HTMLAnchorElement>) => void;
  branding: Branding;
  appName: string;
  t: (key: string) => string;
  className?: string;
}) {
  const showLogo =
    branding.logoUrl !== "" || branding.logoUrlLight !== "" || branding.logoUrlDark !== "";
  return (
    <a href={href} onClick={onClick} className={className}>
      {showLogo ? (
        <>
          {/* VP-007 S3: light/dark logo variants follow the active theme via CSS. */}
          {branding.logoUrlLight !== "" ? (
            <img
              src={branding.logoUrlLight}
              alt=""
              className="size-8 shrink-0 object-contain dark:hidden"
            />
          ) : null}
          {branding.logoUrlDark !== "" ? (
            <img
              src={branding.logoUrlDark}
              alt=""
              className="hidden size-8 shrink-0 object-contain dark:block"
            />
          ) : null}
          {branding.logoUrlLight === "" && branding.logoUrlDark === "" ? (
            <img src={branding.logoUrl} alt="" className="size-8 shrink-0 object-contain" />
          ) : null}
        </>
      ) : (
        <span
          aria-hidden="true"
          className="inline-flex size-8 shrink-0 items-center justify-center rounded-md bg-primary text-xs font-semibold text-primary-foreground"
        >
          {appName.slice(0, 1).toUpperCase()}
        </span>
      )}
      <div className="min-w-0">
        <p className="truncate text-sm font-semibold tracking-tight">{appName}</p>
        <p className="truncate text-[11px] text-muted-foreground">{t("shell.adminConsole")}</p>
      </div>
    </a>
  );
}

/**
 * GOAL-015 semantic breadcrumbs: inner pages reached by row navigation declare
 * their parent page here (web-shell level, no protocol change). The trail is
 * hierarchy, not visit history (user ruling 2026-08-14): 首页 => 一级页 => ...
 * => n级内页, rooted at the manifest homePageRef.
 */
const BREADCRUMB_PAGE_PARENTS: Record<string, string> = {
  "dictionary-entries": "data-dictionary",
  "task-runs": "scheduled-tasks",
  // GOAL-020 (user 2026-08-16): wallet entries is the wallet inner page.
  "wallet-entries": "wallet",
};

// Parses the current URL's query string into a plain record; deep-linked query
// parameters reach $context.route.query.* bindings through the render context.
function parseLocationQuery(): Record<string, string> {
  const search = window.location.search;
  if (search === "") {
    return {};
  }
  return Object.fromEntries(new URLSearchParams(search).entries());
}

function iconFor(name: string | undefined) {
  if (name === undefined) {
    return null;
  }
  const Icon = iconRegistry[name];
  return Icon === undefined ? null : <Icon aria-hidden="true" className="size-4" />;
}

function NavigationLink({
  item,
  onNavigate,
  horizontal = false,
}: {
  item: Extract<ProjectedItem, { type: "link" }>;
  onNavigate: (href: string) => void;
  horizontal?: boolean;
}) {
  // D-004 shell language: Linear/Vercel — rounded side items, subtle active fill.
  const className = item.active
    ? horizontal
      ? "flex min-h-9 items-center gap-2 rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-foreground"
      : "flex min-h-9 items-center gap-3 rounded-md bg-accent px-3 py-2 text-sm font-medium text-accent-foreground"
    : horizontal
      ? "flex min-h-9 items-center gap-2 rounded-md px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent/70 hover:text-accent-foreground"
      : "flex min-h-9 items-center gap-3 rounded-md px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-accent/70 hover:text-accent-foreground";
  const content = (
    <>
      {iconFor(item.icon)}
      <span className="truncate">{item.label}</span>
    </>
  );

  if (item.href === undefined) {
    return (
      <span aria-current={item.active ? "page" : undefined} className={className}>
        {content}
      </span>
    );
  }

  const href = item.href;

  return (
    <a
      href={href}
      aria-current={item.active ? "page" : undefined}
      className={className}
      onClick={(event) => {
        if (
          event.defaultPrevented ||
          event.button !== 0 ||
          event.metaKey ||
          event.ctrlKey ||
          event.shiftKey ||
          event.altKey
        ) {
          return;
        }
        event.preventDefault();
        onNavigate(href);
      }}
    >
      {content}
    </a>
  );
}

function NavigationItems({
  items,
  onNavigate,
  horizontal = false,
}: {
  items: ProjectedItem[];
  onNavigate: (href: string) => void;
  horizontal?: boolean;
}) {
  return (
    <div className={horizontal ? "flex min-w-max items-center gap-1" : "space-y-1"}>
      {items.map((item, index) =>
        item.type === "link" ? (
          <NavigationLink
            key={`${item.href}-${index}`}
            item={item}
            onNavigate={onNavigate}
            horizontal={horizontal}
          />
        ) : (
          <section
            key={`${item.label}-${index}`}
            className={horizontal ? "flex items-center gap-1" : "pt-3"}
          >
            <div className="flex items-center gap-2 px-3 pb-1 text-[11px] font-semibold uppercase tracking-[0.12em] text-muted-foreground">
              {iconFor(item.icon)}
              <span>{item.label}</span>
            </div>
            <div className={horizontal ? "flex items-center gap-1" : "space-y-1 pl-2"}>
              {item.items.map((child, childIndex) => (
                <NavigationLink
                  key={`${child.href}-${childIndex}`}
                  item={child}
                  onNavigate={onNavigate}
                  horizontal={horizontal}
                />
              ))}
            </div>
          </section>
        ),
      )}
    </div>
  );
}


/**
 * T-01 (GOAL-013 D-002): topbar user dropdown — one trigger for every
 * breakpoint (avatar + display name + chevron), menu from projection.user
 * declaration order (个人中心 → 我的钱包[if present] → 设置) plus a
 * divider and the shell-injected 退出登录 last item. Closes on outside
 * click / Escape (parity with the notification bell). The mobile drawer no
 * longer repeats the user chain (D-002 §3).
 */
function UserMenu({
  items,
  user,
  onNavigate,
  onLogout,
}: {
  items: Array<Extract<ProjectedItem, { type: "link" }>>;
  user: { id: string; name?: string; avatarUrl?: string };
  onNavigate: (href: string) => void;
  onLogout?: () => void;
}) {
  const t = useTranslate();
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  const displayName = user.name !== undefined && user.name.trim() !== "" ? user.name : user.id;
  const initials = displayName.slice(0, 1).toUpperCase();

  useEffect(() => {
    if (!open) {
      return;
    }
    const onPointerDown = (event: PointerEvent) => {
      if (rootRef.current !== null && !rootRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
        return;
      }
      if (event.key !== "ArrowDown" && event.key !== "ArrowUp") {
        return;
      }
      const items = rootRef.current?.querySelectorAll<HTMLElement>('[role="menuitem"]');
      if (items === undefined || items.length === 0) {
        return;
      }
      event.preventDefault();
      const list = Array.from(items);
      const current = list.findIndex((el) => el === document.activeElement);
      const delta = event.key === "ArrowDown" ? 1 : -1;
      const next = current < 0 ? 0 : (current + delta + list.length) % list.length;
      list[next]?.focus();
    };
    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        aria-label={t("shell.userMenu")}
        aria-haspopup="menu"
        aria-expanded={open}
        onClick={() => setOpen((value) => !value)}
        className="flex min-h-9 items-center gap-2 rounded-md border border-border bg-card/60 px-2 py-1 transition-colors hover:bg-accent"
      >
        {user.avatarUrl !== undefined && user.avatarUrl !== "" ? (
          <img
            src={user.avatarUrl}
            alt=""
            className="size-6 shrink-0 rounded-full object-cover"
          />
        ) : (
          <span
            aria-hidden="true"
            className="inline-flex size-6 shrink-0 items-center justify-center rounded-full bg-primary text-[11px] font-semibold text-primary-foreground"
          >
            {initials}
          </span>
        )}
        <span className="max-w-[8rem] truncate text-xs text-foreground">
          {displayName}
        </span>
        <ChevronDown
          aria-hidden="true"
          className={`size-3.5 text-muted-foreground transition-transform ${open ? "rotate-180" : ""}`}
        />
      </button>
      {open ? (
        <div
          role="menu"
          aria-label={t("shell.userMenu")}
          className="absolute right-0 top-11 z-50 min-w-52 rounded-md border border-border bg-card py-1 shadow-lg"
        >
          {items.length > 0 ? (
            <div className="p-1">
              {items.map((item, index) => (
                <button
                  key={`${item.href}-${index}`}
                  type="button"
                  role="menuitem"
                  onClick={() => {
                    setOpen(false);
                    if (item.href !== undefined) {
                      onNavigate(item.href);
                    }
                  }}
                  className="flex w-full items-center gap-2 rounded-md px-3 py-1.5 text-left text-sm text-foreground transition-colors hover:bg-accent"
                >
                  {iconFor(item.icon)}
                  <span className="truncate">{item.label}</span>
                </button>
              ))}
            </div>
          ) : null}
          {items.length > 0 ? (
            <div role="separator" className="my-1 border-t border-border" />
          ) : null}
          <div className="p-1">
            <button
              type="button"
              role="menuitem"
              onClick={() => {
                setOpen(false);
                onLogout?.();
              }}
              className="flex w-full items-center gap-2 rounded-md px-3 py-1.5 text-left text-sm text-destructive transition-colors hover:bg-destructive/10"
            >
              <LogOut aria-hidden="true" className="size-4" />
              <span>{t("shell.signOut")}</span>
            </button>
          </div>
        </div>
      ) : null}
    </div>
  );
}

type SchemaSurfaceState =
  | { status: "loading" }
  | { status: "error"; error: PageSchemaError }
  | { status: "ready"; document: unknown };

/** Unified, fail-closed surface for a failed page-schema load or validation. */
function PageSchemaErrorSurface({ error }: { error: PageSchemaError }) {
  const t = useTranslate();
  return (
    <section role="alert" className="space-y-6" aria-labelledby="schema-error-title">
      <div className="space-y-3">
        <p className="text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
          {t("shell.pageSchemaError")}
        </p>
        <h1 id="schema-error-title" className="text-3xl font-semibold tracking-tight">
          {t("shell.pageSchemaError.title")}
        </h1>
        <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
          {t("shell.pageSchemaError.description")}
        </p>
        <button
          type="button"
          onClick={() => window.location.reload()}
          className="inline-flex h-9 cursor-pointer items-center justify-center rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90"
        >
          {t("shell.pageSchemaError.reload")}
        </button>
      </div>
      <details className="rounded-lg border border-border/70 bg-card/85 p-3 text-xs text-muted-foreground">
        <summary className="cursor-pointer text-sm font-medium">
          {t("shell.pageSchemaError.technicalDetails")}
        </summary>
        <p className="mt-2 font-mono">{error.code}</p>
        <p className="mt-1 font-mono">{error.url}</p>
        {error.issues !== undefined && error.issues.length > 0 ? (
          <ul className="mt-2 space-y-1">
            {error.issues.map((issue, index) => (
              <li key={index}>
                <code className="font-mono">{issue.path}</code>: {issue.message}
              </li>
            ))}
          </ul>
        ) : null}
      </details>
    </section>
  );
}

/**
 * Default page surface: the schema-driven path
 * `page.route → page.schemaUrl → loadPageDocument → RenderPage` (GOAL-003).
 * The hand-written EXAMPLE_PAGES are no longer part of the render path (D-003).
 */
function SchemaPageSurface({
  page,
  params,
  query,
  context,
  fetcher,
  resourceFetcher,
  onNavigate,
}: {
  page: PageEntry;
  params: Record<string, string>;
  query: Record<string, string>;
  context: NavigationContext;
  fetcher?: typeof fetch;
  resourceFetcher?: typeof fetch;
  /** Session-internal navigation for schema navigate actions (GOAL-015 F-001). */
  onNavigate?: (url: string) => void;
}) {
  const [state, setState] = useState<SchemaSurfaceState>({ status: "loading" });
  const t = useTranslate();

  useEffect(() => {
    let cancelled = false;
    setState({ status: "loading" });
    loadPageDocument(page, params, { fetcher })
      .then((document) => {
        if (!cancelled) {
          setState({ status: "ready", document });
        }
      })
      .catch((error: unknown) => {
        if (cancelled) {
          return;
        }
        setState({
          status: "error",
          error:
            error instanceof PageSchemaError
              ? error
              : new PageSchemaError(
                  "PAGE_LOAD_FAILED",
                  page.schemaUrl,
                  error instanceof Error ? error.message : "Failed to load the page schema.",
                ),
        });
      });
    return () => {
      cancelled = true;
    };
  }, [page, params, fetcher]);

  if (state.status === "loading") {
    return (
      <p role="status" className="text-sm text-muted-foreground">
        {t("shell.loadingPageSchema")}
      </p>
    );
  }
  if (state.status === "error") {
    return <PageSchemaErrorSurface error={state.error} />;
  }
  return (
    <RenderPage
      document={state.document as RenderPageDocument}
      context={
        {
          ...(context as Record<string, unknown>),
          route: { params, query },
        } as Record<string, unknown>
      }
      tableRenderer={(node) => <SchemaTable node={node} fetcher={resourceFetcher} />}
      dataFetcher={resourceFetcher}
      onNavigate={onNavigate}
    />
  );
}

function PageSurface({
  manifest,
  path,
  query,
  onNavigate,
  navigationContext,
  schemaFetcher,
  resourceFetcher,
}: {
  manifest: AppManifest;
  path: string;
  query: Record<string, string>;
  onNavigate: (href: string) => void;
  navigationContext: NavigationContext;
  schemaFetcher?: typeof fetch;
  resourceFetcher?: typeof fetch;
}) {
  const route = useMemo(() => matchRoute(manifest.pages, path), [manifest, path]);
  const homePage = manifest.pages.find((page) => page.pageId === manifest.app.homePageRef);
  const t = useTranslate();
  const hostOwned = useMemo(() => HOST_OWNED_PATHS.includes(stripPathQuery(path)), [path]);
  // Host-owned path under an authenticated session: return to the app surface
  // (ADR-0036 D3a) — no route-not-found failure for paths the Host owns.
  useEffect(() => {
    if (hostOwned) {
      onNavigate(homePage?.route ?? "/");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [hostOwned]);
  // Stable per-path failureId: redraws of the same unmatched route never
  // re-announce; a different unmatched route is a new occurrence (D1).
  const routeNotFound = useMemo(() => routeNotFoundFailure(), [path]);
  // Recovery actions must land focus on the restored surface's main heading
  // (ADR-0036 D7 rule 5).
  const focusPageTitleOnNextRender = useRef(false);
  useEffect(() => {
    if (!focusPageTitleOnNextRender.current) return;
    focusPageTitleOnNextRender.current = false;
    document.getElementById("page-title")?.focus();
  }, [path]);
  if (route === undefined) {
    if (hostOwned) {
      // Navigation back to the app is in flight (effect above); render nothing
      // while the location resolves.
      return null;
    }
    // Unmatched application route (no manifest page, no host-owned path) →
    // HOST_ROUTE_NOT_FOUND global failure surface (ADR-0036 D3/D3a).
    return (
      <HostFailureScreen
        bare
        failure={routeNotFound}
        onAction={(action) => {
          focusPageTitleOnNextRender.current = true;
          if (action.type === "home") {
            onNavigate(homePage?.route ?? "/");
          } else if (action.type === "back") {
            window.history.back();
          }
        }}
      />
    );
  }

  const pageTitle =
    resolveTextProp(route.page as unknown as Record<string, unknown>, "titleKey", "title", t) ??
    route.page.pageId;
  // GOAL-015: SEMANTIC breadcrumb trail — the page's place in the manifest
  // navigation tree plus declared parents for inner pages (user ruling
  // 2026-08-14: hierarchy, not visit history).
  const trail = resolveBreadcrumbTrail(
    manifest.pages as unknown as Parameters<typeof resolveBreadcrumbTrail>[0],
    route.page as unknown as Parameters<typeof resolveBreadcrumbTrail>[1],
    t,
    {
      navigation: manifest.navigation,
      parents: BREADCRUMB_PAGE_PARENTS,
      homePageId: manifest.app.homePageRef,
    },
  );
  // Semantic back: the nearest ancestor page excluding the home root
  // (group labels are not pages; 首页 is reachable via the trail itself).
  const backRoute = [...trail]
    .reverse()
    .find(
      (entry) =>
        !entry.current &&
        entry.route !== "" &&
        entry.pageId !== manifest.app.homePageRef,
    )?.route;
  return (
    <section className="w-full min-w-0 space-y-8" aria-labelledby="page-title">
      <div className="flex w-full min-w-0 flex-wrap items-start justify-between gap-6 border-b border-border pb-6">
        <div className="min-w-0 flex-1">
          <Breadcrumbs
            entries={trail}
            onNavigate={onNavigate}
            onBack={() => {
              if (backRoute !== undefined) {
                onNavigate(backRoute);
              }
            }}
            showBack={backRoute !== undefined}
          />
          {/* 10px breathing room between the breadcrumb and the page title. */}
          <h1 id="page-title" tabIndex={-1} className="mt-2.5 truncate text-3xl font-semibold tracking-tight outline-none">
            {pageTitle}
          </h1>
        </div>
      </div>
      <div className="w-full min-w-0">
        <SchemaPageSurface
          page={route.page}
          params={route.params}
          query={query}
          context={navigationContext}
          fetcher={schemaFetcher}
          resourceFetcher={resourceFetcher}
          onNavigate={onNavigate}
        />
      </div>
    </section>
  );
}

export function App({
  manifest,
  navigationContext = {},
  accountError,
  schemaFetcher,
  resourceFetcher,
  currentUser,
  onLogout,
  branding: brandingProp,
}: AppProps) {
  const [path, setPath] = useState(() => {
    const requested = currentLocationPath();
    const initial = resolveInitialRoute(manifest, requested);
    if (initial?.source === "home" && requested !== initial.path) {
      window.history.replaceState({}, "", initial.path);
    }
    return initial?.path ?? requested;
  });
  const [routeQuery, setRouteQuery] = useState<Record<string, string>>(() =>
    parseLocationQuery(),
  );
  const [mobileDrawerOpen, setMobileDrawerOpen] = useState(false);
  const t = useTranslate();
  const [branding, setBranding] = useState<Branding>(
    () => brandingProp ?? defaultBranding(),
  );

  useEffect(() => {
    if (brandingProp !== undefined) {
      setBranding(brandingProp);
      applyDocumentBranding(brandingProp);
      return;
    }
    let cancelled = false;
    const load = () => {
      void fetchBranding().then((next) => {
        if (!cancelled) {
          setBranding(next);
          applyDocumentBranding(next);
          applySystemDefaultTheme(next.defaultTheme);
        }
      });
    };
    load();
    const unsubscribe = subscribeToBrandingChanges(load);
    return () => {
      cancelled = true;
      unsubscribe();
    };
  }, [brandingProp]);

  useEffect(() => {
    const handlePopState = () => {
      const requested = currentLocationPath();
      const initial = resolveInitialRoute(manifest, requested);
      if (initial?.source === "home" && requested !== initial.path) {
        window.history.replaceState({}, "", initial.path);
      }
      const resolved = initial?.path ?? requested;
      setPath(resolved);
      setRouteQuery(initial?.query ?? parseLocationQuery());
    };
    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, [manifest]);

  // Mobile navigation drawer focus management (S0 D-003 §8 · F-002)：抽屉打开时
  // 焦点进入首个可聚焦元素，Tab 在抽屉内循环，Escape 关闭，关闭后焦点恢复到
  // 触发元素。
  const drawerNavRef = useRef<HTMLElement>(null);
  const drawerTriggerRef = useRef<HTMLElement | null>(null);
  useEffect(() => {
    if (!mobileDrawerOpen) {
      return;
    }
    drawerTriggerRef.current =
      document.activeElement instanceof HTMLElement ? document.activeElement : null;
    const drawer = drawerNavRef.current;
    const firstFocusable = drawer?.querySelector<HTMLElement>(
      'a[href], button:not([disabled]), [tabindex]:not([tabindex="-1"])',
    );
    (firstFocusable ?? drawer)?.focus();

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        event.stopPropagation();
        setMobileDrawerOpen(false);
        return;
      }
      if (event.key !== "Tab" || !drawer) return;
      const items = Array.from(
        drawer.querySelectorAll<HTMLElement>(
          'a[href], button:not([disabled]), [tabindex]:not([tabindex="-1"])',
        ),
      );
      if (items.length === 0) return;
      const first = items[0];
      const last = items[items.length - 1];
      const active = document.activeElement;
      if (event.shiftKey) {
        if (active === first || !drawer.contains(active)) {
          event.preventDefault();
          last.focus();
        }
      } else if (active === last || !drawer.contains(active)) {
        event.preventDefault();
        first.focus();
      }
    };
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("keydown", handleKeyDown);
      drawerTriggerRef.current?.focus();
    };
  }, [mobileDrawerOpen]);

  const onNavigate = (href: string) => {
    if (!href.startsWith("/")) {
      return;
    }
    window.history.pushState({}, "", href);
    // Keep the path free of the query string (matchRoute expects a clean path);
    // the query lives in routeQuery and reaches the render context (C8).
    const nextPath = stripPathQuery(currentLocationPath());
    setPath(nextPath);
    setRouteQuery(parseLocationQuery());
    setMobileDrawerOpen(false);
  };

  const projection = useMemo(
    () => projectNavigation(manifest, path, navigationContext, t),
    [manifest, navigationContext, path, t],
  );
  const appName = branding.siteTitle || DEFAULT_SITE_TITLE;
  // W13 T-02: shared brand-link handler for the mobile brand bar and the
  // desktop single-row header (home navigation when homePageRef is declared).
  const handleBrandClick = (event: MouseEvent<HTMLAnchorElement>) => {
    if (manifest.app.homePageRef) {
      event.preventDefault();
      const home = manifest.pages.find((page) => page.pageId === manifest.app.homePageRef);
      onNavigate(home?.route ?? "/");
    }
  };

  return (
    <div
      data-shell="admin"
      data-shell-layout="topbar-sidenav"
      className="min-h-screen bg-background text-foreground"
    >
      {accountError !== undefined ? (
        <div
          role="alert"
          className="border-b border-destructive/50 bg-destructive/10 px-4 py-2 text-sm text-destructive"
        >
          {t("shell.accountError")}
        </div>
      ) : null}
      {/* D-004 §3: sticky top bar (desktop shell language) */}
      <header
        data-shell-region="topbar"
        className="sticky top-0 z-20 border-b border-border bg-background/90 shadow-sm backdrop-blur supports-[backdrop-filter]:bg-background/75"
      >
        {/* W13 T-02: mobile brand bar — the logo + site title own a dedicated
            row on small screens (<lg) so the top functional area below keeps
            its full width instead of being squeezed by the brand block. */}
        <div
          data-shell-region="mobile-brandbar"
          className="flex items-center border-b border-border/60 px-4 py-2 lg:hidden"
        >
          <BrandLink
            href={manifest.app.homePageRef ? "/" : "#main"}
            onClick={handleBrandClick}
            branding={branding}
            appName={appName}
            t={t}
            className="flex w-full min-w-0 items-center gap-3"
          />
        </div>
        <div className="flex min-h-14 items-center gap-4 px-4 sm:px-6">
          {/* Desktop brand — back in the single-row layout at lg+. */}
          <BrandLink
            href={manifest.app.homePageRef ? "/" : "#main"}
            onClick={handleBrandClick}
            branding={branding}
            appName={appName}
            t={t}
            className="hidden min-w-0 items-center gap-3 lg:flex"
          />

          <nav className="ml-auto hidden items-center gap-1 lg:flex" aria-label="Primary navigation">
            <NavigationItems items={projection.top} onNavigate={onNavigate} horizontal />
          </nav>

          {/* W13 follow-up (user 2026-08-16): the mobile drawer trigger sits at
              the FAR LEFT of the functional row — the intuitive position for a
              button that opens the left-hand navigation drawer. The remaining
              controls stay right-aligned (ml-auto). Desktop is unaffected: the
              brand link occupies the left slot at lg+. */}
          <button
            type="button"
            aria-label={t("shell.openMenu")}
            aria-expanded={mobileDrawerOpen}
            className="inline-flex size-9 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground lg:hidden"
            onClick={() => setMobileDrawerOpen(true)}
          >
            <Menu aria-hidden="true" className="size-4" />
          </button>

          {/* T-01 (GOAL-013 D-002): user nav + signout folded into the user dropdown. */}
          <div className="ml-auto flex items-center gap-2 lg:ml-4">
            {/* W13 T-04: theme toggle on the left, language switcher on the right. */}
            <ThemeToggle />
            <LocaleSwitcher className="inline-flex" />
            {currentUser !== undefined && currentUser !== null ? (
              <NotificationBell
                fetcher={resourceFetcher}
                onViewAll={() => onNavigate("/notifications")}
                // W13 T-06: a dropdown entry opens its detail on the list page
                // (deep link ?open=<id>; the page expands + marks it read).
                onOpenItem={(id) => onNavigate("/notifications?open=" + encodeURIComponent(id))}
              />
            ) : null}
            {currentUser !== undefined && currentUser !== null ? (
              <UserMenu
                items={projection.user.filter(
                  (item): item is Extract<ProjectedItem, { type: "link" }> => item.type === "link",
                )}
                user={currentUser}
                onNavigate={onNavigate}
                onLogout={onLogout}
              />
            ) : null}
          </div>
        </div>
      </header>

      {/* Mobile navigation drawer */}
      {mobileDrawerOpen ? (
        <>
          <div
            className="fixed inset-0 z-30 bg-overlay lg:hidden"
            aria-hidden="true"
            onClick={() => setMobileDrawerOpen(false)}
          />
          <nav
            ref={drawerNavRef}
            className="fixed inset-y-0 left-0 z-40 flex w-72 flex-col border-r border-border bg-card shadow-lg lg:hidden"
            aria-label="Mobile navigation"
          >
            <div className="flex items-center justify-between border-b border-border px-4 py-3">
              <span className="text-sm font-semibold">{appName}</span>
              <button
                type="button"
                aria-label={t("shell.closeMenu")}
                className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                onClick={() => setMobileDrawerOpen(false)}
              >
                <X aria-hidden="true" className="size-4" />
              </button>
            </div>
            <div className="flex-1 space-y-1 overflow-y-auto px-3 py-4">
              <NavigationItems
                items={[...projection.top, ...projection.sidebar]}
                onNavigate={onNavigate}
              />
            </div>
          </nav>
        </>
      ) : null}

      {/*
        Fluid shell (user gap 2026-08-09): topbar is full viewport width; body must
        also track the browser — do NOT cap sidenav+main in max-w-[1440px]/mx-auto
        (that left a fixed content island while the header stretched).
      */}
      <div
        data-shell-region="body"
        data-shell-width="fluid"
        className="flex w-full min-h-[calc(100vh-3.5rem)]"
      >
        {/* D-004 §3: desktop permanent left nav ~256px (w-64) */}
        <aside
          data-shell-region="sidenav"
          data-shell-sidenav-width="256"
          className="sticky top-14 hidden h-[calc(100vh-3.5rem)] w-64 shrink-0 overflow-y-auto border-r border-border bg-card/40 px-3 py-5 lg:block"
        >
          <div className="mb-3 flex items-center gap-2 px-3 text-[11px] font-semibold uppercase tracking-[0.16em] text-muted-foreground">
            <PanelLeft aria-hidden="true" className="size-3.5" />
            <span>{t("shell.workspace")}</span>
          </div>
          <div className="space-y-0.5">
            <NavigationItems items={projection.sidebar} onNavigate={onNavigate} />
          </div>
        </aside>

        <main
          id="main"
          data-shell-region="main"
          className="min-w-0 w-full flex-1 overflow-x-auto px-4 py-6 sm:px-6 lg:px-8 lg:py-8"
        >
          <div data-shell-region="page" className="w-full min-w-0 max-w-none">
            <PageSurface
              manifest={manifest}
              path={path}
              query={routeQuery}
              onNavigate={onNavigate}
              navigationContext={navigationContext}
              schemaFetcher={schemaFetcher}
              resourceFetcher={resourceFetcher}
            />
          </div>
        </main>
      </div>
    </div>
  );
}