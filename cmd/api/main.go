package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Declare a string containing the application version number. Later in the book we'll // generate this automatically at build time, but for now we'll just store the version // number as a hard-coded global constant.
const version = "1.0.0"

// Define a config struct to hold all the configuration settings for our application. // For now, the only configuration settings will be the network port that we want the // server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these
// configuration settings from command-line flags when the application starts.
type config struct {
	port int
	env  string
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers, // and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	// Declare an instance of the config struct.
	var cfg config
	// Read the value of the port and env command-line flags into the config struct. We // default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()
	// Initialize a new logger which writes messages to the standard out stream, // prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// Declare an instance of the application struct, containing the config struct and // the logger.
	app := &application{
		config: cfg,
		logger: logger}

	// // Declare a new servemux and add a /v1/healthcheck route which dispatches requests // to the healthcheckHandler method (which we will create in a moment).
	// mux := http.NewServeMux()
	// mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// Declare a HTTP server with some sensible timeout settings, which listens on the // port provided in the config struct and uses the servemux we created above as the // handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Start the HTTP
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
