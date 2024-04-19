package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/paemuri/brdoc"
)

func ViaCEP(c1 chan ViaCEP) {

}

func BrasilAPI(c2 chan BrasilAPI) {

}

func CEP(w http.ResponseWriter, r *http.Request) {
	// get cep in query
	cep := r.URL.Query("cep")

	if !brdoc.IsCEP(cep) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invalid CEP"))
	}

	c1 := make(chan ViaCEP)
	c2 := make(chan BrasilAPI)
	go ViaCEP(cep, &c1)
	go BrasilAPI(cep, &c2)

	select {
	case viaCep := <-c1:
		fmt.Printf("ViaCEP")
	case brasilApi := <-c2:
		fmt.Printf("brasilAPI")
	case <-time.After(time.Second):
		fmt.Printf("Timeout!")
	}

}
