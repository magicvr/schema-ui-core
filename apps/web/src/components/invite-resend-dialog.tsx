import { useCallback, useState } from "react";

import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { useSchemaCrud } from "@/renderer/render.tsx";
import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent, type CustomComponentProps } from "@/renderer/custom-components";

/**
 * Resend dialog (workspace-019, protocol modal action): the row-level
 * "Resend" opens this modal with the triggering invite in context.modalRow
 * (renderer injects it for modal content). Sending rotates the token via the
 * admin API and discloses the ONE-TIME link inline — the schema engine cannot
 * stream response bodies into the page, so the modal is the interaction.
 */
export function InviteResendDialog(props: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const row = props.context?.modalRow as Record<string, unknown> | null | undefined;
  const id = typeof row?.id === "string" ? row.id : "";
  const email = typeof row?.email === "string" && row.email !== "" ? row.email : "";
  const roles = Array.isArray(row?.roles) ? (row.roles as unknown[]).map(String).join(", ") : "";

  const [state, setState] = useState<"idle" | "sending" | "done" | "error">(id === "" ? "error" : "idle");
  const [link, setLink] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [copied, setCopied] = useState(false);

  const resend = useCallback(async () => {
    setState("sending");
    setCopied(false);
    try {
      const response = await fetcher(`/api/users/invites/${id}/resend`, { method: "POST" });
      const body = (await response.json().catch(() => ({}))) as { link?: string; message?: string };
      if (!response.ok) {
        setState("error");
        setErrorMessage(body.message ?? t("invite.panel.actionError"));
        return;
      }
      if (typeof body.link !== "string" || body.link === "") {
        setState("error");
        setErrorMessage(t("invite.resendDialog.noLink"));
        return;
      }
      setLink(body.link);
      setState("done");
    } catch {
      setState("error");
      setErrorMessage(t("invite.panel.actionError"));
    }
  }, [fetcher, id, t]);

  const copy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(link);
      setCopied(true);
    } catch {
      setErrorMessage(t("invite.resendDialog.copyFailed"));
    }
  }, [link, t]);

  const finish = useCallback(() => {
    crud?.closeModal?.();
    crud?.reloadList?.();
  }, [crud]);

  return (
    <div data-invite-resend-dialog className="space-y-4">
      <dl className="space-y-1.5 text-sm">
        <div className="flex justify-between gap-3">
          <dt className="text-muted-foreground">{t("schema.users.column.id")}</dt>
          <dd className="font-mono">{id !== "" ? id : "—"}</dd>
        </div>
        {email !== "" ? (
          <div className="flex justify-between gap-3">
            <dt className="text-muted-foreground">{t("invite.panel.email")}</dt>
            <dd>{email}</dd>
          </div>
        ) : null}
        {roles !== "" ? (
          <div className="flex justify-between gap-3">
            <dt className="text-muted-foreground">{t("schema.users.column.roles")}</dt>
            <dd>{roles}</dd>
          </div>
        ) : null}
      </dl>

      {state === "done" ? (
        <div className="space-y-3">
          <Label className="text-xs">{t("invite.resendDialog.newLink")}</Label>
          <p data-resend-link className="break-all rounded-md bg-muted px-3 py-2 font-mono text-xs">
            {link}
          </p>
          <div className="flex justify-end gap-2">
            {copied ? <span data-copied-hint className="self-center text-xs text-success">{t("invite.resendDialog.copied")}</span> : null}
            <Button type="button" variant="outline" data-resend-copy disabled={copied} onClick={() => void copy()}>
              {t("invite.resendDialog.copy")}
            </Button>
            <Button type="button" data-resend-done onClick={finish}>
              {t("invite.resendDialog.done")}
            </Button>
          </div>
        </div>
      ) : (
        <div className="space-y-3">
          <p className="text-sm text-muted-foreground">
            {id === ""
              ? t("invite.resendDialog.noRow")
              : t("invite.resendDialog.confirm", { email: email !== "" ? email : t("invite.panel.filterAll") })}
          </p>
          {state === "error" ? (
            <p role="alert" data-resend-error className="text-sm text-destructive">
              {errorMessage}
            </p>
          ) : null}
          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={crud?.closeModal ?? (() => undefined)}>
              {t("feedback.cancel")}
            </Button>
            <Button type="button" data-resend-send disabled={state === "sending" || id === ""} onClick={() => void resend()}>
              {state === "sending" ? t("login.signingIn") : t("invite.panel.resend")}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

// Self-registration under the schema-declared component key (W25 guard).
registerCustomComponent("invite-resend-dialog", InviteResendDialog);