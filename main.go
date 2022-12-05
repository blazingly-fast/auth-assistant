package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	l := log.New(os.Stdout, " Social Network ", log.LstdFlags)

	v := NewValidation()

	err := godotenv.Load()
	if err != nil {
		l.Fatal("Error loading .env file")
	}

	store, err := NewPostgresStore()
	if err != nil {
		l.Fatal(err)
	}

	if err := store.Init(); err != nil {
		l.Fatal(err)
	}

	ah := NewAccountHandler(l, v, store)

	r := mux.NewRouter()

	postR := r.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/register", makeHTTPHandleFunc(ah.handleCreateAccount)).Methods(http.MethodPost)
	postR.HandleFunc("/login", makeHTTPHandleFunc(ah.handleLogin)).Methods(http.MethodPost)

	imageR := r.PathPrefix("/image/").Methods(http.MethodPost).Subrouter()
	imageR.HandleFunc("/avatar/{uuid}", makeHTTPHandleFunc(ah.handleAvatar))
	imageR.Use(ah.Authenticate)

	getR := r.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/account/{uuid}", makeHTTPHandleFunc(ah.handleGetAccountByID))
	getR.HandleFunc("/account", makeHTTPHandleFunc(ah.handleGetAccounts))
	getR.Use(ah.Authenticate)

	deleteR := r.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/account/{uuid}", makeHTTPHandleFunc(ah.handleDeleteAccount))
	deleteR.Use(ah.Authenticate)

	putR := r.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/account/{uuid}", makeHTTPHandleFunc(ah.handleUpdateAccount))
	putR.Use(ah.Authenticate)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	r.Handle("/docs", sh).Methods(http.MethodGet)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./"))).Methods(http.MethodGet)

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
			log.Println(err)
			WriteJSON(w, http.StatusInternalServerError, &GenericError{Message: "Internal Server Error!"})
		}
	}
}
