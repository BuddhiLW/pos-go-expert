package http_funcs

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

func PostHTTP() {
	c := http.Client{}
	byteJson := []byte(`{"name":"John"}`)
	buf := bytes.NewBuffer(byteJson)

	resp, err := c.Post("http://google.com", "application/json", buf)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.CopyBuffer(os.Stdout, resp.Body, nil)
}
