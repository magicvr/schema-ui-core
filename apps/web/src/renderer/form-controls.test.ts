import { describe, expect, it } from "vitest";

import {
  checkFormCapabilities,
  coerceFieldValue,
  isWhitelistedFormControl,
  wireKindOf,
  type FormControlField,
} from "@/renderer/form-controls";

const baseMeta = { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] };

describe("isWhitelistedFormControl", () => {
  it("accepts §5 whitelist types", () => {
    for (const type of [
      "input",
      "select",
      "textarea",
      "switch",
      "checkbox",
      "radio",
      "cascader",
      "checkboxGroup",
      "richText",
      "password",
    ]) {
      expect(isWhitelistedFormControl(type)).toBe(true);
    }
  });

  it("rejects non-whitelist types", () => {
    expect(isWhitelistedFormControl("upload")).toBe(false);
    expect(isWhitelistedFormControl("chart")).toBe(false);
    expect(isWhitelistedFormControl("slider")).toBe(false);
  });
});

describe("wireKindOf", () => {
  it("maps base and extended/advanced types to wire kinds", () => {
    expect(wireKindOf({ id: "a", type: "input" })).toBe("string");
    expect(wireKindOf({ id: "b", type: "select" })).toBe("string");
    expect(wireKindOf({ id: "c", type: "select", mode: "multiple" })).toBe("string-array");
    expect(wireKindOf({ id: "d", type: "switch" })).toBe("boolean");
    expect(wireKindOf({ id: "e", type: "checkbox" })).toBe("boolean");
    expect(wireKindOf({ id: "f", type: "radio" })).toBe("string");
    expect(wireKindOf({ id: "g", type: "cascader" })).toBe("string-array");
    expect(wireKindOf({ id: "h", type: "checkboxGroup" })).toBe("string-array");
    expect(wireKindOf({ id: "i", type: "richText" })).toBe("string");
    expect(wireKindOf({ id: "j", type: "password" })).toBe("string");
  });
});

describe("coerceFieldValue", () => {
  it("coerces booleans and arrays", () => {
    expect(coerceFieldValue({ id: "s", type: "switch" }, true)).toBe(true);
    expect(coerceFieldValue({ id: "s", type: "switch" }, "yes")).toBe(false);
    expect(
      coerceFieldValue({ id: "m", type: "checkboxGroup" }, ["a", "b"]),
    ).toEqual(["a", "b"]);
    expect(
      coerceFieldValue({ id: "m", type: "checkboxGroup" }, ["a", 1]),
    ).toEqual(["a"]);
  });

  it("applies defaultValue when the value is empty", () => {
    expect(coerceFieldValue({ id: "n", type: "input", defaultValue: "x" }, "")).toBe("x");
    expect(
      coerceFieldValue({ id: "o", type: "checkboxGroup", defaultValue: ["a"] }, []),
    ).toEqual(["a"]);
    expect(coerceFieldValue({ id: "p", type: "switch", defaultValue: true }, undefined)).toBe(true);
  });

  it("ignores defaultValue when a value is present", () => {
    expect(coerceFieldValue({ id: "n", type: "input", defaultValue: "x" }, "y")).toBe("y");
  });

  it("coerces string kinds to string", () => {
    expect(coerceFieldValue({ id: "r", type: "richText" }, 42)).toBe("42");
  });
});

describe("checkFormCapabilities", () => {
  const extendedMeta = {
    protocolVersion: "2.6",
    requiredCapabilities: ["app.manifest", "form.controls.extended"],
  };
  const advancedMeta = {
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "form.controls.extended", "form.controls.advanced"],
  };

  it("passes base controls without a capability gate", () => {
    const errors = checkFormCapabilities(
      baseMeta,
      [{ id: "name", type: "input" }, { id: "kind", type: "select" }],
    );
    expect(errors).toEqual([]);
  });

  it("gates extended controls behind protocol 2.6 + capability", () => {
    // Version ok (2.7) but capability missing → capability error only.
    const missingCap = checkFormCapabilities(baseMeta, [{ id: "x", type: "switch" }]);
    expect(missingCap).toHaveLength(1);
    expect(missingCap[0]!.code).toBe("FORM_CAPABILITY_REQUIRED");
    // Capability present but version too low → version error only.
    const lowVersion = checkFormCapabilities(
      { protocolVersion: "2.5", requiredCapabilities: ["app.manifest", "form.controls.extended"] },
      [{ id: "x", type: "switch" }],
    );
    expect(lowVersion).toHaveLength(1);
    expect(lowVersion[0]!.code).toBe("FORM_VERSION_TOO_LOW");
    // Both satisfied → pass.
    expect(checkFormCapabilities(extendedMeta, [{ id: "x", type: "switch" }])).toEqual([]);
  });

  it("gates advanced controls behind protocol 2.7 + capability", () => {
    const missingCap = checkFormCapabilities(baseMeta, [{ id: "y", type: "cascader" }]);
    expect(missingCap).toHaveLength(1);
    expect(missingCap[0]!.code).toBe("FORM_CAPABILITY_REQUIRED");
    const lowVersion = checkFormCapabilities(
      { protocolVersion: "2.6", requiredCapabilities: ["app.manifest", "form.controls.advanced"] },
      [{ id: "y", type: "cascader" }],
    );
    expect(lowVersion).toHaveLength(1);
    expect(lowVersion[0]!.code).toBe("FORM_VERSION_TOO_LOW");
    expect(checkFormCapabilities(advancedMeta, [{ id: "y", type: "cascader" }])).toEqual([]);
  });

  it("gates select multiple behind extended", () => {
    const missingCap = checkFormCapabilities(baseMeta, [
      { id: "z", type: "select", mode: "multiple" },
    ]);
    expect(missingCap).toHaveLength(1);
    expect(missingCap[0]!.code).toBe("FORM_CAPABILITY_REQUIRED");
    expect(
      checkFormCapabilities(extendedMeta, [{ id: "z", type: "select", mode: "multiple" }]),
    ).toEqual([]);
  });

  it("rejects non-whitelisted types", () => {
    const errors = checkFormCapabilities(baseMeta, [
      { id: "u", type: "upload" as FormControlField["type"] },
    ]);
    expect(errors.some((error) => error.code === "FORM_TYPE_NOT_WHITELISTED")).toBe(true);
  });

  it("gates any defaultValue behind protocol 2.7 + form.controls.advanced", () => {
    // 2.6 (below the defaultValue floor) → version error only.
    const lowVersion = checkFormCapabilities(
      { protocolVersion: "2.6", requiredCapabilities: ["app.manifest", "form.controls.advanced"] },
      [{ id: "d", type: "input", defaultValue: "ok" }],
    );
    expect(lowVersion).toHaveLength(1);
    expect(lowVersion[0]!.code).toBe("FORM_VERSION_TOO_LOW");

    // 2.7 but missing form.controls.advanced → capability error only.
    const missingCap = checkFormCapabilities(baseMeta, [
      { id: "d", type: "input", defaultValue: "ok" },
    ]);
    expect(missingCap).toHaveLength(1);
    expect(missingCap[0]!.code).toBe("FORM_CAPABILITY_REQUIRED");

    // 2.7 + full advanced → passes (wire type matches).
    expect(
      checkFormCapabilities(advancedMeta, [{ id: "d", type: "input", defaultValue: "ok" }]),
    ).toEqual([]);
  });

  it("rejects defaultValue type mismatch", () => {
    const errors = checkFormCapabilities(advancedMeta, [
      { id: "d", type: "input", defaultValue: 42 },
    ]);
    expect(errors.some((error) => error.code === "DEFAULT_VALUE_TYPE_MISMATCH")).toBe(true);
  });
});
