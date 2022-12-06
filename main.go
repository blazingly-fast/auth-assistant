package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/blazingly-fast/social-network/data"
	"github.com/blazingly-fast/social-network/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	l := log.New(os.Stdout, " Social Network ", log.LstdFlags)

	v := data.NewValidation()

	err := godotenv.Load()
	if err != nil {
		l.Fatal("Error loading .env file")
	}

	store, err := data.NewPostgresStore()
	if err != nil {
		l.Fatal(err)
	}

	if err := store.Init(); err != nil {
		l.Fatal(err)
	}

	h := handlers.NewServer(l, v, store)

	r := mux.NewRouter()

	postR := r.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/register", h.MakeHTTPHandleFunc(h.HandleCreateAccount))
	postR.HandleFunc("/login", h.MakeHTTPHandleFunc(h.HandleLogin))

	imageR := r.PathPrefix("/image/").Methods(http.MethodPost).Subrouter()
	imageR.HandleFunc("/avatar/{uuid}", h.MakeHTTPHandleFunc(h.HandleAvatar))
	imageR.Use(h.Authenticate)

	getR := r.Methods(http.MethodGet).Subrouter()
	// getR.HandleFunc("/account/{uuid}", h.MakeHTTPHandleFunc(h.HandleGetAccountByID))
	getR.HandleFunc("/account", h.MakeHTTPHandleFunc(h.HandleGetAccounts))
	getR.Use(h.Authenticate)
	getR.Use(h.Paginate)

	deleteR := r.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/account/{uuid}", h.MakeHTTPHandleFunc(h.HandleDeleteAccount))
	deleteR.Use(h.Authenticate)

	putR := r.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/account/{uuid}", h.MakeHTTPHandleFunc(h.HandleUpdateAccount))
	putR.Use(h.Authenticate)

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
