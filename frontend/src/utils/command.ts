import Bash from '~/assets/images/svg/bash.svg';
import Composer from '~/assets/images/svg/composer.svg';
import Docker from '~/assets/images/svg/docker.svg';
import Go from '~/assets/images/svg/go.svg';
import Laravel from '~/assets/images/svg/laravel.svg';
import NodeJS from '~/assets/images/svg/nodejs.svg';
import Php from '~/assets/images/svg/php.svg';
import Python from '~/assets/images/svg/python.svg';
import Unknown from '~/assets/images/svg/unknown.svg';

export type CommandInfo = { title: string; name: string; icon: string };

export function getCommandIcon(command: string | undefined | null): string {
  if (!command) return Unknown;

  const commandIconMap = {
    'npm ': NodeJS,
    'npx ': NodeJS,
    'node ': NodeJS,
    'yarn ': NodeJS,
    'pnpm ': NodeJS,
    'pnpx ': NodeJS,
    'bun ': NodeJS,
    'bunx ': NodeJS,
    'laravel ': Laravel,
    'php artisan ': Laravel,
    'php ': Php,
    'composer ': Composer,
    'docker ': Docker,
    'docker-compose ': Docker,
    'go ': Go,
    'python ': Python,
    'python3 ': Python,
    'pip ': Python,
    'pip3 ': Python,
    'bash ': Bash,
    'sh ': Bash,
  };

  for (const [prefix, icon] of Object.entries(commandIconMap)) {
    if (command.startsWith(prefix)) {
      return icon;
    }
  }

  return Unknown;
}
