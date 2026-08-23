---
id: D-005-credential-injection
title: 配置键名与凭证注入确认（I-003 闭合）
status: accepted
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-001-object-storage
version: 0.1.0
---

# D-005 · 凭证注入确认（闭合 I-003）

## 决定

1. **键名**：沿用 R1 已冻结合同（GOAL-002 D-001 §3）：`storage.objects.s3.{endpoint,region,bucket,access_key_id,secret_access_key,use_path_style}`；env 覆盖 `STORAGE_OBJECTS_S3_*`。本轮不改名。
2. **凭证来源**：仅两条路径——YAML `${VAR}` 插值（configs/.env 或进程 env 提供 VAR），或直接读进程 env。SDK 侧用 **static credentials provider** 从 Config 字段显式构造，禁用默认链（~/.aws、IMDS、EC2 role 均不读取）：确定性、可测、与 fail-closed 校验一一对应。
3. **fail-closed 面**（R1 已实现并测试）：driver=s3 缺 endpoint/bucket/access_key_id/secret_access_key 任一即 LoadError；误配错误只报键名不带值（A-002 F-001 修复）。
4. **边界**：role/工作负载身份认证不在本 VP 分母；如未来需要属合同扩展，须新决策。

## 证据指回

- 键名冻结与校验实现：GOAL-002 D-001 §3 + config_objects_test.go。
- env 模板：apps/api/configs/env.example object storage 段。
- 泄露面审计修复：GOAL-002 03-audit/A-002 F-001（closed-fixed）。

## 未选方案

- **SDK 默认凭证链**：隐式行为不可测，与"配置缺失启动即拒"的合同冲突。
- **Vault/KMS 动态注入**：A5 密钥轮换边界之外（VP-014 非目标）。
