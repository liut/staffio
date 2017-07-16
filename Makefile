.SILENT :
.PHONY : vet main clean dist release
DATE := `date '+%Y%m%d'`

NAME:=staffio
ROOF:=lcgc/platform/$(NAME)
TAG:=`git describe --tags --always`
LDFLAGS:=-X $(ROOF)/settings.buildVersion=$(TAG)-$(DATE)

main:
	echo "Building $(NAME)"
	go build -ldflags "$(LDFLAGS)" $(ROOF)/cmd/$(NAME)

all: vet main dist release

vet:
	echo "Checking ."
	go tool vet -atomic -bool -copylocks -nilfunc -printf -shadow -rangeloops -unreachable -unsafeptr -unusedresult .

clean:
	echo "Cleaning dist"
	rm -rf dist fe/build
	rm -f $(NAME) $(NAME)-*

dist: clean
	echo "Building $(NAME)"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$(NAME) $(ROOF)/cmd/$(NAME)
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$(NAME) $(ROOF)/cmd/$(NAME)

release: dist
	tar -cvJf $(NAME)-linux-amd64-$(TAG).tar.xz -C dist/linux_amd64 $(NAME)
	tar -cvJf $(NAME)-darwin-amd64-$(TAG).tar.xz -C dist/darwin_amd64 $(NAME)

fetch-exmail:
	echo "Building $@"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$(NAME)-$@ $(ROOF)/cmd/fetch-exmail
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$(NAME)-$@ $(ROOF)/cmd/fetch-exmail
.PHONY: fetch-exmail

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
	gulp clean build
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

test:
	go test

docker-build:
	docker build --rm -t $(NAME) .

