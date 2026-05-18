$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

Write-Host "=== Compile ==="
go build -o portal.exe ./cmd/portal/

Write-Host "=== Stop Old ==="
$old = Get-Process -Name "portal" -ErrorAction SilentlyContinue
if ($old) {
    Write-Host "Found PID $($old.Id), killing..."
    try {
        Stop-Process -Id $old.Id -Force -ErrorAction Stop
        Write-Host "Stopped"
    }
    catch {
        Write-Host "Access denied, trying taskkill..."
        taskkill /F /PID $old.Id /T 2>$null
        Start-Sleep 2
        $check = Get-Process -Name "portal" -ErrorAction SilentlyContinue
        if ($check) {
            Write-Host "WARN: still running, need admin"
        } else {
            Write-Host "Stopped via taskkill"
        }
    }
    Start-Sleep 1
}
else {
    Write-Host "No old process"
}

Write-Host "=== Clear State ==="
$sp = "$env:USERPROFILE\.portal\state.json"
if (Test-Path $sp) { Remove-Item $sp -Force }

Write-Host "=== Start ==="
Start-Process -FilePath ".\portal.exe" -WorkingDirectory $PSScriptRoot
Start-Sleep 2

Write-Host "=== Status ==="
try {
    $h = Invoke-RestMethod -Uri http://127.0.0.1:8747/api/health
    Write-Host "OK v$($h.version) online=$($h.services_online)/$($h.services_configured)"
}
catch {
    Write-Host "Not ready"
}
Write-Host ""
Write-Host "Portal: http://127.0.0.1:8747"
