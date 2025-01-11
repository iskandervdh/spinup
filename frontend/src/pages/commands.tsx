import { useEffect } from 'react';
import { PageTitle } from '~/components/page-title';
import { GetCommands } from 'wjs/go/app/App';
import { useCommandsStore } from '~/stores/commandsStore';
import { Button } from '~/components/button';
import { ArrowPathIcon, InformationCircleIcon, PlusIcon } from '@heroicons/react/20/solid';
import { Page, usePageStore } from '~/stores/pageStore';
import { CommandInfo } from '~/sections/command-info';

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
