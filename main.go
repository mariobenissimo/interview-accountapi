package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

func main() {
	http.HandleFunc("/createAccount", createAccount)
	http.HandleFunc("/getAccount", getAccountByAccountID)
	http.HandleFunc("/deleteAccount/", deleteAccountByAccountId)
	// Avvia il server sulla porta 8080
	log.Println("Server listening on :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func deleteAccountByAccountId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	accountID := r.URL.Path[len("/deleteAccount/"):]
	// Access individual query parameters
	url := "http://localhost:8080/v1/organisation/accounts/" + accountID + "?version=0"
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Verifica lo stato della risposta
	if response.StatusCode == http.StatusNoContent {
		fmt.Println("Risorsa eliminata con successo")
	} else {
		fmt.Println("Errore durante l'eliminazione della risorsa")
	}

}

// /v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea8cc4dc
func getAccountByAccountID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	queryParams := r.URL.Query()

	// Access individual query parameters
	accountID := queryParams.Get("accountID")

	url := "http://localhost:8080/v1/organisation/accounts/" + accountID

	// Send the GET request
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(body)

}

// per la registrazione l'utente deve fornire l'account_number e iban all'altrimenti viene creato un nuovo account
func createAccount(w http.ResponseWriter, r *http.Request) {
	// Verifica il metodo della richiesta che appunto deve essere di tipo Post
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Analizza i dati del corpo della richiesta per vedere se effettivamente ci sono errori di sintasssi
	// Decode the JSON body into a struct
	var accountAttributes AccountAttributes
	err := json.NewDecoder(r.Body).Decode(&accountAttributes)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	url := "http://localhost:8080/v1/organisation/accounts"
	method := "POST"

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

	// Crea una richiesta POST con il payload JSON
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

	// Stampa la risposta
	fmt.Println(string(body))
	w.Write(body)

}
