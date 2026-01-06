package constants

import "time"

// MainViewCheckDelay is the delay before navigating to a selected view in the main menu.
// Set to 1.5 seconds to allow API preloading while showing transition animation.
const MainViewCheckDelay = 1500 * time.Millisecond

// StatusBannerType represents the type of status banner to display at the top of views.
type StatusBannerType int

const (
	// StatusBannerNone indicates no status banner should be displayed.
	StatusBannerNone StatusBannerType = iota
	// StatusBannerDebug indicates debug mode is active and logs are being written.
	StatusBannerDebug
	// StatusBannerNewVersion indicates a new version of Golazo is available.
	StatusBannerNewVersion
	// StatusBannerDev indicates this is a development build.
	StatusBannerDev
)
