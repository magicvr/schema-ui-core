import { useCallback, useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useSchemaCrud } from "@/renderer/render.tsx";
import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent } from "@/renderer/custom-components";

interface PolicyView {
  minLength: number;
  minCategories: number;
  historyDepth: number;
}

const inputClass =
  "h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20";
const buttonClass =
  "inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50";

/**
 * Password-policy configuration block (workspace-019 R3 · GOAL-004 D-001
 * §2): the admin.settings tab extension. GET/PATCH /api/settings/password-
 * policy; the server enforces range validation (8–72 / 0–4 / 0–10).
 */
export function PasswordPolicyTab(_props: unknown) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const [loadState, setLoadState] = useState<"loading" | "ready" | "error">("loading");
  const [minLength, setMinLength] = useState("8");
  const [minCategories, setMinCategories] = useState("0");
  const [historyDepth, setHistoryDepth] = useState("0");
  const [saving, setSaving] = useState(false);
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; message: string } | null>(null);

  const load = useCallback(async () => {
    setLoadState("loading");
    try {
      const response = await fetcher("/api/settings/password-policy");
      if (!response.ok) {
        setLoadState("error");
        return;
      }
      const body = (await response.json()) as PolicyView;
      setMinLength(String(body.minLength));
      setMinCategories(String(body.minCategories));
      setHistoryDepth(String(body.historyDepth));
      setLoadState("ready");
    } catch {
      setLoadState("error");
    }
  }, [fetcher]);

  useEffect(() => {
    void load();
  }, [load]);

  async function save() {
    setSaving(true);
    setFeedback(null);
    try {
      const response = await fetcher("/api/settings/password-policy", {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          minLength: Number(minLength),
          minCategories: Number(minCategories),
          historyDepth: Number(historyDepth),
        }),
      });
      if (!response.ok) {
        const body = (await response.json().catch(() => ({}))) as { message?: string };
        setFeedback({ kind: "error", message: body.message ?? t("invite.policy.saveError") });
        return;
      }
      setFeedback({ kind: "success", message: t("invite.policy.saved") });
    } catch {
      setFeedback({ kind: "error", message: t("invite.policy.saveError") });
    } finally {
      setSaving(false);
    }
  }

  if (loadState === "loading") {
    return (
      <section data-password-policy-tab className="space-y-3 rounded-xl border border-border/70 bg-card/85 p-4">
        <p className="text-sm text-muted-foreground">{t("login.signingIn")}</p>
      </section>
    );
  }
  if (loadState === "error") {
    return (
      <section data-password-policy-tab className="space-y-3 rounded-xl border border-border/70 bg-card/85 p-4">
        <p className="text-sm text-destructive">{t("invite.policy.loadError")}</p>
      </section>
    );
  }
  return (
    <section data-password-policy-tab className="space-y-4 rounded-xl border border-border/70 bg-card/85 p-4">
      <div className="space-y-1">
        <h3 className="text-sm font-semibold">{t("invite.policy.title")}</h3>
        <p className="text-xs text-muted-foreground">{t("invite.policy.description")}</p>
      </div>
      <div className="grid gap-3 sm:grid-cols-3">
        <div className="space-y-1">
          <Label htmlFor="policyMinLength" className="text-xs">{t("invite.policy.minLength")}</Label>
          <Input id="policyMinLength" inputMode="numeric" className={inputClass} value={minLength}
            onChange={(event) => setMinLength(event.target.value)} />
        </div>
        <div className="space-y-1">
          <Label htmlFor="policyMinCategories" className="text-xs">{t("invite.policy.minCategories")}</Label>
          <Input id="policyMinCategories" inputMode="numeric" className={inputClass} value={minCategories}
            onChange={(event) => setMinCategories(event.target.value)} />
        </div>
        <div className="space-y-1">
          <Label htmlFor="policyHistoryDepth" className="text-xs">{t("invite.policy.historyDepth")}</Label>
          <Input id="policyHistoryDepth" inputMode="numeric" className={inputClass} value={historyDepth}
            onChange={(event) => setHistoryDepth(event.target.value)} />
        </div>
      </div>
      {feedback !== null ? (
        <p role="status" data-policy-feedback
          className={feedback.kind === "success" ? "text-sm text-success" : "text-sm text-destructive"}>
          {feedback.message}
        </p>
      ) : null}
      <Button type="button" data-policy-save disabled={saving} className={buttonClass} onClick={() => void save()}>
        {saving ? t("login.signingIn") : t("invite.policy.save")}
      </Button>
    </section>
  );
}

// Self-registration under the schema-declared component key (W25 guard).
registerCustomComponent("password-policy-tab", PasswordPolicyTab);
