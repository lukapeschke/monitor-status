package main

import (
	"context"
	"errors"
	"flag"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/jochenvg/go-udev"
)

func getDeviceStatus(d *udev.Device) string {
	statusFile := d.Syspath() + "/status"
	status, err := ioutil.ReadFile(statusFile)
	if err != nil {
		ErrLog.Println(err)
		return ""
	}
	InfoLog.Printf("Device %s has status %s\n", d.Sysname(), status)
	return string(status)
}

var (
	// ErrLog is a logger for errors
	ErrLog *log.Logger
	// InfoLog is a logger for infos
	InfoLog *log.Logger
)

func initLog() {
	ErrLog = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
	InfoLog = log.New(os.Stdout, "INFO: ", 0)
}

func getConfigFile(configFile string) (filename string, err error) {
	filenames := []string{
		configFile,
		"./config.yml",
		"~/.monitor-status/config.yml",
		"~/.config/monitor-status.yml",
		"~/.config/monitor-status/config.yml",
		build.Default.GOPATH +
			"/src/github.com/lukapeschke/monitor-status/config.yml",
	}

	for _, name := range filenames {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
	}
	return "", errors.New("No config file found")
}

func getDevices() (*map[string]map[string]YamlDevice, error) {
	var configFile string
	var err error

	flag.StringVar(&configFile, "config", "", "Config file to use")
	flag.Parse()
	if configFile, err = getConfigFile(configFile); err != nil {
		ErrLog.Println("Could not find config file")
		return nil, err
	}
	InfoLog.Println("Found config file", configFile)

	devices, err := loadYamlDevices(configFile)
	if err != nil {
		ErrLog.Println("Could not load file", configFile, err)
		return nil, err
	}
	return devices, nil
}

func handleDeviceEvent(dev *udev.Device, action *YamlDevice) {
	InfoLog.Printf("%s: %s\n", dev.Sysname(), dev.Action())
	if dev.Action() == "add" {
		cmd := exec.Command("sh", "-c", action.OnConnect)
		InfoLog.Printf("Executing \"%s\"\n", action.OnConnect)
		if err := cmd.Run(); err != nil {
			ErrLog.Printf("Couldn't exec |%s|", action.OnConnect)
		}
	} else if dev.Action() == "remove" {
		cmd := exec.Command("sh", "-c", action.OnDisconnect)
		InfoLog.Printf("Executing \"%s\"\n", action.OnDisconnect)
		if err := cmd.Run(); err != nil {
			ErrLog.Printf("Couldn't exec |%s|", action.OnDisconnect)
		}
	}
}

func handleEvents(ch <-chan *udev.Device, wg *sync.WaitGroup,
	devices *map[string]YamlDevice) {
	defer wg.Done()
	for d := range ch {
		for key, value := range *devices {
			if key == d.Sysname() {
				handleDeviceEvent(d, &value)
			}
		}
	}
}

func run(config *map[string]map[string]YamlDevice) {

	u := udev.Udev{}
	var wg sync.WaitGroup

	for subsystem, devices := range *config {
		if len(devices) < 1 {
			ErrLog.Printf(
				"No device to watch in subsystem %s, ignoring...", subsystem)
		}

		// Creating context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Creating monitor
		mon := u.NewMonitorFromNetlink("udev")
		if err := mon.FilterAddMatchSubsystem(subsystem); err != nil {
			panic("Couldn't add subsystem device type filter")
		}

		ch, err := mon.DeviceChan(ctx)
		if err != nil {
			panic("Couldn't create chan")
		}

		InfoLog.Printf("Listening to events for subsystem %s...\n", subsystem)
		wg.Add(1)
		go handleEvents(ch, &wg, &devices)
	}
	wg.Wait()
}

func main() {
	initLog()

	devices, err := getDevices()
	if err != nil {
		return
	}
	run(devices)
}
