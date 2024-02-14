package server

import (
	"context"

	"github.com/buddhilw/pos-go-expert/desafios/server-client/client"
)

func MigrateDB() {
	conn := connectDB()
	conn.AutoMigrate(&client.Bid{})
}

func insertDolarBid(ctx context.Context, value float64, resultsCh chan<- error) {
	conn := connectDB()
	bid := client.Bid{Value: value}
	if dbc := conn.Create(&bid); dbc.Error != nil {
		resultsCh <- dbc.Error
	} else {
		resultsCh <- nil
	}
}
