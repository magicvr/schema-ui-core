import {
  Activity,
  Boxes,
  CircleHelp,
  FolderKanban,
  FormInput,
  Home,
  LayoutDashboard,
  LogOut,
  PanelLeft,
  Pencil,
  Search,
  Settings,
  Table2,
  UserRound,
  Zap,
  type LucideIcon,
} from "lucide-react";
import { useEffect, useMemo, useState } from "react";

import {
  applyDocumentBranding,
  DEFAULT_SITE_TITLE,
  fetchBranding,
  subscribeToBrandingChanges,
  type Branding,
} from "@/app/branding";
import { projectNavigation, type ProjectedItem } from "@/app/navigation";
import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import {
  type AppManifest,
  type NavigationContext,
  type PageEntry,
  matchRoute,
  resolveInitialRoute,
} from "@/protocol/app-manifest";
import { PageSchemaError, loadPageDocument } from "@/protocol/load-page";
import type { RenderPageDocument } from "@/renderer/render";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table.tsx";

const iconRegistry: Record<string, LucideIcon> = {
  activity: Activity,
  boxes: Boxes,
  dashboard: LayoutDashboard,
  folder: FolderKanban,
  form: FormInput,
  help: CircleHelp,
  home: Home,
  logout: LogOut,
  menu: PanelLeft,
  pen: Pencil,
  reaction: Zap,
  search: Search,
  settings: Settings,
  table: Table2,
  user: UserRound,
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
  currentUser?: { id: string; name?: string } | null;
  /** Revokes the session (AuthProvider flips to the login page). */
  onLogout?: () => void;
  /** Optional branding override (tests); defaults to live GET /api/branding. */
  branding?: Branding;
}

function currentLocationPath() {
  return `${window.location.pathname}${window.location.search}`;
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
  const className = item.active
    ? horizontal
      ? "flex min-h-10 items-center gap-2 border-b-2 border-primary bg-accent px-3 py-2 text-sm font-medium text-accent-foreground"
      : "flex min-h-10 items-center gap-3 border-l-2 border-primary bg-accent px-3 py-2 text-sm font-medium text-accent-foreground"
    : horizontal
      ? "flex min-h-10 items-center gap-2 border-b-2 border-transparent px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
      : "flex min-h-10 items-center gap-3 border-l-2 border-transparent px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground";
  const content = (
    <>
      {iconFor(item.icon)}
      <span>{item.label}</span>
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

function flattenNavigation(items: ProjectedItem[]): ProjectedItem[] {
  return items.flatMap((item) => (item.type === "group" ? item.items : [item]));
}

type SchemaSurfaceState =
  | { status: "loading" }
  | { status: "error"; error: PageSchemaError }
  | { status: "ready"; document: unknown };

/** Unified, fail-closed surface for a failed page-schema load or validation. */
function PageSchemaErrorSurface({ error }: { error: PageSchemaError }) {
  return (
    <section role="alert" className="space-y-6" aria-labelledby="schema-error-title">
      <div className="space-y-2">
        <p className="text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
          Page schema error
        </p>
        <h1 id="schema-error-title" className="text-3xl font-semibold tracking-tight">
          {error.code}
        </h1>
        <p className="max-w-2xl text-sm leading-6 text-muted-foreground">{error.message}</p>
        <p className="font-mono text-xs text-muted-foreground">{error.url}</p>
      </div>
      {error.issues !== undefined && error.issues.length > 0 ? (
        <ul className="space-y-1 text-sm text-destructive">
          {error.issues.map((issue, index) => (
            <li key={index}>
              <code className="font-mono">{issue.path}</code>: {issue.message}
            </li>
          ))}
        </ul>
      ) : null}
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
  context,
  fetcher,
  resourceFetcher,
}: {
  page: PageEntry;
  params: Record<string, string>;
  context: NavigationContext;
  fetcher?: typeof fetch;
  resourceFetcher?: typeof fetch;
}) {
  const [state, setState] = useState<SchemaSurfaceState>({ status: "loading" });

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
        Loading page schema…
      </p>
    );
  }
  if (state.status === "error") {
    return <PageSchemaErrorSurface error={state.error} />;
  }
  return (
    <RenderPage
      document={state.document as RenderPageDocument}
      context={context as unknown as Record<string, unknown>}
      tableRenderer={(node) => <SchemaTable node={node} fetcher={resourceFetcher} />}
      dataFetcher={resourceFetcher}
    />
  );
}

function PageSurface({
  manifest,
  path,
  onNavigate,
  navigationContext,
  schemaFetcher,
  resourceFetcher,
}: {
  manifest: AppManifest;
  path: string;
  onNavigate: (href: string) => void;
  navigationContext: NavigationContext;
  schemaFetcher?: typeof fetch;
  resourceFetcher?: typeof fetch;
}) {
  const route = useMemo(() => matchRoute(manifest.pages, path), [manifest, path]);
  const homePage = manifest.pages.find((page) => page.pageId === manifest.app.homePageRef);
  if (route === undefined) {
    return (
      <section className="max-w-2xl space-y-6" aria-labelledby="fallback-title">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
            Route fallback
          </p>
          <h1 id="fallback-title" className="text-3xl font-semibold tracking-tight">
            Page not found
          </h1>
          <p className="text-sm leading-6 text-muted-foreground">
            No manifest page matches <code className="font-mono">{path}</code>.
          </p>
        </div>
        <Button
          type="button"
          variant="outline"
          onClick={() => onNavigate(homePage?.route ?? "/")}
        >
          <Home aria-hidden="true" className="size-4" />
          Return to home
        </Button>
      </section>
    );
  }

  const pageTitle = route.page.title ?? route.page.titleKey ?? route.page.pageId;
  return (
    <section className="space-y-8" aria-labelledby="page-title">
      <div className="flex flex-wrap items-start justify-between gap-6 border-b border-border pb-6">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
            Admin workspace
          </p>
          <h1 id="page-title" className="text-3xl font-semibold tracking-tight">
            {pageTitle}
          </h1>
        </div>
        <div className="border border-border bg-card px-4 py-3 text-right text-xs text-muted-foreground">
          <p className="font-medium text-foreground">{route.page.pageId}</p>
          <p className="mt-1 font-mono">{route.page.route}</p>
        </div>
      </div>
      <SchemaPageSurface
        page={route.page}
        params={route.params}
        context={navigationContext}
        fetcher={schemaFetcher}
        resourceFetcher={resourceFetcher}
      />
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
  const [branding, setBranding] = useState<Branding>(
    () => brandingProp ?? { siteTitle: DEFAULT_SITE_TITLE, logoUrl: "" },
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
      setPath(initial?.path ?? requested);
    };
    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, [manifest]);

  const onNavigate = (href: string) => {
    if (!href.startsWith("/")) {
      return;
    }
    window.history.pushState({}, "", href);
    setPath(currentLocationPath());
  };

  const projection = useMemo(
    () => projectNavigation(manifest, path, navigationContext),
    [manifest, navigationContext, path],
  );
  const appName = branding.siteTitle || DEFAULT_SITE_TITLE;
  const showLogo = branding.logoUrl !== "";

  return (
    <div className="min-h-screen bg-background text-foreground">
      {accountError !== undefined ? (
        <div
          role="alert"
          className="border-b border-destructive/50 bg-destructive/10 px-4 py-2 text-sm text-destructive"
        >
          Account session failed to load; permissions and navigation may be incomplete.
        </div>
      ) : null}
      <header className="sticky top-0 z-20 border-b border-border bg-background/95 backdrop-blur">
        <div className="flex min-h-16 items-center gap-4 px-4 sm:px-6">
          <a
            href={manifest.app.homePageRef ? "/" : "#main"}
            className="flex min-w-0 items-center gap-3"
            onClick={(event) => {
              if (manifest.app.homePageRef) {
                event.preventDefault();
                const home = manifest.pages.find(
                  (page) => page.pageId === manifest.app.homePageRef,
                );
                onNavigate(home?.route ?? "/");
              }
            }}
          >
            {showLogo ? (
              <img
                src={branding.logoUrl}
                alt=""
                className="size-9 shrink-0 object-contain"
              />
            ) : null}
            <div className="min-w-0">
              <p className="truncate text-sm font-semibold">{appName}</p>
              <p className="truncate text-xs text-muted-foreground">Admin console</p>
            </div>
          </a>

          <nav className="ml-auto hidden items-center gap-1 lg:flex" aria-label="Primary navigation">
            <NavigationItems items={projection.top} onNavigate={onNavigate} horizontal />
          </nav>

          <div className="ml-auto flex items-center gap-2 lg:ml-4">
            <ThemeToggle />
            {currentUser !== undefined && currentUser !== null ? (
              <div className="flex items-center gap-2">
                <span className="hidden text-xs text-muted-foreground sm:inline">
                  {currentUser.name ?? currentUser.id}
                </span>
                <Button type="button" variant="outline" size="sm" onClick={onLogout}>
                  <LogOut aria-hidden="true" className="size-4" />
                  Sign out
                </Button>
              </div>
            ) : null}
            {projection.user.length > 0 ? (
              <nav className="hidden items-center gap-1 lg:flex" aria-label="User navigation">
                <NavigationItems items={projection.user} onNavigate={onNavigate} horizontal />
              </nav>
            ) : null}
          </div>
        </div>
      </header>

      <nav className="overflow-x-auto border-b border-border px-3 py-2 lg:hidden" aria-label="Workspace navigation">
        <NavigationItems
          items={flattenNavigation([
            ...projection.top,
            ...projection.sidebar,
            ...projection.user,
          ])}
          onNavigate={onNavigate}
          horizontal
        />
      </nav>

      <div className="mx-auto flex max-w-[1440px]">
        <aside className="hidden w-64 shrink-0 border-r border-border px-3 py-6 lg:block">
          <div className="mb-4 flex items-center gap-2 px-3 text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
            <PanelLeft aria-hidden="true" className="size-4" />
            <span>Workspace</span>
          </div>
          <NavigationItems items={projection.sidebar} onNavigate={onNavigate} />
        </aside>

        <main id="main" className="min-w-0 flex-1 px-4 py-8 sm:px-6 lg:px-10 lg:py-10">
          <PageSurface
            manifest={manifest}
            path={path}
            onNavigate={onNavigate}
            navigationContext={navigationContext}
            schemaFetcher={schemaFetcher}
            resourceFetcher={resourceFetcher}
          />
        </main>
      </div>
    </div>
  );
}
