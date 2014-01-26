FROM ubuntu:12.04

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -s /bin/true /sbin/initctl
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys C300EE8C
RUN apt-get update && DEBIAN_FRONTEND=noninteractive
RUN apt-get install     -yq libmysqlclient-dev
RUN apt-get install     -yq python-dev
RUN apt-get install     -yq libssl-dev
RUN apt-get install     -yq libevent-dev
RUN apt-get install     -yq build-essential
RUN apt-get install     -yq liblzma-dev
RUN apt-get install     -yq mysql-server
RUN apt-get install     -yq mysql-client
RUN apt-get install     -yq nano
RUN apt-get install     -yq wget
RUN apt-get install     -yq git
RUN apt-get install     -yq net-tools

RUN wget -P Downloads/ http://python-distribute.org/distribute_setup.py
RUN python Downloads/distribute_setup.py
RUN easy_install pip
RUN pip install Flask
RUN pip install mysql-python
RUN pip install wtforms
RUN pip install iptools
RUN pip install passlib
RUN pip install pbkdf2
RUN pip install peewee
RUN pip install pyzmq
RUN pip install twilio
RUN git clone https://odonnellryan:0okm9ijn@github.com/xutherogue/Neo
RUN /usr/sbin/mysqld & \
    sleep 4s &&\
    echo "CREATE DATABASE users; GRANT ALL ON *.* TO 'root'@'localhost' IDENTIFIED BY '0okm9ijn' WITH GRANT OPTION; FLUSH PRIVILEGES;" | mysql && \
    python Neo/SoftLaunch/createtable.py

