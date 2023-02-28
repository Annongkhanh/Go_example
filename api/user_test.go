package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Annongkhanh/Go_example/db/mock"
	db "github.com/Annongkhanh/Go_example/db/sqlc"
	"github.com/Annongkhanh/Go_example/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetUserAPI(t *testing.T) {

	user, password := randomUser()
	fmt.Printf("random password:%s", password)

	testCases := []struct {
		name          string
		username      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(
					user.Username,
				)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user, password)
			},
		},
		{
			name:     "Not found",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(
					user.Username,
				)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "Internal server error",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(
					user.Username,
				)).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "Invalid username",
			username: " ",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(
					" ",
				)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			//start test server and send request

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", tc.username)

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

			// respDump, err := httputil.DumpResponse(recorder.Result(), true)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// fmt.Printf("RESPONSE:\n%s", string(respDump))

		})

	}

}

func randomUser() (db.User, string) {
	password := util.RandomString(int(util.RandomInt(32, 6)))
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	return db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Fullname:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User, password string) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	require.NoError(t, util.CheckPassword(password, user.HashedPassword))
	user.HashedPassword = ""
	//unmarshal response body
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user, gotUser)
}

// func TestCreateUser(t *testing.T) {

// 	testCases := []struct {
// 		name          string
// 		requestBody   []byte
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:        "OK",
// 			requestBody: []byte(`{"username":"testuser","password":"testpass","fullname":"Test User","email":"testuser@example.com"}`),
// 			buildStubs: func(store *mockdb.MockStore) {
// 				hashedPassword, err := util.HashPassword("testpass")
// 				require.NoError(t, err)
// 				store.EXPECT().CreateUser(gomock.Any(), db.CreateUserParams{
// 					Username:       "testuser",
// 					HashedPassword: hashedPassword,
// 					Fullname:       "Test User",
// 					Email:          "testuser@example.com",
// 				}).Times(1).Return(db.User{
// 					Username:       "testuser",
// 					HashedPassword: hashedPassword,
// 					Fullname:       "Test User",
// 					Email:          "testuser@example.com",
// 				}, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchUser(t, recorder.Body, db.User{
// 					Username:       "testuser",
// 					HashedPassword: "",
// 					Fullname:       "Test User",
// 					Email:          "testuser@example.com",
// 				}, "")
// 			},
// 		},
// 		{
// 			name:        "Invalid request body",
// 			requestBody: []byte(`invalid json`),
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name:        "Missing required fields",
// 			requestBody: []byte(`{"username":"testuser","password":"testpass"}`),
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name:        "Username already exists",
// 			requestBody: []byte(`{"username":"testuser","password":"testpass","fullname":"Test User","email":"testuser@example.com"}`),
// 			buildStubs: func(store *mockdb.MockStore) {
// 				hashedPassword, err := util.HashPassword("testpass")
// 				require.NoError(t, err)
// 				store.EXPECT().CreateUser(gomock.Any(), db.CreateUserParams{
// 					Username:       "testuser",
// 					HashedPassword: hashedPassword,
// 					Fullname:       "Test User",
// 					Email:          "testuser@example.com",
// 				}).Times(1).Return(db.User{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusConflict, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)

// 			tc.buildStubs(store)

// 			//start test server and send request

// 			server := newTestServer(t, store)

// 			recorder := httptest.NewRecorder()

// 			url := fmt.Sprintf("/users/")

// 			request, err := http.NewRequest(http.MethodGet, url, nil)

// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)

// 			tc.checkResponse(t, recorder)

// 			// respDump, err := httputil.DumpResponse(recorder.Result(), true)
// 			// if err != nil {
// 			// 	log.Fatal(err)
// 			// }

// 			// fmt.Printf("RESPONSE:\n%s", string(respDump))

// 		})

// 	}
// }
