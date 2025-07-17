package api

import (
	db "clove/internal/db/sqlc"
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

	arg := db.CreatePropertyParams{
		AgentID:      req.AgentID,
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

	// Bind URI parameter
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Bind JSON body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if property exists first
	existingProperty, err := server.store.GetProperty(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Prepare update parameters - use existing values if not provided
	arg := db.UpdatePropertyParams{
		ID: req.ID,
	}

	// Update title if provided, otherwise keep existing
	if req.Title != nil {
		arg.Title = *req.Title
	} else {
		arg.Title = existingProperty.Title
	}

	// Update description if provided, otherwise keep existing
	if req.Description != nil {
		arg.Description = *req.Description
	} else {
		arg.Description = existingProperty.Description
	}

	// Update price if provided, otherwise keep existing
	if req.Price != nil {
		arg.Price = *req.Price
	} else {
		arg.Price = existingProperty.Price
	}

	// Update location if provided, otherwise keep existing
	if req.Location != nil {
		arg.Location = *req.Location
	} else {
		arg.Location = existingProperty.Location
	}

	// Update property type if provided, otherwise keep existing
	if req.PropertyType != nil {
		arg.PropertyType = *req.PropertyType
	} else {
		arg.PropertyType = existingProperty.PropertyType
	}

	// Update image URL if provided, otherwise keep existing
	if req.ImageUrl != nil {
		arg.ImageUrl = *req.ImageUrl
	} else {
		arg.ImageUrl = existingProperty.ImageUrl
	}

	// Update video URL if provided, otherwise keep existing
	if req.VideoUrl != nil {
		arg.VideoUrl = *req.VideoUrl
	} else {
		arg.VideoUrl = existingProperty.VideoUrl
	}

	property, err := server.store.UpdateProperty(ctx, arg)
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
