/**
 * Breadcrumb navigation for nested admin pages (GOAL-015 D-002 §3.6).
 *
 * Renders a trail of ancestor pages (Home > ... > current) plus a Back button
 * that returns to the previous page in the session's navigation stack. Pure
 * shell-layer UI: the trail is derived from the route stack maintained by the
 * app shell (no schema/protocol dependency), and the back target is the
 * browser history — the route-stack approach confirmed by the user 2026-08-14.
 */

import { useTranslate } from "@/i18n/runtime";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";

export interface BreadcrumbEntry {
  /** Page id (manifest pageId) — used as the React key. */
  pageId: string;
  label: string;
  /** Application route of this ancestor. */
  route: string;
  /** True for the current (deepest) page. */
  current: boolean;
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
  /** Show the back button even for single-level pages (route-stack nav). */
  showBack?: boolean;
}) {
  const t = useTranslate();
  const hasTrail = entries.length > 1;
  if (!hasTrail && !showBack) {
    return null;
  }
  return (
    <nav aria-label="Breadcrumb" className="flex flex-wrap items-center gap-1 text-sm text-muted-foreground">
      {showBack ? (
        <button
          type="button"
          onClick={onBack}
          className="inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        >
          <span aria-hidden="true">←</span>
          <span>{t("shell.back")}</span>
        </button>
      ) : null}
      {hasTrail ? <span aria-hidden="true" className="text-muted-foreground/50">|</span> : null}
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
 * resolveBreadcrumbTrail computes the ancestor chain for a matched page from
 * the route stack (paths visited before the current one).
 *
 * The trail is the sequence of manifest pages that lead to the current page,
 * derived from the visit stack maintained by the shell: the most recent
 * distinct ancestor pages (up to maxAncestors) plus the current page. This is
 * the route-stack approach — no breadcrumbParent protocol key required
 * (P-3 only covers form-field readonly/disabled, not manifest pages).
 */
export function resolveBreadcrumbTrail(
  pages: Array<{ pageId: string; title?: string; titleKey?: string; route: string }>,
  currentPage: { pageId: string; title?: string; titleKey?: string; route: string } | undefined,
  t: (key: string, params?: MessageParams, literalFallback?: string) => string,
  visitStack: string[] = [],
  maxAncestors = 2,
): BreadcrumbEntry[] {
  if (currentPage === undefined) return [];
  const labelOf = (page: { pageId: string; title?: string; titleKey?: string; route: string }) =>
    resolveTextProp(page as unknown as Record<string, unknown>, "titleKey", "title", t) ??
    page.pageId;
  const currentLabel = labelOf(currentPage);
  // Deduplicate the stack, drop the current path, take the most recent
  // ancestors (reverse chronological), cap at maxAncestors, then re-order
  // oldest-first so the trail reads Home > ... > current.
  const seen = new Set<string>();
  const ancestors: BreadcrumbEntry[] = [];
  for (let i = visitStack.length - 1; i >= 0; i--) {
    const path = visitStack[i];
    if (path === currentPage.route || seen.has(path)) continue;
    const page = pages.find((p) => p.route === path);
    if (page === undefined) continue;
    seen.add(path);
    ancestors.push({
      pageId: page.pageId,
      label: labelOf(page),
      route: page.route,
      current: false,
    });
    if (ancestors.length >= maxAncestors) break;
  }
  ancestors.reverse();
  ancestors.push({
    pageId: currentPage.pageId,
    label: currentLabel,
    route: currentPage.route,
    current: true,
  });
  return ancestors;
}
