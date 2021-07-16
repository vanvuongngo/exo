package tui

func DarkMode() bool {
	// MAC: defaults read -g AppleInterfaceStyle
	// WINDOWS: https://gist.github.com/jerblack/1d05bbcebb50ad55c312e4d7cf1bc909
	// LINUX: ???
	return true // This is the default for tview.
}

func LightMode() bool {
	return !DarkMode()
}
