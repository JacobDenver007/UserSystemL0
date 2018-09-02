package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bocheninc/base/log"
)

//BroadcastTx broadcast transaction
func BroadcastTx(param string) {
	paramStr := `{"id":1,"method":"Transaction.Broadcast","params":["` + param + `"]}`
	req, err := http.NewRequest("POST", "http://127.0.0.1:8881", bytes.NewBufferString(paramStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	var client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1000,
		},
		Timeout: time.Duration(60) * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		log.Errorf("broadcast %v failed --- %v", param, err)
	} else {
		defer response.Body.Close()
		body, er := ioutil.ReadAll(response.Body)
		if er != nil {
			log.Errorf("parse response body %v failed --- %v", string(body), err)
		} else {
			log.Info("broadcast %v succeed --- %v", param, string(body))
		}
	}
}
