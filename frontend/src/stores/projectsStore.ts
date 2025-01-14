import { create } from 'zustand';
import { Projects } from '~/types';
import {
  GetProjects,
  RunProject,
  SelectProjectDirectory,
  StopProject,
  AddProject,
  RemoveProject,
  UpdateProject,
} from 'wjs/go/app/App';

interface ProjectsState {
  projects: Projects | null;
  setProjects: (projects: Projects) => void;

  editingProject: string | null;
  setEditingProject: (projectName: string | null) => void;

  runningProjects: string[];
  runProject: (projectName: string) => Promise<void>;
  stopProject: (projectName: string) => Promise<void>;
  selectProjectDir: (projectName: string, defaultDir: string | undefined) => Promise<void>;
  projectFormSubmit: (projectName: string, port: number, commandNames: string[]) => Promise<void>;
  removeProject: (projectName: string) => Promise<void>;

  currentProject: string | null;
  setCurrentProject: (projectName: string | null) => void;
}

export const useProjectsStore = create<ProjectsState>((set, get) => ({
  projects: null,
  setProjects: (projects) => set(() => ({ projects })),

  editingProject: null,
  setEditingProject: (projectName) => set(() => ({ editingProject: projectName })),

  runningProjects: [],
  async runProject(projectName) {
    set((state) => ({ runningProjects: [...state.runningProjects, projectName] }));

    await RunProject(projectName);
  },
  async stopProject(projectName) {
    set((state) => ({ runningProjects: state.runningProjects.filter((p) => p !== projectName) }));

    await StopProject(projectName);
  },
  async selectProjectDir(projectName, defaultDir) {
    await SelectProjectDirectory(projectName, defaultDir ?? '');

    const projects = await GetProjects();
    set(() => ({ projects }));
  },
  async projectFormSubmit(projectName, port, commandNames) {
    if (get().editingProject === projectName) {
      await UpdateProject(projectName, port, commandNames);
      set(() => ({ editingProject: null }));
    } else {
      await AddProject(projectName, port, commandNames);
    }

    const projects = await GetProjects();
    set(() => ({ projects }));
  },
  async removeProject(projectName) {
    await RemoveProject(projectName);

    const projects = await GetProjects();
    set(() => ({ projects }));
  },

  currentProject: null,
  setCurrentProject: (projectName) => set(() => ({ currentProject: projectName })),
}));
