#!/bin/bash

# 构建 Linux 64 位的可执行文件
GOOS=linux GOARCH=amd64 go build -o codegen-linux ./main.go

# 构建 macOS 64 位的可执行文件
GOOS=darwin GOARCH=amd64 go build -o codegen-darwin ./main.go

# 构建 Windows 64 位的可执行文件
GOOS=windows GOARCH=amd64 go build -o codegen-windows.exe ./main.go
