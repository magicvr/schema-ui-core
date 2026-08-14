/**
 * Breadcrumbs component + trail resolution tests (GOAL-015 D-002 §3.6).
 */
import { describe, expect, it } from "vitest";

import { resolveBreadcrumbTrail, type BreadcrumbEntry } from "./breadcrumbs";

const t = (key: string) => key;

describe("resolveBreadcrumbTrail (GOAL-015)", () => {
  const pages = [
    { pageId: "dashboard", title: "Dashboard", route: "/" },
    { pageId: "data-dictionary", title: "Data dictionary", route: "/data-dictionary" },
    {
      pageId: "dictionary-entries",
      title: "Dictionary entries",
      route: "/dictionary-entries",
      breadcrumbParent: "data-dictionary",
    },
  ];

  it("single-level page yields a one-entry trail (no breadcrumb UI)", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[0], t);
    expect(trail).toEqual([
      { pageId: "dashboard", label: "Dashboard", route: "/", current: true },
    ]);
  });

  it("nested page yields the ancestor chain with current marked last", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[2], t);
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

  it("broken breadcrumbParent link is fail-safe (no infinite loop)", () => {
    const broken = {
      pageId: "dictionary-entries",
      title: "Entries",
      route: "/dictionary-entries",
      breadcrumbParent: "does-not-exist",
    };
    const trail = resolveBreadcrumbTrail(pages, broken, t);
    expect(trail.length).toBe(1);
    expect(trail[0].pageId).toBe("dictionary-entries");
  });

  it("cyclic breadcrumbParent chain terminates", () => {
    const cyclicPages = [
      { pageId: "a", title: "A", route: "/a", breadcrumbParent: "b" },
      { pageId: "b", title: "B", route: "/b", breadcrumbParent: "a" },
    ];
    const trail = resolveBreadcrumbTrail(cyclicPages, cyclicPages[0], t);
    expect(trail.length).toBe(2);
  });
});
