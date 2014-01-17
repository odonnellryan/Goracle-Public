FROM		ubuntu
MAINTAINER	Luma

RUN     apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
        git \
        openssh-client \
        openssh-server 
        
PARAM USERNAME username
PARAM PASSWORD password
EXPOSE 22