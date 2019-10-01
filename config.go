package marwind

import (
	"github.com/patrislav/marwind/wm"
)

var Config = wm.Config{
	InnerGap:                4,
	OuterGap:                4,
	Shell:                   "/bin/sh",
	LauncherCommand:         "rofi -show drun",
	TerminalCommand:         "alacritty",
	BorderWidth:             1,
	BorderColor:             0xffa1d1cf,
	TitleBarHeight:          18,
	TitleBarBgColor:         0xffa1d1cf,
	TitleBarFontColorActive: 0xff000000,
	TitleBarFontSize:        12,
}
