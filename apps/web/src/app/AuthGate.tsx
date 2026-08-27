/**
 * Production auth gate (ADR-0035 stage order): renders the anonymous surfaces
 * (login / invite accept / forced password change / terminal failure screens)
 * and mounts the schema-driven shell `<App>` for an authenticated session.
 *
 * W14 F-001: extracted verbatim from main.tsx so the PRODUCTION assembly of
 * <App> is unit-testable. The entry point runs createRoot side effects on
 * import, which made its wiring invisible to tests — the exact gap that let
 * GOAL-013 F-010 (`/api/schema` 挂认证) ship without a matching schemaFetcher
 * Bearer transport, so every page failed D-VAL loading with an anonymous 401
 * ("无法显示此页面"). The production-wiring regression lock lives in
 * auth-gate.wiring.test.tsx.
 */

import { useCallback } from "react";

import { useAuth } from "@/account/AuthContext";
import { App } from "@/app/App";
import { createConfigAwareFetcher } from "@/app/config-events";
import { HostFailureScreen } from "@/app/HostFailureScreen";
import { LoginPage } from "@/app/LoginPage";
import { ForcePasswordChange } from "@/components/force-password-change";
import { InviteAcceptPage } from "@/components/invite-accept";
import { executeBootRecovery, lockedFailure, reauthFailure } from "@/host/boot";
import { useI18n } from "@/i18n/runtime";
import type { AppManifest, NavigationContext } from "@/protocol/app-manifest";

/** Boot placeholder while the session adapter resolves (ADR-0035 D4). */
export function BootScreen() {
  const { t } = useI18n();
  return (
    <div className="flex min-h-screen items-center justify-center bg-background text-sm text-muted-foreground">
      {t("app.boot.checkingSession")}
    </div>
  );
}

/**
 * Host-layer configuration side effects after successful resource writes
 * (A-006 R-002). Modules identify changed namespaces through a response
 * header; the generic Renderer remains unaware of product-specific behavior.
 */
function useResourceFetcher(authFetch: typeof fetch): typeof fetch {
  return useCallback(createConfigAwareFetcher(authFetch), [authFetch]);
}

/** Renders the login page when unauthenticated, the shell when authenticated. */
export function AuthGate({ manifest }: { manifest: AppManifest }) {
  const { status, user, session, login, logout, authFetch } = useAuth();
  const resourceFetcher = useResourceFetcher(authFetch);

  if (status === "loading") {
    return <BootScreen />;
  }
  if (status === "reauth-required") {
    // Post-boot session loss: the adapter reports reauth-required (ADR-0035
    // D4/D7) — a terminal failure surface, not the anonymous login page.
    // `reauth` captures the return intent before leaving for /login.
    return <HostFailureScreen failure={reauthFailure()} onAction={executeBootRecovery} />;
  }
  if (status === "locked") {
    // GOAL-004 S4-6: account-lock terminal (ADR-0035 D7 / ADR-0036 D6) —
    // home/support only, no reauth, no retry loop.
    return <HostFailureScreen failure={lockedFailure()} onAction={executeBootRecovery} />;
  }
  if (status === "unauthenticated") {
    // workspace-019 R3 (GOAL-004 C4): the public invitation acceptance page
    // lives on the pre-auth surface, keyed by /invite/accept?token=…
    try {
      if (window.location.pathname.startsWith("/invite/accept")) {
        return <InviteAcceptPage />;
      }
    } catch {
      // non-browser context — fall through to login
    }
    return <LoginPage onLogin={login} />;
  }
  if (user?.mustChangePassword === true) {
    // W16-F01: force the initial/reset password change before entering the app.
    return <ForcePasswordChange />;
  }

  const context: NavigationContext = {
    user: user === null ? undefined : (user as unknown as Record<string, unknown>),
    features: session?.features,
  };
  return (
    <App
      manifest={manifest}
      navigationContext={context}
      // Page-schema documents (/api/schema/{pageId}) sit behind auth
      // middleware since GOAL-013 F-010: they must ride the same Bearer /
      // refresh-on-401 transport as every other API call, or every page
      // fails D-VAL loading with an anonymous 401. Locked by
      // auth-gate.wiring.test.tsx — do not remove.
      schemaFetcher={authFetch}
      resourceFetcher={resourceFetcher}
      onLogout={logout}
      currentUser={user}
    />
  );
}
