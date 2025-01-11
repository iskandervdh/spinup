import { create } from 'zustand';
import { Commands } from '~/types';
import { AddCommand, GetCommands, RemoveCommand, UpdateCommand } from 'wjs/go/app/App';

interface CommandsState {
  commands: Commands | null;
  setCommands: (commands: Commands) => void;

  editingCommand: string | null;
  setEditingCommand: (commandName: string | null) => void;

  commandFormSubmit: (commandName: string, command: string) => Promise<void>;
  removeCommand: (commandName: string) => Promise<void>;
}

export const useCommandsStore = create<CommandsState>((set, get) => ({
  commands: null,
  setCommands: (projects) => set(() => ({ commands: projects })),

  editingCommand: null,
  setEditingCommand: (commandName) => set(() => ({ editingCommand: commandName })),

  commandFormSubmit: async (commandName, command) => {
    if (get().editingCommand === commandName) {
      await UpdateCommand(commandName, command);
      set(() => ({ editingCommand: null }));
    } else {
      await AddCommand(commandName, command);
    }

    const commands = await GetCommands();
    set(() => ({ commands }));
  },
  removeCommand: async (commandName) => {
    await RemoveCommand(commandName);

    const commands = await GetCommands();
    set(() => ({ commands }));
  },
}));
