package main

import (
	"DIA/internal/api"
	"log"
)


func main(){
	log.Println("Application start")
	api.StartServer()
	log.Println("Application terminated!")
}

