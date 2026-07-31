import { describe, expect, it } from "vitest";

import { runRowAction } from "@/renderer/row-action";

type JsonObject = Record<string, unknown>;

function adminPage(): JsonObject {
  return {
    meta: {
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation", "permissions.inheritance"],
    },
    body: {
      type: "table",
      props: {
        columns: [{ field: "name" }],
        actions: [
          { key: "edit", label: "Edit", permissionIntent: "edit" },
          { key: "delete", label: "Delete", permissionIntent: "delete" },
        ],
      },
    },
  };
}

describe("runRowAction (D-ACT non-batch)", () => {
  it("executes a permitted edit action", () => {
    const result = runRowAction({
      page: adminPage(),
      targetId: "edit",
      context: { user: { roles: ["admin"] }, features: {} },
    });
    expect(result.outcome).toBe("EXECUTED");
    expect(result.permissionDenied).toBe(false);
  });

  it("denies when permission cascade denies edit", () => {
    const page = adminPage();
    page.body = {
      type: "section",
      permissionCascade: { keys: ["edit"] },
      permissions: { edit: false },
      children: [adminPage().body],
    };
    const result = runRowAction({
      page,
      targetId: "edit",
      context: { user: { roles: ["admin"] }, features: {} },
    });
    expect(result.outcome).toBe("BLOCKED");
    expect(result.reason).toBe("PERMISSION_DENIED");
    expect(result.permissionDenied).toBe(true);
  });

  it("blocks hidden actions as NOT_VISIBLE", () => {
    const result = runRowAction({
      page: adminPage(),
      targetId: "edit",
      context: { user: { roles: ["admin"] }, features: {} },
      visible: false,
    });
    expect(result.outcome).toBe("BLOCKED");
    expect(result.reason).toBe("NOT_VISIBLE");
  });

  it("requires confirmation and cancels when unconfirmed", () => {
    const result = runRowAction({
      page: adminPage(),
      targetId: "delete",
      context: { user: { roles: ["admin"] }, features: {} },
      confirm: true,
    });
    expect(result.outcome).toBe("CONFIRM_CANCELLED");
    expect(result.confirmed).toBe(false);

    const confirmed = runRowAction({
      page: adminPage(),
      targetId: "delete",
      context: { user: { roles: ["admin"] }, features: {} },
      confirm: true,
      confirmed: true,
    });
    expect(confirmed.outcome).toBe("EXECUTED");
    expect(confirmed.confirmed).toBe(true);
  });

  it("blocks a disabled action", () => {
    const result = runRowAction({
      page: adminPage(),
      targetId: "edit",
      context: { user: { roles: ["admin"] }, features: {} },
      disabled: true,
    });
    expect(result.outcome).toBe("BLOCKED");
    expect(result.reason).toBe("DISABLED");
  });
});
