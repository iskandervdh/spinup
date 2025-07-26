import { create } from 'zustand';
import { Commands } from '~/types';
import { AddCommand, GetCommands, RemoveCommand, UpdateCommand } from 'wjs/go/app/App';

interface CommandsState {
  commands: Commands | null;
  setCommands: (commands: Commands) => void;

  editingCommand: number | null;
  setEditingCommand: (commandName: number | null) => void;

  commandFormSubmit: (commandName: string, command: string) => Promise<void>;
  removeCommand: (commandID: number) => Promise<void>;
}

export const useCommandsStore = create<CommandsState>((set, get) => ({
  commands: null,
  setCommands: (commands) => set({ commands }),

  editingCommand: null,
  setEditingCommand: (commandName) => set({ editingCommand: commandName }),

  commandFormSubmit: async (commandName, command) => {
    const commandID = get().editingCommand;

    if (commandID !== null) {
      await UpdateCommand(commandID, commandName, command);
      set({ editingCommand: null });
    } else {
      await AddCommand(commandName, command);
    }

    const commands = await GetCommands();
    set({ commands });
  },
  removeCommand: async (commandID) => {
    await RemoveCommand(commandID);

    const commands = await GetCommands();
    set({ commands });
  },
}));
