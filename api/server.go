package api

import (
	"fmt"

	db "github.com/diegoclair/master-class-backend/db/sqlc"
	"github.com/diegoclair/master-class-backend/token"
	"github.com/diegoclair/master-class-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
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
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)
	router.POST("/tokens/refresh", s.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccountByID)
	authRoutes.GET("/accounts", s.getAccounts) // if we put / at final of route, it broke the test.. like this: /accounts/

	authRoutes.POST("/transfers", s.createTransfer)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
