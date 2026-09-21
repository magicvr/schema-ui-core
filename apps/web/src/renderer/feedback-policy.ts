import type { FeedbackNotice } from "@/components/ui/feedback";

export type FeedbackFailureCategory =
  | "validation"
  | "authentication"
  | "forbidden"
  | "not-found"
  | "conflict"
  | "rate-limited"
  | "maintenance"
  | "unavailable"
  | "offline"
  | "timeout"
  | "unknown";

export interface FeedbackFailurePolicy {
  category: FeedbackFailureCategory;
  messageKey: string;
  /** True only when the underlying operation is a safe idempotent read. */
  retryableRead: boolean;
}

export interface FeedbackFailureLike {
  status?: number;
  code?: string;
  message?: string;
  messageKey?: string;
  params?: Record<string, unknown>;
  correlationId?: string;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function failureLikeOf(error: unknown): FeedbackFailureLike {
  if (!isRecord(error)) {
    return {
      message: error instanceof Error ? error.message : undefined,
    };
  }
  // The shared timeout wrapper surfaces an AbortError. It has no HTTP status
  // or application code, so normalize that transport signal before policy
  // classification instead of degrading it to the generic failure copy.
  if (error.name === "AbortError") {
    return {
      code: "REQUEST_TIMEOUT",
      message: typeof error.message === "string" ? error.message : undefined,
    };
  }
  // Browser fetch rejects network loss as a TypeError (often "Failed to
  // fetch"). Normalize that signal so read surfaces expose the explicit
  // offline retry affordance; parse/validation Errors remain non-retryable.
  if (
    error.name === "TypeError" ||
    (typeof error.message === "string" && /failed to fetch|network|offline/i.test(error.message))
  ) {
    return {
      code: "REQUEST_FAILED",
      message: typeof error.message === "string" ? error.message : undefined,
    };
  }
  return {
    ...(typeof error.status === "number" ? { status: error.status } : {}),
    ...(typeof error.code === "string" ? { code: error.code } : {}),
    ...(typeof error.message === "string" ? { message: error.message } : {}),
    ...(typeof error.messageKey === "string" ? { messageKey: error.messageKey } : {}),
    ...(isRecord(error.params) ? { params: error.params } : {}),
    ...(typeof error.correlationId === "string" ? { correlationId: error.correlationId } : {}),
  };
}

/**
 * Maps existing status/code/transport signals to the frozen R1 user-facing
 * category. The retry decision remains read-only: callers must explicitly
 * pass a read recovery callback to `feedbackFromError`.
 */
export function classifyFeedbackFailure(input: {
  status?: number;
  code?: string;
  transport?: "offline" | "timeout";
}): FeedbackFailurePolicy {
  const code = (input.code ?? "").toUpperCase();
  const status = input.status;
  if (input.transport === "timeout" || code.includes("TIMEOUT") || status === 504) {
    return { category: "timeout", messageKey: "feedback.timeout", retryableRead: true };
  }
  if (input.transport === "offline" || code === "REQUEST_FAILED" || code.includes("OFFLINE")) {
    return { category: "offline", messageKey: "feedback.offline", retryableRead: true };
  }
  if (code.includes("MAINTENANCE") || code === "SERVICE_MAINTENANCE") {
    return { category: "maintenance", messageKey: "feedback.maintenance", retryableRead: true };
  }
  if (status === 401 || code.includes("AUTH") || code.includes("REAUTH")) {
    return { category: "authentication", messageKey: "feedback.authenticationRequired", retryableRead: false };
  }
  if (status === 403 || code === "FORBIDDEN" || code.includes("PERMISSION")) {
    return { category: "forbidden", messageKey: "feedback.permissionDenied", retryableRead: false };
  }
  if (status === 404 || code.includes("NOT_FOUND")) {
    return { category: "not-found", messageKey: "feedback.notFound", retryableRead: false };
  }
  if (status === 409 || code.includes("CONFLICT")) {
    return { category: "conflict", messageKey: "feedback.conflict", retryableRead: false };
  }
  if (status === 429 || code.includes("RATE_LIMIT")) {
    return { category: "rate-limited", messageKey: "feedback.rateLimited", retryableRead: true };
  }
  if (
    status === 400 ||
    status === 422 ||
    code.includes("VALIDATION") ||
    code.includes("INVALID") ||
    code.includes("FIELD")
  ) {
    return { category: "validation", messageKey: "feedback.validationFailed", retryableRead: false };
  }
  if (status !== undefined && status >= 500 && status <= 599) {
    return { category: "unavailable", messageKey: "feedback.unavailable", retryableRead: true };
  }
  return { category: "unknown", messageKey: "feedback.operationFailed", retryableRead: false };
}

/** Converts any existing action/resource failure into the shared notice shape. */
export function feedbackFromError(
  error: unknown,
  options: { retry?: () => void | Promise<void>; fallbackCode?: string } = {},
): FeedbackNotice {
  const details = failureLikeOf(error);
  const code = details.code ?? options.fallbackCode;
  const policy = classifyFeedbackFailure({ status: details.status, code });
  return {
    kind: "error",
    code: code ?? "REQUEST_FAILED",
    message: details.message ?? "The operation could not be completed.",
    messageKey: details.messageKey ?? policy.messageKey,
    ...(details.params === undefined ? {} : { params: details.params }),
    ...(details.correlationId === undefined ? {} : { correlationId: details.correlationId }),
    ...(options.retry !== undefined && policy.retryableRead ? { retry: options.retry } : {}),
  };
}
