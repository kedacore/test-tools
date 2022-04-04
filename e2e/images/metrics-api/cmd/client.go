package main

import (
	"flag"
	"fmt"
	"net/http"
)

const (
	BASE_URL = "http://localhost:8080/api"
)

func setValue(value int64) error {

	client := &http.Client{}
	var req *http.Request

	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/value/%d", BASE_URL, value), nil)

	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Println("failed to set value ", err, res.Status)
		return err
	}
	return nil
}

func main() {

	valuePtr := flag.Int64("value", 0, "value")
	flag.Parse()

	setValue(*valuePtr)

}
