package main

import (
	"crud/controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/usuarios", controllers.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/usuarios", controllers.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", controllers.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", controllers.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/usuarios/{id}", controllers.DeleteUser).Methods(http.MethodDelete)

	fmt.Println("Escutando na porta 5010")
	log.Fatal(http.ListenAndServe(":5010", router))
}
