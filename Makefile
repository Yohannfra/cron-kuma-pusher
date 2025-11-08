NAME=cron-kuma-pusher

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME) -v

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)

.PHONY: lint
## lint: Lint code.
lint:
	@golangci-lint run


.PHONY: format
## format: Format code.
format:
	@gofmt -w -l .

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: test
## test: Run tests with verbose mode
test:
	@go test -v ./...

.PHONY: help
all: help
## help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
