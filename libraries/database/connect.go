package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/walterjwhite/go-application/libraries/logging"
)

// TODO: this is specific to Oracle, generalize this
type ConnectionConfiguration struct {
	Username   string
	Password   string
	Host       string
	Port       int
	Service    string
	DriverName string
}

func (configuration *ConnectionConfiguration) Connect() *sqlx.DB {
	db, err := sqlx.Open(configuration.DriverName, configuration.getConnectionString())
	logging.Panic(err)

	testConnection(db)

	return db
}

func (configuration *ConnectionConfiguration) getConnectionString() string {
	return fmt.Sprintf("%v/%v@%v:%d/%v", configuration.Username, configuration.Password, configuration.Host, configuration.Port, configuration.Service)
}

func testConnection(db *sqlx.DB) {
	logging.Panic(db.Ping())
}
