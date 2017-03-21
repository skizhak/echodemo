package resources

import (
	"bytes"
	"database/sql"
	"log"
)

const (
	dbUser     string = "echodemo"
	dbPassword string = "demo123"
	dbName     string = "echodemo"
)

// ConnectDB Connect to mysql DB
func ConnectDB() *sql.DB {
	var dbCon bytes.Buffer
	dbCon.WriteString(dbUser)
	dbCon.WriteString(":")
	dbCon.WriteString(dbPassword)
	dbCon.WriteString("@tcp(127.0.0.1:3306)/")
	dbCon.WriteString(dbName)

	db, err := sql.Open("mysql", dbCon.String())

	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
