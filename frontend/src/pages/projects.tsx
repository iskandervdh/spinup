import {
  PlayIcon,
  StopIcon,
  ExclamationTriangleIcon,
  ArrowPathIcon,
  PlusIcon,
  TrashIcon,
  DocumentTextIcon,
  PencilSquareIcon,
  FolderIcon,
  InformationCircleIcon,
} from '@heroicons/react/20/solid';
import { GetProjects } from 'wjs/go/app/App';
import { useCallback, useEffect, useMemo } from 'react';
import { PageTitle } from '~/components/page-title';
import { useProjectsStore } from '~/stores/projectsStore';
import { BrowserOpenURL } from 'wjs/runtime/runtime';
import { Button } from '~/components/button';
import { Page, usePageStore } from '~/stores/pageStore';
import { LogsPopover } from '~/sections/logs-popover';
import { core } from 'wjs/go/models';
import toast from 'react-hot-toast';

function ProjectInfo({ project }: { project: core.Project }) {
  const {
    runningProjects,
    runProject,
    stopProject,
    selectProjectDir,
    removeProject,
    setCurrentProject,
    setEditingProject,
  } = useProjectsStore();
  const { setCurrentPage } = usePageStore();

  const isRunning = useMemo(() => runningProjects.includes(project.Name), [runningProjects]);
  const commands = useMemo(() => project.Commands.map((c) => c.Name).join(', '), [project.Commands]);
  const variables = useMemo(
    () => project.Variables?.map((v) => `${v.Name}=${v.Value}`).join(', '),
    [project.Variables]
  );
  const domainAliases = useMemo(() => project.DomainAliases?.map((da) => da.Value).join(', '), [project.DomainAliases]);

  const canRunProject = useMemo(
    () => project.Dir.Valid && project.Commands.length > 0,
    [project.Dir, project.Commands]
  );
  const cannotRunProjectReason = useMemo(() => {
    if (!project.Dir.Valid) return 'Project directory is not set';
    if (project.Commands.length === 0) return 'No commands set for this project';
    return '';
  }, [project.Dir, project.Commands]);

  const openProjectInBrowser = useCallback(() => {
    BrowserOpenURL(`http://${project.Domain}`);
  }, [project.Domain]);

  const openProjectDir = useCallback(() => project.Dir.Valid && BrowserOpenURL(project.Dir.String), [project.Dir]);

  const openSelectProjectDir = useCallback(
    () => selectProjectDir(project.Name, project.Dir.String),
    [project.Name, project.Dir.String]
  );

  const startOrStopProject = useCallback(async () => {
    if (isRunning) {
      await stopProject(project.Name);
    } else {
      await runProject(project.Name);
    }
  }, [isRunning]);

  const showLogs = useCallback(() => {
    if (isRunning) {
      setCurrentProject(project.Name);
    } else {
      alert('Project is not running');
    }
  }, [project.Name, isRunning]);

  const edit = useCallback(() => {
    setEditingProject(project.Name);
    setCurrentPage(Page.ProjectForm);
  }, [project.Name, setEditingProject, setCurrentPage]);

  const remove = useCallback(async () => {
    if (confirm(`Are you sure you want to remove project "${project.Name}"?`)) {
      await removeProject(project.Name);

      toast.success(<b>Removed project "{project.Name}"</b>);
    }
  }, [project.Name, removeProject]);

  return (
    <div>
      <div className="flex items-center gap-2 mb-2">
        {canRunProject ? (
          <button
            className="p-2 rounded-lg hover:bg-black/10"
            onClick={startOrStopProject}
            title={isRunning ? 'Stop project' : 'Start project'}
          >
            {isRunning ? (
              <StopIcon width={20} height={20} className="text-red-400" />
            ) : (
              <PlayIcon width={20} height={20} className="text-green-400" />
            )}
          </button>
        ) : (
          <div className="p-2 rounded-lg hover:bg-black/10" title={cannotRunProjectReason}>
            <ExclamationTriangleIcon width={20} height={20} className="text-yellow-400" />
          </div>
        )}

        <div className="flex items-center gap-2">
          <h3 className="pr-2 text-xl font-bold text-primary">{project.Name}</h3>

          <Button onClick={edit} size={'xs'} title="Edit project">
            <PencilSquareIcon width={16} height={16} className="text-current" />
          </Button>

          {isRunning ? (
            <Button onClick={showLogs} size={'icon'} variant={'info'} title="Show logs">
              <DocumentTextIcon width={16} height={16} className="text-current" />
            </Button>
          ) : (
            <Button onClick={remove} size={'icon'} variant={'error'} title="Remove project">
              <TrashIcon width={16} height={16} className="text-current" />
            </Button>
          )}
        </div>
      </div>

      <div className="grid items-center max-w-6xl grid-cols-[16rem,auto] gap-x-4">
        <div>Domain</div>

        {isRunning ? (
          <div
            onClick={openProjectInBrowser}
            className="text-sm underline cursor-pointer text-info hover:text-info-dark"
          >
            {project.Domain}
          </div>
        ) : (
          <div className="text-sm">{project.Domain}</div>
        )}

        <div>Port</div>
        <div className="text-sm">{project.Port}</div>

        <div>Commands</div>
        {commands !== '' ? <div className="text-sm">{commands}</div> : <div className="text-sm text-error">-</div>}

        <div>Directory</div>
        {project.Dir.Valid ? (
          <div className="flex items-center gap-2 py-1">
            <Button onClick={openProjectDir} size={'xs'} title="Open project directory">
              <FolderIcon width={16} height={16} className="text-current" />
            </Button>

            <Button onClick={openSelectProjectDir} size={'xs'} title="Change project directory">
              <PencilSquareIcon width={16} height={16} className="text-current" />
            </Button>
          </div>
        ) : (
          <div className="w-full py-2 min-w-32 max-w-64">
            <Button onClick={openSelectProjectDir} size={'xs'} title="Select project directory">
              Select directory
            </Button>
          </div>
        )}

        <div>Variables</div>
        <div className="text-sm">{variables || '-'}</div>

        <div>Domain aliases</div>
        <div className="text-sm">{domainAliases || '-'}</div>
      </div>
    </div>
  );
}

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
