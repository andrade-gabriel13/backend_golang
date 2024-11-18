package tests

import (
	"log"
	"myapi/config"
	"myapi/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateClient(t *testing.T) {
	// 1. Conectar ao banco
	config.ConnectDB()

	// 2. Criar um cliente de teste
	client := models.Client{
		Name:     "Cliente Teste",
		WeightKg: 70,
		Address:  "Rua Teste, 123",
	}

	// Inserir o cliente no banco
	if err := config.DB.Create(&client).Error; err != nil {
		log.Fatalf("Erro ao inserir cliente no banco: %v", err)
	}
	t.Logf("Cliente criado: %v", client)

	// 3. Atualizar os dados do cliente com o modelo Client
	client.Name = "Cliente Atualizado"
	client.WeightKg = 75
	client.Address = "Rua Teste Atualizada, 456"

	// Atualizar no banco
	if err := config.DB.Save(&client).Error; err != nil {
		log.Fatalf("Erro ao atualizar cliente no banco: %v", err)
	}
	t.Logf("Cliente atualizado: %v", client)

	// 4. Buscar o cliente no banco de dados após a atualização
	var dbClient models.Client
	if err := config.DB.First(&dbClient, client.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			t.Fatalf("Cliente não encontrado no banco com ID %v", client.ID)
		}
		log.Fatalf("Erro ao buscar cliente atualizado: %v", err)
	}

	// 5. Verificar se os dados foram atualizados corretamente
	assert.Equal(t, client.Name, dbClient.Name)
	assert.Equal(t, client.WeightKg, dbClient.WeightKg)
	assert.Equal(t, client.Address, dbClient.Address)
}
