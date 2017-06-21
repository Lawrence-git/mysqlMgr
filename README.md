[![Go Report Card](https://goreportcard.com/badge/github.com/koepkeca/mysqlMgr)](https://goreportcard.com/report/github.com/koepkeca/mysqlMgr)

[![GoDoc](https://godoc.org/github.com/koepkeca/mysqlMgr?status.svg)](https://godoc.org/github.com/koepkeca/mysqlMgr)

# Overview

mysqlMgr is a library which wraps database connections and prepared statments in a single data structure which is safe for concurrent use.

# Installation

To install the library use go get:

```
go get github.com/koepkeca/mysqlMgr
```

# Configuration

Each instance of the library is designed to manage one database connection. The New method takes an [io.Reader](https://godoc.org/io#Reader) which contains json configuration data. An example of the file is here:

```
{
    "database" : "mysqlMgr_Demo",
    "server" : "localhost",
    "port" : "3306",
    "driver" : "mysql",
    "user" : "changeme",
    "pw" : "changeme"
}
```
# Usage

To use the library you can create a new instance with an io.Reader that contains the json configuration. A **very basic** example is here, please note, this example skips error checking and a more complete example is available in the example folder.

```
confRdr, _ := ioutil.ReadFile("config.conf")
confByte := bytes.NewBuffer(confRdr)
s, _ := mysqlMgr.New(confByte)
_ = s.AddStmt("statement","INSERT INTO table (value1, value2) VALUES (?,?)")
//then get the statement
stmt, _ := s.GetStmt("statement")
s.Close()
```

# Example program

There is an example program that demonstrates the library. It is fully functional but requires a few steps to set it up.

## Prereqs

The example program requires a few items.

* A words file (linux dictionary) located in /usr/share/dict/words used to populate the random data pool. You can use any list of strings seperated by new lines.
* A mysql database with the test database created. The schema is located in test_database.sql in the example directory.
* A **properly configured** configuration file. There is an example configuration file which needs to be modified with your specific configuration values in order to function properly.

This example program creates a random data pool and the sequentially load 8192 random records into the example database. You can adjust this number by changing the maxElements in example.go

