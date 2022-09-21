package main

import (
	"currency-converter/model"
	"currency-converter/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var con = service.Currency{}

func init() {
	con.Server = "mongodb://localhost:27017/"
	con.Database = "currencyConversion"
	con.Collection = "exchangeRate"

	con.Connect()
}

func addCurrencyExchange(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var currency model.Currency

	if err := json.NewDecoder(r.Body).Decode(&currency); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if err := con.Insert(currency); err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable To Insert Record")
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": " Record Inserted Successfully",
		})
	}
}

func convertCurrency(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var convert model.Converter

	if err := json.NewDecoder(r.Body).Decode(&convert); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	if convert.Amount == 0 {
		respondWithError(w, http.StatusBadRequest, "Amount provided for conversion should be greater than 0 ")
		return
	}
	fmt.Println(convert)
	if docs, err := con.Convert(convert); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": "Amount In USD : " + fmt.Sprintf("%v", docs),
		})
	}
}
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func main() {
	port := ":8080"
	http.HandleFunc("/add-currency", addCurrencyExchange)
	http.HandleFunc("/convert-currency", convertCurrency)
	fmt.Printf("listening on port %s... \n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
