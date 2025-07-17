package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomReview(t *testing.T) Review {
	// First create a student to reference
	student := CreateRandomUser(t) // Assuming you have this helper function

	arg := CreateReviewParams{
		StudentID: student.ID,
		Rating: pgtype.Int4{
			Int32: int32(randomInt(1, 5)), // Rating between 1-5
			Valid: true,
		},
		Comment: pgtype.Text{
			String: randomString(50),
			Valid:  true,
		},
	}

	review, err := testQueries.CreateReview(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, review)

	require.Equal(t, arg.StudentID, review.StudentID)
	require.Equal(t, arg.Rating.Int32, review.Rating.Int32)
	require.Equal(t, arg.Comment.String, review.Comment.String)

	require.NotZero(t, review.ID)
	require.NotZero(t, review.CreatedAt)

	return review
}

func TestCreateReview(t *testing.T) {
	createRandomReview(t)
}

func TestCreateReviewWithNullValues(t *testing.T) {
	// Test creating review with null rating and comment
	student := CreateRandomUser(t)

	arg := CreateReviewParams{
		StudentID: student.ID,
		Rating: pgtype.Int4{
			Int32: 0,
			Valid: false, // NULL rating
		},
		Comment: pgtype.Text{
			String: "",
			Valid:  false, // NULL comment
		},
	}

	review, err := testQueries.CreateReview(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, review)

	require.Equal(t, arg.StudentID, review.StudentID)
	require.False(t, review.Rating.Valid)
	require.False(t, review.Comment.Valid)
	require.NotZero(t, review.ID)
	require.NotZero(t, review.CreatedAt)
}

func TestGetReviewByID(t *testing.T) {
	review1 := createRandomReview(t)
	review2, err := testQueries.GetReviewByID(context.Background(), review1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, review2)

	require.Equal(t, review1.ID, review2.ID)
	require.Equal(t, review1.StudentID, review2.StudentID)
	require.Equal(t, review1.Rating.Int32, review2.Rating.Int32)
	require.Equal(t, review1.Rating.Valid, review2.Rating.Valid)
	require.Equal(t, review1.Comment.String, review2.Comment.String)
	require.Equal(t, review1.Comment.Valid, review2.Comment.Valid)
	require.WithinDuration(t, review1.CreatedAt.Time, review2.CreatedAt.Time, time.Second)
}

func TestGetReviewByIDNotFound(t *testing.T) {
	_, err := testQueries.GetReviewByID(context.Background(), 999999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows")
}

func TestUpdateReview(t *testing.T) {
	oldReview := createRandomReview(t)

	arg := UpdateReviewParams{
		ID: oldReview.ID,
		Rating: pgtype.Int4{
			Int32: int32(randomInt(1, 5)),
			Valid: true,
		},
		Comment: pgtype.Text{
			String: randomString(100),
			Valid:  true,
		},
	}

	updatedReview, err := testQueries.UpdateReview(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedReview)

	require.Equal(t, oldReview.ID, updatedReview.ID)
	require.Equal(t, oldReview.StudentID, updatedReview.StudentID)
	require.Equal(t, arg.Rating.Int32, updatedReview.Rating.Int32)
	require.Equal(t, arg.Comment.String, updatedReview.Comment.String)
	require.True(t, updatedReview.CreatedAt.Time.After(oldReview.CreatedAt.Time))
}

func TestUpdateReviewWithNullValues(t *testing.T) {
	oldReview := createRandomReview(t)

	arg := UpdateReviewParams{
		ID: oldReview.ID,
		Rating: pgtype.Int4{
			Int32: 0,
			Valid: false, // Set rating to NULL
		},
		Comment: pgtype.Text{
			String: "",
			Valid:  false, // Set comment to NULL
		},
	}

	updatedReview, err := testQueries.UpdateReview(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedReview)

	require.Equal(t, oldReview.ID, updatedReview.ID)
	require.Equal(t, oldReview.StudentID, updatedReview.StudentID)
	require.False(t, updatedReview.Rating.Valid)
	require.False(t, updatedReview.Comment.Valid)
}

func TestUpdateReviewNotFound(t *testing.T) {
	arg := UpdateReviewParams{
		ID: 999999,
		Rating: pgtype.Int4{
			Int32: 5,
			Valid: true,
		},
		Comment: pgtype.Text{
			String: "Updated comment",
			Valid:  true,
		},
	}

	_, err := testQueries.UpdateReview(context.Background(), arg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows")
}

func TestDeleteReview(t *testing.T) {
	review1 := createRandomReview(t)
	err := testQueries.DeleteReview(context.Background(), review1.ID)
	require.NoError(t, err)

	review2, err := testQueries.GetReviewByID(context.Background(), review1.ID)
	require.Error(t, err)
	require.Empty(t, review2)
	require.Contains(t, err.Error(), "no rows")
}

func TestDeleteReviewNotFound(t *testing.T) {
	err := testQueries.DeleteReview(context.Background(), 999999)
	require.NoError(t, err) // DELETE with no matching rows doesn't return error
}

// Test foreign key constraint
func TestCreateReviewWithInvalidStudentID(t *testing.T) {
	arg := CreateReviewParams{
		StudentID: 999999, // Non-existent student ID
		Rating: pgtype.Int4{
			Int32: 5,
			Valid: true,
		},
		Comment: pgtype.Text{
			String: "Great review",
			Valid:  true,
		},
	}

	_, err := testQueries.CreateReview(context.Background(), arg)
	require.Error(t, err)
	// Should fail due to foreign key constraint
	require.Contains(t, err.Error(), "foreign key")
}

// Helper functions (you might already have these)
func randomInt(min, max int) int {
	// Implementation depends on your random utility
	// Example: return rand.Intn(max-min+1) + min
	return min + (max-min)/2 // Simple implementation for demo
}

func randomString(n int) string {
	// Implementation depends on your random utility
	// This is just a simple example
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
