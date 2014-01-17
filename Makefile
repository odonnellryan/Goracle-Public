GORACLE_IMAGE := goracle-build

build:
	docker build -t $(GORACLE_IMAGE) . 
	mkdir bin & true

binary: build
	docker run -rm -i -t -v $(CURDIR)/bin:/go/src/Goracle/bin "$(GORACLE_IMAGE)" ./make.sh

images: binary
	cp -R bin images/goracle/
	cp -R bin images/goracle-standalone/
	docker build -t goracle images/goracle
	docker build -t goracle-standalone images/goracle-standalone

all: binary
