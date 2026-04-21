package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"multithread-cep-api/entities"
	"net/http"
)

func FetchBrasilAPI(ctx context.Context, cep string) (entities.Address, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return entities.Address{}, fmt.Errorf("BrasilAPI: erro ao criar requisição: %w", err)
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

func FetchViaCEP(ctx context.Context, cep string) (entities.Address, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return entities.Address{}, fmt.Errorf("ViaCEP: erro ao criar requisição: %w", err)
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
