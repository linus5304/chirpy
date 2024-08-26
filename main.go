package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/linus5304/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {

	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// if _, err := os.Stat("database.json"); err == nil {
	// 	if err := os.Remove("database.json"); err != nil {
	// 		log.Fatalf("could not delete database.json: %v", err)
	// 	}
	// } else if !os.IsNotExist(err) {
	// 	log.Fatalf("error checking database.json: %v", err)
	// }
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		if err := os.Remove("database.json"); err != nil {
			log.Fatalf("could not delete database.json: %v", err)
		}
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()

	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleChirpGet)

	mux.HandleFunc("POST /api/users", apiCfg.handleUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUsersUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.handleUserLogin)

	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
