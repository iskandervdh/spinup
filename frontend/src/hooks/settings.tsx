import { useSettingsStore } from '~/stores/settingsStore';
import { SETTING_DEFAULTS, SettingKey } from '~/utils/settings';

export function useShowCommandIcons() {
  return useSettingsStore((state) =>
    state.settings ? state.settings[SettingKey.ShowCommandIcons] : SETTING_DEFAULTS[SettingKey.ShowCommandIcons]
  );
}
