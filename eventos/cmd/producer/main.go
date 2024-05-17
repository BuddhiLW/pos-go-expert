package main

import "github.com/buddhilw/pos-go-expert/eventos/pkg/rabbitmq"

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	rabbitmq.Publish(ch, []byte("Hello World"))
}
