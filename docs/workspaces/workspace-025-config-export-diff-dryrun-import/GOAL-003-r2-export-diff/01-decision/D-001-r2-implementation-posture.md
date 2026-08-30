---
id: D-001
title: R2 实现口径：export 管线 / diff 扁平化 / CLI 参数形态（lead · 合同 D-002 派生）
date: 2026-08-30
status: accepted
---

# D-001 · R2 实现口径（2026-08-30）

配置包合同 v0.1.0（GOAL-002 D-002）为分母；本条目记录实现层派生口径（无可核对的调优已在合同中留白，此处落定）：

1. **export 源树 = 内嵌默认 ∪ 显式文件（按键覆盖）**：任何键未在显式文件出现即取内嵌默认（合同 §1「缺失键 = 使用内嵌默认」的落定）；解析用 `KnownFields` 严格模式 + 多文档拒绝（与 `server.LoadConfig` 同纪律）。`server` 包新增只读导出 `DefaultConfigYAML()`（返回拷贝），供 CLI 读取默认树原文——不改变既有装载语义。
2. **env 引用保留**：从源文本读取（不插值），`${VAR}` / `${VAR:-default}` 形态原样进入包（合同 §1/§3）。
3. **敏感键剔除**：登记表（`auth.jwt_secret` / `admin.initial_password`）+ 宽规则不变量（字段名含 secret/password/token 必须命中登记，防新增敏感字段漏登记）；`secrets.exclude` 的 `env` = 源值中首个 `${VAR}` 名（无引用 → 空串）。
4. **产物键命名**：YAML 与 JSON 两侧共享同一 snake_case tag（同构双输出面）。
5. **diff 比较模型 = 扁平叶子**：config 段递归扁平化为「点分路径 → 显示串」（列表 join），比较 add/modify/remove；信息性元数据（package.*）不参与；路径按字母序输出；`--against` 非包源走同一 export 管线（默认 ∪ 显式）。
6. **CLI 参数形态**：`diff` 用类 `cmdCreate` 的手写参数解析（支持 `<pkg> --against <src>` 任意顺序；标准库 `flag` 无 interspersed 语义，`SetInterspersed` 不存在 → 未选）。
7. **dry-run / import**：仅注册占位（R3 实现），返回 `cliError{2}`（错误语义）。
8. **退出码**：export `0/1`；diff `0` 无差 / `1` 有差 / `2` 错误——经 `cliError` 送达 `main`（既有 `err != nil → exit 1` 路径扩展）。

**未选方案**：标准库 flag 全解析（不支持 positional 穿插）；完整包 format 强校验（宽容读取，diff 只取 config 子树）；JSON 输出沿用 Go 字段名（改共用 tag）。