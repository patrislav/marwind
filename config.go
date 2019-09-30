package marwind

import (
	"github.com/patrislav/marwind/wm"
)

var Config = wm.Config{
	InnerGap:        4,
	OuterGap:        4,
	Shell:           "/bin/sh",
	LauncherCommand: "rofi -show drun",
	TerminalCommand: "alacritty",
	BorderWidth:     1,
	BorderColor:     0xff9eeeee,

	TitleBarHeight: 18,
}
