import { Bars3BottomLeftIcon, XMarkIcon } from '@heroicons/react/20/solid';
import { Link, useRouterState } from '@tanstack/react-router';
import { useState } from 'react';
import icon from '~/assets/images/icon-transparent.png';
import { Button } from '~/components/button';
import { cn } from '~/utils/helpers';

const NAV_PAGES = {
  Projects: '/projects',
  Commands: '/commands',
  Settings: '/settings',
} as const;

type NavPage = keyof typeof NAV_PAGES;
type NavRoute = (typeof NAV_PAGES)[keyof typeof NAV_PAGES];

function NavbarItem({ page, route }: { page: NavPage; route: NavRoute }) {
  const router = useRouterState();

  return (
    <li>
      <Button
        className={cn('px-2 py-1 text-lg rounded-lg cursor-pointer text-primary hover:bg-black/10', {
          'bg-black/10': router.location.pathname === route,
        })}
        variant="ghost"
        asChild={true}
      >
        <Link to={route}>{page}</Link>
      </Button>
    </li>
  );
}

export function Navbar() {
  const [showMenu, setShowMenu] = useState(false);

  const toggleMenu = () => {
    setShowMenu((prev) => !prev);
  };

  return (
    <div className="flex items-center h-20 gap-8 px-4 border-b shadow-sm select-none shadow-primary-dark/50 min-h-20 border-primary/80 overflow-clip">
      <Link
        to="/"
        className="rounded-lg focus:outline-offset-2 focus-visible:outline focus:outline-1 focus:outline-primary"
      >
        <img src={icon} id="logo" alt="logo" className="cursor-pointer w-14 h-14" draggable={false} />
      </Link>

      <ul className="hidden gap-4 overflow-x-auto md:flex">
        {Object.entries(NAV_PAGES).map(([page, route]) => (
          <NavbarItem key={page} page={page as keyof typeof NAV_PAGES} route={route} />
        ))}
      </ul>

      <div className="flex justify-end flex-1 md:hidden">
        <Button
          className="px-2 py-1 text-lg rounded-lg cursor-pointer text-primary hover:bg-black/10"
          variant="ghost"
          onClick={toggleMenu}
        >
          {showMenu ? <XMarkIcon width={24} height={24} /> : <Bars3BottomLeftIcon width={24} height={24} />}
        </Button>

        <div>
          {showMenu && (
            <ul className="absolute flex flex-col gap-2 p-4 border rounded-lg shadow-lg bg-background top-16 right-4 md:hidden border-primary/40">
              {Object.entries(NAV_PAGES).map(([page, route]) => (
                <NavbarItem key={page} page={page as keyof typeof NAV_PAGES} route={route} />
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
}
