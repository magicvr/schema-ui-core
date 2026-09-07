/**
 * Host-supported page protocol versions and capabilities (GOAL-041 S2 · F-001).
 *
 * SINGLE SOURCE OF TRUTH: `./host-support.json` (see F2 / GOAL-042 D-001).
 * Both this module (runtime negotiation) and `apps/web/scripts/generate-claim.mjs`
 * (build-time claim) read the same JSON, so the claim can never drift from the
 * runtime support set. Keep the JSON in sync with the upstream
 * `capability-registry.json` (schema-ui-docs@v2.9.0) and the page versions the
 * renderer accepts.
 */

import hostSupportJson from "./host-support.json";

/** Exact page protocol versions the renderer accepts (strict negotiation). */
export const HOST_SUPPORTED_PAGE_VERSIONS: readonly string[] =
  hostSupportJson.supportedPageVersions;

/** Capabilities implemented by this host (superset of served-page requirements). */
export const HOST_SUPPORTED_CAPABILITIES: readonly string[] =
  hostSupportJson.supportedCapabilities;

export type HostSupportedCapability = (typeof HOST_SUPPORTED_CAPABILITIES)[number];
