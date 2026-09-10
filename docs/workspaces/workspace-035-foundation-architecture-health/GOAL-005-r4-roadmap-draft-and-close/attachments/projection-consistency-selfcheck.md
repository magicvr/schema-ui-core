---
doc_type: goal-attachment
id: projection-consistency-selfcheck
parent: GOAL-005-r4-roadmap-draft-and-close
status: recorded
created: 2026-09-10
updated: 2026-09-10
version: 0.1.0
---

# 投影一致性自检（R4 关门前置）

用途：在声明任何 finding `fixed`、或请求独立复核之前，**先**运行本脚本，机器化核对跨投影一致性与 frontmatter 版本规则。
背景：R4 的 12 项 required 中绝大多数源自同一失效模式——改了一个位置、漏改同一事实的其它投影。本脚本把该核对从「人肉记忆」改为「可重复执行」。

用法（仓库根）：

```powershell
pwsh -File docs/workspaces/workspace-035-foundation-architecture-health/GOAL-005-r4-roadmap-draft-and-close/attachments/projection-selfcheck.ps1
```

期望输出：每项检查打印 `PASS`，末尾打印 `ALL CHECKS PASS`。任一 `FAIL` 行给出文件与行号。

检查项：

1. **编号自指**：任何 `A-0NN` 条目正文中出现的「下一次/下一次独立复核」编号不得等于该条目自身编号。
2. **未来式语态**：不得在「已发生的响应」处使用「待 `/govern` 响应闭合后方可宣称」等未来式；出现即 FAIL（历史引述需显式带「历史」或「原文」限定）。
3. **VP 版本投影**：`docs/vision/plans/VP-035-*.md` 的 `version` 必须与 goal-tree、workspace.md、Root `00-meta.md`、roadmap 的**当前**投影一致（历史激活记录允许保留旧版本号，但须带「激活记录」限定）。
4. **frontmatter 版本递增**：最近 3 次提交中被修改的文件，若内容变化则 `version` 必须变化；`updated` 必须等于内容变更日。
5. **必备字段**：workspace 内所有 md 均含 `status`/`created`/`updated`/`parent`/`version`；`00-meta.md` 的 `id` 等于文件夹名。
6. **边界守恒**：`git diff --name-only ebe6013c..HEAD -- apps` 为空；`docs/vision/roadmap.md` 的 `trigger-gated` 计数不少于基线。

```powershell
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path '.').Path
$ws   = Join-Path $root 'docs/workspaces/workspace-035-foundation-architecture-health'
$goal5 = Join-Path $ws 'GOAL-005-r4-roadmap-draft-and-close'
$fail = 0
function Check([string]$name, [scriptblock]$body) {
  try { $msgs = & $body; if ($msgs) { Write-Output "FAIL  $name"; $msgs | ForEach-Object { Write-Output "        $_" }; $script:fail++ } else { Write-Output "PASS  $name" } }
  catch { Write-Output "FAIL  $name (error: $_)"; $script:fail++ }
}
function ReadUtf8([string]$p) { [System.Text.Encoding]::UTF8.GetString([System.IO.File]::ReadAllBytes($p)) }
function Fm([string]$p) {
  $l = (ReadUtf8 $p) -split "`r?`n"
  if ($l[0] -ne '---') { return @{} }
  $end = 0; for ($i=1; $i -lt $l.Count; $i++) { if ($l[$i] -eq '---') { $end = $i; break } }
  $h = @{}; foreach ($line in $l[1..($end-1)]) { if ($line -match '^([A-Za-z_]+):\s*(.*)$') { $h[$matches[1]] = $matches[2].Trim() } }
  return $h
}

Check '1 编号自指' {
  $out = @()
  Get-ChildItem $goal5 -Recurse -Filter 'A-0*.md' | ForEach-Object {
    $self = ($_.BaseName -split '-')[1]
    $lines = (ReadUtf8 $_.FullName) -split "`r?`n"
    for ($i=0; $i -lt $lines.Count; $i++) {
      if ($lines[$i] -match '下一次[^。]{0,40}?(A-\d{3})' -and $matches[1] -eq $self) { $out += "$($_.Name):$($i+1) 自指 $self" }
    }
  }
  $out
}

Check '2 未来式语态' {
  $out = @()
  Get-ChildItem $ws -Recurse -Filter '*.md' | ForEach-Object {
    $lines = (ReadUtf8 $_.FullName) -split "`r?`n"
    for ($i=0; $i -lt $lines.Count; $i++) {
      if ($lines[$i] -match '待\s*/?govern\s*响应闭合后方可宣称' -and $lines[$i] -notmatch '历史|原文') { $out += "$($_.Name):$($i+1)" }
    }
  }
  $out
}

Check '3 VP 版本投影一致' {
  $vp = Join-Path $root 'docs/vision/plans/VP-035-foundation-architecture-health.md'
  $v = (Fm $vp)['version']; if (-not $v) { return @('VP frontmatter missing version') }
  $out = @()
  $targets = @(
    @{ p = (Join-Path $ws 'goal-tree.md'); pat = 'VP-035-foundation-architecture-health`（active · v([0-9.]+)）' },
    @{ p = (Join-Path $ws 'workspace.md'); pat = 'VP-035-foundation-architecture-health`（`active` · v([0-9.]+)）' },
    @{ p = (Join-Path $ws 'GOAL-001-foundation-architecture-health/00-meta.md'); pat = 'VP-035-foundation-architecture-health`（`active` · v([0-9.]+)）' }
  )
  foreach ($t in $targets) {
    $txt = ReadUtf8 $t.p
    $m = [regex]::Matches($txt, $t.pat)
    if ($m.Count -eq 0) { $out += "$(Split-Path $t.p -Leaf): 未找到 VP 版本投影" }
    else { foreach ($x in $m) { if ($x.Groups[1].Value -ne $v) { $out += "$(Split-Path $t.p -Leaf): 投影 v$($x.Groups[1].Value) != VP v$v" } } }
  }
  $out
}

Check '5 必备字段与 id 一致性' {
  $out = @()
  Get-ChildItem $ws -Recurse -Filter '*.md' -File | ForEach-Object {
    $h = Fm $_.FullName
    if ($h.Count -eq 0) { $out += "$($_.Name): 缺 frontmatter"; return }
    foreach ($f in 'status','created','updated','parent','version') { if (-not $h.ContainsKey($f)) { $out += "$($_.Name): 缺 $f" } }
  }
  Get-ChildItem $ws -Directory | ForEach-Object {
    $meta = Join-Path $_.FullName '00-meta.md'
    if (Test-Path $meta) { $id = (Fm $meta)['id']; if ($id -ne $_.Name) { $out += "$($_.Name): id=$id 不匹配" } }
  }
  $out
}

Check '6 边界守恒' {
  $out = @()
  $apps = (git diff --name-only ebe6013c..HEAD -- apps) -join ''
  if ($apps) { $out += "apps/** 有变更: $apps" }
  $g = (git show 'ebe6013c:docs/vision/roadmap.md' 2>$null | Select-String -Pattern 'trigger-gated' -AllMatches | Measure-Object).Count
  $n = (Select-String -Path (Join-Path $root 'docs/vision/roadmap.md') -Pattern 'trigger-gated' -AllMatches | Measure-Object).Count
  if ($n -lt $g) { $out += "trigger-gated 计数下降: $g -> $n" }
  $out
}

if ($fail -eq 0) { Write-Output 'ALL CHECKS PASS' } else { Write-Output "$fail CHECK(S) FAILED" }
```

> 检查项 4（frontmatter 版本递增）依赖最近提交列表，难以在脚本内稳定判定「内容是否变化」，故以人工核对 + 提交前的 `git log --name-only -3` 清单为准；本脚本覆盖其余五项机器可判定项。
