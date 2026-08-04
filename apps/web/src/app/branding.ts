/**
 * Site branding (GOAL-013): loads the Settings contribution's public branding
 * projection for shell title + logo. Empty logoUrl means hide logo in UI;
 * document title always uses siteTitle.
 */

import {
  SETTINGS_BRANDING_NAMESPACE,
  subscribeToConfigChanges,
} from "@/app/config-events";

export const DEFAULT_SITE_TITLE = "Schema UI Core";

export interface Branding {
  siteTitle: string;
  logoUrl: string;
}

export function subscribeToBrandingChanges(listener: () => void): () => void {
  return subscribeToConfigChanges(SETTINGS_BRANDING_NAMESPACE, listener);
}

export async function fetchBranding(fetcher: typeof fetch = fetch): Promise<Branding> {
  try {
    const response = await fetcher("/api/branding");
    if (!response.ok) {
      return { siteTitle: DEFAULT_SITE_TITLE, logoUrl: "" };
    }
    const body = (await response.json()) as { siteTitle?: unknown; logoUrl?: unknown };
    const siteTitle =
      typeof body.siteTitle === "string" && body.siteTitle.trim() !== ""
        ? body.siteTitle.trim()
        : DEFAULT_SITE_TITLE;
    const logoUrl = typeof body.logoUrl === "string" ? body.logoUrl.trim() : "";
    return { siteTitle, logoUrl };
  } catch {
    return { siteTitle: DEFAULT_SITE_TITLE, logoUrl: "" };
  }
}

/** Applies site title to document.title; clears favicon link when logo is empty. */
export function applyDocumentBranding(branding: Branding): void {
  if (typeof document === "undefined") {
    return;
  }
  document.title = branding.siteTitle;
  const existing = document.querySelector<HTMLLinkElement>("link[rel='icon'][data-schema-ui-branding='1']");
  if (branding.logoUrl === "") {
    existing?.remove();
    return;
  }
  const link = existing ?? document.createElement("link");
  link.rel = "icon";
  link.setAttribute("data-schema-ui-branding", "1");
  link.href = branding.logoUrl;
  if (!existing) {
    document.head.appendChild(link);
  }
}
