GORACLE_IMAGE := goracle-build
GORACLE_STANDALONE_IMAGE := goracle-standalone
GORACLE_TEST_IMAGE := goracle-test

build:
	docker build -t $(GORACLE_IMAGE) . 
	mkdir -p bin

binary: build
	docker run -rm -i -t -v /root/workspace/Goracle/bin:/go/src/Goracle/bin "$(GORACLE_IMAGE)" ./make.sh

images: binary
	cp -R bin images/goracle/
	cp -R bin images/goracle-standalone/
	docker build -t goracle images/goracle
	docker build -t $(GORACLE_STANDALONE_IMAGE) images/goracle-standalone
	
test: build
	docker build -t $(GORACLE_TEST_IMAGE) images/goracle-test
	docker run -privileged -i -t -v /root/workspace/Goracle:/go/src/Goracle $(GORACLE_TEST_IMAGE)

web: build
	docker build -t goracle-web images/goracle-web

all: binary
