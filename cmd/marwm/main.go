package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	initCmd     string
)

func main() {
	flag.BoolVar(&flagVersion, "version", false, "show version and exit")
	flag.StringVar(&initCmd, "init", "", "run this executable at startup")
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

	if initCmd != "" {
		cmd := exec.Command(initCmd)
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		go cmd.Wait()
	}

	if err := mgr.Run(); err != nil {
		log.Fatal(err)
	}
}
