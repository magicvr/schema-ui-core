import { executeAction, type ExecutionResult } from "@/renderer/permissions";
import type { NavigationContext } from "@/protocol/app-manifest";

/**
 * R5 D-ACT non-batch row action execution (frozen Q1: batch excluded).
 *
 * Wraps the R4 executeAction engine so UI actions gate through the same
 * permission / confirm / disabled sequence as the frozen permission fixtures.
 * `page` is the R5 example page document; `targetId` is a table row action key.
 */

export type RowActionOutcome = ExecutionResult["outcome"];

export interface RowActionRequest {
  page: Record<string, unknown>;
  targetId: string;
  context: NavigationContext;
  /** Defaults to true; a hidden action fails closed as NOT_VISIBLE. */
  visible?: boolean;
  confirm?: boolean;
  confirmed?: boolean;
  disabled?: boolean;
  requiresSelection?: boolean;
}

export interface RowActionResult {
  outcome: RowActionOutcome;
  reason?: ExecutionResult["reason"];
  permissionDenied: boolean;
  confirmed?: boolean;
}

export function runRowAction(request: RowActionRequest): RowActionResult {
  const result = executeAction(
    request.page,
    {
      targetId: request.targetId,
      visible: request.visible ?? true,
      ...(request.confirm === undefined ? {} : { confirm: request.confirm }),
      ...(request.confirmed === undefined ? {} : { confirmed: request.confirmed }),
      ...(request.disabled === undefined ? {} : { disabled: request.disabled }),
      ...(request.requiresSelection === undefined
        ? {}
        : { requiresSelection: request.requiresSelection }),
    },
    request.context,
  );
  return {
    outcome: result.outcome,
    reason: result.reason,
    permissionDenied: result.outcome === "BLOCKED" && result.reason === "PERMISSION_DENIED",
    ...(request.confirm === true
      ? { confirmed: result.outcome === "EXECUTED" }
      : {}),
  };
}
