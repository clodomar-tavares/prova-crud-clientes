package routes

import (
	"github.com/clodomar/prova/controllers"
	"github.com/gorilla/mux"
)

func ClienteRoute(router *mux.Router) {

	router.HandleFunc("/clientes", controllers.CreateCliente()).Methods("POST")
	router.HandleFunc("/cliente/{clienteId}", controllers.GetCliente()).Methods("GET")
	router.HandleFunc("/cliente/{clienteId}", controllers.EditACliente()).Methods("PUT")
	router.HandleFunc("/cliente/clienteId}", controllers.DeleteAUser()).Methods("DELETE")
	router.HandleFunc("/cliente", controllers.GetAllCliente()).Methods("GET") //add this
	//All routes related to users comes here
}
