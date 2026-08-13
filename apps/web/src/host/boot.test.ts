import { describe, expect, it } from "vitest";

import { reauthFailure } from "@/host/boot";

describe("host reauth-required terminal (ADR-0035 D4/D7)", () => {
  it("produces a closed reauth-required failure with a reauth recovery action", () => {
    const failure = reauthFailure();
    expect(failure.kind).toBe("reauth-required");
    expect(failure.hostCode).toBe("HOST_REAUTH_REQUIRED");
    expect(failure.message.messageKey).toBe("hostFailure.reauthRequired");
    expect(failure.recoveryActions).toEqual([{ type: "reauth" }]);
    // Reauth-required must never carry a stale principal or a retry that would
    // silently resume into a dead session.
    expect(failure.retry).toBeUndefined();
  });

  it("emits a distinct failure id per occurrence (no dedup across sessions)", () => {
    expect(reauthFailure().failureId).not.toBe(reauthFailure().failureId);
  });
});
