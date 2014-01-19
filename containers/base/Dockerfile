FROM		ubuntu
MAINTAINER	Luma

# sets username and password (required for the container)
PARAM USERNAME user
PARAM PASSWORD pass

# installs git (why not), ssh.
RUN     apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
        git \
        openssh-client \
        openssh-server 
        
# adds the start script
ADD start.sh .

# expose ssh
EXPOSE 22