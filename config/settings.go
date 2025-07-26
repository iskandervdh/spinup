package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const settingsFileName = "settings.json"

func (c *Config) writeSettings(settings map[string]interface{}) error {
	settingFilePath := path.Join(c.configDir, settingsFileName)

	data, err := json.MarshalIndent(settings, "", "  ")

	if err != nil {
		return fmt.Errorf("could not marshal settings: %w", err)
	}

	err = os.WriteFile(settingFilePath, data, 0644)

	if err != nil {
		return fmt.Errorf("could not write settings file: %w", err)
	}

	return nil
}

func (c *Config) GetSettings() (map[string]interface{}, error) {
	settingFilePath := path.Join(c.configDir, settingsFileName)
	settings := make(map[string]interface{})

	data, err := os.ReadFile(settingFilePath)

	if os.IsNotExist(err) {
		// If the settings file does not exist we should create it.
		err = c.writeSettings(settings)

		if err != nil {
			return nil, fmt.Errorf("could not create settings file: %w", err)
		}

		// Initialize settings data with an empty JSON object
		data = []byte("{}")
	}

	if err != nil {
		return nil, fmt.Errorf("could not read settings file: %w", err)
	}

	err = json.Unmarshal(data, &settings)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal settings: %w", err)
	}

	return settings, nil
}

func (c *Config) GetSetting(settingKey string) (interface{}, error) {
	settings, err := c.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("could not read settings: %w", err)
	}

	value, exists := settings[settingKey]
	if !exists {
		return nil, fmt.Errorf("setting %s does not exist", settingKey)
	}

	return value, nil
}

func (c *Config) SetSetting(settingKey string, value interface{}) error {
	settings, err := c.GetSettings()

	if err != nil {
		return fmt.Errorf("could not read settings: %w", err)
	}

	settings[settingKey] = value
	err = c.writeSettings(settings)

	if err != nil {
		return fmt.Errorf("could not write settings: %w", err)
	}
	return nil
}
