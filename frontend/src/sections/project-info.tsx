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

export function ProjectInfo({ project }: { project: core.Project }) {
  const navigate = useNavigate();

  const {
    runningProjects,
    runProject,
    stopProject,
    updateProjectDir,
    removeProject,
    setCurrentProject,
    setEditingProject,
  } = useProjectsStore();

  const isRunning = useMemo(() => runningProjects.includes(project.Name), [runningProjects]);
  const commands = useMemo(() => project.Commands?.map((c) => c.Name).join(', '), [project.Commands]);
  // const variables = useMemo(
  //   () => project.Variables?.map((v) => `${v.Name}=${v.Value}`).join(', '),
  //   [project.Variables]
  // );
  const domainAliases = useMemo(() => project.DomainAliases?.map((da) => da.Value).join(', '), [project.DomainAliases]);

  const canRunProject = useMemo(
    () => project.Dir.Valid && project.Commands !== null && project.Commands.length > 0,
    [project.Dir, project.Commands]
  );
  const cannotRunProjectReason = useMemo(() => {
    if (!project.Dir.Valid) return 'Project directory is not set';
    if (project.Commands === null || project.Commands.length === 0) return 'No commands set for this project';
    return '';
  }, [project.Dir, project.Commands]);

  const projectDomain = useMemo(() => `${project.Name}.test`, [project.Name]);

  const openProjectInBrowser = useCallback(() => {
    BrowserOpenURL(`http://${projectDomain}`);
  }, [projectDomain]);

  const openProjectDir = useCallback(() => project.Dir.Valid && BrowserOpenURL(project.Dir.String), [project.Dir]);

  const openUpdateProjectDir = useCallback(
    () => updateProjectDir(project.Name, project.Dir.String),
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
    setEditingProject(project.ID);
    navigate({ to: '/project-form' });
  }, [project.Name, setEditingProject]);

  const remove = useCallback(async () => {
    if (confirm(`Are you sure you want to remove project "${project.Name}"?`)) {
      await removeProject(project.ID);

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
