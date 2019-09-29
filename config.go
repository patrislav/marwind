package marwind

import (
	"github.com/patrislav/marwind/wm"
)

var Config = wm.Config{
	InnerGap:        4,
	OuterGap:        4,
	Shell:           "/bin/sh",
	LauncherCommand: "rofi -show drun",
	TerminalCommand: "kitty",
	BorderWidth:     1,
	BorderColor:     0xff2f343f,
}
