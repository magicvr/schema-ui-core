import { Component, StrictMode, useCallback, useEffect, useState, type ReactNode } from "react";
import { createRoot } from "react-dom/client";

import { AuthProvider, useAuth } from "@/account/AuthContext";
import { createConfigAwareFetcher } from "@/app/config-events";
import { HostFailureScreen } from "@/app/HostFailureScreen";
import { I18nProvider, useI18n } from "@/i18n/runtime";
import type { NavigationContext } from "@/protocol/app-manifest";
import { App } from "@/app/App";
import { LoginPage } from "@/app/LoginPage";
import { ManifestFailure } from "@/app/ManifestFailure";
import {
  loadAppManifestBytes,
  type AppManifest,
} from "@/protocol/app-manifest";
import {
  discoverBootstrapDocument,
  type BootstrapAuth,
  type BootstrapDiscovery,
} from "@/host/bootstrap";
import { bootHost, executeBootRecovery, reauthFailure, type HostBootState } from "@/host/boot";
import { nextFailureId, type HostFailure } from "@/host/failure";
import "./index.css";

// Theme bootstrap is handled by the synchronous inline script in index.html
// (S1 · C3). The call below is kept as a safety net for module-only entry
// points (e.g. Storybook / test harness) that do not serve index.html.
import { initTheme } from "@/theme/theme";
function applyStoredTheme() {
  initTheme();
}

function BootScreen() {
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
function AuthGate({ manifest }: { manifest: AppManifest }) {
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
  if (status === "unauthenticated") {
    return <LoginPage onLogin={login} />;
  }

  const context: NavigationContext = {
    user: user === null ? undefined : (user as unknown as Record<string, unknown>),
    features: session?.features,
  };
  return (
    <App
      manifest={manifest}
      navigationContext={context}
      resourceFetcher={resourceFetcher}
      onLogout={logout}
      currentUser={user}
    />
  );
}

/** Maps the session adapter state to the bootstrap normalized auth input (D4). */
export function bootstrapAuthFor(
  status: string,
  user: { id: string; name?: string } | null,
): BootstrapAuth {
  if (status === "authenticated" && user !== null) {
    return {
      state: "authenticated",
      principal: { id: user.id, name: user.name ?? "", roles: [] },
      provenance: "host-session-adapter",
    };
  }
  if (status === "reauth-required") {
    // The session adapter has a credential that no longer authenticates:
    // reauth-required terminal — never anonymous, never a stale principal.
    return { state: "reauth-required" };
  }
  return { state: "anonymous" };
}

/**
 * Host boot gate (ADR-0035 stage order): availability-gate and
 * auth-resolution terminals render WITHOUT a manifest fetch; manifest
 * failures keep ADR-0025 semantics (ManifestFailure).
 */
function HostBootGate({ discovery }: { discovery: BootstrapDiscovery }) {
  const { status, user } = useAuth();
  const [boot, setBoot] = useState<HostBootState | null>(null);
  const [bootError, setBootError] = useState<unknown>(null);

  useEffect(() => {
    if (status === "loading") return;
    let cancelled = false;
    const auth = bootstrapAuthFor(status, user);
    bootHost({
      documentResult: discovery,
      auth,
      manifestLoader: async () => loadAppManifestBytes(),
    }).then(
      (state) => {
        if (!cancelled) setBoot(state);
      },
      (error: unknown) => {
        if (!cancelled) setBootError(error);
      },
    );
    return () => {
      cancelled = true;
    };
  }, [status, user, discovery]);

  if (status === "loading" || boot === null) {
    if (bootError !== null) {
      return <ManifestFailure error={bootError} />;
    }
    return <BootScreen />;
  }
  if (boot.failure !== null) {
    return <HostFailureScreen failure={boot.failure} onAction={executeBootRecovery} />;
  }
  if (boot.manifest === null) {
    return <ManifestFailure error={new Error("Host boot produced no manifest.")} />;
  }
  return <AuthGate manifest={boot.manifest} />;
}

/** Uncaught renderer/Host exceptions become HOST_RENDER_FAILED (no auto-reload loop). */
class RenderFailureBoundary extends Component<
  { children: ReactNode },
  { failure: HostFailure | null }
> {
  state: { failure: HostFailure | null } = { failure: null };

  static getDerivedStateFromError(): { failure: HostFailure } {
    return {
      failure: {
        failureVersion: "1.0",
        failureId: nextFailureId(),
        scope: "runtime",
        kind: "render-failed",
        hostCode: "HOST_RENDER_FAILED",
        retry: { mode: "manual" },
        message: { messageKey: "hostFailure.renderFailed" },
        recoveryActions: [{ type: "reload" }],
      },
    };
  }

  componentDidCatch(error: unknown): void {
    console.error("[schema-ui] uncaught render failure:", error);
  }

  render(): ReactNode {
    if (this.state.failure !== null) {
      return <HostFailureScreen failure={this.state.failure} onAction={executeBootRecovery} />;
    }
    return this.props.children;
  }
}

const root = document.getElementById("root");
if (!root) {
  throw new Error("root element not found");
}

applyStoredTheme();

discoverBootstrapDocument()
  .then((discovery) => {
    createRoot(root).render(
      <StrictMode>
        <I18nProvider systemDefaultUrl="/api/branding">
          <AuthProvider>
            <RenderFailureBoundary>
              <HostBootGate discovery={discovery} />
            </RenderFailureBoundary>
          </AuthProvider>
        </I18nProvider>
      </StrictMode>,
    );
  })
  .catch((error: unknown) => {
    // Discovery itself threw (unexpected transport class) — surface as an
    // offline-classified bootstrap document failure.
    createRoot(root).render(
      <StrictMode>
        <I18nProvider systemDefaultUrl="/api/branding">
          <AuthProvider>
            <RenderFailureBoundary>
              <HostFailureScreen
                failure={{
                  failureVersion: "1.0",
                  failureId: nextFailureId(),
                  scope: "bootstrap",
                  kind: "offline",
                  hostCode: "HOST_OFFLINE",
                  retry: { mode: "manual" },
                  message: { messageKey: "hostFailure.offline" },
                  recoveryActions: [{ type: "retry" }],
                }}
                onAction={executeBootRecovery}
              />
            </RenderFailureBoundary>
          </AuthProvider>
        </I18nProvider>
      </StrictMode>,
    );
    console.error("[schema-ui] bootstrap discovery failed:", error);
  });
