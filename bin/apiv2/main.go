package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// AuthResponse received from login call
type AuthResponse struct {
	UserName      string `json:"userName"`
	AccountAlias  string `json:"accountAlias"`
	BearerToken   string `json:"bearerToken"`
	LocationAlias string `json:"locationAlias"`
}

func login(user, pass string) (string, error) {
	endpoint := "https://api.ctl.io/v2/authentication/login/"
	payload := fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", user, pass)
	reader := bytes.NewBufferString(payload)
	resp, err := http.Post(endpoint, "application/json", reader)
	defer resp.Body.Close()
	if err != nil {
		log.Panicf("Failed authenticating: %s", resp.Body)
	}
	inst := &AuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(inst)
	if err != nil {
		log.Panicf("Failed unmarshalling response: %s", err)
	}
	return inst.BearerToken, nil
}

func main() {
	un := flag.String("user", "", "username")
	pw := flag.String("pass", "", "password")
	method := flag.String("method", "", "http method")
	endpoint := flag.String("endpoint", "", "endpoint")
	flag.Parse()
	if *un == "" {
		log.Panic("missing flag -user")
	}
	if *pw == "" {
		log.Panic("missing flag -pass")
	}
	token, _ := login(*un, *pw)

	var buf []byte
	var err error
	if *method != "GET" {
		buf, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Panic("Failed reading body content from STDIN")
		}
	}
	*method = strings.ToUpper(*method)
	fmt.Fprintf(os.Stderr, "[%s] %s ... ", *method, *endpoint)
	req, _ := http.NewRequest(*method, *endpoint, bytes.NewReader(buf))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

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
