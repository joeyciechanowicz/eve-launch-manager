package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	CONFIG_FILE = "eve-launch-manager.json"
)

// Config represents the configuration structure
type Config struct {
	ActiveProfile string   `json:"activeProfile"`
	Profiles      []string `json:"profiles"`
}

// DefaultConfig provides the default configuration
var DefaultConfig = Config{
	ActiveProfile: "main",
	Profiles:      []string{"main"},
}

// ProfileManager handles EVE Online profile management
type ProfileManager struct {
	Config Config
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, CONFIG_FILE)
}

// getEvePath returns the path to EVE Online directory
func getEvePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "AppData", "Roaming", "EVE Online")
}

// getStateFilePath returns the path to a profile's state file
func getStateFilePath(profileName string) string {
	return filepath.Join(getEvePath(), "state-"+profileName+".json")
}

// getLauncherStateFilePath returns the path to the launcher's state file
func getLauncherStateFilePath() string {
	return filepath.Join(getEvePath(), "state.json")
}

// NewProfileManager creates a new ProfileManager instance
func NewProfileManager() (*ProfileManager, error) {
	config := DefaultConfig

	configPath := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	} else {
		// Create default config file
		data, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, err
		}

		// Copy initial state file
		sourceFile := getLauncherStateFilePath()
		destFile := getStateFilePath("main")
		if err := copyFile(sourceFile, destFile); err != nil {
			return nil, err
		}
	}

	return &ProfileManager{Config: config}, nil
}

// updateConfig writes the current configuration to disk
func (pm *ProfileManager) updateConfig() error {
	data, err := json.MarshalIndent(pm.Config, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigPath(), data, 0644)
}

// CreateProfile creates a new profile
func (pm *ProfileManager) CreateProfile(name string, basedOff string) error {
	var contents []byte
	var err error

	if basedOff != "" {
		contents, err = os.ReadFile(getStateFilePath(basedOff))
		if err != nil {
			return err
		}
	} else {
		contents = []byte("{}")
	}

	if err := os.WriteFile(getStateFilePath(name), contents, 0644); err != nil {
		return err
	}

	pm.Config.Profiles = append(pm.Config.Profiles, name)
	return pm.updateConfig()
}

// SwitchProfile switches to a different profile
func (pm *ProfileManager) SwitchProfile(profileName string) error {
	// Backup current profile
	if err := copyFile(getLauncherStateFilePath(), getStateFilePath(pm.Config.ActiveProfile)); err != nil {
		return err
	}

	// Copy new state over
	if err := copyFile(getStateFilePath(profileName), getLauncherStateFilePath()); err != nil {
		return err
	}

	pm.Config.ActiveProfile = profileName
	return pm.updateConfig()
}

// copyFile is a helper function to copy files
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}
