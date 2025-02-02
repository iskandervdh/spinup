import { useCallback, useEffect, useMemo, useState } from 'react';
import { Input } from '~/components/input';
import { PageTitle } from '~/components/page-title';
import { useCommandsStore } from '~/stores/commandsStore';
import { Button } from '~/components/button';
import toast from 'react-hot-toast';
import { createFileRoute, useNavigate } from '@tanstack/react-router';

export const Route = createFileRoute('/command-form')({
  component: CommandForm,
});

function CommandForm() {
  const navigate = useNavigate();

  const { commands, commandFormSubmit, editingCommand } = useCommandsStore();

  const [name, setName] = useState('');
  const [command, setCommand] = useState('');

  const pageTitle = useMemo(
    () =>
      editingCommand ? `Edit Command "${commands?.find((c) => c.ID === editingCommand)?.Name || ''}"` : 'Add Command',
    [editingCommand]
  );

  const submitText = useMemo(() => (editingCommand ? 'Save Command' : 'Add Command'), [editingCommand]);

  const submit = useCallback(
    async (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();

      await toast
        .promise(commandFormSubmit(name, command), {
          loading: editingCommand ? 'Saving command...' : 'Creating command...',
          success: editingCommand ? <b>Command saved</b> : <b>Command created</b>,
          error: (err: any) =>
            editingCommand ? (
              <b>
                Failed to save command:
                <br />
                {err}
              </b>
            ) : (
              <b>
                Failed to create command:
                <br />
                {err}
              </b>
            ),
        })
        .then(() => {
          navigate({ to: '/commands' });
        });
    },
    [name, command, editingCommand, commandFormSubmit]
  );

  useEffect(() => {
    if (editingCommand) {
      const command = commands?.find((c) => c.ID === editingCommand);

      if (command) {
        setName(command.Name);
        setCommand(command.Command);
      }
    }
  }, [editingCommand, setName, setCommand]);

  return (
    <form onSubmit={submit} className="flex flex-col w-full max-w-2xl">
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
          <label htmlFor="command" className="w-min">
            Command
          </label>
          <Input
            id="command"
            name="command"
            type="text"
            required
            value={command}
            onChange={(e) => setCommand(e.target.value)}
          />
        </div>

        <Button type="submit" className="mt-2">
          {submitText}
        </Button>
      </div>
    </form>
  );
}
