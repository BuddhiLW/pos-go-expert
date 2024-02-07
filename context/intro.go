package context

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func Context() {
	ctx := context.Background()
	randSecs := rand.Intn(5)
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(randSecs))
	defer cancel()
	bookHotel(ctx)
}

func bookHotel(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("Couldn't book hotel, due to timeout")
	case <-time.After(time.Second * 2):
		fmt.Println("Hotel booked")
	}
}
