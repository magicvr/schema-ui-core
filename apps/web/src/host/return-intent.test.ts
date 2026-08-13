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
