package api

import (
	db "clove/internal/db/sqlc"
	"clove/token"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"strconv"
	"strings"
)

// Inspection request types
type createInspectionRequest struct {
	PropertyID    int32            `json:"property_id" binding:"required,min=1"`
	RequestedDate pgtype.Timestamp `json:"requested_date" binding:"required"`
}

type getInspectionRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type getAllInspectionsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type deleteInspectionRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// Create inspection handler
func (server *Server) createInspection(ctx *gin.Context) {
	var req createInspectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Only students can create inspection requests
	if authPayload.Role != "student" {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("only students can create inspection requests")))
		return
	}

	arg := db.CreateInspectionRequestParams{
		PropertyID:    req.PropertyID,
		StudentID:     authPayload.UserID,
		RequestedDate: req.RequestedDate,
	}

	inspection, err := server.store.CreateInspectionRequest(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inspection)
}

// Get inspection by ID handler
func (server *Server) getInspection(ctx *gin.Context) {
	var req getInspectionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	inspection, err := server.store.GetInspectionRequest(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Authorization check
	switch authPayload.Role {
	case "admin", "agent":
		// Admins and agents can view any inspection
	case "student":
		if inspection.StudentID != authPayload.UserID {
			ctx.JSON(http.StatusForbidden, errorResponse(errors.New("not authorized to access this inspection")))
			return
		}
	default:
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("unauthorized role")))
		return
	}

	ctx.JSON(http.StatusOK, inspection)
}

// Get all inspections handler
func (server *Server) getAllInspections(ctx *gin.Context) {
	var req getAllInspectionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListInspectionsParams{
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
		Column3:   false, // FilterByStudent
		Column5:   false, // FilterByStatus
		StudentID: 0,
		Column6:   []string{},
	}

	switch authPayload.Role {
	case "student":
		// Students can only see their own requests
		arg.Column3 = true
		arg.StudentID = authPayload.UserID

	case "agent":
		// Agents can see all requests but might want to filter by status
		if statusFilter := ctx.Query("status"); statusFilter != "" {
			arg.Column5 = true
			arg.Column6 = strings.Split(statusFilter, ",")
		}

	case "admin":
		// Admins can see everything with optional filters
		if studentID := ctx.Query("student_id"); studentID != "" {
			if sid, err := strconv.Atoi(studentID); err == nil {
				arg.Column3 = true
				arg.StudentID = int32(sid)
			}
		}
		if statusFilter := ctx.Query("status"); statusFilter != "" {
			arg.Column5 = true
			arg.Column6 = strings.Split(statusFilter, ",")
		}
	}

	inspections, err := server.store.ListInspections(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inspections)
}

func (server *Server) deleteInspection(ctx *gin.Context) {
	var req deleteInspectionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	inspection, err := server.store.GetInspectionRequest(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Authorization check
	switch authPayload.Role {
	case "admin", "agent":
		// Both admins and agents can delete any inspection
		break

	case "student":
		// Students can only delete their own pending inspections
		if inspection.StudentID != authPayload.UserID {
			ctx.JSON(http.StatusForbidden,
				errorResponse(errors.New("not authorized to delete this inspection")))
			return
		}

	default:
		ctx.JSON(http.StatusForbidden,
			errorResponse(errors.New("unauthorized role")))
		return
	}

	err = server.store.DeleteInspectionRequest(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "inspection request deleted"})
}
