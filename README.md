# Multithread CEP API Race

Aplicacao em Go que consulta um CEP em duas APIs simultaneamente e usa apenas a resposta mais rapida.

## Requisitos atendidos

- Requisicoes paralelas para `BrasilAPI` e `ViaCEP`
- Vence a primeira resposta válida (race)
- A requisição perdedora é cancelada via `context`
- Timeout global de `1s`
- Saida no terminal com os dados do endereço e nome da API vencedora

## APIs usadas

- `https://brasilapi.com.br/api/cep/v1/{cep}`
- `http://viacep.com.br/ws/{cep}/json/`

## Como executar

```powershell
go run . 01001000
```

Também aceita CEP com hifen:

```powershell
go run . 01001-000
```

## Exemplo de saida

```text
API vencedora: ViaCEP
CEP: 01001-000
Logradouro: Praca da Se
Bairro: Se
Cidade: Sao Paulo
UF: SP
IBGE: 3550308
DDD: 11
```

## Timeout

Se nenhuma API responder em até 1s, a aplicação encerra a execução com o erro abaixo:

```text
erro: timeout de 1s atingido. Nenhuma API respondeu a tempo
```


