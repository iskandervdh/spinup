import icon from '~/assets/images/icon-transparent.png';
import { DEFAULT_PAGE, NAV_PAGES, Page, usePageStore } from '~/stores/pageStore';

function NavbarItem({ page }: { page: Page }) {
  const { currentPage, setCurrentPage } = usePageStore();

  return (
    <li>
      <button
        className="px-2 py-1 text-lg rounded-lg cursor-pointer text-primary hover:bg-black/10 data-[current]:bg-black/10 focus:outline-offset-2 focus-visible:outline focus:outline-1 focus:outline-primary"
        onClick={() => setCurrentPage(page)}
        data-current={page === currentPage ? 'current' : undefined}
      >
        {page}
      </button>
    </li>
  );
}

export function Navbar() {
  const { setCurrentPage } = usePageStore();

  return (
    <div className="flex items-center h-20 gap-8 px-6 border-b select-none min-h-20 border-primary/50 overflow-clip">
      <img
        src={icon}
        id="logo"
        alt="logo"
        className="w-12 h-12 cursor-pointer"
        draggable={false}
        onClick={() => setCurrentPage(DEFAULT_PAGE)}
      />

      <ul className="flex gap-8">
        {NAV_PAGES.map((page) => (
          <NavbarItem key={page} page={page} />
        ))}
      </ul>
    </div>
  );
}
