---
id: GOAL-001-outbound-mail
doc: decision-entry
record_id: D-005
status: accepted
parent: null
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

## D-005 · R4 readyz 探测、显式路径证据与关门叙事（关闭 I-005 / I-006）

### 触发

Root 路线图 R4 门禁：仅显式配置后 `readyz` 扩依赖；显式配置后可核对至少一封投递；I-005/I-006 最晚 R4 关闭。

### 决定

采纳子目标 `GOAL-005-r4-readyz-evidence` D-001/D-002：

1. **readyz**：`(*mail.SMTP).Ping` 仅在显式配置时经 `RegisterWithMFAProbes` 进入 readyz；capture 缺省 probe=nil（未配置不得 not-ready）。装配与 objectStore 同构（NewApp 直接构造）。
2. **证据面**：离线 TLS harness（Ping 全绿/明文必败/envelope+DATA+AUTH 断言）构成"与生产合同等价的 harness"；live 测试 env-gated 留给 operator 复跑。
3. **I-005 → verified**：HTML/MIME 不进本波退出分母（合同仅纯文本 TextBody）。**I-006 → verified（叙事投影闭合 V-F071）**：重启生效，无热加载。

### 为什么

探针与投递共用唯一冻结拨号路径——readyz 绿即投递路径 TLS/证书/可达性全部可核对。

### 未选方案

见 `GOAL-005` D-001/D-002「未选方案」节。
