FROM		stackbrew/ubuntu:12.04
MAINTAINER	Ryan


RUN     apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
        curl \
        build-essential \
	bzr \
	--no-install-recommends

# Install Go
RUN        curl -s https://go.googlecode.com/files/go1.2.src.tar.gz | tar -v -C /usr/local -xz
ENV        PATH        /usr/local/go/bin:$PATH
ENV        GOPATH        /go
RUN        cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Install dependencies
RUN go get "labix.org/v2/mgo"

WORKDIR		/go/src/Goracle

RUN     DEBIAN_FRONTEND=noninteractive apt-get install -yq \
        make \
        --no-install-recommends


# Upload source
ADD	.	/go/src/Goracle
