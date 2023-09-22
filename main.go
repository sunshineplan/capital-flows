package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"

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

	if *debug {
		svc.SetLevel(slog.LevelDebug)
	}

	if err := svc.ParseAndRun(flag.Args()); err != nil {
		svc.Fatal(err)
	}
}
