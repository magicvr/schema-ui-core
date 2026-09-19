// Result-center auto refresh (GOAL-005 R4 · VP-038 C2).
//
// The jobs table is a page-level list surface: like every other admin table it
// keeps its rows until a reload (`refreshList` covers display nodes only —
// statCard/chart — as documented on the SchemaCrudValue seam). A job's status
// and progress change on the SERVER while the operator watches, so the result
// center needs an explicit refresh affordance. This component is the same shape
// as the frozen monitoring auto-refresh control (off / 5s / 10s / 30s) and ticks
// the page-level `reloadList()` seam.
//
// Scope note (GOAL-005 D-001 §5, R-1.2): `reloadList()` clears every table
// selection on the page (ADR-0022 D2). The jobs page declares no
// `props.selection` on its table, so there is no selection to lose — that
// premise is asserted by the interaction test, and if jobs ever gains row
// selection this component must switch to a targeted refresh instead.
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

export function JobsAutoRefresh(_props: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const [intervalMs, setIntervalMs] = useState(0);

  useEffect(() => {
    if (intervalMs <= 0) {
      return;
    }
    const id = window.setInterval(() => {
      crud?.reloadList();
    }, intervalMs);
    return () => window.clearInterval(id);
  }, [intervalMs, crud]);

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
