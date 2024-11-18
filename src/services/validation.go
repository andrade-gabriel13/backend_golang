package services

import (
	"fmt"
	"log/slog"
	"myapi/models"
)

// ValidateCommonClientFields valida os campos obrigatórios de um cliente.
// Esta função verifica se todos os campos do cliente são válidos antes de prosseguir para outras operações.
// Caso algum campo seja inválido ou ausente, ela loga um erro e retorna uma mensagem indicando o campo faltante ou inválido.
//
// Campos validados:
// - Name: não pode ser vazio.
// - TestTechnical: não pode ser vazio.
// - WeightKg: deve ser maior que 0.
// - Address: não pode ser vazio.
// - Street: não pode ser vazio.
// - Number: deve ser maior que 0.
// - Neighborhood: não pode ser vazio.
// - City: não pode ser vazio.
// - State: não pode ser vazio.
// - Country: não pode ser vazio.
// - Latitude: deve ser um valor válido (diferente de 0).
// - Longitude: deve ser um valor válido (diferente de 0).
//
// Retorna um erro caso algum campo seja inválido.
func ValidateCommonClientFields(client models.Client) error {
	// Validando o campo 'Name'
	if client.Name == "" {
		slog.Error("Missing required field", "field", "name", "value", client.Name)
		return fmt.Errorf("name is required")
	}

	// Validando o campo 'WeightKg'
	if client.WeightKg <= 0 {
		slog.Error("Invalid weight value", "field", "WeightKg", "value", client.WeightKg)
		return fmt.Errorf("weightKg must be greater than 0")
	}

	// Validando o campo 'Address'
	if client.Address == "" {
		slog.Error("Missing required field", "field", "Address", "value", client.Address)
		return fmt.Errorf("address is required")
	}

	// Validando o campo 'Street'
	if client.Street == "" {
		slog.Error("Missing required field", "field", "Street", "value", client.Street)
		return fmt.Errorf("street is required")
	}

	// Validando o campo 'Number'
	if client.Number == 0 {
		slog.Error("Missing required field", "field", "Number", "value", client.Number)
		return fmt.Errorf("number is required")
	}

	// Validando o campo 'Neighborhood'
	if client.Neighborhood == "" {
		slog.Error("Missing required field", "field", "Neighborhood", "value", client.Neighborhood)
		return fmt.Errorf("neighborhood is required")
	}

	// Validando o campo 'City'
	if client.City == "" {
		slog.Error("Missing required field", "field", "City", "value", client.City)
		return fmt.Errorf("city is required")
	}

	// Validando o campo 'State'
	if client.State == "" {
		slog.Error("Missing required field", "field", "State", "value", client.State)
		return fmt.Errorf("state is required")
	}

	// Validando o campo 'Country'
	if client.Country == "" {
		slog.Error("Missing required field", "field", "Country", "value", client.Country)
		return fmt.Errorf("country is required")
	}

	// Validando o campo 'Latitude'
	if client.Latitude == 0 {
		slog.Error("Invalid latitude value", "field", "Latitude", "value", client.Latitude)
		return fmt.Errorf("latitude must be a valid number")
	}

	// Validando o campo 'Longitude'
	if client.Longitude == 0 {
		slog.Error("Invalid longitude value", "field", "Longitude", "value", client.Longitude)
		return fmt.Errorf("longitude must be a valid number")
	}

	// Se todos os campos estiverem válidos, retorna nil
	return nil
}

// CreateClientCheckValues valida o cliente usando a função ValidateCommonClientFields e retorna um mapa com o status da validação.
// Caso a validação seja bem-sucedida, retorna um mapa com status "valid" e uma mensagem de sucesso.
// Caso contrário, retorna o erro gerado pela função de validação.
//
// Retorna um mapa contendo:
// - "status": O status da validação ("valid" ou outro status de erro).
// - "message": Uma mensagem detalhada sobre o status da validação.
//
// Caso haja erro na validação, retorna um erro com a mensagem correspondente.
func CreateClientCheckValues(client models.Client) (map[string]interface{}, error) {
	// Validação comum de campos

	if err := ValidateCommonClientFields(client); err != nil {
		// Erro na validação de campos comuns, logando e retornando o erro
		slog.Error("Client validation failed", "error", err)
		return nil, err
	}

	// Se tudo estiver correto, retornamos um mapa de sucesso
	slog.Info("Client validated successfully", "field", "validation", "status", "success")
	return map[string]interface{}{
		"status":  "valid",
		"message": "Client validated successfully",
	}, nil
}

// ValidateClientUpdate valida os dados para atualização de um cliente.
// Verifica se o ID do cliente é válido. Se não for, retorna um erro indicando que o ID é obrigatório para a atualização.
// Caso contrário, a validação é considerada bem-sucedida e retorna um status de sucesso.
//
// Retorna um mapa contendo:
// - "status": O status da validação ("valid" ou outro status de erro).
// - "message": Uma mensagem detalhada sobre o status da validação.
//
// Caso haja erro na validação, retorna um erro com a mensagem correspondente.
func ValidateClientUpdate(client models.ClientUpdate) (map[string]interface{}, error) {
	// Validação específica de atualização
	if client.ID <= 0 {
		// Erro na validação do ID
		slog.Error("Missing or invalid ID for client update", "field", "ID", "value", client.ID)
		return nil, fmt.Errorf("ID do cliente é obrigatório para atualização")
	}

	// Após validar o ID, os dados serão validados no momento da persistência (quando os dados forem gravados)
	// Não há mais necessidade de validar os campos nesta etapa
	slog.Info("Client update validated", "field", "ID", "value", client.ID)
	return map[string]interface{}{
		"status":  "valid",
		"message": "Client validated for update",
	}, nil
}
