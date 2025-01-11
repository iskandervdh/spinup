import { Link, useRouterState } from '@tanstack/react-router';
import icon from '~/assets/images/icon-transparent.png';

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
    <Link to={route}>
      <li>
        <button
          className="px-2 py-1 text-lg rounded-lg cursor-pointer text-primary hover:bg-black/10 data-[current]:bg-black/10 focus:outline-offset-2 focus-visible:outline focus:outline-1 focus:outline-primary"
          data-current={router.location.pathname === route ? 'current' : undefined}
        >
          {page}
        </button>
      </li>
    </Link>
  );
}

export function Navbar() {
  return (
    <div className="flex items-center h-20 gap-8 px-4 border-b shadow-sm select-none shadow-primary-dark/50 min-h-20 border-primary/50 overflow-clip">
      <Link to="/">
        <img src={icon} id="logo" alt="logo" className="cursor-pointer w-14 h-14" draggable={false} />
      </Link>

      <ul className="flex gap-8">
        {Object.entries(NAV_PAGES).map(([page, route]) => (
          <NavbarItem key={page} page={page as keyof typeof NAV_PAGES} route={route} />
        ))}
      </ul>
    </div>
  );
}
