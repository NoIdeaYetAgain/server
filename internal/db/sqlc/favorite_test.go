package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomFavorite(t *testing.T) (Property, User) {
	// Create a student user and a property
	student := CreateRandomUser(t)      // Assuming this creates a student user
	property := createRandomProperty(t) // This should create a property

	// Add the property to student's favorites
	arg := AddFavoriteParams{
		StudentID:  student.ID,
		PropertyID: property.ID,
	}

	err := testQueries.AddFavorite(context.Background(), arg)
	require.NoError(t, err)

	return property, student
}

func TestAddFavorite(t *testing.T) {
	student := CreateRandomUser(t)
	property := createRandomProperty(t)

	arg := AddFavoriteParams{
		StudentID:  student.ID,
		PropertyID: property.ID,
	}

	err := testQueries.AddFavorite(context.Background(), arg)
	require.NoError(t, err)

	// Verify the favorite was added by getting favorites
	getFavoritesArg := GetFavoritesApartmentsParams{
		StudentID: student.ID,
		Limit:     10,
		Offset:    0,
	}

	favorites, err := testQueries.GetFavoritesApartments(context.Background(), getFavoritesArg)
	require.NoError(t, err)
	require.Len(t, favorites, 1)
	require.Equal(t, property.ID, favorites[0].ID)
}

func TestAddFavoriteDuplicate(t *testing.T) {
	student := CreateRandomUser(t)
	property := createRandomProperty(t)

	arg := AddFavoriteParams{
		StudentID:  student.ID,
		PropertyID: property.ID,
	}

	// Add favorite first time
	err := testQueries.AddFavorite(context.Background(), arg)
	require.NoError(t, err)

	// Add the same favorite again - should not error due to ON CONFLICT DO NOTHING
	err = testQueries.AddFavorite(context.Background(), arg)
	require.NoError(t, err)

	// Verify only one favorite exists
	getFavoritesArg := GetFavoritesApartmentsParams{
		StudentID: student.ID,
		Limit:     10,
		Offset:    0,
	}

	favorites, err := testQueries.GetFavoritesApartments(context.Background(), getFavoritesArg)
	require.NoError(t, err)
	require.Len(t, favorites, 1)
}

func TestAddFavoriteWithNonExistentStudent(t *testing.T) {
	property := createRandomProperty(t)

	arg := AddFavoriteParams{
		StudentID:  999999, // Non-existent student ID
		PropertyID: property.ID,
	}

	err := testQueries.AddFavorite(context.Background(), arg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "foreign key") // Should fail foreign key constraint
}

func TestAddFavoriteWithNonExistentProperty(t *testing.T) {
	student := CreateRandomUser(t)

	arg := AddFavoriteParams{
		StudentID:  student.ID,
		PropertyID: 999999, // Non-existent property ID
	}

	err := testQueries.AddFavorite(context.Background(), arg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "foreign key") // Should fail foreign key constraint
}

func TestGetFavoritesApartments(t *testing.T) {
	student := CreateRandomUser(t)

	// Create multiple properties and add them to favorites
	properties := make([]Property, 3)
	for i := 0; i < 3; i++ {
		properties[i] = createRandomProperty(t)

		arg := AddFavoriteParams{
			StudentID:  student.ID,
			PropertyID: properties[i].ID,
		}

		err := testQueries.AddFavorite(context.Background(), arg)
		require.NoError(t, err)
	}

	// Get favorites
	arg := GetFavoritesApartmentsParams{
		StudentID: student.ID,
		Limit:     10,
		Offset:    0,
	}

	favorites, err := testQueries.GetFavoritesApartments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, favorites, 3)

	// Verify all properties are returned and ordered by created_at DESC
	for i, favorite := range favorites {
		require.NotEmpty(t, favorite.ID)
		require.NotEmpty(t, favorite.Title)
		require.NotEmpty(t, favorite.Location)
		require.NotZero(t, favorite.Price)

		// Verify it's one of our created properties
		found := false
		for _, prop := range properties {
			if prop.ID == favorite.ID {
				found = true
				break
			}
		}
		require.True(t, found, "Favorite property should be one of the created properties")

		// Check ordering (created_at DESC)
		if i > 0 {
			require.True(t, favorites[i-1].CreatedAt.Time.After(favorite.CreatedAt.Time) ||
				favorites[i-1].CreatedAt.Time.Equal(favorite.CreatedAt.Time))
		}
	}
}

func TestGetFavoritesApartmentsWithPagination(t *testing.T) {
	student := CreateRandomUser(t)

	// Create 5 properties and add them to favorites
	for i := 0; i < 5; i++ {
		property := createRandomProperty(t)

		arg := AddFavoriteParams{
			StudentID:  student.ID,
			PropertyID: property.ID,
		}

		err := testQueries.AddFavorite(context.Background(), arg)
		require.NoError(t, err)
	}

	// Test pagination - get first 3
	arg := GetFavoritesApartmentsParams{
		StudentID: student.ID,
		Limit:     3,
		Offset:    0,
	}

	favorites, err := testQueries.GetFavoritesApartments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, favorites, 3)

	// Get next 2
	arg.Limit = 3
	arg.Offset = 3

	remainingFavorites, err := testQueries.GetFavoritesApartments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, remainingFavorites, 2)

	// Verify no overlap between pages
	for _, fav1 := range favorites {
		for _, fav2 := range remainingFavorites {
			require.NotEqual(t, fav1.ID, fav2.ID)
		}
	}
}

func TestGetFavoritesApartmentsEmptyResult(t *testing.T) {
	student := CreateRandomUser(t)

	arg := GetFavoritesApartmentsParams{
		StudentID: student.ID,
		Limit:     10,
		Offset:    0,
	}

	favorites, err := testQueries.GetFavoritesApartments(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, favorites)
}

func TestGetFavoritesApartmentsWithNonExistentStudent(t *testing.T) {
	arg := GetFavoritesApartmentsParams{
		StudentID: 999999, // Non-existent student
		Limit:     10,
		Offset:    0,
	}

	favorites, err := testQueries.GetFavoritesApartments(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, favorites) // Should return empty slice, not error
}

func TestRemoveFavorite(t *testing.T) {
	property, student := createRandomFavorite(t)

	// Verify favorite exists
	getFavoritesArg := GetFavoritesApartmentsParams{
		StudentID: student.ID,
		Limit:     10,
		Offset:    0,
	}

	favorites, err := testQueries.GetFavoritesApartments(context.Background(), getFavoritesArg)
	require.NoError(t, err)
	require.Len(t, favorites, 1)

	// Remove the favorite
	removeArg := RemoveFavoriteParams{
		StudentID:  student.ID,
		PropertyID: property.ID,
	}

	err = testQueries.RemoveFavorite(context.Background(), removeArg)
	require.NoError(t, err)

	// Verify favorite was removed
	favorites, err = testQueries.GetFavoritesApartments(context.Background(), getFavoritesArg)
	require.NoError(t, err)
	require.Empty(t, favorites)
}

func TestRemoveFavoriteNotExists(t *testing.T) {
	student := CreateRandomUser(t)
	property := createRandomProperty(t)

	// Try to remove a favorite that doesn't exist
	arg := RemoveFavoriteParams{
		StudentID:  student.ID,
		PropertyID: property.ID,
	}

	err := testQueries.RemoveFavorite(context.Background(), arg)
	require.NoError(t, err) // DELETE with no matching rows doesn't return error
}

func TestRemoveFavoriteWithNonExistentStudent(t *testing.T) {
	property := createRandomProperty(t)

	arg := RemoveFavoriteParams{
		StudentID:  999999, // Non-existent student
		PropertyID: property.ID,
	}

	err := testQueries.RemoveFavorite(context.Background(), arg)
	require.NoError(t, err) // DELETE with no matching rows doesn't return error
}

func TestRemoveFavoriteWithNonExistentProperty(t *testing.T) {
	student := CreateRandomUser(t)

	arg := RemoveFavoriteParams{
		StudentID:  student.ID,
		PropertyID: 999999, // Non-existent property
	}

	err := testQueries.RemoveFavorite(context.Background(), arg)
	require.NoError(t, err) // DELETE with no matching rows doesn't return error
}

func TestFavoriteWorkflow(t *testing.T) {
	student := CreateRandomUser(t)
	property1 := createRandomProperty(t)
	property2 := createRandomProperty(t)

	// Add two favorites
	err := testQueries.AddFavorite(context.Background(), AddFavoriteParams{
		StudentID:  student.ID,
		PropertyID: property1.ID,
	})
	require.NoError(t, err)

	err = testQueries.AddFavorite(context.Background(), AddFavoriteParams{
		StudentID:  student.ID,
		PropertyID: property2.ID,
	})
	require.NoError(t, err)

	// Verify both favorites exist
	favorites, err := testQueries.GetFavoritesApartments(context.Background(), GetFavoritesApartmentsParams{
		StudentID: student.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.Len(t, favorites, 2)

	// Remove one favorite
	err = testQueries.RemoveFavorite(context.Background(), RemoveFavoriteParams{
		StudentID:  student.ID,
		PropertyID: property1.ID,
	})
	require.NoError(t, err)

	// Verify only one favorite remains
	favorites, err = testQueries.GetFavoritesApartments(context.Background(), GetFavoritesApartmentsParams{
		StudentID: student.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.Len(t, favorites, 1)
	require.Equal(t, property2.ID, favorites[0].ID)
}

func TestMultipleStudentsFavorites(t *testing.T) {
	student1 := CreateRandomUser(t)
	student2 := CreateRandomUser(t)
	property := createRandomProperty(t)

	// Both students add the same property to favorites
	err := testQueries.AddFavorite(context.Background(), AddFavoriteParams{
		StudentID:  student1.ID,
		PropertyID: property.ID,
	})
	require.NoError(t, err)

	err = testQueries.AddFavorite(context.Background(), AddFavoriteParams{
		StudentID:  student2.ID,
		PropertyID: property.ID,
	})
	require.NoError(t, err)

	// Verify each student has their own favorite
	favorites1, err := testQueries.GetFavoritesApartments(context.Background(), GetFavoritesApartmentsParams{
		StudentID: student1.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.Len(t, favorites1, 1)

	favorites2, err := testQueries.GetFavoritesApartments(context.Background(), GetFavoritesApartmentsParams{
		StudentID: student2.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.Len(t, favorites2, 1)

	// Remove student1's favorite
	err = testQueries.RemoveFavorite(context.Background(), RemoveFavoriteParams{
		StudentID:  student1.ID,
		PropertyID: property.ID,
	})
	require.NoError(t, err)

	// Verify student1 has no favorites, but student2 still has theirs
	favorites1, err = testQueries.GetFavoritesApartments(context.Background(), GetFavoritesApartmentsParams{
		StudentID: student1.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.Empty(t, favorites1)

	favorites2, err = testQueries.GetFavoritesApartments(context.Background(), GetFavoritesApartmentsParams{
		StudentID: student2.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.Len(t, favorites2, 1)
}
