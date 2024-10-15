generate:
	@echo Embedding assets
	@go generate

install:
	@echo Installing the application
	@go install -ldflags "-s -w"

all: generate install
