/**
 * GOAL-007 follow-up regression: row actions on templated URLs require an
 * explicit requestMapping.path binding (MISSING_PATH_BINDING guard). The
 * module schemas must carry `requestMapping: { path: { id: "$row.id" } }` on
 * every row action whose action URL contains {id} (delete / run now).
 */
import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

import { constructRequest } from "./request-construction";

const __dirname = dirname(fileURLToPath(import.meta.url));
const MODULES = join(__dirname, "../../../../../apps/api/internal/modules");

interface SchemaDoc {
  actions?: Record<string, Record<string, unknown>>;
  body?: { children?: Array<Record<string, unknown>> };
}

function loadSchema(rel: string): SchemaDoc {
  return JSON.parse(readFileSync(join(MODULES, rel), "utf8")) as SchemaDoc;
}

// Every table row action whose referenced action URL contains a {param} slot
// must declare requestMapping.path for that param (users delete precedent).
describe("row actions carry path bindings for templated URLs", () => {
  const suites: Array<{ schema: string; actions: Array<[string, string]> }> = [
    {
      schema: "filelibrary/schema/file-library.json",
      actions: [["delete", "deleteFile"]],
    },
    {
      schema: "datadictionary/schema/data-dictionary.json",
      actions: [["delete", "deleteType"]],
    },
    {
      schema: "datadictionary/schema/dictionary-entries.json",
      actions: [["delete", "deleteEntry"]],
    },
    {
      schema: "scheduledtasks/schema/scheduled-tasks.json",
      actions: [["run", "runTask"], ["delete", "deleteTask"]],
    },
    // W26 (GOAL-038 C3): the users-invites page (workspace-019 R3 follow-up)
    // shipped without its revoke binding — MISSING_PATH_BINDING on click.
    {
      schema: "users/schema/users-invites.json",
      actions: [["revoke", "revokeInvite"]],
    },
  ];

  for (const suite of suites) {
    it(suite.schema + " row actions resolve {id}", () => {
      const doc = loadSchema(suite.schema);
      const table = doc.body?.children?.find((c) => c.type === "table");
      const props = (table?.props ?? {}) as Record<string, unknown>;
      const rowActions = (props.actions ?? []) as Array<Record<string, unknown>>;
      for (const [key, actionRef] of suite.actions) {
        const item = rowActions.find((a) => a.key === key);
        expect(item, key + " row action exists").toBeDefined();
        const mapping = (item?.requestMapping ?? {}) as Record<string, unknown>;
        const path = (mapping.path ?? {}) as Record<string, unknown>;
        expect(path.id, key + " requestMapping.path.id is $row.id").toBe("$row.id");
        // Construction must succeed with a row carrying id.
        const action = (doc.actions ?? {})[actionRef] as Record<string, unknown>;
        const result = constructRequest({
          kind: "rowAction",
          action,
          requestMapping: item?.requestMapping,
          row: { id: "row-1" },
        });
        expect(result.ok, key + " construction: " + JSON.stringify(result)).toBe(true);
      }
    });
  }
});
