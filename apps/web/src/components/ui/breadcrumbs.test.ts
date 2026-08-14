/**
 * Breadcrumbs component + trail resolution tests (GOAL-015).
 * SEMANTIC hierarchy: the trail derives from the manifest navigation
 * tree (group labels → page) plus declared parents for inner pages —
 * never from the visit history (user ruling 2026-08-14).
 */
import { describe, expect, it } from "vitest";

import { resolveBreadcrumbTrail } from "./breadcrumbs";

const t = (key: string) => key;

describe("resolveBreadcrumbTrail (GOAL-015 · semantic hierarchy)", () => {
  const pages = [
    { pageId: "dashboard", title: "Dashboard", route: "/" },
    { pageId: "data-dictionary", title: "Data dictionary", route: "/data-dictionary" },
    { pageId: "dictionary-entries", title: "Dictionary entries", route: "/dictionary-entries/{dictKey}" },
    { pageId: "users", title: "Users", route: "/users" },
    { pageId: "roles", title: "Roles", route: "/roles" },
    { pageId: "scheduled-tasks", title: "Scheduled tasks", route: "/scheduled-tasks" },
    { pageId: "task-runs", title: "Task runs", route: "/task-runs" },
  ];
  const navigation = {
    top: [{ pageRef: "dashboard", label: "Dashboard" }],
    sidebar: [
      { pageRef: "data-dictionary", label: "Data dictionary" },
      {
        label: "Admin",
        items: [
          { pageRef: "users", label: "Users" },
          { pageRef: "roles", label: "Roles" },
        ],
      },
      { pageRef: "scheduled-tasks", label: "Scheduled tasks" },
    ],
  };
  const parents = {
    "dictionary-entries": "data-dictionary",
    "task-runs": "scheduled-tasks",
  };

  it("a nav-root page trails as 首页 => 一级页", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[1], t, { navigation, parents, homePageId: "dashboard" });
    expect(trail).toEqual([
      { pageId: "dashboard", label: "Dashboard", route: "/", current: false },
      { pageId: "data-dictionary", label: "Data dictionary", route: "/data-dictionary", current: true },
    ]);
  });

  it("a grouped page trails as 首页 => 组 => 页面", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[4], t, { navigation, parents, homePageId: "dashboard" });
    expect(trail).toEqual([
      { pageId: "dashboard", label: "Dashboard", route: "/", current: false },
      { pageId: "Admin", label: "Admin", route: "", current: false },
      { pageId: "roles", label: "Roles", route: "/roles", current: true },
    ]);
  });

  it("an inner page trails as 首页 => 一级页 => 内页 (declared parent, not history)", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[2], t, { navigation, parents, homePageId: "dashboard" });
    expect(trail).toEqual([
      { pageId: "dashboard", label: "Dashboard", route: "/", current: false },
      { pageId: "data-dictionary", label: "Data dictionary", route: "/data-dictionary", current: false },
      {
        pageId: "dictionary-entries",
        label: "Dictionary entries",
        route: "/dictionary-entries/{dictKey}",
        current: true,
      },
    ]);
  });

  it("a declared parent that is a nav-root page shows home + parent + current", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[6], t, { navigation, parents, homePageId: "dashboard" });
    expect(trail.map((e) => e.pageId)).toEqual(["dashboard", "scheduled-tasks", "task-runs"]);
  });

  it("an unknown declared parent fails safe (home still leads)", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[2], t, {
      navigation,
      parents: { "dictionary-entries": "no-such-page" },
      homePageId: "dashboard",
    });
    expect(trail.map((e) => e.pageId)).toEqual(["dashboard", "dictionary-entries"]);
    expect(trail[1].current).toBe(true);
  });

  it("no navigation, parents or home yields a single-level trail", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[2], t);
    expect(trail.length).toBe(1);
    expect(trail[0].current).toBe(true);
  });

  it("the current page being home yields a single-level trail (no self root)", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[0], t, { navigation, parents, homePageId: "dashboard" });
    expect(trail.length).toBe(1);
    expect(trail[0].pageId).toBe("dashboard");
    expect(trail[0].current).toBe(true);
  });

  it("the same page always shows the same trail (independent of history)", () => {
    const direct = resolveBreadcrumbTrail(pages, pages[2], t, { navigation, parents, homePageId: "dashboard" });
    const afterVisits = resolveBreadcrumbTrail(pages, pages[2], t, { navigation, parents, homePageId: "dashboard" });
    expect(direct).toEqual(afterVisits);
  });

  it("parent chains loop safely (cycle guard)", () => {
    const trail = resolveBreadcrumbTrail(pages, pages[2], t, {
      navigation,
      parents: { "dictionary-entries": "data-dictionary", "data-dictionary": "dictionary-entries" },
      homePageId: "dashboard",
    });
    expect(trail.map((e) => e.pageId)).toEqual(["dashboard", "data-dictionary", "dictionary-entries"]);
  });
});
