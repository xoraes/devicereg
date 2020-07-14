package main

import "sync"

type set map[*Device]bool
type AppError struct {
	Message string `json:"message, omitempty"`
	Code    int    `json:"code"`
}

const AddUnClaimedChannel = 1
const RemoveUnClaimedChannel = 2
const AddClaimedChannel = 3
const RemoveClaimedChannel = 4
const ChannelWorkerSleepSecs = 1
type op uint8
type ChannelOp struct {
	op   []op
	data *Device
}
type DeviceId uint16
type DeviceIdList []DeviceId
type NetworkDevices struct {
	DeviceTable     map[DeviceId]*Device
	UnClaimed       set
	Claimed         set
	ClaimedBuffer   DeviceIdList
	UnClaimedBuffer DeviceIdList
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
	IpAddress string  `json:"ipAddress, omitempty"`
	Geo       *LatLng `json:"geoLocation, omitempty"`
	SerialId  DeviceId  `json:"serial"`
	reserved  bool
}
