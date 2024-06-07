package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type USDBRL struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Println("tempo excedito no client.go")
	default:
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

		cotacao, err := TreatedResponseBody(res)
		if err != nil {
			return
		}

		filename := "cotacao.txt"
		file, err := CreateFile(filename)
		if err != nil {
			return
		}
		defer file.Close()

		err = WriteFile(file, filename, cotacao)
		if err != nil {
			return
		}
	}
}

func TreatedResponseBody(res []byte) (*USDBRL, error) {
	var cotacao USDBRL
	err := json.Unmarshal(res, &cotacao)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer o parse da resposta: %v\n", err)
		return nil, err
	}

	if cotacao.Bid == "" {
		fmt.Fprintln(os.Stderr, "A cotação está vazia, entre em contato com o suporte.")
		return nil, err
	}

	return &cotacao, nil
}

func CreateFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar o arquivo: %v\n", err)
		return nil, err
	}

	return file, nil
}

func WriteFile(file *os.File, filename string, cotacao *USDBRL) error {
	_, err := file.WriteString(fmt.Sprintf("Dólar: %s", cotacao.Bid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo: %v\n", err)
		return err
	}

	log.Printf("Escrita realizada com sucesso, verifique a cotação do Dolár no arquivo %s.", filename)

	return nil
}
