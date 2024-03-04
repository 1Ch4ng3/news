package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

func ConnectDB() (*sql.DB, error) {
	// Connect to the database
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/laravel_projectnews")
	if err != nil {
		return nil, err
	}

	// Check the connection to the database
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
