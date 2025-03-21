package api

import (
	"fmt"
	"net/http"
	"strings"

	db "tracker/db"
)

func Register_api(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	ip := r.RemoteAddr

	// Sprawdź, czy adres zawiera port i usuń go
	if strings.Contains(ip, ":") {
		ip, _, _ = strings.Cut(ip, ":")
	}

	// Zapisz adres IP w bazie danych
	id, err := db.GetRegisterId(ip)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(id.Id))
}
