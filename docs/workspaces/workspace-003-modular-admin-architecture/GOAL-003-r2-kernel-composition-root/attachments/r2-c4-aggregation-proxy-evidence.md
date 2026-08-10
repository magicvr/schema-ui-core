# R2 C4 Migration and Manifest aggregation evidence

- Migration collection and Store integration:
  `apps/api/internal/migration/collector.go` and
  `apps/api/internal/store/migrate.go`.
- Profile-selected Manifest projection:
  `apps/api/internal/manifest/manifest.go`.
- Public endpoint and ETag: `apps/api/internal/handler/manifest.go`.
- Dev/prod proxy: `apps/web/vite.config.ts` and `apps/web/nginx.conf`.
- Static production removal: `apps/web/Dockerfile`.
- Focused tests: `apps/api/internal/manifest/manifest_test.go` and
  `apps/api/internal/handler/manifest_test.go`.

This is the R2 skeleton. It does not claim full schema, permission,
reconciliation, or business-module contribution migration.
