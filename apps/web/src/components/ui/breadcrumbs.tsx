/**
 * Breadcrumb navigation for nested admin pages (GOAL-015).
 *
 * Semantic hierarchy, not visit history (user ruling 2026-08-14):
 * the trail is the page's place in the manifest navigation tree
 * (slot → group labels → page) plus consumer-declared parents for
 * inner pages reached by row navigation (e.g. dictionary-entries →
 * data-dictionary, task-runs → scheduled-tasks). No protocol change:
 * the parent map is a web-shell declaration (BREADCRUMB_PAGE_PARENTS).
 *
 * Trail shape: 首页 => 一级页 => ... => n级内页 — the home page
 * (manifest homePageRef, the domain-root default) always leads, then nav
 * group labels and declared parents, then the current page. Visual: compact
 * 12px trail — muted clickable ancestors (hover brighten + underline), thin
 * "/" separators, brighter non-clickable current item, and an optional
 * small circular ghost back button (semantic parent) at the far left.
 */

import { ArrowLeft } from "lucide-react";

import { useTranslate } from "@/i18n/runtime";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";

export interface BreadcrumbEntry {
  /** Page id (manifest pageId); group labels use the label text as key. */
  pageId: string;
  label: string;
  /** Application route of this ancestor; empty for group labels. */
  route: string;
  /** True for the current (deepest) page. */
  current: boolean;
}

export interface BreadcrumbPage {
  pageId: string;
  title?: string;
  titleKey?: string;
  route: string;
}

interface BreadcrumbNavItem {
  pageRef?: string;
  label?: string;
  labelKey?: string;
  items?: BreadcrumbNavItem[];
}

export function Breadcrumbs({
  entries,
  onNavigate,
  onBack,
  showBack = false,
}: {
  entries: BreadcrumbEntry[];
  onNavigate: (route: string) => void;
  onBack: () => void;
  /** Compact circular ghost back button (semantic parent), far left. */
  showBack?: boolean;
}) {
  const t = useTranslate();
  const hasTrail = entries.length > 1;
  if (!hasTrail && !showBack) {
    return null;
  }
  const backLabel = t("shell.back");
  return (
    <nav aria-label="Breadcrumb" className="flex items-center gap-2">
      {showBack ? (
        <button
          type="button"
          onClick={onBack}
          aria-label={backLabel}
          title={backLabel}
          className="inline-flex size-6 shrink-0 items-center justify-center rounded-full text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
        >
          <ArrowLeft className="size-3.5" aria-hidden="true" />
        </button>
      ) : null}
      {hasTrail ? (
        <ol className="flex min-w-0 flex-wrap items-center gap-x-1.5 gap-y-1 text-xs font-normal text-muted-foreground">
          {entries.map((entry, index) => {
            const isLast = index === entries.length - 1;
            return (
              <li key={entry.pageId} className="flex min-w-0 items-center gap-x-1.5">
                {isLast ? (
                  <span aria-current="page" className="truncate font-normal text-foreground/90">
                    {entry.label}
                  </span>
                ) : entry.route !== "" ? (
                  <button
                    type="button"
                    onClick={() => onNavigate(entry.route)}
                    className="truncate font-normal text-muted-foreground underline-offset-4 transition-colors hover:text-foreground hover:underline"
                  >
                    {entry.label}
                  </button>
                ) : (
                  <span className="truncate font-normal text-muted-foreground/75">{entry.label}</span>
                )}
                {!isLast ? (
                  <span aria-hidden="true" className="text-muted-foreground/40">/</span>
                ) : null}
              </li>
            );
          })}
        </ol>
      ) : null}
    </nav>
  );
}

/**
 * Resolves the SEMANTIC breadcrumb trail for a matched page.
 *
 * Sources, in order:
 *   1. manifest navigation tree — a page rendered under a group shows the
 *      group labels as non-clickable ancestors (outermost first);
 *   2. declared parents (options.parents) — inner pages reached by row
 *      navigation declare their parent pageId; the chain walks up until a
 *      page with no declared parent (its nav group chain, if any, is
 *      included). Unknown declared parents fail safe (skipped).
 *   3. pages not in the tree and without a declared parent are
 *      single-level (no trail UI).
 *
 * This is NOT the visit history: the same page always shows the same
 * trail regardless of how the user got there (user ruling 2026-08-14).
 */
export function resolveBreadcrumbTrail(
  pages: BreadcrumbPage[],
  currentPage: BreadcrumbPage | undefined,
  t: (key: string, params?: MessageParams, literalFallback?: string) => string,
  options: {
    navigation?: {
      top?: BreadcrumbNavItem[];
      sidebar?: BreadcrumbNavItem[];
      user?: BreadcrumbNavItem[];
    };
    /** Consumer-declared parent pageId per inner page (web-shell level). */
    parents?: Record<string, string>;
    /**
     * Home pageId (manifest app.homePageRef — the domain-root default page,
     * not necessarily the dashboard). The trail ALWAYS starts with it:
     * 首页 => 一级页 => ... => n级内页 (user spec 2026-08-14).
     */
    homePageId?: string;
  } = {},
): BreadcrumbEntry[] {
  if (currentPage === undefined) return [];
  const labelOf = (page: BreadcrumbPage) =>
    resolveTextProp(page as unknown as Record<string, unknown>, "titleKey", "title", t) ??
    page.pageId;
  const byId = new Map(pages.map((page) => [page.pageId, page]));

  // Nav-tree group chain (outermost first) for a pageId; null when absent.
  const groupChain = (pageId: string): string[] | null => {
    const walk = (items: BreadcrumbNavItem[] | undefined, chain: string[]): string[] | null => {
      for (const item of items ?? []) {
        if (item.pageRef === pageId) return chain;
        if (item.items !== undefined) {
          const label =
            resolveTextProp(item as unknown as Record<string, unknown>, "labelKey", "label", t) ?? "";
          const next = label === "" ? chain : [...chain, label];
          const found = walk(item.items, next);
          if (found !== null) return found;
        }
      }
      return null;
    };
    for (const slot of [options.navigation?.top, options.navigation?.sidebar, options.navigation?.user]) {
      const found = walk(slot, []);
      if (found !== null) return found;
    }
    return null;
  };

  // Declared-parent chain: each parent contributes its own nav group
  // labels (outermost first) then the parent page, oldest ancestor first.
  const ancestors: BreadcrumbEntry[] = [];
  const seen = new Set<string>([currentPage.pageId]);
  let cursor = options.parents?.[currentPage.pageId];
  while (cursor !== undefined && !seen.has(cursor)) {
    const page = byId.get(cursor);
    if (page === undefined) break; // unknown declared parent — fail safe
    seen.add(cursor);
    const groups = groupChain(cursor) ?? [];
    for (const label of groups) {
      ancestors.unshift({ pageId: label, label, route: "", current: false });
    }
    ancestors.unshift({
      pageId: page.pageId,
      label: labelOf(page),
      route: page.route,
      current: false,
    });
    cursor = options.parents?.[cursor];
  }

  const ownGroups = groupChain(currentPage.pageId) ?? [];
  const chain: BreadcrumbEntry[] = [
    ...ownGroups.map((label) => ({ pageId: label, label, route: "", current: false })),
    ...ancestors,
    {
      pageId: currentPage.pageId,
      label: labelOf(currentPage),
      route: currentPage.route,
      current: true,
    },
  ];
  // 首页 root: the domain-root default page leads every trail unless the
  // current page IS home (deduplicated when the chain already contains it).
  if (options.homePageId !== undefined && currentPage.pageId !== options.homePageId) {
    const home = byId.get(options.homePageId);
    if (home !== undefined && !chain.some((entry) => entry.pageId === home.pageId)) {
      chain.unshift({
        pageId: home.pageId,
        label: labelOf(home),
        route: home.route,
        current: false,
      });
    }
  }
  return chain;
}
