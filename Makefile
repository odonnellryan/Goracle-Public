GORACLE_IMAGE := goracle-build
GORACLE_STANDALONE_IMAGE := goracle-standalone
GORACLE_TEST_IMAGE := goracle-test
GORACLE_TEST_ENV := goracle-test-env

all: binary

binary: build
	docker run -rm -i -t -v /root/workspace/Goracle/bin:/go/src/Goracle/bin "$(GORACLE_IMAGE)" ./make.sh
	
build: all
	docker build --rm -t $(GORACLE_IMAGE) . 
	mkdir -p bin

images: binary
	cp -R bin images/goracle/
	cp -R bin images/goracle-standalone/
	docker build --rm -t goracle images/goracle
	docker build --rm -t $(GORACLE_STANDALONE_IMAGE) images/goracle-standalone
	
test-full: build
	docker build --rm -t $(GORACLE_TEST_IMAGE) images/goracle-test
	docker run -privileged -t -i -dns 8.8.8.8 -v /root/workspace/Goracle:/go/src/Goracle $(GORACLE_TEST_IMAGE) /start.sh
	
test: all
	docker run -privileged -t -i -dns 8.8.8.8 -v /root/workspace/Goracle:/go/src/Goracle $(GORACLE_TEST_IMAGE) /start.sh
	
web: build
	docker build --rm -t goracle-web images/goracle-web

web-test: all
	docker run -privileged -t -i -dns 8.8.8.8 goracle-web /bin/bash


