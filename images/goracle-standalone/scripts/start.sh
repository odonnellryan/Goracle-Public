mv nginx /etc/nginx/sites-enabled/
mv nginx.conf /etc/nginx/
printf "ryan:$(openssl passwd -crypt test)\n" >> /etc/nginx/htpassword
nginx
usr/bin/mongod & > mongo.log
mysqld & > mysql.log
python write_from_db.py & > pywrite.log
#/opt/goracle/goracle
/bin/bash
