package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	zarinpal "github.com/sinabakh/go-zarinpal-checkout"

	"github.com/gorilla/mux"
)

const (
	SERVER_PORT = ":8383"
	MERCHAND_ID = "111111222222333333444444555555666666"
	SANDBOX     = true
)


func main() {
	router := mux.NewRouter()

	router.Methods("GET").Path("/Bank{price}").HandlerFunc(Bank)
	// From this path, we send the user and the corresponding price to the bank function to perform the desired operation on it.
	router.Methods("GET").Path("/CallBack{price}").HandlerFunc(CallBack)
	log.Fatal(http.ListenAndServe(SERVER_PORT, router))
}

func Bank(w http.ResponseWriter, r *http.Request) {
	// get price
	vars := mux.Vars(r)
	// recieve request as variable
	// This command gives us a map with string and string, and with the key we can access what we need.
	// getting price key by below command
	price, ok := vars["price"]
	// What it returns for us is a price and a boolean flag for presence or absence (ok).
	// checking this command by a loop :
	if !ok {
		fmt.Fprintln(w, "لطفا مبلغ را وارد کنید.")
		// writer : It is the answer that we send to the user
		// Here this response should tell the user that you sent a bad request
		// This is called setting the status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
