import { afterEach, describe, expect, it } from "vitest";

import {
  CONFIG_CHANGED_EVENT,
  notifyConfigChanged,
  publishConfigChangeFromResponse,
  SETTINGS_BRANDING_NAMESPACE,
  subscribeToConfigChanges,
} from "@/app/config-events";

describe("configuration change events", () => {
  afterEach(() => {
    Reflect.deleteProperty(globalThis, "window");
  });

  it("routes a changed namespace to matching subscribers", () => {
    const listeners = new Map<string, EventListener[]>();
    const fakeWindow = {
      addEventListener(type: string, listener: EventListener) {
        listeners.set(type, [...(listeners.get(type) ?? []), listener]);
      },
      removeEventListener(type: string, listener: EventListener) {
        listeners.set(type, (listeners.get(type) ?? []).filter((item) => item !== listener));
      },
      dispatchEvent(event: Event) {
        for (const listener of listeners.get(event.type) ?? []) {
          listener(event);
        }
        return true;
      },
    };
    Object.defineProperty(globalThis, "window", { configurable: true, value: fakeWindow });

    let changes = 0;
    const unsubscribe = subscribeToConfigChanges(SETTINGS_BRANDING_NAMESPACE, () => {
      changes += 1;
    });

    notifyConfigChanged("other.namespace");
    expect(changes).toBe(0);
    notifyConfigChanged(SETTINGS_BRANDING_NAMESPACE);
    expect(changes).toBe(1);

    publishConfigChangeFromResponse(
      new Response(null, {
        status: 200,
        headers: { "X-Schema-UI-Config-Changed": SETTINGS_BRANDING_NAMESPACE },
      }),
    );
    expect(changes).toBe(2);
    publishConfigChangeFromResponse(
      new Response(null, {
        status: 500,
        headers: { "X-Schema-UI-Config-Changed": SETTINGS_BRANDING_NAMESPACE },
      }),
    );
    expect(changes).toBe(2);
    expect(listeners.get(CONFIG_CHANGED_EVENT)).toHaveLength(1);

    unsubscribe();
    notifyConfigChanged(SETTINGS_BRANDING_NAMESPACE);
    expect(changes).toBe(2);
  });
});
