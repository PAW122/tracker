package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	api "tracker/api"

	dbclient "github.com/PAW122/TsunamiDB/lib/dbclient"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func StartServer(port string) {

	dbclient.InitNetworkManager(8765, nil)

	mux := http.NewServeMux()

	staticPath := filepath.Join("src", "web", "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	mux.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tracking..."))
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w) // Ustaw nagłówki CORS przed przekazaniem obsługi do `api.Register_api`
		api.Register_api(w, r)
	})
	mux.HandleFunc("/raport", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		api.Raport_api(w, r)
	})

	fmt.Println("Server running on port", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("Error runing server: ", err)
	}
}
