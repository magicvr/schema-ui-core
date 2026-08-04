[CmdletBinding()]
param([Parameter(ValueFromRemainingArguments = $true)][string[]]$UpdateArgs)

$ErrorActionPreference = 'Stop'
$Updater = Join-Path $PSScriptRoot 'update.py'
$Python = Get-Command python -ErrorAction SilentlyContinue
if (-not $Python) { $Python = Get-Command python3 -ErrorAction SilentlyContinue }
if (-not $Python) {
    Write-Error 'Goal Governance updater requires Python 3.'
    exit 1
}
& $Python.Source $Updater @UpdateArgs
exit $LASTEXITCODE
