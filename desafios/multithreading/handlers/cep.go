package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/paemuri/brdoc"
)

type BrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func ViaCEP(cep string, c1 *chan ViaCep) {
	// get http://viacep.com.br/ws/" + cep + "/json/
	r, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer r.Body.Close()
	var viaCep ViaCep
	if err := json.NewDecoder(r.Body).Decode(&viaCep); err != nil {
		// handle error
		log.Fatal(err)
	}
	*c1 <- viaCep

}

func BrasilAPI(cep string, c2 *chan BrasilApi) {
	// get https://brasilapi.com.br/api/cep/v1/ + cep
	r, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer r.Body.Close()
	var brasilApi BrasilApi

	if err := json.NewDecoder(r.Body).Decode(&brasilApi); err != nil {
		// handle error
		log.Fatal(err)
	}
	*c2 <- brasilApi
}

func CEP(w http.ResponseWriter, r *http.Request) {
	// get cep in query
	cep := r.URL.Query().Get("cep")

	if !brdoc.IsCEP(cep) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invalid CEP"))
	}

	c1 := make(chan ViaCep)
	c2 := make(chan BrasilApi)
	go ViaCEP(cep, &c1)
	go BrasilAPI(cep, &c2)

	select {
	case viaCep := <-c1:
		fmt.Printf("ViaCEP || Estado: %v, Cidade: %v, Bairro: %v, Rua: %v  \n", viaCep.Uf, viaCep.Localidade, viaCep.Bairro, viaCep.Logradouro)
	case brasilApi := <-c2:
		fmt.Printf("brasilAPI || Estado: %v, Cidade: %v, Bairro: %v, Rua: %v  \n", brasilApi.State, brasilApi.City, brasilApi.Neighborhood, brasilApi.Street)
	case <-time.After(time.Second):
		fmt.Printf("Timeout!\n")
	}
}
