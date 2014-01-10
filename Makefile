GORACLE_IMAGE := goracle-build

build:
	mkdir bin & true
	docker build -t $(GORACLE_IMAGE) .

binary: build
	docker run -rm -i -t -v $(CURDIR)/bin:/go/src/Goracle/bin "$(GORACLE_IMAGE)" build/make.sh
