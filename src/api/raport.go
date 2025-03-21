package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	db "tracker/db"

	dbclient "github.com/PAW122/TsunamiDB/lib/dbclient"
)

type RaportBody struct {
	Id   string `json:"id"`
	Time int    `json:"time"`
}

type RaportRecord map[string]Record

type Record struct {
	TotalTime int `json:"total_time"`
}

func Raport_api(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip, _, _ = strings.Cut(ip, ":")
	}

	// get expected id
	id, err := db.GetRegisterId(ip)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Błąd odczytu treści żądania", http.StatusBadRequest)
		log.Printf("Błąd odczytu treści żądania: %v", err)
		return
	}

	var raportBody RaportBody
	err = json.Unmarshal(body, &raportBody)
	if err != nil {
		http.Error(w, "Błąd deserializacji danych", http.StatusBadRequest)
		log.Printf("Błąd deserializacji danych: %v", err)
		return
	}

	if id.Id != raportBody.Id {
		// dont give info for possible attacker
		w.Write([]byte("ok"))
		return
	}

	// try to read today total time
	data, err := dbclient.Read(id.Id, "raport_table")
	if err != nil {
		http.Error(w, "err 5_5", http.StatusBadRequest)
		return
	}

	var raportRecord RaportRecord = make(map[string]Record)
	err = json.Unmarshal(data, &raportRecord)
	if err != nil {
		http.Error(w, "err 5_6", http.StatusBadRequest)
		return
	}

	now := time.Now()
	day := now.Day()
	month := int(now.Month()) // Month() zwraca time.Month, więc konwertujemy na int
	year := now.Year()

	// add time to today total time
	today := fmt.Sprint(year, "-", month, "-", day)
	if _, ok := raportRecord[today]; ok {
		raportRecord[today] = Record{TotalTime: raportRecord[today].TotalTime + raportBody.Time}
	} else {
		raportRecord[today] = Record{TotalTime: raportBody.Time}
	}

	raportRecordJson, err := json.Marshal(raportRecord)
	if err != nil {
		http.Error(w, "err 5_7", http.StatusBadRequest)
		return
	}

	dbclient.Save(id.Id, "raport_table", raportRecordJson)
	w.Write([]byte("ok"))

}
