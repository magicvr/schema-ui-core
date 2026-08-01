import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import { AuthProvider, useAuth } from "@/account/AuthContext";
import type { NavigationContext } from "@/protocol/app-manifest";
import { App } from "@/app/App";
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

/** Renders the login page when unauthenticated, the shell when authenticated. */
function AuthGate({ manifest }: { manifest: AppManifest }) {
  const { status, user, session, login, logout, authFetch } = useAuth();

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
      recordsFetcher={authFetch}
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
