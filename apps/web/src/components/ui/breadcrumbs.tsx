/**
 * Breadcrumb navigation for nested admin pages (GOAL-015 D-002 §3.6).
 *
 * Renders a trail of ancestor pages (Home > ... > current) plus a Back button
 * that returns to the previous page in the session's navigation stack. Pure
 * shell-layer UI: the trail derives from the manifest page tree (matchRoute)
 * and the back target from window.history, so no schema/protocol change is
 * required (route-stack approach, user-confirmed 2026-08-14).
 */

import { useTranslate } from "@/i18n/runtime";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";

export interface BreadcrumbEntry {
  /** Page id (manifest pageId) — used as the React key. */
  pageId: string;
  label: string;
  /** Application route of this ancestor; empty string for the current page. */
  route: string;
  /** True for the current (deepest) page. */
  current: boolean;
}

export function Breadcrumbs({
  entries,
  onNavigate,
  onBack,
}: {
  entries: BreadcrumbEntry[];
  onNavigate: (route: string) => void;
  onBack: () => void;
}) {
  const t = useTranslate();
  if (entries.length <= 1) {
    // No ancestors: a single-level page shows no trail (only the Back button
    // when history has a previous entry is handled by the caller).
    return null;
  }
  return (
    <nav aria-label="Breadcrumb" className="flex flex-wrap items-center gap-1 text-sm text-muted-foreground">
      <button
        type="button"
        onClick={onBack}
        className="inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
      >
        <span aria-hidden="true">←</span>
        <span>{t("shell.back")}</span>
      </button>
      <span aria-hidden="true" className="text-muted-foreground/50">|</span>
      {entries.map((entry, index) => {
        const isLast = index === entries.length - 1;
        return (
          <span key={entry.pageId} className="flex min-w-0 items-center gap-1">
            {isLast ? (
              <span aria-current="page" className="truncate font-medium text-foreground">
                {entry.label}
              </span>
            ) : entry.route !== "" ? (
              <button
                type="button"
                onClick={() => onNavigate(entry.route)}
                className="truncate transition-colors hover:text-foreground"
              >
                {entry.label}
              </button>
            ) : (
              <span className="truncate">{entry.label}</span>
            )}
            {!isLast ? <span aria-hidden="true">/</span> : null}
          </span>
        );
      })}
    </nav>
  );
}

/**
 * resolveBreadcrumbTrail computes the ancestor chain for a matched page.
 *
 * The trail is the sequence of manifest pages that lead to the current page.
 * With a flat page registry this is [current] (single level). Nested inner
 * pages (GOAL-015) extend this by declaring their parent page id via the
 * page's breadcrumbParent key (shell-layer convention; the manifest stays
 * the source of truth for labels/routes).
 */
export function resolveBreadcrumbTrail(
  pages: Array<{ pageId: string; title?: string; titleKey?: string; route: string; breadcrumbParent?: string }>,
  currentPage: { pageId: string; title?: string; titleKey?: string; route: string; breadcrumbParent?: string } | undefined,
  t: (key: string, params?: MessageParams, literalFallback?: string) => string,
): BreadcrumbEntry[] {
  if (currentPage === undefined) return [];
  const byId = new Map(pages.map((p) => [p.pageId, p]));
  const trail: BreadcrumbEntry[] = [];
  const visited = new Set<string>();
  let cursor = currentPage;
  while (cursor !== undefined && !visited.has(cursor.pageId)) {
    visited.add(cursor.pageId);
    const label =
      resolveTextProp(cursor as unknown as Record<string, unknown>, "titleKey", "title", t) ??
      cursor.pageId;
    const parentId = cursor.breadcrumbParent;
    trail.unshift({
      pageId: cursor.pageId,
      label,
      route: cursor.route,
      current: true, // corrected below
    });
    if (parentId !== undefined && parentId !== "") {
      cursor = byId.get(parentId);
    } else {
      cursor = undefined;
    }
  }
  // Mark only the deepest entry as current.
  trail.forEach((entry, index) => {
    entry.current = index === trail.length - 1;
  });
  return trail;
}
