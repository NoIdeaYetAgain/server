package api

import (
	db "clove/internal/db/sqlc"
	"clove/token"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type createPropertyRequest struct {
	AgentID      pgtype.Int4 `json:"agent_id" binding:"required"`
	Title        string      `json:"title" binding:"required"`
	Description  pgtype.Text `json:"description" binding:"required"`
	Price        int32       `json:"price"  binding:"required"`
	Location     string      `json:"location" binding:"required"`
	PropertyType pgtype.Text `json:"property_type" binding:"required"`
	ImageUrl     pgtype.Text `json:"image_url" binding:"required"`
	VideoUrl     pgtype.Text `json:"video_url" binding:"required"`
}

func (server *Server) createProperty(ctx *gin.Context) {
	var req createPropertyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get agent ID from auth payload instead of request body
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Only allow agents to create properties under their own ID
	if authPayload.Role == "agent" && authPayload.UserID != req.AgentID.Int32 {
		ctx.JSON(http.StatusForbidden,
			errorResponse(errors.New("agents can only create properties for themselves")))
		return
	}

	arg := db.CreatePropertyParams{
		AgentID:      pgtype.Int4{Int32: authPayload.UserID, Valid: true},
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		Location:     req.Location,
		PropertyType: req.PropertyType,
		ImageUrl:     req.ImageUrl,
		VideoUrl:     req.VideoUrl,
	}

	property, err := server.store.CreateProperty(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, property)
}

type getPropertyRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getProperty(ctx *gin.Context) {
	var req getPropertyRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	property, err := server.store.GetProperty(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, property)
}

type getAllPropertiesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getAllProperties(ctx *gin.Context) {
	var req getAllPropertiesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListPropertiesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	properties, err := server.store.ListProperties(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, properties)
}

type updatePropertyRequest struct {
	ID           int32        `uri:"id" binding:"required,min=1"`
	Title        *string      `json:"title,omitempty"`
	Description  *pgtype.Text `json:"description,omitempty"`
	Price        *int32       `json:"price,omitempty"`
	Location     *string      `json:"location,omitempty"`
	PropertyType *pgtype.Text `json:"property_type,omitempty"`
	ImageUrl     *pgtype.Text `json:"image_url,omitempty"`
	VideoUrl     *pgtype.Text `json:"video_url,omitempty"`
}

func (server *Server) updateProperty(ctx *gin.Context) {
	var req updatePropertyRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	existingProperty, err := server.store.GetProperty(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Authorization check
	if authPayload.Role == "agent" && existingProperty.AgentID.Int32 != authPayload.UserID {
		ctx.JSON(http.StatusForbidden,
			errorResponse(errors.New("can only update your own properties")))
		return
	}

	// Prepare update parameters
	arg := db.UpdatePropertyParams{
		ID:           req.ID,
		Title:        existingProperty.Title,
		Description:  existingProperty.Description,
		Price:        existingProperty.Price,
		Location:     existingProperty.Location,
		PropertyType: existingProperty.PropertyType,
		ImageUrl:     existingProperty.ImageUrl,
		VideoUrl:     existingProperty.VideoUrl,
	}

	// Apply updates from request
	if req.Title != nil {
		arg.Title = *req.Title
	}
	if req.Description != nil {
		arg.Description = *req.Description
	}
	if req.Price != nil {
		arg.Price = *req.Price
	}
	if req.Location != nil {
		arg.Location = *req.Location
	}
	if req.PropertyType != nil {
		arg.PropertyType = *req.PropertyType
	}
	if req.ImageUrl != nil {
		arg.ImageUrl = *req.ImageUrl
	}
	if req.VideoUrl != nil {
		arg.VideoUrl = *req.VideoUrl
	}

	property, err := server.store.UpdateProperty(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, property)
}
