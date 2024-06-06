package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type USDBRL struct {
	Bid string `json:"bid"`
}

// Inserir os contexts
func main() {
	req, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição: %v\n", err)
		return
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta:: %v\n", err)
		return
	}

	var cotacao USDBRL
	err = json.Unmarshal(res, &cotacao)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer o parse da resposta: %v\n", err)
		return
	}

	if cotacao.Bid == "" {
		fmt.Fprintln(os.Stderr, "A cotação está vazia, entre em contato com o suporte.")
		return
	}

	filename := "cotacao.txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar o arquivo: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", cotacao.Bid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo: %v\n", err)
		return
	}

	log.Printf("Escrita realizada com sucesso, verifique a cotação do Dolár no arquivo %s.", filename)
}
