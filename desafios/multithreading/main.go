package main

import (
	"net/http"

	"github.com/buddhilw/pos-go-expert/desafios/multithreading/handlers"
)

// Os requisitos para este desafio são:
// - Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
// - O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
// - Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

func main() {
	// https://brasilapi.com.br/api/cep/v1/01153000 + cep
	// http://viacep.com.br/ws/" + cep + "/json/
	// cep := os.Args[1]

	// Create a server
	mux := http.NewServeMux()
	mux.HandleFunc("/cep", handlers.CEP)
	http.ListenAndServe(":8001", mux)
}
