package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	reqBody := struct {
		Email    string `json:"customer_email"`
		Password string `json:"customer_pswd"`
	}{"alishan@123", "12345"}

	bdata, err := json.Marshal(&reqBody)
	if err != nil {
		log.Fatalln(err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8989/foodApp/login", bytes.NewReader(bdata))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("roleID", "customer")

	client := http.Client{Timeout: time.Second * 40}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Header: roleID =", resp.Header.Get("roleID"))
	fmt.Println("Header: userID =", resp.Header.Get("userID"))
}
