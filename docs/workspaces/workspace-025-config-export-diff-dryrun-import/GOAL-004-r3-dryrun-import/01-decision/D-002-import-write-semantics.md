---
id: D-002
title: import 写入语义冻结（I-025-004 → verified · 用户裁决方案 A）
date: 2026-08-30
status: accepted
---

# D-002 · import 写入语义（2026-08-30 冻结 · 方案 A）

**I-025-004 用户裁决**（2026-08-30 · GUI 选项确认）：**方案 A —— 原子替换 + 应用前备份**。`I-025-004 → verified`。

## 冻结语义（配置包合同 v0.1.0 §2.4 落定）

1. **预检前置**：`import` 首先执行与 `dry-run` 相同的预检（结构校验 + `secrets.exclude` env fail-closed）；任一 check fail → **拒绝导入**，目标文件零触碰（exit 1）。
2. **生成**：目标配置文本 = 包 `config` 树的 YAML 序列化 + 生成头注释（origin `schema-ui-config-package@v1` · imported_at）；敏感键不入文件（env 注入路径，与 serve 装载语义一致）。
3. **备份**：目标文件已存在 → 应用前复制为 `<file>.pre-import.bak`（**保留**供人工核对；备份失败 = fail-closed 中止，不进入应用）。
4. **应用（原子）**：写 `<file>.tmp` → 以 `server.LoadConfig(tmp)` 做**装载校验**（fail-closed 同纪律：APP_ENV 显式、方言合法、环境密钥要求）→ 校验通过后 `os.Rename(tmp, file)` 原子替换。
5. **失败路径**：装载校验失败 / rename 失败 → 删除 tmp；**原文件全程未被触碰**（rename 原子性），备份保留。既有配置不被破坏（判据 #4 达标）。
6. **成功**：保留 `.bak`；报告 checks + applied + backup；exit 0。
7. **与既有升级前快照的关系**：配置面轻量操作，**不引入** DB 级 `VACUUM INTO` / `pg_dump`（VP-013 方言级快照正交，不在本波）。

## 未选方案

- B（原子替换无备份）：人工核对/回滚手段缺失——用户选 A。
- C（简单覆盖写）：写一半失败破坏既有配置，不满足判据 #4——用户选 A。

## 裁决留痕

用户 2026-08-30 在编排器裁决界面确认「A. 原子替换 + 应用前备份（Recommended）」。若后续改选，以书面为准（三路径，需重开本决策）。