---
id: D-001
doc: decision-entry
goal: GOAL-003-dual-key-jwt
status: accepted
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
---

# D-001 · R2 双密钥语义冻结（关闭 I-003）与实施方案

## I-003 裁定（verified）

| 问题 | 决定 | 理由 / 证据 |
|------|------|-------------|
| 重叠窗语义：旧 access 可验多久 | **重叠窗 = previous 配置存续期**。配置着 previous 且进程以其启动，旧 key 签发的未过期 access 即可验；退役 = 操作者在 ≥`access_ttl`（默认 15m）后移除 previous 并重启。无时间窗状态、无倒计时 | VP 退出 1 默认措辞（previous 可验）；不新增存储/时钟状态；操作程序已在 R1 样例注释写明 |
| 是否使用 JWT `kid` | **不使用**。token 头与 claims 保持逐字节现状；双候选密钥按 current → previous 试验序 | 两把 HMAC 候选即可满足语义；加 kid 改变 token 形状、扩大兼容面，本波退出分母不需要 |
| refresh 是否受签名密钥轮换影响 | **不受影响**（opaque）。refresh 为 256-bit CSPRNG 随机值 + SHA-256 落库比对，签名密钥只作用于 access | `auth.NewOpaqueToken` / `HashToken` / `RefreshTokenByHash`（auth.go:400-409, 242）；Root D-002 §2 |

## 实施方案

1. **auth 包**：
   - `Authenticator` 增加 `previousSecret []byte` 字段；
   - 新构造器 `NewWithRepositoryAndPrevious(current, previous []byte, ...)`；既有 `New` / `NewWithRepository` 签名不变（= 单密钥），全部既有调用点零迁移；
   - 中间件验签改走 `verifyAccess(raw)`：先 `ParseAccessToken(current)`，失败且配置了 previous 再试 `ParseAccessToken(previous)`；两次尝试都强制过期与方法检查——回退不能延长任何 token 的寿命（过期 token 在两次尝试下同样拒绝）；
   - 签发路径（`issue()`）不动，天然只用 current。
2. **composition**：`newAuthenticator` 直接读 `cfg.AuthJWTSecretPrevious`（该函数已有 cfg 参数）传入新构造器；`NewApp` 对外签名不变。
3. **测试矩阵**（auth 包单测，复用 `testsupport.OpenStore` 种子 admin=`user-admin`）：
   - 重叠窗内旧 key token 通过中间件验签；
   - previous 移除（单密钥）后同一旧 token 拒绝（401）；
   - 双密钥下 Login 签发的新 access 只能被 current 验证、不能被 old 验证（签发只用 current）;
   - previous 签发但已过期的 token 在双密钥下仍拒绝（回退不延长寿命）;
   - 单密钥构造器的行为与既有断言不变（回归由全套件兜底）。

## 为什么

- 构造器注入而非 setter：保持 Authenticator 构造后不可变，避免运行期换钥的隐式状态；与现有 `Set*Recorder` 钩子不同，密钥是安全边界不是可选行为钩子。
- 试验序 current→previous：正常态（绝大多数请求）一次验证命中，previous 只在轮换窗承担尾部流量。
- composition 不改 NewApp 签名：previous 无需 main 层 fallback 解析（空 = 单密钥），最小 diff 降低回归面。

## 未选方案

- `kid` 头 + 密钥环查找：token 形状变更、无必要（两把候选封顶）。
- 时间戳式重叠窗（如 previous 仅 N 分钟有效）：需要持久化轮换时刻或推导规则，复杂度与退出分母不匹配；操作者控制退役时机等价且可核对。
- 解析失败一律重试第二把（含格式损坏 token）：实现上即"任何错误再试一次"，无害（结果同为拒绝）且代码最简；不做错误分类分流。
