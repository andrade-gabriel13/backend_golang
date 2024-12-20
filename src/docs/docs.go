// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/deliveries": {
            "get": {
                "description": "Retorna uma lista de clientes paginada, ou um cliente específico se o ID for fornecido, permitindo definir o limite, offset, cidade e ID.",
                "tags": [
                    "deliveries"
                ],
                "summary": "Busca clientes com paginação e filtros de cidade e ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do cliente para busca específica",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 100,
                        "description": "Número máximo de clientes por página",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Número de registros a pular antes de começar a listar os clientes",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Cidade para filtrar os clientes",
                        "name": "city",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Dados da lista de clientes com metadados de paginação",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Erro ao buscar clientes",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Modifica os dados de um cliente existente.",
                "tags": [
                    "deliveries"
                ],
                "summary": "Atualiza um cliente",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do cliente para busca específica",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "description": "Cliente para atualizar",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ClientUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Resposta com os dados do cliente atualizado",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Corpo da requisição inválido",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Cliente não encontrado",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erro ao atualizar cliente",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Recebe um JSON contendo os dados de um cliente e insere o registro no sistema.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "deliveries"
                ],
                "summary": "Cria um novo cliente",
                "parameters": [
                    {
                        "description": "Dados do cliente para criação",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Client"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Cliente criado com sucesso",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Requisição inválida: erro no corpo da requisição ou JSON malformado",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erro ao criar cliente no banco de dados",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Exclui todos os clientes ou um cliente específico pelo ID.",
                "tags": [
                    "deliveries"
                ],
                "summary": "Excluir clientes",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "Excluir todos os clientes (true para excluir todos os clientes)",
                        "name": "deleteAll",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "ID do cliente a ser excluído (se deleteAll não for especificado)",
                        "name": "id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Mensagem de sucesso",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Parâmetros inválidos",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erro interno do servidor",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/deliveries/geoconding/search": {
            "get": {
                "description": "Retorna as coordenadas geográficas (latitude e longitude) de um endereço fornecido.",
                "tags": [
                    "geocoding"
                ],
                "summary": "Busca latitude e longitude de um endereço",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Endereço a ser consultado",
                        "name": "endereco",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Resposta com latitude e longitude do endereço",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "number"
                            }
                        }
                    },
                    "400": {
                        "description": "Parâmetro 'endereco' ausente ou inválido",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Erro ao consultar a API de geocoding",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Client": {
            "type": "object",
            "properties": {
                "address": {
                    "description": "Endereço completo",
                    "type": "string"
                },
                "city": {
                    "description": "Cidade do cliente",
                    "type": "string"
                },
                "complement": {
                    "description": "Complemento do endereço",
                    "type": "string"
                },
                "country": {
                    "description": "País do cliente",
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "latitude": {
                    "description": "Latitude da localização",
                    "type": "number"
                },
                "longitude": {
                    "description": "Longitude da localização",
                    "type": "number"
                },
                "name": {
                    "description": "Nome do cliente",
                    "type": "string"
                },
                "neighborhood": {
                    "description": "Bairro do cliente",
                    "type": "string"
                },
                "number": {
                    "description": "Número da residência",
                    "type": "integer"
                },
                "state": {
                    "description": "Estado do cliente",
                    "type": "string"
                },
                "street": {
                    "description": "Nome da rua",
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "weight_kg": {
                    "description": "Peso do cliente em kg",
                    "type": "number"
                }
            }
        },
        "models.ClientUpdate": {
            "type": "object",
            "properties": {
                "address": {
                    "description": "Endereço completo",
                    "type": "string"
                },
                "city": {
                    "description": "Cidade do cliente",
                    "type": "string"
                },
                "complement": {
                    "description": "Complemento do endereço",
                    "type": "string"
                },
                "country": {
                    "description": "País do cliente",
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "description": "ID do cliente, necessário para a atualização",
                    "type": "integer"
                },
                "latitude": {
                    "description": "Latitude da localização",
                    "type": "number"
                },
                "longitude": {
                    "description": "Longitude da localização",
                    "type": "number"
                },
                "name": {
                    "description": "Nome do cliente",
                    "type": "string"
                },
                "neighborhood": {
                    "description": "Bairro do cliente",
                    "type": "string"
                },
                "number": {
                    "description": "Número da residência",
                    "type": "integer"
                },
                "state": {
                    "description": "Estado do cliente",
                    "type": "string"
                },
                "street": {
                    "description": "Nome da rua",
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "weight_kg": {
                    "description": "Peso do cliente em kg",
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "API - Golang Desafio Backend",
	Description:      "API de exemplo com integração Swagger, oferecendo uma interface interativa para documentação e testes de endpoints. Esta API foi desenvolvida para ilustrar boas práticas de documentação, utilizando Swagger para facilitar a navegação e o entendimento das operações disponíveis, parâmetros de entrada, respostas e exemplos de uso. Indicada para desenvolvimento e demonstração de integração entre APIs e clientes.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
