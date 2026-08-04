---
id: D-002
title: R2 profile and kernel contract implementation boundary
status: accepted
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-003-r2-kernel-composition-root
version: 0.1.0
---

# D-002 R2 profile and kernel contract

The implementation uses one framework-agnostic `kernel.Module` contract and
one deterministic registry. `mvp`, `admin`, and `custom` resolve to explicit
module IDs; `APP_MODULES_ENABLED` replaces the selected Profile defaults, and
custom without explicit modules fails closed. The kernel API range is parsed
and checked against `2.0.0`; declarations are not accepted as opaque strings.

Fx remains an implementation detail of `internal/composition`. Secret values
use distinct private Fx types so the composition root cannot confuse JWT
material with the bootstrap password hash.

This decision does not claim that business modules have been migrated. R3
retains the implementation boundary for the first real module slice.
