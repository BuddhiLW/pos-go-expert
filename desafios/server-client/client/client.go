package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

type Bid struct {
	Value float64 `json:"bid"`
}

func (bid Bid) BidStore() {
	// Save current Dolar Bid in a file called "cotacao.txt"
	// - If the file already exists, it will be appended

	// Open file
	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write to file
	// 1. convert bid.Value (float64) to string
	// 2. write "Dólar: " + bid.Value + "\n" to file
	// 3. if an error occurs, panic
	_, err = file.WriteString("Dólar: " + strconv.FormatFloat(bid.Value, 'E', -1, 64) + "\n")
	if err != nil {
		panic(err)
	}
}

func DolarBidRequest() {
	// Request the Dolar Bid for the local server, in port 8080/cotacao
	// The server will return the Dolar Bid in JSON format
	// - Context should timeout after 300ms
	ctxGet, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/cotacao", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctxGet)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var bid Bid
	err = json.Unmarshal(out, &bid)
	if err != nil {
		panic(err)
	}
	log.Println("Dolar Bid: ", bid.Value)

	var respCh = make(chan *Bid)
	go func() { respCh <- &bid }()

	select {
	case res := <-respCh:
		println("Writing to file...\n")
		res.BidStore()
	case <-ctxGet.Done():
		println("Couldn't get the Dolar Bid, in time: \n", ctxGet.Err(), "\n")
	}
}

func PeriodicDolarBidRequest() {
	// Periodically request the Dolar Bid from the local server
	c := cron.New()
	c.AddFunc("@every 1s30ms", DolarBidRequest)
	c.Start()
}

func StartClient() {
	PeriodicDolarBidRequest()
}
