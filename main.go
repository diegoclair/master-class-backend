package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/diegoclair/master-class-backend/api"
	db "github.com/diegoclair/master-class-backend/db/sqlc"
	"github.com/diegoclair/master-class-backend/gapi"
	"github.com/diegoclair/master-class-backend/proto/pb"
	"github.com/diegoclair/master-class-backend/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGrpcGatewayServer(cfg, store)
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

func runGrpcGatewayServer(cfg util.Config, store db.Store) {
	server, err := gapi.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create http grpc gateway server: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//this options are used to generate the json response with same name defined in protobuf file
	//here is the link of documentation https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/
	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register http grpc gateway handler server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	//the files that are inside of ./docs/swagger we got from https://github.com/swagger-api/swagger-ui inside of dist folder
	//we need to modify the line 6 (url) of file swagger-initializer.js inside of docs/swagger to read our file (simple_bank.swagger.json)
	fs := http.FileServer(http.Dir("./docs/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", cfg.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create http grpc gateway listener: ", err)
	}

	log.Printf("Start http gRPC gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http grpc gateway server: ", err)
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
