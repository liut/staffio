.SILENT :
.PHONY : staffio clean fmt

TAG:=`git describe --tags`
LDFLAGS:=-X main.buildVersion $(TAG)

all: staffio

staffio:
	echo "Building staffio"
	go build -ldflags "$(LDFLAGS)"

dist-clean:
	rm -rf dist dist-tight
	rm -f staffio-linux-*.tar.gz
	rm -f staffio-darwin-*.tar.gz

dist: dist-clean
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux/amd64/staffio
	# mkdir -p dist/linux/i386  && GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -o dist/linux/i386/staffio
	# mkdir -p dist/linux/armel  && GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(LDFLAGS)" -o dist/linux/armel/staffio
	# mkdir -p dist/linux/armhf  && GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "$(LDFLAGS)" -o dist/linux/armhf/staffio
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/staffio
	# mkdir -p dist/darwin/i386  && GOOS=darwin GOARCH=386 go build -ldflags "$(LDFLAGS)" -o dist/darwin/i386/staffio


release: dist
	# glock sync -n < GLOCKFILE
	tar -cvzf staffio-linux-amd64-$(TAG).tar.gz -C dist/linux/amd64 staffio
	# tar -cvzf staffio-linux-i386-$(TAG).tar.gz -C dist/linux/i386 staffio
	# tar -cvzf staffio-linux-armel-$(TAG).tar.gz -C dist/linux/armel staffio
	# tar -cvzf staffio-linux-armhf-$(TAG).tar.gz -C dist/linux/armhf staffio
	tar -cvzf staffio-darwin-amd64-$(TAG).tar.gz -C dist/darwin/amd64 staffio
	# tar -cvzf staffio-darwin-i386-$(TAG).tar.gz -C dist/darwin/i386 staffio

get-deps:
	go get github.com/robfig/glock
	glock sync -n < GLOCKFILE

check-gofmt:
	if [ -n "$(shell gofmt -l .)" ]; then \
		echo 1>&2 'The following files need to be formatted:'; \
		gofmt -l .; \
		exit 1; \
	fi

test:
	go test

docker-build:
	docker build --rm -t staffio .

docker-build-tight:
	mkdir -p dist-tight/linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist-tight/linux/amd64/staffio ./staffio-tight
	docker build --rm -t staffio:tight -f staffio-tight/Dockerfile .

