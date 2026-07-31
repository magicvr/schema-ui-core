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

function labelFor(
  item: NavLink,
  pages: PageEntry[],
): string {
  if (item.label !== undefined) {
    return item.label;
  }
  if (item.labelKey !== undefined) {
    return item.labelKey;
  }
  if (item.pageRef !== undefined) {
    const page = pages.find((entry) => entry.pageId === item.pageRef);
    return page?.title ?? page?.titleKey ?? item.pageRef;
  }
  return "";
}

function groupLabel(item: NavGroup): string {
  return item.label ?? item.labelKey ?? "";
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
): ProjectedLink {
  const href = linkTarget(item, pages, currentPath);
  return {
    type: "link",
    ...(href === undefined ? {} : { href }),
    label: labelFor(item, pages),
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
): ProjectedItem[] {
  const projected: ProjectedItem[] = [];
  for (const item of items) {
    if (!isNavigationItemVisible(item, context)) {
      continue;
    }
    if (isGroup(item)) {
      const children = item.items
        .filter((child) => isNavigationItemVisible(child, context))
        .map((child) => projectLink(child, pages, currentPath));
      if (children.length === 0) {
        continue;
      }
      projected.push({
        type: "group",
        label: groupLabel(item),
        ...(item.icon === undefined ? {} : { icon: item.icon }),
        items: children,
      });
      continue;
    }
    projected.push(projectLink(item, pages, currentPath));
  }
  return projected;
}

export function projectNavigation(
  manifest: AppManifest,
  currentPath: string,
  context: NavigationContext = {},
): NavigationProjection {
  const navigation = manifest.navigation;
  return {
    top: projectItems(navigation?.top ?? [], manifest.pages, currentPath, context),
    sidebar: projectItems(
      navigation?.sidebar ?? [],
      manifest.pages,
      currentPath,
      context,
    ),
    user: projectItems(navigation?.user ?? [], manifest.pages, currentPath, context),
  };
}
