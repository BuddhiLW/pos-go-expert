package http_funcs

import (
	"encoding/json"
	"net/http"

	packages "github.com/buddhilw/pos-go-expert/important-packages"
)

func CEP() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cep", SearchCEP)
	http.ListenAndServe(":8989", mux)
}

func SearchCEP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cep" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cep query param is required"))
		return
	} else {
		data, error := packages.CEPSearch(cepParam)
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(error.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}

}
