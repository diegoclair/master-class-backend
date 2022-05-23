package gapi

import (
	"fmt"

	db "github.com/diegoclair/master-class-backend/db/sqlc"
	"github.com/diegoclair/master-class-backend/proto/pb"
	"github.com/diegoclair/master-class-backend/token"
	"github.com/diegoclair/master-class-backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	cfg        util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

const (
	JWTType    = "jwt"
	PasetoType = "paseto"
)

func NewServer(cfg util.Config, store db.Store) (*Server, error) {

	var tokenMaker token.Maker
	var err error
	if cfg.AccessTokenType == JWTType {
		tokenMaker, err = token.NewJWTMaker(cfg.TokenSymmetricKey)
	} else {
		tokenMaker, err = token.NewPasetoMaker(cfg.TokenSymmetricKey)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		cfg:        cfg,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
