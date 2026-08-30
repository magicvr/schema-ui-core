# R4 证据矩阵 · workspace-025（2026-08-30）

> 对应 VP-025 六条方向级退出判据；本矩阵为关门审计 A-001/A-002 的证据索引。

| # | 判据（VP-025） | 证据链接 | 状态 |
|---|----------------|----------|------|
| 1 | 导出闭环：当前生效配置 → 可移植配置包；往返一致；密钥按冻结规则排除/脱敏（快测 + CLI 实证） | 合同 §1（GOAL-002 D-002）；实现 `configpkg.go` buildPackage/buildExportTree（敏感登记表 + 宽规则不变量 · env 名提取）；测试 TestExport\*、TestRoundtrip（往返一致）；CLI 冒烟 `export -o` exit 0（产物 `${VAR}` 形态 · 无敏感键） | **verified** |
| 2 | diff 可核对：两包 / 包 vs 运行配置差量机器可读可断言（一致/仅差/冲突） | 合同 §2.2；实现 diffLeafMaps/loadConfigLeaf（扁平化 · 忽略信息性元数据 · `--against`）；测试 TestDiff\*（一致 `[]` / modify / add·remove / against / 忽略元数据）；CLI 冒烟 exit 0/1 | **verified** |
| 3 | dry-run 无副作用：预检覆盖校验与影响报告，成功/失败路径快测，无写副作用 | 合同 §2.3；实现 dryRun（结构校验 KnownFields + env fail-closed + 包→目标方向差量）；测试 TestDryRunPass（目标文件快照零副作用）、TestDryRunEnvMissingFailsClosed、TestDryRunChanges、TestDryRunInvalidPackage；CLI 冒烟 exit 0/1（fail-closed 实证） | **verified** |
| 4 | 导入不破坏：预检通过后应用；失败路径不破坏既有配置（I-025-004 冻结语义） | 合同 §2.4 + GOAL-004 D-002（**用户裁决方案 A**：备份 `.pre-import.bak` → `.tmp` → `server.LoadConfig` 装载校验 → `os.Rename` 原子替换；失败清 tmp、原文件不被触碰）；实现 cmdConfigImport；测试 TestImport\*（往返 diff `[]` / 备份内容 = 旧值 / 预检·坏包·装载失败 → 目标原样 + 无 tmp 泄漏 / 缺省文件）；CLI 冒烟（往返闭环 + `.bak` 生成） | **verified** |
| 5 | 边界保持：未改 Charter；未改 Profile 默认集/模块矩阵/Manifest 装配；热加载不进分母；密钥 fail-closed | 红线核账 `git diff --name-only cf68c7ce..HEAD`：代码面仅 4 文件（configpkg.go · configpkg_test.go · main.go · server/config.go 只读导出）——`internal/store` / `kernel/profile.go` / `protocol/upstream` / 迁移面**零触碰**（Root E-004 + GOAL-005 全量复核）；Charter 未动；热加载未引入；管理面未做 | **verified** |
| 6 | 审计闭合：开放 required finding = 0（或已合法闭合） | 子目标关门审计：GOAL-002 A-001 pass · GOAL-003 A-001 pass · GOAL-004 A-001 pass（0 required）；Root 关门双审：A-001 self（本矩阵 + 全链核对）· A-002 grok build independent（后台运行中 · 结果合并于 03-audit） | **进行中**（待 A-002 + VRev-055） |

## 回归与实证汇总

- 单元测试：`apps/api/cmd/schema-ui/configpkg_test.go` **22 用例全绿**（export×4 · dry-run×6 · import×5 · diff×6（5 个 TestDiff* + TestRoundtrip）· 对抗×1+1）——其中 F-001（plaintext secret 拒绝 ×2）与 F-002（type/range + 零副作用 ×2）为 A-002 响应后新增
- 全量回归：`go test ./...`（apps/api 49 包）PASS（R2/R3 各轮复核）+ **web vitest 90 文件 / 1186 用例 PASS**（A-002 F-008 响应 · 变更面不含 web）
- CLI 端到端：export（产物形态）→ diff（0/1 退出码）→ dry-run（0/1 · fail-closed · 类型校验）→ import（往返 diff `[]` · `.bak` 生成 · 失败不破坏 · 明文包拒绝）——全部实证留痕于各 E 条目
- **serve 进程级实证**（A-002 F-006 响应）：import 生成的配置启动 `schema-ui serve`（profile=admin · 8 modules）→ `/healthz` **200** `{"status":"ok",...}`
- Checkpoints：cf68c7ce → 70e9ecd7 → 48ba2ebd → 0f60fbc1 → 9983f206 → 98d33cc2 → f542c677 → 6a495a24（→ R4 A-002 响应批）