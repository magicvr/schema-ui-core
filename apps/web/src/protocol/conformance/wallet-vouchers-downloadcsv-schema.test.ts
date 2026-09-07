// Wallet-vouchers page-document structural regression (workspace-029 · E-007):
// the voucher CSV-export declaration must ride the form node's business
// `props` (pinned node schema allows arbitrary business keys there) and must
// NEVER be placed inside action.onSuccess — the pinned OutcomeBehavior schema
// is strictly closed (additionalProperties: false), so doing that fails the
// runtime D-VAL pass and the shell shows "页面 Schema 错误"
// (PAGE_SCHEMA_INVALID). This test validates the REAL module document served
// at /api/schema/wallet-vouchers and locks both directions.
import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

import { validatePageDocument } from "@/protocol/conformance/runtime-schema-validate";

function walletVouchersDocument(): Record<string, unknown> {
  const raw = readFileSync(
    new URL(
      "../../../../../apps/api/modules/wallet/schema/wallet-vouchers.json",
      import.meta.url,
    ),
    "utf8",
  );
  return JSON.parse(raw) as Record<string, unknown>;
}

function asRecord(value: unknown): Record<string, unknown> {
  if (typeof value !== "object" || value === null || Array.isArray(value)) {
    throw new Error("expected a JSON object");
  }
  return value as Record<string, unknown>;
}

describe("wallet-vouchers page document structural validation", () => {
  it("passes D-VAL with the downloadCsv declaration on the form props (E-007)", () => {
    const doc = walletVouchersDocument();
    const validation = validatePageDocument(doc);
    expect(validation.ok, JSON.stringify(validation.errors)).toBe(true);
    expect(asRecord(asRecord(doc).meta).pageId).toBe("wallet-vouchers");

    const actions = asRecord(asRecord(doc).actions);
    const openGenerate = asRecord(actions.openGenerate);
    // Sanity: the real document carries the declaration on the modal form's
    // business props (not inside action.onSuccess).
    const modalForm = asRecord(openGenerate.content);
    const formProps = asRecord(modalForm.props as Record<string, unknown>);
    expect(formProps.downloadCsv).toBeDefined();
  });

  it("fails D-VAL when downloadCsv is moved back into action.onSuccess (pin guard)", () => {
    const doc = walletVouchersDocument();
    const actions = asRecord(asRecord(doc).actions);
    const generateBatch = asRecord(actions.generateBatch);
    const onSuccess = asRecord(generateBatch.onSuccess);
    onSuccess.downloadCsv = { columns: ["code"], fileName: "vouchers_x.csv" };
    const validation = validatePageDocument(doc);
    expect(validation.ok).toBe(false);
  });
});
