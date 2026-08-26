import { Component, StrictMode, useEffect, useState, type ReactNode } from "react";
import { createRoot } from "react-dom/client";

import { AuthProvider, useAuth } from "@/account/AuthContext";
import { HostFailureScreen } from "@/app/HostFailureScreen";
import { I18nProvider } from "@/i18n/runtime";
// W14 F-001: the production gate lives in its own module so its <App> wiring
// (Bearer transports for page schemas + resources) is locked by unit tests.
import { AuthGate, BootScreen } from "@/app/AuthGate";
// GOAL-018: self-registers custom renderer components (mfa-manager in the
// personal-center MFA block; notification-center on the notifications page).
import "@/components/mfa-manager";
// workspace-018 R3: account email identity binding card (GOAL-004 D-001 §4).
import "@/components/email-identity";
import "@/components/account-session-toolbar";
import "@/components/cron-preview";
import "@/components/monitoring-auto-refresh";
import "@/components/import-template-download";
import "@/components/wallet-ensure";
import "@/components/notification-center";
import "@/components/data-permission-scopes";
import "@/components/activity-export";
// VP-017 R7 UX refinement: settings「邮件」tab console (channel-conditional
// fields, mock-record table under mock only, test composer with subject/body).
import "@/components/mail-admin-tab";
import "@/components/password-policy-tab";
import "@/components/invite-issue-card";
import "@/components/invite-resend-dialog";
import { ManifestFailure } from "@/app/ManifestFailure";
import {
  loadAppManifestBytes,
} from "@/protocol/app-manifest";
import {
  discoverBootstrapDocument,
  type BootstrapAuth,
  type BootstrapDiscovery,
} from "@/host/bootstrap";
import { adapterAuthFor, bootHost, executeBootRecovery, type HostBootState, type SessionAdapterState } from "@/host/boot";
import { nextFailureId, type HostFailure } from "@/host/failure";
import "./index.css";

// Theme bootstrap is handled by the synchronous external /theme-init.js script
// in index.html (S1 · C3 + W8 F-002). The call below is kept as a safety net
// for module-only entry points (e.g. Storybook / test harness) that do not serve
// index.html.
import { initTheme } from "@/theme/theme";
function applyStoredTheme() {
  initTheme();
}

/** Maps the session adapter state to the bootstrap normalized auth input (D4). */
export function bootstrapAuthFor(
  status: string,
  user: { id: string; name?: string } | null,
): BootstrapAuth {
  return adapterAuthFor(status as SessionAdapterState, user);
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
    const auth = adapterAuthFor(status, user);
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
