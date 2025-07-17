package api

import (
	db "clove/internal/db/sqlc"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type createReviewRequest struct {
	StudentID int32  `json:"student_id" binding:"required"`
	Rating    int32  `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required"`
}

func (server *Server) createReview(ctx *gin.Context) {
	var req createReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateReviewParams{
		StudentID: req.StudentID,
		Rating:    pgtype.Int4{Int32: req.Rating, Valid: true},
		Comment:   pgtype.Text{String: req.Comment, Valid: true},
	}

	review, err := server.store.CreateReview(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}

type getReviewRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getReview(ctx *gin.Context) {
	var req getReviewRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.GetReviewByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}

type getAllReviewsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getAllReviews(ctx *gin.Context) {
	var req getAllReviewsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListReviewsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	reviews, err := server.store.ListReviews(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}

type updateReviewRequest struct {
	ID      int32   `uri:"id" binding:"required,min=1"`
	Rating  *int32  `json:"rating,omitempty" binding:"omitempty,min=1,max=5"`
	Comment *string `json:"comment,omitempty"`
}

func (server *Server) updateReview(ctx *gin.Context) {
	var req updateReviewRequest

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

	// Check if review exists first
	existingReview, err := server.store.GetReviewByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Prepare update parameters - use existing values if not provided
	arg := db.UpdateReviewParams{
		ID: req.ID,
	}

	// Update rating if provided, otherwise keep existing
	if req.Rating != nil {
		arg.Rating = pgtype.Int4{Int32: *req.Rating, Valid: true}
	} else {
		arg.Rating = existingReview.Rating
	}

	// Update comment if provided, otherwise keep existing
	if req.Comment != nil {
		arg.Comment = pgtype.Text{String: *req.Comment, Valid: true}
	} else {
		arg.Comment = existingReview.Comment
	}

	review, err := server.store.UpdateReview(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}
