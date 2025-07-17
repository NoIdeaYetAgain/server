package api

import (
	db "clove/internal/db/sqlc"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type createInspectionRequest struct {
	PropertyID    int32            `json:"property_id" binding:"required, min=1"`
	StudentID     int32            `json:"student_id" binding:"required, min=1"`
	RequestedDate pgtype.Timestamp `json:"requested_date" binding:"required"`
}

func (server *Server) createInspection(ctx *gin.Context) {
	var req createInspectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateInspectionRequestParams{
		PropertyID:    req.PropertyID,
		StudentID:     req.StudentID,
		RequestedDate: req.RequestedDate,
	}

	favorite, err := server.store.CreateInspectionRequest(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, favorite)
}

type getInspectionRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getInspection(ctx *gin.Context) {
	var req getInspectionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	property, err := server.store.GetInspectionRequest(ctx, req.ID)
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

type getAllInspectionsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getAllInspections(ctx *gin.Context) {
	var req getAllInspectionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListInspectionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	favorites, err := server.store.ListInspections(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, favorites)
}

type deleteInspectionRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteInspection(ctx *gin.Context) {
	var req deleteInspectionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteInspectionRequest(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "inspection request deleted"})
}
