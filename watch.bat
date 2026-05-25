@echo off
setlocal enabledelayedexpansion

:loop
:: Start the Go application in the background
start /b "" go run main.go

:: Store the current timestamp using a temporary file
echo. > %temp%\last_run.txt

:watch_loop
:: Wait 2 seconds between checks
timeout /t 2 /nobreak >nul

:: Check if any .go files have changed since the timestamp file
xcopy . %temp% /d /e /y /l /m | findstr /i "\.go$" >nul
if errorlevel 1 (
    goto watch_loop
)

:: If changes are found, kill the running go processes and restart
echo Change detected. Restarting...
taskkill /f /im main.exe >nul 2>&1
taskkill /f /im go.exe >nul 2>&1
goto loop
