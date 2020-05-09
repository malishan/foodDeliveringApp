package main

import (
	"log"
	"project/foodDeliveringApp/serve"
)

const (
	version = "0.0.1"
)

func main() {
	log.Println("STARTING FoodApp SERVER, VERSION: ", version)
	log.Printf("Server is running on Port=%s, SubRoute=%s", serve.Port, serve.SubRoute)
	serve.StartRoutes()
}
