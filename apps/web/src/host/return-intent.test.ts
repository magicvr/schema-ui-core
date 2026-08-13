// @vitest-environment jsdom
import { beforeEach, describe, expect, it } from "vitest";

import {
  buildQueryString,
  captureReturnIntent,
  returnIntentStorageKey,
  takeReturnIntent,
} from "@/host/return-intent";

function memoryStorage(): Storage {
  const map = new Map<string, string>();
  return {
    get length() {
      return map.size;
    },
    clear: () => map.clear(),
    getItem: (key: string) => map.get(key) ?? null,
    key: (index: number) => [...map.keys()][index] ?? null,
    removeItem: (key: string) => void map.delete(key),
    setItem: (key: string, value: string) => void map.set(key, value),
  } as Storage;
}

const NOW = "2026-08-13T10:00:00.000Z";

describe("host return-intent lifecycle (ADR-0036 D6)", () => {
  let storage: Storage;
  beforeEach(() => {
    storage = memoryStorage();
  });

  it("captures the current in-app path and consumes it after login", () => {
    captureReturnIntent({ path: "/users", query: { tab: "list" }, nowIso: NOW, storage });
    expect(storage.getItem(returnIntentStorageKey())).not.toBeNull();

    const target = takeReturnIntent({ nowIso: "2026-08-13T10:05:00.000Z", storage });
    expect(target).not.toBeNull();
    expect(target?.path).toBe("/users");
    expect(target?.query).toEqual({ tab: "list" });
    // Single-use: consumed on read.
    expect(storage.getItem(returnIntentStorageKey())).toBeNull();
  });

  it("production call sites: the no-argument capture reads the live location, query included", () => {
    // The production call sites (AuthContext auth-lost, boot reauth action) call
    // captureReturnIntent() with no arguments — the live location (including its
    // query string) must be captured, not an empty query (A-008 F-2 regression).
    window.history.replaceState({}, "", "/users?tab=list&sort=name");
    try {
      captureReturnIntent({ nowIso: NOW, storage });
      const raw = storage.getItem(returnIntentStorageKey());
      expect(raw).not.toBeNull();
      const intent = JSON.parse(raw as string) as { path: string; query: Record<string, string> };
      expect(intent.path).toBe("/users");
      expect(intent.query).toEqual({ tab: "list", sort: "name" });
      // Full consumption round-trip keeps only allowlisted keys.
      const target = takeReturnIntent({ nowIso: "2026-08-13T10:01:00.000Z", storage });
      expect(target?.path).toBe("/users");
      expect(target?.query).toEqual({ tab: "list", sort: "name" });
    } finally {
      window.history.replaceState({}, "", "/");
    }
  });

  it("never captures the login surface (self-loop prevention)", () => {
    captureReturnIntent({ path: "/login", query: {}, nowIso: NOW, storage });
    expect(storage.getItem(returnIntentStorageKey())).toBeNull();
    captureReturnIntent({ path: "/login/", query: {}, nowIso: NOW, storage });
    expect(storage.getItem(returnIntentStorageKey())).toBeNull();
  });

  it("rejects expired intents on consume", () => {
    captureReturnIntent({ path: "/users", query: {}, nowIso: NOW, storage });
    const target = takeReturnIntent({ nowIso: "2026-08-13T11:00:00.000Z", storage });
    expect(target).toBeNull();
    // Expired intent is still consumed (no replay, no stale restore).
    expect(storage.getItem(returnIntentStorageKey())).toBeNull();
  });

  it("drops non-allowlisted and sensitive keys on consume", () => {
    captureReturnIntent(
      { path: "/users", query: { tab: "keep", token: "secret", other: "x" }, nowIso: NOW, storage },
    );
    const target = takeReturnIntent({ nowIso: "2026-08-13T10:01:00.000Z", storage });
    expect(target?.query).toEqual({ tab: "keep" });
  });

  it("returns null when nothing is stored or storage content is malformed", () => {
    expect(takeReturnIntent({ nowIso: NOW, storage })).toBeNull();
    storage.setItem(returnIntentStorageKey(), "{not json");
    expect(takeReturnIntent({ nowIso: NOW, storage })).toBeNull();
  });

  it("builds a query string only from kept keys", () => {
    expect(buildQueryString({ tab: "history", sort: "name" })).toBe("?tab=history&sort=name");
    expect(buildQueryString({})).toBe("");
  });
});
