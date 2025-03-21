package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	api "tracker/api"

	dbclient "github.com/PAW122/TsunamiDB/lib/dbclient"
)

func StartServer(port string) {

	dbclient.InitNetworkManager(8765, nil)

	mux := http.NewServeMux()

	staticPath := filepath.Join("src", "web", "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	mux.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tracking..."))
	})

	mux.HandleFunc("/register", api.Register_api)
	mux.HandleFunc("/raport", api.Raport_api)

	fmt.Println("Server running on port", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("Error runing server: ", err)
	}
}
