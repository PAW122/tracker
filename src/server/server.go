package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	api "tracker/api"

	dbclient "github.com/PAW122/TsunamiDB/lib/dbclient"
)

// Obsługa CORS
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Możesz ustawić konkretną domenę zamiast "*"
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// Globalny handler dla `OPTIONS`
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func StartServer(port string) {

	dbclient.InitNetworkManager(8765, nil)

	mux := http.NewServeMux()

	staticPath := filepath.Join("src", "web", "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	mux.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tracking..."))
	})

	mux.HandleFunc("/register", corsMiddleware(api.Register_api))
	mux.HandleFunc("/raport", corsMiddleware(api.Raport_api))

	fmt.Println("Server running on port", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("Error running server: ", err)
	}
}
