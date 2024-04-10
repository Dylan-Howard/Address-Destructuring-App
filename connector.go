package main

// import "github.com/microsoft/go-mssqldb"
import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type DatabaseSettings struct {
	Hostname string	`json:"hostname"`
	Port 		 string	`json:"port"`
	Username string	`json:"username"`
	Password string	`json:"password"`
}

var db *sql.DB

func GetAddresses(settings DatabaseSettings) {

	// Build connection string
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s", settings.Username, settings.Password, settings.Hostname, settings.Port)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
			log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	// Read addresses
	count, err := ReadAddresses()
	if err != nil {
		log.Fatal("Error reading Addresses: ", err.Error())
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
}

// ReadAddresses reads all address records
func ReadAddresses() (int, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("SELECT %s FROM Address;", "*")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return -1, err
	}

	defer rows.Close()

	var count int

	// Iterate through the result set.
	for rows.Next() {
		var name, location string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &name, &location)
		if err != nil {
			return -1, err
		}

		fmt.Printf("ID: %d, Name: %s, Location: %s\n", id, name, location)
		count++
	}

	return count, nil
}