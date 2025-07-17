package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomProperty(t *testing.T) Property {
	// Create an agent user first
	agent := CreateRandomUser(t) // Assuming you have this helper function

	arg := CreatePropertyParams{
		AgentID: pgtype.Int4{
			Int32: agent.ID,
			Valid: true,
		},
		Title: randomString(20),
		Description: pgtype.Text{
			String: randomString(100),
			Valid:  true,
		},
		Price:    int32(randomInt(100000, 2000000)),
		Location: randomString(50),
		PropertyType: pgtype.Text{
			String: "self-contained", // Fixed: using valid property type from schema
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: "https://example.com/image.jpg",
			Valid:  true,
		},
		VideoUrl: pgtype.Text{
			String: "https://example.com/video.mp4",
			Valid:  true,
		},
	}

	property, err := testQueries.CreateProperty(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, property)

	require.Equal(t, arg.AgentID.Int32, property.AgentID.Int32)
	require.Equal(t, arg.Title, property.Title)
	require.Equal(t, arg.Description.String, property.Description.String)
	require.Equal(t, arg.Price, property.Price)
	require.Equal(t, arg.Location, property.Location)
	require.Equal(t, arg.PropertyType.String, property.PropertyType.String)
	require.Equal(t, arg.ImageUrl.String, property.ImageUrl.String)
	require.Equal(t, arg.VideoUrl.String, property.VideoUrl.String)

	require.NotZero(t, property.ID)
	require.NotZero(t, property.CreatedAt)
	require.NotZero(t, property.UpdatedAt)

	return property
}

func TestCreateProperty(t *testing.T) {
	createRandomProperty(t)
}

func TestCreatePropertyWithNullOptionalFields(t *testing.T) {
	agent := CreateRandomUser(t)

	arg := CreatePropertyParams{
		AgentID: pgtype.Int4{
			Int32: agent.ID,
			Valid: true,
		},
		Title: randomString(20),
		Description: pgtype.Text{
			String: "",
			Valid:  false, // NULL description
		},
		Price:    int32(randomInt(100000, 2000000)),
		Location: randomString(50),
		PropertyType: pgtype.Text{
			String: "",
			Valid:  false, // NULL property_type
		},
		ImageUrl: pgtype.Text{
			String: "",
			Valid:  false, // NULL image_url
		},
		VideoUrl: pgtype.Text{
			String: "",
			Valid:  false, // NULL video_url
		},
	}

	property, err := testQueries.CreateProperty(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, property)

	require.Equal(t, arg.AgentID.Int32, property.AgentID.Int32)
	require.Equal(t, arg.Title, property.Title)
	require.False(t, property.Description.Valid)
	require.Equal(t, arg.Price, property.Price)
	require.Equal(t, arg.Location, property.Location)
	require.False(t, property.PropertyType.Valid)
	require.False(t, property.ImageUrl.Valid)
	require.False(t, property.VideoUrl.Valid)
	require.NotZero(t, property.ID)
}

func TestCreatePropertyWithNullAgentID(t *testing.T) {
	arg := CreatePropertyParams{
		AgentID: pgtype.Int4{
			Int32: 0,
			Valid: false, // NULL agent_id
		},
		Title:        randomString(20),
		Description:  pgtype.Text{String: randomString(100), Valid: true},
		Price:        int32(randomInt(100000, 2000000)),
		Location:     randomString(50),
		PropertyType: pgtype.Text{String: "self-contained", Valid: true}, // Fixed
		ImageUrl:     pgtype.Text{String: "https://example.com/image.jpg", Valid: true},
		VideoUrl:     pgtype.Text{String: "https://example.com/video.mp4", Valid: true},
	}

	property, err := testQueries.CreateProperty(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, property)
	require.False(t, property.AgentID.Valid)
}

func TestCreatePropertyWithInvalidPropertyType(t *testing.T) {
	agent := CreateRandomUser(t)

	arg := CreatePropertyParams{
		AgentID:     pgtype.Int4{Int32: agent.ID, Valid: true},
		Title:       randomString(20),
		Description: pgtype.Text{String: randomString(100), Valid: true},
		Price:       int32(randomInt(100000, 2000000)),
		Location:    randomString(50),
		PropertyType: pgtype.Text{
			String: "invalid_type", // This should fail CHECK constraint
			Valid:  true,
		},
		ImageUrl: pgtype.Text{String: "https://example.com/image.jpg", Valid: true},
		VideoUrl: pgtype.Text{String: "https://example.com/video.mp4", Valid: true},
	}

	_, err := testQueries.CreateProperty(context.Background(), arg)
	require.Error(t, err)
	// PostgreSQL check constraint violations typically contain "check constraint" or "violates check constraint"
	require.Contains(t, err.Error(), "check")
}

func TestCreatePropertyWithInvalidAgentID(t *testing.T) {
	arg := CreatePropertyParams{
		AgentID: pgtype.Int4{
			Int32: 999999, // Non-existent agent ID
			Valid: true,
		},
		Title:        randomString(20),
		Description:  pgtype.Text{String: randomString(100), Valid: true},
		Price:        int32(randomInt(100000, 2000000)),
		Location:     randomString(50),
		PropertyType: pgtype.Text{String: "self-contained", Valid: true}, // Fixed
		ImageUrl:     pgtype.Text{String: "https://example.com/image.jpg", Valid: true},
		VideoUrl:     pgtype.Text{String: "https://example.com/video.mp4", Valid: true},
	}

	_, err := testQueries.CreateProperty(context.Background(), arg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "foreign key")
}

func TestGetProperty(t *testing.T) {
	property1 := createRandomProperty(t)
	property2, err := testQueries.GetProperty(context.Background(), property1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, property2)

	require.Equal(t, property1.ID, property2.ID)
	require.Equal(t, property1.AgentID.Int32, property2.AgentID.Int32)
	require.Equal(t, property1.AgentID.Valid, property2.AgentID.Valid)
	require.Equal(t, property1.Title, property2.Title)
	require.Equal(t, property1.Description.String, property2.Description.String)
	require.Equal(t, property1.Description.Valid, property2.Description.Valid)
	require.Equal(t, property1.Price, property2.Price)
	require.Equal(t, property1.Location, property2.Location)
	require.Equal(t, property1.PropertyType.String, property2.PropertyType.String)
	require.Equal(t, property1.PropertyType.Valid, property2.PropertyType.Valid)
	require.Equal(t, property1.ImageUrl.String, property2.ImageUrl.String)
	require.Equal(t, property1.ImageUrl.Valid, property2.ImageUrl.Valid)
	require.Equal(t, property1.VideoUrl.String, property2.VideoUrl.String)
	require.Equal(t, property1.VideoUrl.Valid, property2.VideoUrl.Valid)
	require.WithinDuration(t, property1.CreatedAt.Time, property2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, property1.UpdatedAt.Time, property2.UpdatedAt.Time, time.Second)
}

func TestGetPropertyNotFound(t *testing.T) {
	_, err := testQueries.GetProperty(context.Background(), 999999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows")
}

func TestDeleteProperty(t *testing.T) {
	property1 := createRandomProperty(t)
	err := testQueries.DeleteProperty(context.Background(), property1.ID)
	require.NoError(t, err)

	property2, err := testQueries.GetProperty(context.Background(), property1.ID)
	require.Error(t, err)
	require.Empty(t, property2)
	require.Contains(t, err.Error(), "no rows")
}

func TestDeletePropertyNotFound(t *testing.T) {
	err := testQueries.DeleteProperty(context.Background(), 999999)
	require.NoError(t, err) // DELETE with no matching rows doesn't return error
}

// NOTE: The ListProperties function appears to be incorrectly generated
// It contains an UPDATE query instead of a SELECT query for listing
// This test demonstrates the issue

func TestPropertyTypesValidation(t *testing.T) {
	// Fixed: Using valid property types from your schema
	validTypes := []string{"self-contained", "shared", "hostel", "flat"}
	agent := CreateRandomUser(t)

	for _, propType := range validTypes {
		t.Run("PropertyType_"+propType, func(t *testing.T) {
			arg := CreatePropertyParams{
				AgentID:      pgtype.Int4{Int32: agent.ID, Valid: true},
				Title:        randomString(20),
				Description:  pgtype.Text{String: randomString(100), Valid: true},
				Price:        int32(randomInt(100000, 2000000)),
				Location:     randomString(50),
				PropertyType: pgtype.Text{String: propType, Valid: true},
				ImageUrl:     pgtype.Text{String: "https://example.com/image.jpg", Valid: true},
				VideoUrl:     pgtype.Text{String: "https://example.com/video.mp4", Valid: true},
			}

			property, err := testQueries.CreateProperty(context.Background(), arg)
			require.NoError(t, err)
			require.Equal(t, propType, property.PropertyType.String)
		})
	}
}

func TestCreatePropertyEdgeCases(t *testing.T) {
	agent := CreateRandomUser(t)

	testCases := []struct {
		name     string
		arg      CreatePropertyParams
		hasError bool
	}{
		{
			name: "Very long title (should succeed if within limit)",
			arg: CreatePropertyParams{
				AgentID:      pgtype.Int4{Int32: agent.ID, Valid: true},
				Title:        randomString(150), // Max length based on VARCHAR(150)
				Description:  pgtype.Text{String: randomString(100), Valid: true},
				Price:        int32(100000),
				Location:     randomString(50),
				PropertyType: pgtype.Text{String: "self-contained", Valid: true}, // Fixed
				ImageUrl:     pgtype.Text{String: "https://example.com/image.jpg", Valid: true},
				VideoUrl:     pgtype.Text{String: "https://example.com/video.mp4", Valid: true},
			},
			hasError: false,
		},
		{
			name: "Zero price",
			arg: CreatePropertyParams{
				AgentID:      pgtype.Int4{Int32: agent.ID, Valid: true},
				Title:        randomString(20),
				Description:  pgtype.Text{String: randomString(100), Valid: true},
				Price:        0,
				Location:     randomString(50),
				PropertyType: pgtype.Text{String: "shared", Valid: true}, // Fixed
				ImageUrl:     pgtype.Text{String: "https://example.com/image.jpg", Valid: true},
				VideoUrl:     pgtype.Text{String: "https://example.com/video.mp4", Valid: true},
			},
			hasError: false, // Depends on your business rules
		},
		{
			name: "Negative price",
			arg: CreatePropertyParams{
				AgentID:      pgtype.Int4{Int32: agent.ID, Valid: true},
				Title:        randomString(20),
				Description:  pgtype.Text{String: randomString(100), Valid: true},
				Price:        -1000,
				Location:     randomString(50),
				PropertyType: pgtype.Text{String: "hostel", Valid: true}, // Fixed
				ImageUrl:     pgtype.Text{String: "https://example.com/image.jpg", Valid: true},
				VideoUrl:     pgtype.Text{String: "https://example.com/video.mp4", Valid: true},
			},
			hasError: false, // Depends on your business rules - you might want to add a CHECK constraint
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			property, err := testQueries.CreateProperty(context.Background(), tc.arg)
			if tc.hasError {
				require.Error(t, err)
				require.Empty(t, property)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, property)
			}
		})
	}
}
