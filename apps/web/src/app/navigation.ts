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
  label: string;
  icon?: string;
  items: ProjectedLink[];
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

function linkActive(item: NavLink, pages: PageEntry[], currentPath: string): boolean {
  const current = stripPathQuery(currentPath);
  if (item.url !== undefined) {
    return stripPathQuery(item.url) === current;
  }
  const page = pages.find((entry) => entry.pageId === item.pageRef);
  return page === undefined ? false : matchRoute([page], current) !== undefined;
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
      projected.push({
        type: "group",
        label: groupLabel(item, t),
        ...(item.icon === undefined ? {} : { icon: item.icon }),
        items: children,
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
