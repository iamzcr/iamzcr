---
name: dev-stop
description: Stop all dev servers started by dev-start
---

# Dev Stop Skill

## Steps

Execute this single PowerShell script:

```powershell
$pidFile = "$PWD\.opencode\running.pids"
if (Test-Path -LiteralPath $pidFile) {
  Get-Content -LiteralPath $pidFile | ForEach-Object {
    $id = $_.Trim()
    if ($id) {
      taskkill /PID $id /T /F 2>&1 | Out-Null
      Write-Host "Killed process tree for PID: $id"
    }
  }
  Remove-Item -LiteralPath $pidFile -Force
  Write-Host "`nAll dev servers stopped."
} else {
  Write-Host "No running.pids file found. Trying fallback cleanup..."

  # Kill Go binaries
  Get-Process -Name "admin","frontend" -ErrorAction SilentlyContinue | Stop-Process -Force

  # Kill processes on dev ports
  $ports = @(8081, 8082, 3000, 3001)
  foreach ($port in $ports) {
    $conn = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
    if ($conn) {
      $conn | ForEach-Object {
        Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue
      }
    }
  }
  Write-Host "Fallback cleanup complete."
}
```

## Notes

- Uses `taskkill /T` to kill the entire process tree (closes the PowerShell window + all child processes)
- If no PIDs file exists, falls back to killing processes by name and port
