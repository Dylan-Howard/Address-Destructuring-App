package main

// import "github.com/microsoft/go-mssqldb"
// "github.com/denisenkom/go-mssqldb"
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	mssql "github.com/denisenkom/go-mssqldb"
)

type ConnectionSettings struct {
	Hostname string	`json:"hostname"`
	Database string `json:"database"`
	Port 		 int		`json:"port"`
	Username string	`json:"username"`
	Password string	`json:"password"`
}

type DatabaseQueries struct {
	Get 	 string `json:"getAddresses"`
	Update string `json:"updateAddress"`
	Remove string `json:"removeAddress"`
}

var db *sql.DB
var _ mssql.Logger // This line uses the package but doesn't do anything meaningful

type InfiniteCampus struct {
	Settings ConnectionSettings
	Queries DatabaseQueries
}

func (ic InfiniteCampus) Connect() bool {
	/* Build connection string */
	connUrl := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(ic.Settings.Username, ic.Settings.Password),
		Host:     fmt.Sprintf("%s:%d", ic.Settings.Hostname, ic.Settings.Port),
	}

	var err error

	/* Create connection pool */
	db, err = sql.Open("mssql", connUrl.String())
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	return true
}

func (ic InfiniteCampus) GetAddresses() []Address {

	ic.Connect()

	/* Read addresses */
	addresses, err := ic.ReadAddresses()
	if err != nil {
		log.Fatal("Error reading Addresses: ", err.Error())
	}
	fmt.Printf("Read %d row(s) successfully.\n", len(addresses))

	return addresses
}

func (ic InfiniteCampus) ReadAddresses() ([]Address, error) {
	ctx := context.Background()

	/* Check if database is alive */
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tsql := ic.Queries.Get

	/* Execute query */
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Address

	/* Iterate through the results */
	for rows.Next() {
		var number, street, tag, apt, city, state, zip, schools string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &number, &street, &tag, &apt, &city, &state, &zip, &schools)
		if err != nil {
			return nil, err
		}

		var address Address
		address.Id = id
		address.StreetNumber = number
		address.StreetName = street + " " + tag
		address.Unit = apt
		address.State = state
		address.Zip = zip
		
		schoolList := strings.Split(schools, ";")

		for i := 0; i < len(schoolList); i++ {
			if strings.Contains(schoolList[i], " Elementary") {
				address.Region.Elementary = strings.Split(schoolList[i], " Elementary")[0]
			}
			if strings.Contains(schoolList[i], " Middle") {
				address.Region.Middle = strings.Split(schoolList[i], " Middle")[0]
			}
			if strings.Contains(schoolList[i], " High") {
				address.Region.High = strings.Split(schoolList[i], " High")[0]
			}
		}

		list = append(list, address)

		// fmt.Printf("ID: %d, Address: %s %s %s %s %s, %s %s; {%s}\n", id, number, street, tag, apt, city, state, zip, schools)
	}

	return list, nil
}