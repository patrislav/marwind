package main

import (
	"log"

	"github.com/patrislav/marwind"
	"github.com/patrislav/marwind/manager"
)

func main() {
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
