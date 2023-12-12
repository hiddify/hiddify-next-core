package global

import (
	"errors"
	"fmt"
	"github.com/hiddify/libcore/shared"
	"github.com/sagernet/sing-box/experimental/libbox"
	"github.com/sagernet/sing-box/log"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var box *libbox.BoxService
var configOptions *shared.ConfigOptions
var logFactory *log.Factory

func StartServiceC(delayStart bool, content string) error {
	options, err := parseConfig(content)
	if err != nil {
		return errors.New(stopAndAlert(EmptyConfiguration, err))
	}
	configOptions = &shared.ConfigOptions{}
	options = shared.BuildConfig(*configOptions, options)

	err = shared.SaveCurrentConfig(sWorkingPath, options)
	if err != nil {
		return err
	}

	err = startCommandServer(*logFactory)
	if err != nil {
		return errors.New(stopAndAlert(StartCommandServer, err))
	}

	instance, err := NewService(options)
	if err != nil {
		return errors.New(stopAndAlert(CreateService, err))
	}

	if delayStart {
		time.Sleep(250 * time.Millisecond)
	}

	err = instance.Start()
	if err != nil {
		return errors.New(stopAndAlert(StartService, err))
	}
	box = instance
	commandServer.SetService(box)

	propagateStatus(Started)
	return nil
}

func StopService() error {
	if status != Started {
		return nil
	}
	if box == nil {
		return errors.New("instance not found")
	}

	propagateStatus(Stopping)
	commandServer.SetService(nil)
	err := box.Close()
	if err != nil {
		return err
	}
	box = nil

	err = commandServer.Close()
	if err != nil {
		return err
	}
	commandServer = nil
	propagateStatus(Stopped)

	return nil
}

func SetupC(baseDir string, workDir string, tempDir string, debug bool) error {
	err := os.MkdirAll("./bin", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll("./work", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll("./temp", os.ModePerm)
	if err != nil {
		return err
	}
	Setup(baseDir, workDir, tempDir)
	var defaultWriter io.Writer
	if !debug {
		defaultWriter = io.Discard
	}
	factory, err := log.New(
		log.Options{
			DefaultWriter: defaultWriter,
			BaseTime:      time.Now(),
			Observable:    false,
		})
	if err != nil {
		return err
	}
	logFactory = &factory
	return nil
}

func MakeConfig(Ipv6 bool, ServerPort int, StrictRoute bool, EndpointIndependentNat bool, Stack string) string {
	var ipv6 string
	if Ipv6 {
		ipv6 = "      \"inet6_address\": \"fdfe:dcba:9876::1/126\",\n"
	} else {
		ipv6 = ""
	}
	base := "{\n  \"inbounds\": [\n    {\n      \"type\": \"tun\",\n      \"tag\": \"tun-in\",\n      \"interface_name\": \"tun0\",\n      \"inet4_address\": \"172.19.0.1/30\",\n" + ipv6 + "      \"mtu\": 9000,\n      \"auto_route\": true,\n      \"strict_route\": " + fmt.Sprintf("%t", StrictRoute) + ",\n      \"endpoint_independent_nat\": " + fmt.Sprintf("%t", EndpointIndependentNat) + ",\n      \"stack\": \"" + Stack + "\"\n    }],\n  \"outbounds\": [\n    {\n      \"type\": \"socks\",\n      \"tag\": \"socks-out\",\n      \"server\": \"127.0.0.1\",\n      \"server_port\": " + fmt.Sprintf("%d", ServerPort) + ",\n      \"version\": \"5\"\n    }\n  ]\n}\n"
	return base
}

func WriteParameters(Ipv6 bool, ServerPort int, StrictRoute bool, EndpointIndependentNat bool, Stack string) error {
	parameters := fmt.Sprintf("%t,%d,%t,%t,%s", Ipv6, ServerPort, StrictRoute, EndpointIndependentNat, Stack)
	err := os.WriteFile("bin/parameters.config", []byte(parameters), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
func ReadParameters() (bool, int, bool, bool, string, error) {
	Data, err := os.ReadFile("bin/parameters.config")
	if err != nil {
		return false, 0, false, false, "", err
	}
	DataSlice := strings.Split(string(Data), ",")
	Ipv6, _ := strconv.ParseBool(DataSlice[0])
	ServerPort, _ := strconv.Atoi(DataSlice[1])
	StrictRoute, _ := strconv.ParseBool(DataSlice[2])
	EndpointIndependentNat, _ := strconv.ParseBool(DataSlice[3])
	Stack := DataSlice[4]
	return Ipv6, ServerPort, StrictRoute, EndpointIndependentNat, Stack, nil
}
