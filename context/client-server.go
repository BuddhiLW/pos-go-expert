package context

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func WorkingHttpContext() {
	http.HandleFunc("/work", workingWindow)
	http.ListenAndServe(":8989", nil)
	fmt.Println("Server is running at localhost:8989")
}

func workingWindow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request initialized")
	defer log.Println("Request finished")

	select {
	case <-time.After(2 * time.Second):
		// Log to Server (stdout)
		log.Println("Successeful request")

		// Log to Client
		w.Write([]byte("Successful request"))

	case <-ctx.Done():
		// Log to Server (stdout)
		log.Println("Request canceled")

		// Log to Client
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	}
}
