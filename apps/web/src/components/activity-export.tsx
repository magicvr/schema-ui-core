// Activity log CSV export (W14 F-03 · GOAL-016): a small custom component on
// the Activity page that downloads /api/operations/export with the same
// structured filters the user has applied in the schema search form.
import { useCallback, useState } from "react";

import { Download } from "lucide-react";

import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent, type CustomComponentProps } from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";
import { buildResourceQuery } from "@/renderer/resource";

const EXPORT_PAGE_SIZE = 10000;

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function ActivityExport({ node }: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;
  const [exporting, setExporting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const targetTable =
    isRecord(node.props) && typeof node.props.targetTable === "string" && node.props.targetTable !== ""
      ? node.props.targetTable
      : "operations-table";

  const exportCsv = useCallback(async () => {
    setExporting(true);
    setError(null);
    try {
      const query = crud?.tableQuery(targetTable) ?? { page: 1, pageSize: 10 };
      const params = new URLSearchParams(
        buildResourceQuery({
          ...query,
          page: 1,
          pageSize: EXPORT_PAGE_SIZE,
        }),
      );
      const url = "/api/operations/export" + (params.size === 0 ? "" : "?" + params.toString());
      const response = await fetcher(url, { headers: { Accept: "text/csv" } });
      if (!response.ok) {
        const body = (await response.json().catch(() => null)) as { message?: string } | null;
        setError(body?.message ?? t("schema.activity.export.error"));
        return;
      }
      const blob = await response.blob();
      const objectUrl = URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = objectUrl;
      anchor.download = "operations.csv";
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();
      URL.revokeObjectURL(objectUrl);
    } catch {
      setError(t("schema.activity.export.error"));
    } finally {
      setExporting(false);
    }
  }, [crud, fetcher, targetTable, t]);

  return (
    <div className="flex items-center justify-end">
      <button
        type="button"
        disabled={exporting}
        onClick={() => void exportCsv()}
        className="inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md border border-input/80 bg-background px-3.5 text-sm font-medium text-muted-foreground shadow-2xs transition-all duration-150 hover:border-muted-foreground/30 hover:bg-accent/40 hover:text-foreground disabled:cursor-not-allowed disabled:opacity-50"
      >
        <Download aria-hidden="true" className="size-3.5" />
        {exporting ? t("feedback.submitting") : t("schema.activity.export.button")}
      </button>
      {error !== null ? (
        <p role="alert" className="text-sm text-destructive">{error}</p>
      ) : null}
    </div>
  );
}

registerCustomComponent("activity-export", ActivityExport);
