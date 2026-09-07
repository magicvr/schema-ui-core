import {
  type AppManifest,
  type NavGroup,
  type NavItem,
  type NavLink,
  type NavigationContext,
  type PageEntry,
  isNavigationItemVisible,
  matchRoute,
  resolveRoutePath,
  stripPathQuery,
} from "@/protocol/app-manifest";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";

/**
 * Shell-level page hierarchy used by both breadcrumbs and navigation active
 * state. It is intentionally outside the protocol: inner/detail pages can
 * inherit the active state of their registered sidebar parent without adding
 * a non-standard field to NavGroup.
 */
export const NAVIGATION_PAGE_PARENTS: Record<string, string> = {
  "dictionary-entries": "data-dictionary",
  "task-runs": "scheduled-tasks",
  "wallet-entries": "wallet",
  "users-invites": "users",
  "telegram-operator": "telegram-settings",
};

export interface ProjectedLink {
  type: "link";
  href?: string;
  label: string;
  pageRef?: string;
  url?: string;
  icon?: string;
  active: boolean;
}

export interface ProjectedGroup {
  type: "group";
  /** Stable UI key; derived from labelKey/literal because protocol NavGroup has no id. */
  key: string;
  label: string;
  icon?: string;
  items: ProjectedLink[];
  active: boolean;
}

export type ProjectedItem = ProjectedLink | ProjectedGroup;

export interface NavigationProjection {
  top: ProjectedItem[];
  sidebar: ProjectedItem[];
  user: ProjectedItem[];
}

function isGroup(item: NavItem): item is NavGroup {
  return "items" in item;
}

/** Translator used for labelKey/titleKey resolution; defaults to identity. */
type Translator = (key: string, params?: MessageParams, literalFallback?: string) => string;

const identityTranslator: Translator = (_key, _params, literalFallback) =>
  literalFallback ?? "";

function labelFor(
  item: NavLink,
  pages: PageEntry[],
  t: Translator,
): string {
  if (item.labelKey !== undefined) {
    return t(item.labelKey, undefined, item.label);
  }
  if (item.label !== undefined) {
    return item.label;
  }
  if (item.pageRef !== undefined) {
    const page = pages.find((entry) => entry.pageId === item.pageRef);
    if (page !== undefined) {
      return resolveTextProp(page as unknown as Record<string, unknown>, "titleKey", "title", t, item.pageRef);
    }
    return item.pageRef;
  }
  return "";
}

function groupLabel(item: NavGroup, t: Translator): string {
  return resolveTextProp(item as unknown as Record<string, unknown>, "labelKey", "label", t);
}

function linkTarget(
  item: NavLink,
  pages: PageEntry[],
  currentPath: string,
): string | undefined {
  if (item.url !== undefined) {
    return item.url;
  }
  const page = pages.find((entry) => entry.pageId === item.pageRef);
  if (page === undefined || !page.route.includes("{")) {
    return page?.route;
  }
  const current = matchRoute([page], currentPath);
  return current === undefined ? undefined : resolveRoutePath(page.route, current.params);
}

function isPageOrAncestorActive(targetPageId: string, currentPageId: string): boolean {
  const visited = new Set<string>();
  let candidate: string | undefined = currentPageId;
  while (candidate !== undefined && !visited.has(candidate)) {
    if (candidate === targetPageId) {
      return true;
    }
    visited.add(candidate);
    candidate = NAVIGATION_PAGE_PARENTS[candidate];
  }
  return false;
}

function linkActive(item: NavLink, pages: PageEntry[], currentPath: string): boolean {
  const current = stripPathQuery(currentPath);
  if (item.url !== undefined) {
    return stripPathQuery(item.url) === current;
  }
  const page = pages.find((entry) => entry.pageId === item.pageRef);
  if (page === undefined) {
    return false;
  }
  const currentPage = matchRoute(pages, current)?.page;
  return currentPage !== undefined && isPageOrAncestorActive(page.pageId, currentPage.pageId);
}

function projectLink(
  item: NavLink,
  pages: PageEntry[],
  currentPath: string,
  t: Translator,
): ProjectedLink {
  const href = linkTarget(item, pages, currentPath);
  return {
    type: "link",
    ...(href === undefined ? {} : { href }),
    label: labelFor(item, pages, t),
    ...(item.pageRef === undefined ? {} : { pageRef: item.pageRef }),
    ...(item.url === undefined ? {} : { url: item.url }),
    ...(item.icon === undefined ? {} : { icon: item.icon }),
    active: linkActive(item, pages, currentPath),
  };
}

function projectItems(
  items: NavItem[],
  pages: PageEntry[],
  currentPath: string,
  context: NavigationContext,
  t: Translator,
): ProjectedItem[] {
  const projected: ProjectedItem[] = [];
  for (const item of items) {
    if (!isNavigationItemVisible(item, context)) {
      continue;
    }
    if (isGroup(item)) {
      const children = item.items
        .filter((child) => isNavigationItemVisible(child, context))
        .map((child) => projectLink(child, pages, currentPath, t));
      if (children.length === 0) {
        continue;
      }
      const label = groupLabel(item, t);
      projected.push({
        type: "group",
        key: item.labelKey ?? item.label ?? label,
        label,
        ...(item.icon === undefined ? {} : { icon: item.icon }),
        items: children,
        active: children.some((child) => child.active),
      });
      continue;
    }
    projected.push(projectLink(item, pages, currentPath, t));
  }
  return projected;
}

export function projectNavigation(
  manifest: AppManifest,
  currentPath: string,
  context: NavigationContext = {},
  t: Translator = identityTranslator,
): NavigationProjection {
  const navigation = manifest.navigation;
  return {
    top: projectItems(navigation?.top ?? [], manifest.pages, currentPath, context, t),
    sidebar: projectItems(
      navigation?.sidebar ?? [],
      manifest.pages,
      currentPath,
      context,
      t,
    ),
    user: projectItems(navigation?.user ?? [], manifest.pages, currentPath, context, t),
  };
}
