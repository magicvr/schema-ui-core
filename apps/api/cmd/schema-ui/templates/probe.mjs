// golden-web 消费探针：仅经 npm 包（@magicvr/schema-ui-protocol）调用协议面功能。
import {
  APP_MANIFEST_PROTOCOL_VERSION,
  APP_MANIFEST_SUPPORTED_PROTOCOL_VERSIONS,
  DEFAULT_MANIFEST_PATH,
  isValidExpression,
  stripPathQuery,
  resolveSchemaUrl,
} from "@magicvr/schema-ui-protocol";
import assert from "node:assert";

assert.equal(APP_MANIFEST_PROTOCOL_VERSION, "2.9");
assert.ok(APP_MANIFEST_SUPPORTED_PROTOCOL_VERSIONS.includes("2.8"));
assert.equal(DEFAULT_MANIFEST_PATH, "/.well-known/schema-ui/app-manifest.json");
assert.equal(isValidExpression("$context.user.role == \"admin\""), true);
assert.equal(isValidExpression("true"), false);
assert.equal(stripPathQuery("/api/users?page=2"), "/api/users");
assert.equal(
  resolveSchemaUrl("http://localhost:25173", "/schemas/pages/{id}.json", { id: "42" }),
  "http://localhost:25173/schemas/pages/42.json",
);

console.log("golden-web protocol probe PASS ·", APP_MANIFEST_PROTOCOL_VERSION);