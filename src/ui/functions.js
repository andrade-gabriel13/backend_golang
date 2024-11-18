
export function searchStreet() {
    // Pega o valor do campo de endereço
    const address = document.getElementById("address").value;

    // Verifica se o campo de endereço está vazio
    if (!address) {
        alert("Por favor, insira um endereço.");
        return;
    }

    document.getElementById("loading").style.display = "block";

    // Faz a requisição para a API com o endereço
    fetch(`http://localhost:8080/deliveries/geoconding/search?endereco=${encodeURIComponent(address)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Erro na requisição');
            }
            return response.json(); // Decodifica a resposta em JSON
        })
        .then(data => {
            console.log('Dados recebidos da API:', data); // Log dos dados recebidos da API

            // Verifique se há resultados para preencher os campos
            if (data && data.address) {
                // Função auxiliar para encontrar o valor de long_name com base no tipo
                const findByType = (type) => {
                    return data.address.find(item => item.types && item.types.includes(type))?.long_name || '';
                };
                    // Preenchendo os campos do formulário com base no tipo e long_name
                const setValueIfExists = (id, value) => {
                    const element = document.getElementById(id);
                    if (element) {
                        element.value = value;
                    }
                };

                const findByTypeAndIndex = (type, index) => {
                    return data.address
                        .filter(item => item.types && item.types.includes(type)) // Filtra os itens pelo tipo
                        .map((item, idx) => ({ item, idx })) // Cria um mapeamento com o índice
                        .find(({ idx }) => idx === index)?.item.long_name || ''; // Busca pelo índice específico
                };


                // Preenchendo os campos do formulário com base no tipo e long_name
                document.getElementById('street').value = findByTypeAndIndex('route', 0); // Rua
                document.getElementById('number').value = findByType('street_number'); // Número
                document.getElementById('neighborhood').value = findByTypeAndIndex('route', 1); // Bairro (se aplicável)
                document.getElementById('complement').value = ''; // Complemento (não presente no JSON)
                document.getElementById('city').value = findByType('locality'); // Cidade
                document.getElementById('state').value = findByType('state'); // Estado
                document.getElementById('country').value = findByType('country'); // País
                
                setValueIfExists('latitude', data.latitude); // Latitude
                setValueIfExists('longitude', data.longitude); // Longitude
            } else {
                console.error('Dados de endereço não encontrados.');
            }

        })
        .catch(error => {
            console.error("Erro ao buscar o endereço:", error);
            alert("Erro ao buscar o endereço.");
        })
        .finally(() => {
        // Esconde o ícone de carregamento após a resposta
        document.getElementById("loading").style.display = "none";
    });

        
        }



export function saveClient() {
    // Captura os valores de todos os campos do formulário
    const clientData = {
        name: document.getElementById("name").value,
        weight_kg: parseFloat(document.getElementById("weight").value),
        address: document.getElementById("address").value,
        street: document.getElementById("street").value,
        number: parseInt(document.getElementById("number").value),
        neighborhood: document.getElementById("neighborhood").value,
        complement: document.getElementById("complement").value,
        city: document.getElementById("city").value,
        state: document.getElementById("state").value,
        country: document.getElementById("country").value,
        latitude: parseFloat(document.getElementById("latitude").value),
        longitude: parseFloat(document.getElementById("longitude").value)
    };

    // Envia os dados do cliente para a API para salvamento no banco de dados
    fetch("/deliveries", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(clientData) // Converte o objeto para JSON
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error("Erro ao salvar cliente: " + response.status);
        }
    })
    .then(data => {
        alert("Cliente salvo com sucesso!");
        // Limpa os campos do formulário após o sucesso
        document.querySelectorAll(".form-control").forEach(input => input.value = '');
    })
    .catch(error => {
        console.error("Erro ao salvar cliente:", error);
        alert("Erro ao salvar os dados. Verifique as informações e tente novamente.");
    });
}

// Anexando funções ao escopo global
window.searchStreet = searchStreet;
window.searchStreet = saveClient;
