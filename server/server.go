package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/xxmdhs/bumpcard/sql"
)

func Server(port int, db *sql.DB) {
	mux := http.NewServeMux()
	mux.HandleFunc("/getforuid", getForUIDH(db))
	s := http.Server{
		Addr:              "127.0.0.1:" + strconv.Itoa(port),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       8 * time.Second,
		Handler:           mux,
	}
	log.Panicln(s.ListenAndServe())
}

func getForUIDH(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cxt := r.Context()

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("cache-control", "max-age=3600")

		uid := r.FormValue("uid")
		if uid == "" {
			d := data{
				Code: 1,
				Msg:  "uid is empty",
			}
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(d)
			return
		}
		uidi, err := strconv.Atoi(uid)
		if err != nil {
			d := data{
				Code: 2,
				Msg:  "uid is not number",
			}
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(d)
			return
		}
		cards, err := db.GetForUID(cxt, uidi)
		if err != nil {
			d := data{
				Code: 3,
				Msg:  "db error",
			}
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(d)
			return
		}
		d := data{
			Code: 0,
			Msg:  "success",
			Data: cards,
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(d)
	}
}

type data struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
