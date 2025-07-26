export enum SettingKey {
  ProjectViewLayout = 'projectViewLayout',
}

export type SettingValues = {
  [SettingKey.ProjectViewLayout]: 'grid' | 'list';
};
