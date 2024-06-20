.PHONY: all linux darwin windows

all: linux darwin windows

linux:
	GOOS=linux GOARCH=amd64 go build -o out/codegen-linux ./main.go

macos:
	GOOS=darwin GOARCH=arm64 go build -o out/codegen-darwin ./main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o out/codegen-windows.exe ./main.go
