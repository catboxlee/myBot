package mydb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Initialize connection constants.
const (
	HOST     = "ec2-184-73-153-64.compute-1.amazonaws.com"
	DATABASE = "dc02s9dkrm6bf"
	USER     = "zmmyjmozivryuq"
	PASSWORD = "d0b0e9a1c1ec22cb440edd82def44ea2085a64b813ee8dfadd8bebaa1a972038"
)

// Db ...
var Db *sql.DB

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// DbStart ...
func DbStart() {

	// Initialize connection string.
	var connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE)

	// Initialize connection object.
	var err error
	Db, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	err = Db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully created connection to database")

}
