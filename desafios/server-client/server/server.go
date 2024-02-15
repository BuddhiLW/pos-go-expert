package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/buddhilw/pos-go-expert/desafios/server-client/client"
	// _ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartServer() {
	// Start the server on port 8080
	// - The server will handle requests to "/cotacao"

	// Migrate the database
	MigrateDB()

	// Creating a tmux and running the server
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", Cotacao)
	http.ListenAndServe(":8080", mux)
}

type USDBRL struct {
	Data struct {
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
	} `json:"USDBRL"`
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

	ctxGet, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*200))
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(ctxGet)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	out, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	// log.Println(string(out))

	var usdbrl USDBRL
	err = json.Unmarshal(out, &usdbrl)

	if err != nil {
		panic(err)
	}

	// Create a channel to receive the result of the request
	// This is necessary to work with the select statement
	// Which is used to handle the context timeout
	resultsCh := make(chan *USDBRL)
	go func() { resultsCh <- &usdbrl }()

	// log if ctxGet is done
	select {

	// Case usdbrl has successefully stored a value from request
	case res := <-resultsCh:
		// Convert usdbrl.Bid to float64
		v, err := strconv.ParseFloat(res.Data.Bid, 64)
		log.Println("Dolar Bid: ", v, " - ", err)
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
		bid := client.Bid{Value: v}
		// json.NewEncoder(w).Encode(bid)
		b, err := json.Marshal(bid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)

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

	// Insert the Dolar Bid into the sql database using a goroutine
	// - The goroutine should send the result of the insert to the results channel
	go insertDolarBid(ctxInsert, bid, resultsCh)

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

func connectDB() *gorm.DB {
	// connect to database
	// user:pass@tcp(localhost:3306)/challenge1?charset=utf8mb4&parseTime=True&loc=Local
	dns := "local.db"
	conn, err := gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return conn
}
