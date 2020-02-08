package marwind

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/patrislav/marwind/keysym"
	"github.com/patrislav/marwind/wm"
)

var Config = wm.Config{
	InnerGap:                4,
	OuterGap:                4,
	Shell:                   "/bin/sh",
	LauncherCommand:         "rofi -show drun",
	TerminalCommand:         "alacritty",
	BorderWidth:             0,
	BorderColor:             0xffa1d1cf,
	TitleBarHeight:          18,
	TitleBarBgColor:         0xffa1d1cf,
	TitleBarFontColorActive: 0xff000000,
	TitleBarFontSize:        12,
	Keybindings: map[xproto.Keysym]string{
		// Brightness control
		keysym.XF86MonBrightnessDown: "light -U 5",
		keysym.XF86MonBrightnessUp:   "light -A 5",
		// Volume control
		keysym.XF86AudioMute: "pactl set-sink-mute @DEFAULT_SINK@ toggle",
		keysym.XF86AudioLowerVolume: "pactl set-sink-volume @DEFAULT_SINK@ -5%",
		keysym.XF86AudioRaiseVolume: "pactl set-sink-volume @DEFAULT_SINK@ +5%",
	},
}
