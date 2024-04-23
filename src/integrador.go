package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/nakagami/firebirdsql"
)

// JSONRequester faz solicitações HTTP e retorna JSON
type JSONRequester struct{}

// RequestJSON faz uma solicitação GET e retorna o JSON da resposta
func (jr *JSONRequester) RequestJSON(uri string) (map[string]interface{}, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer a requisição: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta: %v", err)
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	return jsonData, nil
}

// FirebirdDBConnector faz conexão com o banco de dados Firebird
type FirebirdDBConnector struct {
	host     string
	database string
	user     string
	password string
	conn     *sql.DB
}

// NewFirebirdDBConnector cria uma nova instância de FirebirdDBConnector
func NewFirebirdDBConnector(host, database, user, password string) *FirebirdDBConnector {
	return &FirebirdDBConnector{
		host:     host,
		database: database,
		user:     user,
		password: password,
	}
}

// Connect conecta ao banco de dados Firebird
func (dbc *FirebirdDBConnector) Connect() error {
	//sysdba:masterkey@127.0.0.1:3050/D:/caminho/para/o/arquivo.fdb
	path := dbc.user + ":" + dbc.password + "@" + dbc.host + ":3050/" + dbc.database

	conn, _ := sql.Open("firebirdsql", path)
	dbc.conn = conn
	fmt.Println("Conexão bem-sucedida com o banco de dados Firebird.")
	return nil
}

// Disconnect desconecta do banco de dados Firebird
func (dbc *FirebirdDBConnector) Disconnect() {
	if dbc.conn != nil {
		dbc.conn.Close()
		fmt.Println("Conexão com o banco de dados Firebird fechada.")
	}
}

// InsertData insere dados no banco de dados Firebird
func (dbc *FirebirdDBConnector) InsertData(data map[string]interface{}) error {
	nome := data["pessoa"].(map[string]interface{})["nome"].(string)
	estado := data["pessoa"].(map[string]interface{})["Estado"].(string)
	cidade := data["pessoa"].(map[string]interface{})["cidade"].(string)
	pais := data["pessoa"].(map[string]interface{})["pais"].(string)
	sistema := "Go"

	sqlCommand := fmt.Sprintf("INSERT INTO USUARIO (NOME, CIDADE, ESTADO, PAIS, SISTEMA) VALUES ('%s', '%s', '%s', '%s','%s')", nome, cidade, estado, pais, sistema)

	_, err := dbc.conn.Exec(sqlCommand)
	if err != nil {
		return fmt.Errorf("erro ao inserir dados: %v", err)
	}

	fmt.Println("Dados inseridos com sucesso.")
	return nil
}

func main() {
	requester := &JSONRequester{}
	uri := "https://gist.githubusercontent.com/BoscoBecker/b343b480631ca61b0b06f4dca6b23139/raw/440f560f86627871789eabdc4c86b4e819ddc9b1/data.json"
	jsonData, err := requester.RequestJSON(uri)
	if err != nil {
		log.Fatalf("Erro ao obter JSON: %v", err)
	}

	connector := NewFirebirdDBConnector("localhost", "D:\\Projetos\\Python\\Integrador\\src\\Data\\DADOS.FDB", "SYSDBA", "masterkey")
	err = connector.Connect()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer connector.Disconnect()

	err = connector.InsertData(jsonData)
	if err != nil {
		log.Fatalf("Erro ao inserir dados: %v", err)
	}
}
