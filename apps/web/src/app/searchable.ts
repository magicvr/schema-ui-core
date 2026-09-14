import {
  evaluateExpression,
  isValidExpression,
  matchRoute,
  type AppManifest,
  type NavigationContext,
  type PageEntry,
} from "@/protocol/app-manifest";
import { loadPageDocument } from "@/protocol/load-page";
import {
  evaluatePermissionTargets,
  validatePermissions,
} from "@/renderer/permissions";
import { projectNavigation, type ProjectedItem, type ProjectedLink } from "@/app/navigation";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";

/** Version of the frontend SearchableItem/provider contract. */
export const SEARCHABLE_PROVIDER_VERSION = "1" as const;

export type SearchableItemKind = "page" | "navigation" | "action";

type Translator = (key: string, params?: MessageParams, literalFallback?: string) => string;

export interface SearchableActionReference {
  /** Owner page whose validated Schema document contains the action. */
  pageId: string;
  /** The validated top-level action type, used for post-selection focus policy. */
  actionType: string;
  /** Exact toolbar/actionButton trigger passed to the page executor. */
  trigger: Record<string, unknown>;
}

/**
 * A permission-safe, already-localized item exposed by a SearchableProvider.
 * Providers own visibility decisions; the palette only matches and presents
 * their output. `id` is a global stable identity and must not be reused with a
 * different payload by another provider.
 */
export interface SearchableItem {
  id: string;
  kind: SearchableItemKind;
  label: string;
  keywords?: readonly string[];
  group?: string;
  href?: string;
  action?: SearchableActionReference;
  /** Empty-query / tie-break order inside the provider contract. */
  rank?: number;
}

export interface SearchableProviderError {
  providerId: string;
  code: string;
  message: string;
  pageId?: string;
}

export interface SearchableProviderContext {
  manifest: AppManifest;
  navigationContext: NavigationContext;
  currentPath: string;
  t: Translator;
  /** Authenticated, D-VAL validating page loader supplied by the host. */
  loadPage: (page: PageEntry) => Promise<unknown>;
  /** The App shell owns this entry; the provider never invents a shell route. */
  includeShellNotifications?: boolean;
}

export interface SearchableProviderResult {
  items: SearchableItem[];
  errors?: SearchableProviderError[];
}

export interface SearchableProvider {
  id: string;
  version: typeof SEARCHABLE_PROVIDER_VERSION;
  getItems: (
    context: SearchableProviderContext,
  ) => SearchableProviderResult | Promise<SearchableProviderResult>;
}

export interface SearchableAggregation {
  items: SearchableItem[];
  errors: SearchableProviderError[];
}

export const DEFAULT_SEARCHABLE_RESULT_LIMIT = 12;

/** Normalizes case, accents and repeated whitespace for deterministic matching. */
export function normalizeSearchText(value: string): string {
  return value
    .normalize("NFKD")
    .replace(/\p{Diacritic}/gu, "")
    .toLocaleLowerCase()
    .trim()
    .replace(/\s+/g, " ");
}

function compareStableItems(left: SearchableItem, right: SearchableItem): number {
  const leftRank = left.rank ?? 0;
  const rightRank = right.rank ?? 0;
  if (leftRank !== rightRank) {
    return leftRank - rightRank;
  }
  return left.id.localeCompare(right.id);
}

function matchScore(item: SearchableItem, query: string): number | null {
  const normalizedQuery = normalizeSearchText(query);
  if (normalizedQuery === "") {
    return 0;
  }
  const label = normalizeSearchText(item.label);
  const keywords = (item.keywords ?? []).map(normalizeSearchText).filter(Boolean);
  if (label === normalizedQuery) {
    return 400;
  }
  if (label.startsWith(normalizedQuery)) {
    return 300;
  }
  if (keywords.some((keyword) => keyword.startsWith(normalizedQuery))) {
    return 250;
  }
  if (label.includes(normalizedQuery)) {
    return 200;
  }
  if (keywords.some((keyword) => keyword.includes(normalizedQuery))) {
    return 100;
  }
  return null;
}

/** Filters, scores, sorts and caps an already aggregated provider list. */
export function filterSearchableItems(
  items: readonly SearchableItem[],
  query: string,
  limit = DEFAULT_SEARCHABLE_RESULT_LIMIT,
): SearchableItem[] {
  const safeLimit = Number.isFinite(limit) && limit > 0 ? Math.floor(limit) : 0;
  if (safeLimit === 0) {
    return [];
  }
  return items
    .map((item) => ({ item, score: matchScore(item, query) }))
    .filter((entry): entry is { item: SearchableItem; score: number } => entry.score !== null)
    .sort((left, right) => right.score - left.score || compareStableItems(left.item, right.item))
    .slice(0, safeLimit)
    .map((entry) => entry.item);
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function stringValue(value: unknown): string {
  return typeof value === "string" ? value : "";
}

function safeRoute(value: unknown): value is string {
  return (
    typeof value === "string" &&
    /^\/(?!\/)[^\s\\?#{}]*$/.test(value) &&
    !value.includes("//")
  );
}

function pageForRef(manifest: AppManifest, pageRef: string | undefined): PageEntry | undefined {
  if (pageRef === undefined || pageRef === "") {
    return undefined;
  }
  return manifest.pages.find((page) => page.pageId === pageRef);
}

function flattenProjected(
  items: ProjectedItem[],
  slot: "top" | "sidebar" | "user",
  into: Array<{ link: ProjectedLink; group?: string }>,
  group?: string,
): void {
  for (const item of items) {
    if (item.type === "group") {
      flattenProjected(item.items, slot, into, item.label);
      continue;
    }
    // Keep the slot in the stable order even though it is currently only useful
    // as a keyword. This also makes the traversal explicit for future providers.
    void slot;
    into.push({ link: item, ...(group === undefined ? {} : { group }) });
  }
}

function evaluateTriggerGate(
  value: unknown,
  context: NavigationContext,
  absentDefault: boolean,
): { ok: true; value: boolean } | { ok: false } {
  if (value === undefined || value === null) {
    return { ok: true, value: absentDefault };
  }
  if (typeof value === "boolean") {
    return { ok: true, value };
  }
  if (typeof value !== "string") {
    return { ok: false };
  }
  // The shared validator is the gate. Malformed expressions and invalid JSON
  // string escapes must never be evaluated into a visible global command.
  if (!isValidExpression(value)) {
    return { ok: false };
  }
  try {
    return { ok: true, value: evaluateExpression(value, context) };
  } catch {
    return { ok: false };
  }
}

function actionTargetId(trigger: Record<string, unknown>): string {
  const key = stringValue(trigger.key);
  if (key !== "") {
    return key;
  }
  const actionRef = stringValue(trigger.actionRef);
  if (actionRef !== "") {
    return actionRef;
  }
  return stringValue(trigger.actionId);
}

function hasLocalPermissionMetadata(trigger: Record<string, unknown>): boolean {
  if (trigger.permissionIntent !== undefined) {
    return true;
  }
  return isRecord(trigger.permissions) && Object.keys(trigger.permissions).length > 0;
}

function triggerPermissionAllowed(
  document: Record<string, unknown>,
  trigger: Record<string, unknown>,
  context: NavigationContext,
): boolean {
  const permissionErrors = validatePermissions(document);
  const targetId = actionTargetId(trigger);
  const targets = evaluatePermissionTargets(document, context).filter(
    (target) => target.targetId === targetId,
  );
  const declared = hasLocalPermissionMetadata(trigger);
  // A declared permission target that cannot be resolved is malformed and must
  // fail closed. Existing unmarked actions retain the renderer's default-allow
  // semantics, but any target found by the evaluator still has to be allowed.
  if (declared && (targetId === "" || targets.length === 0)) {
    return false;
  }
  if (targets.some((target) => !target.effectivePermission)) {
    return false;
  }
  if (permissionErrors.length > 0 && (declared || targets.length > 0)) {
    return false;
  }
  return true;
}

interface TriggerCandidate {
  trigger: Record<string, unknown>;
  actionRef: string;
  path: string;
}

function collectDirectTriggers(
  node: unknown,
  path: string,
  into: TriggerCandidate[],
): void {
  if (!isRecord(node)) {
    return;
  }
  if (node.type === "actionButton" && isRecord(node.props)) {
    const actionRef = stringValue(node.props.actionId) || stringValue(node.props.actionRef);
    if (actionRef !== "") {
      const trigger = { ...node.props };
      // The renderer's permission target falls back to node.id. Preserve that
      // identity in the programmatic trigger passed to invokeAction.
      if (stringValue(trigger.key) === "" && stringValue(node.id) !== "") {
        trigger.key = node.id;
      }
      into.push({ trigger, actionRef, path: `${path}.props` });
    }
  }
  if (node.type === "table" && isRecord(node.props) && Array.isArray(node.props.toolbar)) {
    node.props.toolbar.forEach((raw, index) => {
      if (!isRecord(raw)) {
        return;
      }
      const actionRef = stringValue(raw.actionRef) || stringValue(raw.actionId);
      if (actionRef !== "") {
        into.push({ trigger: { ...raw }, actionRef, path: `${path}.props.toolbar[${index}]` });
      }
    });
  }
  if (Array.isArray(node.children)) {
    node.children.forEach((child, index) => collectDirectTriggers(child, `${path}.children[${index}]`, into));
  }
  if (isRecord(node.props) && Array.isArray(node.props.items)) {
    node.props.items.forEach((item, index) => {
      if (isRecord(item)) {
        collectDirectTriggers(item.content, `${path}.props.items[${index}].content`, into);
      }
    });
  }
}

function actionCanBeInvokedGlobally(
  manifest: AppManifest,
  action: Record<string, unknown>,
  trigger: Record<string, unknown>,
): boolean {
  const type = stringValue(action.type);
  if (type === "upload") {
    return false;
  }
  if (type === "request" || type === "custom") {
    // Requests/custom handlers with a row placeholder would otherwise execute
    // with an empty row and are therefore not global page commands.
    return !stringValue(action.url).includes("{") && !stringValue(action.handler).includes("{");
  }
  if (type === "modal") {
    return isRecord(action.content) || stringValue(action.modalId) !== "";
  }
  if (type === "navigate") {
    const url = stringValue(action.url);
    return safeRoute(url) && matchRoute(manifest.pages, url) !== undefined;
  }
  void trigger;
  return false;
}

function manifestDestinationItems(
  context: SearchableProviderContext,
): { items: SearchableItem[]; pages: Map<string, { page: PageEntry; destination: SearchableItem }> } {
  const projection = projectNavigation(
    context.manifest,
    context.currentPath,
    context.navigationContext,
    context.t,
  );
  const links: Array<{ link: ProjectedLink; group?: string }> = [];
  flattenProjected(projection.top, "top", links);
  flattenProjected(projection.sidebar, "sidebar", links);
  flattenProjected(projection.user, "user", links);

  const items: SearchableItem[] = [];
  const pages = new Map<string, { page: PageEntry; destination: SearchableItem }>();
  const seenPages = new Set<string>();
  const seenUrls = new Set<string>();
  let rank = 0;

  for (const { link, group } of links) {
    const href = link.href;
    if (!safeRoute(href)) {
      continue;
    }
    if (link.pageRef !== undefined) {
      const page = pageForRef(context.manifest, link.pageRef);
      if (page === undefined || page.route.includes("{") || seenPages.has(page.pageId)) {
        continue;
      }
      const destination: SearchableItem = {
        id: `page:${page.pageId}`,
        kind: "page",
        label: link.label || page.title || page.pageId,
        keywords: [page.pageId, page.title ?? "", link.secondary ?? "", href],
        ...(group === undefined ? {} : { group }),
        href,
        rank: rank++,
      };
      seenPages.add(page.pageId);
      pages.set(page.pageId, { page, destination });
      items.push(destination);
      continue;
    }

    const normalizedUrl = href.split(/[?#]/, 1)[0];
    if (seenUrls.has(normalizedUrl)) {
      continue;
    }
    seenUrls.add(normalizedUrl);
    items.push({
      id: `navigation:${normalizedUrl}`,
      kind: "navigation",
      label: link.label || normalizedUrl,
      keywords: [normalizedUrl, link.secondary ?? "", group ?? ""],
      ...(group === undefined ? {} : { group }),
      href,
      rank: rank++,
    });
  }

  if (context.includeShellNotifications && context.manifest.pages.some((page) => page.pageId === "notifications")) {
    const page = pageForRef(context.manifest, "notifications");
    if (page !== undefined && !page.route.includes("{") && !seenPages.has(page.pageId)) {
      const destination: SearchableItem = {
        id: `page:${page.pageId}`,
        kind: "page",
        label: context.t("manifest.nav.notifications", undefined, page.title ?? page.pageId),
        keywords: [page.pageId, page.title ?? "", context.t("shell.notifications.bell")],
        group: context.t("shell.notifications.bell"),
        href: page.route,
        rank: rank++,
      };
      pages.set(page.pageId, { page, destination });
      items.push(destination);
    }
  }

  return { items, pages };
}

/**
 * Built-in provider. It derives destinations from the already projected,
 * permission-safe navigation and lazily reads authenticated page Schemas to
 * expose only direct page-level toolbar/actionButton commands.
 */
export function createManifestSearchProvider(): SearchableProvider {
  return {
    id: "manifest",
    version: SEARCHABLE_PROVIDER_VERSION,
    async getItems(context): Promise<SearchableProviderResult> {
      const destination = manifestDestinationItems(context);
      const items = [...destination.items];
      const errors: SearchableProviderError[] = [];
      let actionRank = 1000;

      const pages = [...destination.pages.values()];
      const loaded = await Promise.all(
        pages.map(async ({ page, destination: pageDestination }) => {
          try {
            const document = await context.loadPage(page);
            return { page, pageDestination, document };
          } catch (error) {
            errors.push({
              providerId: "manifest",
              code: "PAGE_SCHEMA_UNAVAILABLE",
              message: "Some page commands are temporarily unavailable.",
              pageId: page.pageId,
            });
            void error;
            return null;
          }
        }),
      );

      for (const entry of loaded) {
        if (entry === null || !isRecord(entry.document) || !isRecord(entry.document.body)) {
          continue;
        }
        const document = entry.document;
        const actions = isRecord(document.actions) ? document.actions : {};
        // A malformed permission tree must never make a global command visible.
        if (validatePermissions(document).length > 0) {
          continue;
        }
        const candidates: TriggerCandidate[] = [];
        collectDirectTriggers(document.body, "body", candidates);
        for (const candidate of candidates) {
          const action = actions[candidate.actionRef];
          if (!isRecord(action)) {
            continue;
          }
          if (candidate.trigger.requiresSelection === true || candidate.trigger.batchMapping !== undefined) {
            continue;
          }
          const visible = evaluateTriggerGate(candidate.trigger.visibleWhen, context.navigationContext, true);
          const enabled = evaluateTriggerGate(candidate.trigger.disabledWhen, context.navigationContext, false);
          if (!visible.ok || !visible.value || !enabled.ok || enabled.value) {
            continue;
          }
          if (!triggerPermissionAllowed(document, candidate.trigger, context.navigationContext)) {
            continue;
          }
          if (!actionCanBeInvokedGlobally(context.manifest, action, candidate.trigger)) {
            continue;
          }
          const label = resolveTextProp(
            candidate.trigger,
            "labelKey",
            "label",
            context.t,
            candidate.actionRef,
          );
          const targetId = actionTargetId(candidate.trigger) || candidate.actionRef;
          items.push({
            id: `action:${entry.page.pageId}:${targetId}`,
            kind: "action",
            label,
            keywords: [
              label,
              candidate.actionRef,
              entry.page.pageId,
              entry.pageDestination.label,
              entry.pageDestination.group ?? "",
              stringValue(action.type),
            ],
            ...(entry.pageDestination.group === undefined
              ? {}
              : { group: entry.pageDestination.group }),
            href: entry.pageDestination.href,
            action: {
              pageId: entry.page.pageId,
              actionType: stringValue(action.type),
              trigger: candidate.trigger,
            },
            rank: actionRank++,
          });
        }
      }

      return { items, ...(errors.length === 0 ? {} : { errors }) };
    },
  };
}

/**
 * Aggregates providers in deterministic provider-id order. Duplicate ids are
 * treated as a contract conflict: both candidates are removed rather than one
 * silently overriding the other.
 */
export async function aggregateSearchableItems(
  providers: readonly SearchableProvider[],
  context: SearchableProviderContext,
): Promise<SearchableAggregation> {
  const errors: SearchableProviderError[] = [];
  const accepted = new Map<string, { item: SearchableItem; providerId: string }>();
  const conflicts = new Set<string>();
  const orderedProviders = [...providers].sort((left, right) =>
    String(left.id ?? "").localeCompare(String(right.id ?? "")),
  );

  for (const provider of orderedProviders) {
    if (
      typeof provider.id !== "string" ||
      provider.id.trim() === "" ||
      provider.version !== SEARCHABLE_PROVIDER_VERSION ||
      typeof provider.getItems !== "function"
    ) {
      errors.push({ providerId: provider.id || "unknown", code: "INVALID_PROVIDER", message: "Search provider contract is invalid." });
      continue;
    }
    let result: SearchableProviderResult;
    try {
      const output = await provider.getItems(context);
      result = Array.isArray(output)
        ? { items: output as unknown as SearchableItem[] }
        : output;
    } catch (error) {
      errors.push({ providerId: provider.id, code: "PROVIDER_FAILED", message: "Search provider failed." });
      void error;
      continue;
    }
    if (!isRecord(result) || !Array.isArray(result.items)) {
      errors.push({ providerId: provider.id, code: "INVALID_PROVIDER_RESULT", message: "Search provider returned an invalid result." });
      continue;
    }
    if (Array.isArray(result.errors)) {
      for (const error of result.errors) {
        if (!isRecord(error) || typeof error.code !== "string" || typeof error.message !== "string") {
          errors.push({ providerId: provider.id, code: "INVALID_PROVIDER_ERROR", message: "Search provider returned an invalid error." });
          continue;
        }
        errors.push({
          providerId: provider.id,
          code: error.code,
          message: error.message,
          ...(typeof error.pageId === "string" ? { pageId: error.pageId } : {}),
        });
      }
    }
    for (const rawItem of result.items) {
      if (!isRecord(rawItem) || typeof rawItem.id !== "string" || rawItem.id.trim() === "" || typeof rawItem.label !== "string" || rawItem.label.trim() === "") {
        errors.push({ providerId: provider.id, code: "INVALID_ITEM", message: "Search provider returned an invalid item." });
        continue;
      }
      const item = rawItem as unknown as SearchableItem;
      if (conflicts.has(item.id)) {
        continue;
      }
      if (accepted.has(item.id)) {
        accepted.delete(item.id);
        conflicts.add(item.id);
        errors.push({ providerId: provider.id, code: "DUPLICATE_ITEM_ID", message: `Duplicate searchable item id: ${item.id}.` });
        continue;
      }
      accepted.set(item.id, { item, providerId: provider.id });
    }
  }

  return {
    items: [...accepted.values()]
      .sort((left, right) => left.providerId.localeCompare(right.providerId) || compareStableItems(left.item, right.item))
      .map(({ item }) => item),
    errors,
  };
}

/** Host helper used by App so the provider can share its shell-owned cache. */
export function createPageSchemaLoader(options: {
  schemaFetcher?: typeof fetch;
  cache?: Map<string, unknown>;
}): (page: PageEntry) => Promise<unknown> {
  return (page) => loadPageDocument(page, {}, { fetcher: options.schemaFetcher, cache: options.cache });
}
