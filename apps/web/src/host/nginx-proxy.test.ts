import { describe, expect, it } from "vitest";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";

/**
 * W9 F-002 regression: the production nginx config must proxy BOTH protocol
 * discovery documents to the API. The SPA fallback answers 200 text/html for
 * any unproxied path, and the Host classifies a non-JSON bootstrap response as
 * a terminal protocol failure (retry:none) — the documented "not provided"
 * fallback requires a real 404/410, which only an exact-match proxy (or no
 * fallback hit) can deliver behind this server.
 */
describe("nginx protocol discovery proxying (W9 F-002)", () => {
  const root = fileURLToPath(new URL("../../", import.meta.url));
  const conf = readFileSync(`${root}nginx.conf`, "utf8");

  it("proxies the app manifest", () => {
    expect(conf).toContain("location = /.well-known/schema-ui/app-manifest.json");
  });

  it("proxies the host bootstrap document", () => {
    expect(conf).toContain("location = /.well-known/schema-ui/host-bootstrap.json");
  });

  it("sends both discovery documents to the api upstream", () => {
    const blocks = conf.split(/(?=location = \/\.well-known\/)/).filter((chunk) =>
      chunk.startsWith("location = /.well-known/"),
    );
    expect(blocks.length).toBe(2);
    for (const block of blocks) {
      expect(block).toContain("proxy_pass http://api:25080;");
    }
  });
});
