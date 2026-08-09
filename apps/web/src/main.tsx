import { StrictMode, useCallback } from "react";
import { createRoot } from "react-dom/client";

import { AuthProvider, useAuth } from "@/account/AuthContext";
import { createConfigAwareFetcher } from "@/app/config-events";
import { I18nProvider, useI18n } from "@/i18n/runtime";
import type { NavigationContext } from "@/protocol/app-manifest";
import { App } from "@/app/App";
import { LoginPage } from "@/app/LoginPage";
import { ManifestFailure } from "@/app/ManifestFailure";
import { loadAppManifest, type AppManifest } from "@/protocol/app-manifest";
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

const root = document.getElementById("root");
if (!root) {
  throw new Error("root element not found");
}

applyStoredTheme();

loadAppManifest()
  .then((manifest) => {
    createRoot(root).render(
      <StrictMode>
        <I18nProvider>
          <AuthProvider>
            <AuthGate manifest={manifest} />
          </AuthProvider>
        </I18nProvider>
      </StrictMode>,
    );
  })
  .catch((error: unknown) => {
    createRoot(root).render(
      <StrictMode>
        <I18nProvider>
          <ManifestFailure error={error} />
        </I18nProvider>
      </StrictMode>,
    );
  });
