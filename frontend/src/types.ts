import { GetCommands, GetProjects } from 'wjs/go/app/App';

export type Projects = Awaited<ReturnType<typeof GetProjects>>;

export type Commands = Awaited<ReturnType<typeof GetCommands>>;
