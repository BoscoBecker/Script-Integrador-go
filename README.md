Script integrador + Firebird

## Versão do Go go version go1.22.1 windows/amd64

* ##Instalando dependencias
```` go
  go get github.com/nakagami/firebirdsql
````
* Alterar String de conexao
```` go
connector := NewFirebirdDBConnector("localhost", "D:\\Projetos\\Python\\Integrador\\src\\Data\\DADOS.FDB", "SYSDBA", "masterkey")
````
* Descompactar base

* Executar arquivo integrador.go
```` go
 go run integrador.go

````
