package keysym

// Known KeySyms from /usr/include/X11/keysymdef.h
// Copied from https://github.com/driusan/dewm/blob/master/keysym/keysym.go
const (
	// TTY function keys, cleverly chosen to map to ASCII, for convenience of
	// programming, but could have been arbitrary (at the cost of lookup
	// tables in client code).
	XKBackSpace  = 0xff08 // Back space, back char
	XKTab        = 0xff09
	XKLinefeed   = 0xff0a // Linefeed, LF
	XKClear      = 0xff0b
	XKReturn     = 0xff0d // Return, enter
	XKPause      = 0xff13 // Pause, hold
	XKScrollLock = 0xff14
	XKSysReq     = 0xff15
	XKEscape     = 0xff1b
	XKDelete     = 0xffff // Delete, rubout
	// Latin 1
	// (ISO/IEC 8859-1 = Unicode U+0020..U+00FF)
	// Byte 3 = 0
	XKSpace        = 0x0020 // U+0020 SPACE
	XKExclam       = 0x0021 // U+0021 EXCLAMATION MARK
	XKQuotedbl     = 0x0022 // U+0022 QUOTATION MARK
	XKNumberSign   = 0x0023 // U+0023 NUMBER SIGN
	XKDollar       = 0x0024 // U+0024 DOLLAR SIGN
	XKPercent      = 0x0025 // U+0025 PERCENT SIGN
	XKAmpersand    = 0x0026 // U+0026 AMPERSAND
	XKApostrophe   = 0x0027 // U+0027 APOSTROPHE
	XKQuoteRight   = 0x0027 // deprecated
	XKParenLeft    = 0x0028 // U+0028 LEFT PARENTHESIS
	XKParenRight   = 0x0029 // U+0029 RIGHT PARENTHESIS
	XKAsterisk     = 0x002a // U+002A ASTERISK
	XKPlus         = 0x002b // U+002B PLUS SIGN
	XKComma        = 0x002c // U+002C COMMA
	XKMinus        = 0x002d // U+002D HYPHEN-MINUS
	XKPeriod       = 0x002e // U+002E FULL STOP
	XKSlash        = 0x002f // U+002F SOLIDUS
	XK0            = 0x0030 // U+0030 DIGIT ZERO
	XK1            = 0x0031 // U+0031 DIGIT ONE
	XK2            = 0x0032 // U+0032 DIGIT TWO
	XK3            = 0x0033 // U+0033 DIGIT THREE
	XK4            = 0x0034 // U+0034 DIGIT FOUR
	XK5            = 0x0035 // U+0035 DIGIT FIVE
	XK6            = 0x0036 // U+0036 DIGIT SIX
	XK7            = 0x0037 // U+0037 DIGIT SEVEN
	XK8            = 0x0038 // U+0038 DIGIT EIGHT
	XK9            = 0x0039 // U+0039 DIGIT NINE
	XKColon        = 0x003a // U+003A COLON
	XKSemicolon    = 0x003b // U+003B SEMICOLON
	XKLess         = 0x003c // U+003C LESS-THAN SIGN
	XKEqual        = 0x003d // U+003D EQUALS SIGN
	XKGreater      = 0x003e // U+003E GREATER-THAN SIGN
	XKQuestion     = 0x003f // U+003F QUESTION MARK
	XKAt           = 0x0040 // U+0040 COMMERCIAL AT
	XKA            = 0x0041 // U+0041 LATIN CAPITAL LETTER A
	XKB            = 0x0042 // U+0042 LATIN CAPITAL LETTER B
	XKC            = 0x0043 // U+0043 LATIN CAPITAL LETTER C
	XKD            = 0x0044 // U+0044 LATIN CAPITAL LETTER D
	XKE            = 0x0045 // U+0045 LATIN CAPITAL LETTER E
	XKF            = 0x0046 // U+0046 LATIN CAPITAL LETTER F
	XKG            = 0x0047 // U+0047 LATIN CAPITAL LETTER G
	XKH            = 0x0048 // U+0048 LATIN CAPITAL LETTER H
	XKI            = 0x0049 // U+0049 LATIN CAPITAL LETTER I
	XKJ            = 0x004a // U+004A LATIN CAPITAL LETTER J
	XKK            = 0x004b // U+004B LATIN CAPITAL LETTER K
	XKL            = 0x004c // U+004C LATIN CAPITAL LETTER L
	XKM            = 0x004d // U+004D LATIN CAPITAL LETTER M
	XKN            = 0x004e // U+004E LATIN CAPITAL LETTER N
	XKO            = 0x004f // U+004F LATIN CAPITAL LETTER O
	XKP            = 0x0050 // U+0050 LATIN CAPITAL LETTER P
	XKQ            = 0x0051 // U+0051 LATIN CAPITAL LETTER Q
	XKR            = 0x0052 // U+0052 LATIN CAPITAL LETTER R
	XKS            = 0x0053 // U+0053 LATIN CAPITAL LETTER S
	XKT            = 0x0054 // U+0054 LATIN CAPITAL LETTER T
	XKU            = 0x0055 // U+0055 LATIN CAPITAL LETTER U
	XKV            = 0x0056 // U+0056 LATIN CAPITAL LETTER V
	XKW            = 0x0057 // U+0057 LATIN CAPITAL LETTER W
	XKX            = 0x0058 // U+0058 LATIN CAPITAL LETTER X
	XKY            = 0x0059 // U+0059 LATIN CAPITAL LETTER Y
	XKZ            = 0x005a // U+005A LATIN CAPITAL LETTER Z
	XKBracketLeft  = 0x005b // U+005B LEFT SQUARE BRACKET
	XKBackslash    = 0x005c // U+005C REVERSE SOLIDUS
	XKBracketRight = 0x005d // U+005D RIGHT SQUARE BRACKET
	XKAsciiCircum  = 0x005e // U+005E CIRCUMFLEX ACCENT
	XKUnderscore   = 0x005f // U+005F LOW LINE
	XKGrave        = 0x0060 // U+0060 GRAVE ACCENT
	XKQuoteLeft    = 0x0060 // deprecated
	XKa            = 0x0061 // U+0061 LATIN SMALL LETTER A
	XKb            = 0x0062 // U+0062 LATIN SMALL LETTER B
	XKc            = 0x0063 // U+0063 LATIN SMALL LETTER C
	XKd            = 0x0064 // U+0064 LATIN SMALL LETTER D
	XKe            = 0x0065 // U+0065 LATIN SMALL LETTER E
	XKf            = 0x0066 // U+0066 LATIN SMALL LETTER F
	XKg            = 0x0067 // U+0067 LATIN SMALL LETTER G
	XKh            = 0x0068 // U+0068 LATIN SMALL LETTER H
	XKi            = 0x0069 // U+0069 LATIN SMALL LETTER I
	XKj            = 0x006a // U+006A LATIN SMALL LETTER J
	XKk            = 0x006b // U+006B LATIN SMALL LETTER K
	XKl            = 0x006c // U+006C LATIN SMALL LETTER L
	XKm            = 0x006d // U+006D LATIN SMALL LETTER M
	XKn            = 0x006e // U+006E LATIN SMALL LETTER N
	XKo            = 0x006f // U+006F LATIN SMALL LETTER O
	XKp            = 0x0070 // U+0070 LATIN SMALL LETTER P
	XKq            = 0x0071 // U+0071 LATIN SMALL LETTER Q
	XKr            = 0x0072 // U+0072 LATIN SMALL LETTER R
	XKs            = 0x0073 // U+0073 LATIN SMALL LETTER S
	XKt            = 0x0074 // U+0074 LATIN SMALL LETTER T
	XKu            = 0x0075 // U+0075 LATIN SMALL LETTER U
	XKv            = 0x0076 // U+0076 LATIN SMALL LETTER V
	XKw            = 0x0077 // U+0077 LATIN SMALL LETTER W
	XKx            = 0x0078 // U+0078 LATIN SMALL LETTER X
	XKy            = 0x0079 // U+0079 LATIN SMALL LETTER Y
	XKz            = 0x007a // U+007A LATIN SMALL LETTER Z
	XKBraceLeft    = 0x007b // U+007B LEFT CURLY BRACKET
	XKBar          = 0x007c // U+007C VERTICAL LINE
	XKBraceRight   = 0x007d // U+007D RIGHT CURLY BRACKET
	XKAsciiTilde   = 0x007e // U+007E TILDE

	// Cursor control & motion
	XKHome     = 0xff50
	XKLeft     = 0xff51 // Move left, left arrow
	XKUp       = 0xff52 // Move up, up arrow
	XKRight    = 0xff53 // Move right, right arrow
	XKDown     = 0xff54 // Move down, down arrow
	XKPrior    = 0xff55 // Prior, previous
	XKPageUp   = 0xff55
	XKNext     = 0xff56 // Next
	XKPageDown = 0xff56
	XKEnd      = 0xff57 // EOL
	XKBegin    = 0xff58 // BOL

	XF86MonBrightnessUp   = 0x1008ff02
	XF86MonBrightnessDown = 0x1008ff03
	XF86AudioLowerVolume  = 0x1008ff11
	XF86AudioMute         = 0x1008ff12
	XF86AudioRaiseVolume  = 0x1008ff13
)
