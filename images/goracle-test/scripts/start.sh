#!/bin/bash
# creates the password for httpbasic
printf "ryan:$(openssl passwd -crypt test)\n" >> /etc/nginx/htpassword
# start nginx, mongo, mysql, python script, goracle, docker
nginx
mongod --smallfiles > /var/log/mongo.log 2>&1 &
mysqld > /var/log/mysql.log 2>&1 &
sleep 2
python write_from_db.py > /var/log/pywrite.log 2>&1 &
goracle > /var/log/goracle.log 2>&1 &
wrapdocker > /var/log/docker.log 2>&1 &
# nginx needs to own the socket for docker so we con
# wrap it in http basic
# give docker a bit to initialize...
sleep 1
chown www-data /var/run/docker.sock
# go coverage/test tools
cp /gtest /usr/bin/gtest
chmod +x /usr/bin/gtest
go test -coverprofile=coverage.out 
go tool cover -func=coverage.out
# uncomment the below if you wish to enter the console 
# for the test container at the end of testing
/bin/bash