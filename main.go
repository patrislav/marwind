package main

import (
	"log"

	"github.com/patrislav/marwind-wm/manager"
)

func main() {
	mgr, err := manager.New(manager.Config{
		InnerGap: innerGap,
		OuterGap: outerGap,
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
