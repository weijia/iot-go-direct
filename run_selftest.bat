@echo off
REM 不连 MQTT broker，直接演示网关 <-> 虚拟玻璃(串口)回包闭环：
REM   -mock       用内存 MockBoard 扮演控制板
REM   -selftest    启动后对默认节点发串口查询，打印虚拟玻璃回包与节点状态回填
REM 适合确认"串口回复"链路是否完整，无需任何外部服务。
cd /d %~dp0
go run ./cmd/iot -mock -selftest
