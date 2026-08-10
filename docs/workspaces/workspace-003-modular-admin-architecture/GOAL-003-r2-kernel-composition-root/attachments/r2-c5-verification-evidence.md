# R2 C5 verification evidence

- API: `go test ./...` in `apps/api` passed on 2026-08-04 after the R2
  verification repair.
- Web full test: `npm test -- --run` in `apps/web` passed on 2026-08-04 with
  23 test files and 492 tests.
- Web build: `npm run build` in `apps/web` passed on 2026-08-04.
- The earlier three Web pinned-SHA failures were caused by CRLF checkout bytes
  under `core.autocrlf=true`; LF-normalized bytes match the unchanged
  provenance values. The checks now canonicalize CRLF to LF before hashing.
- The exact source map, revision boundary, and command record are in
  `attachments/audit-A-002-r2-evidence-snapshot.md`.
- Verification ran in dirty snapshot
  `HEAD=b1b7650b3202de7a7a7ce6c0bdffe212093fe75f`; it is not CI, clean-revision,
  deployment, or release evidence.
- Grok A-003 re-audit and high-effort A-004 Root-response audit are recorded;
  child A-005 records the F-003 `fixed` response and C5 close-out. Root R2
  stage response remains separate.
