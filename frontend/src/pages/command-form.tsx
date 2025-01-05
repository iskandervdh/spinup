import { useCallback, useEffect, useMemo, useState } from 'react';
import { Input } from '~/components/input';
import { PageTitle } from '~/components/page-title';
import { useCommandsStore } from '~/stores/commandsStore';
import { Button } from '~/components/button';
import { Page, usePageStore } from '~/stores/pageStore';
import toast from 'react-hot-toast';

export function CommandFormPage() {
  const { commands, commandFormSubmit, editingCommand, setEditingCommand } = useCommandsStore();
  const { setCurrentPage } = usePageStore();

  const [name, setName] = useState('');
  const [command, setCommand] = useState('');

  const pageTitle = useMemo(
    () => (editingCommand ? `Edit Command "${editingCommand}"` : 'Add Command'),
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
          error: (err) =>
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
          setCurrentPage(Page.Commands);
        });
    },
    [name, command, editingCommand, commandFormSubmit]
  );

  useEffect(() => {
    if (editingCommand) {
      const command = commands?.find((c) => c.Name === editingCommand);

      if (command) {
        setName(editingCommand);
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
          <Input id="name" name="name" type="text" value={name} onChange={(e) => setName(e.target.value)} />
        </div>

        <div className="flex flex-col gap-2">
          <label htmlFor="command" className="w-min">
            Command
          </label>
          <Input id="command" name="command" type="text" value={command} onChange={(e) => setCommand(e.target.value)} />
        </div>

        <Button type="submit" className="mt-2">
          {submitText}
        </Button>
      </div>
    </form>
  );
}
