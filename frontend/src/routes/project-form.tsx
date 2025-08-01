import { useCallback, useEffect, useMemo, useState } from 'react';
import { Input } from '~/components/input';
import { PageTitle } from '~/components/page-title';
import { useCommandsStore } from '~/stores/commandsStore';
import { useProjectsStore } from '~/stores/projectsStore';
import { GetCommands } from 'wjs/go/app/App';
import { Button } from '~/components/button';
import { SelectMultiple } from '~/components/select-multiple';
import toast from 'react-hot-toast';
import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { PencilSquareIcon } from '@heroicons/react/20/solid';
import { getCommandIcon } from '~/utils/command';
import { useShowCommandIcons } from '~/hooks/settings';

export const Route = createFileRoute('/project-form')({
  component: ProjectFormPage,
});

export function ProjectFormPage() {
  const navigate = useNavigate();

  const { commands, setCommands } = useCommandsStore();
  const { projects, projectFormSubmit, editingProject, selectProjectDir } = useProjectsStore();

  const showCommandIcons = useShowCommandIcons();

  const [name, setName] = useState('');
  const [port, setPort] = useState(3000);
  const [commandNames, setCommandNames] = useState<string[]>([]);
  const [projectDir, setProjectDir] = useState<string | null>(null);

  const pageTitle = useMemo(
    () =>
      editingProject ? `Edit Project "${projects?.find((p) => p.ID === editingProject)?.Name || ''}"` : 'Add Project',
    [projects, editingProject]
  );

  const submitText = useMemo(() => (editingProject ? 'Save Project' : 'Add Project'), [editingProject]);

  const submit = useCallback(
    async (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();

      await toast
        .promise(projectFormSubmit(name, port, commandNames, projectDir), {
          loading: editingProject ? 'Saving project...' : 'Creating project...',
          success: editingProject ? <b>Project saved</b> : <b>Project created</b>,
          error: (err: Error) =>
            editingProject ? (
              <b>
                Failed to save project:
                <br />
                {err.message}
              </b>
            ) : (
              <b>
                Failed to create project:
                <br />
                {err.message}
              </b>
            ),
        })
        .then(() => {
          navigate({ to: '/projects' });
        });
    },
    [name, port, commandNames, projectDir, editingProject, projectFormSubmit]
  );

  const openSelectProjectDir = useCallback(() => {
    selectProjectDir(name, projectDir).then(setProjectDir);
  }, [name, projectDir]);

  useEffect(() => {
    GetCommands().then(setCommands);
  }, []);

  useEffect(() => {
    if (editingProject) {
      const project = projects?.find((p) => p.ID === editingProject);

      if (project) {
        setName(project.Name);
        setPort(project.Port);
        setCommandNames(
          project.Commands !== null ? project.Commands.map((c) => c.Name).sort((a, b) => a.localeCompare(b)) : []
        );
        setProjectDir(project.Dir.Valid ? project.Dir.String : null);
      }
    }
  }, [editingProject, setName, setPort, setCommandNames]);

  return (
    <form id="project-form" onSubmit={submit} className="flex flex-col w-full max-w-6xl mx-auto">
      <div className="flex items-center pb-4 h-14">
        <PageTitle>{pageTitle}</PageTitle>
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <label htmlFor="name" className="w-min">
            Name
          </label>
          <Input id="name" name="name" type="text" required value={name} onChange={(e) => setName(e.target.value)} />
        </div>

        <div className="flex flex-col gap-2">
          <label htmlFor="port" className="w-min">
            Port
          </label>
          <Input
            id="port"
            name="port"
            type="number"
            required
            min={1}
            max={65536}
            value={port}
            onChange={(e) => setPort(parseInt(e.target.value))}
          />
        </div>

        <div className="flex flex-col gap-2">
          <label className="w-min">Commands</label>
          <SelectMultiple
            id="commands"
            name="commands"
            options={
              commands
                ? commands.map((c) => ({
                    label: `${c.Name}: ${c.Command}`,
                    value: c.Name,
                    icon: showCommandIcons ? getCommandIcon(c.Command) : null,
                  }))
                : []
            }
            value={commandNames}
            onChanged={setCommandNames}
          />
        </div>

        <div className="flex flex-col gap-2">
          <label className="w-min">Directory</label>

          <div className="flex items-center gap-4">
            <div className="text-sm">{projectDir ?? 'Not set'}</div>

            <Button type="button" onClick={openSelectProjectDir} size="xs" title="Change project directory">
              <PencilSquareIcon width={16} height={16} className="text-current" />
            </Button>
          </div>
        </div>

        <Button type="submit" className="mt-2">
          {submitText}
        </Button>
      </div>
    </form>
  );
}
