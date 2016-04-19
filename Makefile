.SILENT :
.PHONY : main clean gofmt dist-tight
DATE := `date '+%Y%m%d'`

NAME:=staffio
ROOF:=lcgc/platform/$(NAME)
TAG:=`git describe --tags --always`
LDFLAGS:=-X $(ROOF)/settings.buildVersion=$(TAG)-$(DATE)

main:
	echo "Building $(NAME)"
	go build -ldflags "$(LDFLAGS)"

all: main dist dist-tight release

clean:
	rm -rf dist dist-tight
	rm -f $(NAME) $(NAME)-*.?z

dist: clean
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$(NAME)
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$(NAME)

release: dist
	tar -cvJf $(NAME)-linux-amd64-$(TAG).tar.xz -C dist/linux_amd64 $(NAME)
	tar -cvJf $(NAME)-darwin-amd64-$(TAG).tar.xz -C dist/darwin_amd64 $(NAME)

release-clean:
	rm -f *.tar.xz

get-deps:
	go get github.com/robfig/glock
	glock sync -n < GLOCKFILE

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

dist-tight:
	echo "Building tight version"
	rm -rf dist-tight
	mkdir -p dist-tight/linux_amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o dist-tight/linux_amd64/$(NAME) -a -installsuffix nocgo ./staffio-tight

docker-build-tight:
	strip dist-tight/linux/amd64/$(NAME)
	docker build --rm -t lcgc/$(NAME):tight -f staffio-tight/Dockerfile .

