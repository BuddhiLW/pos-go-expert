package webserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
	fmt.Printf("Added handler for path: %s\n", path)
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	// Add a simple test route
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Order System API is running"}`))
	})

	fmt.Printf("Registering %d handlers\n", len(s.Handlers))
	for path, handler := range s.Handlers {
		// Register both GET and POST methods for flexibility
		s.Router.Get(path, handler)
		s.Router.Post(path, handler)
		fmt.Printf("Registered GET and POST for path: %s\n", path)
	}
	fmt.Printf("Starting web server on :%s\n", s.WebServerPort)
	http.ListenAndServe(":"+s.WebServerPort, s.Router)
}
