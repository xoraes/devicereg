package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (devices *NetworkDevices) Reserve(id string) *AppError {
	i, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return appErr("invalid id", http.StatusBadRequest)
	}
	d, ok := devices.DeviceTable[DeviceId(i)]
	if !ok {
		return appErr(fmt.Sprintf("device id %v not found", id), http.StatusBadRequest)
	}
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	if d.reserved {
		return appErr("device unavailable for reservation", http.StatusConflict)
	}
	d.reserved = true
	devices.Channel <- &ChannelOp{op: []op{AddClaimedChannel, RemoveUnClaimedChannel}, data: d}
	return nil
}

func (devices *NetworkDevices) Release(id string) *AppError {
	i, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return appErr("invalid id", http.StatusBadRequest)
	}
	d, ok := devices.DeviceTable[DeviceId(i)]
	if !ok {
		return appErr(fmt.Sprintf("device id %v not found", id), http.StatusBadRequest)
	}
	if d.reserved == false {
		return nil
	}
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	d.reserved = false
	devices.Channel <- &ChannelOp{op: []op{RemoveClaimedChannel, AddUnClaimedChannel}, data: d}
	return nil
}

func (devices *NetworkDevices) GetDevices(claimed bool) DeviceIdList {
	if !claimed {
		return devices.UnClaimedBuffer
	}
	return devices.ClaimedBuffer
}

func (devices *NetworkDevices) GetDevice(id string) (*Device, *AppError) {
	i, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return nil, appErr("invalid id", http.StatusBadRequest)
	}
	deviceId := DeviceId(i)
	d := devices.DeviceTable[deviceId]
	if d == nil {
		return nil, appErr(fmt.Sprintf("device id %v not found", id), http.StatusBadRequest)
	}
	return d, nil
}
