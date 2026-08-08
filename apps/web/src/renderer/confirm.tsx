/**
 * Confirm dialog for row actions that declare a `confirm` string (S4 · GOAL-007).
 *
 * One-time renderer completion allowed by I-007-003 §9.5 (`confirm*.tsx`). The
 * message text comes from the schema row action entry (`table.props.actions[].
 * confirm`), so the dialog stays fixture-driven; cancelling mirrors the frozen
 * `executeAction` CONFIRM_CANCELLED path and never issues a request.
 */

export function ConfirmDialog({
  message,
  onConfirm,
  onCancel,
}: {
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
}) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-overlay p-4"
      role="dialog"
      aria-modal="true"
      aria-label="Confirm action"
    >
      <div className="mt-24 w-full max-w-sm rounded-lg border border-border bg-card p-6 shadow-md">
        <p className="text-sm leading-6 text-foreground">{message}</p>
        <div className="mt-5 flex justify-end gap-2">
          <button
            type="button"
            onClick={onCancel}
            className="h-9 rounded-md border border-input bg-background px-3 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className="h-9 rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground transition-opacity hover:opacity-90"
          >
            Confirm
          </button>
        </div>
      </div>
    </div>
  );
}
