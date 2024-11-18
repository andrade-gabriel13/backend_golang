package config

import (
	"fmt"
	"log"
	"myapi/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB realiza a conexão com o banco de dados MySQL e configura a migração automática do modelo `Client`.
//
// Esta função é responsável por estabelecer a conexão com o banco de dados MySQL utilizando as credenciais e
// configurações definidas na string DSN (Data Source Name). Se a conexão for bem-sucedida, ela realiza a migração
// automática do modelo `Client` para garantir que a tabela esteja atualizada conforme o modelo. Caso contrário,
// a função registra um erro detalhado e interrompe a execução do programa.
//
// Exemplo de uso:
//
//	config.ConnectDB()
func ConnectDB() {
	// String de conexão com o banco de dados MySQL.
	dsn := "golang:golang@tcp(127.0.0.1:3306)/golang?charset=utf8&parseTime=True&loc=Local"

	// Tentativa de abrir a conexão com o banco de dados.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Caso ocorra um erro na conexão, loga o erro e encerra a execução do programa.
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Realiza a migração automática das tabelas `Client` e `ArchivedClient` para o banco de dados.
	if err := db.AutoMigrate(&models.Client{}, &models.ArchivedClient{}); err != nil {
		// Caso ocorra um erro durante a migração, loga o erro e encerra a execução do programa.
		log.Fatalf("Erro ao migrar os modelos: %v", err)
	}

	// Atribui a conexão bem-sucedida ao banco de dados à variável global `DB`.
	DB = db

	// Log de sucesso indicando que a conexão foi bem-sucedida.
	fmt.Println("Conectado com sucesso ao MySQL!")
}
