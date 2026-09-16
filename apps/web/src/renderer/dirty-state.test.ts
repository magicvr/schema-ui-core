import { afterEach, describe, expect, it } from "vitest";

import {
  equalDirtyValues,
  hasDirtyState,
  registerDirtyStateSource,
  resetDirtyStateSources,
} from "@/renderer/dirty-state";

afterEach(() => resetDirtyStateSources());

describe("dirty-state registry", () => {
  it("tracks live predicates and unregisters them", () => {
    let dirty = false;
    const dispose = registerDirtyStateSource(() => dirty);
    expect(hasDirtyState()).toBe(false);
    dirty = true;
    expect(hasDirtyState()).toBe(true);
    dispose();
    expect(hasDirtyState()).toBe(false);
  });

  it("fails closed when a source throws", () => {
    registerDirtyStateSource(() => {
      throw new Error("broken source");
    });
    expect(hasDirtyState()).toBe(true);
  });

  it("compares JSON-shaped form values structurally", () => {
    expect(equalDirtyValues({ a: [1, { b: "x" }] }, { a: [1, { b: "x" }] })).toBe(true);
    expect(equalDirtyValues({ a: [1] }, { a: [2] })).toBe(false);
    expect(equalDirtyValues({ a: undefined }, {})).toBe(false);
  });
});
