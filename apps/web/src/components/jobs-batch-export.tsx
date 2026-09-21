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
import { downloadJobResultDocument } from "@/lib/job-result-download";
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

/**
 * Whether the signed-in principal can actually use the async batch export.
 *
 * The submit route enforces TWO independent gates server-side — `jobs.write`
 * (authorizes mutating the async runtime) and `data.export` (authorizes moving
 * row data out). Both are contributed by modules that are absent from the
 * `mvp` / `demo` presets (`admin.jobs` contributes jobs.write,
 * `admin.data-transfer` contributes data.export), while the schema node itself
 * ships with the users/roles module and is therefore served in every profile.
 * Without this check the trigger looks usable where its route is not even
 * mounted, and pressing it produces a bare 404 (NOT_FOUND / 「未找到」) instead of
 * an export (user report 2026-09-19). Requiring the same two grants the route
 * requires covers both failure shapes at once: route absent (module not
 * enabled) and caller unauthorized (missing grant).
 *
 * The control is DISABLED rather than hidden. W33 D-001 §3 froze this entry
 * point as fail-open: silently removing an operation entry point is a worse
 * failure mode than showing it unavailable, and hiding it would also leave the
 * slot host rendered empty (the e2e list-surface contract requires every
 * page-action child to carry the same control height).
 *
 * Fail-open on an unknown context: when permissions cannot be read at all (a
 * bare test harness, an older host), the component keeps its pre-existing
 * behaviour rather than disabling a working operation entry point.
 */
function canBatchExport(context: Record<string, unknown> | undefined): boolean {
  if (context === undefined) {
    return true;
  }
  const user = context.user;
  if (!isRecord(user)) {
    return true;
  }
  const permissions = user.permissions;
  if (!Array.isArray(permissions)) {
    return true;
  }
  return permissions.includes("jobs.write") && permissions.includes("data.export");
}

export function JobsBatchExport({ node, context }: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const targetTable = stringProp(node.props, "targetTable", "users-table");
  const resource = stringProp(node.props, "resource", "users");
  const permitted = canBatchExport(context);

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
    if (keys.length === 0 || !permitted) {
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
        // A 404 here is not a data problem: the route only exists when the
        // `admin.jobs` module is enabled. Say so, instead of surfacing the
        // server's bare NOT_FOUND (「未找到」), which reads as a missing row and
        // sent the 2026-09-19 investigation down the wrong path.
        setError(
          response.status === 404
            ? t("schema.jobs.batchExport.unavailable")
            : (body?.message ?? t("schema.jobs.batchExport.error")),
        );
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
  }, [fetcher, keys, permitted, resource, t]);

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
      // R4: the envelope → file decision is shared with the result center's row
      // action (GOAL-005 D-001 §4) so both export entry points produce the same
      // file, under the same name, from the same document.
      const payload: unknown = await response.json();
      downloadJobResultDocument(payload, `${resource}-selection.csv`);
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
        disabled={!permitted || count === 0 || submitting || running}
        title={permitted ? undefined : t("schema.jobs.batchExport.unavailable")}
        className="inline-flex h-8 items-center gap-2 rounded-md border border-input bg-background px-3 text-sm disabled:cursor-not-allowed disabled:opacity-50"
        data-jobs-batch-export-submit
        data-jobs-batch-export-unavailable={permitted ? undefined : "true"}
      >
        {submitting || running ? <Loader2 className="h-4 w-4 animate-spin" aria-hidden /> : null}
        {t("schema.jobs.batchExport.action")}
        {count > 0 ? ` (${count})` : ""}
      </button>

      {!permitted ? (
        <span className="text-muted-foreground" role="status" data-jobs-batch-export-unavailable-note>
          {t("schema.jobs.batchExport.unavailable")}
        </span>
      ) : null}

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
