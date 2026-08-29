# Semver / Breaking 流程 v0.1.0（草案）

适用范围：契约冻结面 A 层（kernel）与 B 层（模块装配面）（freeze-face §0/§3）。内部实现面（C 层）不适用。

## 1. 版本号语义

| 段位 | 触发 |
|------|------|
| **major** | 任何 breaking：删除/改名导出符号；变更签名（含构造签名）；`KernelAPIRange` 语义收窄或主号变化；行为语义破坏（校验、权限、格式语义反转）；协议兼容边界变化（upstream pin 升 major 且影响公开行为） |
| **minor** | 只增不改：新增导出符号/能力；现有语义只扩展；模块新增贡献键 |
| **patch** | 缺陷修复、文档、不变更行为契约的实现细节 |

## 2. Breaking 判定（拿不准时）

1. 缺省按 breaking 处理（保守）；
2. 判定需要证据时 → 打开 `semver-breaking` 决策：列出受影响消费面（下游组合根 / 既有模块 DependsOn）与替代路径；
3. 仍不可判定 → **cross 审计**（self + grok build independent，项目默认执行路径）后由编排器提议、用户裁决（P-004）。
4. 分类结论写回 D 条目与 changelog。

## 3. Breaking 发布流程（major）

1. 冻结面变更清单（符号 + 调用面 + 迁移映射）；
2. 版本号 major +1；`KernelAPIVersion`（内核）同步主号；受影响模块 `KernelAPIRange` 更新；
3. changelog 追加 major 节（breaking 列表 + **迁移说明**必须完整：旧 → 新调用改写示例）；
4. 依赖窗更新：受影响模块 `DependsOn` 兼容窗；消费方文档（QUICKSTART/playbook）同步；
5. 回归：golden consumer（示例下游）升级演练全绿（R4 建成后执行）；
6. 发布记录：Go tag + npm 版本（R5 流水线落地后自动核验）。

## 4. Deprecation 先例（删符号）

1. 宽限 ≥ 1 个 minor 周期：保留符号 + `Deprecated:` 注释 + changelog 提示；
2. 宽限满后 major 发布移除，迁移说明必须存在。

## 5. 门禁

- registry 对 `KernelAPIRange` 的 fail-closed 校验（存在，保持）；
- 新流程自 v1.0.0（用户确认冻结面清单）起生效；生效前主线变更不受追溯。