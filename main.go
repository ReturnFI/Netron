package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"netron/cmdtools"
	"netron/handlers"

	"github.com/gorilla/mux"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	run := flag.Bool("run", false, "Run the server")
	port := flag.String("port", "8080", "Port to run server on")
	removeDeps := flag.Bool("remove-deps", false, "Remove installed dependencies")

	flag.StringVar(port, "p", "8080", "Port to run server on (shorthand)")
	flag.Parse()

	if *removeDeps {
		cmdtools.RemoveDependency()
		os.Exit(0)
	}

	if !*run {
		fmt.Println("Usage:")
		fmt.Println("  --run              : Run the server (prompts for dependency install if needed)")
		fmt.Println("  --remove-deps      : Remove dependencies")
		fmt.Println("  --port or -p [port]: Specify port (default: 8080)")
		os.Exit(1)
	}

	if !cmdtools.EnsureDependency() {
		os.Exit(1)
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/system", handlers.GetSystemInfo).Methods("GET")
	r.HandleFunc("/api/speedtest", handlers.GetSpeedTest).Methods("GET")
	r.HandleFunc("/api/speedtest/start", handlers.StartSpeedTest).Methods("POST")

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.FS(staticFS)))

	fmt.Printf("Server starting on :%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, r))
}