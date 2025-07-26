import { create } from 'zustand';
import { GetSettings, SetSetting } from 'wjs/go/app/App';
import { SETTING_DEFAULTS, SettingKey, Settings, SettingValues } from '~/utils/settings';

interface SettingsState {
  settings: Settings | null;
  fetchSettings: () => Promise<void>;
  setSetting: <T extends SettingKey>(settingKey: T, value: SettingValues[T]) => Promise<void>;
  getSetting: <T extends SettingKey>(settingKey: T) => SettingValues[T];
}

export const useSettingsStore = create<SettingsState>((set, get) => ({
  settings: null,
  fetchSettings: async () => {
    const settings = await GetSettings();
    set({ settings: settings as Settings });
  },
  setSetting: async (settingKey, value) => {
    await SetSetting(settingKey, value).catch((error) => {
      console.error(`Failed to set setting ${settingKey}:`, error);
    });

    await get().fetchSettings();
  },
  getSetting: (settingKey) => {
    const settings = get().settings;

    const setting = settings ? settings[settingKey] || null : null;

    return setting || SETTING_DEFAULTS[settingKey];
  },
}));
