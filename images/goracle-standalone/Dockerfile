from goracle

RUN echo "deb http://archive.ubuntu.com/ubuntu/ precise universe" >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get -y install dialog net-tools lynx nano wget
RUN apt-get -y install python-software-properties
RUN add-apt-repository -y ppa:nginx/stable
RUN apt-get update
RUN apt-get -y install nginx

RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | tee /etc/apt/sources.list.d/mongodb.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys C300EE8C
RUN 	apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
	linux-generic \
	mongodb-10gen \
	mysql-server mysql-client

ENV mongo-server localhost

ADD scripts /

ENTRYPOINT /scripts/.start.sh
