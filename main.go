package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/service"
	"github.com/vharitonsky/iniflags"
)

var meta metadata.Server
var tz = time.FixedZone("CST", 8*60*60)
var debug bool

var svc = service.Service{
	Name:     "Flows",
	Desc:     "Auto fetching capital flows service",
	Exec:     run,
	TestExec: test,
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
		cmd := flag.Arg(0)
		var ok bool
		if ok, err = svc.Command(cmd); !ok {
			log.Fatalln("Unknown argument:", cmd)
		}
	default:
		log.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		log.Fatalf("failed to %s: %v", flag.Arg(0), err)
	}
}
