generate:
	@echo Embedding assets
	@go generate

install: generate
	@echo Installing the application
	@go install -ldflags "-s -w"

all: install
