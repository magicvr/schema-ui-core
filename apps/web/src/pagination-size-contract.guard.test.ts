/**
 * GOAL-011 · pagination page-size contract guard.
 *
 * WHY THIS EXISTS
 * ---------------
 * The web UI's `DEFAULT_PAGE_SIZE` and the API's `handler.DefaultPageSize` are
 * not independent numbers: `buildResourceQuery` omits `pageSize` when the value
 * equals the client constant and lets the server apply its own default. When the
 * two disagreed (10 in the web app, 20 in the API) the "10" option became a
 * no-op and the control displayed a default that was not in effect — a defect
 * the user found in the browser, because the regression test of the day modelled
 * the wrong contract (`?? "10"`).
 *
 * This guard pins the two constants together so the next drift fails a test
 * instead of a user.
 */

import { readFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

import { DEFAULT_PAGE_SIZE, buildResourceQuery } from "@/renderer/resource";

const HERE = dirname(fileURLToPath(import.meta.url));
/** apps/web — this guard lives at src/, so the package root is one level up. */
const WEB_ROOT = resolve(HERE, "..");
/** Repository root (apps/web → apps → repo). */
const REPO_ROOT = resolve(WEB_ROOT, "../..");

const API_RESOURCES = join(REPO_ROOT, "apps/api/internal/handler/resources.go");
const MESSAGES = {
  "zh-CN": join(WEB_ROOT, "src/i18n/messages/zh-CN.json"),
  "en-US": join(WEB_ROOT, "src/i18n/messages/en-US.json"),
} as const;

/** `DefaultPageSize = 20` in the handler package. */
function apiDefaultPageSize(): number {
  const source = readFileSync(API_RESOURCES, "utf8");
  const match = /DefaultPageSize\s*=\s*(\d+)/.exec(source);
  expect(match, `could not read DefaultPageSize from ${API_RESOURCES}`).not.toBeNull();
  return Number(match![1]);
}

describe("GOAL-011 · pagination page-size contract", () => {
  it("keeps the web default page size equal to the API default", () => {
    expect(
      DEFAULT_PAGE_SIZE,
      "apps/web DEFAULT_PAGE_SIZE must match handler.DefaultPageSize, because an omitted pageSize falls back to the server default",
    ).toBe(apiDefaultPageSize());
  });

  it("sends an explicit pageSize for every size the control offers except the default", () => {
    // The control offers 10 / 20 / 50 / 100. Only the default may be omitted —
    // everything else must reach the wire, or the option is a lie.
    for (const size of [10, 50, 100]) {
      expect(buildResourceQuery({ page: 1, pageSize: size })).toContain(`pageSize=${size}`);
    }
    expect(buildResourceQuery({ page: 1, pageSize: DEFAULT_PAGE_SIZE })).not.toContain("pageSize=");
    // 10 is the size the user could not make effective; it must never equal the
    // default again without this test failing first.
    expect(DEFAULT_PAGE_SIZE).not.toBe(10);
  });

  it("labels the page-jump confirm action as a jump in both locales", () => {
    for (const [locale, path] of Object.entries(MESSAGES)) {
      const catalog = JSON.parse(readFileSync(path, "utf8")) as Record<string, string>;
      expect(catalog["feedback.jumpToPage"], `${locale} is missing feedback.jumpToPage`).toBeTypeOf(
        "string",
      );
      expect(catalog["feedback.jumpToPage"]).not.toBe(catalog["feedback.search"]);
      // The form/input keep the "go to page" wording; only the confirm action
      // says "jump", so the two keys must stay distinct.
      expect(catalog["feedback.goToPage"]).not.toBe(catalog["feedback.jumpToPage"]);
    }
    const zh = JSON.parse(readFileSync(MESSAGES["zh-CN"], "utf8")) as Record<string, string>;
    const en = JSON.parse(readFileSync(MESSAGES["en-US"], "utf8")) as Record<string, string>;
    expect(zh["feedback.jumpToPage"]).toBe("跳转");
    expect(en["feedback.jumpToPage"]).toBe("Go");
  });
});
