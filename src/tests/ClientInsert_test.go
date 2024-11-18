package tests

import (
	"bytes"
	"encoding/json"
	"myapi/config"
	"myapi/controller"

	"myapi/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateClient(t *testing.T) {
	// Conecta ao banco de dados.
	config.ConnectDB()

	// Instancia o controller
	controller := &controller.APIController{}

	// Cria um cliente fictício para o teste.
	client := models.Client{
		Name:         "Teste Cliente",
		WeightKg:     70,
		Address:      "Rua Teste, 123",
		Street:       "Rua Teste",
		Number:       123,
		Neighborhood: "Bairro Teste",
		City:         "Cidade Teste",
		State:        "Estado Teste",
		Country:      "Brasil",
		Latitude:     -22.619,
		Longitude:    -43.164,
	}

	// Cria o corpo da requisição JSON para enviar para o método CreateClient.
	clientData, err := json.Marshal(client)
	if err != nil {
		t.Fatalf("Erro ao marshaling client: %v", err)
	}

	// Cria a requisição HTTP simulando o que o CreateClient esperaria.
	req, err := http.NewRequest("POST", "/clients", bytes.NewBuffer(clientData))
	if err != nil {
		t.Fatalf("Erro ao criar requisição: %v", err)
	}

	// Cria um ResponseWriter fictício.
	rr := httptest.NewRecorder()

	// Chama o método CreateClient da APIController
	controller.CreateClient(rr, req)

	// Verifica se a resposta HTTP foi bem-sucedida
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verifica se o cliente foi criado no banco.
	var createdClient models.Client
	err = config.DB.First(&createdClient, "name = ?", "Teste Cliente").Error
	assert.NoError(t, err)
	assert.Equal(t, "Teste Cliente", createdClient.Name)
	assert.Equal(t, float64(70), createdClient.WeightKg)

}
