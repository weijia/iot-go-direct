set GOOS=windows
set GOARCH=amd64
go build -o int_test.exe -ldflags="-s -w" cmd/int_test/main.go
