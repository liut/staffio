.SILENT :
.PHONY : dep vet main clean dist package
DATE := `date '+%Y%m%d'`

WITH_ENV = env `cat .env 2>/dev/null | xargs`

NAME:=staffio
ROOF:=github.com/liut/$(NAME)
SOURCES=$(shell find client cmd pkg templates -type f \( -name "*.go" ! -name "*_test.go" \) -print )
TAG:=`git describe --tags --always`
LDFLAGS:=-X $(ROOF)/pkg/settings.buildVersion=$(TAG)-$(DATE)

main:
	echo "Building $(NAME)"
	go build -ldflags "$(LDFLAGS)" .

all: vet dist package

dep:
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow

vet:
	echo "Checking ./pkg/..."
	go vet -vettool=$(which shadow) -atomic -bool -copylocks -nilfunc -printf -rangeloops -unreachable -unsafeptr -unusedresult ./pkg/...

clean:
	echo "Cleaning dist"
	rm -rf dist fe/build
	rm -f $(NAME) $(NAME)-*

dist/linux_amd64/$(NAME): $(SOURCES)
	echo "Building $(NAME) of linux"
	mkdir -p dist/linux_amd64 && cd dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS) -s -w" $(ROOF)

dist/darwin_amd64/$(NAME): $(SOURCES)
	echo "Building $(NAME) of darwin"
	mkdir -p dist/darwin_amd64 && cd dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS) -w" $(ROOF)

dist/windows_amd64/$(NAME): $(SOURCES)
	echo "Building $(NAME) of windows"
	mkdir -p dist/windows_amd64 && cd dist/windows_amd64 && GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS) -s -w" $(ROOF)

dist: vet dist/linux_amd64/$(NAME) dist/darwin_amd64/$(NAME) dist/windows_amd64/$(NAME)


package: dist
	tar -cvJf $(NAME)-linux-amd64-$(TAG).tar.xz -C dist/linux_amd64 $(NAME)
	tar -cvJf $(NAME)-darwin-amd64-$(TAG).tar.xz -C dist/darwin_amd64 $(NAME)
	tar -cvJf $(NAME)-templates-$(TAG).tar.xz templates

fetch-exmail:
	echo "Building $@"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$(NAME)-$@ $(ROOF)/cmd/$@
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$(NAME)-$@ $(ROOF)/cmd/$@
.PHONY: fetch-exmail

wechat-work:
	echo "Building $@"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$(NAME)-$@ $(ROOF)/cmd/$@
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$(NAME)-$@ $(ROOF)/cmd/$@
.PHONY: wechat-work

demo:
	echo "Building $@"
	go build -ldflags "$(LDFLAGS)" $(ROOF)/cmd/$(NAME)-$@
.PHONY: demo

gen-key:
	echo "Building $@"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$(NAME)-$@ $(ROOF)/cmd/gen-key
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$(NAME)-$@ $(ROOF)/cmd/gen-key
.PHONY: gen-key

js-deps:
	npm install
.PHONY: js-deps

js-build:
	./node_modules/.bin/gulp clean build
.PHONY: js-build

statik:
	statik -src htdocs -dest ./pkg/web
.PHONY: statik

gofmt:
	if [ -n "$(shell gofmt -l .)" ]; then \
		echo 1>&2 'The following files need to be formatted:'; \
		gofmt -l .; \
		exit 1; \
	fi

test-ldap: vet
	mkdir -p tests
	@DEBUG=staffio:ldap go test -v -cover -coverprofile tests/cover_ldap.out ./pkg/backends/ldap
	@go tool cover -html=tests/cover_ldap.out -o tests/cover_ldap.out.html


docker-db-build:
	docker build --rm -t liut7/$(NAME)-db:$(TAG) database/
	docker tag liut7/$(NAME)-db:$(TAG) liut7/$(NAME)-db:latest

docker-auto-build:
	docker build --rm -t $(NAME) .

docker-local-build: dist/linux_amd64/$(NAME)
	echo "Building docker image"
	cp -rf templates entrypoint.sh dist/
	cp -rf Dockerfile.local dist/Dockerfile
	docker build --rm -t liut7/$(NAME):$(TAG) dist/
	docker tag liut7/$(NAME):$(TAG) liut7/$(NAME):latest
.PHONY: $@

