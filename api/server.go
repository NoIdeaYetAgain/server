package api

import (
	db "clove/internal/db/sqlc"
	"clove/token"
	"clove/util"
	"clove/worker"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config          util.Config
	store           *db.SQLStore
	tokenMaker      token.Maker
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store *db.SQLStore, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Public routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	// Authenticated routes (require valid token)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker, []string{"student", "agent", "admin"}))
	{
		// User routes
		authRoutes.GET("/users/me", server.getCurrentUser)
		authRoutes.PUT("/users/me", server.updateCurrentUser)

		// Review routes
		authRoutes.POST("/reviews", server.createReview)
		authRoutes.GET("/reviews", server.getAllReviews)
		authRoutes.GET("/reviews/:id", server.getReview)
		authRoutes.PUT("/reviews/:id", server.updateReview)
		authRoutes.DELETE("/reviews/:id", server.deleteReview)

		// Property routes
		authRoutes.POST("/properties", server.createProperty)
		authRoutes.GET("/properties/:id", server.getProperty)
		authRoutes.GET("/properties", server.getAllProperties)
		authRoutes.PUT("/properties/:id", server.updateProperty)

		// Favorite routes
		authRoutes.POST("/favorites", server.addFavorite)
		authRoutes.GET("/favorites", server.getAllFavorites)
		authRoutes.DELETE("/favorites/:property_id", server.removeFavorite)

		// Inspection routes
		authRoutes.POST("/inspections", server.createInspection)
		authRoutes.GET("/inspections", server.getAllInspections)
		authRoutes.GET("/inspections/:id", server.getInspection)
		authRoutes.DELETE("/inspections/:id", server.deleteInspection)
	}

	// Admin-only routes
	adminRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker, []string{"admin"}))
	{
		adminRoutes.GET("/users/:id", server.getUser)
		adminRoutes.PUT("/users/:id", server.updateUser)
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
