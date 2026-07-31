import { describe, expect, it } from "vitest";

import { loadAccountContext } from "@/account/context";

describe("loadAccountContext", () => {
  it("maps a session body to a NavigationContext", async () => {
    const { context, error } = await loadAccountContext(
      (async () =>
        new Response(
          JSON.stringify({
            user: { id: "dev-001", name: "Dev Admin", roles: ["admin"] },
            features: { beta: true },
          }),
          { status: 200 },
        )) as typeof fetch,
    );
    expect(error).toBeNull();
    expect(context.user).toEqual({ id: "dev-001", name: "Dev Admin", roles: ["admin"] });
    expect(context.features).toEqual({ beta: true });
  });

  it("fails closed to an empty context on HTTP error", async () => {
    const { context, error } = await loadAccountContext(
      (async () => new Response("missing", { status: 401 })) as typeof fetch,
    );
    expect(context).toEqual({});
    expect(error).not.toBeNull();
  });

  it("fails closed to an empty context on network error", async () => {
    const { context, error } = await loadAccountContext(
      (async () => {
        throw new Error("network down");
      }) as typeof fetch,
    );
    expect(context).toEqual({});
    expect(error).not.toBeNull();
  });
});
