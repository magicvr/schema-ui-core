export declare const DEFAULT_MANIFEST_PATH = "/.well-known/schema-ui/app-manifest.json";
/**
 * Host manifest-version support set (strict negotiation, ADR-0009): 2.7 for
 * existing production manifests, 2.8 for Host/App interoperability manifests
 * (returnIntentQueryKeys etc.), 2.9 for ADR-0039/ADR-0040 (data.route-binding /
 * form.controls.readonly). Kept additive — older manifests stay accepted.
 */
export declare const APP_MANIFEST_SUPPORTED_PROTOCOL_VERSIONS: readonly ["2.7", "2.8", "2.9"];
export declare const APP_MANIFEST_PROTOCOL_VERSION: "2.9";
export declare const MANIFEST_SOURCE_HEADER = "X-Schema-UI-Manifest-Source";
export declare const APP_MANIFEST_SOURCE = "https://github.com/magicvr/schema-ui-docs/tree/81aa1d8";
/** True when `expression` matches the frozen $context expression grammar. */
export declare function isValidExpression(expression: string): boolean;
export type ManifestErrorCode = "MANIFEST_LOAD_FAILED" | "MISSING_PROTOCOL_VERSION" | "INVALID_PROTOCOL_VERSION" | "UNSUPPORTED_PROTOCOL_VERSION" | "PROTOCOL_VERSION_TOO_LOW" | "MISSING_REQUIRED_CAPABILITY" | "CAPABILITY_REQUIRED" | "UNKNOWN_MANIFEST_FIELD" | "INVALID_MANIFEST" | "INVALID_PATH" | "MANIFEST_HOME_PAGE_UNKNOWN" | "MANIFEST_HOME_ROUTE_PARAMETRIC" | "PAGE_REF_WITH_EMPTY_PAGES" | "MISSING_PATH_BINDING" | "MANIFEST_PAGE_ID_MISMATCH" | "INVALID_LOGO_URL" | "UNKNOWN_NAV_SLOT" | "NAV_LINK_MUTEX" | "NAV_GROUP_NESTED" | "NAV_PAGE_REF_UNKNOWN" | "INVALID_RETURN_INTENT_QUERY_KEYS" | "FORBIDDEN_VARIABLE" | "SYNTAX";
export declare class ManifestError extends Error {
    readonly code: ManifestErrorCode;
    readonly path: string;
    /** Optional machine detail (e.g. the missing capability id). */
    readonly detail?: string;
    constructor(code: ManifestErrorCode, path: string, message: string, detail?: string);
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
    /** v2.8+ (ADR-0036): auth-return intent allowlist extension for this page. */
    returnIntentQueryKeys?: string[];
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
export declare function validateAppManifest(value: unknown): AppManifest;
export declare function matchRoute(pages: PageEntry[], path: string): RouteMatch | undefined;
export declare function stripPathQuery(path: string): string;
export declare function resolveInitialRoute(manifest: AppManifest, requestedPath: string): ResolvedRoute | undefined;
export declare function resolveSchemaUrl(baseURL: string, schemaUrl: string, params: Record<string, string>): string;
export declare function resolveRoutePath(route: string, params: Record<string, string>): string;
export declare function resolveLogoUrl(baseURL: string, logoUrl: string): string;
export declare function pageIdMatches(page: PageEntry, schemaPageId: string): boolean;
export declare function loadAppManifest(options?: {
    url?: string;
    fetcher?: typeof fetch;
}): Promise<AppManifest>;
/** Loads the manifest with its raw 200 bytes (bootstrap integrity, ADR-0035 D6). */
export declare function loadAppManifestBytes(options?: {
    url?: string;
    fetcher?: typeof fetch;
}): Promise<{
    manifest: AppManifest;
    bytes: Uint8Array;
}>;
export declare function evaluateExpression(expression: string, context: NavigationContext): boolean;
export declare function isNavigationItemVisible(item: NavLink | NavGroup, context: NavigationContext): boolean;
/** Normalizes a page identifier for contribution-key matching (trim + lowercase).
 *  Added in the R4 zero-conflict upgrade drill as a protocol additive sample. */
export declare function normalizePageID(id: string): string;
