set GOOS=linux
set GOARCH=arm
go build  -ldflags="-s -w" cmd/iot/main.go
