# S5 launch evidence (durable in-repo)

Copied from implementer scratch on 2026-08-09 to close A-001 F-001.

| File | Claim |
|------|-------|
| run1-admin.json / run2-admin.json | Two independent API starts, GET /api/branding bodies |
| compare-admin.log | Bodies identical + field sanity |
| go-build.log | go build ./cmd/server |
| web-build.log | npm run build |
| e2e-localization-admin.log | playwright localization.spec.ts (admin) |
| s5-settings-zh.png | e2e screenshot after settings save |
| web-i18n-related.log / api-handler-related.log | matrix suite refresh |
| run1-mvp.json / run2-mvp.json / compare-mvp.log | mvp profile dual-run (F-003) |
