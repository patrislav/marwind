package marwind

import (
	"github.com/patrislav/marwind/wm"
)

var Config = wm.Config{
	InnerGap:        4,
	OuterGap:        2,
	Shell:           "/bin/sh",
	LauncherCommand: "rofi -show drun",
	TerminalCommand: "kitty",
}
