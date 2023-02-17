package main

import (
	"database/sql"
	"log"
	"github.com/Annongkhanh/Go_example/api"
	db "github.com/Annongkhanh/Go_example/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource) 

	if (err != nil){
		log.Fatal("Can not connect to database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if (err != nil){
		log.Fatal("Can not start server: ", err)
	}

}