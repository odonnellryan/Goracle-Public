import MySQLdb
import os
import time

file_path = "testing/"

def search_and_write():

    #
    # file name MUST ALWAYS be unique. otherwise silly things happen of course.
    #

    conn = MySQLdb.connect(host = '127.0.0.1', user = 'neo', passwd = 'CA5EfraCra6ABuga', db = 'configs')
    cursor = conn.cursor()
    rows  = cursor.execute("SELECT * FROM configs;")
    data = cursor.fetchall()
    for con in range(rows):
        name, content, hash = data[con][1], data[con][2], str(data[con][3])
        hash_line = "".join(("#" + hash))
        file = os.path.join(file_path, name)
        try:
            with open(file) as f:
                if hash_line == f.readline():
                    pass
                else:
                    with open(file, "w") as f:
                        content = "".join(["#", hash, "\n", content])
                        f.write(content)
        except IOError:
            with open(file, "w") as f:
                content = "".join(["#", hash, "\n", content])
                f.write(content)
        # put some sleeps below to ensure this doesn't eat the entire disk always
        time.sleep(.1)
    conn.close()
    time.sleep(1)

while True:
    search_and_write()

