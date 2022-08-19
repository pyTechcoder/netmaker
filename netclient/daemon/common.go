package daemon

import (
	"errors"
	"runtime"
	"time"

	"github.com/gravitl/netmaker/netclient/config"
)

// InstallDaemon - Calls the correct function to install the netclient as a daemon service on the given operating system.
func InstallDaemon(cfg config.ClientConfig) error {
	os := runtime.GOOS
	var err error

	interval := "15"
	if cfg.Server.CheckinInterval != "" {
		interval = cfg.Server.CheckinInterval
	}

	switch os {
	case "windows":
		err = SetupWindowsDaemon()
	case "darwin":
		err = SetupMacDaemon(interval)
	case "linux":
		err = SetupSystemDDaemon(interval)
	case "freebsd":
		err = SetupFreebsdDaemon()
	default:
		err = errors.New("this os is not yet supported for daemon mode. Run join cmd with flag '--daemon off'")
	}
	return err
}

// Restart - restarts a system daemon
func Restart() error {
	os := runtime.GOOS
	var err error

	time.Sleep(time.Second)

	switch os {
	case "windows":
		RestartWindowsDaemon()
	case "darwin":
		RestartLaunchD()
	case "linux":
		RestartSystemD()
	case "freebsd":
		FreebsdDaemon("restart")
	default:
		err = errors.New("this os is not yet supported for daemon mode. Run join cmd with flag '--daemon off'")
	}
	return err
}

// Stop - stops a system daemon
func Stop() error {
	os := runtime.GOOS
	var err error

	time.Sleep(time.Second)

	switch os {
	case "windows":
		StopWindowsDaemon()
	case "darwin":
		StopLaunchD()
	case "linux":
		StopSystemD()
	case "freebsd":
		FreebsdDaemon("stop")
	default:
		err = errors.New("no OS daemon to stop")
	}
	return err
}
