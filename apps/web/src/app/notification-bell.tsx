import { Bell } from "lucide-react";
import { useEffect, useRef, useState } from "react";

import { NOTIFICATIONS_READ_NAMESPACE, subscribeToConfigChanges } from "@/app/config-events";
import { useTranslate } from "@/i18n/runtime";

/**
 * F-04 header notification bell (GOAL-006 D-002 `5): unread-count badge +
 * lightweight dropdown with the 5 latest notifications and a "view all" link.
 * Fail-open: any fetch failure hides the badge/panel silently — the shell must
 * never break because the notifications surface is unavailable.
 */

interface NotificationItem {
  id: string;
  event: string;
  title: string;
  body: string;
  titleKey?: string;
  bodyKey?: string;
  read: boolean;
  createdAt: string;
}

export interface NotificationBellProps {
  /** Authed same-origin fetch (defaults to globalThis.fetch for bare tests). */
  fetcher?: typeof fetch;
  /** Navigation callback for the "view all" link (kept out of router). */
  onViewAll: () => void;
  /** W13 T-06: opens one notification — navigate to the list page with the
   * detail target (/notifications?open=<id>); the page expands + marks read. */
  onOpenItem: (id: string) => void;
}

export function NotificationBell({ fetcher, onViewAll, onOpenItem }: NotificationBellProps) {
  const t = useTranslate();
  const [unread, setUnread] = useState<number | null>(null);
  const [items, setItems] = useState<NotificationItem[] | null>(null);
  const [open, setOpen] = useState(false);
  const [error, setError] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  const active = fetcher ?? globalThis.fetch;

  const loadCount = () => {
    active("/api/notifications/unread-count")
      .then((response) => {
        if (!response.ok) {
          setError(true);
          return;
        }
        return response.json() as Promise<{ unread: number }>;
      })
      .then((body) => {
        if (body !== undefined) {
          setUnread(body.unread);
          setError(false);
        }
      })
      .catch(() => setError(true));
  };

  useEffect(() => {
    loadCount();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [active]);

  // W13 T-06: any read/read-all (header or list page) refreshes the badge —
  // and an open dropdown refetches its items — immediately.
  useEffect(() => {
    return subscribeToConfigChanges(NOTIFICATIONS_READ_NAMESPACE, () => {
      loadCount();
      if (open) {
        setItems(null);
        active("/api/notifications?pageSize=5")
          .then((response) => {
            if (!response.ok) {
              setItems(null);
              return;
            }
            return response.json() as Promise<{ items: NotificationItem[] }>;
          })
          .then((body) => {
            if (body !== undefined) {
              setItems(body.items);
            }
          })
          .catch(() => setItems(null));
      }
    });
  }, [active, open]);

  useEffect(() => {
    if (!open) {
      return;
    }
    active("/api/notifications?pageSize=5")
      .then((response) => {
        if (!response.ok) {
          setError(true);
          setItems(null);
          return;
        }
        return response.json() as Promise<{ items: NotificationItem[] }>;
      })
      .then((body) => {
        if (body !== undefined) {
          setItems(body.items);
          setError(false);
        }
      })
      .catch(() => {
        setError(true);
        setItems(null);
      });
  }, [open, active]);

  // Close on outside click / Escape (a11y parity with the mobile drawer).
  useEffect(() => {
    if (!open) {
      return;
    }
    const onPointerDown = (event: PointerEvent) => {
      if (rootRef.current !== null && !rootRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
      }
    };
    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  const badge = unread !== null && unread > 0 ? unread : null;

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        aria-label={
          unread !== null && unread > 0
            ? t("shell.notifications.bell") + " (" + (unread > 99 ? "99+" : String(unread)) + ")"
            : t("shell.notifications.bell")
        }
        aria-expanded={open}
        className="relative inline-flex size-9 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        onClick={() => setOpen((value) => !value)}
      >
        <Bell aria-hidden="true" className="size-4" />
        {badge !== null ? (
          <span
            aria-hidden="true"
            className="absolute right-1 top-1 inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-destructive px-1 text-[10px] font-semibold leading-none text-destructive-foreground"
          >
            {badge > 99 ? "99+" : badge}
          </span>
        ) : null}
      </button>
      {open ? (
        <div
          role="menu"
          aria-label={t("shell.notifications.bell")}
          className="absolute right-0 top-11 z-50 w-80 rounded-md border border-border bg-card shadow-lg"
        >
          <div className="max-h-80 overflow-y-auto p-1">
            {error ? (
              <p role="alert" className="px-3 py-2 text-xs text-muted-foreground">
                {t("shell.notifications.unavailable")}
              </p>
            ) : items === null ? (
              <p role="status" className="px-3 py-2 text-xs text-muted-foreground">
                {t("feedback.loading")}
              </p>
            ) : items.length === 0 ? (
              <p className="px-3 py-2 text-xs text-muted-foreground">{t("feedback.noItemsMatch")}</p>
            ) : (
              <ul className="divide-y divide-border">
                {items.map((item) => {
                  const titleText = item.titleKey !== undefined && item.titleKey !== "" ? t(item.titleKey) : item.title;
                  const bodyText = item.bodyKey !== undefined && item.bodyKey !== "" ? t(item.bodyKey) : item.body;
                  // W13 T-06: a dropdown entry is actionable — marking the
                  // notification read (best-effort) and opening it on the
                  // list page with its detail expanded.
                  return (
                  <li key={item.id} className="px-1.5 py-1">
                    <button
                      type="button"
                      role="menuitem"
                      className="flex w-full items-start gap-2.5 rounded-md px-2 py-1.5 text-left transition-colors hover:bg-accent"
                      onClick={() => {
                        setOpen(false);
                        void active("/api/notifications/" + encodeURIComponent(item.id) + "/read", { method: "POST" }).catch(() => undefined);
                        onOpenItem(item.id);
                      }}
                    >
                      <span
                        aria-hidden="true"
                        className={"mt-1.5 size-2 shrink-0 rounded-full " + (item.read ? "bg-transparent" : "bg-primary")}
                      />
                      <span className="min-w-0 flex-1">
                        <span className="block text-sm font-medium text-foreground">{titleText}</span>
                        <span className="mt-0.5 line-clamp-2 block text-xs text-muted-foreground">{bodyText}</span>
                      </span>
                    </button>
                  </li>
                  );
                })}
              </ul>
            )}
          </div>
          <div className="border-t border-border p-1">
            <button
              type="button"
              role="menuitem"
              className="w-full rounded-md px-3 py-2 text-left text-sm text-foreground transition-colors hover:bg-accent"
              onClick={() => {
                setOpen(false);
                onViewAll();
              }}
            >
              {t("shell.notifications.viewAll")}
            </button>
          </div>
        </div>
      ) : null}
    </div>
  );
}