@echo off
setlocal enabledelayedexpansion
set port=9090
for /f "tokens=1-5" %%a in ('netstat -ano ^| find ":%port%"') do (
    if "%%e%" == "" (
        set pid=%%d
    ) else (
        set pid=%%e
    )
    echo !pid!
    taskkill /f /pid !pid!
    goto out1
)
:out1
for /f "tokens=1-5" %%a in ('tasklist ^| findstr "DcAgent.exe"') do (
    if "%%b%" == "" (
        goto put
    ) else (
        set pid=%%b
    )
    echo !pid!
    taskkill /f /pid !pid!
    goto out
)
:out
%~dp0sendIP.exe %~dp0
rd/s/q "%~dp0"