package main

import (
	"fmt"
	"net/http"
)

func main() {

	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	DB = db

	lancement := true

	if lancement == true {
		GetAllDocker()
		lancement = false
	}

	go WatchContainers()

	mux := RegisterRoutes()
	fmt.Println("Serveur démarré sur :8080")
	http.ListenAndServe(":8080", mux)
}
