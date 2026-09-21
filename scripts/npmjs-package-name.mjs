const packageNames = new Set(["protocol", "lib", "theme", "ui", "renderer", "shell"]);

/** Map an internal package name or pnpm tarball name to its public npmjs name. */
export function npmjsPackageName(sourceName, scope = "@magicvr") {
  const leaf = String(sourceName).split("/").at(-1) ?? "";
  const withoutVersion = leaf.replace(/-\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$/, "");
  const withoutLegacyScope = withoutVersion.replace(/^magicvr-/, "");
  const packageName = withoutLegacyScope.startsWith("schema-ui-")
    ? withoutLegacyScope.slice("schema-ui-".length)
    : withoutLegacyScope;

  if (!packageNames.has(packageName)) {
    throw new Error(`unsupported schema-ui package name: ${sourceName}`);
  }
  return `${scope}/schema-ui-${packageName}`;
}
