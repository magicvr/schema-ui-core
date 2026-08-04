---
id: D-005-r2-pinned-artifact-hash-boundary
title: R2 pinned artifact hash boundary
status: accepted
parent: GOAL-003-r2-kernel-composition-root
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
---

# D-005 · R2 pinned artifact hash boundary

## Decision

Pinned JSON integrity checks hash canonical UTF-8 bytes with CRLF converted to
LF before comparison with the recorded provenance value. The expected hashes
and provenance records remain unchanged.

## Reason

The recorded upstream values match the LF bytes of the vendored artifacts.
This checkout has `core.autocrlf=true`, so the same files are presented to the
tests as CRLF. The mismatch was a checkout representation difference, not an
upstream content change.

## Scope

The boundary is applied only to the integrity checks in
`apps/web/src/protocol/upstream-fixtures.test.ts`,
`apps/web/src/protocol/conformance/stage3-fixtures.test.ts`, and
`apps/web/src/renderer/permissions-inheritance.test.ts`. JSON parsing and
application behavior continue to use the same file content. A future change
to the pinned source, provenance, or newline policy requires the full fixture
integrity suite to be rerun.

## Unselected alternative

The expected constants are not replaced with CRLF hashes. That would make the
checks depend on the local checkout policy and break the recorded upstream
provenance contract.
