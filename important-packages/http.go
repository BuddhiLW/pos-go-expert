package packages

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Http() {
	req, err := http.Get("https://www.google.com.br")
	if err != nil {
		panic(err)

	}
	defer req.Body.Close()

	// res, err := io.ReadAll(req.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Data from google (in binary form)", res)
}

type ViaCEP struct {
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

func CEPSearch(cep string) (*ViaCEP, error) {
	if cep == "" {
		// fmt.Println("Fetching data from ViaCEP... (default cep: 71540043)\n")
		// cep = "71540043"
		return &ViaCEP{}, errors.New("cep query param is required")
	}

	req, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var data ViaCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("CEP: %s, UF: %s, Localidade: %s, Bairro: %s, Lougradouro: %s", data.Cep, data.Uf, data.Localidade, data.Bairro, data.Logradouro)
	return &data, nil
}

// func (v ViaCEP) Json() []byte {
// 	res, err := json.Marshal(v)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return res
// }
