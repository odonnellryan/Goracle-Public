import MySQLdb
from _mysql_exceptions import OperationalError
import os
import time
import subprocess

try:
    WindowsError
except NameError:
    WindowsError = None

#
# TODO: please, please fix the fact that this can break if a config breaks but it isn't one of the fiels in the
# below paths.
#

file_path = "nginx/"

def set_write_check_null(cursor, _file):
    try:
        conn = MySQLdb.connect(host = '127.0.0.1', user = 'root', passwd = '0okm9ijn', db = 'nginx')
        cursor = conn.cursor()
        search_file = _file.split("/")
        e_file = _file.split(".")
        n_file = ".".join((e_file[0],"ERROR"))
        db_write_file = n_file.split("/")
        print "Renaming file to: " + n_file
        os.rename(_file, n_file)
        print "Searching in db for: " + search_file[1]
        cursor.execute("UPDATE configs SET filename = %s, write_check = 0 WHERE filename = %s", (db_write_file[1], search_file[1]))
        if cursor.rowcount:
            print "Successfully updated DB"
        conn.commit()
        conn.close()
    except OperationalError:
        print "Database Down"
        return False

def write_file_(_file, content, file_hash, write_check, cursor):
    with open(_file, "w") as f:
        if write_check:
            content = "".join(["#", file_hash, "\n", content])
            f.write(content)
            try:
                import pwd
                uid = pwd.getpwnam('www-data')[2]
                os.chown(_file,uid,uid)
                cmd = ['nginx -s reload']
                process = subprocess.Popen(cmd, shell=True,
                           stdout=subprocess.PIPE, 
                           stderr=subprocess.PIPE)
                out, err = process.communicate()
                errcode = process.returncode
                if err:
                    # major error if we can't reload...no configs will reload...
                    print "Major error - nginx configs not reloading. Error is: " + err
                    print "Attempting to set config to incorrect..."
                    e_file = err.split("/")
                    if (e_file[-2] + "/") == file_path:
                        error_file = ("".join((file_path,e_file[-1]))).split(":")
                        error_file = error_file[0]
                        print "File with error:" + error_file
                        set_write_check_null(cursor, error_file)
                        return False
            except (ImportError, ReferenceError, WindowsError):
                pass
    return True


def search_and_write():
    #
    # file name MUST ALWAYS be unique. otherwise silly things happen of course.
    #
    try:
        conn = MySQLdb.connect(host = '127.0.0.1', user = 'root', passwd = '0okm9ijn', db = 'nginx')
        cursor = conn.cursor()
        rows  = cursor.execute("SELECT * FROM configs;")
        data = cursor.fetchall()
        for con in range(rows):
            name, content, file_hash, write_check = data[con][1], data[con][2], str(data[con][3]), data[con][4]
            hash_line = "".join(("#" + file_hash))
            _file = os.path.join(file_path, name)
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
    except OperationalError, e:
        print 'Server Offline' + str(e)
        time.sleep(5)

while True:
    search_and_write()