export const CONFIG_CHANGED_EVENT = "schema-ui:config-changed";
export const CONFIG_CHANGED_HEADER = "X-Schema-UI-Config-Changed";
export const SETTINGS_BRANDING_NAMESPACE = "settings.branding";
/** Account self-service profile saves (W13 T-05): the session refreshes /me. */
export const ACCOUNT_PROFILE_NAMESPACE = "account.profile";

export interface ConfigChangedDetail {
  namespace: string;
}

/** Wraps a resource fetcher with the host-level configuration change hook. */
export function createConfigAwareFetcher(authFetch: typeof fetch): typeof fetch {
  return async (input: RequestInfo | URL, init?: RequestInit) => {
    const response = await authFetch(input, init);
    publishConfigChangeFromResponse(response);
    return response;
  };
}

/** Converts a successful module response into a namespaced host event. */
export function publishConfigChangeFromResponse(response: Response): void {
  if (!response.ok) {
    return;
  }
  const namespace = response.headers.get(CONFIG_CHANGED_HEADER)?.trim();
  if (namespace !== undefined && namespace !== "") {
    notifyConfigChanged(namespace);
  }
}

/** Publishes a host configuration change without naming a product endpoint. */
export function notifyConfigChanged(namespace: string): void {
  const normalized = namespace.trim();
  if (typeof window !== "undefined" && normalized !== "") {
    window.dispatchEvent(
      new CustomEvent<ConfigChangedDetail>(CONFIG_CHANGED_EVENT, {
        detail: { namespace: normalized },
      }),
    );
  }
}

/** Subscribes to one configuration namespace and returns its cleanup hook. */
export function subscribeToConfigChanges(
  namespace: string,
  listener: () => void,
): () => void {
  if (typeof window === "undefined") {
    return () => undefined;
  }
  const normalized = namespace.trim();
  const handleChange = (event: Event) => {
    const detail = (event as CustomEvent<ConfigChangedDetail>).detail;
    if (detail?.namespace === normalized) {
      listener();
    }
  };
  window.addEventListener(CONFIG_CHANGED_EVENT, handleChange);
  return () => window.removeEventListener(CONFIG_CHANGED_EVENT, handleChange);
}
