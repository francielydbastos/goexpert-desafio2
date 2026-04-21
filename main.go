package main

import (
	"context"
	"fmt"
	"multithread-cep-api/entities"
	"multithread-cep-api/services"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("uso: go run . <cep>")
		os.Exit(1)
	}

	cep := sanitizeCEP(os.Args[1])
	if !isValidCEP(cep) {
		fmt.Println("erro: CEP inválido. Informe 8 dígitos, por exemplo 01001000 ou 01001-000")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := make(chan entities.Result, 2)

	go func() {
		address, err := services.FetchBrasilAPI(ctx, cep)
		results <- entities.Result{Address: address, Err: err}
	}()

	go func() {
		address, err := services.FetchViaCEP(ctx, cep)
		results <- entities.Result{Address: address, Err: err}
	}()

	errs := make([]error, 0, 2)
	for i := 0; i < 2; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("erro: timeout de 1s atingido. Nenhuma API respondeu a tempo")
			os.Exit(1)
		case res := <-results:
			if res.Err == nil {
				cancel()
				printAddress(res.Address)
				return
			}
			errs = append(errs, res.Err)
		}
	}

	fmt.Println("erro: nenhuma API retornou um resultado válido")
	for _, err := range errs {
		fmt.Printf("- %v\n", err)
	}
	os.Exit(1)
}

func printAddress(address entities.Address) {
	fmt.Printf("API vencedora: %s\n", address.SourceAPI)
	fmt.Printf("CEP: %s\n", address.CEP)
	fmt.Printf("Logradouro: %s\n", address.Logradouro)
	if address.Complemento != "" {
		fmt.Printf("Complemento: %s\n", address.Complemento)
	}
	fmt.Printf("Bairro: %s\n", address.Bairro)
	fmt.Printf("Cidade: %s\n", address.Localidade)
	fmt.Printf("UF: %s\n", address.UF)
	if address.IBGE != "" {
		fmt.Printf("IBGE: %s\n", address.IBGE)
	}
	if address.DDD != "" {
		fmt.Printf("DDD: %s\n", address.DDD)
	}
}

func sanitizeCEP(raw string) string {
	return strings.ReplaceAll(raw, "-", "")
}

func isValidCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}
