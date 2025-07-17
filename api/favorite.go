package api

import (
	db "clove/internal/db/sqlc"
	"clove/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type addFavoriteRequest struct {
	PropertyID int32 `json:"property_id" binding:"required,min=1"`
}

func (server *Server) addFavorite(ctx *gin.Context) {
	var req addFavoriteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.UserID

	arg := db.AddFavoriteParams{
		StudentID:  userID,
		PropertyID: req.PropertyID,
	}

	favorite, err := server.store.AddFavorite(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, favorite)
}

type getAllFavoritesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getAllFavorites(ctx *gin.Context) {
	var req getAllFavoritesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.UserID

	arg := db.GetFavoritesApartmentsParams{
		StudentID: userID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	favorites, err := server.store.GetFavoritesApartments(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, favorites)
}

func (server *Server) removeFavorite(ctx *gin.Context) {
	propertyIDStr := ctx.Param("property_id")
	propertyID, err := strconv.Atoi(propertyIDStr)
	if err != nil || propertyID < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.UserID

	arg := db.RemoveFavoriteParams{
		StudentID:  userID,
		PropertyID: int32(propertyID),
	}

	err = server.store.RemoveFavorite(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}
