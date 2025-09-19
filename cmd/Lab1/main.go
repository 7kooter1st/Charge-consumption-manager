package main

import (
	"log"
	"RIP/iternal/api"
)

func main(){
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminnated!")
}