import {
  DocumentTextIcon,
  ExclamationTriangleIcon,
  FolderIcon,
  PencilSquareIcon,
  PlayIcon,
  StopIcon,
  TrashIcon,
} from '@heroicons/react/20/solid';
import { useNavigate } from '@tanstack/react-router';
import { useCallback, useMemo } from 'react';
import toast from 'react-hot-toast';
import { core } from 'wjs/go/models';
import { BrowserOpenURL } from 'wjs/runtime/runtime';
import { Button } from '~/components/button';
import { useProjectsStore } from '~/stores/projectsStore';
import { useSettingsStore } from '~/stores/settingsStore';
import { cn } from '~/utils/helpers';
import { SettingKey } from '~/utils/settings';

export function ProjectInfoHeader({ project, isRunning }: { project: core.Project; isRunning: boolean }) {
  const navigate = useNavigate();

  const { runProject, stopProject, removeProject, setCurrentProject, setEditingProject } = useProjectsStore();

  const projectViewLayout = useSettingsStore((state) => state.getSetting(SettingKey.ProjectViewLayout) || 'grid');

  const canRunProject = useMemo(
    () => project.Dir.Valid && project.Commands !== null && project.Commands.length > 0,
    [project.Dir, project.Commands]
  );
  const cannotRunProjectReason = useMemo(() => {
    if (!project.Dir.Valid) return 'Project directory is not set';
    if (project.Commands === null || project.Commands.length === 0) return 'No commands set for this project';
    return '';
  }, [project.Dir, project.Commands]);

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
    setEditingProject(project.ID);
    navigate({ to: '/project-form' });
  }, [project.Name, setEditingProject]);

  const remove = useCallback(async () => {
    if (confirm(`Are you sure you want to remove project "${project.Name}"?`)) {
      await removeProject(project.ID);

      toast.success(<b>Removed project "{project.Name}"</b>);
    }
  }, [project.Name, removeProject]);

  if (projectViewLayout === 'grid') {
    return (
      <div className="flex gap-2 mb-2">
        {canRunProject ? (
          <button
            className={cn(
              'p-2 rounded-lg hover:bg-black/10 focus:outline-offset-2 focus-visible:outline focus:outline-1',
              isRunning ? 'focus:outline-error' : 'focus:outline-success'
            )}
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
        </div>

        <div className="flex justify-end flex-1 gap-2">
          {isRunning ? (
            <Button onClick={showLogs} size="square" variant="info" title="Show logs">
              <DocumentTextIcon width={16} height={16} className="text-current" />
            </Button>
          ) : (
            <>
              <Button onClick={edit} size="square" title="Edit project">
                <PencilSquareIcon width={16} height={16} className="text-current" />
              </Button>
              <Button onClick={remove} size="square" variant="error" title="Remove project">
                <TrashIcon width={16} height={16} className="text-current" />
              </Button>
            </>
          )}
        </div>
      </div>
    );
  }

  return (
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

        {isRunning ? (
          <Button onClick={showLogs} size="icon" variant="info" title="Show logs">
            <DocumentTextIcon width={16} height={16} className="text-current" />
          </Button>
        ) : (
          <>
            <Button onClick={edit} size="xs" title="Edit project">
              <PencilSquareIcon width={16} height={16} className="text-current" />
            </Button>
            <Button onClick={remove} size="icon" variant="error" title="Remove project">
              <TrashIcon width={16} height={16} className="text-current" />
            </Button>
          </>
        )}
      </div>
    </div>
  );
}

export function ProjectInfo({ project }: { project: core.Project }) {
  const { runningProjects, updateProjectDir } = useProjectsStore();

  const projectViewLayout = useSettingsStore((state) => state.getSetting(SettingKey.ProjectViewLayout) || 'grid');

  const isRunning = useMemo(() => runningProjects.includes(project.Name), [runningProjects]);
  const commands = useMemo(() => project.Commands?.map((c) => c.Name).join(', '), [project.Commands]);
  // const variables = useMemo(
  //   () => project.Variables?.map((v) => `${v.Name}=${v.Value}`).join(', '),
  //   [project.Variables]
  // );
  const domainAliases = useMemo(() => project.DomainAliases?.map((da) => da.Value).join(', '), [project.DomainAliases]);

  const projectDomain = useMemo(() => `${project.Name}.test`, [project.Name]);

  const openProjectInBrowser = useCallback(() => {
    BrowserOpenURL(`http://${projectDomain}`);
  }, [projectDomain]);

  const openProjectDir = useCallback(() => project.Dir.Valid && BrowserOpenURL(project.Dir.String), [project.Dir]);

  const openUpdateProjectDir = useCallback(
    () => updateProjectDir(project.Name, project.Dir.String),
    [project.Name, project.Dir.String]
  );

  if (projectViewLayout === 'grid') {
    return (
      <div className="flex-1 p-2 border rounded-lg border-primary/80 min-w-72">
        <ProjectInfoHeader project={project} isRunning={isRunning} />

        <div className="grid items-center max-w-6xl grid-cols-1 gap-x-4">
          <div className="block">
            {isRunning ? (
              <span onClick={openProjectInBrowser} className="underline cursor-pointer text-info hover:text-info-dark">
                {projectDomain}
              </span>
            ) : (
              <span>{projectDomain}</span>
            )}{' '}
            {project.DomainAliases ? (
              <span className="px-2 py-1 text-sm rounded-lg select-none bg-black/10">
                +{project.DomainAliases.length}
              </span>
            ) : (
              ''
            )}
          </div>

          <div>
            Port: <span className="text-sm">{project.Port}</span>
          </div>

          <div className="flex items-center justify-between gap-2">
            <div>
              Commands:{' '}
              <span className="text-sm">
                {commands !== undefined && commands !== '' ? commands : <span className="text-sm text-error">-</span>}
              </span>
            </div>

            {project.Dir.Valid ? (
              <div className="flex items-center gap-2">
                <Button onClick={openProjectDir} size="square" title="Open project directory">
                  <FolderIcon width={16} height={16} className="text-current" />
                </Button>
              </div>
            ) : (
              <Button onClick={openUpdateProjectDir} size="xs" title="Select project directory">
                Select directory
              </Button>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div>
      <ProjectInfoHeader project={project} isRunning={isRunning} />

      <div className="grid items-center max-w-6xl grid-cols-[16rem,auto] gap-x-4">
        <div>Domain</div>

        {isRunning ? (
          <div
            onClick={openProjectInBrowser}
            className="text-sm underline cursor-pointer text-info hover:text-info-dark"
          >
            {projectDomain}
          </div>
        ) : (
          <div className="text-sm">{projectDomain}</div>
        )}

        <div>Port</div>
        <div className="text-sm">{project.Port}</div>

        <div>Commands</div>
        {commands !== undefined && commands !== '' ? (
          <div className="text-sm">{commands}</div>
        ) : (
          <div className="text-sm text-error">-</div>
        )}

        <div>Directory</div>
        {project.Dir.Valid ? (
          <div className="flex items-center gap-2 py-1">
            <Button onClick={openProjectDir} size="xs" title="Open project directory">
              <FolderIcon width={16} height={16} className="text-current" />
            </Button>
          </div>
        ) : (
          <div className="w-full py-2 min-w-32 max-w-64">
            <Button onClick={openUpdateProjectDir} size="xs" title="Select project directory">
              Select directory
            </Button>
          </div>
        )}

        {/* <div>Variables</div>
        <div className="text-sm">{variables || '-'}</div> */}

        <div>Domain aliases</div>
        <div className="text-sm">{domainAliases || '-'}</div>
      </div>
    </div>
  );
}
