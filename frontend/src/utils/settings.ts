export enum SettingKey {
  ProjectViewLayout = 'projectViewLayout',
  ShowCommandIcons = 'showCommandIcons',
}

export type SettingValues = {
  [SettingKey.ProjectViewLayout]: 'grid' | 'list';
  [SettingKey.ShowCommandIcons]: boolean;
};

export type Settings = { [K in SettingKey]: SettingValues[K] };

export const SETTING_DEFAULTS = {
  [SettingKey.ProjectViewLayout]: 'grid',
  [SettingKey.ShowCommandIcons]: true,
} as const satisfies Settings;
