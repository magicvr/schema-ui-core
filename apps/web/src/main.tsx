import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

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
  .then((manifest) => {
    createRoot(root).render(
      <StrictMode>
        <App manifest={manifest} />
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
