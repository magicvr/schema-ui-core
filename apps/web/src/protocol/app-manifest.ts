export const DEFAULT_MANIFEST_PATH = "/.well-known/schema-ui/app-manifest.json";
export const APP_MANIFEST_PROTOCOL_VERSION = "2.7" as const;
export const APP_MANIFEST_SOURCE =
  "https://github.com/magicvr/schema-ui-docs/tree/ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b";

const APP_ID_PATTERN = /^[a-z][a-z0-9_-]*$/;
const CAPABILITY_PATTERN = /^[a-z][a-z0-9]*(?:\.[a-z][a-z0-9]*)*$/;
const ICON_PATTERN = /^[a-z][a-z0-9-]*$/;
const PATH_PATTERN = /^\/(?!\/)[^\s\\]*$/;
const TEMPLATE_NAME_PATTERN = /^[a-zA-Z_][a-zA-Z0-9_]*$/;
const EXPRESSION_PATTERN =
  /^\$context\.(user|features)\.([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)\s+(==|!=|contains)\s+(.+)$/;

export type ManifestErrorCode =
  | "MANIFEST_LOAD_FAILED"
  | "MISSING_PROTOCOL_VERSION"
  | "INVALID_PROTOCOL_VERSION"
  | "UNSUPPORTED_PROTOCOL_VERSION"
  | "PROTOCOL_VERSION_TOO_LOW"
  | "MISSING_REQUIRED_CAPABILITY"
  | "UNKNOWN_MANIFEST_FIELD"
  | "INVALID_MANIFEST"
  | "INVALID_PATH"
  | "MANIFEST_HOME_PAGE_UNKNOWN"
  | "MANIFEST_HOME_ROUTE_PARAMETRIC"
  | "PAGE_REF_WITH_EMPTY_PAGES"
  | "MISSING_PATH_BINDING"
  | "MANIFEST_PAGE_ID_MISMATCH"
  | "INVALID_LOGO_URL"
  | "UNKNOWN_NAV_SLOT"
  | "NAV_LINK_MUTEX"
  | "NAV_GROUP_NESTED"
  | "NAV_PAGE_REF_UNKNOWN"
  | "FORBIDDEN_VARIABLE"
  | "SYNTAX";

export class ManifestError extends Error {
  readonly code: ManifestErrorCode;
  readonly path: string;

  constructor(code: ManifestErrorCode, path: string, message: string) {
    super(message);
    this.name = "ManifestError";
    this.code = code;
    this.path = path;
  }
}

export interface AppInfo {
  appId: string;
  name?: string;
  nameKey?: string;
  homePageRef?: string;
  logo?: {
    light: string;
    dark?: string;
  };
  description?: string;
  descriptionKey?: string;
}

export interface PageEntry {
  pageId: string;
  title?: string;
  titleKey?: string;
  schemaUrl: string;
  route: string;
}

export interface VisibleWhen {
  when: string;
}

export interface NavPermissions {
  view?: string;
}

export interface NavLink {
  pageRef?: string;
  url?: string;
  label?: string;
  labelKey?: string;
  icon?: string;
  visibleWhen?: VisibleWhen;
  permissions?: NavPermissions;
}

export interface NavGroup {
  label?: string;
  labelKey?: string;
  icon?: string;
  items: NavLink[];
  visibleWhen?: VisibleWhen;
  permissions?: NavPermissions;
}

export type NavItem = NavLink | NavGroup;

export interface Navigation {
  top?: NavItem[];
  sidebar?: NavItem[];
  user?: NavItem[];
}

export interface AppManifest {
  protocolVersion: string;
  requiredCapabilities: string[];
  app: AppInfo;
  pages: PageEntry[];
  navigation?: Navigation;
}

export interface RouteMatch {
  page: PageEntry;
  index: number;
  params: Record<string, string>;
}

export interface ResolvedRoute extends RouteMatch {
  path: string;
  query: Record<string, string>;
  source: "home" | "deepLink";
}

export interface NavigationContext {
  user?: Record<string, unknown>;
  features?: Record<string, unknown>;
}

type JsonRecord = Record<string, unknown>;

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function fail(
  code: ManifestErrorCode,
  path: string,
  message: string,
): never {
  throw new ManifestError(code, path, message);
}

function requireRecord(value: unknown, path: string): JsonRecord {
  if (!isRecord(value)) {
    return fail("INVALID_MANIFEST", path, "Expected an object.");
  }
  return value;
}

function requireString(value: unknown, path: string, minLength = 1): string {
  if (typeof value !== "string" || value.length < minLength) {
    return fail("INVALID_MANIFEST", path, "Expected a non-empty string.");
  }
  return value;
}

function ensureKeys(value: JsonRecord, allowed: readonly string[], path: string) {
  for (const key of Object.keys(value)) {
    if (!allowed.includes(key)) {
      fail(
        "UNKNOWN_MANIFEST_FIELD",
        path === "$" ? key : `${path}.${key}`,
        `Unknown manifest field: ${key}.`,
      );
    }
  }
}

function requireArray(value: unknown, path: string): unknown[] {
  if (!Array.isArray(value)) {
    return fail("INVALID_MANIFEST", path, "Expected an array.");
  }
  return value;
}

function requireStringArray(value: unknown, path: string): string[] {
  const values = requireArray(value, path);
  const result = values.map((item, index) =>
    requireString(item, `${path}[${index}]`),
  );
  if (new Set(result).size !== result.length) {
    fail("INVALID_MANIFEST", path, "Array values must be unique.");
  }
  return result;
}

function validateRelativePath(
  value: unknown,
  path: string,
  allowPlaceholders: boolean,
): string {
  const result = requireString(value, path);
  if (
    !PATH_PATTERN.test(result) ||
    result.includes("?") ||
    result.includes("#") ||
    (!allowPlaceholders && result.includes("{"))
  ) {
    fail("INVALID_PATH", path, "Expected an application-relative path.");
  }
  return result;
}

function parseTemplate(route: string, path: string): string[] {
  if (route === "/") {
    return [];
  }
  if (route.includes("//")) {
    fail("INVALID_PATH", path, "Route templates cannot contain empty segments.");
  }

  const names: string[] = [];
  for (const [index, segment] of route.slice(1).split("/").entries()) {
    if (segment.length === 0) {
      fail(
        "INVALID_PATH",
        `${path}.route[${index}]`,
        "Route templates cannot contain empty segments.",
      );
    }
    if (segment.startsWith("{") || segment.endsWith("}")) {
      if (!/^\{[^{}]+\}$/.test(segment)) {
        fail("INVALID_PATH", path, "Invalid route placeholder.");
      }
      const name = segment.slice(1, -1);
      if (!TEMPLATE_NAME_PATTERN.test(name) || names.includes(name)) {
        fail("INVALID_PATH", path, "Invalid or repeated route placeholder.");
      }
      names.push(name);
    } else if (segment.includes("{") || segment.includes("}")) {
      fail("INVALID_PATH", path, "Invalid route placeholder.");
    }
  }
  return names;
}

function validateExpression(value: unknown, path: string): string {
  const expression = requireString(value, path);
  if (!expression.startsWith("$context.")) {
    fail("FORBIDDEN_VARIABLE", path, "Only $context variables are allowed.");
  }
  const match = EXPRESSION_PATTERN.exec(expression);
  if (!match) {
    fail("SYNTAX", path, "Invalid navigation expression.");
  }
  const literal = match[4].trim();
  if (
    !/^true$|^false$|^-?\d+(?:\.\d+)?$/.test(literal) &&
    !/^"(?:[^"\\]|\\.)*"$/.test(literal)
  ) {
    fail("SYNTAX", path, "Invalid navigation expression literal.");
  }
  return expression;
}

function parseVisibleWhen(value: unknown, path: string): VisibleWhen {
  const record = requireRecord(value, path);
  ensureKeys(record, ["when"], path);
  return { when: validateExpression(record.when, `${path}.when`) };
}

function parsePermissions(value: unknown, path: string): NavPermissions {
  const record = requireRecord(value, path);
  ensureKeys(record, ["view"], path);
  return record.view === undefined
    ? {}
    : { view: validateExpression(record.view, `${path}.view`) };
}

function parseIcon(value: unknown, path: string): string | undefined {
  if (value === undefined) {
    return undefined;
  }
  const icon = requireString(value, path);
  if (!ICON_PATTERN.test(icon)) {
    fail("INVALID_MANIFEST", path, "Invalid semantic icon name.");
  }
  return icon;
}

function parseLink(
  value: unknown,
  path: string,
  pages: PageEntry[],
): NavLink {
  const record = requireRecord(value, path);
  ensureKeys(
    record,
    ["pageRef", "url", "label", "labelKey", "icon", "visibleWhen", "permissions"],
    path,
  );
  const hasPageRef = record.pageRef !== undefined;
  const hasUrl = record.url !== undefined;
  if (hasPageRef === hasUrl) {
    fail("NAV_LINK_MUTEX", path, "A navigation link requires pageRef xor url.");
  }

  const pageRef = hasPageRef
    ? requireString(record.pageRef, `${path}.pageRef`)
    : undefined;
  if (pageRef !== undefined) {
    if (pages.length === 0) {
      fail("PAGE_REF_WITH_EMPTY_PAGES", `${path}.pageRef`, "No pages are registered.");
    }
    if (!pages.some((page) => page.pageId === pageRef)) {
      fail("NAV_PAGE_REF_UNKNOWN", `${path}.pageRef`, `Unknown pageRef: ${pageRef}.`);
    }
  }

  const url = hasUrl
    ? validateRelativePath(record.url, `${path}.url`, false)
    : undefined;
  if (url !== undefined && record.label === undefined && record.labelKey === undefined) {
    fail("INVALID_MANIFEST", path, "URL navigation links require a label.");
  }

  const label =
    record.label === undefined
      ? undefined
      : requireString(record.label, `${path}.label`);
  const labelKey =
    record.labelKey === undefined
      ? undefined
      : requireString(record.labelKey, `${path}.labelKey`);
  const visibleWhen =
    record.visibleWhen === undefined
      ? undefined
      : parseVisibleWhen(record.visibleWhen, `${path}.visibleWhen`);
  const permissions =
    record.permissions === undefined
      ? undefined
      : parsePermissions(record.permissions, `${path}.permissions`);

  const icon = parseIcon(record.icon, `${path}.icon`);
  return {
    ...(pageRef === undefined ? {} : { pageRef }),
    ...(url === undefined ? {} : { url }),
    ...(label === undefined ? {} : { label }),
    ...(labelKey === undefined ? {} : { labelKey }),
    ...(icon === undefined ? {} : { icon }),
    ...(visibleWhen === undefined ? {} : { visibleWhen }),
    ...(permissions === undefined ? {} : { permissions }),
  };
}

function parseItem(value: unknown, path: string, pages: PageEntry[]): NavItem {
  const record = requireRecord(value, path);
  if (record.items !== undefined) {
    ensureKeys(
      record,
      ["label", "labelKey", "icon", "items", "visibleWhen", "permissions"],
      path,
    );
    const label =
      record.label === undefined
        ? undefined
        : requireString(record.label, `${path}.label`);
    const labelKey =
      record.labelKey === undefined
        ? undefined
        : requireString(record.labelKey, `${path}.labelKey`);
    if (label === undefined && labelKey === undefined) {
      fail("INVALID_MANIFEST", path, "Navigation groups require a label.");
    }
    const items = requireArray(record.items, `${path}.items`);
    if (items.length === 0) {
      fail("INVALID_MANIFEST", `${path}.items`, "Navigation groups cannot be empty.");
    }
    const visibleWhen =
      record.visibleWhen === undefined
        ? undefined
        : parseVisibleWhen(record.visibleWhen, `${path}.visibleWhen`);
    const permissions =
      record.permissions === undefined
        ? undefined
        : parsePermissions(record.permissions, `${path}.permissions`);
    const icon = parseIcon(record.icon, `${path}.icon`);
    return {
      ...(label === undefined ? {} : { label }),
      ...(labelKey === undefined ? {} : { labelKey }),
      ...(icon === undefined ? {} : { icon }),
      items: items.map((item, index) => {
        const childPath = `${path}.items[${index}]`;
        const child = requireRecord(item, childPath);
        if (child.items !== undefined) {
          fail("NAV_GROUP_NESTED", childPath, "Navigation groups cannot be nested.");
        }
        return parseLink(item, childPath, pages);
      }),
      ...(visibleWhen === undefined ? {} : { visibleWhen }),
      ...(permissions === undefined ? {} : { permissions }),
    };
  }
  return parseLink(value, path, pages);
}

function parsePages(value: unknown): PageEntry[] {
  const entries = requireArray(value, "pages");
  const pageIds = new Set<string>();
  const routes = new Set<string>();
  return entries.map((entry, index) => {
    const path = `pages[${index}]`;
    const record = requireRecord(entry, path);
    ensureKeys(record, ["pageId", "title", "titleKey", "schemaUrl", "route"], path);
    const pageId = requireString(record.pageId, `${path}.pageId`);
    if (pageIds.has(pageId)) {
      fail("INVALID_MANIFEST", `${path}.pageId`, "pageId values must be unique.");
    }
    pageIds.add(pageId);
    const title =
      record.title === undefined ? undefined : requireString(record.title, `${path}.title`);
    const titleKey =
      record.titleKey === undefined
        ? undefined
        : requireString(record.titleKey, `${path}.titleKey`);
    if (title === undefined && titleKey === undefined) {
      fail("INVALID_MANIFEST", path, "Pages require a title or titleKey.");
    }
    const schemaUrl = validateRelativePath(record.schemaUrl, `${path}.schemaUrl`, true);
    const route = validateRelativePath(record.route, `${path}.route`, true);
    const routeNames = parseTemplate(route, `${path}.route`);
    const schemaNames = parseTemplate(schemaUrl, `${path}.schemaUrl`);
    if (schemaNames.some((name) => !routeNames.includes(name))) {
      fail(
        "INVALID_PATH",
        `${path}.schemaUrl`,
        "schemaUrl placeholders must be bound by the page route.",
      );
    }
    if (routes.has(route)) {
      fail("INVALID_MANIFEST", `${path}.route`, "Route templates must be unique.");
    }
    routes.add(route);
    return {
      pageId,
      ...(title === undefined ? {} : { title }),
      ...(titleKey === undefined ? {} : { titleKey }),
      schemaUrl,
      route,
    };
  });
}

function parseApp(value: unknown, pages: PageEntry[]): AppInfo {
  const record = requireRecord(value, "app");
  ensureKeys(
    record,
    ["appId", "name", "nameKey", "homePageRef", "logo", "description", "descriptionKey"],
    "app",
  );
  const appId = requireString(record.appId, "app.appId");
  if (!APP_ID_PATTERN.test(appId)) {
    fail("INVALID_MANIFEST", "app.appId", "Invalid appId.");
  }
  const name = record.name === undefined ? undefined : requireString(record.name, "app.name");
  const nameKey =
    record.nameKey === undefined ? undefined : requireString(record.nameKey, "app.nameKey");
  if (name === undefined && nameKey === undefined) {
    fail("INVALID_MANIFEST", "app", "App requires a name or nameKey.");
  }
  if (pages.length > 0 && record.homePageRef === undefined) {
    fail("INVALID_MANIFEST", "app.homePageRef", "Non-empty pages require homePageRef.");
  }
  if (pages.length === 0 && record.homePageRef !== undefined) {
    fail(
      "PAGE_REF_WITH_EMPTY_PAGES",
      "app.homePageRef",
      "An empty page registry cannot declare a homePageRef.",
    );
  }
  const homePageRef =
    record.homePageRef === undefined
      ? undefined
      : requireString(record.homePageRef, "app.homePageRef");
  if (homePageRef !== undefined) {
    const home = pages.find((page) => page.pageId === homePageRef);
    if (!home) {
      fail("MANIFEST_HOME_PAGE_UNKNOWN", "app.homePageRef", "Unknown home page.");
    }
    if (home.route.includes("{")) {
      fail(
        "MANIFEST_HOME_ROUTE_PARAMETRIC",
        "app.homePageRef",
        "Home page route cannot be parametric.",
      );
    }
  }

  let logo: AppInfo["logo"];
  if (record.logo !== undefined) {
    const logoRecord = requireRecord(record.logo, "app.logo");
    ensureKeys(logoRecord, ["light", "dark"], "app.logo");
    const light = validateLogoUrl(logoRecord.light, "app.logo.light");
    const dark =
      logoRecord.dark === undefined
        ? undefined
        : validateLogoUrl(logoRecord.dark, "app.logo.dark");
    logo = dark === undefined ? { light } : { light, dark };
  }

  const description =
    record.description === undefined
      ? undefined
      : requireString(record.description, "app.description", 0);
  const descriptionKey =
    record.descriptionKey === undefined
      ? undefined
      : requireString(record.descriptionKey, "app.descriptionKey", 0);
  return {
    appId,
    ...(name === undefined ? {} : { name }),
    ...(nameKey === undefined ? {} : { nameKey }),
    ...(homePageRef === undefined ? {} : { homePageRef }),
    ...(logo === undefined ? {} : { logo }),
    ...(description === undefined ? {} : { description }),
    ...(descriptionKey === undefined ? {} : { descriptionKey }),
  };
}

function parseNavigation(value: unknown, pages: PageEntry[]): Navigation {
  const record = requireRecord(value, "navigation");
  const navigation: Navigation = {};
  for (const key of Object.keys(record)) {
    if (key !== "top" && key !== "sidebar" && key !== "user") {
      fail("UNKNOWN_NAV_SLOT", `navigation.${key}`, `Unknown navigation slot: ${key}.`);
    }
  }
  for (const slot of ["top", "sidebar", "user"] as const) {
    if (record[slot] !== undefined) {
      const items = requireArray(record[slot], `navigation.${slot}`);
      navigation[slot] = items.map((item, index) =>
        parseItem(item, `navigation.${slot}[${index}]`, pages),
      );
    }
  }
  return navigation;
}

export function validateAppManifest(value: unknown): AppManifest {
  const record = requireRecord(value, "$" );
  ensureKeys(record, ["protocolVersion", "requiredCapabilities", "app", "pages", "navigation"], "$" );

  if (record.protocolVersion === undefined) {
    fail("MISSING_PROTOCOL_VERSION", "protocolVersion", "Manifest protocolVersion is required.");
  }
  const protocolVersion = requireString(record.protocolVersion, "protocolVersion");
  if (!/^\d+\.\d+$/.test(protocolVersion)) {
    fail("INVALID_PROTOCOL_VERSION", "protocolVersion", "Expected MAJOR.MINOR.");
  }
  const [majorText, minorText] = protocolVersion.split(".");
  const major = Number(majorText);
  const minor = Number(minorText);
  if (major < 2 || (major === 2 && minor < 5)) {
    fail("PROTOCOL_VERSION_TOO_LOW", "protocolVersion", "App manifest requires protocol >= 2.5.");
  }
  if (protocolVersion !== APP_MANIFEST_PROTOCOL_VERSION) {
    fail(
      "UNSUPPORTED_PROTOCOL_VERSION",
      "protocolVersion",
      `This host supports ${APP_MANIFEST_PROTOCOL_VERSION}.`,
    );
  }

  const requiredCapabilities = requireStringArray(
    record.requiredCapabilities,
    "requiredCapabilities",
  );
  if (requiredCapabilities.some((capability) => !CAPABILITY_PATTERN.test(capability))) {
    fail(
      "INVALID_MANIFEST",
      "requiredCapabilities",
      "Capabilities must use dotted lowercase names.",
    );
  }
  if (!requiredCapabilities.includes("app.manifest")) {
    fail(
      "MISSING_REQUIRED_CAPABILITY",
      "requiredCapabilities",
      "app.manifest is required.",
    );
  }

  const pages = parsePages(record.pages);
  const app = parseApp(record.app, pages);
  let navigation: Navigation | undefined;
  if (record.navigation !== undefined) {
    if (!requiredCapabilities.includes("app.navigation")) {
      fail(
        "MISSING_REQUIRED_CAPABILITY",
        "requiredCapabilities",
        "app.navigation is required when navigation is present.",
      );
    }
    navigation = parseNavigation(record.navigation, pages);
  }

  return {
    protocolVersion,
    requiredCapabilities,
    app,
    pages,
    ...(navigation === undefined ? {} : { navigation }),
  };
}

function validateLogoUrl(value: unknown, path: string): string {
  const logo = requireString(value, path);
  if (
    !(
      (/^\/(?!\/)[^\s\\{}]*$/.test(logo) && !logo.includes("?") && !logo.includes("#")) ||
      /^https:\/\/[^\s\\]+$/.test(logo)
    )
  ) {
    fail("INVALID_LOGO_URL", path, "Logo must be a site-relative or https URL.");
  }
  return logo;
}

function decodeSegment(value: string): string | undefined {
  try {
    return decodeURIComponent(value);
  } catch {
    return undefined;
  }
}

function routeSegments(path: string): string[] | undefined {
  const cleanPath = stripPathQuery(path);
  if (!cleanPath.startsWith("/") || cleanPath.includes("//")) {
    return undefined;
  }
  if (cleanPath === "/") {
    return [];
  }
  const segments = cleanPath.slice(1).split("/");
  return segments.some((segment) => segment === "") ? undefined : segments;
}

function matchSinglePage(page: PageEntry, path: string): Record<string, string> | undefined {
  const inputSegments = routeSegments(path);
  const templateSegments = routeSegments(page.route);
  if (!inputSegments || !templateSegments || inputSegments.length !== templateSegments.length) {
    return undefined;
  }

  const params: Record<string, string> = {};
  for (const [index, templateSegment] of templateSegments.entries()) {
    const decoded = decodeSegment(inputSegments[index]);
    const expected = decodeSegment(templateSegment);
    if (decoded === undefined || expected === undefined) {
      return undefined;
    }
    if (/^\{[^{}]+\}$/.test(templateSegment)) {
      params[templateSegment.slice(1, -1)] = decoded;
    } else if (decoded !== expected) {
      return undefined;
    }
  }
  return params;
}

export function matchRoute(pages: PageEntry[], path: string): RouteMatch | undefined {
  const inputSegments = routeSegments(path);
  if (!inputSegments) {
    return undefined;
  }
  const candidates = pages.flatMap((page, index) => {
    const params = matchSinglePage(page, path);
    if (!params) {
      return [];
    }
    const templateSegments = routeSegments(page.route) ?? [];
    const literalCount = templateSegments.filter(
      (segment) => !/^\{[^{}]+\}$/.test(segment),
    ).length;
    return [{ page, index, params, literalCount }];
  });
  candidates.sort((left, right) => {
    if (left.literalCount !== right.literalCount) {
      return right.literalCount - left.literalCount;
    }
    if (left.page.route.length !== right.page.route.length) {
      return right.page.route.length - left.page.route.length;
    }
    return left.index - right.index;
  });
  const winner = candidates[0];
  return winner === undefined
    ? undefined
    : { page: winner.page, index: winner.index, params: winner.params };
}

export function stripPathQuery(path: string): string {
  const queryIndex = path.search(/[?#]/);
  return queryIndex === -1 ? path : path.slice(0, queryIndex);
}

function parseQuery(path: string): Record<string, string> {
  const queryIndex = path.indexOf("?");
  if (queryIndex === -1) {
    return {};
  }
  const query = new URLSearchParams(path.slice(queryIndex + 1).split("#", 1)[0]);
  return Object.fromEntries(query.entries());
}

export function resolveInitialRoute(
  manifest: AppManifest,
  requestedPath: string,
): ResolvedRoute | undefined {
  const requested = stripPathQuery(requestedPath);
  if (requested !== "/") {
    const deepLink = matchRoute(manifest.pages, requested);
    if (deepLink) {
      return {
        ...deepLink,
        path: requested,
        query: parseQuery(requestedPath),
        source: "deepLink",
      };
    }
    return undefined;
  }
  const homePage = manifest.pages.find((page) => page.pageId === manifest.app.homePageRef);
  if (!homePage) {
    return undefined;
  }
  const home = matchRoute(manifest.pages, homePage.route);
  if (!home) {
    return undefined;
  }
  return { ...home, path: homePage.route, query: {}, source: "home" };
}

export function resolveSchemaUrl(
  baseURL: string,
  schemaUrl: string,
  params: Record<string, string>,
): string {
  const resolved = resolveTemplate(schemaUrl, params, "schemaUrl");
  return joinBaseURL(baseURL, resolved);
}

export function resolveRoutePath(
  route: string,
  params: Record<string, string>,
): string {
  return resolveTemplate(route, params, "route");
}

export function resolveLogoUrl(baseURL: string, logoUrl: string): string {
  if (/^https:\/\//.test(logoUrl)) {
    return logoUrl;
  }
  if (!/^\/(?!\/)[^\s\\{}]*$/.test(logoUrl)) {
    fail("INVALID_LOGO_URL", "logoUrl", "Logo must be a site-relative or https URL.");
  }
  return joinBaseURL(baseURL, logoUrl);
}

function resolveTemplate(
  template: string,
  params: Record<string, string>,
  field: string,
): string {
  return template.replace(/\{([a-zA-Z_][a-zA-Z0-9_]*)\}/g, (_match, name: string) => {
    const value = params[name];
    if (value === undefined) {
      fail("MISSING_PATH_BINDING", field + ".{" + name + "}", "Missing path binding: " + name + ".");
    }
    return encodeURIComponent(value);
  });
}

function joinBaseURL(baseURL: string, relativePath: string): string {
  const base = baseURL.endsWith("/") ? baseURL : `${baseURL}/`;
  return new URL(relativePath.replace(/^\/+/, ""), base).toString();
}

export function pageIdMatches(page: PageEntry, schemaPageId: string): boolean {
  if (page.pageId !== schemaPageId) {
    throw new ManifestError(
      "MANIFEST_PAGE_ID_MISMATCH",
      `pages[${page.pageId}].pageId`,
      "The page schema pageId does not match the manifest pageId.",
    );
  }
  return true;
}

export async function loadAppManifest(options: {
  url?: string;
  fetcher?: typeof fetch;
} = {}): Promise<AppManifest> {
  const url = options.url ?? DEFAULT_MANIFEST_PATH;
  const fetcher = options.fetcher ?? globalThis.fetch;
  if (!fetcher) {
    throw new ManifestError("MANIFEST_LOAD_FAILED", url, "Fetch is unavailable.");
  }
  try {
    const response = await fetcher(url);
    if (!response.ok) {
      throw new ManifestError(
        "MANIFEST_LOAD_FAILED",
        url,
        `Manifest request failed with HTTP ${response.status}.`,
      );
    }
    const payload: unknown = await response.json();
    return validateAppManifest(payload);
  } catch (error) {
    if (error instanceof ManifestError) {
      throw error;
    }
    throw new ManifestError(
      "MANIFEST_LOAD_FAILED",
      url,
      "Manifest could not be fetched or parsed.",
    );
  }
}

function getContextValue(
  context: NavigationContext,
  root: "user" | "features",
  path: string,
): unknown {
  let current: unknown = context[root];
  for (const part of path.split(".")) {
    if (!isRecord(current)) {
      return undefined;
    }
    current = current[part];
  }
  return current;
}

function parseLiteral(value: string): unknown {
  if (value === "true") return true;
  if (value === "false") return false;
  if (/^-?\d+(?:\.\d+)?$/.test(value)) return Number(value);
  return JSON.parse(value) as unknown;
}

export function evaluateExpression(
  expression: string,
  context: NavigationContext,
): boolean {
  const match = EXPRESSION_PATTERN.exec(expression);
  if (!match) {
    return false;
  }
  const actual = getContextValue(context, match[1] as "user" | "features", match[2]);
  const expected = parseLiteral(match[4].trim());
  switch (match[3]) {
    case "contains":
      return Array.isArray(actual)
        ? actual.includes(expected)
        : typeof actual === "string" && typeof expected === "string"
          ? actual.includes(expected)
          : false;
    case "==":
      return Object.is(actual, expected);
    case "!=":
      return !Object.is(actual, expected);
    default:
      return false;
  }
}

export function isNavigationItemVisible(
  item: NavLink | NavGroup,
  context: NavigationContext,
): boolean {
  const permission = item.permissions?.view;
  const condition = item.visibleWhen?.when;
  return (
    (permission === undefined || evaluateExpression(permission, context)) &&
    (condition === undefined || evaluateExpression(condition, context))
  );
}
