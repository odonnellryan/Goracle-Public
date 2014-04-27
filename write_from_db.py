import MySQLdb
from _mysql_exceptions import OperationalError
import os
import time
import subprocess
import pwd

USERNAME = "ryan"
PASSWORD = "test"
DATABASE = "nginx"
HOSTNAME = "127.0.0.1"
TABLNAME = "configs"
FILEPATH = "/etc/nginx/sites-enabled"

TABLSTRC = """CREATE TABLE IF NOT EXISTS {0} (
          `id` int(11) NOT NULL AUTO_INCREMENT,
          `name` varchar(2555) DEFAULT NULL,
          `content` varchar(2555) DEFAULT NULL,
          `hostname` varchar(255) NOT NULL UNIQUE,
          `hash` VARCHAR(2555) NOT NULL,
          `write` int(11) NOT NULL DEFAULT '0',
           PRIMARY KEY (`id`)
           ) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;""".format(TABLNAME)

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


class NginxWriter:
    def __init__(self):
        self.rows = None
        self.data = None
        self.conn = None
        self.cursor = None

    def get_db_table(self):
        try:
            rows = self.cursor.execute("SELECT `name`, `content`, `hostname`, `hash`, `write` FROM {0};".format(TABLNAME))
            data = self.cursor.fetchall()
            self.rows = rows
            self.data = data
        except OperationalError as e:
            print ("Server Offline" + str(e))
            time.sleep(5)

    def search_through_db_rows(self):
        data = self.data
        if self.rows:
            for con in range(self.rows):
                self.check_to_write_file(name=data[con][0], content=data[con][1], hostname=data[con][2],
                                    hash=str(data[con][3]), write=data[con][4])
                time.sleep(.01)
            self.conn.commit()

    def check_to_write_file(self, name=None, content=None, hostname=None, hash=None, write=None):
        hash_line = "".join(("#" + hash))
        _file = os.path.join(FILEPATH, name)
        try:
            with open(_file) as f:
                if hash_line == f.readline():
                    pass
                elif write == 0:
                    pass
                else:
                    print("writing to: " + f.name)
                    self.write_file_(_file, content, hash, write)
        except (IOError, OSError):
            self.write_file_(_file, content, hash, write)

    def set_file_not_to_be_written(self, config_file):
        try:
            e_file = config_file.split(".")
            n_file = ".".join((e_file[0], "ERROR"))
            print("Renaming file to: " + "/".join((FILEPATH, n_file)))
            os.rename("/".join((FILEPATH, config_file)), "/".join((FILEPATH, n_file)))
            print("Searching in db for: ", config_file)
            # updates file in database
            #self.execute("UPDATE configs SET name = %s, write = 0 WHERE name = %s", (db_write_file[1], search_file[1]))
            self.cursor.execute("UPDATE configs SET `name` = %s, `write` = 0 WHERE `name` = %s",
                           (n_file, config_file))
            if self.cursor.rowcount:
                print ("Successfully updated DB")
            else:
                print("Could not update DB")
        except OperationalError:
            print ("Database Down")
            return False
        except IOError:
            print ("File Operation Unsuccessful")

    def change_file_owner(self, config_file):
        uid = pwd.getpwnam('www-data')[2]
        os.chown(config_file, uid, uid)

    def start_nginx(self):
        print ("Trying to start Nginx")
        cmd = ['nginx']
        process = subprocess.Popen(cmd, shell=True,
                                   stdout=subprocess.PIPE,
                                   stderr=subprocess.PIPE)
        return process.communicate()

    def reload_nginx(self):
        print("Reloading Nginx")
        cmd = ['nginx -s reload']
        process = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        out, err = process.communicate()
        if ("""open() "/run/nginx.pid" failed""" in err):
            return self.start_nginx()
        print ("Major error - nginx configs not reloading. Error is: " + err)
        if self.fix_nginx_config(err):
            return None, None
        return out, err

    def fix_nginx_config(self, err):
        e_message = err.split("/")
        e_file_path = "/".join(("",e_message[-4],e_message[-3],e_message[-2]))
        e_filename = (e_message[-1]).split(":")[0]
        print("error in file: ", e_filename, "at path: ", e_file_path, "\n")
        # magic string parsing. nginx error is:
        # nginx: [emerg] unknown directive "localhost:9999" in /etc/nginx/sites-enabled/hostname.GOOD:3
        if e_file_path == FILEPATH:
            print ("Attempting to set config to incorrect...")
            # parse nginx error message
            print ("File with error: " + e_filename)
            self.set_file_not_to_be_written(e_filename)
            return False
        else:
            print("Paths do not match. MAJOR NGINX ERROR. Alert...alert...")

    def write_file_(self, config_file, content, hash, write):
        if not write:
            return False
        with open(config_file, "w") as f:
            f.write(content)
        try:
            self.change_file_owner(config_file)
        except (IOError, OSError):
            print("Can not change file permissions")
        out, err = self.reload_nginx()
        if err or out:
            print(out + err)
        return True

    def start(self):
        self.conn = MySQLdb.connect(host=HOSTNAME, user=USERNAME, passwd=PASSWORD, db=DATABASE)
        self.cursor = self.conn.cursor()
        self.cursor.execute(TABLSTRC)
        self.conn.close()
        while True:
            print("Still alive...")
            self.conn = MySQLdb.connect(host=HOSTNAME, user=USERNAME, passwd=PASSWORD, db=DATABASE)
            self.cursor = self.conn.cursor()
            self.get_db_table()
            self.search_through_db_rows()
            self.conn.close()
            time.sleep(1)
            pass

nginx_writer = NginxWriter()
if __name__ == "__main__":
    nginx_writer.start()


