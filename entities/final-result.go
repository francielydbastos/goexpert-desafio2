package entities

type Address struct {
	CEP         string
	Logradouro  string
	Complemento string
	Bairro      string
	Localidade  string
	UF          string
	IBGE        string
	DDD         string
	SourceAPI   string
}

type Result struct {
	Address Address
	Err     error
}
