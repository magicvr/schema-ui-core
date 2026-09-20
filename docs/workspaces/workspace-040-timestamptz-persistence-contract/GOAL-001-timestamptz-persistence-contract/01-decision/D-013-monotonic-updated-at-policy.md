---
id: D-013-monotonic-updated-at-policy
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-013 · 微秒单调 updated_at 裁决

用户选择保留 users/roles 现有单调更新语义，升级为微秒：`max(now.UTC().Truncate(time.Microsecond), old.Add(time.Microsecond))`。该规则用于 C2 codec/runtime mapping，保持 cache/ETag/排序不回退，不依赖 wall-clock 必然向前。
