import { useCallback, useEffect, useRef, useState } from "react";
import { ChevronDown } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useSchemaCrud } from "@/renderer/render.tsx";
import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent } from "@/renderer/custom-components";

interface RoleOption {
  key: string;
  name: string;
}

const inputClass =
  "h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20";
const buttonClass =
  "inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50";

/**
 * Dropdown multi-select (shared with the legacy panel): a trigger showing
 * selected options as badges, a checkbox popover to toggle, outside-click to
 * close. Role assignment is fixed at issuance per user adjudication.
 */
function RoleMultiSelect({
  options,
  selected,
  onChange,
  disabled,
  placeholder,
}: {
  options: RoleOption[];
  selected: string[];
  onChange: (keys: string[]) => void;
  disabled?: boolean;
  placeholder: string;
}) {
  const t = useTranslate();
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) {
      return;
    }
    const onDown = (event: MouseEvent) => {
      if (rootRef.current !== null && !rootRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", onDown);
    return () => document.removeEventListener("mousedown", onDown);
  }, [open]);

  const toggle = (key: string) => {
    onChange(selected.includes(key) ? selected.filter((k) => k !== key) : [...selected, key]);
  };
  const nameOf = (key: string) => {
    const found = options.find((o) => o.key === key);
    return found !== undefined ? found.name : key;
  };

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        data-role-multiselect-trigger
        disabled={disabled}
        aria-haspopup="listbox"
        aria-expanded={open}
        className="flex h-9 min-h-9 w-full cursor-pointer items-center justify-between gap-2 rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20 disabled:cursor-not-allowed disabled:opacity-50"
        onClick={() => setOpen((v) => !v)}
      >
        <span className="flex min-w-0 flex-wrap items-center gap-1">
          {selected.length === 0 ? (
            <span className="text-muted-foreground">{placeholder}</span>
          ) : (
            selected.map((key) => (
              <span
                key={key}
                data-role-multiselect-badge={key}
                className="inline-flex items-center rounded-md bg-primary/10 px-1.5 py-0.5 text-xs font-medium text-primary"
              >
                {nameOf(key)}
              </span>
            ))
          )}
        </span>
        <ChevronDown aria-hidden="true" className="size-4 shrink-0 text-muted-foreground" />
      </button>
      {open ? (
        <div
          data-role-multiselect-popover
          role="listbox"
          aria-multiselectable="true"
          className="absolute z-20 mt-1 max-h-56 w-full overflow-auto rounded-md border border-border/70 bg-popover py-1 shadow-md"
        >
          {options.length === 0 ? (
            <p className="px-3 py-2 text-xs text-muted-foreground">{t("invite.panel.rolesEmpty")}</p>
          ) : (
            options.map((option) => (
              <label
                key={option.key}
                role="option"
                aria-selected={selected.includes(option.key)}
                className="flex cursor-pointer items-center gap-2 px-3 py-1.5 text-sm transition-colors hover:bg-muted"
              >
                <input
                  type="checkbox"
                  data-role-multiselect-option={option.key}
                  checked={selected.includes(option.key)}
                  onChange={() => toggle(option.key)}
                />
                <span className="min-w-0">{option.name}</span>
              </label>
            ))
          )}
        </div>
      ) : null}
    </div>
  );
}


/**
 * Invitation issue card (workspace-019): the create operation PLUS one-time
 * link disclosure for creation and resend. The records list lives on the
 * schema table of the same page; resend accepts the invite id from the table
 * (column "id") because the schema engine cannot disclose response bodies.
 */
export function InviteIssueCard(_props: unknown) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const [roleOptions, setRoleOptions] = useState<RoleOption[]>([]);
  const [rolesLoaded, setRolesLoaded] = useState(false);
  const [email, setEmail] = useState("");
  const [roles, setRoles] = useState<string[]>(["viewer"]);
  const [days, setDays] = useState("7");
  const [resendId, setResendId] = useState("");
  const [disclosed, setDisclosed] = useState<{ link: string } | null>(null);
  const [busy, setBusy] = useState(false);
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; message: string } | null>(null);

  const loadRoles = useCallback(async () => {
    try {
      const response = await fetcher("/api/roles?pageSize=100");
      if (!response.ok) {
        return;
      }
      const body = (await response.json()) as { items?: RoleOption[] };
      setRoleOptions(Array.isArray(body.items) ? body.items : []);
    } catch {
      // fail-open: create stays disabled until roles are selectable
    } finally {
      setRolesLoaded(true);
    }
  }, [fetcher]);

  useEffect(() => {
    void loadRoles();
  }, [loadRoles]);

  async function create() {
    setBusy(true);
    setFeedback(null);
    setDisclosed(null);
    try {
      const response = await fetcher("/api/users/invites", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: email.trim() === "" ? undefined : email.trim(),
          roles,
          expiresInDays: Number(days) || 7,
        }),
      });
      const body = (await response.json().catch(() => ({}))) as { link?: string; message?: string };
      if (!response.ok) {
        setFeedback({ kind: "error", message: body.message ?? t("invite.panel.actionError") });
        return;
      }
      if (typeof body.link === "string") {
        setDisclosed({ link: body.link });
      }
      setEmail("");
    } catch {
      setFeedback({ kind: "error", message: t("invite.panel.actionError") });
    } finally {
      setBusy(false);
    }
  }

  async function resend() {
    const id = resendId.trim();
    if (id === "") {
      setFeedback({ kind: "error", message: t("invite.card.resendIdRequired") });
      return;
    }
    setBusy(true);
    setFeedback(null);
    setDisclosed(null);
    try {
      const response = await fetcher(`/api/users/invites/${id}/resend`, { method: "POST" });
      const body = (await response.json().catch(() => ({}))) as { link?: string; message?: string };
      if (!response.ok) {
        setFeedback({ kind: "error", message: body.message ?? t("invite.panel.actionError") });
        return;
      }
      if (typeof body.link === "string") {
        setDisclosed({ link: body.link });
      }
      setResendId("");
    } catch {
      setFeedback({ kind: "error", message: t("invite.panel.actionError") });
    } finally {
      setBusy(false);
    }
  }

  return (
    <section data-invite-issue-card className="space-y-4 rounded-xl border border-border/70 bg-card/85 p-4">
      <div className="space-y-1">
        <h3 className="text-sm font-semibold">{t("invite.panel.issueTitle")}</h3>
        <p className="text-xs text-muted-foreground">{t("invite.panel.description")}</p>
      </div>
      <div className="grid gap-3 sm:grid-cols-4">
        <div className="space-y-1">
          <Label htmlFor="inviteEmail" className="text-xs">{t("invite.panel.email")}</Label>
          <Input id="inviteEmail" className={inputClass} value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label htmlFor="inviteRoles" className="text-xs">{t("invite.panel.roles")}</Label>
          <RoleMultiSelect
            options={roleOptions}
            selected={roles}
            onChange={setRoles}
            disabled={!rolesLoaded}
            placeholder={t("invite.panel.rolesPlaceholder")}
          />
        </div>
        <div className="space-y-1">
          <Label htmlFor="inviteDays" className="text-xs">{t("invite.panel.days")}</Label>
          <Input id="inviteDays" inputMode="numeric" className={inputClass} value={days} onChange={(e) => setDays(e.target.value)} />
        </div>
        <div className="flex items-end">
          <Button
            type="button"
            data-invite-create
            disabled={busy || !rolesLoaded || roles.length === 0}
            className={`${buttonClass} w-full`}
            onClick={() => void create()}
          >
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
      {/* Resend: the schema table cannot disclose the rotated link, so the
          invite id (table column) is pasted here and the fresh link lands in
          the disclosure band above (workspace-019 rendering decision). */}
      <div className="flex flex-wrap items-end gap-3 border-t border-border/60 pt-3">
        <div className="min-w-52 flex-1 space-y-1">
          <Label htmlFor="inviteResendId" className="text-xs">{t("invite.card.resendLabel")}</Label>
          <Input
            id="inviteResendId"
            data-invite-resend-id
            className={inputClass}
            value={resendId}
            placeholder={t("invite.card.resendPlaceholder")}
            onChange={(e) => setResendId(e.target.value)}
          />
        </div>
        <Button
          type="button"
          data-invite-resend
          variant="outline"
          disabled={busy || resendId.trim() === ""}
          onClick={() => void resend()}
        >
          {t("invite.panel.resend")}
        </Button>
      </div>
    </section>
  );
}

// Self-registration under the schema-declared component key (W25 guard).
registerCustomComponent("invite-issue-card", InviteIssueCard);