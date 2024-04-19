package main

import (
	"log"
	"net/http"

	"github.com/buddhilw/pos-go-expert/desafios/multithreading/handlers"
)

// Os requisitos para este desafio são:
// - Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
// - O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
// - Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

func main() {
	//source 1: https://brasilapi.com.br/api/cep/v1/ + cep
	//source 2: http://viacep.com.br/ws/ + cep + /json/

	mux := http.NewServeMux()
	mux.HandleFunc("/cep", handlers.CEP)
	log.Println("Running at: 8001")
	http.ListenAndServe(":8001", mux)
}
