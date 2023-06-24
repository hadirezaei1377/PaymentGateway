package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sinabakh/go-zarinpal-checkout"

	"github.com/gorilla/mux"
)

// todo :
// review code
// add feature

const (
	SERVER_PORT = ":8383"
	MERCHAND_ID = "111111222222333333444444555555666666" // sellerport
	SANDBOX     = true                                   // sandbox env
)

func main() {
	router := mux.NewRouter()
	router.Methods("GET").Path("/Bank{price}").HandlerFunc(Bank)
	router.Methods("GET").Path("/CallBack{price}").HandlerFunc(CallBack)
	log.Fatal(http.ListenAndServe(SERVER_PORT, router))
}

func Bank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // extract the "price" parameter from the URL
	price, ok := vars["price"]

	if !ok {
		fmt.Fprintln(w, "لطفا مبلغ را وارد کنید.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	zarinpal, err := zarinpal.NewZarinpal(MERCHAND_ID, SANDBOX)

	if err != nil {
		fmt.Fprintln(w, "خطا در پرداخت.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	intPrice, err := strconv.Atoi(price)
	if err != nil {

		fmt.Fprintln(w, "لطفا مبلغ را بصورت عدد وارد کنید.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	paymentUrl, authority, statusCode, err := zarinpal.NewPaymentRequest(intPrice, "http://localhost"+SERVER_PORT+"/CallBack"+price, "پرداخت تست  توسط توسعه دهنده", "", "")

	if err != nil {

		if statusCode == -3 {

			fmt.Fprintln(w, "مبلغ قابل پرداخت نیست.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, "خطایی در پرداخت رخ داده است")
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	fmt.Println("PaymentURL: ", paymentUrl, " statusCode : ", statusCode, " Authority: ", authority)

	http.Redirect(w, r, paymentUrl, http.StatusFound)

}

func CallBack(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "برگشت از درگاه")

	authority := r.URL.Query().Get("Authority")
	status := r.URL.Query().Get("Status")

	if authority == "" || status == "" || status != "OK" {
		fmt.Fprintln(w, "خطایی در پرداخت رخ داده است")
		w.WriteHeader(http.StatusOK)

		return
	}

	vars := mux.Vars(r)
	price, ok := vars["price"]
	if !ok {
		fmt.Fprintln(w, "لطفا مبلغ را وارد کنید.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	intPrice, err := strconv.Atoi(price)
	if err != nil {
		fmt.Fprintln(w, "لطفا مبلغ را بصورت عدد وارد کنید.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	zarinpal, err := zarinpal.NewZarinpal(MERCHAND_ID, SANDBOX)
	if err != nil {
		fmt.Fprintln(w, "خطا در پرداخت.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	verified, refId, statusCode, err := zarinpal.PaymentVerification(intPrice, authority)
	if err != nil {
		if statusCode == 101 {
			fmt.Fprintln(w, "این پرداخت موفق بوده و قبلا این عملیات انجام شده است.")
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Fprintln(w, "خطا در پرداخت.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "پرداخت موفقیت آمیز بود . شماره پیگیری : ", refId)

	fmt.Println(w, "Payment Verified : ", verified, " ,  refId: ", refId, " statusCode: ", statusCode)
}

/*
test :
localhost:8383/Bank5000
localhost:8383/callback5000
*/
