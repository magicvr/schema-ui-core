/**
 * Host-supported page protocol versions and capabilities (GOAL-041 S2 · F-001).
 *
 * Single source of truth for the production host's support set:
 *  - `HOST_SUPPORTED_PAGE_VERSIONS` gates page-document negotiation
 *    (load-page.ts UNSUPPORTED_PROTOCOL_VERSION).
 *  - `HOST_SUPPORTED_CAPABILITIES` gates page-document capability negotiation
 *    (load-page.ts MISSING_REQUIRED_CAPABILITY) and backs `boot.ts` HOST_SUPPORT.
 *
 * Keep in sync with `apps/web/scripts/generate-claim.mjs` `support.capabilities`
 * (the build claim must attest the same set), and with the upstream
 * `capability-registry.json` (schema-ui-docs@v2.9.0). All 19 capabilities are
 * implemented by this host; the claim lists their mandatory suites (all green).
 */

/** Exact page protocol versions the renderer accepts (strict negotiation). */
export const HOST_SUPPORTED_PAGE_VERSIONS = ["2.7", "2.8", "2.9"] as const;

/** Capabilities implemented by this host (superset of served-page requirements). */
export const HOST_SUPPORTED_CAPABILITIES = [
  "app.manifest",
  "app.navigation",
  "host.bootstrap",
  "host.failure-recovery",
  "host.conformance-claim",
  "actions.upload",
  "actions.row.request",
  "actions.page.trigger",
  "actions.row.navigate",
  "actions.batch.request",
  "form.record.load",
  "form.controls.extended",
  "form.controls.advanced",
  "form.controls.readonly",
  "table.selection",
  "table.sort",
  "record.view.load",
  "permissions.inheritance",
  "data.route-binding",
] as const;

export type HostSupportedCapability = (typeof HOST_SUPPORTED_CAPABILITIES)[number];
