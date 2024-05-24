# Market Vibranium

Market Vibranium é um serviço de gerenciamento de ordens de alta performance, projetado para facilitar a compra e venda de vibraniums. Este serviço opera totalmente em memória para atingir a meta de processar 5000 requisições por segundo.

## Diagrama da aplicacao
Se preferir, [o diagrama demonstra](docs/diagram.png) os passos que é feito para cada rota e processo dentro da aplicação

## Índice

1. [Visão Geral](#visao-geral)
2. [Funcionalidades](#funcionalidades)
3. [Instalação](#instalacao)
4. [Uso](#uso)
5. [Endpoints](#endpoints)
    - [Criar Ordem](#criar-ordem)
    - [Obter Ordem](#obter-ordem)
    - [Criar Carteira](#criar-carteira)
    - [Depositar na Carteira](#depositar-na-carteira)
    - [Obter Carteira](#obter-carteira)
    - [Métricas](#metricas)
6. [Persistência de Dados](#persistencia-de-dados)
7. [Monitoramento com Grafana e Prometheus](#monitoramento-com-grafana-e-prometheus)
    - [Acessando o Prometheus](#acessando-o-prometheus)
    - [Acessando o Grafana](#acessando-o-grafana)
    - [Métricas Disponíveis](#métricas-disponíveis)

## Visão Geral

Market Vibranium é um serviço de ordens de vendas onde você pode comprar e vender vibraniums de/para outros usuários. Inicialmente, é necessário criar uma carteira onde será registrado o usuário, e essa carteira indicará quantos vibraniums e reais você possui. É preciso chamar uma rota HTTP para depositar dinheiro ou vibraniums na carteira antes de poder comprar ou vender.

## Funcionalidades

- **Operações em memória**: Todas as operações são realizadas em memória.
- **Persistência de dados**: Utiliza snapshots para salvar o estado das ordens e carteiras durante o graceful shutdown.
- **Métricas**: Integração com Prometheus para monitoramento.
- **Processamento em ordem**: Os pedidos de compra/venda são processados utilizando channels do Go para processamento na mesma ordem que foi recebido.

## Instalação

Para instalar e rodar o projeto localmente, siga os passos abaixo:

1. Clone o repositório:
    ```sh
    git clone https://github.com/HunnTeRUS/vibranium-market-ml
    cd market-vibranium
    ```

2. Rodar o projeto com Docker:
    ```sh
    docker-compose up --build
    ```

3. Configure as variáveis de ambiente conforme necessário, onde hoje temos:
    4. **SNAPSHOT_DIR:** Utilizada para guardar os snapshots das mensagens restantes nas filas de processamentos (channels),

    5. **WALLETS_SNAPSHOT_FILE:** Utilizada para indicar em qual arquivo salvar as wallets existentes dos clientes

    6. **ORDERS_SNAPSHOT_FILE:** Utilizada para guardar as ordens de compra e venda existentes no serviço atualmente
    7. **LOG_OUTPUT:** Utilizado para caso você queira mudar a localização de onde os logs vão ser salvos (default "stdout").
    8. **LOG_LEVEL:** Utilizado para mudar o nível de logs que você quer registrado (default "info").

## Uso

Para usar o serviço, siga as instruções abaixo para interagir com os endpoints disponíveis.

## Endpoints

### Criar Ordem

- **Endpoint**: `/orders`
- **Método**: POST
- **Descrição**: Cria uma nova ordem.
- **Corpo da Requisição**:
    ```json
    {
        "userId": "UUID string",
        "type": 1 | 2,
        "amount": <quantidade de vibranium em int64>,
        "price": <valor à ser pago/vendido por cada vibranium em float64>
    }
    ```
- **Resposta**: `201 Created` em caso de sucesso.

### Obter Ordem

- **Endpoint**: `/orders/:id`
- **Método**: GET
- **Descrição**: Obtém os detalhes de uma ordem específica.
- **Resposta**: Detalhes da ordem.

### Criar Carteira

- **Endpoint**: `/wallets`
- **Método**: POST
- **Descrição**: Cria uma nova carteira para um usuário.
- **Resposta**: Detalhes da carteira criada.
    ```json
    {
      "user_id": "a5d40e13-95fb-4be5-9c69-bfda40e41267",
      "balance": 0,
      "vibranium": 0
    }
    ```

### Depositar na Carteira

- **Endpoint**: `/wallets/deposit`
- **Método**: POST
- **Descrição**: Deposita dinheiro ou vibraniums na carteira de um usuário.
- **Corpo da Requisição**:
    ```json
    {
        "userId": "string",
        "amount": "float64",
        "vibranium": "int64"
    }
    ```
- **Resposta**: Detalhes da carteira atualizada.

### Obter Carteira

- **Endpoint**: `/wallets/:userId`
- **Método**: GET
- **Descrição**: Obtém os detalhes da carteira de um usuário específico.
- **Resposta**: Detalhes da carteira.
    ```json
    {
        "userId": "a5d40e13-95fb-4be5-9c69-bfda40e41267",
        "amount": 0,
        "vibranium": 0
    }
    ```

### Métricas

- **Endpoint**: `/metrics`
- **Método**: GET
- **Descrição**: Obtém as métricas da aplicação para monitoramento via Prometheus.

## Persistência de Dados

Para garantir que as ordens e carteiras não sejam perdidas, o projeto implementa um mecanismo de graceful shutdown que salva snapshots do estado atual da aplicação em arquivos. Esses arquivos são recarregados quando a aplicação é iniciada novamente.


## Monitoramento com Grafana e Prometheus

O projeto Market Vibranium utiliza Prometheus para coleta de métricas e Grafana para visualização dessas métricas. Ambas as ferramentas já estão configuradas e instaladas através do Docker Compose. A seguir, estão as instruções para acessar e utilizar essas ferramentas.

### Acessando o Prometheus

1. **Inicie o Docker Compose**:
   Certifique-se de que todos os serviços estão em execução:
    ```sh
    docker-compose up -d
    ```

2. **Acesse o Prometheus**:
   Abra um navegador e vá para `http://localhost:9090`.

### Acessando o Grafana

1. **Acesse o Grafana**:
   Abra um navegador e vá para `http://localhost:3000`.

2. **Faça login no Grafana**:
   Use as credenciais padrão:
    - Usuário: `admin`
    - Senha: `admin`

3. **Acesse o dashboard**:
    4. Acesse a aba Dashboards e em seguida selecione **Order Processing and Go Performance Dashboard**

### Métricas Disponíveis

- `processed_orders_total`: Total de ordens processadas.
- `pending_orders_total`: Total de ordens pendentes.
- `canceled_orders_total`: Total de ordens canceladas.
- `queue_length`: Tamanho da fila de ordens.
- `order_processing_duration_seconds`: Duração do processamento de ordens.