---
id: D-005-r1-wire-input-compat
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-005 · 公共 wire 输入兼容裁决

用户在统一公共输出为 6 位微秒 RFC3339 UTC `Z` 的基础上选择：

- **输出严格规范化**：所有公共时间输出为 `YYYY-MM-DDTHH:MM:SS.ffffffZ`。
- **输入兼容旧格式**：合法 RFC3339 输入允许 0/3/6/9 位小数与 `+00:00` 等等价 offset；解析后统一转 UTC `time.Time`，写入按目标 storage codec 规范化。
- 非法时间、非零 offset 解析失败；不接受模糊本地时间字符串。

该裁决扩展 C2 的 formatter/parser/fixture 清单，不改变 DB 分母、不把 JSON/TEXT payload 内嵌时间加入列分母。
