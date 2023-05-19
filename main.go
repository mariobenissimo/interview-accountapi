package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/createAccount", createAccount)
	http.HandleFunc("/getAccount/", getAccountByAccountID)
	http.HandleFunc("/deleteAccount/", deleteAccountByAccountId)
	// Start server port 80888
	log.Println("Server listening on :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func deleteAccountByAccountId(w http.ResponseWriter, r *http.Request) {
	// check if the method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// get accountId from url
	accountID := r.URL.Path[len("/deleteAccount/"):]

	url := "http://localhost:8080/v1/organisation/accounts/" + accountID + "?version=0"
	body := makeRequest(url, "DELETE", nil)
	w.Write(body)
}

func makeRequest(url string, method string, payload []byte) []byte {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Esegui la richiesta
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Leggi il corpo della risposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func getAccountByAccountID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	accountID := r.URL.Path[len("/getAccount/"):]
	url := "http://localhost:8080/v1/organisation/accounts/" + accountID
	body := makeRequest(url, "GET", nil)
	w.Write(body)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	// Verifica il metodo della richiesta che appunto deve essere di tipo Post
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body into a struct
	var accountAttributes AccountAttributes
	err := json.NewDecoder(r.Body).Decode(&accountAttributes)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	url := "http://localhost:8080/v1/organisation/accounts"

	organisation_id := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	uuid := uuid.New()

	accountData := AccountData{
		ID:             uuid.String(),
		OrganisationID: organisation_id,
		Type:           "accounts",
		Attributes:     &accountAttributes,
	}

	data := map[string]interface{}{
		"data": accountData,
	}

	// Serializza i dati come JSON
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	body := makeRequest(url, "POST", payload)
	w.Write(body)

}
