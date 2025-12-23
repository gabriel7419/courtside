package data

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const settingsFileName = "settings.yaml"

// LeagueInfo contains league metadata for display purposes.
type LeagueInfo struct {
	ID      int
	Name    string
	Country string
}

// AllSupportedLeagues contains all leagues that Golazo supports.
// This is the source of truth for available leagues.
var AllSupportedLeagues = []LeagueInfo{
	// Top 5 European Leagues
	{ID: 47, Name: "Premier League", Country: "England"},
	{ID: 87, Name: "La Liga", Country: "Spain"},
	{ID: 54, Name: "Bundesliga", Country: "Germany"},
	{ID: 55, Name: "Serie A", Country: "Italy"},
	{ID: 53, Name: "Ligue 1", Country: "France"},
	// European Competitions
	{ID: 42, Name: "UEFA Champions League", Country: "Europe"},
	{ID: 73, Name: "UEFA Europa League", Country: "Europe"},
	{ID: 50, Name: "UEFA Euro", Country: "Europe"},
	// South America
	{ID: 268, Name: "Brasileirão Série A", Country: "Brazil"},
	{ID: 112, Name: "Liga Profesional", Country: "Argentina"},
	{ID: 14, Name: "Copa Libertadores", Country: "South America"},
	{ID: 44, Name: "Copa America", Country: "South America"},
	// Other
	{ID: 130, Name: "MLS", Country: "USA"},
	{ID: 77, Name: "FIFA World Cup", Country: "International"},
}

// Settings represents user preferences stored in settings.yaml.
type Settings struct {
	// SelectedLeagues contains the IDs of leagues the user wants to follow.
	// If empty, all supported leagues are used.
	SelectedLeagues []int `yaml:"selected_leagues"`
}

// SettingsPath returns the path to the settings file.
func SettingsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsFileName), nil
}

// LoadSettings reads settings from the settings.yaml file.
// Returns default settings (empty selection = all leagues) if file doesn't exist.
func LoadSettings() (*Settings, error) {
	path, err := SettingsPath()
	if err != nil {
		return &Settings{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No settings file - return empty settings (will use all leagues)
			return &Settings{}, nil
		}
		return &Settings{}, err
	}

	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		// Invalid YAML - return empty settings
		return &Settings{}, nil
	}

	return &settings, nil
}

// SaveSettings writes settings to the settings.yaml file.
func SaveSettings(settings *Settings) error {
	path, err := SettingsPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetActiveLeagueIDs returns the league IDs that should be used for API calls.
// If no leagues are selected in settings, returns all supported league IDs.
func GetActiveLeagueIDs() []int {
	settings, err := LoadSettings()
	if err != nil || len(settings.SelectedLeagues) == 0 {
		// Return all league IDs as default
		ids := make([]int, len(AllSupportedLeagues))
		for i, league := range AllSupportedLeagues {
			ids[i] = league.ID
		}
		return ids
	}

	return settings.SelectedLeagues
}

// GetAllLeagueIDs returns all supported league IDs (used as fallback).
func GetAllLeagueIDs() []int {
	ids := make([]int, len(AllSupportedLeagues))
	for i, league := range AllSupportedLeagues {
		ids[i] = league.ID
	}
	return ids
}

// IsLeagueSelected checks if a league ID is in the selected list.
func (s *Settings) IsLeagueSelected(leagueID int) bool {
	for _, id := range s.SelectedLeagues {
		if id == leagueID {
			return true
		}
	}
	return false
}
