// W19: lazy-open the session wallet via POST so GET /me stays read-only
// (W15-F11) while "我的钱包" still feels automatic (GOAL-020/022).
//
// Performance (2026-08): the previous implementation POSTed /api/wallet/me on
// EVERY mount and then reloaded the whole page — even when the wallet already
// existed — doubling the page's request wave. Now the mount probe goes through
// the page fetch cache (`crud.fetchList` with the shared DISPLAY_LIST_QUERY),
// so it coalesces with the statCards' GET /me into one network request, and
// the reload wave fires only when the probe proved the wallet was missing and
// the POST actually created it.
import { useCallback, useEffect, useRef, useState } from "react";

import { useAuth } from "@/account/AuthContext";
import { Button } from "@/components/ui/button";
import { useTranslate } from "@/i18n/runtime";
import {
  DISPLAY_LIST_QUERY,
  isWalletNotFoundError,
} from "@/renderer/resource";
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

  /** True when the session wallet does not exist yet (404 / WALLET_NOT_FOUND). */
  const ensure = useCallback(async () => {
    setBusy(true);
    setError(null);
    try {
      let missing = false;
      if (crud !== null) {
        // Join the statCards' /api/wallet/me request (shared page cache);
        // a missing wallet surfaces as WALLET_NOT_FOUND.
        try {
          await crud.fetchList("/api/wallet/me", DISPLAY_LIST_QUERY, undefined, authFetch);
        } catch (err) {
          if (isWalletNotFoundError(err)) {
            missing = true;
          } else {
            throw err;
          }
        }
      } else {
        // Bare hostless render: probe and create directly.
        const probe = await authFetch("/api/wallet/me", {
          headers: { Accept: "application/json" },
        });
        if (probe.ok) {
          missing = false;
        } else if (probe.status === 404) {
          missing = true;
        } else {
          throw new Error(`wallet probe failed: ${probe.status}`);
        }
      }
      if (!missing) {
        // Wallet exists: the statCards/table already loaded it. No POST, no
        // page reload — this was the silent per-visit duplicate before.
        return;
      }
      const response = await authFetch("/api/wallet/me", {
        method: "POST",
        headers: { Accept: "application/json" },
      });
      if (!response.ok) {
        setError(t("schema.myWallet.ensure.failed"));
        return;
      }
      // The wallet was just created: the 404'd statCards/table need data.
      reloadList?.();
    } catch {
      setError(t("schema.myWallet.ensure.failed"));
    } finally {
      setBusy(false);
    }
  }, [authFetch, crud, reloadList, t]);

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