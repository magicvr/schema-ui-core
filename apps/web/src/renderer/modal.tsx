import { useEffect, useRef, type ReactNode } from "react";

/**
 * Minimal modal host for Schema-driven create/edit forms (S4 · GOAL-007).
 *
 * One-time renderer completion allowed by I-007-003 §9.5 (`modal*.tsx`). The
 * host is deliberately generic: the schema action (`ModalAction.content`)
 * supplies the body, the title comes from the triggering affordance label, and
 * `onClose` is wired to the page's action executor. No record-specific logic
 * lives here — a fixture-only page change never touches this file.
 *
 * Shared-host accessibility floor (S0 D-003 §8 · F-002)：模态焦点进入/约束/恢复
 * 与 Escape 关闭成立。实现：
 * - 打开时记录触发元素，进入后聚焦容器内首个可聚焦元素；
 * - Tab / Shift+Tab 在容器内循环（focus trap），焦点不逃逸到页面背景；
 * - Escape 触发 onClose；
 * - 卸载时把焦点恢复到打开前的触发元素。
 */

const FOCUSABLE_SELECTOR =
  'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])';

function getFocusable(container: HTMLElement): HTMLElement[] {
  // The selector already excludes disabled and tabindex="-1" elements. No
  // offsetParent/layout check here: jsdom has no layout (offsetParent is always
  // null), so a layout-based visibility filter would return an empty list in
  // tests and break the focus trap. The modal is always rendered on open.
  return Array.from(container.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR));
}

export function ModalHost({
  title,
  onClose,
  children,
}: {
  title?: string;
  onClose: () => void;
  children: ReactNode;
}) {
  const containerRef = useRef<HTMLDivElement>(null);
  const previouslyFocusedRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    // Focus entry: remember the trigger, move focus into the dialog.
    previouslyFocusedRef.current =
      document.activeElement instanceof HTMLElement ? document.activeElement : null;
    const focusable = getFocusable(container);
    (focusable[0] ?? container).focus();

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        event.stopPropagation();
        onClose();
        return;
      }
      if (event.key !== "Tab") return;
      const items = getFocusable(container);
      if (items.length === 0) return;
      const first = items[0];
      const last = items[items.length - 1];
      const active = document.activeElement;
      if (event.shiftKey) {
        if (active === first || !container.contains(active)) {
          event.preventDefault();
          last.focus();
        }
      } else if (active === last || !container.contains(active)) {
        event.preventDefault();
        first.focus();
      }
    };
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("keydown", handleKeyDown);
      // Focus restore: return to the element that opened the modal.
      previouslyFocusedRef.current?.focus();
    };
  }, [onClose]);

  return (
    <div
      ref={containerRef}
      className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-overlay p-4"
      role="dialog"
      aria-modal="true"
      aria-label={title ?? "Dialog"}
    >
      <div className="mt-10 w-full max-w-lg rounded-lg border border-border bg-card p-6 shadow-lg">
        <div className="mb-4 flex items-center justify-between gap-4">
          <h2 className="text-lg font-semibold text-foreground">{title ?? "Dialog"}</h2>
          <button
            type="button"
            aria-label="Close dialog"
            onClick={onClose}
            className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            ×
          </button>
        </div>
        {children}
      </div>
    </div>
  );
}
