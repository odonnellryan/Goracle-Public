GORACLE_IMAGE := goracle-build
GORACLE_STANDALONE_IMAGE := goracle-standalone
GORACLE_TEST_IMAGE := goracle-test
GORACLE_TEST_ENV := goracle-test-env

build:
	docker build --rm -t $(GORACLE_IMAGE) . 
	mkdir -p bin

binary: build
	docker run -rm -i -t -v /root/workspace/Goracle/bin:/go/src/Goracle/bin "$(GORACLE_IMAGE)" ./make.sh

images: binary
	cp -R bin images/goracle/
	cp -R bin images/goracle-standalone/
	docker build --rm -t goracle images/goracle
	docker build --rm -t $(GORACLE_STANDALONE_IMAGE) images/goracle-standalone
	
test: build
	docker build --rm -t $(GORACLE_TEST_IMAGE) images/goracle-test
	docker run -privileged -t -i -dns 8.8.8.8 -v /root/workspace/Goracle:/go/src/Goracle $(GORACLE_TEST_IMAGE) /start.sh
	
web: build
	docker build --rm -t goracle-web images/goracle-web

all: binary
