printf "ryan:$(openssl passwd -crypt test)\n" >> /etc/nginx/htpassword
nginx
mongod > /var/log/mongo.log 2>&1 &
mysqld > /var/log/mysql.log 2>&1 &
python write_from_db.py > /var/log/pywrite.log 2>&1 &
goracle > /var/log/goracle.log 2>&1 &
wrapdocker > /var/log/docker.log 2>&1 &
/bin/bash
