package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/hiddify/libcore/global"
	"github.com/hiddify/libcore/utils"
	"github.com/hiddify/libcore/web"
	"github.com/kardianos/service"
	"os"
)

//go:embed bin/*
var bin embed.FS

type hiddifyNextCore struct{}

func (m *hiddifyNextCore) Start(s service.Service) error {
	go m.run()
	return nil
}
func (m *hiddifyNextCore) Stop(s service.Service) error {
	err := global.StopService()
	if err != nil {
		return err
	}
	return nil
}
func (m *hiddifyNextCore) run() {
	Port := flag.Int("port", 6548, "")
	web.StartWebServer(*Port)
}

func main() {
	args := os.Args
	svcConfig := &service.Config{
		Name:        "hiddify_next_core",
		DisplayName: "hiddify next core",
		Description: "@hiddify_com set this",
	}
	prg := &hiddifyNextCore{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if len(args) > 1 {
		switch args[1] {
		case "start-service":
			err = svc.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "install-service":
			err = svc.Install()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "gen-cert":
			err := os.MkdirAll("cert", os.ModePerm)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			utils.GenerateCertificate("cert/server-cert.pem", "cert/server-key.pem", true)
			utils.GenerateCertificate("cert/client-cert.pem", "cert/client-key.pem", false)
		}
	} else {
		fmt.Println("Error:", "not enough parameters")
		usage, _ := bin.ReadFile("bin/usage.txt")
		fmt.Println(string(usage))
	}

}
