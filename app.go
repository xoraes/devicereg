package main

import (
    "expvar"
    _ "expvar"
    "flag"
    "fmt"
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
	devices.DeviceTable = make(map[string]*Device)
	devices.UnClaimed = make(map[*Device]bool)
	devices.Claimed = make(map[*Device]bool)
	devices.Channel = make(chan *ChannelOp, 1000)

	go func() {
		for i := 0; i < 20; i++ {
			devId := fmt.Sprintf("dev-%v", i+1)
			d := &Device{SerialId: devId, Geo: &LatLng{Lat: 123.00001, Lng: 123.000001}, IpAddress: "10.1.10.1"}
			devices.DeviceTable[devId] = d
			devices.Channel <- &ChannelOp{data: d, op: []uint8{AddUnClaimedChannel}}
		}
	}()
	go devices.channelWorker()
	go devices.stateWorker()
	logrus.SetLevel(logrus.InfoLevel)
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	router := httprouter.New()
    // Define route and call expvar http handler
    router.GET("/debug/vars", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        expvar.Handler().ServeHTTP(w, r)
    })

	router.GET("/ping", Ping)
    router.GET("/devices/:name", basicAuth(devices.DevicesHandler))
	router.POST("/devices/:name/release", basicAuth(devices.ReleaseHandler))
	router.POST("/devices/:name/reserve", basicAuth(devices.ReserveHandler))


	if *nossl {
		logrus.Info("starting http server on port 8080")
		logrus.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		logrus.Info("starting https server on port 8081")
		logrus.Fatal(http.ListenAndServeTLS(":8081", "server.crt", "server.key", router))
	}
}
