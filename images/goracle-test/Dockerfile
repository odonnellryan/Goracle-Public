FROM		goracle-standalone
MAINTAINER	Ryan

RUN         apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
            curl \
            build-essential \
	        bzr \
	        git \
	        --no-install-recommends

# Install Go
RUN         curl -s https://go.googlecode.com/files/go1.2.src.tar.gz | tar -v -C /usr/local -xz
ENV         PATH        /usr/local/go/bin:$PATH
ENV         GOPATH        /go
RUN         cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Install go dependencies
RUN         go get "labix.org/v2/mgo"
RUN         go get "labix.org/v2/mgo/bson"
RUN         apt-get install mercurial -yq
RUN         go get "code.google.com/p/go.tools/cmd/cover"
RUN         go get "github.com/ziutek/mymysql/thrsafe"
RUN         go get "github.com/ziutek/mymysql/autorc"
RUN         go get "github.com/ziutek/mymysql/godrv"

WORKDIR		/go/src/Goracle

RUN         DEBIAN_FRONTEND=noninteractive apt-get install -yq \
            make \
            --no-install-recommends

# Upload source
ADD	        .	/go/src/Goracle
ADD         scripts /
ADD         nginx /etc/nginx/
RUN         chmod +x /start.sh
RUN         chmod +x /pull.sh
