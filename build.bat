set GOOS=windows
set GOARCH=amd64
go build  -ldflags="-s -w" cmd/iot/main.go
go build  -ldflags="-s -w" cmd/int_test/int_test_main.go
@REM go build  -ldflags="-s -w" cmd/lora_test/lora_service_test_main.go
@REM go build  -ldflags="-s -w" cmd/lora_service/lora_service.go

