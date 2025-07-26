import { useEffect, useState } from 'react';
import { PageTitle } from '~/components/page-title';
import { GetSpinupVersion } from 'wjs/go/app/App';
import { createFileRoute } from '@tanstack/react-router';
import { Checkbox } from '~/components/checkbox';
import { useSettingsStore } from '~/stores/settingsStore';
import { SettingKey } from '~/utils/settings';
import { useShowCommandIcons } from '~/hooks/settings';

export const Route = createFileRoute('/settings')({
  component: Settings,
});

function Settings() {
  const [spinupVersion, setSpinupVersion] = useState<string | null>(null);

  const { setSetting } = useSettingsStore();

  const changeShowCommandIcons = (checked: boolean) => {
    setSetting(SettingKey.ShowCommandIcons, checked);
  };

  const showCommandIcons = useShowCommandIcons();

  useEffect(() => {
    GetSpinupVersion().then(setSpinupVersion);
  }, []);

  return (
    <div id="settings" className="max-w-6xl mx-auto">
      <div className="flex items-center pb-4 h-14">
        <PageTitle>Settings</PageTitle>
      </div>

      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          <h2 className="text-xl font-bold text-primary">Layout</h2>

          <div className="grid items-center max-w-6xl grid-cols-2">
            <div>Show command icons</div>

            <div className="w-full min-w-32 max-w-64">
              <Checkbox
                id="show-command-icons"
                name="Show command icons"
                checked={showCommandIcons}
                onChange={(event) => changeShowCommandIcons(event.target.checked)}
              />
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <h2 className="text-xl font-bold text-primary">Info</h2>

          <div className="grid max-w-6xl grid-cols-[16rem,auto]">
            <div>Spinup version</div>
            <div>{spinupVersion}</div>
          </div>
        </div>
      </div>
    </div>
  );
}
