package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"myapi/config"
	"myapi/handlers"
	"myapi/models"
	"myapi/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type APIController struct{}

// CreateClient lida com a criação de um cliente a partir do corpo da requisição.
// @Summary Cria um novo cliente
// @Tags deliveries
// @Description Recebe um JSON contendo os dados de um cliente e insere o registro no sistema.
// @Accept json
// @Produce json
// @Param client body models.Client true "Dados do cliente para criação"
// @Success 200 {object} map[string]interface{} "Cliente criado com sucesso"
// @Failure 400 {string} string "Requisição inválida: erro no corpo da requisição ou JSON malformado"
// @Failure 500 {string} string "Erro ao criar cliente no banco de dados"
// @Router /deliveries [post]

func (c *APIController) CreateClient(w http.ResponseWriter, r *http.Request) {
	slog.Info("Iniciando o processo de criação de cliente", slog.String("endpoint", "CreateClient"))

	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Erro ao ler o corpo da requisição", slog.String("error", err.Error()))
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}
	slog.Info("Corpo da requisição lido com sucesso")

	// Decodifica o JSON para um map de dados
	var client models.Client
	if err := json.Unmarshal(body, &client); err != nil {
		slog.Error("Erro ao decodificar JSON para cliente", slog.String("error", err.Error()))
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	// Processa a criação do cliente
	insertResponse, err := handlers.ProcessClient(client)
	if err != nil {
		slog.Error("Erro ao criar o cliente", slog.String("error", err.Error()))
		http.Error(w, "Erro ao criar o cliente", http.StatusBadRequest)
		return
	}

	// Responde com sucesso
	c.respondWithJSON(w, insertResponse)
	slog.Info("Cliente criado e resposta enviada com sucesso", slog.String("client_name", client.Name))
}

// GetClients lida com a requisição GET para buscar clientes com paginação e filtros de cidade e ID.
// @Summary Busca clientes com paginação e filtros de cidade e ID
// @Tags deliveries
// @Description Retorna uma lista de clientes paginada, ou um cliente específico se o ID for fornecido, permitindo definir o limite, offset, cidade e ID.
// @Param id query int false "ID do cliente para busca específica"
// @Param limit query int false "Número máximo de clientes por página" default(100)
// @Param offset query int false "Número de registros a pular antes de começar a listar os clientes" default(0)
// @Param city query string false "Cidade para filtrar os clientes"
// @Success 200 {object} map[string]interface{} "Dados da lista de clientes com metadados de paginação"
// @Failure 500 {string} string "Erro ao buscar clientes"
// @Router /deliveries [get]

func (c *APIController) GetClients(w http.ResponseWriter, r *http.Request) {
	// Define valores padrão para `limit` e `offset`
	limit := 100
	offset := 0

	// Extrai `limit`, `offset`, `city` e `id` dos parâmetros de consulta
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
		slog.Info("Limit recebido", "limit", limit)
	} else {
		slog.Info("Valor default de limit utilizado", "limit", limit)
	}

	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o >= 0 {
		offset = o
		slog.Info("Offset recebido", "offset", offset)
	} else {
		slog.Info("Valor default de offset utilizado", "offset", offset)
	}

	city := r.URL.Query().Get("city")
	if city != "" {
		slog.Info("Filtro de cidade recebido", "city", city)
	}

	// Extrai o `id` para busca específica de cliente
	var id int
	if idParam := r.URL.Query().Get("id"); idParam != "" {
		if idVal, err := strconv.Atoi(idParam); err == nil && idVal > 0 {
			id = idVal
			slog.Info("Filtro de ID recebido", "id", id)
		} else {
			slog.Error("ID inválido fornecido", "id", idParam)
			http.Error(w, "ID inválido fornecido", http.StatusBadRequest)
			return
		}
	}

	// Se um ID foi fornecido, busca apenas o cliente específico
	if id > 0 {
		var client models.Client
		if err := config.DB.First(&client, id).Error; err != nil {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			slog.Error("Erro ao buscar cliente específico", "error", err, "id", id)
			return
		}
		slog.Info("Cliente específico encontrado", "client", client)

		// Retorna o cliente encontrado	 como resposta
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"client": client}); err != nil {
			slog.Error("Erro ao enviar resposta do cliente específico", slog.String("error", err.Error()))
			http.Error(w, "Erro ao enviar resposta", http.StatusInternalServerError)
		}
		return
	}

	// Conta o total de clientes no banco de dados com filtro de cidade, se fornecido
	var total int64
	dbQuery := config.DB.Model(&models.Client{})
	if city != "" {
		dbQuery = dbQuery.Where("city = ?", city)
	}
	if err := dbQuery.Count(&total).Error; err != nil {
		http.Error(w, "Failed to count clients", http.StatusInternalServerError)
		slog.Error("Erro ao contar clientes", "error", err)
		return
	}
	slog.Info("Total de clientes", "total", total)

	// Busca os clientes com o limit e offset definidos, e aplica filtro de cidade, se fornecido
	var clients []models.Client
	dbQuery = config.DB.Limit(limit).Offset(offset)
	if city != "" {
		dbQuery = dbQuery.Where("city = ?", city)
	}
	if err := dbQuery.Find(&clients).Error; err != nil {
		http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
		slog.Error("Erro ao buscar clientes", "error", err)
		return
	}
	slog.Info("Clientes encontrados", "num_clients", len(clients))

	// Calcula total de páginas, página atual e próximo offset
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // Arredonda para cima
	currentPage := (offset / limit) + 1
	nextOffset := offset + limit

	// Monta a URL para a próxima página se houver
	var nextPageURL *string
	if nextOffset < int(total) {
		url := fmt.Sprintf("%s?limit=%d&offset=%d", r.URL.Path, limit, nextOffset)
		if city != "" {
			url += fmt.Sprintf("&city=%s", city)
		}
		nextPageURL = &url
		slog.Info("URL da próxima página", "nextPageURL", *nextPageURL)
	}

	// Monta a resposta com os clientes e metadados
	response := map[string]interface{}{
		"clients":     clients,
		"total":       total,
		"totalPages":  totalPages,
		"currentPage": currentPage,
		"nextPageURL": nextPageURL,
	}

	// Envia a resposta
	c.respondWithJSON(w, response)
	slog.Info("Resposta de clientes enviada com sucesso")
}

// DeleteClientsHandler lida com a exclusão de clientes.
// @Summary Excluir clientes
// @Description Exclui todos os clientes ou um cliente específico pelo ID.
// @Tags deliveries
// @Param deleteAll query bool false "Excluir todos os clientes (true para excluir todos os clientes)"
// @Param id query int false "ID do cliente a ser excluído (se deleteAll não for especificado)"
// @Success 200 {string} string "Mensagem de sucesso"
// @Failure 400 {string} string "Parâmetros inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /deliveries [delete]

func (c *APIController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Parse os parâmetros da URL
	queryParams := r.URL.Query()
	deleteAll := queryParams.Get("deleteAll")
	idParam := queryParams.Get("id")

	// Valida se ambos os parâmetros foram fornecidos
	if deleteAll != "" && idParam != "" {
		http.Error(w, "Somente um parâmetro pode ser fornecido: 'deleteAll' ou 'id'. Não forneça ambos.", http.StatusBadRequest)
		return
	}

	// Caso deleteAll seja verdadeiro
	if deleteAll == "true" {
		err := services.DeleteAllClients()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Todos os clientes foram excluídos com sucesso."))
		return
	}

	// Caso id seja especificado
	if idParam != "" {
		clientID, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID inválido: o valor deve ser um número", http.StatusBadRequest)
			return
		}

		// Converte o ID para uint8 (caso o serviço use uint8)
		clientIDUint8 := int(clientID)

		err = services.DeleteClientByID(clientIDUint8)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cliente excluído com sucesso."))
		return
	}

	// Caso nenhum parâmetro seja especificado
	http.Error(w, "Parâmetros inválidos. Use ?deleteAll=true ou ?id=<ID>", http.StatusBadRequest)
}

// UpdateClient lida com a atualização de um cliente a partir do corpo da requisição.
// @Summary Atualiza um cliente
// @Tags deliveries
// @Description Modifica os dados de um cliente existente.
// @Param id query int false "ID do cliente para busca específica"
// @Param client body models.ClientUpdate true "Cliente para atualizar"
// @Success 200 {object} map[string]interface{} "Resposta com os dados do cliente atualizado"
// @Failure 400 {string} string "Corpo da requisição inválido"
// @Failure 404 {string} string "Cliente não encontrado"
// @Failure 500 {string} string "Erro ao atualizar cliente"
// @Router /deliveries [put]

func (c *APIController) UpdateClient(w http.ResponseWriter, r *http.Request) {
	slog.Info("Iniciando o processo de atualização de cliente", slog.String("endpoint", "UpdateClient"))

	// Obtém o ID do cliente dos parâmetros da URL
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		slog.Error("ID do cliente não fornecido na URL")
		http.Error(w, "ID do cliente não fornecido", http.StatusBadRequest)
		return
	}

	// Converte o ID para inteiro
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error("ID do cliente inválido", slog.String("error", err.Error()))
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Busca o cliente no banco de dados para verificar se existe
	var existingClient models.Client
	if err := config.DB.First(&existingClient, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("Cliente não encontrado", slog.Int("id", id))
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}
		slog.Error("Erro ao buscar cliente", slog.String("error", err.Error()))
		http.Error(w, "Erro interno ao buscar cliente", http.StatusInternalServerError)
		return
	}
	slog.Info("Cliente encontrado no banco de dados", slog.Int("client_id", id))

	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Erro ao ler o corpo da requisição", slog.String("error", err.Error()))
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}
	slog.Info("Corpo da requisição lido com sucesso")

	// Decodifica o corpo como um único ClientUpdate
	var clientUpdate models.ClientUpdate
	if err := json.Unmarshal(body, &clientUpdate); err != nil {
		slog.Error("Erro ao decodificar JSON", slog.String("error", err.Error()))
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	// Combina o ID do cliente da URL com o payload recebido
	clientUpdate.ID = uint(id)

	// Processa a atualização do cliente
	updatedClient, err := handlers.ProcessClientUpdate(clientUpdate)
	if err != nil {
		slog.Error("Erro ao processar atualização do cliente", slog.String("error", err.Error()))
		http.Error(w, "Erro ao atualizar o cliente", http.StatusInternalServerError)
		return
	}

	// Responde com sucesso para o caso de atualização
	c.respondWithJSON(w, updatedClient)

	slog.Info("Cliente atualizado e resposta enviada com sucesso", slog.String("client_name", clientUpdate.Name))
}

// SearchAddress lida com a busca de latitude e longitude a partir de um endereço fornecido.
// @Summary Busca latitude e longitude de um endereço
// @Tags geocoding
// @Description Retorna as coordenadas geográficas (latitude e longitude) de um endereço fornecido.
// @Param endereco query string true "Endereço a ser consultado"
// @Success 200 {object} map[string]float64 "Resposta com latitude e longitude do endereço"
// @Failure 400 {string} string "Parâmetro 'endereco' ausente ou inválido"
// @Failure 500 {string} string "Erro ao consultar a API de geocoding"
// @Router /deliveries/geoconding/search [get]

func (c *APIController) SearchAddress(w http.ResponseWriter, r *http.Request) {
	slog.Info("Iniciando o processo de busca de lat e long por endereço", slog.String("endpoint", "Geocoding"))

	// Obtém o parâmetro 'endereco' da URL
	endereco := r.URL.Query().Get("endereco")
	if endereco == "" {
		slog.Error("Parâmetro 'endereco' não fornecido na URL")
		http.Error(w, "Parâmetro 'endereco' é obrigatório", http.StatusBadRequest)
		return
	}
	slog.Info("Parâmetro 'endereco' recebido", slog.String("endereco", endereco))

	// Log do método HTTP e detalhes da requisição
	slog.Info("Detalhes da requisição recebida",
		slog.String("method", r.Method),
		slog.String("url", r.URL.String()),
		slog.String("remote_addr", r.RemoteAddr),
	)

	// Chama a função para obter os dados do endereço
	slog.Info("Chamando a função de geocoding para obter os dados do endereço")
	locationData, err := handlers.GetLocationFromAddress(endereco)
	if err != nil {
		slog.Error("Erro ao consultar a API de geocoding", slog.String("error", err.Error()))
		// Verifica se o erro é do tipo "nenhum resultado encontrado"
		if err.Error() == "nenhum resultado encontrado na resposta da API" {
			http.Error(w, "Nenhum resultado encontrado para o endereço fornecido", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao consultar a API de geocoding", http.StatusInternalServerError)
		}
		return
	}
	slog.Info("Geocoding realizado com sucesso", slog.Any("location_data", locationData))

	// Configura o cabeçalho da resposta
	slog.Info("Configurando cabeçalho da resposta para JSON")
	w.Header().Set("Content-Type", "application/json")

	// Responde com os dados em JSON
	slog.Info("Enviando resposta com os dados do endereço em JSON")
	if err := json.NewEncoder(w).Encode(locationData); err != nil {
		slog.Error("Erro ao codificar a resposta JSON", slog.String("error", err.Error()))
		http.Error(w, "Erro ao gerar a resposta JSON", http.StatusInternalServerError)
		return
	}
	slog.Info("Resposta enviada com sucesso", slog.String("status", "200 OK"))
}

// respondWithJSON envia uma resposta JSON ao cliente.
// Esta função é usada internamente para padronizar a resposta da API.
// @Description Envia resposta JSON ao cliente
// @Param data body map[string]interface{} true "Dados a serem enviados ao cliente"
// @Failure 500 {object} map[string]string "Erro ao gerar ou enviar a resposta JSON"

func (c *APIController) respondWithJSON(w http.ResponseWriter, data map[string]interface{}) {
	slog.Info("Iniciando o envio de resposta JSON", slog.String("endpoint", "respondWithJSON"))

	// Define o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	// Tenta converter os dados para JSON
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		// Log de erro se ocorrer um problema ao gerar a resposta JSON
		slog.Error("Erro ao gerar a resposta JSON", slog.String("error", err.Error()))
		http.Error(w, "Erro ao enviar resposta", http.StatusInternalServerError)
		return
	}

	// Envia a resposta com o JSON gerado
	if _, err := w.Write(jsonResponse); err != nil {
		// Log de erro se falhar ao escrever a resposta
		slog.Error("Erro ao enviar a resposta JSON", slog.String("error", err.Error()))
		http.Error(w, "Erro ao enviar resposta", http.StatusInternalServerError)
		return
	}

	slog.Info("Resposta JSON enviada com sucesso")
}

// Método para registrar as rotas da API
// Este método define as rotas de endpoint para as operações de entrega. Cada rota é associada a um
// método do controlador correspondente para lidar com as requisições.
func (c *APIController) RegisterRoutes(r *mux.Router) {
	slog.Info("Registrando rotas da API", slog.String("endpoint", "RegisterRoutes"))

	// Definindo a rota para criar um novo cliente
	r.HandleFunc("/deliveries", c.CreateClient).Methods("POST")
	slog.Info("Rota '/deliveries' registrada para POST")

	// Definindo a rota para obter todos os clientes
	r.HandleFunc("/deliveries", c.GetClients).Methods("GET")
	slog.Info("Rota '/deliveries' registrada para GET")

	//Definindo a rota para deletar um cliente
	r.HandleFunc("/deliveries", c.DeleteHandler).Methods("DELETE")
	slog.Info("Rota '/deliveries' registrada para DELETE")

	// Definindo a rota para atualizar um cliente existente
	r.HandleFunc("/deliveries", c.UpdateClient).Methods("PUT")
	slog.Info("Rota '/deliveries' registrada para PUT")

	// Definindo a rota para buscar usar o geoconding reverse
	r.HandleFunc("/deliveries/geoconding/search", c.SearchAddress).Methods("GET")
	slog.Info("Rota '/deliveries/geoconding/search' registrada para GET")

	// Definindo a rota para buscar usar o geoconding reverse
	r.HandleFunc("/deliveries/geoconding/search", c.SearchAddress).Methods("GET")
	slog.Info("Rota '/deliveries/geoconding/search' registrada para GET")

	// Definindo a rota para deletar um cliente com base no id
	slog.Info("Todas as rotas da API foram registradas com sucesso")
}
