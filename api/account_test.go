package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Annongkhanh/Go_example/db/mock"
	db "github.com/Annongkhanh/Go_example/db/sqlc"
	"github.com/Annongkhanh/Go_example/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T){
	account :=  randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
 
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	//start test server and send request
	
	server := NewServer(store)

	recorder := httptest.NewRecorder()

	uri := fmt.Sprintf("/accounts/%d", account.ID)

	request, err :=  http.NewRequest(http.MethodGet, uri, nil)

	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

	account = db.Account{}

	requireBodyMatchAccount(t, recorder.Body,account)

}

func randomAccount() db.Account{
	return db.Account{
		ID: util.RandomInt(1000, 1),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: db.Currency(util.RandomCurrency()),
	}
}

func requireBodyMatchAccount(t *testing.T,body *bytes.Buffer, account db.Account){
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	//unmarshal response body
	err = json.Unmarshal(data, &gotAccount) 
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}