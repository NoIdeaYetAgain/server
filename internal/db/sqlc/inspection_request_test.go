package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomInspectionRequest(t *testing.T) InspectionRequest {
	// Create a property and student first since they are foreign keys
	property := createRandomProperty(t)
	student := CreateRandomUser(t)

	// Use UTC timezone for consistent testing
	requestedTime := time.Now().UTC().Add(24 * time.Hour)
	arg := CreateInspectionRequestParams{
		PropertyID:    property.ID,
		StudentID:     student.ID,
		RequestedDate: pgtype.Timestamp{Time: requestedTime, Valid: true},
	}

	inspectionRequest, err := testQueries.CreateInspectionRequest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, inspectionRequest)

	require.Equal(t, arg.PropertyID, inspectionRequest.PropertyID)
	require.Equal(t, arg.StudentID, inspectionRequest.StudentID)
	require.True(t, inspectionRequest.RequestedDate.Valid)

	// Compare times in UTC and round to minute to avoid microsecond differences
	expectedTime := arg.RequestedDate.Time.UTC().Round(time.Minute)
	actualTime := inspectionRequest.RequestedDate.Time.UTC().Round(time.Minute)
	require.Equal(t, expectedTime, actualTime)

	// Properly compare pgtype.Text status
	require.True(t, inspectionRequest.Status.Valid)
	require.Equal(t, "pending", inspectionRequest.Status.String)

	require.NotZero(t, inspectionRequest.ID)
	require.True(t, inspectionRequest.CreatedAt.Valid)

	return inspectionRequest
}

func TestCreateInspectionRequest(t *testing.T) {
	createRandomInspectionRequest(t)
}

func TestGetInspectionRequest(t *testing.T) {
	inspectionRequest1 := createRandomInspectionRequest(t)
	inspectionRequest2, err := testQueries.GetInspectionRequest(context.Background(), inspectionRequest1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, inspectionRequest2)

	require.Equal(t, inspectionRequest1.ID, inspectionRequest2.ID)
	require.Equal(t, inspectionRequest1.PropertyID, inspectionRequest2.PropertyID)
	require.Equal(t, inspectionRequest1.StudentID, inspectionRequest2.StudentID)

	// Compare status properly
	require.True(t, inspectionRequest1.Status.Valid)
	require.True(t, inspectionRequest2.Status.Valid)
	require.Equal(t, inspectionRequest1.Status.String, inspectionRequest2.Status.String)

	// Compare times in UTC and round to minute
	expectedTime := inspectionRequest1.RequestedDate.Time.UTC().Round(time.Minute)
	actualTime := inspectionRequest2.RequestedDate.Time.UTC().Round(time.Minute)
	require.Equal(t, expectedTime, actualTime)

	// Compare created_at timestamps
	require.True(t, inspectionRequest1.CreatedAt.Valid)
	require.True(t, inspectionRequest2.CreatedAt.Valid)
	require.WithinDuration(t,
		inspectionRequest1.CreatedAt.Time.UTC(),
		inspectionRequest2.CreatedAt.Time.UTC(),
		time.Second,
	)
}

func TestDeleteInspectionRequest(t *testing.T) {
	inspectionRequest1 := createRandomInspectionRequest(t)

	err := testQueries.DeleteInspectionRequest(context.Background(), inspectionRequest1.ID)
	require.NoError(t, err)

	inspectionRequest2, err := testQueries.GetInspectionRequest(context.Background(), inspectionRequest1.ID)
	require.Error(t, err)
	require.Empty(t, inspectionRequest2)
}
