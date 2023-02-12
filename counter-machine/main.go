package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/karuppiah7890/deebees/counter-machine/pkg/config"
	"github.com/karuppiah7890/deebees/counter-machine/pkg/db"
)

var version = "dev"

func defaultHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func statusHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func isInvalidMethod(method string) bool {
	// POST is usually used for create, in our case, the only
	// operation we do is update counter.The counter is already initialized to zero, always.
	// So, we won't use POST.
	// No DELETE for now.
	// We won't use PUT also, as that's for overwriting a value.
	// We won't use any other HTTP request methods too, like OPTIONS etc
	if method != http.MethodGet && method != http.MethodPatch {
		return true
	}

	return false
}

func isInvalidRequest(req *http.Request) bool {
	if isInvalidMethod(req.Method) {
		return true
	}

	// TODO: Should content type be sent for get request too
	if req.Header.Get("Content-Type") != "application/json" {
		return true
	}

	// TODO: Should we check accept header too and give error if it's not json? As we can't support
	// anything other than json for now

	return false
}

type JsonError struct {
	Error string `json:"error"`
}

type JsonRequest struct {
	IncrementBy int `json:"incrementBy"`
}

type JsonResponse struct {
	Counter int `json:"counter"`
}

func createJsonError(err string) ([]byte, error) {
	return json.Marshal(JsonError{Error: err})
}

func getCounter(res http.ResponseWriter, req *http.Request, dbChannel chan int) {
	counter := db.GetCounter()

	jsonResponseData := JsonResponse{
		Counter: counter,
	}

	jsonResponse, err := json.Marshal(jsonResponseData)
	if err != nil {
		// TODO: Handle this error
		jsonError, _ := createJsonError(fmt.Sprintf("error occurred while forming JSON response: %v", err))
		// TODO: Handle this error
		_, _ = res.Write(jsonError)
		// TODO: The assumption here is some server side error - or it's possible
		// that it's client side too - client closed connection in between and caused
		// some error. Giving the benefit of doubt to the user for now though
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	_, err = res.Write(jsonResponse)
	if err != nil {
		// TODO: Handle this error
		jsonError, _ := createJsonError(fmt.Sprintf("error occurred while sending JSON response: %v", err))
		// TODO: Handle this error
		_, _ = res.Write(jsonError)
		// TODO: The assumption here is some server side error - or it's possible
		// that it's client side too - client closed connection in between and caused
		// some error. Giving the benefit of doubt to the user for now though
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func incrementCounter(res http.ResponseWriter, req *http.Request, dbChannel chan int) {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		// TODO: Handle this error
		jsonError, _ := createJsonError(fmt.Sprintf("error occurred while reading JSON request: %v", err))
		// TODO: Handle this error
		_, _ = res.Write(jsonError)
		// TODO: The assumption here is some server side error - or it's possible
		// that it's client side too - client closed connection in between and caused
		// some error. Giving the benefit of doubt to the user for now though
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonRequestData := JsonRequest{}

	// Parse the JSON request
	err = json.Unmarshal(requestBody, &jsonRequestData)
	if err != nil {
		// TODO: Handle this error
		jsonError, _ := createJsonError(fmt.Sprintf("error occurred while parsing JSON request: %v", err))
		// TODO: Handle this error
		_, _ = res.Write(jsonError)
		// TODO: The assumption here is some server side error - or it's possible
		// that it's client side too - client closed connection in between and caused
		// some error. Giving the benefit of doubt to the user for now though
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbChannel <- jsonRequestData.IncrementBy

	res.WriteHeader(http.StatusOK)
}

func counterHandler(dbChannel chan int) func(res http.ResponseWriter, req *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		if isInvalidRequest(req) {
			// TODO: Give more details on why it's a bad request with some error in the response body
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodPatch {
			incrementCounter(res, req, dbChannel)
		}

		if req.Method == http.MethodGet {
			getCounter(res, req, dbChannel)
		}
	}

}

func listenForShutdownSignal(sigs chan os.Signal, done chan bool) {
	<-sigs
	fmt.Printf("Received SIGTERM signal. Shutting down ...")
	done <- true
}

func main() {
	// TODO: Show command line help with description on config (env vars) etc
	// and the version etc

	// TODO: Show version for `-v` (or use it for verbose?), `-V`, `--version` maybe.
	// What about `version`? Gotta think on that

	log.Printf("version: %v", version)
	c, err := config.NewConfigFromEnvVars()
	if err != nil {
		log.Fatalf("error occurred while getting configuration from environment variables: %v", err)
	}

	dbChannel := make(chan int, 100)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	done := make(chan bool, 1)

	db.Start(dbChannel, done)

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/counter", counterHandler(dbChannel))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", c.GetPort()),
		Handler: http.DefaultServeMux,
	}

	go listenForShutdownSignal(sigs, done)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("error occurred while running server: %v", err)
	}
}
