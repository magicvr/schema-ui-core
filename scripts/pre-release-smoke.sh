#!/usr/bin/env bash
#
# scripts/pre-release-smoke.sh · 发版前完整冒烟（生产 CSP + 真实浏览器 + 隔离种子）
#
# 组合现有验收：
#   1. 在独立 Compose project（不自复用开发库）构建并启动生产栈。
#   2. 运行 scripts/smoke.sh --disposable（SM-001~005 + SM-006 种子可重复性 +
#      可选 SM-007 Profile/Manifest + SMOKE_CSP=1 → 真实浏览器 + 生产 CSP 头 SM-008）。
#   3. 结束后默认 docker compose -p <隔离 project> down -v 清理（保留镜像/缓存）。
#
# 说明：compose.yaml 按 W7 F-008 不发布 API 宿主端口；本脚本为本地/CI 冒烟，
#   通过临时 override 仅把 API 端口 127.0.0.1:25080 发布到本机 loopback，供 smoke.sh
#   的 SM-002/006 readiness 检查使用（生产部署仍遵循 compose.yaml 不发布）。
#
# 退出码：透传 smoke.sh（0=完整绿；2/3/4/5/6/7/8/70 含义见 smoke.sh 头）。
#
# 前置：
#   docker（Compose v2）、bash、node（受控于 apps/web 的 @playwright/test Chromium）
#   .env 或环境变量提供 AUTH_JWT_SECRET / ADMIN_INITIAL_PASSWORD（fail-closed）
#   Playwright Chromium 已安装（apps/web 下 npx playwright install chromium）
#
# 输入（env，均可选）：
#   PRERELEASE_PROFILE   期望 Profile（默认从 apps/api/configs/config.yaml 读取）
#   WEB_HOST_PORT        覆盖宿主 Web 端口（默认从 .env WEB_HOST_PORT 或 25081）
#   SMOKE_USERNAME       默认 admin
#   SMOKE_PASSWORD_NEW   强制改密后的新密码（默认 = 初始密码-<profile>-smoke）
#   KEEP_STACK=1         结束时保留 Compose 栈（默认 down -v）
#
# 用法：
#   bash scripts/pre-release-smoke.sh
set -u

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

PROJECT="schema-ui-prerelease-$(date +%s)-$$"

env_val() {
  local key="$1" val
  val="$(grep -m1 "^${key}=" .env 2>/dev/null | cut -d= -f2- || true)"
  printf '%s' "$val"
}

fail() { printf 'PRE-RELEASE FAIL: %s\n' "$*" >&2; exit 2; }

# ---- 工具与前提 ----
for tool in docker node curl; do command -v "$tool" >/dev/null 2>&1 || fail "缺少工具 ${tool}"; done
[ -f scripts/smoke.sh ] || fail "缺少 scripts/smoke.sh"
[ -f apps/web/scripts/check-prod-csp.mjs ] || fail "缺少 apps/web/scripts/check-prod-csp.mjs"

# ---- 密钥 / 端口 / profile ----
ADMIN_INITIAL_PASSWORD="${ADMIN_INITIAL_PASSWORD:-$(env_val ADMIN_INITIAL_PASSWORD)}"
[ -n "$ADMIN_INITIAL_PASSWORD" ] || fail "缺少 ADMIN_INITIAL_PASSWORD（.env 或环境变量；禁止猜测默认 secret）"
AUTH_JWT_SECRET="${AUTH_JWT_SECRET:-$(env_val AUTH_JWT_SECRET)}"
[ -n "$AUTH_JWT_SECRET" ] || fail "缺少 AUTH_JWT_SECRET（.env 或环境变量；Compose fail-closed）"

WEB_HOST_PORT="${WEB_HOST_PORT:-$(env_val WEB_HOST_PORT)}"
WEB_HOST_PORT="${WEB_HOST_PORT:-25081}"
WEB_BASE_URL="http://127.0.0.1:${WEB_HOST_PORT}"

PROFILE="${PRERELEASE_PROFILE:-}"
if [ -z "$PROFILE" ]; then
  PROFILE="$(grep -m1 '^  profile:' apps/api/configs/config.yaml 2>/dev/null | sed -E 's/^[[:space:]]*profile:[[:space:]]*//; s/"//g' || true)"
fi
PROFILE="${PROFILE:-mvp}"
# W16-F01: fresh seed 的 admin 强制改密（must_change_password=1）。冒烟前用
# 初始密码登录并换成新密码；后续 smoke 全部用新密码（隔离卷，删除即清）。
NEW_PASSWORD="${SMOKE_PASSWORD_NEW:-${ADMIN_INITIAL_PASSWORD}-${PROFILE}-smoke}"

# 临时 override：仅为本机/CI 冒烟把 API 端口 loopback 发布给 smoke.sh（生产仍不发布）
OVERRIDE_DIR="$(mktemp -d)"
OVERRIDE_FILE="$OVERRIDE_DIR/compose.smoke.yaml"
cat > "$OVERRIDE_FILE" <<'EOF'
services:
  api:
    ports:
      - "127.0.0.1:25080:25080"
EOF

printf 'PRE-RELEASE SMOKE | project=%s | profile=%s | web=%s\n' "$PROJECT" "$PROFILE" "$WEB_BASE_URL"

cleanup() {
  if [ "${KEEP_STACK:-}" = "1" ]; then
    printf 'PRE-RELEASE: KEEP_STACK=1，保留 %s 栈供排查。\n' "$PROJECT"
  else
    docker compose -p "$PROJECT" down -v >/dev/null 2>&1 || true
    printf 'PRE-RELEASE: 已清理 %s 栈。\n' "$PROJECT"
  fi
  rm -rf "$OVERRIDE_DIR"
}
trap cleanup EXIT

# ---- 启动隔离生产栈 ----
docker compose -p "$PROJECT" -f compose.yaml -f "$OVERRIDE_FILE" down -v >/dev/null 2>&1 || true
docker compose -p "$PROJECT" -f compose.yaml -f "$OVERRIDE_FILE" up -d --build || fail "docker compose up -d --build 失败"

# ---- 等待 web 可访问 ----
printf 'PRE-RELEASE: 等待 %s 可访问 ...\n' "$WEB_BASE_URL"
up=0
for _ in $(seq 1 60); do
  if curl -fsS --max-time 3 "$WEB_BASE_URL/" >/dev/null 2>&1; then up=1; break; fi
  sleep 2
done
[ "$up" = "1" ] || fail "web 未在 120s 内可访问（${WEB_BASE_URL}）"

# ---- 强制改密 bootstrap（W16-F01：fresh admin 强制改密） ----
printf 'PRE-RELEASE: admin 强制改密 bootstrap ...\n'
SMOKE_USERNAME="${SMOKE_USERNAME:-admin}"
BOOT_TOKEN="$( \
  SMU="$SMOKE_USERNAME" SMP="$ADMIN_INITIAL_PASSWORD" WEB="$WEB_BASE_URL" node -e '
    fetch(process.env.WEB+"/api/auth/login",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({username:process.env.SMU,password:process.env.SMP})})
      .then(async r=>{const b=await r.json();process.stdout.write(r.ok?(b.accessToken||""):"")})
      .catch(()=>{})
  ' || true)"
[ -n "$BOOT_TOKEN" ] || fail "初始 admin 登录失败（无法获取 bootstrap token）"

SMU="$SMOKE_USERNAME" SMP="$ADMIN_INITIAL_PASSWORD" NEW="$NEW_PASSWORD" TOK="$BOOT_TOKEN" WEB="$WEB_BASE_URL" node -e '
  fetch(process.env.WEB+"/api/account/password",{method:"POST",headers:{Authorization:"Bearer "+process.env.TOK,"Content-Type":"application/json"},body:JSON.stringify({currentPassword:process.env.SMP,newPassword:process.env.NEW})})
    .then(async r=>{process.exit(r.ok?0:1)})
    .catch(()=>process.exit(1))
' || fail "强制改密失败（/api/account/password）"

# ---- 完整冒烟（disposable + CSP 真实浏览器，使用新密码） ----
SMOKE_USERNAME="$SMOKE_USERNAME" \
SMOKE_PASSWORD="$NEW_PASSWORD" \
SMOKE_EXPECTED_PROFILE="$PROFILE" \
SMOKE_ISOLATION_ID="$PROJECT" \
SMOKE_DISPOSABLE_CONFIRM=yes \
SMOKE_CSP=1 \
API_BASE_URL="http://127.0.0.1:25080" \
WEB_BASE_URL="$WEB_BASE_URL" \
bash scripts/smoke.sh --disposable
rc=$?

if [ "$rc" = "0" ]; then
  printf 'PRE-RELEASE SMOKE RESULT: PASS\n'
else
  printf 'PRE-RELEASE SMOKE RESULT: FAIL (exit %s)\n' "$rc"
fi
exit "$rc"