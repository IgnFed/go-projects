# WEB API
---

Creating an API using [gorilla/mux package](https://github.com/gorilla/mux).

First at all, you will need to do <code>go get .</code> to install all the following dependecies:
  * github.com/gorilla/mux
  * github.com/go-sql-driver/mysql

You can change some enviroment variable in the [main.go file](./main.go).
You can change the following properties:
  * DBNAME: the database name (default 'recordings')
  * USERDB: the database user (default 'root')
  * PASSDB: the database password(no password default)
  * DBDRIVER: the database driver (default 'mysql')
