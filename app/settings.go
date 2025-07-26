package app

func (a *App) GetSettings() (map[string]interface{}, error) {
	return a.core.GetConfig().GetSettings()
}

func (a *App) GetSetting(settingKey string) (interface{}, error) {
	return a.core.GetConfig().GetSetting(settingKey)
}

func (a *App) SetSetting(settingKey string, value interface{}) error {
	return a.core.GetConfig().SetSetting(settingKey, value)
}
