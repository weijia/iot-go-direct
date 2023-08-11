set GOOS=windows
set GOARCH=amd64
go build  -ldflags="-s -w" cmd/iot/main.go
