GORACLE_IMAGE := goracle-build

build:
	docker build -t $(GORACLE_IMAGE) . 
	mkdir -p bin

binary: build
	docker run -rm -i -t -v /root/workspace/Goracle/bin:/go/src/Goracle/bin "$(GORACLE_IMAGE)" ./make.sh

images: binary
	cp -R bin images/goracle/
	cp -R bin images/goracle-standalone/
	docker build -t goracle images/goracle
	docker build -t goracle-standalone images/goracle-standalone
	
test: build
	docker run -rm -i -t -v /root/workspace/Goracle:/go/src/Goracle $(GORACLE_IMAGE) go test

web: build
	docker build -t goracle-web images/goracle-web

all: binary
