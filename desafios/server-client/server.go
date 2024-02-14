package servcli

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func StartServer() {
	// Start the server on port 8080
	// - The server will handle requests to "/cotacao"

	// Creating a tmux and running the server
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", Cotacao)
	http.ListenAndServe(":8080", mux)
}

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

// - The server will return the Dolar Bid in JSON format (consume: https://economia.awesomeapi.com.br/json/last/USD-BRL)
// - The server will use sql to store the Dolar Bid
func Cotacao(w http.ResponseWriter, r *http.Request) {
	// Request the Dolar Bid from the API
	// - Extra: If an error occurs, the server should return a 500 status code
	// - If the request is successful, the server should return the Dolar Bid in JSON format
	// - The server should also store the Dolar Bid in a sql database
	// - Use context to timeout the request after 200ms
	// - Use context to timeout the sql insert after 10ms

	ctxGet := r.Context() // log.Println("Request initialized")
	// defer log.Println("Request finished")

	req, err := http.NewRequestWithContext(ctxGet, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer req.Body.Close()

	// var usdbrl USDBRL
	// err = json.NewDecoder(req.Body).Decode(&usdbrl)
	resultsCh := make(chan *USDBRL)
	err = json.NewDecoder(req.Body).Decode(&resultsCh)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// log if ctxGet is done
	select {

	// Case usdbrl has successefully stored a value from request
	case usdbrl := <-resultsCh:
		// Convert usdbrl.Bid to float64
		v, err := strconv.ParseFloat(usdbrl.Bid, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// store the Dolar Bid in a sql database
		err = StoreDolarBid(v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			log.Println(err.Error())
			return
		}

		// Success request and convertion to float
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Return the Dolar Bid in JSON format
		bid := Bid{Value: v}
		json.NewEncoder(w).Encode(bid)

	// Case context for consuming the API runs out
	case <-ctxGet.Done():
		// Log to Server (stdout)
		log.Println("Request canceled due to API-Consumption timeout")

		// Log to Client
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	}
}

func StoreDolarBid(bid float64) error {
	// Store the Dolar Bid in a sql database
	// - Use context to timeout the sql insert after 10ms
	// - If the insert is successful, return nil
	// - If the insert fails, return an error
	// - Extra: Use a real sql database (e.g. postgres, mysql, etc)

	// Create a channel to receive the result of the insert
	resultsCh := make(chan error)

	// Create a context with a timeout of 10ms
	ctxInsert, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Insert the Dolar Bid into the sql database
	// go insertDolarBid(ctxInsert, bid, resultsCh)

	// Log if ctxInsert is done
	select {
	case err := <-resultsCh:
		if err != nil {
			return err
		} else {
			log.Println("Dolar-bid inserted in the DB successfully")
			return nil
		}
	case <-ctxInsert.Done():
		return errors.New("Dolar-bid Insert timeout")
	}
}

// func connectDB() *sql.DB {
// 	// connect to database
// 	conn, err := sql.Open("mysql", "buddhilw:pass@tcp(localhost:3306)/challenge1")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return conn
// }

// func insertDolarBid(ctx context.Context, bid float64, resultsCh chan<- error) {

// }
