import { describe, expect, it } from "vitest";

import { resolveListObjectLabel } from "@/renderer/list-surface";

describe("resolveListObjectLabel", () => {
  const translate = (key: string): string =>
    ({
      "feedback.listObject.users": "Users",
      "feedback.listObject.roles": "Roles",
      "feedback.listObject.sessions": "Sessions",
      "feedback.listObject.records": "Records",
    })[key] ?? key;

  it("uses the page semantic noun for users and roles", () => {
    expect(resolveListObjectLabel({ pageId: "users", translate })).toBe("Users");
    expect(resolveListObjectLabel({ pageId: "roles", translate })).toBe("Roles");
  });

  it("uses the table semantic key for account sessions", () => {
    expect(
      resolveListObjectLabel({
        pageId: "account",
        tableTitle: "Signed-in sessions",
        tableTitleKey: "schema.account.session.title",
        translate,
      }),
    ).toBe("Sessions");
  });

  it("falls back to a reliable title before the generic records noun", () => {
    expect(resolveListObjectLabel({ tableTitle: "Audit events", translate })).toBe("Audit events");
    expect(resolveListObjectLabel({ translate })).toBe("Records");
  });
});

