// @vitest-environment jsdom
//
// GOAL-020 · 用户反馈（2026-08-16）：钱包行操作「流水」必须把行 id 绑定到
// /wallet-entries/{id}（navigateMapping path 绑定）。

import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { dirname } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { RenderPageDocument } from "@/renderer/render";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";

const __dir = dirname(fileURLToPath(import.meta.url));

function loadDoc(rel: string): unknown {
	return JSON.parse(
		readFileSync(resolve(__dir, "../../../api/internal/modules/" + rel), "utf8"),
	);
}

function json(body: unknown, status = 200): Response {
	return new Response(JSON.stringify(body), {
		status,
		headers: { "Content-Type": "application/json" },
	});
}

function walletFetcher() {
	return (async (input: RequestInfo | URL) => {
		const raw = String(input);
		if (raw.includes("/api/wallet/accounts")) {
			return json({
				items: [
					{
						id: "acct-1", ownerType: "user", ownerId: "u1", currency: "CNY",
						balanceTotal: 1000, balanceAvailable: 700, balanceFrozen: 300,
						status: "active", version: 1, updatedAt: "2026-08-16T00:00:00.000Z",
					},
				],
				total: 1,
				page: 1,
				pageSize: 10,
			});
		}
		return json({ error: "NOT_FOUND", message: "no such route" }, 404);
	}) as typeof fetch;
}

const CONTEXT = {
	user: { id: "admin", name: "Admin", permissions: ["wallet.read", "wallet.write", "wallet.adjust"] },
	features: {},
} as never;

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
	Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
		configurable: true,
		value: true,
	});
});

afterEach(async () => {
	for (const { root, container } of activeRoots.splice(0)) {
		await act(async () => root.unmount());
		container.remove();
	}
});

describe("wallet entries row navigation (GOAL-020)", () => {
	it("binds the row id into /wallet-entries/{id} via navigateMapping", async () => {
		const doc = loadDoc("wallet/schema/wallet.json") as RenderPageDocument;
		const container = document.createElement("div");
		document.body.appendChild(container);
		const root = createRoot(container);
		activeRoots.push({ root, container });
		const navigated: string[] = [];
		await act(async () => {
			root.render(
				<RenderPage
					document={doc}
					context={CONTEXT}
					tableRenderer={(node) => <SchemaTable node={node} fetcher={walletFetcher()} />}
					onNavigate={(route: string) => navigated.push(route)}
				/> as never,
			);
		});
		const entries = Array.from(container.querySelectorAll("button")).find(
			(b) => b.textContent?.trim() === "Entries",
		);
		expect(entries, "Entries row action must render").toBeDefined();
		await act(async () => {
			entries!.click();
		});
		// The row id (acct-1) must be bound — never the literal {id}.
		expect(navigated).toEqual(["/wallet-entries/acct-1"]);
	});
});
