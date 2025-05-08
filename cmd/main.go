package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"1337b04rd/internal/adapters/db"
	"1337b04rd/internal/app"
)

var (
	port     int
	showHelp bool
)

func initFlags() {
	flag.IntVar(&port, "port", 8080, "Port number to listen on")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.Parse()
}

func main() {
	initFlags()

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	if logger == nil {
		panic("logger is nil")
	}

	database := db.ConnectToDB()
	defer database.Close()

	mux := app.NewApp(database, logger)

	log.Printf("App running on http://localhost:%d/", port)
	log.Println("Minio running on http://localhost:9001/")

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func printHelp() {
	fmt.Println("hacker board")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  1337b04rd [--port <N>]")
	fmt.Println("  1337b04rd --help")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help       Show this screen.")
	fmt.Println("  --port N     Port number.")
}
