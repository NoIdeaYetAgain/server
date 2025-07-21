package api

import (
	"bytes"
	"clove/internal/db/sqlc"
	"clove/util"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// mockStore implements the minimal db.SQLStore interface needed for testing
type mockStore struct {
	createUserFunc func(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	getUserFunc    func(ctx context.Context, id int32) (db.User, error)
}

func (m *mockStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, arg)
	}
	return db.User{}, errors.New("createUserFunc not implemented")
}

func (m *mockStore) GetUser(ctx context.Context, id int32) (db.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, id)
	}
	return db.User{}, errors.New("getUserFunc not implemented")
}

// Implement other required SQLStore methods with empty implementations
func (m *mockStore) BeginTx(ctx context.Context) (db.SQLStore, error) {
	return m, nil
}

func (m *mockStore) Commit() error {
	return nil
}

func (m *mockStore) Rollback() error {
	return nil
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStore    func() *mockStore
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"full_name":    user.FullName,
				"email":        user.Email,
				"phone_number": user.PhoneNumber,
				"password":     password,
				"role":         user.Role,
			},
			buildStore: func() *mockStore {
				return &mockStore{
					createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
						require.Equal(t, user.FullName, arg.FullName)
						require.Equal(t, user.Email, arg.Email)
						require.Equal(t, user.PhoneNumber, arg.PhoneNumber)
						require.Equal(t, user.Role, arg.Role)
						require.NotEmpty(t, arg.PasswordHash)
						return user, nil
					},
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"full_name":    user.FullName,
				"email":        "invalid-email",
				"phone_number": user.PhoneNumber,
				"password":     password,
				"role":         user.Role,
			},
			buildStore: func() *mockStore {
				return &mockStore{
					createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
						t.Error("createUser should not be called for invalid email")
						return db.User{}, nil
					},
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ShortPassword",
			body: gin.H{
				"full_name":    user.FullName,
				"email":        user.Email,
				"phone_number": user.PhoneNumber,
				"password":     "123",
				"role":         user.Role,
			},
			buildStore: func() *mockStore {
				return &mockStore{
					createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
						t.Error("createUser should not be called for short password")
						return db.User{}, nil
					},
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"full_name":    user.FullName,
				"email":        user.Email,
				"phone_number": user.PhoneNumber,
				"password":     password,
				"role":         user.Role,
			},
			buildStore: func() *mockStore {
				return &mockStore{
					createUserFunc: func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
						return db.User{}, errors.New("internal server error")
					},
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := tc.buildStore()
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name          string
		userID        string
		buildStore    func() *mockStore
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: strconv.Itoa(int(user.ID)),
			buildStore: func() *mockStore {
				return &mockStore{
					getUserFunc: func(ctx context.Context, id int32) (db.User, error) {
						require.Equal(t, user.ID, id)
						return user, nil
					},
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userID: strconv.Itoa(int(user.ID)),
			buildStore: func() *mockStore {
				return &mockStore{
					getUserFunc: func(ctx context.Context, id int32) (db.User, error) {
						return db.User{}, pgx.ErrNoRows
					},
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			userID: "0",
			buildStore: func() *mockStore {
				return &mockStore{
					getUserFunc: func(ctx context.Context, id int32) (db.User, error) {
						t.Error("getUser should not be called for invalid ID")
						return db.User{}, nil
					},
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userID: strconv.Itoa(int(user.ID)),
			buildStore: func() *mockStore {
				return &mockStore{
					getUserFunc: func(ctx context.Context, id int32) (db.User, error) {
						return db.User{}, errors.New("internal server error")
					},
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := tc.buildStore()
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/users/" + tc.userID
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID:           int32(util.RandomInt(1, 1000)),
		FullName:     util.RandomString(6),
		Email:        util.RandomGmail(),
		PhoneNumber:  util.RandomPhoneNumber(),
		PasswordHash: hashedPassword,
		Role:         "user",
		CreatedAt:    pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
	}
	return user, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.PhoneNumber, gotUser.PhoneNumber)
	require.Equal(t, user.Role, gotUser.Role)
	require.Empty(t, gotUser.PasswordHash)
}

func newTestServer(t *testing.T, store db.SQLStore) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}
