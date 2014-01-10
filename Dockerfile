FROM		stackbrew/ubuntu:12.04
MAINTAINER	Ryan


RUN        apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
        curl \
        --no-install-recommends

# Install Go
RUN        curl -s https://go.googlecode.com/files/go1.2.src.tar.gz | tar -v -C /usr/local -xz
ENV        PATH        /usr/local/go/bin:$PATH
ENV        GOPATH        /go:/go/src/github.com/dotcloud/docker/vendor
RUN        cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Install dependencies
RUN go get "labix.org/v2/mgo"

WORKDIR		/go/src/Goracle

# Upload source
ADD	.	/go/src/Goracle

ENTRYPOINT ["make binary"]