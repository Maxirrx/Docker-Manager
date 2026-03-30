package main

import (
	"net/http"
	"fmt"
)

func main() {

	db, err := ConnectDB()
	if err != nil {
	    panic(err)
	}
	DB = db


	GetAllDocker()

	go WatchContainers()	

	mux := RegisterRoutes()
	fmt.Println("Serveur démarré sur :8080")
	http.ListenAndServe(":8080", mux)}
