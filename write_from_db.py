import MySQLdb
from _mysql_exceptions import OperationalError
import os
import time
import subprocess

USERNAME = "ryan"
PASSWORD = "test"
DATABASE = "nginx"
HOSTNAME = "127.0.0.1"
TABLNAME = "configs"
FILEPATH = "/etc/nginx/sites-enabled"

TABLSTRC = """CREATE TABLE IF NOT EXISTS %s (
          `id` int(11) NOT NULL AUTO_INCREMENT,
          `name` varchar(2555) DEFAULT NULL,
          `content` varchar(25555) DEFAULT NULL,
          `hash` VARCHAR(2555) CHARACTER SET utf8 COLLATE utf8_cs NOT NULL UNIQUE,
          `write` int(11) NOT NULL DEFAULT '0',
           PRIMARY KEY (`id`)
           ) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;""" % (TABLNAME)

""" 	
	This file currently loops through the DB, comparing hash-values
	(the first line in the contents) to existing files. This is
	so that it's possible to write to a DB and have that replicated to a filesystem. 
	This is so we can easily add nginx configurations.
	the config files MUST end in .GOOD (per the below nginx.conf config) and this script
	DOES NOT remove the files. 
	However, it does rename files that are giving nginx problems. 
	NGINX line to edit: include /etc/nginx/sites-enabled/*.GOOD;

"""

def set_write_check_null(cursor, _file):
    try:
        conn = MySQLdb.connect(host=HOSTNAME, user=USERNAME, passwd=PASSWORD, db=DATABASE)
        cursor = conn.cursor()
        search_file = _file.split("/")
        e_file = _file.split(".")
        n_file = ".".join((e_file[0],"ERROR"))
        db_write_file = n_file.split("/")
        print("Renaming file to: " + n_file)
        os.rename(_file, n_file)
        print("Searching in db for: " + search_file[1])
        #updates file in database
        cursor.execute("UPDATE %s SET filename = %s, write_check = 0 WHERE filename = %s", (TABLNAME, db_write_file[1], search_file[1]))
        if cursor.rowcount:
            print ("Successfully updated DB")
        conn.commit()
        conn.close()
    except OperationalError:
        print ("Database Down")
        return False
    except IOError:
        print ("File Operation Unsuccessful")

def write_file_(_file, content, file_hash, write_check, cursor):
    if not write_check:
	return False
    with open(_file, "w") as f:
        content = "".join(["#", file_hash, "\n", content])
        f.write(content)
        try:
            import pwd
            uid = pwd.getpwnam('www-data')[2]
            os.chown(_file,uid,uid)
            cmd = ['nginx && nginx -s reload']
            process = subprocess.Popen(cmd, shell=True,
                       stdout=subprocess.PIPE, 
                       stderr=subprocess.PIPE)
            out, err = process.communicate()
	    if ("""open() "/run/nginx.pid" failed""" in err):
                print ("Nginx not running...")
                return False
            if err:
                # major error if we can't reload...
                print ("Major error - nginx configs not reloading. Error is: " + err)
                e_file = err.split("/")
                if (e_file[-2] + "/") == FILEPATH:
                    print ("Attempting to set config to incorrect...")
                    error_file = ("".join((FILEPATH,e_file[-1]))).split(":")
                    error_file = error_file[0]
                    print ("File with error:" + error_file)
                    set_write_check_null(cursor, error_file)
                    return False
        except (ImportError, ReferenceError):
            print("Was not able to change the file permissions")
            return False
    return True


def search_and_write():
    #
    # file name MUST ALWAYS be unique. otherwise silly things happen of course.
    #
    try:
	conn = MySQLdb.connect(host=HOSTNAME, user=USERNAME, passwd=PASSWORD, db=DATABASE)
        cursor = conn.cursor()
	cursor.execute(TABLSTRC)
        rows  = cursor.execute("SELECT * FROM {0};".format(TABLNAME))
        data = cursor.fetchall()
        for con in range(rows):
            name, content, file_hash, write_check = data[con][1], data[con][2], str(data[con][3]), data[con][4]
            hash_line = "".join(("#" + file_hash))
            _file = os.path.join(FILEPATH, name)
            try:
                with open(_file) as f:
                    if hash_line == f.readline():
                        pass
                    else:
                        write_file_(_file,content,file_hash,write_check,cursor)
            except IOError:
                        write_file_(_file,content,file_hash,write_check,cursor)
            time.sleep(.01)
        conn.commit()
        conn.close()
        time.sleep(1)
    except OperationalError as e:
        print ("Server Offline" + str(e))
        time.sleep(5)

while True:
    search_and_write()
