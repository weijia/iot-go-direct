rem PATH=%PATH%;D:\迅雷下载\gcc-linaro-7.5.0-2019.12-i686-mingw32_arm-linux-gnueabihf.tar\gcc-linaro-7.5.0-2019.12-i686-mingw32_arm-linux-gnueabihf\bin\
PATH=%PATH%;D:\codes\gcc-linaro-7.5.0-2019.12-i686-mingw32_arm-linux-gnueabihf\bin\
@REM copy version\version-tpl.txt version\version.go
@REM echo var BuildDate = "%TIME%" >> version\version.go

set CGO_ENABLED=1
set GOARCH=arm
set GOOS=linux
set CGO_LDFLAGS="-static"
set CC=arm-linux-gnueabihf-gcc
arm-linux-gnueabihf-gcc -c -o pan3028.o pan3028-app\pan3028.c 
arm-linux-gnueabihf-gcc -c -o pan3028_port.o pan3028-app\pan3028_port.c
arm-linux-gnueabihf-gcc -c -o spi-3028.o pan3028-app\spi-3028.c 
arm-linux-gnueabihf-gcc -c -o radio.o pan3028-app\radio.c
rem .a must be in the same folder as the go file really using the .a lib
arm-linux-gnueabihf-ar rcs pkg\lora\libpan3028.a pan3028.o pan3028_port.o spi-3028.o radio.o
rem go build -a -v

for /f "delims=" %%i in ('git rev-parse HEAD') do set LATEST_COMMIT=%%i
echo The latest commit hash is: %LATEST_COMMIT%

go build -a -v  -ldflags="-s -w  -X main.GitCommit=%LATEST_COMMIT% -X main.BuildDate=%DATE:~0,10%-%TIME%" ./cmd/cobra_main/main.go ./cmd/cobra_main/main_loop.go ./cmd/cobra_main/version_and_supervisord.go ./cmd/cobra_main/lora_service.go ./cmd/cobra_main/version.go
@REM go build -a -v  -ldflags="-s -w" cmd/iot/main.go
@REM go build -a -v  -ldflags="-s -w" cmd/lora_service/lora_service.go
move main iot_go
sync.bat
