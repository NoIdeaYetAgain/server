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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateReviewParams{
		StudentID: authPayload.UserID,
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

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	existingReview, err := server.store.GetReviewByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Authorization check - only owner or admin can update
	if authPayload.Role != "admin" && existingReview.StudentID != authPayload.UserID {
		ctx.JSON(http.StatusForbidden,
			errorResponse(errors.New("can only update your own reviews")))
		return
	}

	// Prepare update parameters
	arg := db.UpdateReviewParams{
		ID:      req.ID,
		Rating:  existingReview.Rating,
		Comment: existingReview.Comment,
	}

	// Apply updates from request
	if req.Rating != nil {
		arg.Rating = pgtype.Int4{Int32: *req.Rating, Valid: true}
	}
	if req.Comment != nil {
		arg.Comment = pgtype.Text{String: *req.Comment, Valid: true}
	}

	review, err := server.store.UpdateReview(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}

type deleteReviewRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteReview(ctx *gin.Context) {
	var req deleteReviewRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	review, err := server.store.GetReviewByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Only allow admin or review owner to delete
	if authPayload.Role != "admin" && review.StudentID != authPayload.UserID {
		ctx.JSON(http.StatusForbidden,
			errorResponse(errors.New("only admin or review owner can delete reviews")))
		return
	}

	err = server.store.DeleteReview(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}
