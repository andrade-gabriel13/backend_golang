// config/ConnDB_test.go
package tests

import (
	"fmt"
	"testing"

	// Importa o pacote 'config' para acessar a função ConnectDB e a variável DB
	"myapi/config"
)

// Testa a conexão com o banco de dados MySQL
func TestDBConnection(t *testing.T) {
	// Chama a função de conexão com o banco de dados do pacote config
	config.ConnectDB()

	// Verifica se a variável global DB foi configurada corretamente
	if config.DB == nil {
		t.Fatalf("A conexão com o banco de dados não foi estabelecida.")
	}

	// Verifica se a conexão pode ser "pingada" (confirmando se está ativa)
	err := config.DB.Exec("SELECT 1").Error
	if err != nil {
		t.Fatalf("Erro ao tentar fazer ping no banco de dados: %v", err)
	}

	fmt.Println("Conexão com o banco de dados foi bem-sucedida!")
}
