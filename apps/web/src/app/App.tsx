import {
  Activity,
  Boxes,
  CircleHelp,
  FolderKanban,
  Home,
  LayoutDashboard,
  LogOut,
  PanelLeft,
  Settings,
  UserRound,
  type LucideIcon,
} from "lucide-react";
import { useEffect, useMemo, useState } from "react";

import { projectNavigation, type ProjectedItem } from "@/app/navigation";
import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import {
  type AppManifest,
  type NavigationContext,
  matchRoute,
  resolveInitialRoute,
} from "@/protocol/app-manifest";

const iconRegistry: Record<string, LucideIcon> = {
  activity: Activity,
  boxes: Boxes,
  dashboard: LayoutDashboard,
  folder: FolderKanban,
  help: CircleHelp,
  home: Home,
  logout: LogOut,
  menu: PanelLeft,
  settings: Settings,
  user: UserRound,
};

export interface AppProps {
  manifest: AppManifest;
  navigationContext?: NavigationContext;
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

function PageSurface({
  manifest,
  path,
  onNavigate,
}: {
  manifest: AppManifest;
  path: string;
  onNavigate: (href: string) => void;
}) {
  const route = matchRoute(manifest.pages, path);
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
          <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
            This surface is selected from the application manifest. The page renderer remains a
            later protocol boundary.
          </p>
        </div>
        <div className="border border-border bg-card px-4 py-3 text-right text-xs text-muted-foreground">
          <p className="font-medium text-foreground">{route.page.pageId}</p>
          <p className="mt-1 font-mono">{route.page.route}</p>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <div className="border border-border bg-card p-5">
          <p className="text-xs font-medium uppercase tracking-[0.12em] text-muted-foreground">
            Route status
          </p>
          <p className="mt-3 text-2xl font-semibold">Matched</p>
          <p className="mt-1 text-sm text-muted-foreground">D4a template resolution</p>
        </div>
        <div className="border border-border bg-card p-5">
          <p className="text-xs font-medium uppercase tracking-[0.12em] text-muted-foreground">
            Protocol
          </p>
          <p className="mt-3 text-2xl font-semibold">{manifest.protocolVersion}</p>
          <p className="mt-1 text-sm text-muted-foreground">App manifest contract</p>
        </div>
        <div className="border border-border bg-card p-5">
          <p className="text-xs font-medium uppercase tracking-[0.12em] text-muted-foreground">
            Navigation
          </p>
          <p className="mt-3 text-2xl font-semibold">{manifest.pages.length}</p>
          <p className="mt-1 text-sm text-muted-foreground">Registered pages</p>
        </div>
      </div>

      <div className="border border-dashed border-border bg-background p-6">
        <div className="flex items-start gap-3">
          <Activity aria-hidden="true" className="mt-0.5 size-5 text-muted-foreground" />
          <div>
            <h2 className="font-medium">Manifest-driven shell is ready</h2>
            <p className="mt-1 text-sm leading-6 text-muted-foreground">
              Navigation, route matching, default landing, and fallback behavior are active. A
              later renderer can own the page schema at <code className="font-mono">{route.page.schemaUrl}</code>.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}

export function App({ manifest, navigationContext = {} }: AppProps) {
  const [path, setPath] = useState(() => {
    const requested = currentLocationPath();
    const initial = resolveInitialRoute(manifest, requested);
    if (initial?.source === "home" && requested !== initial.path) {
      window.history.replaceState({}, "", initial.path);
    }
    return initial?.path ?? requested;
  });

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
  const appName = manifest.app.name ?? manifest.app.nameKey ?? manifest.app.appId;

  return (
    <div className="min-h-screen bg-background text-foreground">
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
            <div className="flex size-9 shrink-0 items-center justify-center bg-primary text-primary-foreground">
              <PanelLeft aria-hidden="true" className="size-4" />
            </div>
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
          <PageSurface manifest={manifest} path={path} onNavigate={onNavigate} />
        </main>
      </div>
    </div>
  );
}
