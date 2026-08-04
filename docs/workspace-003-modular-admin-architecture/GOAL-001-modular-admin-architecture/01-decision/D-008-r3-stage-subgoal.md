---
id: D-008-r3-stage-subgoal
title: R3 bounded pilot child and I-006 gate
status: accepted
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# D-008 · R3 bounded pilot child and I-006 gate

## Decision

Create the flat child `GOAL-004-r3-bounded-pilot` to carry R3 progressively.
The first checkpoint is information collection for I-006; implementation and
V-1-V-4 verification remain gated until the deletion, compatibility, warning,
and rollback boundaries are recorded and verified.

## Scope

The pilot focuses on operationlog/activity/settings and the four VP-003 draft
lesions. `core.operationlog` remains always-on; Activity and Settings are
optional Profile-controlled surfaces. The child must prove all five A
deliveries, all four B lesions, V-1-V-4, and the D-gate before R4 planning.

## Boundary

Creating the child does not mark R3 started as a passed stage, close I-006,
authorize R4, or close VP-003. The current centralized `handler.Register`
surface and missing cross-layer V-1-V-4 tests are recorded as implementation
gaps for the child.
