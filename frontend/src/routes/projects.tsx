import { ArrowPathIcon, PlusIcon, InformationCircleIcon } from '@heroicons/react/20/solid';
import { GetProjects } from 'wjs/go/app/App';
import { useCallback, useEffect } from 'react';
import { PageTitle } from '~/components/page-title';
import { useProjectsStore } from '~/stores/projectsStore';
import { Button } from '~/components/button';
import { LogsPopover } from '~/sections/logs-popover';
import { ProjectInfo } from '~/sections/project-info';
import { createFileRoute, Link } from '@tanstack/react-router';

export const Route = createFileRoute('/projects')({
  component: Projects,
});

export function Projects() {
  const { projects, setProjects, setEditingProject } = useProjectsStore();

  const fetchProjects = useCallback(() => {
    GetProjects().then((projects) => setProjects(projects || []));
  }, [setProjects]);

  useEffect(() => {
    fetchProjects();
  }, []);

  return (
    <div id="projects">
      <div className="flex items-center gap-4 pb-4">
        <PageTitle>Projects</PageTitle>

        <LogsPopover />

        <div className="flex gap-2">
          <Button onClick={fetchProjects} size="icon-lg" title="Refresh projects">
            <ArrowPathIcon width={24} height={24} className="text-current" />
          </Button>

          <Link to="/project-form" onClick={() => setEditingProject(null)}>
            <Button size="icon-lg" variant="success" title="Add project">
              <PlusIcon width={24} height={24} className="text-current" />
            </Button>
          </Link>
        </div>
      </div>

      <div className="flex flex-col gap-4">
        {projects ? (
          projects.length === 0 ? (
            <div className="flex flex-col gap-4 py-2">
              <div className="flex items-center gap-2 text-lg text-gray-300">
                <InformationCircleIcon width={24} height={24} className="text-info" />
                <span>No projects found.</span>
              </div>

              <div>
                <Link to="/project-form" onClick={() => setEditingProject(null)}>
                  <Button size="xs" variant={'success'}>
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
