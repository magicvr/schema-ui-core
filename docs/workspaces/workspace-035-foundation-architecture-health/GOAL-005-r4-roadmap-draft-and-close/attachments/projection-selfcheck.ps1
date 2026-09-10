# projection-selfcheck.ps1 · workspace-035 投影一致性自检（只读）
#
# 用途：在声明任何 finding `fixed`、请求独立复核、或关闭目标之前运行。
#       把「同一事实的其它投影是否同步」从人肉记忆改为可重复执行的检查。
# 用法（仓库根；Windows PowerShell 5.1 或 PowerShell 7）：
#   powershell -NoProfile -ExecutionPolicy Bypass -File docs/workspaces/workspace-035-foundation-architecture-health/GOAL-005-r4-roadmap-draft-and-close/attachments/projection-selfcheck.ps1
# 注意：文件带 UTF-8 BOM；Windows PowerShell 5.1 依赖 BOM 才能正确解析中文，请勿去掉 BOM。
# 退出码：0 = 全部通过；1 = 存在失败项。
# 只读：不写入任何文件。

$ErrorActionPreference = 'Stop'
$repo = (Get-Location).Path
$ws   = Join-Path $repo 'docs/workspaces/workspace-035-foundation-architecture-health'
$goal5 = Join-Path $ws 'GOAL-005-r4-roadmap-draft-and-close'
$script:fail = 0

function Read-Text([string]$p) { [System.Text.Encoding]::UTF8.GetString([System.IO.File]::ReadAllBytes($p)) }

function Get-Front([string]$p) {
  $lines = (Read-Text $p) -split "`r?`n"
  $h = @{}
  if ($lines[0] -ne '---') { return $h }
  for ($i = 1; $i -lt $lines.Count; $i++) {
    if ($lines[$i] -eq '---') { break }
    if ($lines[$i] -match '^([A-Za-z_]+):\s*(.*)$') { $h[$matches[1]] = $matches[2].Trim() }
  }
  return $h
}

function Report([string]$name, [string[]]$problems) {
  if ($problems.Count -eq 0) { Write-Output "PASS  $name" }
  else { Write-Output "FAIL  $name"; $problems | ForEach-Object { Write-Output "        $_" }; $script:fail++ }
}

# ---- 检查 1：审计编号自指（条目正文不得把“下一次独立复核”写成本条目编号）----
$p1 = @()
Get-ChildItem $goal5 -Recurse -Filter 'A-0*.md' | ForEach-Object {
  $self = ($_.BaseName -split '-')[1]
  $lines = (Read-Text $_.FullName) -split "`r?`n"
  for ($i = 0; $i -lt $lines.Count; $i++) {
    if ($lines[$i] -match '下一次[^。]{0,40}?(A-\d{3})' -and $matches[1] -eq $self) {
      $p1 += "$($_.Name):$($i+1) 自指 $self"
    }
  }
}
Report '1 审计编号自指' $p1

# ---- 检查 2：未来式语态（不得把已发生的响应写成未来条件）----
$p2 = @()
Get-ChildItem $ws -Recurse -Filter '*.md' -File | ForEach-Object {
  $lines = (Read-Text $_.FullName) -split "`r?`n"
  for ($i = 0; $i -lt $lines.Count; $i++) {
    if ($lines[$i] -match '待\s*/?govern\s*响应闭合后方可宣称' -and $lines[$i] -notmatch '历史|原文') {
      $p2 += "$($_.Name):$($i+1)"
    }
  }
}
Report '2 未来式语态' $p2

# ---- 检查 3：VP-035 当前版本投影一致（覆盖 roadmap 全部命中；显式历史记录允许旧版本）----
$p3 = @()
$vpPath = Join-Path $repo 'docs/vision/plans/VP-035-foundation-architecture-health.md'
$vpVer = (Get-Front $vpPath)['version']
if (-not $vpVer) { $p3 += 'plan: frontmatter 缺 version' }
else {
  Write-Output "      VP-035 当前 version = $vpVer"
  $scan = @()
  # 审计台账（03-audit/ 下的条目）是历史意见记录，正文引用的是当时时点，天然属历史：整目录排除。
  $scan += Get-ChildItem $ws -Recurse -Filter '*.md' -File | Where-Object { $_.DirectoryName -notmatch '\\03-audit($|\\)' }
  $scan += Get-Item (Join-Path $repo 'docs/vision/roadmap.md')
  $scan += Get-Item (Join-Path $repo 'docs/vision/workspaces.md')
  foreach ($f in $scan) {
    $lines = (Read-Text $f.FullName) -split "`r?`n"
    for ($i = 0; $i -lt $lines.Count; $i++) {
      $line = $lines[$i]
      if ($line -notmatch 'VP-035') { continue }
      # 显式历史/激活时点记录允许保留旧版本号
      if ($line -match '激活记录|历史|时点记录|2026-09-0\d 激活|activation|VRev-08[0-9]') { continue }
      foreach ($m in [regex]::Matches($line, 'v(\d+\.\d+\.\d+)')) {
        if ($m.Groups[1].Value -eq $vpVer) { continue }
        $p3 += "$($f.Name):$($i+1) 当前投影出现 v$($m.Groups[1].Value)（VP 当前 v$vpVer，且未标为历史）"
      }
    }
  }
}
Report '3 VP-035 当前版本投影' $p3

# ---- 检查 4：workspace 内 frontmatter 必备字段 + 00-meta id 与目录名一致 ----
$p4 = @()
Get-ChildItem $ws -Recurse -Filter '*.md' -File | ForEach-Object {
  $h = Get-Front $_.FullName
  if ($h.Count -eq 0) { $p4 += "$($_.Name): 缺 frontmatter"; return }
  foreach ($k in 'status','created','updated','parent','version') {
    if (-not $h.ContainsKey($k)) { $p4 += "$($_.Name): 缺 $k" }
  }
}
Get-ChildItem $ws -Directory | ForEach-Object {
  $meta = Join-Path $_.FullName '00-meta.md'
  if (Test-Path $meta) { $id = (Get-Front $meta)['id']; if ($id -ne $_.Name) { $p4 += "$($_.Name): id=$id 与目录名不一致" } }
}
Report '4 frontmatter 必备字段与 id 一致性' $p4

# ---- 检查 5：边界守恒（apps 零变更；roadmap trigger-gated 计数不下降）----
$p5 = @()
$apps = git diff --name-only ebe6013c..HEAD -- apps
if ($apps) { $p5 += "apps/** 有变更: $($apps -join ', ')" }
$base = (git show 'ebe6013c:docs/vision/roadmap.md' | Select-String -Pattern 'trigger-gated' -AllMatches | Measure-Object).Count
$now = (Select-String -Path (Join-Path $repo 'docs/vision/roadmap.md') -Pattern 'trigger-gated' -AllMatches | Measure-Object).Count
Write-Output "      trigger-gated 基线=$base 现=$now"
if ($now -lt $base) { $p5 += "trigger-gated 计数下降: $base -> $now" }
Report '5 边界守恒' $p5

if ($script:fail -eq 0) { Write-Output 'ALL CHECKS PASS'; exit 0 } else { Write-Output "$script:fail CHECK(S) FAILED"; exit 1 }
