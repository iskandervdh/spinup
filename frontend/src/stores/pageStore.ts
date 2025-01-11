import { useMemo } from 'react';
import { create } from 'zustand';
import { CommandFormPage } from '~/pages/command-form';
import { ProjectFormPage } from '~/pages/project-form';
import { CommandsPage } from '~/pages/commands';
import { ProjectsPage } from '~/pages/projects';
import { SettingsPage } from '~/pages/settings';

export enum Page {
  Projects = 'Projects',
  Commands = 'Commands',
  Settings = 'Settings',
  ProjectForm = 'ProjectForm',
  CommandForm = 'CommandForm',
}

export const DEFAULT_PAGE = Page.Projects;

export const NAV_PAGES = [Page.Projects, Page.Commands, Page.Settings] as const;

interface PageState {
  currentPage: Page;
  setCurrentPage: (page: Page) => void;
  getPageComponent: (currentPage: Page) => () => React.JSX.Element;
}

export const usePageStore = create<PageState>((set) => ({
  currentPage: DEFAULT_PAGE,
  setCurrentPage: (page) => set(() => ({ currentPage: page })),
  getPageComponent: (currentPage) => {
    switch (currentPage) {
      case Page.Projects:
        return ProjectsPage;
      case Page.Commands:
        return CommandsPage;
      case Page.Settings:
        return SettingsPage;
      case Page.ProjectForm:
        return ProjectFormPage;
      case Page.CommandForm:
        return CommandFormPage;
    }
  },
}));
