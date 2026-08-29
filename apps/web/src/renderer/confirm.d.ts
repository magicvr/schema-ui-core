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
export declare function ConfirmDialog({ message, onConfirm, onCancel, }: {
    message: string;
    onConfirm: () => void;
    onCancel: () => void;
}): import("react").JSX.Element;
