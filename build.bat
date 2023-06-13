@echo off

echo Building Windows version
set GOARCH=amd64
set GOOS=windows
go build -o build/tic-80-bundler-win-64.exe

echo Building Mac version
set GOOS=darwin
set GOARCH=amd64
go build -o build/tic-80-bundler-mac-64.exe

echo Building Linux version
set GOOS=linux
set GOARCH=amd64
go build -o build/tic-80-bundler-linux-64.exe

Rem Should return to original state, but it's late and I'm tired
set GOARCH=amd64
set GOOS=windows
