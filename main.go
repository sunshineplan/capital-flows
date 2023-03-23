package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/flags"
)

var (
	svc  = service.New()
	meta metadata.Server
)

func init() {
	svc.Name = "Flows"
	svc.Desc = "Auto fetching capital flows service"
	svc.Exec = run
	svc.TestExec = test
	svc.Options = service.Options{
		Dependencies: []string{"Wants=network-online.target", "After=network.target"},
	}
}

var debug = flag.Bool("debug", false, "debug")

func main() {
	self, err := os.Executable()
	if err != nil {
		svc.Fatalln("Failed to get self path:", err)
	}

	flag.StringVar(&meta.Addr, "server", "", "Metadata Server Address")
	flag.StringVar(&meta.Header, "header", "", "Verify Header Header Name")
	flag.StringVar(&meta.Value, "value", "", "Verify Header Value")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	flags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	flags.Parse()

	if service.IsWindowsService() {
		svc.Run()
		return
	}

	switch flag.NArg() {
	case 0:
		svc.Run()
	case 1:
		cmd := flag.Arg(0)
		var ok bool
		if ok, err = svc.Command(cmd); !ok {
			svc.Fatalln("Unknown argument:", cmd)
		}
	default:
		svc.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		svc.Fatalf("failed to %s: %v", flag.Arg(0), err)
	}
}
