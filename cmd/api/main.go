package main

import (
	"log"
	"net/http"

	"github.com/clodomar/prova/configs"
	"github.com/clodomar/prova/routes"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	configs.ConnectDB()

	routes.ClienteRoute(router)

	log.Fatal(http.ListenAndServe(":6000", router))
}
