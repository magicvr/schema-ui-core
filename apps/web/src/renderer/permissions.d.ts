import { type NavigationContext } from "@/protocol/app-manifest";
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
export type L2ErrorCode = "PROTOCOL_VERSION_TOO_LOW" | "CAPABILITY_REQUIRED" | "PERMISSION_CASCADE_TYPE_INVALID" | "PERMISSION_CASCADE_KEYS_INVALID" | "PERMISSION_CASCADE_SOURCE_MISSING" | "PERMISSION_INTENT_FORBIDDEN" | "PERMISSION_INTENT_INVALID";
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
    events: Array<{
        type: "confirmShown" | "actionExecuted";
    }>;
}
type JsonRecord = Record<string, unknown>;
/** Validates the permission fields of a page; returns all L2 errors. */
export declare function validatePermissions(page: JsonRecord): L2Error[];
/**
 * Evaluates the effective permission of every approved target in the page
 * against the frozen $context snapshot. Targets in modal content and
 * navigated pages start new roots (ADR-0023 D2a).
 */
export declare function evaluatePermissionTargets(page: JsonRecord, context: NavigationContext): PermissionTarget[];
/** Gates an action: visible → permission → disabled/requiresSelection → confirm. */
export declare function executeAction(page: JsonRecord, request: ExecutionRequest, context: NavigationContext): ExecutionResult;
export {};
