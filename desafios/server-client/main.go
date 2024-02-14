package servcli

func main() {
	// Start the server
	go StartServer()
	// Start the client
	StartClient()
}
