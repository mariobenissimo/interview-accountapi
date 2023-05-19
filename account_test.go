package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var AccountID string

func TestCreateAccount(t *testing.T) {
	//Create an account
	country := "GB"

	accountAttrTest := AccountAttributes{
		BankID:        "GBDSC",
		Bic:           "NWBKGB22",
		Country:       &country,
		Iban:          "GB11NWBK40030041426819",
		AccountNumber: "41426819",
		Name:          []string{"Mario", "Benissimo"},
	}

	// Serializza i dati come JSON
	payload, err := json.Marshal(accountAttrTest)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/createAccount/", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	createAccount(res, req)
	var data map[string]interface{}
	// fmt.Println(res.Body.String())
	err = json.Unmarshal(res.Body.Bytes(), &data)

	if err != nil {
		t.Fatal(err)
	}

	if account, ok := data["links"].(map[string]interface{}); ok {
		if str, ok := account["self"].(string); ok {
			AccountID = strings.Split(str, "/")[4]
		}
	}
	// fmt.Println(AccountID)
	assert.Equal(t, http.StatusOK, res.Code, "Expected status 200")
	assert.Equal(t, len(AccountID), 36)
}

func TestGetAccount(t *testing.T) {
	getAccount(t, false)
}
func getAccount(t *testing.T, delete bool) {
	url := "/getAccount/" + AccountID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	getAccountByAccountID(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "Expected status 200")
	expectedBody := `{"error_message":"record ` + AccountID + ` does not exist"}`
	if delete == false {
		assert.NotEqual(t, expectedBody, res.Body.String(), "Response body mismatch")
	} else {
		assert.Equal(t, expectedBody, res.Body.String(), "Response body mismatch")
	}
}
func TestDeleteAccount(t *testing.T) {

	url := "/deleteAccount/" + AccountID + "?version=0"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	deleteAccountByAccountId(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "Expected status 200")
	getAccount(t, true)
}
