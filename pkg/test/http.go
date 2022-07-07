package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	client := &http.Client{}

	req, err := http.NewRequest("", "http://123.56.62.34:80/graph", strings.NewReader(""))
	if err != nil {
		fmt.Println(err)
	}

	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "data-platform-prometheus.com")
	req.Host = "data-platform-prometheus.com"
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}
