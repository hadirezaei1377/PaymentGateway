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

	// In this section, we want to send a request to Zarin Pal and get permission to use this port.
	// zarinpal is an object
	// merchandID is The ID that the portal gives us for identification operations
	// it must be 36 character
	// The sandbox is designed for the developer to check payment operations while developing the app.
	//By default, we set it to be correct so that we can test our port without exchanging any money.

	zarinpal, err := zarinpal.NewZarinpal(MERCHAND_ID, SANDBOX)
	// error in here means we have a problem to connection and We have not been able to connect to Zarin Pal for any reason.

	if err != nil {
		fmt.Fprintln(w, "خطا در پرداخت.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	intPrice, err := strconv.Atoi(price)
	if err != nil {
		// The error here is for when the user cannot enter the amount correctly, which requires error handling.
		fmt.Fprintln(w, "لطفا مبلغ را بصورت عدد وارد کنید.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// using zarinpal
	paymentUrl, authority, statusCode, err := zarinpal.NewPaymentRequest(intPrice, "http://localhost"+SERVER_PORT+"/CallBack"+price, "پرداخت تست تاپ لرن", "", "")
	// "" It means they are optional and we don't need to fill this field at the moment.
	// PaymentURL, Authority, SattusCode, err are zarinpal.NewPaymentRequest return for us and intPrice, "http://localhost"+SERCER_PORT+"/callBack"+price,"پرداخت تست توسط توسعه دهنده","",""
	// are zarinpal.NewPaymentRequest gives as inputs
	if err != nil {
		// If there is an error or not, we know that the statuscode has
		//  a value and we can use this parameter to manage the error
		if statusCode == -3 {
			// This number means that one of the limitations of the payment network has not been observed,
			// for example, a payment of less than 1000 Tomans has been requested by the user
			fmt.Fprintln(w, "مبلغ قابل پرداخت نیست.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// If the error was not the type of errors we mentioned:
		fmt.Fprintln(w, "خطایی در پرداخت رخ داده است")
		w.WriteHeader(http.StatusBadRequest)
		return
		// We should always use return in error handling so that the work is not continued.
	}

	// NewPaymentRequest recieve us the price as int but we have it as string so we must convert it
	// as well description and email and mobile and callback URL whole as string
	// It returns the authority for us.
	// This phrase is required to track payment

	// Checking a series of reports by the developer :

	//Create Record in DB
	fmt.Println("PaymentURL: ", paymentUrl, " statusCode : ", statusCode, " Authority: ", authority)
	// we can store the payment info in db
	// Sending the user to ZarinPal portal to perform payment operations

	http.Redirect(w, r, paymentUrl, 302)

	// 302 is an arbitrary number chosen by the developer to retrieve data
}

func CallBack(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "برگشت از درگاه")
	// In localhost, there are two functional parameters, one is authority, which is the payment identifier,
	// and the other is statuscode, which shows ok for successful payment and shows nook or nok for unsuccessful payment.
	// the parameters that we have to recieve
	authority := r.URL.Query().Get("Authority")
	status := r.URL.Query().Get("Status")

	// If none of these are present, it will automatically return an empty string
	if authority == "" || status == "" || status != "OK" {
		fmt.Fprintln(w, "خطایی در پرداخت رخ داده است")
		w.WriteHeader(http.StatusOK)
		// StatusOk It means that the connection to the server was made, but an error occurred
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

	// To use the price after copying the above, we need to reconnect to Zarin Pal
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
	// 101 is an arbitrary number chosen by the developer, which was created to maintain the security of the site so that the user does not perform the payment twice after being verified and creating the status of the operation code.
	fmt.Fprintln(w, "پرداخت موفقیت آمیز بود . شماره پیگیری : ", refId)

	//Create Transaction in DB
	fmt.Println(w, "Payment Verified : ", verified, " ,  refId: ", refId, " statusCode: ", statusCode)
}
