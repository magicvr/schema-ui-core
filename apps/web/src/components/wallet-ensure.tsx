// W19: lazy-open the session wallet via POST so GET /me stays read-only
// (W15-F11) while "我的钱包" still feels automatic (GOAL-020/022).
import { useCallback, useEffect, useRef, useState } from "react";

import { useAuth } from "@/account/AuthContext";
import { Button } from "@/components/ui/button";
import { useTranslate } from "@/i18n/runtime";
import {
  registerCustomComponent,
  type CustomComponentProps,
} from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";

export function WalletEnsure(_props: CustomComponentProps) {
  const t = useTranslate();
  const { authFetch } = useAuth();
  const crud = useSchemaCrud();
  const reloadList = crud?.reloadList;
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const opened = useRef(false);

  const ensure = useCallback(async () => {
    setBusy(true);
    setError(null);
    try {
      const response = await authFetch("/api/wallet/me", {
        method: "POST",
        headers: { Accept: "application/json" },
      });
      if (!response.ok) {
        setError(t("schema.myWallet.ensure.failed"));
        return;
      }
      reloadList?.();
    } catch {
      setError(t("schema.myWallet.ensure.failed"));
    } finally {
      setBusy(false);
    }
  }, [authFetch, reloadList, t]);

  useEffect(() => {
    if (opened.current) {
      return;
    }
    opened.current = true;
    void ensure();
  }, [ensure]);

  if (error === null) {
    return <div data-wallet-ensure hidden={busy ? undefined : true} aria-busy={busy || undefined} />;
  }

  return (
    <div className="space-y-2 rounded-md border border-destructive/30 bg-destructive/5 p-3" data-wallet-ensure>
      <p role="alert" className="text-sm text-destructive">
        {error}
      </p>
      <Button type="button" variant="outline" className="h-8 text-xs" disabled={busy} onClick={() => void ensure()}>
        {busy ? t("schema.myWallet.ensure.retrying") : t("schema.myWallet.action.open")}
      </Button>
    </div>
  );
}

registerCustomComponent("wallet-ensure", WalletEnsure);
