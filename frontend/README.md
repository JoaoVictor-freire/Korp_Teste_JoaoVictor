# Frontend Angular

Aplicacao Angular para consumir o backend em Go:

- `stock-service`: `http://localhost:8081`
- `billing-service`: `http://localhost:8082`

## Rodar

```bash
npm install
npm start
```

A aplicacao sobe em `http://localhost:4200`.

## Fluxo

1. Crie uma conta em `/auth`
2. Faca login para receber o JWT
3. Cadastre produtos
4. Crie notas fiscais com itens

## Observacao

Se a build do Angular falhar localmente com Node `v25`, prefira usar uma versao LTS par, como Node `20` ou `22`.
