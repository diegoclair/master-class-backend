package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/diegoclair/master-class-backend/api"
	db "github.com/diegoclair/master-class-backend/db/sqlc"
	"github.com/diegoclair/master-class-backend/gapi"
	"github.com/diegoclair/master-class-backend/proto/pb"
	"github.com/diegoclair/master-class-backend/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGrpcServer(cfg, store)

}
func runGrpcServer(cfg util.Config, store db.Store) {
	server, err := gapi.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer) //its optional but it is like a documentation for how to call the server

	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create grpc listener: ", err)
	}

	log.Printf("Start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server: ", err)
	}

}
func runGinServer(cfg util.Config, store db.Store) {
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(cfg.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
