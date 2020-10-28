package watch

import (
	"log"
	"strconv"
	"time"
)

// SetInterval sets the interval for the watch daemon
func (srv *Watch) SetInterval(newInterval int) error {
	_, err := srv.data.SettingsPut("watch", "interval", strconv.Itoa(newInterval))
	return err
}

// GetInterval retrieves the interval for the watch daemon
func (srv *Watch) GetInterval() int {
	setting, err := srv.data.SettingsGet("watch", "interval")
	if err != nil {
		// use the default
		return tickInterval
	}
	interval, err := strconv.Atoi(setting.Data)
	if err != nil {
		// use the default
		return tickInterval
	}
	if interval < minInterval {
		// use the minimum interval
		return minInterval
	}
	return interval
}

func (srv *Watch) setLastRun() error {
	_, err := srv.data.SettingsPut("watch", "last-run", time.Now().String())
	return err
}

// LastRun retrieves the time of the last watch daemon run
func (srv *Watch) LastRun() string {
	setting, err := srv.data.SettingsGet("watch", "last-run")
	if err != nil {
		log.Printf("Error getting last run: %v", err)
		return ""
	}
	return setting.Data
}
