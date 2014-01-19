FROM jpetazzo/dind

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -s /bin/true /sbin/initctl

RUN echo "deb http://ppa.launchpad.net/nginx/stable/ubuntu lucid main" > /etc/apt/sources.list.d/nginx-stable-lucid.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys C300EE8C

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
	libmysqlclient-dev \
	python-pip \ 
	python-dev \
	linux-headers-generic \
	nginx \
	mysql-server \
	mysql-client \
	nano \
	wget \
	net-tools

RUN mkdir -p /data/db
RUN /usr/sbin/mysqld & \
    sleep 4s &&\
    echo "CREATE DATABASE nginx; GRANT ALL ON *.* TO 'ryan'@'localhost' IDENTIFIED BY 'test' WITH GRANT OPTION; FLUSH PRIVILEGES;" | mysql

RUN pip install mysql-python
EXPOSE 6000
ENTRYPOINT /start.sh

ADD bin /usr/local/bin
ADD scripts /
ADD nginx/ /etc/
