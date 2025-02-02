import { create } from 'zustand';
import { Projects } from '~/types';
import {
  GetProjects,
  RunProject,
  UpdateProjectDirectory,
  StopProject,
  AddProject,
  RemoveProject,
  UpdateProject,
  SelectProjectDirectory,
} from 'wjs/go/app/App';

interface ProjectsState {
  projects: Projects | null;
  setProjects: (projects: Projects) => void;

  editingProject: number | null;
  setEditingProject: (projectName: number | null) => void;

  runningProjects: string[];
  runProject: (projectName: string) => Promise<void>;
  stopProject: (projectName: string) => Promise<void>;
  updateProjectDir: (projectName: string, defaultDir: string | undefined) => Promise<void>;
  selectProjectDir: (projectName: string, defaultDir: string | null) => Promise<string>;
  projectFormSubmit: (
    projectName: string,
    port: number,
    commandNames: string[],
    projectDir: string | null
  ) => Promise<void>;
  removeProject: (projectID: number) => Promise<void>;

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
  async updateProjectDir(projectName, defaultDir) {
    await UpdateProjectDirectory(projectName, defaultDir ?? '');

    const projects = await GetProjects();
    set(() => ({ projects }));
  },
  async selectProjectDir(projectName, defaultDir) {
    return await SelectProjectDirectory(projectName, defaultDir ?? '');
  },
  async projectFormSubmit(projectName, port, commandNames, projectDir) {
    if (projectName.includes(' ')) {
      throw new Error('Project name can not include a space');
    }

    const projectID = get().editingProject;

    if (projectID !== null) {
      await UpdateProject(projectID, projectName, port, commandNames, projectDir ?? '');
      set(() => ({ editingProject: null }));
    } else {
      await AddProject(projectName, port, commandNames, projectDir ?? '');
    }

    const projects = await GetProjects();
    set(() => ({ projects }));
  },
  async removeProject(projectID) {
    await RemoveProject(projectID);

    const projects = await GetProjects();
    set(() => ({ projects }));
  },

  currentProject: null,
  setCurrentProject: (projectName) => set(() => ({ currentProject: projectName })),
}));
