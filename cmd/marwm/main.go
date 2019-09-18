package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	flag "github.com/spf13/pflag"

	"github.com/patrislav/marwind"
	"github.com/patrislav/marwind/manager"
)

var (
	version   string // program version
	buildTime string // when the executable was built
)

var (
	flagVersion bool
)

func main() {
	flag.BoolVar(&flagVersion, "version", false, "show version and exit")
	flag.Parse()

	if flagVersion {
		fmt.Printf("marwm version:\t%s (%s)\n", version, buildTime)
		fmt.Printf("go version:\t%s\n", runtime.Version())
		os.Exit(0)
	}

	mgr, err := manager.New(manager.Config{
		InnerGap: marwind.InnerGap,
		OuterGap: marwind.OuterGap,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer mgr.Close()
	if err := mgr.Init(); err != nil {
		log.Fatal(err)
	}
	if err := mgr.Run(); err != nil {
		log.Fatal(err)
	}
}
