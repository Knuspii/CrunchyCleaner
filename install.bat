@echo off
:: CrunchyCleaner Install Script

set "REPO=Knuspii/CrunchyCleaner"
set "BINARY_NAME=crunchycleaner.exe"
set "INSTALL_PATH=%SystemRoot%\System32\%BINARY_NAME%"

echo Installing CrunchyCleaner...

:: Download using PowerShell (Windows-native way to curl)
echo Downloading %BINARY_NAME%...
powershell -Command "Invoke-WebRequest -Uri 'https://github.com/%REPO%/releases/latest/download/%BINARY_NAME%' -OutFile '%BINARY_NAME%'"

echo Installing to %INSTALL_PATH%...
move /Y "%BINARY_NAME%" "%INSTALL_PATH%"

echo.
echo Done! You can now use 'crunchycleaner'
echo Type: 'crunchycleaner -h'
