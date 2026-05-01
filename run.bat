@echo off
if "%~1"=="" (
    echo Usage: run.bat ^<PORT^>
    exit /b 1
)

set PORT=%1
docker compose up -d --build
