import { evaluateExpression, type NavigationContext } from "@/protocol/app-manifest";

/**
 * R4 D-PERM permission evaluation engine (ADR-0023, frozen by GOAL-006 D-004).
 *
 * Coverage scope: the minimal renderer subset exercised by the frozen
 * permissions-inheritance fixtures — permissionCascade / permissionIntent
 * validation (fail-closed error codes) and effectivePermission evaluation for
 * the approved target kinds (formField / formSubmit / rowAction /
 * toolbarTrigger / actionButton). Execution-time gating (visibleWhen →
 * permission → disabled/requiresSelection → confirm → action) is modeled for
 * rowAction / toolbarTrigger / actionButton / default form submit.
 */

export type PermissionKey = "view" | "edit" | "delete";

const CASCADE_TYPES = new Set(["section", "grid", "form", "tabs", "table"]);
const CASCADE_KEYS = new Set(["edit", "delete"]);
const INTENT_KEYS = new Set(["edit", "delete"]);
const FORM_EDIT_FIELD_TYPES = new Set([
  "input",
  "inputNumber",
  "datePicker",
  "dateRangePicker",
  "select",
  "upload",
]);
const PERMISSION_CAPABILITY = "permissions.inheritance";

export type L2ErrorCode =
  | "PROTOCOL_VERSION_TOO_LOW"
  | "CAPABILITY_REQUIRED"
  | "PERMISSION_CASCADE_TYPE_INVALID"
  | "PERMISSION_CASCADE_KEYS_INVALID"
  | "PERMISSION_CASCADE_SOURCE_MISSING"
  | "PERMISSION_INTENT_FORBIDDEN"
  | "PERMISSION_INTENT_INVALID";

export interface L2Error {
  code: L2ErrorCode;
  path: string;
}

export interface PermissionTarget {
  targetId: string;
  kind: "formField" | "formSubmit" | "rowAction" | "toolbarTrigger" | "actionButton" | "column";
  key: PermissionKey;
  cascadeApplied: boolean;
  cascadedBy: string[];
  effectivePermission: boolean;
}

export interface ExecutionRequest {
  targetId: string;
  visible: boolean;
  confirm?: boolean;
  confirmed?: boolean;
  disabled?: boolean;
  requiresSelection?: boolean;
}

export interface ExecutionResult {
  outcome: "EXECUTED" | "BLOCKED" | "CONFIRM_CANCELLED";
  reason?: "NOT_VISIBLE" | "PERMISSION_DENIED" | "DISABLED";
  events: Array<{ type: "confirmShown" | "actionExecuted" }>;
}

type JsonRecord = Record<string, unknown>;

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function asNode(value: unknown): PermissionsNode | undefined {
  if (!isRecord(value) || typeof value.type !== "string") {
    return undefined;
  }
  return value as PermissionsNode;
}

interface PermissionsNode extends JsonRecord {
  type: string;
  id?: string;
  permissions?: Record<string, unknown>;
  permissionCascade?: { keys?: unknown };
  props?: Record<string, unknown>;
  children?: JsonRecord[];
}

function stringValue(value: unknown): string {
  return typeof value === "string" ? value : "";
}

// --- L2 validation (fail-closed) ---

/** Validates the permission fields of a page; returns all L2 errors. */
export function validatePermissions(page: JsonRecord): L2Error[] {
  const meta = isRecord(page.meta) ? page.meta : {};
  const metaVersion = stringValue(meta.protocolVersion);
  const capabilities = Array.isArray(meta.requiredCapabilities)
    ? meta.requiredCapabilities.filter((c): c is string => typeof c === "string")
    : [];
  const hasPermissionFields = pageHasPermissionFields(page);

  const errors: L2Error[] = [];
  if (hasPermissionFields && !hasMinimumVersion(metaVersion)) {
    errors.push({ code: "PROTOCOL_VERSION_TOO_LOW", path: "meta.protocolVersion" });
  }
  if (hasPermissionFields && !capabilities.includes(PERMISSION_CAPABILITY)) {
    errors.push({ code: "CAPABILITY_REQUIRED", path: "meta.requiredCapabilities" });
  }

  const body = asNode(page.body);
  if (body !== undefined) {
    walkNodes(body, [], errors);
  }
  if (isRecord(page.actions)) {
    for (const [actionName, action] of Object.entries(page.actions)) {
      if (isRecord(action) && isRecord(action.content) && asNode(action.content) !== undefined) {
        walkNodes(asNode(action.content)!, [`actions.${actionName}.content`], errors);
      }
    }
  }
  return errors;
}

function hasMinimumVersion(version: string): boolean {
  const match = /^(\d+)\.(\d+)$/.exec(version);
  if (!match) {
    return false;
  }
  const major = Number(match[1]);
  const minor = Number(match[2]);
  return major > 2 || (major === 2 && minor >= 3);
}

function pageHasPermissionFields(page: JsonRecord): boolean {
  if (isRecord(page.body) && nodeHasPermissionFields(page.body)) {
    return true;
  }
  if (isRecord(page.actions)) {
    for (const action of Object.values(page.actions)) {
      if (isRecord(action) && isRecord(action.content) && nodeHasPermissionFields(action.content)) {
        return true;
      }
    }
  }
  return false;
}

function nodeHasPermissionFields(node: JsonRecord): boolean {
  if (node.permissionCascade !== undefined) {
    return true;
  }
  if (isRecord(node.props) && node.props.permissionIntent !== undefined) {
    return true;
  }
  if (isRecord(node.props)) {
    for (const key of ["actions", "toolbar"]) {
      if (Array.isArray(node.props[key])) {
        for (const entry of node.props[key]) {
          if (isRecord(entry) && entry.permissionIntent !== undefined) {
            return true;
          }
        }
      }
    }
  }
  if (Array.isArray(node.children) && node.children.some(nodeHasPermissionFields)) {
    return true;
  }
  if (isRecord(node.props) && Array.isArray(node.props.items)) {
    for (const item of node.props.items) {
      if (isRecord(item) && isRecord(item.content) && nodeHasPermissionFields(item.content)) {
        return true;
      }
    }
  }
  return false;
}

function walkNodes(node: PermissionsNode, stack: string[], errors: L2Error[]): void {
  const here = stack.join(".");
  const pathOf = (tail: string) => (here === "" ? `body.${tail}` : `body.${here}.${tail}`);

  // permissionCascade constraints (ADR-0023 D2a).
  if (node.permissionCascade !== undefined) {
    if (!CASCADE_TYPES.has(node.type)) {
      errors.push({ code: "PERMISSION_CASCADE_TYPE_INVALID", path: pathOf("permissionCascade") });
    }
    const keys = node.permissionCascade.keys;
    const keyList = Array.isArray(keys) ? keys.filter((k): k is string => typeof k === "string") : [];
    const uniqueKeys = new Set(keyList);
    const keysValid =
      Array.isArray(keys) &&
      keys.length > 0 &&
      uniqueKeys.size === keys.length &&
      keyList.every((k) => CASCADE_KEYS.has(k));
    if (!keysValid) {
      errors.push({ code: "PERMISSION_CASCADE_KEYS_INVALID", path: pathOf("permissionCascade.keys") });
    }
    for (const key of uniqueKeys) {
      if (keysValid && (!isRecord(node.permissions) || node.permissions[key] === undefined)) {
        errors.push({ code: "PERMISSION_CASCADE_SOURCE_MISSING", path: pathOf(`permissions.${key}`) });
      }
    }
  }

  // permissionIntent mount-point matrix (ADR-0023 D4b).
  if (isRecord(node.props) && node.props.permissionIntent !== undefined) {
    const intent = node.props.permissionIntent;
    if (node.type === "actionButton") {
      if (!INTENT_KEYS.has(stringValue(intent))) {
        errors.push({ code: "PERMISSION_INTENT_INVALID", path: pathOf("props.permissionIntent") });
      }
    } else {
      errors.push({ code: "PERMISSION_INTENT_FORBIDDEN", path: pathOf("props.permissionIntent") });
    }
  }

  if (node.type === "table" && isRecord(node.props)) {
    // Columns are not on the permission structure tree (D2a): intent there is
    // forbidden. Checked before actions/toolbar to match the fixture order.
    const columns = Array.isArray(node.props.columns) ? node.props.columns : [];
    columns.forEach((column, index) => {
      if (isRecord(column) && column.permissionIntent !== undefined) {
        errors.push({
          code: "PERMISSION_INTENT_FORBIDDEN",
          path: pathOf(`props.columns[${index}].permissionIntent`),
        });
      }
    });
    const actions = Array.isArray(node.props.actions) ? node.props.actions : [];
    actions.forEach((action, index) => {
      if (!isRecord(action) || action.permissionIntent === undefined) {
        return;
      }
      if (!INTENT_KEYS.has(stringValue(action.permissionIntent))) {
        errors.push({
          code: "PERMISSION_INTENT_INVALID",
          path: pathOf(`props.actions[${index}].permissionIntent`),
        });
      }
    });
    const toolbar = Array.isArray(node.props.toolbar) ? node.props.toolbar : [];
    toolbar.forEach((trigger, index) => {
      if (!isRecord(trigger) || trigger.permissionIntent === undefined) {
        return;
      }
      if (!INTENT_KEYS.has(stringValue(trigger.permissionIntent))) {
        errors.push({
          code: "PERMISSION_INTENT_INVALID",
          path: pathOf(`props.toolbar[${index}].permissionIntent`),
        });
      }
    });
  }

  if (Array.isArray(node.children)) {
    node.children.forEach((child, index) => {
      if (asNode(child) !== undefined) {
        walkNodes(asNode(child)!, [...stack, `children[${index}]`], errors);
      }
    });
  }
  if (node.type === "tabs" && isRecord(node.props) && Array.isArray(node.props.items)) {
    node.props.items.forEach((item, index) => {
      if (isRecord(item) && isRecord(item.content) && asNode(item.content) !== undefined) {
        walkNodes(asNode(item.content)!, [...stack, `props.items[${index}].content`], errors);
      }
    });
  }
}

// --- effectivePermission evaluation (ADR-0023 D3) ---

/**
 * Evaluates the effective permission of every approved target in the page
 * against the frozen $context snapshot. Targets in modal content and
 * navigated pages start new roots (ADR-0023 D2a).
 */
export function evaluatePermissionTargets(
  page: JsonRecord,
  context: NavigationContext,
): PermissionTarget[] {
  const targets: PermissionTarget[] = [];
  const body = asNode(page.body);
  if (body !== undefined) {
    collectTargets(body, [], targets, context, false);
  }
  if (isRecord(page.actions)) {
    for (const [actionName, action] of Object.entries(page.actions)) {
      if (isRecord(action) && isRecord(action.content) && asNode(action.content) !== undefined) {
        collectTargets(asNode(action.content)!, [], targets, context, true);
      }
      void actionName;
    }
  }
  if (isRecord(page.navigatedPage) && asNode(page.navigatedPage.body) !== undefined) {
    collectTargets(asNode(page.navigatedPage.body)!, [], targets, context, true);
  }
  return targets;
}

/**
 * W9 F-008: the gate target id mirrors the renderer's consumers
 * (render.tsx invokeAction / schema-table.tsx row + toolbar lookups): an
 * explicit `key` wins, otherwise the action ref. Registering only `key`
 * left an intent-marked action without a key unmatchable (targetId ""), so
 * the client-side permission gate was silently skipped for it.
 */
function actionGateTargetId(record: Record<string, unknown>): string {
  const key = stringValue(record.key);
  if (key !== "") {
    return key;
  }
  return stringValue(record.actionRef);
}

function collectTargets(
  node: PermissionsNode,
  ancestors: PermissionsNode[],
  targets: PermissionTarget[],
  context: NavigationContext,
  newRoot: boolean,
): void {
  const nodeId = node.id ?? "unnamed";

  // Form field targets (D4a whitelist inside default-mode forms).
  if (FORM_EDIT_FIELD_TYPES.has(node.type) && inDefaultForm(ancestors)) {
    const fieldId = (node.id ?? (isRecord(node.props) ? stringValue(node.props.field) : "")) || "unnamed";
    targets.push({
      targetId: fieldId,
      kind: "formField",
      key: "edit",
      cascadeApplied: cascadingAncestors(ancestors, "edit").length > 0,
      cascadedBy: cascadingAncestors(ancestors, "edit"),
      effectivePermission: effectivePermission(ancestors, "edit", context, node),
    });
  }

  // Default form submit targets (implicit edit intent, D4a). The form node is
  // the submit entry's ancestor (D2a form edge); its own permissions count
  // only when it declares a matching cascade (D3a).
  if (node.type === "form" && isDefaultForm(node) && typeof node.props?.submitAction === "string") {
    const submitAncestors = [...ancestors, node];
    targets.push({
      targetId: `${nodeId}:submit`,
      kind: "formSubmit",
      key: "edit",
      cascadeApplied: cascadingAncestors(submitAncestors, "edit").length > 0,
      cascadedBy: cascadingAncestors(submitAncestors, "edit"),
      effectivePermission: effectivePermission(submitAncestors, "edit", context),
    });
  }

  // Table action / toolbar intent targets and actionButton (D4b). The table
  // node is the mount ancestor of its actions/toolbar (D2a mount edges).
  if (node.type === "table" && isRecord(node.props)) {
    const mountAncestors = [...ancestors, node];
    const actions = Array.isArray(node.props.actions) ? node.props.actions : [];
    actions.forEach((action) => {
      if (!isRecord(action)) {
        return;
      }
      if (action.permissionIntent !== undefined) {
        const key = stringValue(action.permissionIntent) as PermissionKey;
        targets.push({
          targetId: actionGateTargetId(action),
          kind: "rowAction",
          key,
          cascadeApplied: cascadingAncestors(mountAncestors, key).length > 0,
          cascadedBy: cascadingAncestors(mountAncestors, key),
          effectivePermission: effectivePermission(mountAncestors, key, context),
        });
        return;
      }
      // Unmarked intents do not participate in cascade (D3b): local
      // permissions only.
      for (const key of INTENT_KEYS) {
        const local = isRecord(action.permissions) ? action.permissions[key] : undefined;
        if (local !== undefined) {
          targets.push({
            targetId: actionGateTargetId(action),
            kind: "rowAction",
            key: key as PermissionKey,
            cascadeApplied: false,
            cascadedBy: [],
            effectivePermission: evaluatePermissionValue(local, context),
          });
        }
      }
    });
    const toolbar = Array.isArray(node.props.toolbar) ? node.props.toolbar : [];
    toolbar.forEach((trigger) => {
      if (!isRecord(trigger)) {
        return;
      }
      if (trigger.permissionIntent !== undefined) {
        const key = stringValue(trigger.permissionIntent) as PermissionKey;
        targets.push({
          targetId: actionGateTargetId(trigger),
          kind: "toolbarTrigger",
          key,
          cascadeApplied: cascadingAncestors(mountAncestors, key).length > 0,
          cascadedBy: cascadingAncestors(mountAncestors, key),
          effectivePermission: effectivePermission(mountAncestors, key, context),
        });
        return;
      }
      for (const key of INTENT_KEYS) {
        const local = isRecord(trigger.permissions) ? trigger.permissions[key] : undefined;
        if (local !== undefined) {
          targets.push({
            targetId: actionGateTargetId(trigger),
            kind: "toolbarTrigger",
            key: key as PermissionKey,
            cascadeApplied: false,
            cascadedBy: [],
            effectivePermission: evaluatePermissionValue(local, context),
          });
        }
      }
    });
    // Columns are not on the permission structure tree (D2a): local
    // permissions only.
    const columns = Array.isArray(node.props.columns) ? node.props.columns : [];
    columns.forEach((column) => {
      if (!isRecord(column)) {
        return;
      }
      for (const key of ["view", "edit", "delete"] as const) {
        const local = isRecord(column.permissions) ? column.permissions[key] : undefined;
        if (local !== undefined) {
          targets.push({
            targetId: stringValue(column.field) || "unnamed",
            kind: "column",
            key,
            cascadeApplied: false,
            cascadedBy: [],
            effectivePermission: evaluatePermissionValue(local, context),
          });
        }
      }
    });
  }
  if (node.type === "actionButton" && isRecord(node.props) && node.props.permissionIntent !== undefined) {
    const key = stringValue(node.props.permissionIntent) as PermissionKey;
    targets.push({
      targetId: stringValue(node.props.key) || nodeId,
      kind: "actionButton",
      key,
      cascadeApplied: cascadingAncestors(ancestors, key).length > 0,
      cascadedBy: cascadingAncestors(ancestors, key),
      effectivePermission: effectivePermission(ancestors, key, context),
    });
  }

  if (Array.isArray(node.children)) {
    node.children.forEach((child) => {
      if (asNode(child) !== undefined) {
        collectTargets(asNode(child)!, [...ancestors, node], targets, context, newRoot);
      }
    });
  }
  if (node.type === "tabs" && isRecord(node.props) && Array.isArray(node.props.items)) {
    node.props.items.forEach((item) => {
      if (isRecord(item) && isRecord(item.content) && asNode(item.content) !== undefined) {
        collectTargets(asNode(item.content)!, [...ancestors, node], targets, context, newRoot);
      }
    });
  }
}

function isDefaultForm(node: PermissionsNode): boolean {
  return node.props?.mode === undefined || node.props.mode === "default";
}

function inDefaultForm(ancestors: PermissionsNode[]): boolean {
  // Nearest form ancestor must be a default-mode form.
  for (let i = ancestors.length - 1; i >= 0; i--) {
    if (ancestors[i].type === "form") {
      return isDefaultForm(ancestors[i]);
    }
  }
  return false;
}

function cascadingAncestors(ancestors: PermissionsNode[], key: PermissionKey): string[] {
  const ids: string[] = [];
  for (const ancestor of ancestors) {
    if (!isRecord(ancestor.permissionCascade) || !Array.isArray(ancestor.permissionCascade.keys)) {
      continue;
    }
    if ((ancestor.permissionCascade.keys as string[]).includes(key)) {
      ids.push(ancestor.id ?? "unnamed");
    }
  }
  return ids;
}

/** Permission values may be booleans or $context expressions; unknown values deny. */
function evaluatePermissionValue(value: unknown, context: NavigationContext): boolean {
  if (typeof value === "boolean") {
    return value;
  }
  if (typeof value === "string") {
    return evaluateExpression(value, context);
  }
  return false;
}

/** AND formula (ADR-0023 D3): cascade boundaries of ancestors + target local. */
function effectivePermission(
  ancestors: PermissionsNode[],
  key: PermissionKey,
  context: NavigationContext,
  targetNode?: PermissionsNode,
): boolean {
  for (const ancestor of ancestors) {
    if (!isRecord(ancestor.permissionCascade) || !Array.isArray(ancestor.permissionCascade.keys)) {
      continue;
    }
    if (!(ancestor.permissionCascade.keys as string[]).includes(key)) {
      continue;
    }
    const source = isRecord(ancestor.permissions) ? ancestor.permissions[key] : undefined;
    // W9 F-009: a cascade-declared key WITHOUT a permission source fails
    // closed (deny). The previous skip-and-continue evaluated malformed
    // structures as allow, rendering with the gate open.
    if (source === undefined || !evaluatePermissionValue(source, context)) {
      return false;
    }
  }
  if (targetNode !== undefined && isRecord(targetNode.permissions) && targetNode.permissions[key] !== undefined) {
    return evaluatePermissionValue(targetNode.permissions[key], context);
  }
  return true;
}

// --- Execution sequence (ADR-0023 D4c) ---

/** Gates an action: visible → permission → disabled/requiresSelection → confirm. */
export function executeAction(
  page: JsonRecord,
  request: ExecutionRequest,
  context: NavigationContext,
): ExecutionResult {
  const entry = evaluatePermissionTargets(page, context).find(
    (candidate) => candidate.targetId === request.targetId,
  );
  if (entry === undefined) {
    return { outcome: "BLOCKED", reason: "NOT_VISIBLE", events: [] };
  }

  if (!request.visible) {
    return { outcome: "BLOCKED", reason: "NOT_VISIBLE", events: [] };
  }
  if (entry.effectivePermission === false) {
    return { outcome: "BLOCKED", reason: "PERMISSION_DENIED", events: [] };
  }
  if (request.disabled === true || request.requiresSelection === true) {
    return { outcome: "BLOCKED", reason: "DISABLED", events: [] };
  }
  const events: Array<{ type: "confirmShown" | "actionExecuted" }> = [];
  if (request.confirm === true) {
    events.push({ type: "confirmShown" });
    if (request.confirmed !== true) {
      return { outcome: "CONFIRM_CANCELLED", events };
    }
  }
  events.push({ type: "actionExecuted" });
  return { outcome: "EXECUTED", events };
}
