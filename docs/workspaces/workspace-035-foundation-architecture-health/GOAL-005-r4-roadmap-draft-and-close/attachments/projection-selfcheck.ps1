﻿﻿# projection-selfcheck.ps1 · workspace-035 投影一致性自检（只读）
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

# ---- 检查 2：未来式语态（不得把已发生的响应写成未来条件；03-audit/ 台账属历史，排除）----
$p2 = @()
Get-ChildItem $ws -Recurse -Filter '*.md' -File | Where-Object { $_.DirectoryName -notmatch '\\03-audit($|\\)' } | ForEach-Object {
  $lines = (Read-Text $_.FullName) -split "`r?`n"
  for ($i = 0; $i -lt $lines.Count; $i++) {
    if ($lines[$i] -match '待\s*/?govern\s*响应闭合后方可宣称' -and $lines[$i] -notmatch '历史|原文') {
      $p2 += "$($_.Name):$($i+1)"
    }
  }
}
Report '2 未来式语态' $p2

# ---- 检查 3：VP-035 当前版本投影一致（子句级：只检查「VP-035 的当前状态子句」内的版本）----
# 判定规则（A-016 F-013 反例驱动）：
#   - 一行可含多个 VP/多个版本；只取「显式状态标记」（`active` vX / active vX / 激活记录 vX / 激活 vX）关联的版本。
#   - 若该版本 token 前 10 字符内含「激活/历史/时点」，视为历史记录，跳过。
#   - 与 VP-035 的距离：任一 VP-035 token 之后 400 字符内视为同一描述子句。
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
      # VP-035 出现位置（最后一个即可覆盖「同一子句内随后出现版本」的常见写法）
      $vpHits = [regex]::Matches($line, 'VP-035') | ForEach-Object { $_.Index }
      $statusRe = [regex]'(?:`active`|(?<![\w`])active|激活记录?|时点记录)\s*`?v(\d+\.\d+\.\d+)'
      foreach ($m in $statusRe.Matches($line)) {
        $ver = $m.Groups[1].Value
        if ($ver -eq $vpVer) { continue }
        $pre = $line.Substring([Math]::Max(0, $m.Index - 10), [Math]::Min(10, $m.Index))
        if ($pre -match '激活|历史|时点') { continue }   # 形如「2026-09-09 激活 v0.2.0」
        $near = $false
        foreach ($h in $vpHits) { if ($m.Index -ge $h -and ($m.Index - $h) -le 400) { $near = $true } }
        if (-not $near) { continue }                      # 属其它 VP 的描述子句
        $p3 += "$($f.Name):$($i+1) VP-035 当前子句出现 v$ver（VP 当前 v$vpVer）"
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
# 更强的断言：基线中每个 trigger-gated 的 RT-* 行 ID 必须仍为 trigger-gated（防「释放一行 + 新增一行」抵消计数）
$gatedNow = @()
Select-String -Path (Join-Path $repo 'docs/vision/roadmap.md') -Pattern '^\|\s*(RT-[A-Z0-9]+)\s*\|.*trigger-gated' | ForEach-Object { $gatedNow += $matches[1] }
$gatedBase = @()
(git show 'ebe6013c:docs/vision/roadmap.md') | Select-String -Pattern '^\|\s*(RT-[A-Z0-9]+)\s*\|.*trigger-gated' | ForEach-Object { $gatedBase += $matches[1] }
$lost = $gatedBase | Where-Object { $gatedNow -notcontains $_ }
Write-Output "      trigger-gated RT-* ID: 基线=$($gatedBase.Count) 现=$($gatedNow.Count) 丢失=$($lost.Count)"
if ($lost.Count -gt 0) { $p5 += "基线 trigger-gated 行被释放: $($lost -join ', ')" }
Report '5 边界守恒' $p5

# ---- 检查 6：最近提交的 version/updated 核账（脚本化，输出逐文件结论）----
$p6 = @()
$commits = (git log --format=%H -6)
$rootPath = $repo
$seen = @{}
foreach ($c in $commits) {
  $files = (git show --name-only --format= --diff-filter=AM $c) | Where-Object { $_ -match '\.md$' }
  foreach ($rel in $files) {
    if ($seen.ContainsKey($rel)) { continue }
    $seen[$rel] = $true
    $abs = Join-Path $rootPath $rel
    if (-not (Test-Path $abs)) { continue }          # 已删除/重命名
    $h = Get-Front $abs
    if ($h.Count -eq 0) { continue }                  # 非治理 md（无 frontmatter）
    $commitDate = (git show -s --format=%ad --date=short $c)
    if ($h['updated'] -ne $commitDate) {
      $p6 += "${rel} -> updated=$($h['updated']) 但内容变更于 $commitDate（$($c.Substring(0,8))）"
    }
  }
}
Report '6 最近提交 version/updated 核账' $p6

if ($script:fail -eq 0) { Write-Output 'ALL CHECKS PASS'; exit 0 } else { Write-Output "$script:fail CHECK(S) FAILED"; exit 1 }
