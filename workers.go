package main

import (
    "github.com/sirupsen/logrus"
    "time"
)
func init() {

}
func (devices *NetworkDevices) channelWorker() {
	for {
		select {
		case c := <-devices.Channel:
			for _, v := range c.op {
				if v == AddUnClaimedChannel {
					devices.mUnclaimed.Lock()
					if !devices.UnClaimed[c.data] {
						devices.UnClaimed[c.data] = true
					}
					devices.mUnclaimed.Unlock()
                    logrus.Debug("Added to Unclaimed: %s",c.data)
				} else if v == AddClaimedChannel {
					devices.mClaimed.Lock()
					if !devices.Claimed[c.data] {
						devices.Claimed[c.data] = true
					}
					devices.Claimed[c.data] = true
					devices.mClaimed.Unlock()
					logrus.Debug("Added to claimed: %s",c.data)
				} else if v == RemoveClaimedChannel {
					devices.mClaimed.Lock()
					delete(devices.Claimed, c.data)
					devices.mClaimed.Unlock()
					logrus.Debug("Removed from claimed: %s",c.data)
				} else if v == RemoveUnClaimedChannel {
					devices.mUnclaimed.Lock()
					delete(devices.UnClaimed, c.data)
					devices.mUnclaimed.Unlock()
					logrus.Debug("Removed from Unclaimed: %s",c.data)
				}
			}
		}
	}
}
func (devices *NetworkDevices) stateWorker() {
	for {
		var bufferClaimed DeviceIdList
		var bufferUnClaimed DeviceIdList
		devices.mUnclaimed.Lock()
		for k, _ := range devices.UnClaimed {
			bufferUnClaimed = append(bufferUnClaimed, k.SerialId)
		}
		devices.mUnclaimed.Unlock()
		devices.mClaimed.Lock()
		for k := range devices.Claimed {
			bufferClaimed = append(bufferClaimed, k.SerialId)
		}
		devices.mClaimed.Unlock()
		devices.UnClaimedBuffer = bufferUnClaimed
		devices.ClaimedBuffer = bufferClaimed
		time.Sleep(ChannelWorkerSleepSecs * time.Second)
	}

}
