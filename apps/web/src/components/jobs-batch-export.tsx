// Async batch export (GOAL-004 R3 · VP-038 exit criterion 3): the first real
// batch operation carried by the durable Job runtime.
//
// Design notes (GOAL-004 D-001 §2):
// - The selection comes from `useSchemaCrud().selection(tableId)`, NOT from the
//   custom-node `context` (which carries only host/nav data and no table id).
//   The target table is therefore declared on the node as `props.targetTable`.
// - The submit deliberately does NOT go through ADR-0022's `runBatchRequest`:
//   that path only checks `response.ok` and then calls `reloadList()`, which
//   would discard the 202 body (the jobId) and wipe the selection.
// - It never calls `reloadList()` here either: any successful reload clears
//   every table selection, and the user still needs the selection while the
//   async job runs.
// - The table must declare `props.selection.mode = "multiple"`, otherwise the
//   renderer never publishes a selection for it.
import { useCallback, useEffect, useRef, useState } from "react";

import { Download, Loader2 } from "lucide-react";

import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent, type CustomComponentProps } from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";

// Poll cadence. D-001 §2.4 froze the 5/10/30s family used by
// monitoringAutoRefresh; a batch export settles in seconds, so a flat 5s start
// would feel unresponsive. The initial 1s step is a bounded deviation: it
// exists only for the first few ticks, after which the loop settles onto the
// frozen 5s cadence (A-002 F-003). Every value stays inside the frozen family's
// lower bound once the job is past its first seconds.
const POLL_FAST_MS = 1000;
const POLL_STEADY_MS = 5000;
const POLL_FAST_TICKS = 5;

function pollDelayMs(tick: number): number {
  return tick < POLL_FAST_TICKS ? POLL_FAST_MS : POLL_STEADY_MS;
}

/** Terminal statuses: stop polling once the job settles. */
const TERMINAL = new Set(["succeeded", "failed", "cancelled", "expired"]);

interface JobProjection {
  id?: string;
  status?: string;
  progress?: number;
  error?: { code?: string; message?: string };
  resultUrl?: string;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function stringProp(props: unknown, key: string, fallback: string): string {
  if (isRecord(props) && typeof props[key] === "string" && props[key] !== "") {
    return props[key] as string;
  }
  return fallback;
}

export function JobsBatchExport({ node }: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const targetTable = stringProp(node.props, "targetTable", "users-table");
  const resource = stringProp(node.props, "resource", "users");

  const selection = crud?.selection(targetTable);
  const keys = selection?.keys ?? [];
  const count = selection?.count ?? 0;

  const [submitting, setSubmitting] = useState(false);
  const [job, setJob] = useState<JobProjection | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [downloaded, setDownloaded] = useState(false);
  const timer = useRef<number | null>(null);
  const tickCount = useRef(0);

  // Poll the submitted job until it settles. The job id is the only state we
  // need; the read route is the R2 management-scope surface. A single
  // self-rescheduling timeout (rather than setInterval) lets the delay grow
  // from the fast first ticks onto the frozen steady cadence.
  useEffect(() => {
    const jobId = job?.id;
    if (jobId === undefined || (job?.status !== undefined && TERMINAL.has(job.status))) {
      return;
    }
    let cancelled = false;
    const schedule = () => {
      if (cancelled) {
        return;
      }
      timer.current = window.setTimeout(() => {
        void tick();
      }, pollDelayMs(tickCount.current));
    };
    const tick = async () => {
      tickCount.current += 1;
      try {
        const response = await fetcher(`/api/jobs/${jobId}`, { headers: { Accept: "application/json" } });
        if (response.ok) {
          const body = (await response.json()) as JobProjection;
          if (!cancelled) {
            setJob((current) => ({ ...current, ...body }));
          }
        }
      } catch {
        // A transient poll failure is not fatal: the loop reschedules and the
        // job keeps running server-side regardless of this component.
      }
      schedule();
    };
    schedule();
    return () => {
      cancelled = true;
      if (timer.current !== null) {
        window.clearTimeout(timer.current);
        timer.current = null;
      }
    };
  }, [fetcher, job?.id, job?.status]);

  const submit = useCallback(async () => {
    if (keys.length === 0) {
      return;
    }
    setSubmitting(true);
    setError(null);
    setDownloaded(false);
    try {
      const response = await fetcher("/api/jobs/batch-export", {
        method: "POST",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify({ resource, ids: keys }),
      });
      const body = (await response.json().catch(() => null)) as
        | (JobProjection & { message?: string })
        | null;
      if (response.status !== 202 || body?.id === undefined) {
        setError(body?.message ?? t("schema.jobs.batchExport.error"));
        return;
      }
      // 202 + jobId: hand over to the poller. No reloadList() — it would clear
      // the selection the user is still looking at.
      setJob({ id: body.id, status: body.status ?? "queued", progress: body.progress ?? 0 });
    } catch {
      setError(t("schema.jobs.batchExport.error"));
    } finally {
      setSubmitting(false);
    }
  }, [fetcher, keys, resource, t]);

  const download = useCallback(async () => {
    const jobId = job?.id;
    if (jobId === undefined) {
      return;
    }
    try {
      const response = await fetcher(`/api/jobs/${jobId}/result`, { headers: { Accept: "application/json" } });
      if (!response.ok) {
        setError(t("schema.jobs.batchExport.error"));
        return;
      }
      const payload = (await response.json()) as { csv?: string; fileName?: string };
      if (typeof payload.csv !== "string") {
        setError(t("schema.jobs.batchExport.error"));
        return;
      }
      const blob = new Blob([payload.csv], { type: "text/csv;charset=utf-8" });
      const objectUrl = URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = objectUrl;
      anchor.download = payload.fileName ?? `${resource}-selection.csv`;
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();
      URL.revokeObjectURL(objectUrl);
      setDownloaded(true);
    } catch {
      setError(t("schema.jobs.batchExport.error"));
    }
  }, [fetcher, job?.id, resource, t]);

  const status = job?.status;
  const running = status !== undefined && !TERMINAL.has(status);
  const succeeded = status === "succeeded";
  const progress = job?.progress ?? 0;

  return (
    <div className="flex flex-wrap items-center gap-3 text-sm" data-jobs-batch-export>
      <button
        type="button"
        onClick={() => void submit()}
        disabled={count === 0 || submitting || running}
        className="inline-flex h-8 items-center gap-2 rounded-md border border-input bg-background px-3 text-sm disabled:cursor-not-allowed disabled:opacity-50"
        data-jobs-batch-export-submit
      >
        {submitting || running ? <Loader2 className="h-4 w-4 animate-spin" aria-hidden /> : null}
        {t("schema.jobs.batchExport.action")}
        {count > 0 ? ` (${count})` : ""}
      </button>

      {running ? (
        <span className="text-muted-foreground" role="status" data-jobs-batch-export-progress>
          {t("schema.jobs.batchExport.running")} {progress}%
        </span>
      ) : null}

      {succeeded ? (
        <button
          type="button"
          onClick={() => void download()}
          className="inline-flex h-8 items-center gap-2 rounded-md border border-input bg-background px-3 text-sm"
          data-jobs-batch-export-download
        >
          <Download className="h-4 w-4" aria-hidden />
          {t("schema.jobs.batchExport.download")}
        </button>
      ) : null}

      {downloaded ? (
        <span className="text-muted-foreground" role="status">
          {t("schema.jobs.batchExport.done")}
        </span>
      ) : null}

      {status !== undefined && !running && !succeeded && status !== undefined ? (
        <span className="text-destructive" role="alert" data-jobs-batch-export-failed>
          {job?.error?.message ?? t("schema.jobs.batchExport.failed")}
        </span>
      ) : null}

      {error !== null ? (
        <span className="text-destructive" role="alert">
          {error}
        </span>
      ) : null}
    </div>
  );
}

registerCustomComponent("jobs-batch-export", JobsBatchExport);
