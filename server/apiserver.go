package main

import (
	"fmt"
	//"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/onrik/ethrpc"
)

func main() {
	go hub.start()
	router := mux.NewRouter().StrictSlash(true)
	//EJEMPLOS RUTAS
	router.HandleFunc("/todos", TodoIndex)
	router.HandleFunc("/", serveHome2)
	router.HandleFunc("/css/bootstrap.min.css", serveHome2)
	router.HandleFunc("/todos/{todoId}", TodoShow)
	router.HandleFunc("/ws", wsIndex)
	//API PARA PEDIR INFO A USUARIO
	router.HandleFunc("/balance/{account}", Balance)

	//DASHBOARD
	router.HandleFunc("/", serveHome2)
	router.HandleFunc("/css/bootstrap.min.css", serveHome2)
	router.HandleFunc("/css/font-awesome.min.css", serveHome2)
	router.HandleFunc("/d3.v3.min.js", serveHome2)
	router.HandleFunc("/static/css/main.fa8e7bf9.css", serveHome2)
	router.HandleFunc("/static/js/main.1d47eba2.js", serveHome2)
	router.HandleFunc("/fonts/{font}", serveFont)
	log.Fatal(http.ListenAndServe(":80", router))
}

func serveFont(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	vars := mux.Vars(r)
	fuente := vars["font"]
	urlfile := "../client/build/fonts/" + fuente
	fmt.Println("retorna fichero", urlfile)

	http.ServeFile(w, r, urlfile)
}
func serveHome2(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	urlfile := "../client/build" + r.URL.Path
	fmt.Println("retorna fichero", urlfile)

	http.ServeFile(w, r, urlfile)
}
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RESPONDIENDO index")
	fmt.Fprintln(w, "Todo Index!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

//METODOS DE LA API QUE INTERACTUAN CON ETH
func Balance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tratando de obtener balance")
	vars := mux.Vars(r)
	account := vars["account"]
	rpc := ethrpc.New(*selfserver)
	balance, err3 := rpc.EthGetBalance(account, "latest")
	if err3 != nil {
		fmt.Fprintln(w, "ERROR BUSCANDO:", err3)
	} else {
		fmt.Fprintln(w, "Todo show:", balance)

	}

}
