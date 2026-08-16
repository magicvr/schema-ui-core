/**
 * D-VAL 全量结构回归（2026-08-16 用户反馈：数据字典 PAGE_SCHEMA_INVALID）。
 *
 * 运行时 load-page 会对每个页面文档执行 validatePageDocument（fail closed）。
 * 此前只有 wallet 被临时验证过；data-dictionary 的 openEntries 曾残留非法字段
 * 而测试全绿——本测试遍历全部模块 schema 文档做同一结构验证，
 * 让「测试全绿但页面打不开」的类别永远无法再发生。
 */

import { readFileSync, readdirSync } from "node:fs";
import { resolve } from "node:path";

import { describe, expect, it } from "vitest";

import { validatePageDocument } from "@/protocol/conformance/runtime-schema-validate";

const MODULES = resolve(__dirname, "../../../api/internal/modules");

function collectSchemaFiles(): Array<{ module: string; file: string; abs: string }> {
	const out: Array<{ module: string; file: string; abs: string }> = [];
	for (const entry of readdirSync(MODULES, { withFileTypes: true })) {
		if (!entry.isDirectory()) continue;
		const schemaDir = resolve(MODULES, entry.name, "schema");
		let files: string[];
		try {
			files = readdirSync(schemaDir).filter((f) => f.endsWith(".json"));
		} catch {
			continue; // no schema dir
		}
		for (const f of files) {
			out.push({ module: entry.name, file: f, abs: resolve(schemaDir, f) });
		}
	}
	return out;
}

describe("D-VAL · every module page document passes structural validation", () => {
	const docs = collectSchemaFiles();
	expect(docs.length).toBeGreaterThan(10);

	for (const { module, file, abs } of docs) {
		it(`${module}/${file}`, () => {
			const doc = JSON.parse(readFileSync(abs, "utf8"));
			const result = validatePageDocument(doc);
			expect(result.ok, JSON.stringify(result.errors, null, 2)).toBe(true);
		});
	}
});
