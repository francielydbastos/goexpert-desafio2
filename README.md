# Multithread CEP API Race Condition

Aplicação em Go que consulta um CEP em duas APIs simultaneamente e usa apenas a resposta mais rapida.

## Requisitos atendidos

- Requisições paralelas para `BrasilAPI` e `ViaCEP`
- Vence a primeira resposta válida
- A requisição perdedora é cancelada via `context`
- Timeout global de `1s`
- Saída no terminal com os dados do endereço e nome da API vencedora

## APIs usadas

- `https://brasilapi.com.br/api/cep/v1/{cep}`
- `http://viacep.com.br/ws/{cep}/json/`

## Como executar

```powershell
go run . 01001000
```

Também aceita CEP com hífen:

```powershell
go run . 01001-000
```

## Exemplo de saída

```text
API vencedora: ViaCEP
CEP: 01001-000
Logradouro: Praca da Sé
Bairro: Sé
Cidade: São Paulo
UF: SP
IBGE: 3550308
DDD: 11
```

## Timeout

Se nenhuma API responder em até 1s, a aplicação encerra a execução com o erro abaixo:

```text
erro: timeout de 1s atingido. Nenhuma API respondeu a tempo
```


