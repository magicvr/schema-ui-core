# schema-ui-core Web

React 19 + Vite + TypeScript + Tailwind application shell for the MVP Admin
workspace.

## Requirements

- Node.js 20+
- npm with the committed `package-lock.json`

## Run

```bash
npm install
npm run dev
# http://localhost:5173

npm test
npm run build
```

## R3 boundary

The shell fetches the pinned app manifest from
`/.well-known/schema-ui/app-manifest.json`, validates the 2.7 contract, projects
manifest navigation, and resolves History API routes. R4 account/permission
behavior and R5 page rendering remain separate goals.
