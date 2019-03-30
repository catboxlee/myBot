package mydb

import (
	"database/sql"
	"log"
	"os"

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

func init() {
	DbStart()
}

// DbStart ...
func DbStart() {

	// Initialize connection string.
	//var connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", HOST, USER, PASSWORD, DATABASE)

	// Initialize connection object.
	var err error
	//log.Println("Start creat connection to database...")
	Db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully created connection to database")

}
