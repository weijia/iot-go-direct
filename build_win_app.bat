set GOOS=windows
set GOARCH=amd64
for /f "delims=" %%i in ('git rev-parse HEAD') do set LATEST_COMMIT=%%i
echo The latest commit hash is: %LATEST_COMMIT%

python -c "from datetime import datetime; print(datetime.now().strftime('%%Y-%%m-%%d-%%H:%%M:%%S'))" > temp.txt
for /f "delims=" %%i in (temp.txt) do set CURRENT_DATE_TIME=%%i  
echo Current date is %CURRENT_DATE_TIME%
del temp.txt

python increase_version.py

for /f "delims=" %%i in (version.txt) do set CURRENT_VERSION=%%i  
echo Current version is %CURRENT_VERSION%


@REM go build -a -v  -ldflags="-s -w  -X main.GitCommit=%LATEST_COMMIT% -X main.BuildDate=%DATE:~0,10%-%TIME%" ./cmd/cobra_main/main.go ./cmd/cobra_main/main_loop.go ./cmd/cobra_main/version_and_supervisord.go ./cmd/cobra_main/lora_service.go ./cmd/cobra_main/version.go
go build -a -v  -ldflags="-s -w -X bsp.SwVersion=%CURRENT_VERSION% -X main.GitCommit=%LATEST_COMMIT% -X main.BuildDate=%CURRENT_DATE_TIME%" ./cmd/cobra_main/main.go ./cmd/cobra_main/main_loop.go ./cmd/cobra_main/version_and_supervisord.go ./cmd/cobra_main/lora_service.go ./cmd/cobra_main/version.go

@REM go build  -ldflags="-s -w" cmd/iot/main.go
@REM go build  -ldflags="-s -w" cmd/int_test/int_test_main.go
@REM go build  -ldflags="-s -w" cmd/lora_test/lora_service_test_main.go
@REM go build  -ldflags="-s -w" cmd/lora_service/lora_service.go

