// Data-permission scope assignment editor (W14 F-02 · GOAL-016): the
// schema-driven data-permission page previously only registered policies;
// this custom component gives an admin a user-facing entry to set each user's
// resource scope (all/self). It reuses the existing scopes GET/PATCH endpoints
// and the users list endpoint for the user picker.
import { useCallback, useEffect, useMemo, useState } from "react";

import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent, type CustomComponentProps } from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";

interface UserRow {
  id: string;
  username: string;
  name?: string;
}

interface PolicyRow {
  resource: string;
  ownerColumn: string;
  defaultScope: string;
  enabled: boolean;
}

interface ScopeRow {
  resource: string;
  scopeType: string;
}

interface ListResponse<T> {
  items?: T[];
  total?: number;
}

const SCOPE_OPTIONS = [
  { value: "all", labelKey: "schema.dataPermission.scope.all" },
  { value: "self", labelKey: "schema.dataPermission.scope.self" },
];

export function DataPermissionScopes(_props: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const [users, setUsers] = useState<UserRow[]>([]);
  const [policies, setPolicies] = useState<PolicyRow[]>([]);
  const [selectedUserId, setSelectedUserId] = useState("");
  const [scopes, setScopes] = useState<Record<string, string>>({});
  const [status, setStatus] = useState<"idle" | "loading" | "ready" | "error">("idle");
  const [saving, setSaving] = useState(false);
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; message: string } | null>(null);

  const loadUsers = useCallback(async () => {
    try {
      const response = await fetcher("/api/users?pageSize=100");
      if (!response.ok) {
        return;
      }
      const body = (await response.json()) as ListResponse<UserRow>;
      setUsers(Array.isArray(body.items) ? body.items : []);
    } catch {
      // Users are optional for the editor; the page still shows policies.
      setUsers([]);
    }
  }, [fetcher]);

  const loadPolicies = useCallback(async () => {
    try {
      const response = await fetcher("/api/data-permission/policies");
      if (!response.ok) {
        setStatus("error");
        return;
      }
      const body = (await response.json()) as ListResponse<PolicyRow>;
      const rows = Array.isArray(body.items) ? body.items : [];
      setPolicies(rows);
      setScopes((prev) => {
        const next: Record<string, string> = {};
        for (const policy of rows) {
          next[policy.resource] = prev[policy.resource] ?? policy.defaultScope;
        }
        return next;
      });
      setStatus("ready");
    } catch {
      setStatus("error");
    }
  }, [fetcher]);

  const loadScopesForUser = useCallback(
    async (userId: string) => {
      if (userId === "") {
        return;
      }
      try {
        const response = await fetcher("/api/data-permission/scopes?userId=" + encodeURIComponent(userId));
        if (!response.ok) {
          setFeedback({ kind: "error", message: t("schema.dataPermission.scopes.loadError") });
          return;
        }
        const body = (await response.json()) as { items?: ScopeRow[] };
        const assignments = new Map<string, string>();
        for (const item of Array.isArray(body.items) ? body.items : []) {
          if (typeof item.resource === "string" && typeof item.scopeType === "string") {
            assignments.set(item.resource, item.scopeType);
          }
        }
        const next: Record<string, string> = {};
        for (const policy of policies) {
          next[policy.resource] = assignments.get(policy.resource) ?? policy.defaultScope;
        }
        setScopes(next);
        setFeedback(null);
      } catch {
        setFeedback({ kind: "error", message: t("schema.dataPermission.scopes.loadError") });
      }
    },
    [fetcher, policies, t],
  );

  useEffect(() => {
    void loadUsers();
    void loadPolicies();
  }, [loadUsers, loadPolicies]);

  useEffect(() => {
    if (selectedUserId !== "") {
      void loadScopesForUser(selectedUserId);
    }
  }, [selectedUserId, loadScopesForUser]);

  const save = useCallback(async () => {
    if (selectedUserId === "") {
      setFeedback({ kind: "error", message: t("schema.dataPermission.scopes.selectUserFirst") });
      return;
    }
    setSaving(true);
    setFeedback(null);
    try {
      const response = await fetcher("/api/data-permission/scopes", {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userId: selectedUserId, scopes }),
      });
      if (!response.ok) {
        const body = (await response.json().catch(() => null)) as { error?: string; message?: string } | null;
        setFeedback({
          kind: "error",
          message: body?.message ?? t("schema.dataPermission.scopes.saveError"),
        });
        return;
      }
      setFeedback({ kind: "success", message: t("schema.dataPermission.scopes.saved") });
    } catch {
      setFeedback({ kind: "error", message: t("schema.dataPermission.scopes.saveError") });
    } finally {
      setSaving(false);
    }
  }, [fetcher, selectedUserId, scopes, t]);

  const selectedUserLabel = useMemo(() => {
    const user = users.find((entry) => entry.id === selectedUserId);
    if (user === undefined) {
      return "";
    }
    return user.name && user.name !== "" ? `${user.username} (${user.name})` : user.username;
  }, [users, selectedUserId]);

  if (status === "loading" || status === "idle") {
    return (
      <section className="space-y-3 rounded-xl border border-border/70 bg-card/85 p-4">
        <h2 className="text-sm font-semibold">{t("schema.dataPermission.scopes.title")}</h2>
        <p className="text-sm text-muted-foreground">{t("feedback.loading")}</p>
      </section>
    );
  }

  if (status === "error") {
    return (
      <section className="space-y-3 rounded-xl border border-border/70 bg-card/85 p-4">
        <h2 className="text-sm font-semibold">{t("schema.dataPermission.scopes.title")}</h2>
        <p role="alert" className="text-sm text-destructive">{t("schema.dataPermission.scopes.loadError")}</p>
      </section>
    );
  }

  return (
    <section data-data-permission-scopes className="space-y-3 rounded-xl border border-border/70 bg-card/85 p-4">
      <div className="flex items-center justify-between gap-3">
        <h2 className="text-sm font-semibold">{t("schema.dataPermission.scopes.title")}</h2>
        {selectedUserId !== "" ? <p className="text-xs text-muted-foreground">{selectedUserLabel}</p> : null}
      </div>

      <div className="grid gap-2 sm:grid-cols-2">
        <label className="text-sm font-medium" htmlFor="data-permission-user">
          {t("schema.dataPermission.scopes.user")}
        </label>
        <select
          id="data-permission-user"
          value={selectedUserId}
          onChange={(event) => setSelectedUserId(event.target.value)}
          className="h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20"
        >
          <option value="">{t("schema.dataPermission.scopes.selectUser")}</option>
          {users.map((user) => (
            <option key={user.id} value={user.id}>
              {user.name && user.name !== "" ? `${user.username} (${user.name})` : user.username}
            </option>
          ))}
        </select>
      </div>

      {selectedUserId === "" ? (
        <p className="text-sm text-muted-foreground">{t("schema.dataPermission.scopes.selectUserPrompt")}</p>
      ) : policies.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t("schema.dataPermission.scopes.noPolicies")}</p>
      ) : (
        <div className="overflow-x-auto rounded-lg border border-border/60">
          <table className="w-full text-sm">
            <thead className="bg-muted/50 text-left text-xs uppercase tracking-wide text-muted-foreground">
              <tr>
                <th scope="col" className="px-3 py-2">{t("schema.dataPermission.column.resource")}</th>
                <th scope="col" className="px-3 py-2">{t("schema.dataPermission.column.ownerColumn")}</th>
                <th scope="col" className="px-3 py-2">{t("schema.dataPermission.column.defaultScope")}</th>
                <th scope="col" className="px-3 py-2">{t("schema.dataPermission.scopes.scope")}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {policies.map((policy) => {
                const labelKey =
                  scopes[policy.resource] === "self"
                    ? "schema.dataPermission.scope.self"
                    : "schema.dataPermission.scope.all";
                return (
                  <tr key={policy.resource}>
                    <td className="px-3 py-2 font-medium">{policy.resource}</td>
                    <td className="px-3 py-2">{policy.ownerColumn}</td>
                    <td className="px-3 py-2">{t(policy.defaultScope === "self" ? "schema.dataPermission.scope.self" : "schema.dataPermission.scope.all")}</td>
                    <td className="px-3 py-2">
                      <select
                        aria-label={t("schema.dataPermission.scopes.scopeFor", { resource: policy.resource })}
                        value={scopes[policy.resource] ?? policy.defaultScope}
                        onChange={(event) =>
                          setScopes((prev) => ({ ...prev, [policy.resource]: event.target.value }))
                        }
                        className="h-9 rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20"
                      >
                        {SCOPE_OPTIONS.map((option) => (
                          <option key={option.value} value={option.value}>
                            {t(option.labelKey)}
                          </option>
                        ))}
                      </select>
                      <span className="sr-only">{t(labelKey)}</span>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {feedback !== null ? (
        <p role={feedback.kind === "error" ? "alert" : "status"} className={"text-sm " + (feedback.kind === "error" ? "text-destructive" : "text-emerald-600")}>
          {feedback.message}
        </p>
      ) : null}

      <button
        type="button"
        disabled={selectedUserId === "" || saving}
        onClick={() => void save()}
        className="inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {saving ? t("feedback.submitting") : t("schema.dataPermission.scopes.save")}
      </button>
    </section>
  );
}

registerCustomComponent("data-permission-scopes", DataPermissionScopes);
