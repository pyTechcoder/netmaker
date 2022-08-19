package ncutils

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

//go:embed windowsdaemon/winsw.exe
var winswContent embed.FS

// RunCmd - runs a local command
func RunCmd(command string, printerr bool) (string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Wait()
	//cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "/C \"" + command + "\""}
	out, err := cmd.CombinedOutput()
	if err != nil && printerr {
		log.Println("error running command:", command)
		log.Println(strings.TrimSuffix(string(out), "\n"))
	}
	return string(out), err
}

// RunCmd - runs a local command
func RunCmdFormatted(command string, printerr bool) (string, error) {
	var comSpec = os.Getenv("COMSPEC")
	if comSpec == "" {
		comSpec = os.Getenv("SystemRoot") + "\\System32\\cmd.exe"
	}
	cmd := exec.Command(comSpec)
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "/C \"" + command + "\""}
	cmd.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil && printerr {
		log.Println("error running command:", command)
		log.Println(strings.TrimSuffix(string(out), "\n"))
	}
	return string(out), err
}

// GetEmbedded - Gets the Windows daemon creator
func GetEmbedded() error {
	data, err := winswContent.ReadFile("windowsdaemon/winsw.exe")
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%swinsw.exe", GetNetclientPathSpecific())
	err = os.WriteFile(fileName, data, 0700)
	if err != nil {
		Log("could not mount winsw.exe")
		return err
	}
	return nil
}
