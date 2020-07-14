package main

import "sync"

type set map[*Device]bool
type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

const AddUnClaimedChannel = 1
const RemoveUnClaimedChannel = 2
const AddClaimedChannel = 3
const RemoveClaimedChannel = 4
const ChannelWorkerSleepSecs = 1

type ChannelOp struct {
	op   []uint8
	data *Device
}

type NetworkDevices struct {
	DeviceTable     map[string]*Device
	UnClaimed       set
	Claimed         set
	ClaimedBuffer   []string
	UnClaimedBuffer []string
	Channel         chan *ChannelOp
	mClaimed        sync.RWMutex
	mUnclaimed      sync.RWMutex
}
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type Device struct {
	sync.Mutex
	IpAddress  string  `json:"ipAddress"`
	Geo        *LatLng `json:"geoLocation"`
	SerialId   string  `json:"serial"`
	customerId string
	reserved   bool
}
