@echo off
REM 启动 iot 网关：使用 EMQX 免费公共 MQTT broker，无硬件即可端到端测试控制板/MQTT 链路。
REM
REM 用法:
REM   run_emqx_test.bat          -> mock 串口(内存模拟控制板) + emqx broker（无需串口设备）
REM   run_emqx_test.bat real     -> 真实串口 COM1/9600 + emqx broker（需接设备，端口可在下方修改）
REM
REM EMQX 免费公共服务器: broker.emqx.io:1883 (匿名，无需账号密码)
REM 可用 MQTTX / mosquitto_sub 订阅主题 device/<GatewayNodeId>/in 观察上下行消息。

setlocal
cd /d %~dp0

set BROKER=broker.emqx.io:1883
set ARGS=-mock -broker %BROKER%

if /I "%1"=="real" (
    set ARGS=-port COM1 -baud 9600 -broker %BROKER%
)

echo Starting iot gateway with args: %ARGS%
echo (Ctrl+C to stop)
go run ./cmd/iot %ARGS%

endlocal
