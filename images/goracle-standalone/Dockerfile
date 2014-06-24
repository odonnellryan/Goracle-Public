FROM jpetazzo/dind

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -s /bin/true /sbin/initctl

RUN echo "deb http://ppa.launchpad.net/nginx/stable/ubuntu lucid main" > /etc/apt/sources.list.d/nginx-stable-lucid.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys C300EE8C

RUN apt-get update && DEBIAN_FRONTEND=noninteractive
RUN apt-get install	-yq libmysqlclient-dev
RUN apt-get install	-yq python-pip 
RUN apt-get install	-yq python-dev
RUN apt-get install	-yq libssl-dev 
RUN apt-get install	-yq libevent-dev
RUN apt-get install	-yq build-essential
RUN apt-get install	-yq liblzma-dev
RUN apt-get install	-yq linux-headers-generic
RUN apt-get install	-yq nginx
RUN apt-get install	-yq mysql-server
RUN apt-get install	-yq mysql-client
RUN apt-get install	-yq nano
RUN apt-get install	-yq wget
RUN apt-get install	-yq git
RUN apt-get install	-yq net-tools
RUN apt-get install	-yq iputils-ping
RUN apt-get install	-yq arping
RUN apt-get install	-yq curl
RUN pip install mysql-python

# install latest docker
RUN curl -s https://get.docker.io/ubuntu/ | sh

# registry stuff and generate docker login to local directory
RUN cd /opt && git clone https://github.com/dotcloud/docker-registry.git \
    && cp /opt/docker-registry/config/config_sample.yml /opt/docker-registry/config/config.yml
# install registry requirements
RUN cd /opt/docker-registry \
    && pip install -r requirements.txt 

# make mongo directories
RUN mkdir -p /data/db

# mysql server setup
RUN /usr/sbin/mysqld & \
    sleep 4s &&\
    echo "CREATE DATABASE nginx; GRANT ALL ON *.* TO 'ryan'@'localhost' IDENTIFIED BY 'test' WITH GRANT OPTION; FLUSH PRIVILEGES;" | mysql
ADD bin /usr/local/bin
ADD scripts /
RUN cp /pipework /usr/bin/pipework
RUN chmod +x /usr/bin/pipework
ADD nginx /etc/nginx/
EXPOSE 6000
ENTRYPOINT /start.sh
