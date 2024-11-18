package services

import (
	"fmt"
	"log"
	"log/slog"
	"myapi/config"
	"myapi/models"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// InsertData insere um novo cliente no banco de dados.
//
// Parâmetros:
// - client (models.Client): Estrutura contendo os dados do cliente a ser inserido.
//
// Retorno:
// - models.Client: Estrutura preenchida com os dados do cliente, incluindo os campos gerados pelo banco de dados, como ID e timestamps.
// - error: Retorna um erro caso a operação de inserção falhe.
//
// Detalhes:
// - A função utiliza o método `Create` do GORM para salvar os dados do cliente na tabela correspondente no banco de dados.
// - Caso ocorra um erro durante a inserção, um log será gerado e o erro será retornado.
// func InsertData(client models.Client) (models.Client, error) {
//	if err := config.DB.Create(&client).Error; err != nil {
//		log.Println("Erro ao inserir o cliente no MySQL:", err)
//		return models.Client{}, err // Retorna estrutura vazia e erro
//	}

// Retorna o objeto client com todos os campos preenchidos após a inserção
//	return client, nil
//}

func InsertData(client models.Client) (models.Client, error) {
	if err := config.DB.Create(&client).Error; err != nil {
		log.Println("Erro ao inserir o cliente no MySQL:", err)
		return models.Client{}, err // Retorna estrutura vazia e erro
	}

	// Retorna o objeto client com todos os campos preenchidos após a inserção
	return client, nil
}

// UpdateClientData atualiza os dados de um cliente existente no banco de dados.
//
// Parâmetros:
// - client (models.ClientUpdate): Estrutura contendo os campos a serem atualizados. O campo "ID" é obrigatório
//   e identifica o cliente que será atualizado.
//
// Retorno:
// - models.ClientResponse: Estrutura contendo os dados do cliente atualizados no formato correto.
// - error: Um erro será retornado nos seguintes casos:
//   - O campo "ID" não foi fornecido no parâmetro client.
//   - O cliente com o ID especificado não foi encontrado no banco de dados.
//   - Não há campos válidos para atualizar.
//   - A operação de atualização ou busca dos dados atualizados falhou.
//
// Detalhes:
// - A função verifica se o cliente existe antes de realizar a atualização.
// - Apenas os campos que não possuem valor zero (não inicializados) serão atualizados, ignorando os campos de metadados
//   como "CreatedAt", "UpdatedAt" e "DeletedAt".
// - Após a atualização, os dados do cliente são buscados novamente para retornar as informações atualizadas no formato esperado.

func UpdateClientData(client models.ClientUpdate) (models.ClientResponse, error) {

	// Verifica se o cliente com o ID fornecido existe
	var existingClient models.Client
	if err := config.DB.Where("id = ?", client.ID).First(&existingClient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.ClientResponse{}, fmt.Errorf("cliente com ID %d não encontrado", client.ID)
		}
		return models.ClientResponse{}, fmt.Errorf("erro ao verificar cliente: %v", err)
	}

	updateData := map[string]interface{}{}

	// Usa reflexão para iterar sobre os campos do struct ClientUpdate
	clientValue := reflect.ValueOf(client)
	clientType := reflect.TypeOf(client)

	for i := 0; i < clientValue.NumField(); i++ {
		field := clientValue.Field(i)
		fieldName := clientType.Field(i).Name

		// Ignora campos de gorm.Model (ID, CreatedAt, etc)
		if fieldName == "ID" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "DeletedAt" {
			continue
		}

		// Verifica se o campo tem valor e não é zero
		if !isZeroValue(field) {
			updateData[fieldName] = field.Interface()
		}
	}

	// Se não houver dados a serem atualizados, retorne erro
	if len(updateData) == 0 {
		return models.ClientResponse{}, fmt.Errorf("nenhum campo válido foi enviado para atualização")
	}

	// Executa a atualização no banco de dados
	err := config.DB.Model(&models.Client{}).Where("id = ?", client.ID).Updates(updateData).Error
	if err != nil {
		return models.ClientResponse{}, err
	}

	// Busca os dados atualizados do cliente
	var updatedClient models.Client
	if err := config.DB.Where("id = ?", client.ID).First(&updatedClient).Error; err != nil {
		return models.ClientResponse{}, fmt.Errorf("erro ao buscar cliente atualizado: %v", err)
	}

	// Cria um struct de resposta com a ordem correta dos campos
	response := models.ClientResponse{
		Name:         updatedClient.Name,
		WeightKg:     updatedClient.WeightKg,
		Address:      updatedClient.Address,
		Street:       updatedClient.Street,
		Number:       updatedClient.Number,
		Neighborhood: updatedClient.Neighborhood,
		Complement:   updatedClient.Complement,
		City:         updatedClient.City,
		State:        updatedClient.State,
		Country:      updatedClient.Country,
		Latitude:     updatedClient.Latitude,
		Longitude:    updatedClient.Longitude,
	}

	// Retorna os dados formatados
	slog.Info("Cliente atualizado com sucesso", slog.Int("client_id", int(client.ID)))
	return response, nil
}

// GetClientByID busca um cliente pelo ID e retorna as informações do cliente encontrado.
//
// Parâmetros:
// - id (uint8): O ID do cliente que será buscado no banco de dados.
//
// Retorno:
// - models.ClientUpdate: Estrutura contendo as informações do cliente encontrado.
// - error: Retorna um erro se o cliente não for encontrado ou ocorrer algum problema na consulta.
//
// Comportamento:
// - Realiza uma consulta no banco de dados para buscar o cliente com o ID especificado.
// - Se o cliente for encontrado, retorna os dados do cliente e `nil` para o erro.
// - Se não for encontrado ou ocorrer um erro, retorna um cliente vazio e o erro correspondente.
//
// Exemplo de uso:
//
//	client, err := GetClientByID(1)
//	if err != nil {
//		log.Println("Erro ao buscar cliente:", err)
//	} else {
//		log.Printf("Cliente encontrado: %+v\n", client)
//	}

func GetClientByID(id uint8) (models.ClientUpdate, error) {
	// Declara uma variável client para armazenar o cliente encontrado
	var client models.ClientUpdate

	// Tenta encontrar o cliente no banco de dados pelo id
	if err := config.DB.First(&client, id).Error; err != nil {
		log.Println("Erro ao buscar o cliente:", err)
		return client, err // Retorna o cliente vazio e o erro
	}

	// Retorna o cliente encontrado e nil, indicando que não houve erro
	return client, nil
}

// isZeroValue verifica se o valor de um campo é considerado zero ou não inicializado.
//
// Esta função é utilizada para determinar se um campo deve ser ignorado em operações como atualizações
// parciais, onde campos não inicializados não devem sobrescrever valores existentes.
//
// Parâmetros:
// - value (reflect.Value): O valor do campo a ser avaliado.
//
// Retorno:
// - bool: Retorna `true` se o valor for considerado zero ou não inicializado, ou `false` caso contrário.
//
// Comportamento:
// - Para tipos Slice, Map e Array, verifica se o comprimento é igual a zero.
// - Para tipos Pointer e Interface, verifica se o valor é `nil`.
// - Para outros tipos, utiliza o método `IsZero()` para determinar o estado.
//
// Exemplo de uso:
//
//	fieldValue := reflect.ValueOf(myStruct.SomeField)
//	if isZeroValue(fieldValue) {
//		log.Println("O campo está vazio ou não inicializado")
//	} else {
//		log.Println("O campo possui um valor válido")
//	}

func isZeroValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		return value.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return value.IsNil()
	}
	return value.IsZero()
}

// DeleteAllClients arquiva e exclui todos os clientes da tabela principal (`clients`).
//
// Esta função realiza as seguintes etapas:
// 1. Busca todos os registros de clientes na tabela principal.
// 2. Converte os clientes encontrados para o formato de arquivamento (`archived_clients`),
//    incluindo o timestamp de arquivamento (`DeletedAt`).
// 3. Insere os registros convertidos na tabela de arquivados em lote.
// 4. Remove todos os registros da tabela principal (`clients`) em lote.
//
// Retorno:
// - error: Retorna `nil` se a operação for bem-sucedida ou um erro descritivo caso ocorra falha
//   em qualquer etapa do processo.
//
// Erros possíveis:
// - Falha ao buscar os clientes na tabela principal.
// - Falha ao inserir os registros na tabela de arquivados.
// - Falha ao deletar os registros da tabela principal.
//
// Observação:
// - Caso não existam clientes na tabela principal, a função retorna `nil` e registra uma mensagem
//   indicando que não há clientes para excluir.
//
// Exemplo de uso:
//
//	err := DeleteAllClients()
//	if err != nil {
//		log.Printf("Erro ao arquivar e excluir todos os clientes: %v", err)
//	} else {
//		log.Println("Todos os clientes foram processados com sucesso")
//	}

func DeleteAllClients() error {
	var clients []models.Client

	// Buscar todos os clientes na tabela principal
	if err := config.DB.Find(&clients).Error; err != nil {
		log.Println("Erro ao buscar todos os clientes:", err)
		return fmt.Errorf("erro ao buscar todos os clientes")
	}

	// Verificar se existem clientes para arquivar
	if len(clients) == 0 {
		log.Println("Nenhum cliente encontrado para excluir")
		return nil
	}

	// Converter os clientes para o formato de arquivamento
	var archivedClients []models.ArchivedClient
	for _, client := range clients {
		archivedClients = append(archivedClients, models.ArchivedClient{
			ID:           int(client.ID), // Conversão de uint para int
			Name:         client.Name,
			WeightKg:     client.WeightKg,
			Address:      client.Address,
			Street:       client.Street,
			Number:       client.Number, // Conversão de int para string
			Neighborhood: client.Neighborhood,
			Complement:   client.Complement,
			City:         client.City,
			State:        client.State,
			Country:      client.Country,
			Latitude:     client.Latitude,
			Longitude:    client.Longitude,
			CreatedAt:    client.CreatedAt,
			UpdatedAt:    client.UpdatedAt,
			DeletedAt:    time.Now(),
		})
	}

	// Inserir todos os clientes na tabela de arquivados em lote
	for _, archivedClient := range archivedClients {
		if err := config.DB.Save(&archivedClient).Error; err != nil {
			log.Println("Erro ao arquivar cliente:", archivedClient.ID, err)
			return fmt.Errorf("erro ao arquivar o cliente com ID %d", archivedClient.ID)
		}
	}

	// Deletar todos os clientes da tabela principal em lote
	if err := config.DB.Exec("DELETE FROM clients").Error; err != nil {
		log.Println("Erro ao deletar todos os clientes:", err)
		return fmt.Errorf("erro ao deletar todos os clientes")
	}

	log.Println("Todos os clientes foram arquivados e excluídos com sucesso")
	return nil
}

// DeleteClientByID arquiva e exclui um cliente específico com base no ID fornecido.
//
// Esta função realiza as seguintes etapas:
// 1. Busca um cliente na tabela principal (`clients`) com o ID especificado.
// 2. Cria um registro do cliente na tabela de arquivados (`archived_clients`), incluindo
//    informações completas do cliente e o timestamp de arquivamento (`DeletedAt`).
// 3. Remove o cliente da tabela principal (`clients`).
//
// Parâmetros:
// - clientID (int): ID do cliente a ser arquivado e excluído.
//
// Retorno:
// - error: Retorna `nil` se a operação for bem-sucedida ou um erro descritivo caso ocorra falha
//   em qualquer etapa do processo.
//
// Erros possíveis:
// - Caso o cliente não seja encontrado na tabela principal.
// - Falha ao inserir o registro na tabela de arquivados.
// - Falha ao deletar o cliente da tabela principal.
//
// Exemplo de uso:
//
//	err := DeleteClientByID(123)
//	if err != nil {
//		log.Printf("Erro ao arquivar e excluir cliente: %v", err)
//	}

func DeleteClientByID(clientID int) error {
	var client models.Client

	// Buscar o cliente pelo ID na tabela principal
	if err := config.DB.Table("clients").First(&client, "id = ?", clientID).Error; err != nil {
		log.Println("Erro ao encontrar o cliente:", err)
		return fmt.Errorf("erro ao encontrar o cliente com ID %d", clientID)
	}

	// Criar um registro do cliente arquivado
	archivedClient := models.ArchivedClient{
		ID:           int(client.ID), // Garante compatibilidade caso ID do ArchivedClient seja int
		Name:         client.Name,
		WeightKg:     client.WeightKg,
		Address:      client.Address,
		Street:       client.Street,
		Number:       client.Number,
		Neighborhood: client.Neighborhood,
		Complement:   client.Complement,
		City:         client.City,
		State:        client.State,
		Country:      client.Country,
		Latitude:     client.Latitude,
		Longitude:    client.Longitude,
		CreatedAt:    client.CreatedAt,
		UpdatedAt:    client.UpdatedAt,
		DeletedAt:    time.Now(), // Adiciona o timestamp atual para arquivamento
	}

	// Inserir na tabela de arquivados
	if err := config.DB.Table("archived_clients").Create(&archivedClient).Error; err != nil {
		log.Println("Erro ao excluir cliente:", err)
		return fmt.Errorf("erro ao excluir o cliente com ID %d", clientID)
	}

	// Deletar o cliente da tabela principal
	if err := config.DB.Table("clients").Delete(&models.Client{}, "id = ?", clientID).Error; err != nil {
		log.Println("Erro ao deletar o cliente:", err)
		return fmt.Errorf("erro ao deletar o cliente com ID %d", clientID)
	}

	log.Printf("Cliente com ID %d foi arquivado e excluído com sucesso", clientID)
	return nil
}
