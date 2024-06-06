package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type CotacaoService interface {
	Cotacao()
}

type Dolar struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	Bid string `json:"bid"`
}

func BuscaUSDBRLHandler(w http.ResponseWriter, r *http.Request) {
	cotacao := Dolar{}

	resp, err := cotacao.Cotacao()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	log.Printf("A cotação do Dólar neste momento é de: $%s", resp.Bid)
}

func (dr *Dolar) Cotacao() (*USDBRL, error) {
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &dr)
	if err != nil {
		return nil, err
	}

	return &dr.USDBRL, nil
}

func main() {
	log.Println("Servidor iniciado com sucesso...")
	http.HandleFunc("/cotacao", BuscaUSDBRLHandler)
	http.ListenAndServe(":8080", nil)
}
