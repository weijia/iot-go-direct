@echo off
REM 向 EMQX broker 发送一条 MQTT 控制消息，触发网关经(mock)串口与虚拟玻璃交互，
REM 并在 out 主题上打印网关回复，用于验证完整闭环。
REM
REM 用法:
REM   run_publish_test.bat                              -> get_glass_status_request (node1，触发串口 cmd=2 回包)
REM   run_publish_test.bat update <node_id> <color>    -> 改色(需 8 位 hex node_id 才过 schema)
REM   run_publish_test.bat get <node_id>
cd /d %~dp0
set GATEWAY=F12309150001
set METHOD=%1
if "%METHOD%"=="" set METHOD=get_glass_status_request
set NODE=%2
if "%NODE%"=="" set NODE=node1
set COLOR=%3
if "%COLOR%"=="" set COLOR=1022344052647082

echo Publishing [%METHOD%] node=%NODE% to EMQX broker...
go run ./cmd/mqtt_pub_test -gateway %GATEWAY% -method %METHOD% -node %NODE% -color %COLOR%
