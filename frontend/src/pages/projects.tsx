import { ArrowPathIcon, PlusIcon, InformationCircleIcon } from '@heroicons/react/20/solid';
import { GetProjects } from 'wjs/go/app/App';
import { useEffect } from 'react';
import { PageTitle } from '~/components/page-title';
import { useProjectsStore } from '~/stores/projectsStore';
import { Button } from '~/components/button';
import { Page, usePageStore } from '~/stores/pageStore';
import { LogsPopover } from '~/sections/logs-popover';
import { ProjectInfo } from '~/sections/project-info';

export function ProjectsPage() {
  const { setCurrentPage } = usePageStore();
  const { projects, setProjects, setEditingProject } = useProjectsStore();

  useEffect(() => {
    GetProjects().then((projects) => setProjects(projects || []));
  }, []);

  return (
    <div id="projects">
      <div className="flex items-center gap-4 pb-4">
        <PageTitle>Projects</PageTitle>

        <LogsPopover />

        <div className="flex gap-2">
          <Button onClick={() => GetProjects().then(setProjects)} size={'icon-lg'} title="Refresh projects">
            <ArrowPathIcon width={24} height={24} className="text-current" />
          </Button>

          <Button
            onClick={() => {
              setEditingProject(null);
              setCurrentPage(Page.ProjectForm);
            }}
            size={'icon-lg'}
            variant={'success'}
            title="Add project"
          >
            <PlusIcon width={24} height={24} className="text-current" />
          </Button>
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
                <Button
                  onClick={() => {
                    setEditingProject(null);
                    setCurrentPage(Page.ProjectForm);
                  }}
                  size={'xs'}
                  variant={'success'}
                >
                  Add a project
                </Button>
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
