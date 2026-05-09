---
name: dev-start
description: Start all project dev servers (Go backends + Vue frontends) and open browser tabs
---

# Dev Start Skill

## Services

| Service | Port | Command | Dir |
|---|---|---|---|
| Admin API | 8081 | `air -c .air_admin.toml` | project root |
| Frontend API | 8082 | `air -c .air_frontend.toml` | project root |
| Admin SPA | 3001 | `npm run dev` | `web/admin` |
| Blog SPA | 3000 | `npm run dev` | `web/frontend` |

## Steps

Execute this single PowerShell script (everything in ONE bash call — do NOT split across multiple calls):

```powershell
$pids = @()

$proc = Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PWD'; air -c .air_admin.toml" -WindowStyle Minimized -PassThru
$pids += $proc.Id
Write-Host "Admin API (8081) PID: $($proc.Id)"

$proc = Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PWD'; air -c .air_frontend.toml" -WindowStyle Minimized -PassThru
$pids += $proc.Id
Write-Host "Frontend API (8082) PID: $($proc.Id)"

$proc = Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PWD\web\admin'; npm run dev" -WindowStyle Minimized -PassThru
$pids += $proc.Id
Write-Host "Admin SPA (3001) PID: $($proc.Id)"

$proc = Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PWD\web\frontend'; npm run dev" -WindowStyle Minimized -PassThru
$pids += $proc.Id
Write-Host "Blog SPA (3000) PID: $($proc.Id)"

$pids | Out-File -LiteralPath "$PWD\.opencode\running.pids" -Encoding utf8
Write-Host "`nAll 4 dev servers started. PIDs saved to .opencode\running.pids"

Start-Sleep -Seconds 3
Start-Process "http://localhost:3000"
Start-Process "http://localhost:3001"
Write-Host "Browser opened: http://localhost:3000 (Blog) | http://localhost:3001 (Admin)"
```

## Notes

- Each service runs in its own minimized PowerShell window — close them individually or use `dev-stop`
- PIDs are saved to `.opencode\running.pids` for dev-stop
- Admin login: `test` / `admin123` at `http://localhost:3001/login`
