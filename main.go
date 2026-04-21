package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"multithread-cep-api/entities"
	"net/http"
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
		fmt.Println("erro: CEP invalido. Informe 8 digitos, por exemplo 01001000")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := make(chan entities.Result, 2)

	go func() {
		address, err := fetchBrasilAPI(ctx, cep)
		results <- entities.Result{Address: address, Err: err}
	}()

	go func() {
		address, err := fetchViaCEP(ctx, cep)
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

	fmt.Println("erro: nenhuma API retornou um resultado valido")
	for _, err := range errs {
		fmt.Printf("- %v\n", err)
	}
	os.Exit(1)
}

func fetchBrasilAPI(ctx context.Context, cep string) (entities.Address, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return entities.Address{}, fmt.Errorf("BrasilAPI: erro ao criar requisicao: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return entities.Address{}, err
		}
		return entities.Address{}, fmt.Errorf("BrasilAPI: erro na chamada HTTP: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return entities.Address{}, fmt.Errorf("BrasilAPI: status HTTP inesperado: %d", resp.StatusCode)
	}

	var payload entities.BrasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return entities.Address{}, fmt.Errorf("BrasilAPI: erro ao decodificar resposta: %w", err)
	}

	return entities.Address{
		CEP:        payload.CEP,
		Logradouro: payload.Street,
		Bairro:     payload.Neighborhood,
		Localidade: payload.City,
		UF:         payload.State,
		SourceAPI:  "BrasilAPI",
	}, nil
}

func fetchViaCEP(ctx context.Context, cep string) (entities.Address, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return entities.Address{}, fmt.Errorf("ViaCEP: erro ao criar requisicao: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return entities.Address{}, err
		}
		return entities.Address{}, fmt.Errorf("ViaCEP: erro na chamada HTTP: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return entities.Address{}, fmt.Errorf("ViaCEP: status HTTP inesperado: %d", resp.StatusCode)
	}

	var payload entities.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return entities.Address{}, fmt.Errorf("ViaCEP: erro ao decodificar resposta: %w", err)
	}
	if payload.Erro {
		return entities.Address{}, fmt.Errorf("ViaCEP: CEP nao encontrado")
	}

	return entities.Address{
		CEP:         payload.CEP,
		Logradouro:  payload.Logradouro,
		Complemento: payload.Complemento,
		Bairro:      payload.Bairro,
		Localidade:  payload.Localidade,
		UF:          payload.UF,
		IBGE:        payload.IBGE,
		DDD:         payload.DDD,
		SourceAPI:   "ViaCEP",
	}, nil
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
