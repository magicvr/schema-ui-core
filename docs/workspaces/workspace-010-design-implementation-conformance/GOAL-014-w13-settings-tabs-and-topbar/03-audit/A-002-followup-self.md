---
id: A-002
doc: audit
source: self
status: recorded
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# A-002 · 追加范围（汉堡靠左 + T-05 头像上传）自审（source: self）

- **日期**：2026-08-16
- **scope**：T-02 修正（汉堡靠左）、T-05 头像上传（Go 存储/端点/迁移/持久化 + Web schema/渲染）
- **verdict**：**pass**

## 核对清单

| 项 | 结论 | 证据 |
|----|------|------|
| 汉堡靠左与 D-002 一致 | ✅ | App.tsx 汉堡在主行首元素（<lg）；e2e 390px 视口首元素断言 + 快照（Open navigation menu 在功能区最左） |
| 头像安全模型与品牌资产一致 | ✅ | 共享 RasterAssetStore：服务端重编码（≤256px）、嗅探拒绝（SVG/HTML 415）、公开 GET nosniff/sandbox/immutable；单元测试覆盖 |
| 上传仅认证（自服务） | ✅ | POST /api/account/avatar 无权限键；匿名 401 测试；无 settings.write 依赖 |
| avatarUrl 提交校验 | ✅ | PATCH profile 仅接受空值或头像存储 URL（品牌 URL 等一律 400 测试）；非本存储 URL 不会被删除 |
| 替换/清空清理 | ✅ | 上传替换 + PATCH 清空均 best-effort 删除旧文件（单元测试验证文件消失） |
| 迁移台账 | ✅ | 0035（users.avatar_url）+ 0036（operation_log 事件 CHECK）追加于全局台账末尾；catalog/operations/restart/fresh 测试更新并全绿 |
| /me 快照与头部展示 | ✅ | account.User.AvatarURL + accountFromUser + parseAuthUser 解析；e2e reload 后用户菜单头像 img 可见 |
| i18n/D-VAL/键集 | ✅ | schema.account.field.avatarUrl 双语文案键；schema-keys.structural 与 all-module-schemas-dval 全绿 |
| 回归 | ✅ | Go 0 FAIL；vitest 1029/1029；tsc 0；e2e admin/mvp 8/8 |
| go 判定 | ✅ | 能力追加（两条 admin.account 路由 + users 追加列 + 事件扩展），不改变 Profile 默认集/模块矩阵/Manifest 装配语义 → 无影响、不暂挂（D-002） |

## Findings

无 required/必改 findings。

## 附带说明（不阻断）

- 无头像启动 GC：替换/清空已覆盖常见路径；崩溃遗留最多一个孤儿文件（附属功能可接受残余，D-002 已记录）。
- 头像不做服务端裁剪/圆角（客户端展示层圆角；附属功能最小化）。
