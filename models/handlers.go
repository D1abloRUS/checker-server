package models

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/D1abloRUS/checker-server/config"

	"github.com/julienschmidt/httprouter"
)

//Client Hash
type Client struct {
	Hash string
}

//Status task
type Status struct {
	ID     int
	Status bool
}

//FailOnHtpp func
func FailOnHtpp(err error, w http.ResponseWriter, msg string, httpcode int) {
	if err != nil {
		http.Error(w, msg, httpcode)
		return
	}
}

//GetTask /api/v1/gettask/:id return json []
func GetTask(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, _ := strconv.Atoi(p.ByName("id"))

		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		bks, err := GetTasksReq(env.DB, id)
		FailOnHtpp(err, w, "Requst not found", 500)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err = json.NewEncoder(w).Encode(bks); err != nil {
			w.WriteHeader(500)
		}
	}
}

//Activate post /api/v1/activate json {"hash": "123"}
func Activate(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var c Client

		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&c)
		FailOnHtpp(err, w, "json error", 400)

		err = InsertHash(env.DB, c.Hash)
		FailOnHtpp(err, w, "hash key not found", 500)

		bks, err := GetCheckerID(env.DB, c.Hash)
		FailOnHtpp(err, w, "checker id not found", 500)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err = json.NewEncoder(w).Encode(bks); err != nil {
			w.WriteHeader(500)
		}
	}
}

//StatusUpdate post /api/v1/statusupdate  json [{"id": 1, "status": 0},{N}]
func StatusUpdate(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var status []Status

		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&status)
		FailOnHtpp(err, w, "json error", 400)

		for s := range status {
			err = UpdateStatus(env.DB, status[s].ID, status[s].Status)
			FailOnHtpp(err, w, "task not found", 500)
		}
		w.WriteHeader(200)
	}
}
