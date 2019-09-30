package wm

type Config struct {
	InnerGap uint32 // Gap around each window, in pixels
	OuterGap uint32 // Additional gap around the entire workspace, in pixels

	Shell string // Name of the program to use for executing commands ("/bin/sh" by default)

	// Shell command to execute after using the "Launcher" binding (Win + D by default)
	LauncherCommand string
	// Shell command to execute after using the "Terminal" binding (Win + Shift + Enter by default)
	TerminalCommand string

	BorderWidth uint8
	BorderColor uint32

	TitleBarHeight            uint8
	TitleBarBgColor           uint32
	TitleBarFontColorActive   uint32
	TitleBarFontColorInactive uint32
	TitleBarFontSize          uint8
}
