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
	docker build -t $(GORACLE_TEST_IMAGE) images/goracle-test
	 docker run -privileged -d -v /root/workspace/Goracle:/go/src/Goracle -dns 8.8.8.8 $(GORACLE_TEST_IMAGE) /pull.sh >> docker commit $(GORACLE_TEST_IMAGE)
	
test: build
	docker run -privileged -t -i -v /root/workspace/Goracle:/go/src/Goracle -dns 8.8.8.8 $(GORACLE_TEST_IMAGE) /start.sh
web: build
	docker build -t goracle-web images/goracle-web

all: binary
