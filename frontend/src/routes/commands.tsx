import { useCallback, useEffect } from 'react';
import { PageTitle } from '~/components/page-title';
import { GetCommands } from 'wjs/go/app/App';
import { useCommandsStore } from '~/stores/commandsStore';
import { Button } from '~/components/button';
import { ArrowPathIcon, InformationCircleIcon, PlusIcon } from '@heroicons/react/20/solid';
import { CommandInfo } from '~/sections/command-info';
import { createFileRoute, Link } from '@tanstack/react-router';

export const Route = createFileRoute('/commands')({
  component: Commands,
});

function Commands() {
  const { commands, setCommands, setEditingCommand } = useCommandsStore();

  const fetchCommands = useCallback(() => {
    GetCommands().then((commands) => setCommands(commands || []));
  }, [setCommands]);

  useEffect(() => {
    fetchCommands();
  }, []);

  return (
    <div id="commands">
      <div className="flex items-center gap-4 pb-4">
        <PageTitle>Commands</PageTitle>

        <div className="flex gap-2">
          <Button onClick={fetchCommands} size="icon-lg" title="Refresh commands">
            <ArrowPathIcon width={24} height={24} className="text-current" />
          </Button>

          <Link to="/command-form" onClick={() => setEditingCommand(null)}>
            <Button size="icon-lg" variant="success" title="Add command">
              <PlusIcon width={24} height={24} className="text-current" />
            </Button>
          </Link>
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
                <Link to="/command-form" onClick={() => setEditingCommand(null)}>
                  <Button size="xs" variant="success">
                    Add a command
                  </Button>
                </Link>
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
