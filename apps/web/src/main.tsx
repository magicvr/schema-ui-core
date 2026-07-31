import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import { loadAccountContext } from "@/account/context";
import { App } from "@/app/App";
import { ManifestFailure } from "@/app/ManifestFailure";
import { loadAppManifest } from "@/protocol/app-manifest";
import "./index.css";

function applyStoredTheme() {
  const stored = localStorage.getItem("theme");
  const preferDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  if (stored === "dark" || (!stored && preferDark)) {
    document.documentElement.classList.add("dark");
  }
}

const root = document.getElementById("root");
if (!root) {
  throw new Error("root element not found");
}

applyStoredTheme();

loadAppManifest()
  .then(async (manifest) => {
    // R4: attach the account $context snapshot before first render so the
    // shell's navigation permission checks evaluate against real identity.
    // A failed load is not fatal: surface it as a banner instead of dropping it.
    const { context, error: accountError } = await loadAccountContext();
    if (accountError !== null) {
      console.error("[account] failed to load session snapshot:", accountError);
    }
    createRoot(root).render(
      <StrictMode>
        <App
          manifest={manifest}
          navigationContext={context}
          accountError={accountError ?? undefined}
        />
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
