package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/hiddify/libcore/utils"
	"github.com/hiddify/libcore/web"
	"os"
)

//go:embed bin/*
var bin embed.FS

func main() {
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "start-service":
			Port := flag.Int("port", 6548, "")
			web.StartWebServer(*Port)
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
