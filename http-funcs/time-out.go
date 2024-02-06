package http_funcs

import (
	"io/ioutil"
	"net/http"
	"time"
)

func TimeOutHTTP() {
	c := http.Client{Timeout: time.Duration(1) * time.Microsecond}
	resp, err := c.Get("http://google.com")
	if err != nil {
		println("Couldn't get the page, in time: \n", err.Error(), "\n")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	println(string(body))
}
