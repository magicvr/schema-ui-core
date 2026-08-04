---
id: E-003-r2-verification-repair
title: R2 verification repair and rerun
status: recorded
parent: GOAL-003-r2-kernel-composition-root
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
---

# E-003 · R2 verification repair and rerun

## Facts

- The three Web pinned-SHA test failures were reproduced as CRLF checkout
  bytes while the recorded LF-normalized bytes matched the existing
  provenance values.
- The integrity checks now canonicalize CRLF to LF before hashing. Constants,
  provenance, fixture values, and production manifest content were not
  replaced with local CRLF hashes.
- `go test ./...` in `apps/api` passed.
- `npm test -- --run` in `apps/web` passed: 23 test files and 492 tests.
- `npm run build` in `apps/web` passed: TypeScript compilation and Vite
  production build completed.

## Evidence boundary

The verification commands ran against dirty working tree snapshot
`HEAD=b1b7650b3202de7a7a7ce6c0bdffe212093fe75f`. The result is local snapshot
evidence only; it is not CI, clean-revision, deployment, or release evidence.
The source map and command record are in
`attachments/audit-A-002-r2-evidence-snapshot.md`.

## Gate state

The implementation and local self-verification portion of C5 is complete.
Grok A-003 and A-004 independent audits and the Root D-006/E-006/A-005
response are recorded. The child F-003 closure trail and C5 close-out were
subsequently recorded in E-004/A-005; Root R2 stage response remains separate.
