import { StrictMode, useCallback } from "react";
import { createRoot } from "react-dom/client";

import { AuthProvider, useAuth } from "@/account/AuthContext";
import type { NavigationContext } from "@/protocol/app-manifest";
import { App } from "@/app/App";
import { notifyBrandingChanged } from "@/app/branding";
import { LoginPage } from "@/app/LoginPage";
import { ManifestFailure } from "@/app/ManifestFailure";
import { loadAppManifest, type AppManifest } from "@/protocol/app-manifest";
import "./index.css";

function applyStoredTheme() {
  const stored = localStorage.getItem("theme");
  const preferDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  if (stored === "dark" || (!stored && preferDark)) {
    document.documentElement.classList.add("dark");
  }
}

function BootScreen() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background text-sm text-muted-foreground">
      Checking session…
    </div>
  );
}

/**
 * Host-layer side effects after successful resource writes (A-006 R-002).
 * Keeps product endpoints (settings branding) out of the generic Renderer.
 */
function useResourceFetcher(authFetch: typeof fetch): typeof fetch {
  return useCallback(
    async (input: RequestInfo | URL, init?: RequestInit) => {
      const response = await authFetch(input, init);
      if (response.ok) {
        try {
          const path = new URL(String(input), window.location.origin).pathname;
          const method = (init?.method ?? "GET").toUpperCase();
          if (method === "PATCH" && path.startsWith("/api/settings")) {
            notifyBrandingChanged();
          }
        } catch {
          // Malformed URL: ignore side-effect; response still returned.
        }
      }
      return response;
    },
    [authFetch],
  );
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
        <AuthProvider>
          <AuthGate manifest={manifest} />
        </AuthProvider>
      </StrictMode>,
    );
  })
  .catch((error: unknown) => {
    createRoot(root).render(
      <StrictMode>
        <ManifestFailure error={error} />
      </StrictMode>,
    );
  });
