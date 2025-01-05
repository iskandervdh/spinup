import { useCallback, useEffect } from 'react';
import { PageTitle } from '../components/page-title';
import { GetCommands } from '../../wailsjs/go/app/App';
import { useCommandsStore } from '~/stores/commandsStore';
import { Button } from '~/components/button';
import { ArrowPathIcon, InformationCircleIcon, PencilSquareIcon, PlusIcon, TrashIcon } from '@heroicons/react/20/solid';
import { Page, usePageStore } from '~/stores/pageStore';
import toast from 'react-hot-toast';
import { sqlc } from 'wjs/go/models';

function CommandInfo({ command }: { command: sqlc.Command }) {
  const { setCurrentPage } = usePageStore();
  const { removeCommand, setEditingCommand } = useCommandsStore();

  const edit = useCallback(() => {
    setEditingCommand(command.Name);
    setCurrentPage(Page.CommandForm);
  }, [command.Name, setEditingCommand, setCurrentPage]);

  const remove = useCallback(async () => {
    if (confirm(`Are you sure you want to remove command "${command.Name}"?`)) {
      await removeCommand(command.Name);

      toast.success(<b>Removed command "{command.Name}"</b>);
    }
  }, [command.Name, removeCommand]);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <h3 className="pr-2 text-xl font-bold text-primary">{command.Name}</h3>

        <Button onClick={edit} size={'xs'} title="Edit command">
          <PencilSquareIcon width={16} height={16} className="text-current" />
        </Button>

        <Button onClick={remove} size={'icon'} variant={'error'} title="Remove command">
          <TrashIcon width={16} height={16} className="text-current" />
        </Button>
      </div>

      <div className="grid grid-cols-[16rem,auto]">
        <div>Command</div>
        <div className="text-sm">{command.Command}</div>
      </div>
    </div>
  );
}

export function CommandsPage() {
  const { setCurrentPage } = usePageStore();
  const { commands, setCommands, setEditingCommand } = useCommandsStore();

  useEffect(() => {
    GetCommands().then((commands) => setCommands(commands || []));
  }, []);

  return (
    <div id="commands">
      <div className="flex items-center gap-4 pb-4">
        <PageTitle>Commands</PageTitle>

        <div className="flex gap-2">
          <Button onClick={() => GetCommands().then(setCommands)} size={'icon-lg'} title="Refresh commands">
            <ArrowPathIcon width={24} height={24} className="text-current" />
          </Button>

          <Button
            onClick={() => {
              setEditingCommand(null);
              setCurrentPage(Page.CommandForm);
            }}
            size={'icon-lg'}
            variant={'success'}
            title="Add command"
          >
            <PlusIcon width={24} height={24} className="text-current" />
          </Button>
        </div>
      </div>

      <div className="flex flex-col gap-4">
        {commands ? (
          commands.length === 0 ? (
            <div className="flex flex-col gap-4 py-2">
              <div className="flex items-center gap-2 text-lg text-gray-300">
                <InformationCircleIcon width={24} height={24} className="text-info" />
                <span>No commands found.</span>
              </div>

              <div>
                <Button
                  onClick={() => {
                    setEditingCommand(null);
                    setCurrentPage(Page.CommandForm);
                  }}
                  size={'xs'}
                  variant={'success'}
                >
                  Add a command
                </Button>
              </div>
            </div>
          ) : (
            commands.map((command) => <CommandInfo key={command.ID} command={command} />)
          )
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </div>
  );
}
