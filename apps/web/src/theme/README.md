# Theme & fork brand customization

`src/index.css` is the single source of truth for the Design Token system
(Color / Typography / Elevation·Shadow / Radius). Every primitive
(`components/ui/*`), the Schema Renderer, and the Admin Shell consume these
tokens through Tailwind v4's `@theme inline` mapping — nothing reads raw hex
or HSL values.

## Forking with a custom brand

A fork does **not** need a second Token system, a new theme layer, or a
design-tokens package. It only needs to override the same CSS custom
properties `index.css` already declares, *after* `index.css` loads:

```ts
// src/main.tsx
import "./index.css";
import "./theme/brand.css"; // your fork's overrides (copy brand.example.css)
```

See [`brand.example.css`](./brand.example.css) for a minimal, working example
that swaps the brand accent color, the 5-stop chart palette, and the corner
radius for both light and dark mode. `brand-example.test.ts` asserts the
example only overrides token names that `index.css` already declares, so
copying it into a fork can never introduce an unknown token the rest of the
app doesn't already consume.

## Dark/light boot (no FOUC)

`index.html` inlines a synchronous script that reads `localStorage.theme` /
`prefers-color-scheme` and applies the `dark` class before first paint. The
decision logic itself is the pure, directly-testable `resolveTheme` function
in [`theme.ts`](./theme.ts) — see `theme.test.ts` for its unit tests.
