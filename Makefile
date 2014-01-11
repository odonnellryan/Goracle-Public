GORACLE_IMAGE := goracle-build

build:
	mkdir bin & true

binary: build
	docker build -t $(GORACLE_IMAGE) .
	docker run -rm -i -t -v $(CURDIR)/bin:/go/src/Goracle/bin "$(GORACLE_IMAGE)" ./make.sh

images: build
	docker build -t goracle images/goracle
	docker build -t goracle-standalone images/goracle-standalone

all: binary
