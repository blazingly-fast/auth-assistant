package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	l := log.New(os.Stdout, " Social Network ", log.LstdFlags)
	v := NewValidation()

	store, err := NewPostgresStore()
	if err != nil {
		l.Fatal(err)
	}

	if err := store.Init(); err != nil {
		l.Fatal(err)
	}
	ah := NewAccountHandler(l, v, store)

	r := mux.NewRouter()

	PostR := r.Methods(http.MethodPost).Subrouter()
	PostR.HandleFunc("/register", makeHTTPHandleFunc(ah.handleCreateAccount)).Methods(http.MethodPost)
	PostR.HandleFunc("/login", makeHTTPHandleFunc(ah.handleLogin)).Methods(http.MethodPost)

	getR := r.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/account/{id:[0-9]+}", makeHTTPHandleFunc(ah.handleGetAccountByID))
	getR.HandleFunc("/account", makeHTTPHandleFunc(ah.handleGetAccounts))
	getR.Use(ah.IsAdmin)
	getR.Use(ah.Authenticate)

	DeleteR := r.Methods(http.MethodDelete).Subrouter()
	DeleteR.HandleFunc("/account/{id:[0-9]+}", makeHTTPHandleFunc(ah.handleDeleteAccount))
	DeleteR.Use(ah.IsAdmin)
	DeleteR.Use(ah.Authenticate)

	// create a new server
	s := http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      r,                 // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: err.Error()})
		}
	}
}
