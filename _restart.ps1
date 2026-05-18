Set-StrictMode -Off
Set-Location d:\workspaces\portal
go build -o portal.exe ./cmd/portal/
$old = Get-Process -Name portal -ErrorAction SilentlyContinue
if ($old) { Stop-Process -Id $old.Id -Force -ErrorAction SilentlyContinue; Start-Sleep 1 }
Start-Process -FilePath .\portal.exe -WorkingDirectory d:\workspaces\portal
Start-Sleep 2
try { $h = Invoke-RestMethod -Uri http://127.0.0.1:8747/api/health; Write-Host "OK v$($h.version) online=$($h.services_online)/$($h.services_configured)" } catch { Write-Host "NOT READY" }
