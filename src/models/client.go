package models

import (
	"log/slog"
	"time"

	"github.com/jinzhu/gorm"
)

// Client representa o modelo de um cliente, utilizado no banco de dados.
// Ele define todos os campos necessários para armazenar informações de um cliente.
// Cada cliente possui informações como nome, endereço, peso, localização geográfica, etc.
// A estrutura utiliza GORM como ORM para a manipulação dos dados.
// Tipos de dados:
// - Nome, endereço e dados geográficos são strings ou floats.
// - O ID é gerado automaticamente por GORM.
// A tabela associada a este modelo no banco de dados é chamada "clients".
type Client struct {
	gorm.Model
	Name         string  `json:"name"`         // Nome do cliente
	WeightKg     float64 `json:"weight_kg"`    // Peso do cliente em kg
	Address      string  `json:"address"`      // Endereço completo
	Street       string  `json:"street"`       // Nome da rua
	Number       int     `json:"number"`       // Número da residência
	Neighborhood string  `json:"neighborhood"` // Bairro do cliente
	Complement   string  `json:"complement"`   // Complemento do endereço
	City         string  `json:"city"`         // Cidade do cliente
	State        string  `json:"state"`        // Estado do cliente
	Country      string  `json:"country"`      // País do cliente
	Latitude     float64 `json:"latitude"`     // Latitude da localização
	Longitude    float64 `json:"longitude"`    // Longitude da localização
}

// ClientUpdate representa um cliente com os campos atualizáveis.
// Ele é utilizado para receber dados de atualização, incluindo o ID do cliente.
// Esse modelo é usado quando os dados de um cliente existente precisam ser atualizados no banco de dados.
// Ele possui os mesmos campos que o modelo Client, mas inclui um campo adicional para o ID do cliente.
type ClientUpdate struct {
	gorm.Model
	ID           uint    `json:"id"`           // ID do cliente, necessário para a atualização
	Name         string  `json:"name"`         // Nome do cliente
	WeightKg     float64 `json:"weight_kg"`    // Peso do cliente em kg
	Address      string  `json:"address"`      // Endereço completo
	Street       string  `json:"street"`       // Nome da rua
	Number       int     `json:"number"`       // Número da residência
	Neighborhood string  `json:"neighborhood"` // Bairro do cliente
	Complement   string  `json:"complement"`   // Complemento do endereço
	City         string  `json:"city"`         // Cidade do cliente
	State        string  `json:"state"`        // Estado do cliente
	Country      string  `json:"country"`      // País do cliente
	Latitude     float64 `json:"latitude"`     // Latitude da localização
	Longitude    float64 `json:"longitude"`    // Longitude da localização
}

type ArchivedClient struct {
	ID           int    `gorm:"primaryKey"`
	Name         string `gorm:"size:255"`
	WeightKg     float64
	Address      string `gorm:"size:255"`
	Street       string `gorm:"size:255"`
	Number       int    `gorm:"size:50"`
	Neighborhood string `gorm:"size:255"`
	Complement   string `gorm:"size:255"`
	City         string `gorm:"size:255"`
	State        string `gorm:"size:255"`
	Country      string `gorm:"size:255"`
	Latitude     float64
	Longitude    float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

// Struct auxiliar para garantir a ordem dos campos
type ClientResponse struct {
	gorm.Model
	Name         string  `json:"name"`         // Nome do cliente
	WeightKg     float64 `json:"weight_kg"`    // Peso do cliente em kg
	Address      string  `json:"address"`      // Endereço completo
	Street       string  `json:"street"`       // Nome da rua
	Number       int     `json:"number"`       // Número da residência
	Neighborhood string  `json:"neighborhood"` // Bairro do cliente
	Complement   string  `json:"complement"`   // Complemento do endereço
	City         string  `json:"city"`         // Cidade do cliente
	State        string  `json:"state"`        // Estado do cliente
	Country      string  `json:"country"`      // País do cliente
	Latitude     float64 `json:"latitude"`     // Latitude da localização
	Longitude    float64 `json:"longitude"`    // Longitude da localização
}

type GeocodingResponse struct {
	Status  string   `json:"status"`
	Results []Result `json:"results"`
}

type Result struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          Geometry           `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	PlusCode          PlusCode           `json:"plus_code"`
	Types             []string           `json:"types"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Location     Location `json:"location"`
	LocationType string   `json:"location_type"`
	Viewport     Viewport `json:"viewport"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Viewport struct {
	Northeast Location `json:"northeast"`
	Southwest Location `json:"southwest"`
}

type PlusCode struct {
	// Campos do PlusCode (se existirem)
}

// TableName define explicitamente o nome da tabela associada ao modelo Client no banco de dados.
// GORM utiliza essa função para mapear o modelo para a tabela correta no banco de dados.
func (Client) TableName() string {
	slog.Info("Definindo nome da tabela para modelo Client", slog.String("table_name", "clients"))
	return "clients"
}
