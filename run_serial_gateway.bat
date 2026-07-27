@echo off
chcp 65001 >nul
REM Start the MQTT->serial gateway in local debug mode.
REM Connects to broker.emqx.io with the prefilled account (needs auth).
REM serial_backend=simulate acts as the glass board and replies per protocol, so no real serial port is needed.
REM No web server required: just open cmd\mqtt_serial_gateway\web\index.html in a browser (file:// works).
cd /d %~dp0
cd cmd\mqtt_serial_gateway
go run . mqtt_serial_gateway-local.json
pause
