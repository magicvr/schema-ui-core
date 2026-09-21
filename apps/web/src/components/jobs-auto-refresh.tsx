// Result-center auto refresh (GOAL-005 R4 · VP-038 C2 implementation; refined by
// GOAL-044 W32 D-001 §2/§3).
//
// The jobs table is a page-level list surface: like every other admin table it
// keeps its rows until something refreshes them. A job's status and progress
// change on the SERVER while the operator watches, so the result center needs an
// explicit refresh affordance. This component is the same shape as the frozen
// monitoring auto-refresh control (off / 5s / 10s / 30s).
//
// W32 changes (the two R4 residuals this component carried):
// - It now ticks `crud.refreshTable(targetTable)` — a refresh of THIS table that
//   does NOT clear the page's table selections. The previous `reloadList()` tick
//   was inert only because the jobs table declares no selection; the seam is now
//   correct rather than accidentally harmless (D-001 §2).
// - It skips a tick entirely when the table holds no row in an `activeStatuses`
//   state, so an operator who leaves refreshing on stops generating requests once
//   every job has settled (D-001 §3). Without a readable status field, or when
//   the table has not published rows yet, it refreshes conservatively.
import { useEffect, useState } from "react";

import { useTranslate } from "@/i18n/runtime";
import {
  registerCustomComponent,
  type CustomComponentProps,
} from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";

const OPTIONS = [
  { value: 0, labelKey: "jobsRefresh.off" },
  { value: 5000, labelKey: "jobsRefresh.5s" },
  { value: 10000, labelKey: "jobsRefresh.10s" },
  { value: 30000, labelKey: "jobsRefresh.30s" },
];

/** Default target/field names, matching the jobs page document. */
const DEFAULT_TABLE = "jobs-table";
const DEFAULT_STATUS_FIELD = "status";

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function stringProp(props: Record<string, unknown> | undefined, key: string, fallback: string): string {
  const value = props?.[key];
  return typeof value === "string" && value !== "" ? value : fallback;
}

/** The statuses that keep the poll alive; [] disables the idle check. */
function activeStatuses(props: Record<string, unknown> | undefined): string[] {
  const raw = props?.activeStatuses;
  if (!Array.isArray(raw)) {
    return [];
  }
  return raw.filter((entry): entry is string => typeof entry === "string" && entry !== "");
}

export function JobsAutoRefresh({ node }: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const [intervalMs, setIntervalMs] = useState(0);

  const props = isRecord(node.props) ? node.props : undefined;
  const targetTable = stringProp(props, "targetTable", DEFAULT_TABLE);
  const statusField = stringProp(props, "statusField", DEFAULT_STATUS_FIELD);
  const active = activeStatuses(props);

  useEffect(() => {
    if (intervalMs <= 0) {
      return;
    }
    const id = window.setInterval(() => {
      // Idle check (D-001 §3): with an explicitly declared active set, only poll
      // while at least one rendered row is still in it. Unknown rows (not
      // published yet / non-table host) refresh conservatively.
      if (active.length > 0) {
        const rows = crud?.tableRows(targetTable);
        if (rows !== undefined) {
          const busy = rows.some((row) => active.includes(String(row[statusField] ?? "")));
          if (!busy) {
            return;
          }
        }
      }
      crud?.refreshTable(targetTable);
    }, intervalMs);
    return () => window.clearInterval(id);
  }, [active, crud, intervalMs, statusField, targetTable]);

  return (
    <div className="flex items-center gap-2 text-sm" data-jobs-refresh>
      <span className="text-muted-foreground">{t("jobsRefresh.label")}</span>
      <select
        aria-label={t("jobsRefresh.label")}
        value={String(intervalMs)}
        onChange={(event) => setIntervalMs(Number(event.target.value))}
        className="h-8 rounded-md border border-input bg-background px-2 text-sm"
        data-jobs-refresh-select
      >
        {OPTIONS.map((option) => (
          <option key={option.value} value={String(option.value)}>
            {t(option.labelKey)}
          </option>
        ))}
      </select>
    </div>
  );
}

registerCustomComponent("jobs-auto-refresh", JobsAutoRefresh);
