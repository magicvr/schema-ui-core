import { useCallback, useEffect, useRef, useState } from "react";
import { ChevronDown } from "lucide-react";

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

interface RoleOption {
  key: string;
  name: string;
}

const inputClass =
  "h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20";
const buttonClass =
  "inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50";

/**
 * Dropdown multi-select (workspace-019 UX polish): a trigger showing the
 * selected options as badges, a checkbox popover to toggle, outside-click to
 * close. Used by the invitation panel for role assignment at issuance.
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
    onChange(
      selected.includes(key) ? selected.filter((k) => k !== key) : [...selected, key],
    );
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
  const [roleOptions, setRoleOptions] = useState<RoleOption[]>([]);
  const [rolesLoaded, setRolesLoaded] = useState(false);
  const [email, setEmail] = useState("");
  const [roles, setRoles] = useState<string[]>(["viewer"]);
  const [days, setDays] = useState("7");
  const [disclosed, setDisclosed] = useState<{ link: string } | null>(null);
  const [busy, setBusy] = useState(false);
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; message: string } | null>(null);
  // List filter: server-side narrowing so paged filters never mislead.
  const [filter, setFilter] = useState("all");

  const load = useCallback(async () => {
    try {
      const response = await fetcher(`/api/users/invites?page=1&pageSize=50${filter === "all" ? "" : `&status=${filter}`}`);
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
  }, [fetcher, filter]);

  // Role catalog for the multi-select (same /api/roles source the schema
  // checkboxGroup uses). Fail-open: an empty catalog leaves create disabled.
  const loadRoles = useCallback(async () => {
    try {
      const response = await fetcher("/api/roles?pageSize=100");
      if (!response.ok) {
        return;
      }
      const body = (await response.json()) as { items?: RoleOption[] };
      setRoleOptions(Array.isArray(body.items) ? body.items : []);
    } catch {
      // fail-open
    } finally {
      setRolesLoaded(true);
    }
  }, [fetcher]);

  useEffect(() => {
    void load();
    void loadRoles();
  }, [load, loadRoles]);

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

  const roleName = (key: string) => {
    const found = roleOptions.find((o) => o.key === key);
    return found !== undefined ? found.name : key;
  };

  return (
    <div data-user-invites-panel className="space-y-4">
      {/* Issue card: the create operation lives apart from the records list. */}
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
      </section>

      {/* Records card: filterable, server-side narrowed list. */}
      <section data-invite-records-card className="space-y-4 rounded-xl border border-border/70 bg-card/85 p-4">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div className="space-y-1">
            <h3 className="text-sm font-semibold">{t("invite.panel.recordsTitle")}</h3>
            <p className="text-xs text-muted-foreground">{t("invite.panel.recordsHint")}</p>
          </div>
          <div className="flex items-center gap-2">
            <Label htmlFor="inviteStatusFilter" className="text-xs">{t("invite.panel.filter")}</Label>
            <select
              id="inviteStatusFilter"
              data-invite-status-filter
              value={filter}
              onChange={(event) => setFilter(event.target.value)}
              className="h-9 cursor-pointer rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20"
            >
              <option value="all">{t("invite.panel.filterAll")}</option>
              <option value="pending">{t("invite.panel.filterPending")}</option>
              <option value="consumed">{t("invite.panel.filterConsumed")}</option>
              <option value="revoked">{t("invite.panel.filterRevoked")}</option>
              <option value="expired">{t("invite.panel.filterExpired")}</option>
            </select>
          </div>
        </div>
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
                  <td className="py-2 pr-2">{row.roles.map(roleName).join(", ") || "—"}</td>
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
    </div>
  );
}

// Self-registration under the schema-declared component key (W25 guard).
registerCustomComponent("user-invites-panel", UserInvitesPanel);
