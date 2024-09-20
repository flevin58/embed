generate:
	@echo Embedding assets
	@go generate

install: generate
	@echo Installing the application
	@go install

all: install
