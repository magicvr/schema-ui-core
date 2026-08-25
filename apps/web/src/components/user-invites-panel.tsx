import { useCallback, useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useSchemaCrud } from "@/renderer/render.tsx";
import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent } from "@/renderer/custom-components";

interface InviteRow {
  id: string;
  roles: string[];
  invitedBy?: string;
  email?: string;
  expiresAt?: string;
  status: string;
}

const inputClass =
  "h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20";
const buttonClass =
  "inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50";

/**
 * Invitation management panel (workspace-019 R3 · GOAL-004 D-001 §3): the
 * admin.users page block — issue (roles fixed at issuance per user
 * adjudication), list, revoke, resend. The raw token/link is disclosed ONCE
 * per issuance/resend in the response panel.
 */
export function UserInvitesPanel(_props: unknown) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const [rows, setRows] = useState<InviteRow[]>([]);
  const [loadState, setLoadState] = useState<"loading" | "ready" | "error">("loading");
  const [email, setEmail] = useState("");
  const [roles, setRoles] = useState("viewer");
  const [days, setDays] = useState("7");
  const [disclosed, setDisclosed] = useState<{ link: string } | null>(null);
  const [busy, setBusy] = useState(false);
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; message: string } | null>(null);

  const load = useCallback(async () => {
    try {
      const response = await fetcher("/api/users/invites?page=1&pageSize=50");
      if (!response.ok) {
        setLoadState("error");
        return;
      }
      const body = (await response.json()) as { items?: InviteRow[] };
      setRows(Array.isArray(body.items) ? body.items : []);
      setLoadState("ready");
    } catch {
      setLoadState("error");
    }
  }, [fetcher]);

  useEffect(() => {
    void load();
  }, [load]);

  async function create() {
    setBusy(true);
    setFeedback(null);
    setDisclosed(null);
    try {
      const roleKeys = roles.split(",").map((s) => s.trim()).filter(Boolean);
      const response = await fetcher("/api/users/invites", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: email.trim() === "" ? undefined : email.trim(),
          roles: roleKeys,
          expiresInDays: Number(days) || 7,
        }),
      });
      const body = (await response.json().catch(() => ({}))) as { link?: string; message?: string; error?: string };
      if (!response.ok) {
        setFeedback({ kind: "error", message: body.message ?? t("invite.panel.actionError") });
        return;
      }
      if (typeof body.link === "string") {
        setDisclosed({ link: body.link });
      }
      setEmail("");
      await load();
    } catch {
      setFeedback({ kind: "error", message: t("invite.panel.actionError") });
    } finally {
      setBusy(false);
    }
  }

  async function act(id: string, action: "revoke" | "resend") {
    setBusy(true);
    setFeedback(null);
    setDisclosed(null);
    try {
      const response = await fetcher(`/api/users/invites/${id}${action === "resend" ? "/resend" : ""}`, {
        method: action === "revoke" ? "DELETE" : "POST",
      });
      const body = (await response.json().catch(() => ({}))) as { link?: string };
      if (!response.ok) {
        setFeedback({ kind: "error", message: t("invite.panel.actionError") });
        return;
      }
      if (action === "resend" && typeof body.link === "string") {
        setDisclosed({ link: body.link });
      }
      await load();
    } catch {
      setFeedback({ kind: "error", message: t("invite.panel.actionError") });
    } finally {
      setBusy(false);
    }
  }

  return (
    <section data-user-invites-panel className="space-y-4 rounded-xl border border-border/70 bg-card/85 p-4">
      <div className="space-y-1">
        <h3 className="text-sm font-semibold">{t("invite.panel.title")}</h3>
        <p className="text-xs text-muted-foreground">{t("invite.panel.description")}</p>
      </div>
      <div className="grid gap-3 sm:grid-cols-4">
        <div className="space-y-1">
          <Label htmlFor="inviteEmail" className="text-xs">{t("invite.panel.email")}</Label>
          <Input id="inviteEmail" className={inputClass} value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label htmlFor="inviteRoles" className="text-xs">{t("invite.panel.roles")}</Label>
          <Input id="inviteRoles" className={inputClass} value={roles} onChange={(e) => setRoles(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label htmlFor="inviteDays" className="text-xs">{t("invite.panel.days")}</Label>
          <Input id="inviteDays" inputMode="numeric" className={inputClass} value={days} onChange={(e) => setDays(e.target.value)} />
        </div>
        <div className="flex items-end">
          <Button type="button" data-invite-create disabled={busy || roles.trim() === ""} className={`${buttonClass} w-full`}
            onClick={() => void create()}>
            {t("invite.panel.create")}
          </Button>
        </div>
      </div>
      {disclosed !== null ? (
        <p role="status" data-invite-link className="break-all rounded-md bg-muted px-3 py-2 text-xs">
          {t("invite.panel.disclose")}{" "}
          <code className="font-mono">{disclosed.link}</code>
        </p>
      ) : null}
      {feedback !== null ? (
        <p role="alert" data-invite-feedback className={feedback.kind === "success" ? "text-sm text-success" : "text-sm text-destructive"}>
          {feedback.message}
        </p>
      ) : null}
      {loadState === "loading" ? (
        <p className="text-sm text-muted-foreground">{t("login.signingIn")}</p>
      ) : loadState === "error" ? (
        <p className="text-sm text-destructive">{t("invite.panel.loadError")}</p>
      ) : rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t("invite.panel.empty")}</p>
      ) : (
        <table data-invite-table className="w-full text-left text-xs">
          <thead>
            <tr className="border-b border-border/60">
              <th className="py-2 pr-2 font-medium">{t("schema.users.column.roles")}</th>
              <th className="py-2 pr-2 font-medium">{t("invite.panel.email")}</th>
              <th className="py-2 pr-2 font-medium">{t("invite.panel.status")}</th>
              <th className="py-2 pr-2 font-medium">{t("invite.panel.expires")}</th>
              <th className="py-2 font-medium">{t("invite.panel.actions")}</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.id} className="border-b border-border/40 last:border-b-0">
                <td className="py-2 pr-2 font-mono">{row.roles.join(", ")}</td>
                <td className="py-2 pr-2">{row.email ?? "—"}</td>
                <td className="py-2 pr-2" data-invite-status>{row.status}</td>
                <td className="py-2 pr-2">{row.expiresAt?.slice(0, 10) ?? "—"}</td>
                <td className="py-2">
                  {row.status === "pending" ? (
                    <span className="flex gap-2">
                      <Button type="button" variant="outline" disabled={busy} onClick={() => void act(row.id, "resend")}>
                        {t("invite.panel.resend")}
                      </Button>
                      <Button type="button" variant="outline" disabled={busy} onClick={() => void act(row.id, "revoke")}>
                        {t("invite.panel.revoke")}
                      </Button>
                    </span>
                  ) : (
                    "—"
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
}

// Self-registration under the schema-declared component key (W25 guard).
registerCustomComponent("user-invites-panel", UserInvitesPanel);
