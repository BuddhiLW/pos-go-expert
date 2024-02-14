package servcli

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Bid struct {
	Value float64 `json:"bid"`
}

func DolarBidRequest() {
	// Request the Dolar Bid for the local server, in port 8080/cotacao
	// The server will return the Dolar Bid in JSON format
	// - Context should timeout after 300ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	defer req.Body.Close()

	if err != nil {
		panic(err)
	}

	bid := Bid{}
	err = json.NewDecoder(req.Body).Decode(&bid)
	if err != nil {
		panic(err)
	}

	// Save current Dolar Bid in a file called "cotacao.txt"
	// - If the file already exists, it will be appended

	// Open file
	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Write to file
	// 1. convert bid.Value (float64) to string
	// 2. write "Dólar: " + bid.Value + "\n" to file
	// 3. if an error occurs, panic
	_, err = file.WriteString("Dólar: " + strconv.FormatFloat(bid.Value, 'E', -1, 64) + "\n")
	if err != nil {
		panic(err)
	}

}
