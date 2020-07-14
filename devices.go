package main

import (
	"fmt"
	"net/http"
)

func (devices *NetworkDevices) Reserve(devId, cusId string) *AppError {
	d, ok := devices.DeviceTable[devId]
	if !ok {
		return appErr(fmt.Sprintf("device id %v not found", devId), http.StatusNotFound)
	}
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	if d.reserved && d.customerId != cusId {
		return appErr("device unavailable for reservation", http.StatusConflict)
	}
	d.reserved = true
	d.customerId = cusId
	devices.Channel <- &ChannelOp{op: []uint8{AddClaimedChannel, RemoveUnClaimedChannel}, data: d}
	return nil
}

func (devices *NetworkDevices) Release(devId, cusId string) *AppError {
	d, ok := devices.DeviceTable[devId]
	if !ok {
		return appErr(fmt.Sprintf("customer id %v not found", cusId), http.StatusNotFound)
	}
	if d.reserved == false {
		return nil
	}
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	if d.customerId != cusId {
		return appErr(fmt.Sprintf("unauthorized: customer does not own the device"), http.StatusUnauthorized)
	}
	d.reserved = false
	d.customerId = ""
	devices.Channel <- &ChannelOp{op: []uint8{RemoveClaimedChannel, AddUnClaimedChannel}, data: d}
	return nil
}

func (devices *NetworkDevices) GetDevices(claimed bool) []string {
	if !claimed {
		return devices.UnClaimedBuffer
	}
	return devices.ClaimedBuffer
}

func (devices *NetworkDevices) GetDevice(devId, cusId string) (*Device, *AppError) {
	d := devices.DeviceTable[devId]

	if d == nil {
		return nil, appErr(fmt.Sprintf("device id %v not found", devId), http.StatusNotFound)
	}

	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	if d.customerId != "" && d.customerId != cusId {
		return nil, appErr("customer does not own the device", http.StatusUnauthorized)
	}
	return d, nil
}
