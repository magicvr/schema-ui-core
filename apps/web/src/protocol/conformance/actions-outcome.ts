/**
 * actions fixture adapter — transport outcome → host events (non-batch).
 */

type Transport =
  | { type: "success"; status: number }
  | { type: "httpError"; status: number; body?: { message?: string; errors?: unknown } }
  | { type: "timeout" }
  | { type: "networkError" }
  | { type: "abort" };

interface Behavior {
  behavior: string;
  message?: string;
  url?: string;
}

const SAFE_MESSAGES: Record<number, string> = {
  403: "无权限访问",
  404: "资源不存在",
};

const SERVER_ERROR_DEFAULT = "系统异常，请稍后重试";
const TIMEOUT_DEFAULT = "请求超时，请稍后重试";
const NETWORK_DEFAULT = "网络异常，请检查网络连接";

function emitBehavior(
  behavior: Behavior | undefined,
  context: { tableId?: string },
  events: Array<Record<string, unknown>>,
): void {
  if (!behavior) {
    return;
  }
  switch (behavior.behavior) {
    case "toast":
      events.push({ type: "toast", message: behavior.message ?? "" });
      break;
    case "reload":
      events.push({ type: "reloadTable", tableId: context.tableId });
      break;
    case "navigate":
      events.push({ type: "navigate", url: behavior.url });
      break;
    case "closeModal":
      events.push({ type: "closeModal" });
      break;
    default:
      break;
  }
}

export function runActionOutcome(input: Record<string, unknown>): Record<string, unknown> {
  const transport = input.transport as Transport;
  const onSuccess = input.onSuccess as Behavior | undefined;
  const onError = input.onError as Behavior | undefined;
  const context = (input.context as { tableId?: string }) ?? {};
  const events: Array<Record<string, unknown>> = [];

  if (transport.type === "abort") {
    return { ok: false, events: [] };
  }

  if (transport.type === "success") {
    events.push({ type: "requestSucceeded", status: transport.status });
    emitBehavior(onSuccess, context, events);
    return { ok: true, events };
  }

  if (transport.type === "timeout") {
    events.push({
      type: "errorState",
      display: TIMEOUT_DEFAULT,
      retryable: true,
      outcome: "unknown",
    });
    emitBehavior(onError, context, events);
    return { ok: false, events };
  }

  if (transport.type === "networkError") {
    events.push({
      type: "errorState",
      display: NETWORK_DEFAULT,
      retryable: true,
      outcome: "unknown",
    });
    emitBehavior(onError, context, events);
    return { ok: false, events };
  }

  if (transport.type === "httpError") {
    const status = transport.status;
    const body = transport.body ?? {};

    if (status === 400 && Array.isArray(body.errors)) {
      events.push({ type: "fieldErrors", errors: body.errors });
      if (onError?.behavior === "toast") {
        events.push({ type: "toast", message: onError.message ?? body.message ?? "" });
      } else {
        // validation errors suppress navigate; default toast from response message
        events.push({ type: "toast", message: body.message ?? "" });
      }
      return { ok: false, events };
    }

    if (status === 401 || status === 403) {
      events.push({ type: "authFailure", status });
      events.push({
        type: "errorState",
        display: status === 403 ? SAFE_MESSAGES[403] : null,
      });
      // auth hook suppresses onError
      return { ok: false, events };
    }

    if (status === 404) {
      events.push({
        type: "errorState",
        display: SAFE_MESSAGES[404],
      });
      emitBehavior(onError, context, events);
      return { ok: false, events };
    }

    // 5xx / other
    events.push({
      type: "errorState",
      display: SERVER_ERROR_DEFAULT,
    });
    emitBehavior(onError, context, events);
    return { ok: false, events };
  }

  return { ok: false, events: [] };
}
