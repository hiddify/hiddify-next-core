package main

import (
	"embed"
	"fmt"
	"github.com/hiddify/libcore/shared"
	"github.com/sagernet/sing-box/experimental/libbox"
	"github.com/sagernet/sing-box/log"
	"os"
	"time"
)

//go:embed bin/*
var bin embed.FS
var box *libbox.BoxService
var configOptions *shared.ConfigOptions
var logFactory *log.Factory

func startService(delayStart bool) error {
	content, err := bin.ReadFile("bin/config.json")
	if err != nil {
		return stopAndAlert(EmptyConfiguration, err)
	}
	options, err := parseConfig(string(content))
	if err != nil {
		return stopAndAlert(EmptyConfiguration, err)
	}
	options = shared.BuildConfig(*configOptions, options)

	err = startCommandServer(*logFactory)
	if err != nil {
		return stopAndAlert(StartCommandServer, err)
	}

	instance, err := NewService(options)
	if err != nil {
		return stopAndAlert(CreateService, err)
	}

	if delayStart {
		time.Sleep(250 * time.Millisecond)
	}

	err = instance.Start()
	if err != nil {
		return stopAndAlert(StartService, err)
	}
	box = instance
	commandServer.SetService(box)

	propagateStatus(Started)
	return nil
}

func main() {
	args := os.Args
	switch args[1] {
	case "service-mode":
		err := startService(false)
		if err != nil {
			fmt.Println("Error:", err)
		}
	default:
		usage, _ := bin.ReadFile("bin/usage.txt")
		fmt.Println("Usage:", string(usage))
	}
}
