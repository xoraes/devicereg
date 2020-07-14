package main

import (
	"expvar"
	_ "expvar"
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

func appErr(message string, code int) *AppError {
	return &AppError{Message: message, Code: code}
}

func main() {

	nossl := flag.Bool("nossl", false, "run server in nossl mode")
	debug := flag.Bool("debug", false, "run server in debug mode")

	flag.Parse()
	devices := NetworkDevices{}
	devices.DeviceTable = make(map[DeviceId]*Device)
	devices.UnClaimed = make(map[*Device]bool)
	devices.Claimed = make(map[*Device]bool)
	devices.Channel = make(chan *ChannelOp, 1000)

	go func() {
		for i := 0; i < 20; i++ {
			id := DeviceId(i + 1)
			d := &Device{SerialId: id, Geo: &LatLng{Lat: 123.00001, Lng: 123.000001}, IpAddress: "10.1.10.1"}
			devices.DeviceTable[id] = d
			devices.Channel <- &ChannelOp{data: d, op: []op{AddUnClaimedChannel}}
		}
	}()
	go devices.channelWorker()
	go devices.stateWorker()
	logrus.SetLevel(logrus.InfoLevel)
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	router := httprouter.New()
	router.GET("/devices/:name", devices.devicesHandler)
	router.GET("/ping", Ping)
	router.POST("/devices/:name/release", devices.releaseHandler)
	router.POST("/devices/:name/reserve", devices.reserveHandler)
	// Define route and call expvar http handler
	router.GET("/debug/vars", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		expvar.Handler().ServeHTTP(w, r)
	})

	if *nossl {
		logrus.Info("starting http server on port 8080")
		logrus.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		logrus.Info("starting https server on port 8081")
		logrus.Fatal(http.ListenAndServeTLS(":8081", "server.crt", "server.key", router))
	}
}
