TEST ?= $(shell go list ./... | grep -v vendor)
VERSION = $(shell grep 'const Version' version.go | sed -e 's/.*= //g' -e 's/"//g')
REVISION = $(shell git describe --always)

INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

default: build
ci: depsdev test vet lint ## Run test and more...
depsdev: ## Installing dependencies for development
	go get -u github.com/tcnksm/ghr
	go get github.com/mitchellh/gox

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	go test -v $(TEST) -timeout=30s -parallel=4
	go test -race $(TEST)

vet: ## Exec go vet
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Vetting$(RESET)"
	go vet $(TEST)

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -set_exit_status $(TEST)

ghr: ## Upload to Github releases without token check
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Releasing for Github$(RESET)"
	ghr -u pyama86 v$(VERSION)-$(REVISION) pkg

build: ## Build as linux binary
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Building$(RESET)"
	./misc/build $(VERSION) $(REVISION)

dist: build ## Upload to Github releases
	@test -z $(GITHUB_TOKEN) || $(MAKE) ghr

dev: build
	docker build -t gothree .
	docker run --privileged --rm -e AWS_ACCESS_KEY_ID=$$AWS_ACCESS_KEY_ID -e AWS_REGION=$$AWS_REGION -e AWS_BUCKET=$$AWS_BUCKET -e AWS_SECRET_ACCESS_KEY=$$AWS_SECRET_ACCESS_KEY -it gothree
