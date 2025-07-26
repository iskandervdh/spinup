import {
  ArrowPathIcon,
  PlusIcon,
  InformationCircleIcon,
  ListBulletIcon,
  Squares2X2Icon,
} from '@heroicons/react/20/solid';
import { GetProjects } from 'wjs/go/app/App';
import { useCallback, useEffect } from 'react';
import { PageTitle } from '~/components/page-title';
import { useProjectsStore } from '~/stores/projectsStore';
import { Button } from '~/components/button';
import { LogsPopover } from '~/sections/logs-popover';
import { createFileRoute, Link } from '@tanstack/react-router';
import { useSettingsStore } from '~/stores/settingsStore';
import { SettingKey } from '~/utils/settings';
import { cn } from '~/utils/helpers';
import { ProjectInfo } from '~/sections/project-info';

export const Route = createFileRoute('/projects')({
  component: Projects,
});

export function Projects() {
  const { projects, setProjects, setEditingProject } = useProjectsStore();

  const { setSetting } = useSettingsStore();

  const projectViewLayout = useSettingsStore((state) => state.getSetting(SettingKey.ProjectViewLayout) || 'grid');

  const fetchProjects = useCallback(() => {
    GetProjects().then((projects) => setProjects(projects || []));
  }, [setProjects]);

  useEffect(() => {
    fetchProjects();
  }, []);

  return (
    <div id="projects" className="max-w-6xl">
      <LogsPopover />

      <div className="flex items-center gap-4 pb-4">
        <div className="flex items-center flex-1 gap-4">
          <PageTitle>Projects</PageTitle>

          <div className="flex gap-2">
            <Button onClick={fetchProjects} size="icon-lg" title="Refresh projects">
              <ArrowPathIcon width={24} height={24} className="text-current" />
            </Button>

            <Button size="icon-lg" variant="success" title="Add project" asChild={true}>
              <Link to="/project-form" onClick={() => setEditingProject(null)}>
                <PlusIcon width={24} height={24} className="text-current" />
              </Link>
            </Button>
          </div>
        </div>
        <div className="flex items-center gap-1">
          <Button
            size="icon-lg"
            variant="ghost"
            onClick={() => setSetting(SettingKey.ProjectViewLayout, 'list')}
            className={cn('hover:bg-black/10', {
              'bg-black/10': projectViewLayout === 'list',
            })}
            title="Change view layout to list"
          >
            <ListBulletIcon width={24} height={24} className="text-primary" />
          </Button>
          <Button
            size="icon-lg"
            variant="ghost"
            onClick={() => setSetting(SettingKey.ProjectViewLayout, 'grid')}
            className={cn('hover:bg-black/10', {
              'bg-black/10': projectViewLayout === 'grid',
            })}
            title="Change view layout to grid"
          >
            <Squares2X2Icon width={24} height={24} className="text-primary" />
          </Button>
        </div>
      </div>

      <div
        className={cn({
          'flex flex-col gap-4': projectViewLayout === 'list',
          'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4': projectViewLayout === 'grid',
        })}
      >
        {projects ? (
          projects.length === 0 ? (
            <div className="flex flex-col gap-4 py-2">
              <div className="flex items-center gap-2 text-lg text-gray-300">
                <InformationCircleIcon width={24} height={24} className="text-info" />
                <span>No projects found.</span>
              </div>

              <div>
                <Link to="/project-form" onClick={() => setEditingProject(null)}>
                  <Button size="xs" variant="success">
                    Add a project
                  </Button>
                </Link>
              </div>
            </div>
          ) : (
            projects.map((project) => <ProjectInfo key={project.ID} project={project} />)
          )
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </div>
  );
}
