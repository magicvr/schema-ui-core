// @vitest-environment jsdom

import { beforeEach, describe, expect, it } from "vitest";

import {
  clearTokens,
  getAccessToken,
  getRefreshToken,
  hasSession,
  setAccessToken,
  setRefreshToken,
} from "@/account/tokens";

describe("tokens", () => {
  beforeEach(() => {
    clearTokens();
    window.localStorage.clear();
  });

  it("keeps the access token in memory only (never persisted)", () => {
    expect(getAccessToken()).toBeNull();
    setAccessToken("access-1");
    expect(getAccessToken()).toBe("access-1");
    expect(window.localStorage.getItem("schema-ui.accessToken")).toBeNull();
  });

  it("persists the refresh token in localStorage and clears it", () => {
    setRefreshToken("refresh-1");
    expect(getRefreshToken()).toBe("refresh-1");
    expect(window.localStorage.getItem("schema-ui.refreshToken")).toBe("refresh-1");
    setRefreshToken(null);
    expect(getRefreshToken()).toBeNull();
  });

  it("reports a session when either token exists", () => {
    expect(hasSession()).toBe(false);
    setRefreshToken("refresh-1");
    expect(hasSession()).toBe(true);
    setRefreshToken(null);
    setAccessToken("access-1");
    expect(hasSession()).toBe(true);
  });

  it("clearTokens drops both tokens", () => {
    setAccessToken("access-1");
    setRefreshToken("refresh-1");
    clearTokens();
    expect(getAccessToken()).toBeNull();
    expect(getRefreshToken()).toBeNull();
    expect(hasSession()).toBe(false);
  });
});
