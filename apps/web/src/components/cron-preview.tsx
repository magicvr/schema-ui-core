// Cron expression preview (W16-F05): accepts a 5-field cron expression and
// shows the server-computed description and next three run times.
import { useCallback, useEffect, useState, type FormEvent } from "react";

import { useAuth } from "@/account/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useTranslate } from "@/i18n/runtime";
import {
  registerCustomComponent,
  type CustomComponentProps,
} from "@/renderer/custom-components";

interface CronPreviewResponse {
  description?: string;
  nextRuns?: string[];
}

export function CronPreview(_props: CustomComponentProps) {
  const t = useTranslate();
  const { authFetch } = useAuth();
  const [cron, setCron] = useState("");
  const [preview, setPreview] = useState<CronPreviewResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const runPreview = useCallback(
    async (expr: string) => {
      const trimmed = expr.trim();
      if (trimmed === "") {
        setPreview(null);
        setError(null);
        return;
      }
      setBusy(true);
      setError(null);
      try {
        const response = await authFetch("/api/scheduled-tasks/cron/preview", {
          method: "POST",
          headers: { "Content-Type": "application/json", Accept: "application/json" },
          body: JSON.stringify({ cron: trimmed }),
        });
        if (!response.ok) {
          setError(t("cronPreview.invalid"));
          return;
        }
        setPreview((await response.json()) as CronPreviewResponse);
      } catch {
        setError(t("cronPreview.failed"));
      } finally {
        setBusy(false);
      }
    },
    [authFetch, t],
  );

  // W16-F05: debounced auto-preview while typing.
  useEffect(() => {
    const id = window.setTimeout(() => {
      void runPreview(cron);
    }, 400);
    return () => window.clearTimeout(id);
  }, [cron, runPreview]);

  const submit = (event: FormEvent) => {
    event.preventDefault();
    void runPreview(cron);
  };

  return (
    <div className="space-y-2 rounded-md border p-3" data-cron-preview>
      <h3 className="text-base font-semibold">{t("cronPreview.title")}</h3>
      <form onSubmit={submit} className="flex items-end gap-2">
        <div className="min-w-0 flex-1 space-y-1.5">
          <Label htmlFor="cronPreviewInput">{t("cronPreview.cron")}</Label>
          <Input
            id="cronPreviewInput"
            value={cron}
            onChange={(event) => setCron(event.target.value)}
            placeholder="0 0 2 * * *"
          />
        </div>
        <Button type="submit" disabled={busy || cron.trim() === ""}>
          {busy ? t("cronPreview.loading") : t("cronPreview.preview")}
        </Button>
      </form>
      {error !== null ? <p className="text-sm text-destructive">{error}</p> : null}
      {preview !== null ? (
        <div className="space-y-1 text-sm">
          <p>
            <span className="text-muted-foreground">{t("cronPreview.description")}:</span>{" "}
            {preview.description ?? "—"}
          </p>
          <p className="text-muted-foreground">{t("cronPreview.nextRuns")}:</p>
          <ul className="list-inside list-disc font-mono text-xs">
            {(preview.nextRuns ?? []).map((run) => (
              <li key={run}>{run}</li>
            ))}
          </ul>
        </div>
      ) : null}
    </div>
  );
}

registerCustomComponent("cron-preview", CronPreview);
