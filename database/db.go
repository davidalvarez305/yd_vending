package database

import (
	"database/sql"
	"fmt"

	"github.com/davidalvarez305/yd_vending/constants"
	_ "github.com/lib/pq"
)

type connection struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

var DB *sql.DB

func Connect() (*sql.DB, error) {
	conn := connection{
		host:     constants.PostgresHost,
		port:     constants.PostgresPort,
		user:     constants.PostgresUser,
		password: constants.PostgresPassword,
		dbName:   constants.PostgresDBName,
	}

	connectionString := connToString(conn)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to verify connection: %v", err)
	}

	DB = db
	return db, nil
}

func connToString(info connection) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		info.user, info.password, info.host, info.port, info.dbName)
}
