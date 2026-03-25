# Order Management System

Sistema de gerenciamento de pedidos construído em Go com arquitetura limpa, processamento assíncrono e logs estruturados.

## Tecnologias

- **Go** — linguagem principal
- **Gin** — framework HTTP
- **MongoDB** — banco de dados
- **RabbitMQ** — mensageria assíncrona
- **Zap** — logs estruturados
- **Docker** — containerização

## Arquitetura

O projeto segue os princípios de **Clean Architecture**, com separação clara de responsabilidades entre as camadas:

```
cmd/
  api/        → entrypoint da API HTTP
  worker/     → entrypoint do worker assíncrono

internal/
  domain/         → entidades de negócio (Order, Customer, Payment...)
  usecase/        → regras de negócio
  repository/     → interfaces de repositório
  delivery/http/  → handlers HTTP (Gin)
  infrastructure/
    database/     → conexão com MongoDB
    messaging/    → publisher e consumer RabbitMQ
    repository/   → implementação MongoDB dos repositórios
  worker/         → processamento assíncrono de mensagens
  logger/         → logger centralizado (Zap)
  events/         → definição de eventos
```

### Decisões técnicas

**MongoDB** foi escolhido por se adequar bem ao domínio de pedidos — um pedido contém itens, descontos, impostos e dados de pagamento, que naturalmente se representam como um documento aninhado. Isso elimina a necessidade de múltiplas tabelas e JOINs que seriam necessários no modelo relacional.

**RabbitMQ** é utilizado para desacoplar a criação do pedido do seu processamento. Ao criar um pedido, um evento `order.created` é publicado e consumido assincronamente pelo worker, garantindo que a API responda rapidamente sem bloquear no processamento.

**Dead Letter Queue (DLQ)** foi implementada para garantir que mensagens que falham repetidamente não sejam perdidas. Após 3 tentativas, a mensagem é movida para a fila `order.created.dlq` para análise posterior.

## Como rodar

### Pré-requisitos

- Go 1.23+
- Docker e Docker Compose

### 1. Clone o repositório

```bash
git clone https://github.com/matheusleo17/order-management-system.git
cd order-management-system
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
```

O arquivo `.env.example` já contém os valores padrão para desenvolvimento local.

### 3. Suba a infraestrutura

```bash
docker-compose up -d
```

Isso irá subir o MongoDB e o RabbitMQ.

### 4. Rode a API

```bash
go run cmd/api/main.go
```

### 5. Rode o Worker (em outro terminal)

```bash
go run cmd/worker/main.go
```

A API estará disponível em `http://localhost:8080`.

## Endpoints

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | /health | Healthcheck da aplicação |
| POST | /orders | Criar pedido |
| GET | /orders | Listar pedidos |
| GET | /orders/:id | Buscar pedido por ID |
| PUT | /orders/:id | Atualizar pedido |
| DELETE | /orders/:id | Deletar pedido |

### Exemplo de criação de pedido

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer": {
      "id": "customer-1",
      "name": "João Silva",
      "email": "joao@email.com"
    },
    "items": [
      {
        "product_id": "prod-1",
        "product": "Notebook",
        "quantity": 1,
        "price": 4500.00
      }
    ],
    "payment": {
      "method": "credit_card",
      "status": "pending"
    },
    "shipment": {
      "address": "Rua das Flores, 123",
      "city": "São Paulo",
      "state": "SP",
      "zip_code": "01310-100"
    }
  }'
```

## Fluxo assíncrono

```
POST /orders
    │
    ▼
CreateOrderUseCase
    │
    ├── Salva no MongoDB
    │
    └── Publica evento order.created no RabbitMQ
                │
                ▼
            Worker consome o evento
                │
                ├── Processa o pedido
                ├── Em caso de falha → retry (até 3x)
                └── Após 3 falhas → DLQ (order.created.dlq)
```

## Testes

```bash
go test ./...
```

Os testes cobrem todos os use cases com mocks da camada de repositório, sem necessidade de infraestrutura rodando.

## Variáveis de ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `APP_ENV` | Ambiente (`development` ou `production`) | `development` |
| `MONGO_URI` | URI de conexão com MongoDB | `mongodb://localhost:27017` |
| `MONGO_DB` | Nome do banco de dados | `order_management` |
| `RABBITMQ_URI` | URI de conexão com RabbitMQ | `amqp://guest:guest@localhost:5672/` |
