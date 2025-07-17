package api

import (
	db "clove/internal/db/sqlc"
	"clove/token"
	"clove/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/api").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:id", server.getUser)

	authRoutes.POST("/reviews", server.createReview)
	authRoutes.GET("/reviews", server.getAllReviews)
	authRoutes.GET("/reviews/:id", server.getReview)
	authRoutes.PUT("/reviews/:id", server.updateReview)

	authRoutes.POST("/listings", server.createProperty)
	authRoutes.GET("/listings/:id", server.getProperty)
	authRoutes.GET("/listings", server.getAllProperties)
	authRoutes.PUT("/listings/:id", server.updateProperty)

	authRoutes.POST("/favorites", server.addFavorite)
	authRoutes.GET("/favorites", server.getAllFavorites)
	authRoutes.DELETE("/favorites/:property_id", server.removeFavorite)

	authRoutes.POST("/book-inspection", server.createInspection)
	authRoutes.GET("/inspection-requests", server.getAllInspections)
	authRoutes.GET("/inspection-requests/:id", server.getInspection)
	authRoutes.DELETE("inspection-requests/:id", server.deleteInspection)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
