/**
 * Site branding / startup configuration (GOAL-013 + VP-007 S3): loads the
 * Settings contribution's public startup projection for shell title + logos,
 * the browser favicon, and the site-wide defaults (locale / timezone / theme).
 * Empty logoUrl means hide logo in UI; document title always uses siteTitle.
 */

import {
  SETTINGS_BRANDING_NAMESPACE,
  subscribeToConfigChanges,
} from "@/app/config-events";

export const DEFAULT_SITE_TITLE = "Schema UI Core";

export interface Branding {
  siteTitle: string;
  logoUrl: string;
  logoUrlLight: string;
  logoUrlDark: string;
  faviconUrl: string;
  defaultLocale: string;
  supportedLocales: string[];
  siteTimezone: string;
  defaultTheme: string;
}

export function subscribeToBrandingChanges(listener: () => void): () => void {
  return subscribeToConfigChanges(SETTINGS_BRANDING_NAMESPACE, listener);
}

export async function fetchBranding(fetcher: typeof fetch = fetch): Promise<Branding> {
  try {
    const response = await fetcher("/api/branding");
    if (!response.ok) {
      return defaultBranding();
    }
    const body = (await response.json()) as Record<string, unknown>;
    const str = (key: string): string =>
      typeof body[key] === "string" ? (body[key] as string).trim() : "";
    return {
      siteTitle: str("siteTitle") !== "" ? str("siteTitle") : DEFAULT_SITE_TITLE,
      logoUrl: str("logoUrl"),
      logoUrlLight: str("logoUrlLight"),
      logoUrlDark: str("logoUrlDark"),
      faviconUrl: str("faviconUrl"),
      defaultLocale: str("defaultLocale") !== "" ? str("defaultLocale") : "auto",
      supportedLocales: Array.isArray(body.supportedLocales)
        ? body.supportedLocales.filter((entry): entry is string => typeof entry === "string")
        : [],
      siteTimezone: str("siteTimezone") !== "" ? str("siteTimezone") : "auto",
      defaultTheme: str("defaultTheme") !== "" ? str("defaultTheme") : "auto",
    };
  } catch {
    return defaultBranding();
  }
}

export function defaultBranding(): Branding {
  return {
    siteTitle: DEFAULT_SITE_TITLE,
    logoUrl: "",
    logoUrlLight: "",
    logoUrlDark: "",
    faviconUrl: "",
    defaultLocale: "auto",
    supportedLocales: [],
    siteTimezone: "auto",
    defaultTheme: "auto",
  };
}

/**
 * Applies site title + favicon to the document; clears the favicon link when
 * no favicon/logo URL is available. VP-007 S3: the favicon comes from
 * `faviconUrl` (falling back to `logoUrl` for backward compatibility).
 */
export function applyDocumentBranding(branding: Branding): void {
  if (typeof document === "undefined") {
    return;
  }
  document.title = branding.siteTitle;
  const favicon = branding.faviconUrl !== "" ? branding.faviconUrl : branding.logoUrl;
  const existing = document.querySelector<HTMLLinkElement>("link[rel='icon'][data-schema-ui-branding='1']");
  if (favicon === "") {
    existing?.remove();
    return;
  }
  const link = existing ?? document.createElement("link");
  link.rel = "icon";
  link.setAttribute("data-schema-ui-branding", "1");
  link.href = favicon;
  if (!existing) {
    document.head.appendChild(link);
  }
}
