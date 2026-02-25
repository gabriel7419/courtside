package constants

// Menu items
const (
	MenuStats       = "Finished Games"
	MenuLiveMatches = "Live Games"
	MenuSettings    = "Settings"
)

// Panel titles
const (
	PanelLiveMatches       = "Live Games"
	PanelFinishedMatches   = "Finished Games"
	PanelMatchDetails      = "Game Details"
	PanelMatchList         = "Game List"
	PanelUpcomingMatches   = "Upcoming Games"
	PanelPlayByPlay        = "Play-by-play"
	PanelGameStatistics    = "Game Statistics"
	PanelUpdates           = "Live Updates"
	PanelLeaguePreferences = "Conference Preferences"
)

// Backward-compat aliases (used in older callers)
const (
	PanelMinuteByMinute  = PanelPlayByPlay
	PanelMatchStatistics = PanelGameStatistics
)

// Empty state messages
const (
	EmptyNoLiveMatches     = "No live games right now"
	EmptyNoFinishedMatches = "No finished games"
	EmptySelectMatch       = "Select a game"
	EmptyNoUpdates         = "No play-by-play yet"
	EmptyNoMatches         = "No games available"
)

// Help text
const (
	HelpMainMenu           = "↑/↓: navigate  Enter: select  q: quit"
	HelpMatchesView        = "↑/↓: navigate  r: refresh  /: filter  Esc: back  q: quit"
	HelpSettingsView       = "↑/↓: navigate  ←/→: switch tabs  Space: toggle  /: filter  Enter: save  Esc: back"
	HelpStatsView          = "h/l: date range  j/k: navigate  Tab: focus details  ↑/↓: scroll when focused  r: refresh  /: filter  Esc: back"
	HelpStatsViewUnfocused = "Tab: focus details"
	HelpStatsViewFocused   = "Tab: unfocus  s: standings  x: all statistics  ↑/↓: scroll"
	HelpStandingsDialog    = "Esc: close"
	HelpFormationsDialog   = "Tab/←/→: switch team  Esc: close"
	HelpStatisticsDialog   = "↑/↓: navigate  Esc: close"
)

// Status text
const (
	StatusLive            = "LIVE"
	StatusFinished        = "Final"
	StatusNotStarted      = "VS"
	StatusNotStartedShort = "NS"
	StatusFinishedText    = "Final"
)

// Loading text
const (
	LoadingFetching = "Fetching..."
)

// Notification text
const (
	// NotificationTitleGoal is shown in scoring notifications.
	NotificationTitleGoal  = "🏀 Courtside!"
	NotificationTitleScore = "🏀 Score!"
)

// Stats labels
const (
	LabelStatus = "Status: "
	LabelScore  = "Score: "
	LabelLeague = "Conference: "
	LabelDate   = "Date: "
	LabelVenue  = "Arena: "
)
