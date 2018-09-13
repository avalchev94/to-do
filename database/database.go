package database

import (
	"database/sql"
	"fmt"
)

// ConnectionInfo wraps all the needed information for establishing connection with PostgresSQL.
type ConnectionInfo struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (conn *ConnectionInfo) connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conn.Host, conn.Port, conn.User, conn.Password, conn.DBName)
}

// Connect establish connection with postgres database using the given connection data.
func Connect(connInfo ConnectionInfo) (*sql.DB, error) {
	return sql.Open("postgres", connInfo.connectionString())
}
