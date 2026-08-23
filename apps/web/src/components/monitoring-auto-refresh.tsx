// System monitoring auto-refresh control (W16-F06): lets the operator poll the
// monitoring surface every 5/10/30 seconds without a full page reload.
// W25 (2026-08-23): the tick now performs a TARGETED refresh of the /status
// statCards only (six cards share one coalesced request), instead of a
// full-page reloadList wave that also refetched the events table and every
// other surface. The events table stays current on manual reload; the
// operator-visible contract "poll the monitoring surface" is unchanged.
import { useEffect, useState } from "react";

import { useTranslate } from "@/i18n/runtime";
import {
  registerCustomComponent,
  type CustomComponentProps,
} from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";

const OPTIONS = [
  { value: 0, labelKey: "monitoringRefresh.off" },
  { value: 5000, labelKey: "monitoringRefresh.5s" },
  { value: 10000, labelKey: "monitoringRefresh.10s" },
  { value: 30000, labelKey: "monitoringRefresh.30s" },
];

/** The display dataSource refreshed by the polling tick (system-monitoring.json). */
const STATUS_SOURCE = "/api/system-monitoring/status";

export function MonitoringAutoRefresh(_props: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const [intervalMs, setIntervalMs] = useState(0);

  useEffect(() => {
    if (intervalMs <= 0) {
      return;
    }
    const id = window.setInterval(() => {
      crud?.refreshList(STATUS_SOURCE);
    }, intervalMs);
    return () => window.clearInterval(id);
  }, [intervalMs, crud]);

  return (
    <div className="flex items-center gap-2 text-sm">
      <span className="text-muted-foreground">{t("monitoringRefresh.label")}</span>
      <select
        aria-label={t("monitoringRefresh.label")}
        value={String(intervalMs)}
        onChange={(event) => setIntervalMs(Number(event.target.value))}
        className="h-8 rounded-md border border-input bg-background px-2 text-sm"
        data-monitoring-refresh
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

registerCustomComponent("monitoring-auto-refresh", MonitoringAutoRefresh);
