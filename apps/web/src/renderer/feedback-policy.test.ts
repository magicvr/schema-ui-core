import { describe, expect, it, vi } from "vitest";

import { classifyFeedbackFailure, feedbackFromError } from "@/renderer/feedback-policy";

describe("R4 feedback policy", () => {
  it("classifies read transport failures with an explicit retry policy", () => {
    expect(classifyFeedbackFailure({ status: 503, code: "SERVICE_MAINTENANCE" })).toEqual({
      category: "maintenance",
      messageKey: "feedback.maintenance",
      retryableRead: true,
    });
    expect(classifyFeedbackFailure({ status: 503 })).toEqual({
      category: "unavailable",
      messageKey: "feedback.unavailable",
      retryableRead: true,
    });
    expect(classifyFeedbackFailure({ code: "REQUEST_FAILED" })).toEqual({
      category: "offline",
      messageKey: "feedback.offline",
      retryableRead: true,
    });
    expect(classifyFeedbackFailure({ status: 504 })).toEqual({
      category: "timeout",
      messageKey: "feedback.timeout",
      retryableRead: true,
    });
    const timeout = feedbackFromError(new DOMException("aborted", "AbortError"), { retry: vi.fn() });
    expect(timeout.messageKey).toBe("feedback.timeout");
    expect(timeout.retry).toBeTypeOf("function");
    const offline = feedbackFromError(new TypeError("Failed to fetch"), { retry: vi.fn() });
    expect(offline.messageKey).toBe("feedback.offline");
    expect(offline.retry).toBeTypeOf("function");
    const fallbackOffline = feedbackFromError(new Error("transport failed"), {
      fallbackCode: "REQUEST_FAILED",
      retry: vi.fn(),
    });
    expect(fallbackOffline.messageKey).toBe("feedback.offline");
    expect(fallbackOffline.retry).toBeTypeOf("function");
  });

  it("keeps auth, permission, conflict and validation failures non-retryable", () => {
    for (const input of [
      { status: 401 },
      { status: 403 },
      { status: 409 },
      { status: 422 },
    ]) {
      expect(classifyFeedbackFailure(input).retryableRead).toBe(false);
    }
  });

  it("attaches retry only when the read owner explicitly supplies it", () => {
    const retry = vi.fn();
    const readFeedback = feedbackFromError(
      {
        status: 503,
        code: "READ_UNAVAILABLE",
        message: "server detail",
      },
      { retry },
    );
    expect(readFeedback.messageKey).toBe("feedback.unavailable");
    expect(readFeedback.retry).toBe(retry);

    const writeFeedback = feedbackFromError({
      status: 503,
      code: "WRITE_UNAVAILABLE",
      message: "server detail",
    });
    expect(writeFeedback.retry).toBeUndefined();

    const forbidden = feedbackFromError({ status: 403, code: "FORBIDDEN", message: "denied" }, { retry });
    expect(forbidden.messageKey).toBe("feedback.permissionDenied");
    expect(forbidden.retry).toBeUndefined();
  });

  it("preserves an explicit catalog key and safe correlation metadata", () => {
    const feedback = feedbackFromError({
      status: 400,
      code: "INVALID_SITE_TITLE",
      message: "siteTitle must not be empty",
      messageKey: "error.siteTitleEmpty",
      params: { field: "siteTitle" },
      correlationId: "req-r4-001",
    });
    expect(feedback.messageKey).toBe("error.siteTitleEmpty");
    expect(feedback.correlationId).toBe("req-r4-001");
    expect(feedback.retry).toBeUndefined();
  });
});
