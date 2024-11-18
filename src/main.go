// main.go

package main

import (
	"fmt"
	"log"
	"log/slog"
	"myapi/config"
	"myapi/controller"
	"net/http"
	"os"

	_ "myapi/docs" // Importa a documentação gerada pelo swag

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger" // Importa o Swagger UI
)

// logger global
var logger *slog.Logger

// init inicializa o logger global com um handler de texto padrão.
// Ele configura a saída de logs no formato de texto simples para o terminal.
func init() {
	// Inicializando o logger global com opções padrão
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger = slog.New(handler)
}

// enableCORS adiciona os cabeçalhos necessários para habilitar o CORS (Cross-Origin Resource Sharing).
//
// Parâmetros:
// - next (http.Handler): Próximo handler na cadeia.
//
// Retorno:
// - http.Handler: Um handler que adiciona os cabeçalhos CORS e processa requisições OPTIONS.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// logRequest é um middleware que registra logs detalhados de requisições HTTP,
// incluindo informações sobre o método, URL, endereço remoto e status da resposta.
//
// Parâmetros:
// - next (http.Handler): Próximo handler na cadeia.
//
// Retorno:
// - http.Handler: Um handler que adiciona o registro de logs antes e depois de processar a requisição.
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logando informações da requisição antes de processar
		logger.Info("Recebendo requisição",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.String("remoteAddr", r.RemoteAddr),
		)

		// Criar um responseRecorder para capturar o status code
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)

		// Logando informações após processar a requisição
		logger.Info("Resposta enviada",
			slog.Int("statusCode", recorder.statusCode),
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
		)
	})
}

// responseRecorder é um wrapper para http.ResponseWriter, que permite capturar
// e registrar o status code da resposta HTTP enviada.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader sobrescreve o método WriteHeader original do http.ResponseWriter,
// permitindo capturar e armazenar o status code da resposta HTTP.
//
// Parâmetros:
// - statusCode (int): Código de status da resposta HTTP.
func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// @title API - Golang Desafio Backend
// @version 1.0
// @description API de exemplo com integração Swagger, oferecendo uma interface interativa para documentação e testes de endpoints. Esta API foi desenvolvida para ilustrar boas práticas de documentação, utilizando Swagger para facilitar a navegação e o entendimento das operações disponíveis, parâmetros de entrada, respostas e exemplos de uso. Indicada para desenvolvimento e demonstração de integração entre APIs e clientes.
// @host localhost:8080
// @BasePath /

func main() {

	// Criar o roteador
	r := mux.NewRouter()

	// Conectar ao banco de dados
	config.ConnectDB()

	// Criar uma instância do controlador
	controller := &controller.APIController{}

	// Registrar as rotas no controlador
	controller.RegisterRoutes(r)

	// Aplicar o middleware CORS globalmente
	r.Use(enableCORS)

	// Aplicar o middleware de log para todas as requisições
	r.Use(logRequest)

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent) // Retorna 204 No Content
	})

	// Rota para o Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Servir os arquivos estáticos (HTML, JS, CSS) da pasta 'ui'
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui"))))

	// Alterar a rota principal para exibir automaticamente o index.html se o caminho for "/"
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html") // Serve a página principal (formulário + mapa)
	})

	// Configurar resposta global para OPTIONS em qualquer rota
	// r.PathPrefix("/geoconding/search/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// })
	// Configurar resposta global para OPTIONS em qualquer rota
	r.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Iniciar o servidor
	fmt.Println("Servidor iniciado em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
