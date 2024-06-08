package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CotacaoService interface {
	TreatedResponseBody()
}

type Dolar struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	Bid string `json:"bid"`
}

type Database struct {
	Database *sql.DB
}

func NewConnectionDB() (*Database, error) {
	database, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		fmt.Println("Erro ao abrir o banco de dados:", err)
		return nil, err
	}

	return &Database{
		Database: database,
	}, nil
}

func BuscaUSDBRLHandler(w http.ResponseWriter, r *http.Request) {
	cotacao := &Dolar{}
	ctx := context.Background()

	resp, err := Cotacao(ctx)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dolar, err := cotacao.TreatedResponseBody(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dolar)
	log.Printf("A cotação do Dólar neste momento é de: $%s", dolar.Bid)
}

func Cotacao(ctx context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("request na URL %s passou do tempo limite", url)
		default:
			return nil, err
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (d *Dolar) TreatedResponseBody(body []byte) (*USDBRL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("inserção da cotação no banco de dados passou do tempo limite")
	default:
		err := json.Unmarshal(body, &d)
		if err != nil {
			return nil, err
		}

		db, err := NewConnectionDB()
		if err != nil {
			return nil, err
		}
		defer db.Database.Close()

		err = createTableSQL(db)
		if err != nil {
			return nil, err
		}

		err = insertCotacao(db, d.USDBRL.Bid)
		if err != nil {
			return nil, err
		}

		return &d.USDBRL, nil
	}
}

// CREATE TABLE IF NOT EXISTS
func createTableSQL(db *Database) error {
	createTableStr := `CREATE TABLE IF NOT EXISTS cotacao_dolar (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"cotacao" NUMBER,
		"processing_date" DATE
	)`
	stmt, err := db.Database.Prepare(createTableStr)
	if err != nil {
		log.Println("Erro ao criar a tabela no SQLite: ", err)
		return err
	}
	stmt.Exec()

	return nil
}

func insertCotacao(db *Database, cotacao string) error {
	insertCotacao := `INSERT INTO cotacao_dolar (cotacao, processing_date) values (?, ?)`
	stmt, err := db.Database.Prepare(insertCotacao)
	if err != nil {
		log.Println("Erro ao preparar a query do SQLite: ", err)
		return err
	}

	currentDate := generateCurrentDate()

	_, err = stmt.Exec(cotacao, currentDate)
	if err != nil {
		log.Println("Erro ao executar a insercao: ", err)
		return err
	}

	log.Println("Cotação inserida com sucesso!!!")

	return nil
}

func generateCurrentDate() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	return formattedTime
}

func main() {
	log.Println("Servidor iniciado com sucesso...")
	http.HandleFunc("/cotacao", BuscaUSDBRLHandler)
	http.ListenAndServe(":8080", nil)
}
