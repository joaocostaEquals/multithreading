package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Address struct {
	Cep          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

type AddressViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

func requestBrAPI(cep string, ch chan<- Address) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erro na requisição da API BrasilAPI:", err)
		return
	}
	defer resp.Body.Close()

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		log.Println("Erro ao decodificar resposta da API BrasilAPI:", err)
		return
	}

	ch <- address
}

func requestViaCEP(cep string, ch chan<- AddressViaCep) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erro na requisição da API ViaCEP:", err)
		return
	}
	defer resp.Body.Close()

	var address AddressViaCep
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		log.Println("Erro ao decodificar resposta da API ViaCEP:", err)
		return
	}

	ch <- address
}

func main() {
	if len(os.Args) != 2 {
		log.Println("É preciso informar o parâmetro cep")
		return
	}

	cep := os.Args[1]
	brAPICh := make(chan Address)
	viaCEPCh := make(chan AddressViaCep)
	start := time.Now()
	go requestBrAPI(cep, brAPICh)
	go requestViaCEP(cep, viaCEPCh)

	var brAPIResponse Address
	var viaCEPResponse AddressViaCep
	select {
	case brAPIResponse = <-brAPICh:
		log.Printf("Resultado da API BrasilAPI: (%s)\n", time.Since(start))
		log.Println(brAPIResponse)
	case viaCEPResponse = <-viaCEPCh:
		log.Printf("Resultado da API ViaCEP: (%s)\n", time.Since(start))
		log.Println(viaCEPResponse)
	case <-time.After(1 * time.Second):
		log.Println("Nenhuma resposta recebida dentro do tempo limite.")
	}
}
