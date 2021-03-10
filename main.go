package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/metadata"
	"github.com/vharitonsky/iniflags"
)

var meta metadata.Server
var tz = time.FixedZone("CST", 8*60*60)
var debug bool

var svc = service.Service{
	Name: "Flows",
	Desc: "Auto fetching capital flows service",
	Exec: run,
	Options: service.Options{
		Dependencies: []string{"Wants=network-online.target", "After=network.target"},
	},
}

func main() {
	self, err := os.Executable()
	if err != nil {
		log.Fatalln("Failed to get self path:", err)
	}

	flag.StringVar(&meta.Addr, "server", "", "Metadata Server Address")
	flag.StringVar(&meta.Header, "header", "", "Verify Header Header Name")
	flag.StringVar(&meta.Value, "value", "", "Verify Header Value")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	flag.BoolVar(&debug, "debug", false, "debug")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.SetAllowUnknownFlags(true)
	iniflags.Parse()

	if service.IsWindowsService() {
		svc.Run(false)
		return
	}

	switch flag.NArg() {
	case 0:
		run()
	case 1:
		switch flag.Arg(0) {
		case "run", "debug":
			run()
		case "install":
			err = svc.Install()
		case "remove":
			err = svc.Remove()
		case "start":
			err = svc.Start()
		case "stop":
			err = svc.Stop()
		case "restart":
			err = svc.Restart()
		case "update":
			err = svc.Update()
		default:
			log.Fatalln("Unknown argument:", flag.Arg(0))
		}
	default:
		log.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		log.Fatalf("failed to %s: %v", flag.Arg(0), err)
	}
}
