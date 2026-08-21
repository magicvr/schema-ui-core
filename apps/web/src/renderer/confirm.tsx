/**
 * Confirm dialog for row actions that declare a `confirm` string (S4 · GOAL-007).
 *
 * One-time renderer completion allowed by I-007-003 §9.5 (`confirm*.tsx`). The
 * message text comes from the schema row action entry (`table.props.actions[].
 * confirm` / `confirmKey`), so the dialog stays fixture-driven; cancelling
 * mirrors the frozen `executeAction` CONFIRM_CANCELLED path and never issues
 * a request.
 *
 * W14 F-12 (GOAL-018): keyboard/accessibility — initial focus on Cancel, ESC
 * cancels, and Tab is trapped inside the dialog.
 */

import { useEffect, useRef } from "react";

import { useTranslate } from "@/i18n/runtime";

export function ConfirmDialog({
  message,
  onConfirm,
  onCancel,
}: {
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
}) {
  const t = useTranslate();
  const rootRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const root = rootRef.current;
    if (root === null) {
      return;
    }
    const focusable = () =>
      Array.from(
        root.querySelectorAll<HTMLElement>(
          'button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
        ),
      );
    const initial = focusable().find((element) => element.textContent?.includes(t("feedback.cancel")));
    (initial ?? focusable()[0])?.focus();

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        event.preventDefault();
        onCancel();
        return;
      }
      if (event.key !== "Tab") {
        return;
      }
      const items = focusable();
      if (items.length === 0) {
        return;
      }
      const first = items[0];
      const last = items[items.length - 1];
      const active = document.activeElement;
      if (event.shiftKey && (active === first || !root.contains(active))) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && (active === last || !root.contains(active))) {
        event.preventDefault();
        first.focus();
      }
    };

    root.addEventListener("keydown", onKeyDown);
    return () => root.removeEventListener("keydown", onKeyDown);
  }, [onCancel, t]);

  return (
    <div
      ref={rootRef}
      className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-overlay p-4"
      role="dialog"
      aria-modal="true"
      aria-label={t("feedback.confirmAction")}
    >
      <div className="mt-24 w-full max-w-sm rounded-lg border border-border bg-card p-6 shadow-md">
        <p className="text-sm leading-6 text-foreground">{message}</p>
        <div className="mt-5 flex justify-end gap-2">
          <button
            type="button"
            onClick={onCancel}
            className="h-9 rounded-md border border-input bg-background px-3 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            {t("feedback.cancel")}
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className="h-9 rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground transition-opacity hover:opacity-90"
          >
            {t("feedback.confirm")}
          </button>
        </div>
      </div>
    </div>
  );
}
