// Notification center (W13 T-06 · GOAL-014): replaces the schema-driven table
// on the notifications page with an interactive list whose row click opens
// the detail inline AND marks the notification read — matching user
// intuition (the per-row "Mark read" action is removed from the schema).
//
// Behaviour:
//   - reads the shared search-query state (targetTable) so the schema search
//     form (q + read filter) keeps binding
//   - row click → expand detail + mark read (POST via the authed fetcher;
//     the notifications.read config-change header refreshes the bell badge)
//   - deep link ?open=<id> (bell dropdown entries) expands + marks that one
//   - toolbar keeps "Mark all read"
//
// The component is registered in the GOAL-018 custom-component registry;
// main.tsx imports it for the side-effect registration.

import { useCallback, useEffect, useMemo, useState } from "react";

import { Check } from "lucide-react";

import { formatDisplayTime } from "@/lib/datetime";
import { registerCustomComponent, type CustomComponentProps } from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";
import { useTranslate } from "@/i18n/runtime";

interface NotificationRow {
  id: string;
  event: string;
  title: string;
  body: string;
  titleKey?: string;
  bodyKey?: string;
  read: boolean;
  createdAt: string;
}

interface NotificationListResponse {
  items: NotificationRow[];
  total: number;
  page: number;
  pageSize: number;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

const DEFAULT_PAGE_SIZE = 10;

export function NotificationCenter({ node, context }: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const targetTable =
    isRecord(node.props) && typeof node.props.targetTable === "string" && node.props.targetTable !== ""
      ? node.props.targetTable
      : "notifications-table";
  const query = crud?.tableQuery(targetTable) ?? { page: 1, pageSize: DEFAULT_PAGE_SIZE };
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const [state, setState] = useState<{
    status: "loading" | "ready" | "error";
    items: NotificationRow[];
    total: number;
  }>({ status: "loading", items: [], total: 0 });
  const [expandedId, setExpandedId] = useState<string | null>(null);

  // Deep link from the bell dropdown: /notifications?open=<id>.
  const routeOpen = useMemo(() => {
    const route = isRecord(context.route) ? (context.route as { query?: Record<string, unknown> }).query : undefined;
    return isRecord(route) && typeof route.open === "string" && route.open !== "" ? route.open : null;
  }, [context]);

  // Marks one notification read (best-effort POST) and flips the local row.
  // The notifications.read config-change header on the response refreshes
  // the header bell badge immediately.
  const markReadLocal = useCallback(
    (id: string) => {
      void fetcher("/api/notifications/" + encodeURIComponent(id) + "/read", { method: "POST" }).catch(() => undefined);
      setState((prev) => ({
        ...prev,
        items: prev.items.map((item) => (item.id === id ? { ...item, read: true } : item)),
      }));
    },
    [fetcher],
  );

  const load = useCallback(async () => {
    const params = new URLSearchParams();
    const q = typeof query.q === "string" && query.q.trim() !== "" ? query.q.trim() : "";
    if (q !== "") params.set("q", q);
    const readFilter = isRecord(query.filters) && typeof query.filters.read === "string" ? query.filters.read : "";
    if (readFilter !== "") params.set("read", readFilter);
    params.set("page", String(query.page ?? 1));
    params.set("pageSize", String(query.pageSize ?? DEFAULT_PAGE_SIZE));
    // Silent refresh when rows are already visible (row-click mark-read keeps
    // the list stable; the first load still shows the loading state).
    setState((prev) => (prev.items.length === 0 ? { ...prev, status: "loading" } : prev));
    try {
      const response = await fetcher("/api/notifications?" + params.toString());
      if (!response.ok) {
        setState({ status: "error", items: [], total: 0 });
        return;
      }
      const body = (await response.json()) as NotificationListResponse;
      const items = Array.isArray(body.items) ? body.items : [];
      setState({ status: "ready", items, total: typeof body.total === "number" ? body.total : 0 });
      // Deep link: expand + mark the targeted notification after load.
      if (routeOpen !== null && items.some((item) => item.id === routeOpen)) {
        setExpandedId(routeOpen);
        markReadLocal(routeOpen);
      }
    } catch {
      setState({ status: "error", items: [], total: 0 });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [fetcher, query.q, query.page, query.pageSize, query.filters, routeOpen]);

  useEffect(() => {
    void load();
  }, [load]);

  const markAllRead = useCallback(() => {
    void fetcher("/api/notifications/read-all", { method: "POST" }).catch(() => undefined);
    setState((prev) => ({ ...prev, items: prev.items.map((item) => ({ ...item, read: true })) }));
    void load();
  }, [fetcher, load]);

  const toggleRow = (row: NotificationRow) => {
    setExpandedId((current) => (current === row.id ? null : row.id));
    if (!row.read) {
      markReadLocal(row.id);
      void load();
    }
  };

  const page = query.page ?? 1;
  const pageSize = query.pageSize ?? DEFAULT_PAGE_SIZE;
  const pages = Math.max(1, Math.ceil(state.total / pageSize));
  const setPage = (next: number) => {
    if (crud !== null && next >= 1 && next <= pages) {
      crud.setTableQuery(targetTable, { ...query, page: next });
    }
  };

  return (
    <div data-notification-center className="space-y-3">
      <div className="flex items-center justify-between gap-3">
        <p className="text-sm text-muted-foreground">
          {state.total} {state.total === 1 ? t("feedback.item") : t("feedback.items")}
        </p>
        <button
          type="button"
          onClick={markAllRead}
          className="inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md border border-input/80 bg-background px-3.5 text-sm font-medium text-muted-foreground shadow-2xs transition-all duration-150 hover:border-muted-foreground/30 hover:bg-accent/40 hover:text-foreground"
        >
          <Check aria-hidden="true" className="size-3.5" />
          {t("schema.notifications.toolbar.markAllRead")}
        </button>
      </div>

      {state.status === "loading" ? (
        <p role="status" className="text-sm text-muted-foreground">{t("feedback.loading")}</p>
      ) : state.status === "error" ? (
        <p role="alert" className="text-sm text-destructive">{t("shell.notifications.unavailable")}</p>
      ) : state.items.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t("schema.notifications.empty")}</p>
      ) : (
        <ul className="divide-y divide-border rounded-xl border border-border/70 bg-card/85">
          {state.items.map((item) => {
            const expanded = expandedId === item.id;
            const titleText = item.titleKey !== undefined && item.titleKey !== "" ? t(item.titleKey) : item.title;
            const bodyText = item.bodyKey !== undefined && item.bodyKey !== "" ? t(item.bodyKey) : item.body;
            return (
              <li key={item.id}>
                <button
                  type="button"
                  data-notification-row={item.id}
                  aria-expanded={expanded}
                  className="flex w-full cursor-pointer items-start gap-3 px-4 py-3 text-left transition-colors hover:bg-accent/40"
                  onClick={() => toggleRow(item)}
                >
                  <span
                    aria-hidden="true"
                    className={"mt-1.5 size-2 shrink-0 rounded-full " + (item.read ? "bg-transparent" : "bg-primary")}
                  />
                  <span className="min-w-0 flex-1">
                    <span className={"block truncate text-sm " + (item.read ? "text-muted-foreground" : "font-semibold text-foreground")}>
                      {titleText}
                    </span>
                    <span className="mt-0.5 block truncate text-xs text-muted-foreground">{bodyText}</span>
                  </span>
                  <span className="shrink-0 text-xs text-muted-foreground/70">
                    {formatDisplayTime(item.createdAt) ?? ""}
                  </span>
                </button>
                {expanded ? (
                  <div data-notification-detail={item.id} className="border-t border-border/60 bg-muted/30 px-4 py-3 pl-9 text-sm">
                    <p className="whitespace-pre-wrap break-words text-foreground">{bodyText}</p>
                    <dl className="mt-2 flex flex-wrap gap-x-6 gap-y-1 text-xs text-muted-foreground">
                      <div>
                        <dt className="inline">{t("schema.notifications.column.event")}: </dt>
                        <dd className="inline">{item.event}</dd>
                      </div>
                      <div>
                        <dt className="inline">{t("schema.notifications.column.read")}: </dt>
                        <dd className="inline">{item.read ? t("schema.notifications.filter.read.read") : t("schema.notifications.filter.read.unread")}</dd>
                      </div>
                      <div>
                        <dt className="inline">{t("schema.notifications.column.created")}: </dt>
                        <dd className="inline">{formatDisplayTime(item.createdAt) ?? ""}</dd>
                      </div>
                    </dl>
                  </div>
                ) : null}
              </li>
            );
          })}
        </ul>
      )}

      {state.status === "ready" && pages > 1 ? (
        <nav className="flex items-center gap-1.5" aria-label={t("feedback.pagination")}>
          <button
            type="button"
            disabled={page <= 1}
            aria-label={t("feedback.previousPage")}
            onClick={() => setPage(page - 1)}
            className="flex h-7 w-7 items-center justify-center rounded-md border border-input bg-background text-sm shadow-sm transition-colors hover:bg-accent hover:text-foreground disabled:opacity-40"
          >
            {"‹"}
          </button>
          <span className="px-1.5 text-xs text-muted-foreground">
            {page} / {pages}
          </span>
          <button
            type="button"
            disabled={page >= pages}
            aria-label={t("feedback.nextPage")}
            onClick={() => setPage(page + 1)}
            className="flex h-7 w-7 items-center justify-center rounded-md border border-input bg-background text-sm shadow-sm transition-colors hover:bg-accent hover:text-foreground disabled:opacity-40"
          >
            {"›"}
          </button>
        </nav>
      ) : null}
    </div>
  );
}

registerCustomComponent("notification-center", NotificationCenter);
