import type { ReactNode } from "react";

/**
 * Minimal modal host for Schema-driven create/edit forms (S4 · GOAL-007).
 *
 * One-time renderer completion allowed by I-007-003 §9.5 (`modal*.tsx`). The
 * host is deliberately generic: the schema action (`ModalAction.content`)
 * supplies the body, the title comes from the triggering affordance label, and
 * `onClose` is wired to the page's action executor. No record-specific logic
 * lives here — a fixture-only page change never touches this file.
 */

export function ModalHost({
  title,
  onClose,
  children,
}: {
  title?: string;
  onClose: () => void;
  children: ReactNode;
}) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-black/40 p-4"
      role="dialog"
      aria-modal="true"
      aria-label={title ?? "Dialog"}
    >
      <div className="mt-10 w-full max-w-lg rounded-lg border border-border bg-card p-6 shadow-xl">
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
