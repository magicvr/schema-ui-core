import { useEffect, useState } from "react";

import { useTranslate } from "@/i18n/runtime";

export const QUICKSTART_UPGRADE_URL =
  "https://github.com/magicvr/schema-ui-core/blob/main/QUICKSTART.md";

const STATUS_URL = "/api/system-monitoring/status";

export function hasMonitoringRead(permissions: readonly string[] | undefined): boolean {
  return Array.isArray(permissions) && permissions.includes("monitoring.read");
}

export function VersionChip({
  canReadMonitoring,
  fetcher,
  onOpenDiagnostics,
}: {
  canReadMonitoring: boolean;
  fetcher?: typeof fetch;
  onOpenDiagnostics?: () => void;
}) {
  const t = useTranslate();
  const [version, setVersion] = useState<string | null>(null);

  useEffect(() => {
    if (!canReadMonitoring) {
      setVersion(null);
      return;
    }
    const load = fetcher ?? globalThis.fetch;
    let cancelled = false;
    void load(STATUS_URL)
      .then(async (response) => {
        if (!response.ok) {
          return null;
        }
        const body = (await response.json()) as { items?: Array<{ version?: unknown }> };
        const value = body.items?.[0]?.version;
        return typeof value === "string" && value !== "" ? value : null;
      })
      .then((value) => {
        if (!cancelled) {
          setVersion(value);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setVersion(null);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [canReadMonitoring, fetcher]);

  if (!canReadMonitoring || version === null) {
    return null;
  }

  return (
    <div
      data-version-chip
      className="hidden items-center gap-2 text-xs text-muted-foreground sm:flex"
    >
      <span data-version-chip-value aria-label={t("versionChip.label")}>
        {version}
      </span>
      <a
        href={QUICKSTART_UPGRADE_URL}
        target="_blank"
        rel="noopener noreferrer"
        className="underline-offset-2 hover:text-foreground hover:underline"
      >
        {t("versionChip.upgrade")}
      </a>
      {onOpenDiagnostics !== undefined ? (
        <button
          type="button"
          className="underline-offset-2 hover:text-foreground hover:underline"
          onClick={onOpenDiagnostics}
        >
          {t("versionChip.diagnostics")}
        </button>
      ) : null}
    </div>
  );
}
