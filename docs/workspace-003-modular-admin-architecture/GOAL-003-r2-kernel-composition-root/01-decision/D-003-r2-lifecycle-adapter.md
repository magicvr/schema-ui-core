---
id: D-003
title: R2 lifecycle adapter at the Fx composition root
status: accepted
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-003-r2-kernel-composition-root
version: 0.1.0
---

# D-003 R2 lifecycle adapter

The composition root adapts the resolved plan to `kernel.Runtime`. It binds
the HTTP listener before module start, runs `Start` then `Ready` in dependency
order, closes resources on any failed start/readiness path, and runs module
`Stop` in reverse order before closing the Store. The kernel package remains
free of Fx and HTTP implementation types.

The adapter is a lifecycle gate and observability boundary for R2. It does not
claim that all future business modules already expose their final runtime
resources; R3 must replace no-op module hooks with the first migrated module
implementations.
