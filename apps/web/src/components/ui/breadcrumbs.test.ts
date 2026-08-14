/**
 * Breadcrumbs component + trail resolution tests (GOAL-015 D-002 §3.6).
 * Route-stack approach: the trail derives from the session visit stack, no
 * manifest/protocol dependency (user-confirmed 2026-08-14).
 */
import { describe, expect, it } from "vitest";

import { resolveBreadcrumbTrail, type BreadcrumbEntry } from "./breadcrumbs";

const t = (key: string) => key;

describe("resolveBreadcrumbTrail (GOAL-015 · route-stack)", () => {
  const pages = [
    { pageId: "dashboard", title: "Dashboard", route: "/" },
    { pageId: "data-dictionary", title: "Data dictionary", route: "/data-dictionary" },
    { pageId: "dictionary-entries", title: "Dictionary entries", route: "/dictionary-entries" },
  ];

  it("empty stack yields a one-entry trail (single-level page)", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[1], t, []);
    expect(trail).toEqual([
      { pageId: "data-dictionary", label: "Data dictionary", route: "/data-dictionary", current: true },
    ]);
  });

  it("stack ancestors produce the trail with current last", () => {
    const trail = resolveBreadcrumbTrail(
      pages,
      pages[2],
      t,
      ["/data-dictionary", "/dictionary-entries"],
    );
    expect(trail).toEqual([
      { pageId: "data-dictionary", label: "Data dictionary", route: "/data-dictionary", current: false },
      {
        pageId: "dictionary-entries",
        label: "Dictionary entries",
        route: "/dictionary-entries",
        current: true,
      },
    ]);
  });

  it("unknown stack paths are skipped (fail-safe)", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[2], t, ["/not-a-page", "/dictionary-entries"]);
    expect(trail.length).toBe(1);
    expect(trail[0].pageId).toBe("dictionary-entries");
    expect(trail[0].current).toBe(true);
  });

  it("duplicate stack paths are deduplicated", () => {
    const trail = resolveBreadcrumbTrail(
      pages,
      pages[2],
      t,
      ["/data-dictionary", "/", "/data-dictionary", "/dictionary-entries"],
    );
    // maxAncestors = 2: the two most recent distinct ancestors are / and
    // /data-dictionary; oldest-first yields [/, /data-dictionary, current].
    expect(trail.map((e) => e.pageId)).toEqual([
      "dashboard",
      "data-dictionary",
      "dictionary-entries",
    ]);
  });

  it("deep stacks are capped at maxAncestors", () => {
    const many = ["/a", "/b", "/c", "/dictionary-entries"];
    const trail = resolveBreadcrumbTrail(pages, pages[2], t, many, 2);
    // None of /a /b /c resolve to known pages, so only the current remains.
    expect(trail.length).toBe(1);
    expect(trail[0].pageId).toBe("dictionary-entries");
  });
});
