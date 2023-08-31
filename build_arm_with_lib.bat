rem PATH=%PATH%;D:\迅雷下载\gcc-linaro-7.5.0-2019.12-i686-mingw32_arm-linux-gnueabihf.tar\gcc-linaro-7.5.0-2019.12-i686-mingw32_arm-linux-gnueabihf\bin\
set CGO_ENABLED=1
set GOARCH=arm
set GOOS=linux
set CGO_LDFLAGS="-static"
set CC=arm-linux-gnueabihf-gcc
arm-linux-gnueabihf-gcc -c -o pan3028.o pan3028-app\pan3028.c 
arm-linux-gnueabihf-gcc -c -o pan3028_port.o pan3028-app\pan3028_port.c
arm-linux-gnueabihf-gcc -c -o spi-3028.o pan3028-app\spi-3028.c 
arm-linux-gnueabihf-gcc -c -o radio.o pan3028-app\radio.c 
arm-linux-gnueabihf-ar rcs libpan3028.a pan3028.o pan3028_port.o spi-3028.o radio.o
rem go build -a -v

go build -a -v  -ldflags="-s -w" cmd/iot/main.go
