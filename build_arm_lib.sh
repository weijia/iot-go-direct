#!/bin/bash
export CC=arm-linux-gnueabihf-gcc
arm-linux-gnueabihf-gcc -c -o pan3028.o pan3028-app\pan3028.c 
arm-linux-gnueabihf-gcc -c -o pan3028_port.o pan3028-app\pan3028_port.c
arm-linux-gnueabihf-gcc -c -o spi-3028.o pan3028-app\spi-3028.c 
arm-linux-gnueabihf-gcc -c -o radio.o pan3028-app\radio.c
rem .a must be in the same folder as the go file really using the .a lib
arm-linux-gnueabihf-ar rcs pkg\lora\libpan3028.a pan3028.o pan3028_port.o spi-3028.o radio.o
