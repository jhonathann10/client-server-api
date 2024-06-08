# client-server-api
Desafio da pós gradução na full cycle (client-server-api).

## Requisitos
O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.

## Instruções para execução
- Inicializar o `./server/server.go`. 
  - Execute:
    - go run server/server.go
- Em outro terminal, inicialize o `./client/client.go`.
  - Execute:
    - go run client/client.go

## Regras
- O arquivo server.go é responsável por inicializar uma API do tipo HTTP com o metódo GET, com o objetivo de retornar a cotação do dolár nesse exato momento se baseando nas informações do site `https://economia.awesomeapi.com.br/json/last/USD-BRL`.
- Os dados serão salvos em um banco de dados do tipo SQLite3, contendo a data de processamento da requisição.
  - Após a primeira execução será criado um arquivo chamado `cotacao.db`, que será o banco de dados.
  - Também será criado uma tabela caso ela não exista, com o nome de `cotacao_dolar`.
- Exemplo de requisição manual:
```curl
curl  -X GET \
  'http://localhost:8080/cotacao' \
  --header 'Content-Type: application/json'
```
- O arquivo client.go é responsável por fazer o papel do cliente que irá fazer as requisições na API do server.go.
- Na primeira execução, será responsável por criar um arquivo chamado `cotacao.txt` e escrever o valor da cotação do Dólar. 

## Banco de dados
Para visualizar os dados no banco de dados, podemos seguir as seguintes instruções pelo próprio terminal:
1. `sqlite3 cotacao.db`.
2. `select * from cotacao_dolar;`.
