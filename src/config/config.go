package config

import (
	"fmt"
	"os"
)

/* Database config */
var (
	connectionName = os.Getenv("CLOUDSQL_CONNECTION_NAME")
	user           = os.Getenv("CLOUDSQL_USER")
	password       = os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty
	database       = os.Getenv("CLOUDSQL_DATABASE")
	Mysql          = "mysql"
	Dbconnection   = fmt.Sprintf("%s:%s@cloudsql(%s)/%s", user, password, connectionName, database)
)
