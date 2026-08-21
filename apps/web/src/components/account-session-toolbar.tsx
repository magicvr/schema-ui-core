// Personal-center "revoke other sessions" toolbar (W16-F07): a self-service
// security action that invalidates every other device's refresh/access tokens
// and reissues a fresh pair for the current browser.
import { useState } from "react";

import { useAuth } from "@/account/AuthContext";
import { setAccessToken, setRefreshToken } from "@/account/tokens";
import { Button } from "@/components/ui/button";
import { useTranslate } from "@/i18n/runtime";
import {
  registerCustomComponent,
  type CustomComponentProps,
} from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";

interface RevokeOthersResponse {
  accessToken?: string;
  refreshToken?: string;
}

export function AccountSessionToolbar(_props: CustomComponentProps) {
  const t = useTranslate();
  const { authFetch, refreshSession } = useAuth();
  const crud = useSchemaCrud();
  const [busy, setBusy] = useState(false);
  const [notice, setNotice] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const revokeOthers = async () => {
    if (busy) {
      return;
    }
    if (!window.confirm(t("schema.account.session.revokeOthersConfirm"))) {
      return;
    }
    setBusy(true);
    setNotice(null);
    setError(null);
    try {
      const response = await authFetch("/api/account/sessions/revoke-others", {
        method: "POST",
        headers: { Accept: "application/json" },
      });
      if (!response.ok) {
        setError(t("schema.account.session.revokeOthersFailed"));
        return;
      }
      const body = (await response.json()) as RevokeOthersResponse;
      if (typeof body.accessToken === "string" && typeof body.refreshToken === "string") {
        setAccessToken(body.accessToken);
        setRefreshToken(body.refreshToken);
      }
      await refreshSession();
      crud?.reloadList();
      setNotice(t("schema.account.session.revokeOthersDone"));
    } catch {
      setError(t("schema.account.session.revokeOthersFailed"));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="flex items-center justify-between gap-3">
      <div>
        <h3 className="text-base font-semibold">{t("schema.account.session.title")}</h3>
        <p className="text-sm text-muted-foreground">{t("schema.account.session.revokeOthersHint")}</p>
      </div>
      <Button
        type="button"
        variant="outline"
        disabled={busy}
        onClick={() => void revokeOthers()}
        data-revoke-others
      >
        {t("schema.account.session.revokeOthers")}
      </Button>
      {notice !== null ? <p className="text-xs text-success">{notice}</p> : null}
      {error !== null ? <p className="text-xs text-destructive">{error}</p> : null}
    </div>
  );
}

registerCustomComponent("account-session-toolbar", AccountSessionToolbar);
