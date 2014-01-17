from goracle

RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | tee /etc/apt/sources.list.d/10gen.list

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -s /bin/true /sbin/initctl

RUN apt-get update
RUN apt-get -y install dialog net-tools lynx nano wget
RUN apt-get -y install python-software-properties
RUN add-apt-repository -y ppa:nginx/stable
RUN 	apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
	libmysqlclient-dev \
	python-pip \ 
	python-dev \
	linux-generic \
	nginx \
	mysql-server mysql-client

ADD scripts /

RUN mkdir -p /data/db
RUN /usr/sbin/mysqld & \
    sleep 10s &&\
    echo "CREATE DATABASE nginx; GRANT ALL ON *.* TO 'ryan'@'localhost' IDENTIFIED BY 'test' WITH GRANT OPTION; FLUSH PRIVILEGES;" | mysql
RUN pip install mysql-python
EXPOSE 6000

RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | tee /etc/apt/sources.list.d/mongodb.list
RUN apt-get update
RUN apt-get -y install mongodb-10gen

ENTRYPOINT ./start.sh
