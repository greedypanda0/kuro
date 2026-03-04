.PHONY: dev build run clean

dev:
	@echo Building development environment...
	@air

build:
	@echo Building production environment...
	@if not exist bin mkdir bin
	@go build -o bin/remote.exe cmd/remote/main.go

run: build
	@echo Running production environment...
	@bin\remote.exe

clean:
	@echo Cleaning...
	@if exist bin rmdir /s /q bin