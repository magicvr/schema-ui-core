#!/usr/bin/env bash
#
# scripts/smoke.sh · S4 可复现 smoke 验收（GOAL-008 / I-008-002 v0.1.2 §5）
#
# 机器可判定的非破坏性冒烟：SM-001 参数/工具/安全前提 → SM-002 API readiness →
# SM-003 代理登录 → SM-004 当前身份 → SM-005 代表页路由。
# 仅在显式 --disposable 下执行 SM-006（种子可重复性，要求隔离 Compose project/volume）。
#
# 退出码：0=完整绿（SM-001～005 且 disposable SM-006 通过）｜2=参数、工具或安全
# 前提不满足（含不安全 destructive/隔离校验失败）｜3=readiness 30s 超时｜
# 4=登录/身份失败｜5=路由/数据失败｜6=SM-006 种子断言失败｜
# 8=部分绿（非 disposable，SM-006 未运行——不是 S4 完整绿，不得作为种子可重复证据）｜
# 70=未分类内部错误。
#
# 输入（env）：
#   API_BASE_URL         默认 http://localhost:8080
#   WEB_BASE_URL         默认 http://localhost:8081
#   SMOKE_USERNAME       默认 admin
#   SMOKE_PASSWORD       必填（无默认，禁止猜测 secret）
#   SMOKE_SEED_ID         默认 user-admin
#   SMOKE_EXPECTED_SEED_TOTAL  默认 1
#   SMOKE_ISOLATION_ID   仅 --disposable 必填：隔离 Compose project 名（机器校验
#                        运行中 project 与 db-data 卷均绑定该身份；不得指向默认开发库）
#   SMOKE_DISPOSABLE_CONFIRM   仅 --disposable 必填：必须为 yes（书面确认 disposable 语义）
# 安全：脚本不输出 token/password/secret；无 --disposable 不得执行种子 reset；
#       --disposable 时重启由脚本以显式隔离 project 执行（拒绝外部注入命令）。

set -u

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
WEB_BASE_URL="${WEB_BASE_URL:-http://localhost:8081}"
SMOKE_USERNAME="${SMOKE_USERNAME:-admin}"
SMOKE_PASSWORD="${SMOKE_PASSWORD:-}"
SMOKE_SEED_ID="${SMOKE_SEED_ID:-user-admin}"
SMOKE_EXPECTED_SEED_TOTAL="${SMOKE_EXPECTED_SEED_TOTAL:-1}"
SMOKE_ISOLATION_ID="${SMOKE_ISOLATION_ID:-}"
SMOKE_DISPOSABLE_CONFIRM="${SMOKE_DISPOSABLE_CONFIRM:-}"
DISPOSABLE=0

usage() {
  sed -n '2,32p' "$0"
  exit 2
}

if [ "$#" -gt 1 ]; then usage; fi
if [ "${1:-}" = "--disposable" ]; then DISPOSABLE=1; elif [ "${1:-}" != "" ]; then usage; fi

PASS=0; FAIL=0
smoke_line() { printf 'SM-%s=%s\n' "$1" "$2"; }

json_field() {
  # usage: json_field <json> <key>  → writes raw value to stdout (no trailing newline)
  local json="$1" key="$2"
  node -e 'let s="";process.stdin.on("data",d=>s+=d);process.stdin.on("end",()=>{try{const o=JSON.parse(s);const v=o[process.argv[1]];process.stdout.write(v===undefined||v===null?"":String(v))}catch(e){process.exit(1)}})' "$key" <<< "$json"
}

fail_check() {
  local id="$1" detail="$2" code="${3:-70}"
  smoke_line "$id" "FAIL"
  printf '  detail: %s\n' "$detail"
  printf '  (secret 已脱敏；见脚本说明)\n'
  FAIL=$((FAIL+1))
  exit "$code"
}

# F-008：disposable 必须机器可判定地运行在显式隔离的 Compose project/volume 上；
# 不满足任何安全前提一律按协议退出码 2（安全前提失败），且绝不 reset 普通开发库。
check_isolation() {
  local api_cid mount_info
  if ! command -v docker >/dev/null 2>&1; then
    printf 'SM-001=FAIL\n  detail: disposable 模式缺少工具 docker\n'
    return 1
  fi
  api_cid="$(docker compose -p "${SMOKE_ISOLATION_ID}" ps -q api 2>/dev/null | head -n1)"
  if [ -z "$api_cid" ]; then
    printf 'SM-001=FAIL\n  detail: disposable 需要运行中的隔离 compose project「%s」（docker compose -p %s ps 无 api 容器）\n' "$SMOKE_ISOLATION_ID" "$SMOKE_ISOLATION_ID"
    return 1
  fi
  mount_info="$(docker inspect -f '{{range .Mounts}}{{.Type}}|{{.Name}}|{{.Source}}{{end}}' "$api_cid" 2>/dev/null || true)"
  case "$mount_info" in
    *"${SMOKE_ISOLATION_ID}_db-data"*|*"|bind|"*"${SMOKE_ISOLATION_ID}"*)
      printf 'isolation: project=%s container=%s\n' "$SMOKE_ISOLATION_ID" "$api_cid"
      return 0 ;;
    *)
      printf 'SM-001=FAIL\n  detail: api 容器卷未绑定隔离 project「%s」的 db-data 卷（mounts=%s）\n' "$SMOKE_ISOLATION_ID" "$mount_info"
      return 1 ;;
  esac
}

# ---------------------------------------------------------------------------
# SM-001 · 参数、工具与安全前提
# ---------------------------------------------------------------------------
SM001=""
if [ -z "${SMOKE_PASSWORD}" ]; then SM001="缺少 SMOKE_PASSWORD（禁止猜测默认 secret）"; fi
for tool in bash curl node; do
  command -v "$tool" >/dev/null 2>&1 || SM001="${SM001}缺少工具 ${tool};"
done
if [ -z "$API_BASE_URL" ] || [ -z "$WEB_BASE_URL" ]; then SM001="${SM001}缺少 API_BASE_URL/WEB_BASE_URL;"; fi
if [ "$DISPOSABLE" = "1" ] && [ "$SMOKE_DISPOSABLE_CONFIRM" != "yes" ]; then
  SM001="${SM001}disposable 模式缺少 SMOKE_DISPOSABLE_CONFIRM=yes（书面确认 disposable 语义）;"
fi
if [ "$DISPOSABLE" = "1" ] && [ -z "$SMOKE_ISOLATION_ID" ]; then
  SM001="${SM001}disposable 模式缺少 SMOKE_ISOLATION_ID（隔离 compose project）;"
fi
if [ -n "$SM001" ]; then
  printf 'SM-001=FAIL\n  detail: %s\n' "$SM001"
  exit 2
fi
if [ "$DISPOSABLE" = "1" ] && ! check_isolation; then
  exit 2
fi
smoke_line "001" "PASS"

# ---------------------------------------------------------------------------
# SM-002 · API readiness（30s 内 /healthz 200 + status=ok）
# ---------------------------------------------------------------------------
ready=0
for _ in $(seq 1 30); do
  body="$(curl -fsS --max-time 3 "${API_BASE_URL}/healthz" 2>/dev/null || true)"
  if [ -n "$body" ]; then
    st="$(json_field "$body" status)"
    if [ "$st" = "ok" ]; then ready=1; break; fi
  fi
  sleep 1
done
if [ "$ready" != "1" ]; then
  printf 'SM-002=FAIL\n  detail: /healthz 未在 30s 内返回 status=ok (%s)\n' "$API_BASE_URL"
  exit 3
fi
smoke_line "002" "PASS"

login_ok=0; ACCESS_TOKEN=""
for _ in $(seq 1 3); do
  login_body="$(curl -fsS --max-time 5 -X POST "${WEB_BASE_URL}/api/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"username\":\"${SMOKE_USERNAME}\",\"password\":\"${SMOKE_PASSWORD}\"}" 2>/dev/null || true)"
  if [ -n "$login_body" ]; then
    tok="$(json_field "$login_body" accessToken || true)"
    if [ -n "$tok" ]; then ACCESS_TOKEN="$tok"; login_ok=1; break; fi
  fi
  sleep 1
done

# ---------------------------------------------------------------------------
# SM-003 · 代理登录
# ---------------------------------------------------------------------------
if [ "$login_ok" != "1" ]; then
  printf 'SM-003=FAIL\n  detail: 经 WEB_BASE_URL /api/auth/login 登录失败或未返回非空 accessToken（用户名 %s）\n' "$SMOKE_USERNAME"
  exit 4
fi
smoke_line "003" "PASS"

# ---------------------------------------------------------------------------
# SM-004 · 当前身份（Bearer /api/accounts/me）
# ---------------------------------------------------------------------------
me="$(curl -fsS --max-time 5 "${WEB_BASE_URL}/api/accounts/me" -H "Authorization: Bearer ${ACCESS_TOKEN}" 2>/dev/null || true)"
if [ -n "$me" ]; then
  u="$(json_field "$me" user || true)"
  f="$(json_field "$me" features || true)"
  if [ -n "$u" ] && [ -n "$f" ]; then
    smoke_line "004" "PASS"
  else
    printf 'SM-004=FAIL\n  detail: /api/accounts/me 200 但缺少 user 或 features 投影\n'
    exit 4
  fi
else
  printf 'SM-004=FAIL\n  detail: /api/accounts/me 未返回 200（Bearer 无效或已过期）\n'
  exit 4
fi

# ---------------------------------------------------------------------------
# SM-005 · 代表页路由（SPA root 挂载标记）
# ---------------------------------------------------------------------------
spa="$(curl -fsS --max-time 5 "${WEB_BASE_URL}/users" 2>/dev/null || true)"
if [ -n "$spa" ] && printf '%s' "$spa" | grep -q 'id="root"'; then
  smoke_line "005" "PASS"
else
  printf 'SM-005=FAIL\n  detail: %s 未返回含 id="root" 的 SPA 文档\n' "${WEB_BASE_URL}/users"
  exit 5
fi

# ---------------------------------------------------------------------------
# SM-006 · 种子重复性（仅 disposable/隔离环境，S4 必检）
# ---------------------------------------------------------------------------
if [ "$DISPOSABLE" = "1" ]; then
  check_seed() {
    local expect="$1" detail="$2"
    local list total has_user
    list="$(curl -fsS --max-time 5 "${WEB_BASE_URL}/api/users?pageSize=100" -H "Authorization: Bearer ${ACCESS_TOKEN}" 2>/dev/null || true)"
    if [ -z "$list" ]; then
      printf 'SM-006=FAIL\n  detail: %s\n' "$detail（列表请求失败）"
      exit 6
    fi
    total="$(json_field "$list" total || true)"
    has_user=0
    printf '%s' "$list" | node -e 'let s="";process.stdin.on("data",d=>s+=d);process.stdin.on("end",()=>{try{const o=JSON.parse(s);const items=o.items||[];const hit=items.some(r=>r.id===process.argv[1]||r.username==="admin");process.exit(hit?0:1)}catch(e){process.exit(2)}})' "$SMOKE_SEED_ID"
    [ "$?" = "0" ] && has_user=1
    if [ "$total" != "$expect" ] || [ "$has_user" != "1" ]; then
      printf 'SM-006=FAIL\n  detail: %s（期望 total=%s 且含 %s/admin，实际 total=%s）\n' "$detail" "$expect" "$SMOKE_SEED_ID" "$total"
      exit 6
    fi
  }

  # 首次断言：空库种子后 total == 期望且含 user-admin / admin
  check_seed "$SMOKE_EXPECTED_SEED_TOTAL" "首次种子断言失败"

  # 重启由脚本以显式隔离 project 执行（F-008：拒绝外部注入任意命令）
  if ! docker compose -p "${SMOKE_ISOLATION_ID}" restart api >/dev/null 2>&1; then
    printf 'SM-006=FAIL\n  detail: docker compose -p %s restart api 执行失败\n' "$SMOKE_ISOLATION_ID"
    exit 70
  fi
  # R-011：重启后必须重新判定 readiness，失败按退出码 3，不沿用初始 ready
  ready=0
  for _ in $(seq 1 30); do
    body="$(curl -fsS --max-time 3 "${API_BASE_URL}/healthz" 2>/dev/null || true)"
    if [ -n "$body" ] && [ "$(json_field "$body" status)" = "ok" ]; then ready=1; break; fi
    sleep 1
  done
  if [ "$ready" != "1" ]; then
    printf 'SM-006=FAIL\n  detail: 重启后 /healthz 未在 30s 内恢复 status=ok（%s）\n' "$API_BASE_URL"
    exit 3
  fi
  login_body="$(curl -fsS --max-time 5 -X POST "${WEB_BASE_URL}/api/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"username\":\"${SMOKE_USERNAME}\",\"password\":\"${SMOKE_PASSWORD}\"}" 2>/dev/null || true)"
  ACCESS_TOKEN="$(json_field "$login_body" accessToken || true)"
  if [ -z "$ACCESS_TOKEN" ]; then
    printf 'SM-006=FAIL\n  detail: 重启后重新登录失败\n'
    exit 4
  fi
  check_seed "$SMOKE_EXPECTED_SEED_TOTAL" "重启后种子断言失败（重复播种或用户丢失）"

  smoke_line "006" "PASS"
else
  printf 'SM-006=SKIP（非 disposable；S4 完整绿需 --disposable 且 SM-006=PASS）\n'
fi

if [ "$DISPOSABLE" = "1" ]; then
  printf 'SMOKE RESULT: PASS (SM-001~005 + SM-006)\n'
  exit 0
else
  printf 'SMOKE RESULT: PARTIAL (SM-001~005; SM-006 未运行——非 S4 完整绿，不得作为种子可重复证据；需 --disposable + 隔离环境)\n'
  exit 8
fi
