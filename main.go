package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/sinabakh/go-zarinpal-checkout"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
)

// todo :
// review completly code code step by step
// add feature

const (
	SERVER_PORT   = ":8383"
	MERCHAND_ID   = "111111222222333333444444555555666666"                         // sellerport
	SANDBOX       = true                                                           // sandbox env
	USERNAME      = "user"                                                         // for authentication
	PASSWORD_HASH = "$2a$10$q7OQK8cYUQLkpf3I8utM9eyVCxybLzRgLWz6hQf.hrwfXgA.4rk5S" // hashed password for authentication
)

func main() {
	router := mux.NewRouter()
	router.Methods("GET").Path("/Bank{price}").HandlerFunc(authenticate(Bank))
	router.Methods("GET").Path("/CallBack{price}").HandlerFunc(authenticate(CallBack))
	log.Fatal(http.ListenAndServe(SERVER_PORT, router))
}

func Bank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // extract the "price" parameter from the URL
	price, ok := vars["price"]

	if !ok {
		http.Error(w, "لطفا مبلغ را وارد کنید.", http.StatusBadRequest)
		return
	}

	// Validate the price input
	match, _ := regexp.MatchString(`^[1-9]\d*$`, price)
	if !match {
		http.Error(w, "لطفا مبلغ را بصورت عدد وارد کنید.", http.StatusBadRequest)
		return
	}

	intPrice, err := strconv.Atoi(price)
	if err != nil {
		http.Error(w, "لطفا مبلغ را بصورت عدد وارد کنید.", http.StatusBadRequest)
		return
	}

	// Check if the price is within the allowed range
	MAX_PAYMENT_AMOUNT := 1000000000
	if intPrice > MAX_PAYMENT_AMOUNT || intPrice <= 0 {
		http.Error(w, fmt.Sprintf("\u0645\u0628\u0644\u063a \u067e\u0631\u062f\u0627\u062e\u062a \u0646\u0627\u0645\u0639\u062a\u0628\u0631 \u0627\u0633\u062a. \u0645\u06cc\u200c\u0628\u0627\u06cc\u0633\u062a \u0628\u06cc\u0646 %d \u0648 %d \u0628\u0627\u0634\u062f.", 1, MAX_PAYMENT_AMOUNT), http.StatusBadRequest)

		return
	}

	zarinpal, err := zarinpal.NewZarinpal(MERCHAND_ID, SANDBOX)

	if err != nil {
		http.Error(w, "خطا در پرداخت: خطای اتصال به Zarinpal", http.StatusInternalServerError)
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

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != USERNAME || !CheckPasswordHash(password, PASSWORD_HASH) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized access.\n")
			return
		}
		next.ServeHTTP(w, r)
	}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*
test :
localhost:8383/Bank5000
localhost:8383/callback5000
*/
