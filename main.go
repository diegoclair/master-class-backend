package main

import (
	"database/sql"
	"log"

	"github.com/diegoclair/master-class-backend/api"
	db "github.com/diegoclair/master-class-backend/db/sqlc"
	"github.com/diegoclair/master-class-backend/util"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
