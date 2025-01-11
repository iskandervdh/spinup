import { PencilSquareIcon, TrashIcon } from '@heroicons/react/20/solid';
import { useCallback } from 'react';
import toast from 'react-hot-toast';
import { sqlc } from 'wjs/go/models';
import { Button } from '~/components/button';
import { useCommandsStore } from '~/stores/commandsStore';
import { Page, usePageStore } from '~/stores/pageStore';

export function CommandInfo({ command }: { command: sqlc.Command }) {
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
