# Challenge para Desenvolvedor Back End - Golang
## Engenheiro Responsável
## Gabriel Andrade
# Projeto de Validação de Campos de Cliente

Este projeto fornece uma série de funcionalidades para validar e processar dados relacionados aos clientes. Ele foi desenvolvido utilizando Go e visa garantir que os dados fornecidos para criação e atualização de clientes sejam válidos e estejam no formato correto antes de serem persistidos no banco de dados.

## Funcionalidades

### 1. **Validação de Campos Comuns de Cliente**

A função `ValidateCommonClientFields` valida os campos essenciais de um cliente. Ela verifica se todos os campos obrigatórios estão presentes e se os valores são válidos (por exemplo, se um campo numérico não é zero ou negativo). Caso algum campo esteja ausente ou inválido, a função retorna um erro e loga as informações de falha, como o nome do campo e o valor que causou o erro.

#### Campos Validados:
- **Name**: Não pode ser vazio.
- **TestTechnical**: Não pode ser vazio.
- **WeightKg**: Deve ser maior que 0.
- **Address**: Não pode ser vazio.
- **Street**: Não pode ser vazio.
- **Number**: Deve ser maior que 0.
- **Neighborhood**: Não pode ser vazio.
- **City**: Não pode ser vazio.
- **State**: Não pode ser vazio.
- **Country**: Não pode ser vazio.
- **Latitude**: Deve ser um valor válido (diferente de 0).
- **Longitude**: Deve ser um valor válido (diferente de 0).

### 2. **Criação de Cliente e Validação**

A função `CreateClientCheckValues` utiliza a função `ValidateCommonClientFields` para garantir que os dados do cliente estão corretos antes de prosseguir. Se a validação for bem-sucedida, ela retorna um mapa com um status de sucesso e uma mensagem.

**Status de Retorno:**
- `"status": "valid"`
- `"message": "Client validated successfully"`

### 3. **Validação de Atualização de Cliente**

A função `ValidateClientUpdate` é usada quando um cliente já existe e há a necessidade de atualizar seus dados. Ela valida se o ID do cliente está presente e é válido (maior que 0). Após essa validação, a atualização dos dados será realizada na persistência (banco de dados).

**Status de Retorno:**
- `"status": "valid"`
- `"message": "Client validated for update"`

Se o ID não for válido, retorna o erro:

```txt
"ID do cliente é obrigatório para atualização"
```

## Logs de Validação

Durante o processo de validação, os logs são utilizados para registrar falhas ou sucessos nas validações. O log será registrado com a severidade apropriada:

- **Logs de erro**: Caso algum campo esteja ausente ou com um valor inválido, um erro será logado com as informações detalhadas do campo problemático.
- **Logs de sucesso**: Caso a validação seja bem-sucedida, um log informando a conclusão bem-sucedida da validação será gerado.

### Exemplos de Log:

- **Erro**:  
  Caso o campo `Name` não seja preenchido, o log gerado será:

  ```txt
  Missing required field, field: name, value: 
  ```

- **Sucesso**:  
  Após a validação bem-sucedida de todos os campos, o log gerado será:

  ```txt
  Client validated successfully, field: validation, status: success
  ```

## Estrutura do Projeto

O projeto é composto por duas funções principais para a validação de dados de clientes:

1. **`ValidateCommonClientFields`**: Valida os campos essenciais de um cliente.
2. **`CreateClientCheckValues`**: Chama `ValidateCommonClientFields` e retorna um status de validação bem-sucedida.
3. **`ValidateClientUpdate`**: Valida a atualização dos dados de um cliente, verificando a existência de um ID válido.

Essas funções podem ser reutilizadas de forma flexível em diferentes partes do seu sistema para garantir a integridade dos dados de clientes, seja na criação ou atualização de informações.

## Como Utilizar

### Dependências

Este projeto foi desenvolvido com a linguagem Go e um container em docker do banco de dados Mysql. Certifique-se de ter o Go instalado em sua máquina.

## 1. Instalar o Go:  
   Se você ainda não tiver o Go instalado, siga as instruções em [https://golang.org/doc/install](https://golang.org/doc/install).

## 2. Instalar dependências:  
   Para instalar as dependências do seu projeto, navegue até o diretório do projeto e execute:

   ```bash
   go mod tidy
   ```
e o docker com o devido container


## 3. Configurar o Docker com o Docker Compose

Este projeto utiliza o Docker Compose para configurar e executar o banco de dados MySQL. Siga as instruções abaixo para configurar o container do MySQL.

### Passo 1: Instalar o Docker

Antes de começar a configurar o Docker Compose, você precisa ter o Docker instalado em seu sistema. O Docker permite a criação, o gerenciamento e a execução de containers, proporcionando um ambiente de desenvolvimento isolado.

#### Para sistemas **Linux** (Ubuntu):

1. Atualize o repositório de pacotes:
   ```bash
   sudo apt update
   ```

2. Instale pacotes necessários para permitir o uso de repositórios `https` e o repositório oficial do Docker:
   ```bash
   sudo apt install apt-transport-https ca-certificates curl software-properties-common
   ```

3. Baixe a chave oficial do Docker:
   ```bash
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
   ```

4. Adicione o repositório Docker ao APT sources:
   ```bash
   sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
   ```

5. Atualize novamente o repositório de pacotes:
   ```bash
   sudo apt update
   ```

6. Instale o Docker:
   ```bash
   sudo apt install docker-ce
   ```

7. Verifique se o Docker foi instalado corretamente:
   ```bash
   sudo docker --version
   ```

#### Para sistemas **Windows**:

1. Baixe o **Docker Desktop** no link [https://desktop.docker.com](https://desktop.docker.com).
2. Execute o instalador e siga os passos de instalação.
3. Após a instalação, inicie o Docker Desktop e verifique se o Docker está rodando corretamente no seu sistema.

#### Para sistemas **macOS**:

1. Baixe o **Docker Desktop** no link [https://desktop.docker.com](https://desktop.docker.com).
2. Abra o arquivo `.dmg` e siga as instruções para instalar o Docker Desktop.
3. Após a instalação, inicie o Docker Desktop e verifique se o Docker está funcionando corretamente.

---

### Passo 2: Instalar o Docker Compose

Docker Compose é uma ferramenta que permite definir e executar aplicativos multi-container com o Docker. É particularmente útil quando seu projeto envolve mais de um container, como no caso do MySQL.

#### Para sistemas **Linux**:

1. Baixe a versão mais recente do Docker Compose:
   ```bash
   sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   ```

2. Torne o Docker Compose executável:
   ```bash
   sudo chmod +x /usr/local/bin/docker-compose
   ```

3. Verifique se o Docker Compose foi instalado corretamente:
   ```bash
   docker-compose --version
   ```

#### Para sistemas **Windows** e **macOS**:

O Docker Compose já está incluído no Docker Desktop para Windows e macOS, então você não precisa instalar separadamente. Basta verificar se o Docker Compose está funcionando com o comando:

```bash
docker-compose --version
```

---

### Passo 3: Configurar o Docker Compose com MySQL

Agora que você tem o Docker e o Docker Compose instalados, o próximo passo é configurar o serviço do MySQL usando o Docker Compose. O arquivo `docker-compose.yml` vai definir as configurações do container, como a imagem do MySQL, variáveis de ambiente, volumes e rede.

1. **Criar o arquivo `docker-compose.yml`**

Crie o arquivo `docker-compose.yml` na raiz do seu projeto com o seguinte conteúdo:

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: mysql
    environment:
      - MYSQL_ROOT_PASSWORD=golang # Senha do root
      - MYSQL_DATABASE=golang     # Nome do banco de dados
      - MYSQL_USER=golang         # Usuário para o banco de dados
      - MYSQL_PASSWORD=golang     # Senha do usuário
    ports:
      - "3306:3306"  # Mapeia a porta do MySQL para a porta do host
    volumes:
      - mysql-data:/var/lib/mysql  # Persiste dados em um volume local
    networks:
      - app-network  # Conecta o serviço à rede do aplicativo

volumes:
  mysql-data:
    driver: local  # Definição do volume local para persistência dos dados

networks:
  app-network:
    driver: bridge  # Define a rede de comunicação entre os containers
```

2. **Subir o Docker Compose**

Após criar o arquivo `docker-compose.yml`, o próximo passo é subir os containers definidos nele. No terminal, execute o comando:

```bash
docker-compose up -d
```

Esse comando vai iniciar o container MySQL em segundo plano (`-d` para detached). O Docker Compose vai baixar a imagem do MySQL (se necessário) e configurar o banco de dados conforme as variáveis de ambiente definidas no arquivo `docker-compose.yml`.

3. **Verificar se o MySQL está em execução**

Para verificar se o MySQL está rodando, use o seguinte comando:

```bash
docker ps
```

Isso irá listar todos os containers em execução. O container MySQL deve aparecer na lista.

4. **Conectar ao MySQL**

Agora que o MySQL está em funcionamento, você pode conectar-se ao banco de dados usando um cliente MySQL de sua preferência. Utilize as credenciais definidas no arquivo `docker-compose.yml`:

- **Host**: `localhost`
- **Porta**: `3306`
- **Usuário**: `golang`
- **Senha**: `golang`
- **Banco de Dados**: `golang`

Por exemplo, utilizando o cliente MySQL via terminal:

```bash
mysql -h localhost -u golang -p
```

Quando solicitado, insira a senha `golang`.

5. **Parar o Docker Compose**

Quando terminar de usar o ambiente Docker, pare os containers com o seguinte comando:

```bash
docker-compose down
```

Isso irá parar os containers e remover as redes e volumes criados durante a execução.

---

Com esses passos, você terá configurado o Docker, Docker Compose e o MySQL para o seu projeto. Isso permite um ambiente de desenvolvimento limpo, isolado e facilmente replicável por outros desenvolvedores.



### Execução do teste em Go

Para testar a validação de dados do cliente, você pode criar uma função de teste que utilize as funções de validação descritas. Exemplo de código para testar a função de validação:

```go
package main

import (
	"fmt"
	"myapi/services"
	"myapi/models"
)

func main() {
	client := models.Client{
		Name:         "",
		TestTechnical: "TechnicalTest123",
		WeightKg:     75,
		Address:      "123 Main St",
		Street:       "Main St",
		Number:       10,
		Neighborhood: "Downtown",
		City:         "Metropolis",
		State:        "StateX",
		Country:      "CountryY",
		Latitude:     0,
		Longitude:    0,
	}

	// Testando a validação
	if err := services.ValidateCommonClientFields(client); err != nil {
		fmt.Println("Erro de validação:", err)
	} else {
		fmt.Println("Cliente validado com sucesso")
	}
}
```

### Saída Esperada:

Se os campos não forem preenchidos corretamente, você verá as mensagens de erro no console.

Exemplo de saída:

```txt
Erro de validação: name is required
```

## Contribuição

Se você deseja contribuir com o projeto, fique à vontade para enviar pull requests. Certifique-se de que o código esteja bem testado e com a formatação correta. Para contribuir, siga estas etapas:

1. Faça um fork do repositório.
2. Crie uma nova branch para a sua funcionalidade.
3. Faça as alterações necessárias.
4. Envie um pull request com uma descrição clara das mudanças.

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).
