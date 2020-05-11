package main

import (
"context"
"database/sql"
"os"
_ "github.com/vertica/vertica-sql-go"
"github.com/vertica/vertica-sql-go/logger"
)

// https://www.vertica.com/blog/the-vertica-sql-driver-for-go/
// 1. connect to a vertica db and run a query
connDB, err := sql.Open("vertica", "vertica://dbadmin:@localhost:5433/dbadmin")
if err != nil {
testLogger.Fatal(err.Error())
os.Exit(1)
}
defer connDB.Close()

