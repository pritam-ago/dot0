# PowerShell build script for the relay project
# Usage: .\build.ps1 [command]

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Available commands:" -ForegroundColor Green
    Write-Host "  install    - Install all dependencies"
    Write-Host "  build      - Build all applications"
    Write-Host "  up         - Start all services locally"
    Write-Host "  dev        - Start development servers"
    Write-Host "  down       - Stop all services"
    Write-Host "  logs       - Show logs for all services"
    Write-Host "  clean      - Clean build artifacts"
    Write-Host "  restart    - Restart all services"
    Write-Host "  rebuild    - Rebuild and start all services"
    Write-Host "  clean-all  - Remove all build artifacts and dependencies"
    Write-Host "  clean-go   - Clean Go cache"
    Write-Host "  build-go   - Build Go Relay only"
    Write-Host "  build-web  - Build Web Client only"
    Write-Host "  build-pc   - Build PC App only"
    Write-Host "  db-setup   - Database setup instructions"
    Write-Host "  health     - Check service health"
    Write-Host "  ports      - Show service ports"
}

function Install-Dependencies {
    Write-Host "Installing dependencies..." -ForegroundColor Yellow
    Set-Location "web-client"
    npm install
    Set-Location "../pc-app"
    npm install
    Set-Location "../go-relay"
    go mod download
    Set-Location ".."
    Write-Host "Dependencies installed successfully!" -ForegroundColor Green
}

function Build-All {
    Write-Host "Building all applications..." -ForegroundColor Yellow
    
    Write-Host "Building Go Relay..." -ForegroundColor Cyan
    Set-Location "go-relay"
    go build -o main.exe ./cmd/server
    Set-Location ".."
    
    Write-Host "Building Web Client..." -ForegroundColor Cyan
    Set-Location "web-client"
    npm run build
    Set-Location ".."
    
    Write-Host "Building PC App..." -ForegroundColor Cyan
    Set-Location "pc-app"
    npm run build
    Set-Location ".."
    
    Write-Host "All applications built successfully!" -ForegroundColor Green
}

function Start-Services {
    Write-Host "Starting all services..." -ForegroundColor Yellow
    
    Write-Host "Starting Go Relay server..." -ForegroundColor Cyan
    Start-Process -FilePath "powershell" -ArgumentList "-Command", "cd go-relay/cmd/server; fresh" -WindowStyle Minimized
    
    Write-Host "Starting Web Client..." -ForegroundColor Cyan
    Start-Process -FilePath "powershell" -ArgumentList "-Command", "cd web-client; npm start" -WindowStyle Minimized
    
    Write-Host "Starting PC App..." -ForegroundColor Cyan
    Start-Process -FilePath "powershell" -ArgumentList "-Command", "cd pc-app; npm run tauri dev" -WindowStyle Minimized
    
    Write-Host "All services started!" -ForegroundColor Green
    Write-Host "Check the opened PowerShell windows for logs." -ForegroundColor Yellow
}

function Stop-Services {
    Write-Host "Stopping all services..." -ForegroundColor Yellow
    
    # Stop processes by port (more reliable)
    try {
        # Stop process on port 8080 (Go Relay)
        $process8080 = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess
        if ($process8080) {
            Get-Process -Id $process8080 -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
            Write-Host "Stopped Go Relay (port 8080)" -ForegroundColor Green
        }
    } catch {}
    
    try {
        # Stop process on port 3001 (Web Client)
        $process3001 = Get-NetTCPConnection -LocalPort 3001 -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess
        if ($process3001) {
            Get-Process -Id $process3001 -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
            Write-Host "Stopped Web Client (port 3001)" -ForegroundColor Green
        }
    } catch {}
    
    try {
        # Stop process on port 3002 (PC App)
        $process3002 = Get-NetTCPConnection -LocalPort 3002 -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess
        if ($process3002) {
            Get-Process -Id $process3002 -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
            Write-Host "Stopped PC App (port 3002)" -ForegroundColor Green
        }
    } catch {}
    
    # Also stop common process names as backup
    $processesToStop = @("go", "main", "fresh", "node", "npm", "tauri", "cargo")
    foreach ($procName in $processesToStop) {
        Get-Process -Name $procName -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
    }
    
    Write-Host "All services stopped!" -ForegroundColor Green
}

function Start-Dev {
    Write-Host "Starting development servers..." -ForegroundColor Yellow
    Write-Host "Go Relay: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "Web Client: http://localhost:3001" -ForegroundColor Cyan
    Write-Host "PC App: http://localhost:3002" -ForegroundColor Cyan
    Start-Services
}

function Clean-Build {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow
    
    if (Test-Path "go-relay/main.exe") {
        Remove-Item "go-relay/main.exe" -Force
    }
    if (Test-Path "web-client/build") {
        Remove-Item "web-client/build" -Recurse -Force
    }
    if (Test-Path "pc-app/dist") {
        Remove-Item "pc-app/dist" -Recurse -Force
    }
    
    Write-Host "Build artifacts cleaned!" -ForegroundColor Green
    Write-Host "Note: To remove node_modules, run '.\build.ps1 clean-all'" -ForegroundColor Yellow
}

function Clean-All {
    Write-Host "Removing all build artifacts and dependencies..." -ForegroundColor Yellow
    
    Clean-Build
    
    if (Test-Path "web-client/node_modules") {
        Remove-Item "web-client/node_modules" -Recurse -Force
    }
    if (Test-Path "pc-app/node_modules") {
        Remove-Item "pc-app/node_modules" -Recurse -Force
    }
    
    Set-Location "go-relay"
    go clean -cache -modcache
    Set-Location ".."
    
    Write-Host "All artifacts and dependencies removed!" -ForegroundColor Green
}

function Clean-Go {
    Write-Host "Cleaning Go cache..." -ForegroundColor Yellow
    Set-Location "go-relay"
    go clean -cache -modcache
    Set-Location ".."
    Write-Host "Go cache cleaned!" -ForegroundColor Green
}

function Build-Go {
    Write-Host "Building Go Relay..." -ForegroundColor Yellow
    Set-Location "go-relay"
    go build -o main.exe ./cmd/server
    Set-Location ".."
    Write-Host "Go Relay built successfully!" -ForegroundColor Green
}

function Build-Web {
    Write-Host "Building Web Client..." -ForegroundColor Yellow
    Set-Location "web-client"
    npm run build
    Set-Location ".."
    Write-Host "Web Client built successfully!" -ForegroundColor Green
}

function Build-PC {
    Write-Host "Building PC App..." -ForegroundColor Yellow
    Set-Location "pc-app"
    npm run build
    Set-Location ".."
    Write-Host "PC App built successfully!" -ForegroundColor Green
}

function Show-DbSetup {
    Write-Host "Database Setup Instructions:" -ForegroundColor Yellow
    Write-Host "1. Install PostgreSQL from https://www.postgresql.org/download/windows/" -ForegroundColor White
    Write-Host "2. Create database: dot0" -ForegroundColor White
    Write-Host "3. Update connection string in go-relay/internal/db/database.go" -ForegroundColor White
    Write-Host "4. Default connection: postgresql://postgres:dhanus17@localhost:5432/dot0" -ForegroundColor White
}

function Show-Health {
    Write-Host "Checking service health..." -ForegroundColor Yellow
    Write-Host "Go Relay: http://localhost:8080/health" -ForegroundColor Cyan
    Write-Host "Web Client: http://localhost:3001" -ForegroundColor Cyan
    Write-Host "PC App: http://localhost:3002" -ForegroundColor Cyan
}

function Show-Ports {
    Write-Host "Service ports:" -ForegroundColor Yellow
    Write-Host "  Go Relay Server: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "  Web Client: http://localhost:3001" -ForegroundColor Cyan
    Write-Host "  PC App: http://localhost:3002" -ForegroundColor Cyan
    Write-Host "  PostgreSQL: localhost:5432" -ForegroundColor Cyan
}

function Restart-Services {
    Write-Host "Restarting all services..." -ForegroundColor Yellow
    Stop-Services
    Start-Sleep -Seconds 2
    Start-Services
}

function Rebuild-All {
    Write-Host "Rebuilding and starting all services..." -ForegroundColor Yellow
    Clean-Build
    Build-All
    Start-Services
}

# Main command dispatcher
switch ($Command.ToLower()) {
    "help" { Show-Help }
    "install" { Install-Dependencies }
    "build" { Build-All }
    "up" { Start-Services }
    "dev" { Start-Dev }
    "down" { Stop-Services }
    "clean" { Clean-Build }
    "clean-all" { Clean-All }
    "clean-go" { Clean-Go }
    "build-go" { Build-Go }
    "build-web" { Build-Web }
    "build-pc" { Build-PC }
    "db-setup" { Show-DbSetup }
    "health" { Show-Health }
    "ports" { Show-Ports }
    "restart" { Restart-Services }
    "rebuild" { Rebuild-All }
    "logs" { Write-Host "Logs are displayed in the PowerShell windows where services are running" -ForegroundColor Yellow }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Show-Help
    }
} 