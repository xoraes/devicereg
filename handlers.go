package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func JSONError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.Code)
	b, _ := json.Marshal(err)
	w.Write(b)
}
func (devices *NetworkDevices) reserveHandler(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	p := ps.ByName("name")
	err := devices.Reserve(p)
	if err != nil {
		JSONError(w, err)
		return
	}
}

func (devices *NetworkDevices) releaseHandler(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	p := ps.ByName("name")
	err := devices.Release(p)
	if err != nil {
		JSONError(w, err)
		return
	}
}

func (devices *NetworkDevices) devicesHandler(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	p := ps.ByName("name")
	var b []byte
	switch p {
	case "unclaimed":
		{
			d := devices.GetDevices(false)
			if d == nil {
			    d = DeviceIdList{}
            }
			b, _ = json.Marshal(d)
		}
	case "claimed":
		{
			d := devices.GetDevices(true)
            if d == nil {
                d = DeviceIdList{}
            }
			b, _ = json.Marshal(d)
		}
	default:
		d, err := devices.GetDevice(p)
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
