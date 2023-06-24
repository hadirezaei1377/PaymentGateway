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
// review completly code code step by step
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
		http.Error(w, "لطفا مبلغ را وارد کنید.", http.StatusBadRequest)
		return
	}

	zarinpal, err := zarinpal.NewZarinpal(MERCHAND_ID, SANDBOX)

	if err != nil {
		http.Error(w, "خطا در پرداخت: خطای اتصال به Zarinpal", http.StatusInternalServerError)
		return
	}

	intPrice, err := strconv.Atoi(price)
	if err != nil {
		http.Error(w, "لطفا مبلغ را بصورت عدد وارد کنید.", http.StatusBadRequest)
		return
	}

	paymentUrl, authority, statusCode, err := zarinpal.NewPaymentRequest(intPrice, "http://localhost"+SERVER_PORT+"/CallBack"+price, "پرداخت تست  توسط توسعه دهنده", "", "")

	if err != nil {
		if statusCode == -3 {
			http.Error(w, "خطا در پرداخت: مبلغ قابل پرداخت نیست", http.StatusBadRequest)
			return
		}
		http.Error(w, "خطا در پرداخت: خطای سیستمی", http.StatusInternalServerError)
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
		http.Error(w, "خطایی در پرداخت رخ داده است", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	price, ok := vars["price"]
	if !ok {
		http.Error(w, "لطفا مبلغ را وارد کنید.", http.StatusBadRequest)
		return
	}
	intPrice, err := strconv.Atoi(price)
	if err != nil {
		http.Error(w, "لطفا مبلغ را بصورت عدد وارد کنید.", http.StatusBadRequest)
		return
	}

	zarinpal, err := zarinpal.NewZarinpal(MERCHAND_ID, SANDBOX)
	if err != nil {
		http.Error(w, "خطا در ساخت شی Zarinpal", http.StatusInternalServerError)
		return
	}

	verified, refId, statusCode, err := zarinpal.PaymentVerification(intPrice, authority)
	if err != nil {
		switch statusCode {
		case 101:
			http.Error(w, "این پرداخت موفق بوده و قبلا این عملیات انجام شده است.", http.StatusOK)
			return
		case -21:
			http.Error(w, "پرداخت توسط کاربر لغو شده است.", http.StatusNoContent)
			return
		case -22:
			http.Error(w, "پرداخت انجام نشده است.", http.StatusNoContent)
			return
		default:
			http.Error(w, "خطا در پرداخت", http.StatusInternalServerError)
			return
		}
	}

	if !verified {
		http.Error(w, "پرداخت تایید نشده است", http.StatusOK)
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
