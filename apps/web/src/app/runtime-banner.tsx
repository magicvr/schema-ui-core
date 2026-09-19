import { useTranslate } from "@/i18n/runtime";
import type { RuntimeMode } from "@/account/auth-client";

const BANNER_KEYS: Record<Exclude<RuntimeMode, "normal">, string> = {
  maintenance: "runtimeBanner.maintenance",
  degraded: "runtimeBanner.degraded",
  "read-only": "runtimeBanner.readOnly",
};

export function RuntimeBanner({ runtimeMode }: { runtimeMode: RuntimeMode | undefined }) {
  const t = useTranslate();
  if (runtimeMode !== "maintenance" && runtimeMode !== "degraded" && runtimeMode !== "read-only") {
    return null;
  }
  return (
    <div
      role="status"
      data-runtime-banner={runtimeMode}
      className="border-b border-amber-500/40 bg-amber-500/10 px-4 py-2 text-sm text-amber-950 dark:text-amber-100"
    >
      {t(BANNER_KEYS[runtimeMode])}
    </div>
  );
}
