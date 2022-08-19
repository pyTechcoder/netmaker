package daemon

import (
	//"github.com/davecgh/go-spew/spew"

	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gravitl/netmaker/netclient/ncutils"
)

const EXEC_DIR = "/sbin/"

// SetupSystemDDaemon - sets system daemon for supported machines
func SetupSystemDDaemon(interval string) error {

	if ncutils.IsWindows() {
		return nil
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	binarypath := dir + "/netclient"

	_, err = os.Stat("/etc/netclient/config")
	if os.IsNotExist(err) {
		os.MkdirAll("/etc/netclient/config", 0744)
	} else if err != nil {
		log.Println("couldnt find or create /etc/netclient")
		return err
	}
	//install binary
	//should check if the existing binary is the corect version -- for now only copy if file doesn't exist
	if !ncutils.FileExists(EXEC_DIR + "netclient") {
		err = ncutils.Copy(binarypath, EXEC_DIR+"netclient")
		if err != nil {
			log.Println(err)
			return err
		}
	}

	systemservice := `[Unit]
Description=Netclient Daemon
Documentation=https://docs.netmaker.org https://k8s.netmaker.org
After=network-online.target
Wants=network-online.target systemd-networkd-wait-online.service

[Service]
User=root
Type=simple
ExecStart=/sbin/netclient daemon
Restart=on-failure
RestartSec=15s

[Install]
WantedBy=multi-user.target
`

	servicebytes := []byte(systemservice)

	if !ncutils.FileExists("/etc/systemd/system/netclient.service") {
		err = os.WriteFile("/etc/systemd/system/netclient.service", servicebytes, 0644)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	_, _ = ncutils.RunCmd("systemctl enable netclient.service", true)
	_, _ = ncutils.RunCmd("systemctl daemon-reload", true)
	_, _ = ncutils.RunCmd("systemctl start netclient.service", true)
	return nil
}

// RestartSystemD - restarts systemd service
func RestartSystemD() {
	ncutils.PrintLog("restarting netclient.service", 1)
	time.Sleep(time.Second)
	_, _ = ncutils.RunCmd("systemctl restart netclient.service", true)
}

// CleanupLinux - cleans up neclient configs
func CleanupLinux() {
	if err := os.RemoveAll(ncutils.GetNetclientPath()); err != nil {
		ncutils.PrintLog("Removing netclient configs: "+err.Error(), 1)
	}
	if err := os.Remove(EXEC_DIR + "netclient"); err != nil {
		ncutils.PrintLog("Removing netclient binary: "+err.Error(), 1)
	}
}

// StopSystemD - tells system to stop systemd
func StopSystemD() {
	ncutils.RunCmd("systemctl stop netclient.service", false)
}

// RemoveSystemDServices - removes the systemd services on a machine
func RemoveSystemDServices() error {
	//sysExec, err := exec.LookPath("systemctl")
	var err error
	if !ncutils.IsWindows() && isOnlyService() {
		if err != nil {
			log.Println(err)
		}
		ncutils.RunCmd("systemctl disable netclient.service", false)
		ncutils.RunCmd("systemctl disable netclient.timer", false)
		if ncutils.FileExists("/etc/systemd/system/netclient.service") {
			err = os.Remove("/etc/systemd/system/netclient.service")
			if err != nil {
				ncutils.Log("Error removing /etc/systemd/system/netclient.service. Please investigate.")
			}
		}
		if ncutils.FileExists("/etc/systemd/system/netclient.timer") {
			err = os.Remove("/etc/systemd/system/netclient.timer")
			if err != nil {
				ncutils.Log("Error removing /etc/systemd/system/netclient.timer. Please investigate.")
			}
		}
		ncutils.RunCmd("systemctl daemon-reload", false)
		ncutils.RunCmd("systemctl reset-failed", false)
		ncutils.Log("removed systemd remnants if any existed")
	}
	return nil
}

func isOnlyService() bool {
	files, err := filepath.Glob("/etc/netclient/config/netconfig-*")
	if err != nil {
		return false
	}
	return len(files) == 0
}
