package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func login(key, pass string) (*http.Cookie, error) {
	endpoint := "https://api.ctl.io/REST/Auth/Logon/"
	payload := fmt.Sprintf("{\"APIKey\":\"%s\", \"Password\":\"%s\"}", key, pass)
	reader := bytes.NewBufferString(payload)
	resp, err := http.Post(endpoint, "application/json", reader)
	if err != nil {
		log.Panicf("Failed authenticating: %s", resp.Body)
	}
	for _, c := range resp.Cookies() {
		if c.Name == "Tier3.API.Cookie" {
			return c, nil
		}
	}
	return nil, nil
}

func main() {
	key := flag.String("key", "", "api key")
	pw := flag.String("pass", "", "password")
	method := flag.String("method", "", "http method")
	endpoint := flag.String("endpoint", "", "endpoint")
	flag.Parse()
	if *key == "" {
		log.Panic("missing flag -key")
	}
	if *pw == "" {
		log.Panic("missing flag -pass")
	}
	cookie, _ := login(*key, *pw)

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Panic("Failed reading body content")
	}
	*method = strings.ToUpper(*method)
	fmt.Fprintf(os.Stderr, "[%s] %s ... ", *method, *endpoint)
	req, _ := http.NewRequest(*method, *endpoint, bytes.NewReader(buf))
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, resp.StatusCode)
		log.Panic("API Call Failed")
	}
	defer resp.Body.Close()
	fmt.Fprintln(os.Stderr, resp.StatusCode)
	io.Copy(os.Stdout, resp.Body)
	fmt.Println("")
}