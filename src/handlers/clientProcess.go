package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"myapi/models"
	"myapi/services"
	"net/http"
	"net/url"
)

// processClient é responsável por processar a criação de um novo cliente.
// Recebe um objeto `client` e um `validateType` para validar o tipo de operação.
// Retorna um map contendo o ID da operação ou um erro caso haja falha.
func ProcessClient(payload models.Client) (map[string]interface{}, error) {
	// Valida os dados do cliente antes de criar
	response, err := services.CreateClientCheckValues(payload)
	if err != nil || response["status"] != "valid" {
		// Retorna um map com erro de validação
		return map[string]interface{}{"error": "dados inválidos para criação do cliente"}, fmt.Errorf("dados inválidos para criação do cliente")
	}

	// Insere o cliente no banco de dados
	operationResponse, err := services.InsertData(payload)
	if err != nil {
		// Retorna o erro com as informações da operação falha
		return map[string]interface{}{"error": fmt.Sprintf("erro ao inserir cliente: %v", err)}, err
	}

	// Retorna a resposta de sucesso com a operação realizada
	// Você pode adicionar informações do cliente inserido ou algo mais relevante
	return map[string]interface{}{"operationID": operationResponse.ID}, nil
}

// processClientUpdate processa a atualização de um cliente existente.
// Recebe um objeto `clientUpdate` e um `validateType` para verificar o tipo da operação.
// Retorna um map com o ID da operação, ID do cliente e o tipo de operação, ou um erro em caso de falha.
func ProcessClientUpdate(clientUpdate models.ClientUpdate) (map[string]interface{}, error) {
	// Valida a atualização do cliente
	response, err := services.ValidateClientUpdate(clientUpdate)
	if err != nil || response["status"] != "valid" {
		return map[string]interface{}{
			"error":         "dados inválidos para atualização do cliente",
			"original_data": clientUpdate, // Retorna o payload original para depuração
		}, fmt.Errorf("dados inválidos para atualização do cliente")
	}

	// Tenta atualizar os dados do cliente no banco de dados
	operationResponse, err := services.UpdateClientData(clientUpdate)
	if err != nil {
		return map[string]interface{}{
			"error":         fmt.Sprintf("erro ao atualizar cliente: %v", err),
			"original_data": clientUpdate, // Retorna o payload original para depuração
		}, fmt.Errorf("erro ao atualizar cliente: %v", err)
	}

	// Converte operationResponse (models.ClientUpdate) para map[string]interface{}
	result := map[string]interface{}{
		"name":         operationResponse.Name,
		"weight_kg":    operationResponse.WeightKg,
		"address":      operationResponse.Address,
		"street":       operationResponse.Street,
		"number":       operationResponse.Number,
		"neighborhood": operationResponse.Neighborhood,
		"complement":   operationResponse.Complement,
		"city":         operationResponse.City,
		"state":        operationResponse.State,
		"country":      operationResponse.Country,
		"latitude":     operationResponse.Latitude,
		"longitude":    operationResponse.Longitude,
	}

	return result, nil
}

// Chama a API DistanceMatrix para buscar o endereço com base na latitude e longitude
func GetLocationFromAddress(address string) (map[string]interface{}, error) {
	// Log do endereço recebido para consulta
	slog.Info("Recebendo endereço para consulta", slog.String("endereco", address))

	// Codificando o endereço para ser usado na URL
	encodedAddress := url.QueryEscape(address)
	slog.Info("Endereço codificado para consulta", slog.String("endereco_codificado", encodedAddress))

	// URL da API DistanceMatrix para fazer a consulta com o endereço
	apiURL := fmt.Sprintf("https://api.distancematrix.ai/maps/api/geocode/json?address=%s&key=ByEWZK6VJV8AY5BcsLihzoH1xGCcrInVoc7OR0tZLFwTFksNOCxSZqjbZPiQxbos", encodedAddress)
	slog.Info("URL da API configurada", slog.String("url", apiURL))

	// Fazendo a requisição GET
	resp, err := http.Get(apiURL)
	if err != nil {
		slog.Error("Erro ao fazer requisição para API DistanceMatrix", slog.String("error", err.Error()))
		return nil, fmt.Errorf("falha ao acessar a API de geocoding: %w", err)
	}
	defer func() {
		slog.Info("Fechando o corpo da resposta da API")
		resp.Body.Close()
	}()

	// Lendo e decodificando a resposta JSON
	var jsonResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		slog.Error("Erro ao decodificar resposta JSON da API", slog.String("error", err.Error()))
		return nil, fmt.Errorf("falha ao decodificar resposta JSON: %w", err)
	}

	// Formata o JSON decodificado com indentação para log legível
	formattedJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
	if err != nil {
		slog.Error("Erro ao formatar resposta JSON para log", slog.String("error", err.Error()))
	} else {
		// Usa slog.Any para logar a string formatada como está, sem caracteres de escape
		slog.Info("Resposta JSON decodificada com sucesso", slog.Any("resposta", string(formattedJSON)))
	}

	// Verifica se o status da API é OK
	status, ok := jsonResponse["status"].(string)
	if !ok || status != "OK" {
		slog.Error("Erro na resposta da API, status não é OK", slog.String("status", status))
		return nil, fmt.Errorf("erro na resposta da API de geocoding, status: %s", status)
	}

	// Verifica se há resultados válidos na resposta
	results, ok := jsonResponse["result"].([]interface{})
	if !ok || len(results) == 0 {
		slog.Error("Nenhum resultado encontrado na resposta da API")
		return nil, fmt.Errorf("nenhum resultado encontrado na resposta da API")
	}

	// Extraindo dados do primeiro resultado
	result := results[0].(map[string]interface{})
	addressComponents, ok := result["address_components"].([]interface{})
	if !ok || len(addressComponents) == 0 {
		slog.Error("Não foi possível extrair componentes do endereço")
		return nil, fmt.Errorf("não foi possível extrair os componentes do endereço")
	}

	// Extraindo a latitude e longitude
	geometry, ok := result["geometry"].(map[string]interface{})
	if !ok {
		slog.Error("Não foi possível extrair dados de geometria")
		return nil, fmt.Errorf("não foi possível extrair dados de geometria")
	}

	location, ok := geometry["location"].(map[string]interface{})
	if !ok {
		slog.Error("Não foi possível extrair localização geográfica")
		return nil, fmt.Errorf("não foi possível extrair localização geográfica")
	}

	lat, latOk := location["lat"].(float64)
	lng, lngOk := location["lng"].(float64)
	if !latOk || !lngOk {
		slog.Error("Não foi possível extrair latitude ou longitude")
		return nil, fmt.Errorf("não foi possível extrair latitude ou longitude")
	}

	// Cria o mapa para os dados essenciais
	output := map[string]interface{}{
		"latitude":     lat,
		"longitude":    lng,
		"display_name": result["formatted_address"],
		"address":      addressComponents,
	}

	// Retorna os dados extraídos
	slog.Info("Dados do endereço extraídos com sucesso", slog.Any("dados_extraidos", output))
	return output, nil
}
