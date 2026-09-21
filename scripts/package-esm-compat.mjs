const packageNames = new Set(["protocol", "lib", "theme", "ui", "renderer", "shell"]);
const explicitExtensions = /\.(?:js|json|css|svg|png|mjs|cjs|woff2?)$/i;

function toPublicPackageSpecifier(specifier) {
  const schemaAsset = specifier.match(/^@schemas\/([a-z-]+)\.schema\.json$/);
  if (schemaAsset) {
    return `@magicvr/schema-ui-protocol/schemas/${schemaAsset[1]}.schema.json`;
  }

  const internalPackage = specifier.match(/^@schema-ui\/([^/]+)(\/.*)?$/);
  if (internalPackage && packageNames.has(internalPackage[1])) {
    return `@magicvr/schema-ui-${internalPackage[1]}${internalPackage[2] ?? ""}`;
  }

  return specifier;
}

function addEsmFileExtension(specifier) {
  const normalized = specifier.replace(/\.tsx?\.js$/i, ".js").replace(/\.tsx?$/i, ".js");
  const relative = normalized.startsWith("./") || normalized.startsWith("../");
  const packageSubpath = /^@magicvr\/schema-ui-(?:protocol|lib|theme|ui|renderer|shell)\/.+/.test(normalized);

  if (!relative && !packageSubpath) return normalized;
  if (explicitExtensions.test(normalized)) return normalized;
  return `${normalized}.js`;
}

function hasJsonImportAttribute(attributes = "") {
  return /\btype\s*:\s*(["'])json\1/.test(attributes);
}

/** Normalize static ESM import/export specifiers in generated package JavaScript. */
export function normalizePackageEsmImports(source) {
  return source.replace(
    /(\bfrom\s*|\bimport\s+)(["'])([^"']+)\2(\s+(?:with|assert)\s+\{[^}]*\})?/g,
    (_match, prefix, quote, rawSpecifier, rawAttributes = "") => {
      const publicSpecifier = toPublicPackageSpecifier(rawSpecifier);
      const specifier = addEsmFileExtension(publicSpecifier);
      const isJson = specifier.endsWith(".json");

      let attributes = rawAttributes;
      if (isJson) {
        if (hasJsonImportAttribute(attributes)) {
          attributes = attributes.replace(/^\s+assert\b/, " with");
        } else {
          attributes = ' with { type: "json" }';
        }
      }

      return `${prefix}${quote}${specifier}${quote}${attributes}`;
    },
  );
}
