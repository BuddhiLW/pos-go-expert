package main

import (
	"runtime"
	"time"

	"github.com/buddhilw/pos-go-expert/desafios/server-client/client"
	"github.com/buddhilw/pos-go-expert/desafios/server-client/server"
)

func main() {
	// Start the server
	go server.StartServer()
	defer runtime.Goexit()
	time.Sleep(100 * time.Millisecond)
	// Start the client
	// client.DolarBidRequest()
	client.StartClient()
}
