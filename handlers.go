package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func basicAuth(pass httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, ok := r.BasicAuth()
		if !ok || !validate(user, password) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		pass(w, r, ps)
	}
}

func validate(username, password string) bool {
	//this should call a auth service or performance from a secure db
	if username != "" && password == "test" {
		return true
	}
	return false
}

func JSONError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.Code)
	b, _ := json.Marshal(err)
	w.Write(b)
}
func (devices *NetworkDevices) ReserveHandler(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	user, _, _ := request.BasicAuth()
	p := ps.ByName("name")
	err := devices.Reserve(p, user)
	if err != nil {
		JSONError(w, err)
		return
	}
}

func (devices *NetworkDevices) ReleaseHandler(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	user, _, _ := request.BasicAuth()
	p := ps.ByName("name")
	err := devices.Release(p, user)
	if err != nil {
		JSONError(w, err)
		return
	}
}

func (devices *NetworkDevices) DevicesHandler(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	user, _, _ := request.BasicAuth()
	p := ps.ByName("name")
	var b []byte
	switch p {
	case "unclaimed":
		{
			d := devices.GetDevices(false)
			if d == nil {
				d = []string{}
			}
			b, _ = json.Marshal(d)
		}
	case "claimed":
		{
			d := devices.GetDevices(true)
			if d == nil {
				d = []string{}
			}
			b, _ = json.Marshal(d)
		}
	default:
		d, err := devices.GetDevice(p, user)
		if err != nil {
			JSONError(w, err)
			return
		}
		b, _ = json.Marshal(d)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}
func Ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Pong")
}
