# Korp_Teste_JoaoVictor

Sistema full stack para controle de estoque, faturamento e insights com IA.

## Stack

- Frontend: Angular 21, Signals, RxJS e Angular Material
- Backend: Go, Gin, GORM e PostgreSQL
- Autenticacao: JWT
- IA: Gemini via SDK `google.golang.org/genai`
- Documentacao: Swagger

## Arquitetura

O backend foi separado em dois microsservicos:

- `stock-service`: cadastro de produtos, controle de estoque e insights com IA
- `billing-service`: cadastro e fechamento de notas fiscais

Ambos compartilham autenticacao via JWT e se comunicam por HTTP quando o fluxo de faturamento precisa interagir com estoque.

## Principais funcionalidades

- cadastro e autenticacao de usuarios
- cadastro, listagem, edicao e exclusao de produtos
- controle de baixa de estoque
- cadastro, listagem, edicao, exclusao e fechamento de notas fiscais
- tela de insights com IA no frontend
- recomendacoes de produtos em alta com base no contexto do estoque
- documentacao Swagger por servico

## Requisitos

- Docker e Docker Compose
- Node.js LTS se quiser rodar o frontend fora do Docker
- Go 1.26+ se quiser rodar os servicos manualmente
- PostgreSQL acessivel pela `DATABASE_URL`
- chave do Gemini para testar os insights com IA

## Configuracao

As configuracoes do backend ficam em `backend/.env`.

Campos importantes:

- `DATABASE_URL`
- `JWT_SECRET`
- `AUTO_MIGRATE`
- `GEMINI_API_KEY`
- `GEMINI_MODEL`
- `GEMINI_TIMEOUT_MS`
- `AI_LOW_STOCK_THRESHOLD`

## Como executar

### Backend com Docker

Na pasta `backend`:

```bash
docker compose up --build
```

Servicos:

- `http://localhost:8081` -> `stock-service`
- `http://localhost:8082` -> `billing-service`

### Frontend

Na pasta `frontend`:

```bash
npm install
npm start
```

Aplicacao:

- `http://localhost:4200`

Observacao:
Se houver problema de build com Node muito novo, prefira uma versao LTS como `20` ou `22`.

## Fluxo rapido

1. Crie uma conta ou faca login
2. Cadastre produtos no estoque
3. Crie notas fiscais com os itens desejados
4. Feche notas para baixar estoque
5. Acesse a tela de insights para consultar a analise com IA

## Endpoints uteis

### Auth

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

### Produtos

- `POST /api/v1/products`
- `GET /api/v1/products`
- `GET /api/v1/products/:code`
- `PUT /api/v1/products/:code`
- `PATCH /api/v1/products/:code/decrease`
- `DELETE /api/v1/products/:code`

### Notas fiscais

- `POST /api/v1/invoices`
- `GET /api/v1/invoices`
- `GET /api/v1/invoices/:number`
- `PUT /api/v1/invoices/:number`
- `PATCH /api/v1/invoices/:number/close`
- `DELETE /api/v1/invoices/:number`

### IA

- `GET /api/v1/ai/insights`

Todos os endpoints `/api/v1/*`, exceto autenticacao, exigem:

```txt
Authorization: Bearer <token>
```

## Swagger

- Stock: `http://localhost:8081/swagger/index.html`
- Billing: `http://localhost:8082/swagger/index.html`

## Concorrencia e resiliencia

O projeto possui algumas protecoes importantes:

- baixa de estoque atomica no banco
- fechamento condicional de nota para reduzir processamento duplicado
- `Circuit Breaker` na comunicacao entre `billing-service` e `stock-service`

Isso melhora a consistencia do fluxo principal e a tolerancia a falhas entre servicos.

## IA no projeto

Os insights com IA analisam:

- situacao atual do estoque
- alertas operacionais
- observacoes sobre faturamento
- recomendacoes de compra com base em produtos em alta

Para testar essa funcionalidade, basta configurar a `GEMINI_API_KEY` no `backend/.env` e acessar a tela de insights no frontend.
