import { createRootRoute, Outlet, useRouterState } from '@tanstack/react-router';
import { useEffect, useRef } from 'react';
import { Toaster } from 'react-hot-toast';
import { Navbar } from '~/sections/navbar';
import { useSettingsStore } from '~/stores/settingsStore';

export const Route = createRootRoute({
  component: () => {
    const { location } = useRouterState();

    const { fetchSettings } = useSettingsStore();

    const routeOutletContainerRef = useRef<HTMLDivElement | null>(null);

    // Scroll to top when changing pages
    useEffect(() => {
      if (routeOutletContainerRef.current) {
        routeOutletContainerRef.current.scrollTop = 0;
      }
    }, [location.pathname]);

    useEffect(() => {
      fetchSettings();
    }, []);

    return (
      <div id="App" className="flex flex-col h-screen text-white bg-background font-azeret">
        <Navbar />

        <Toaster position="bottom-right" reverseOrder={false} />

        <div className="p-4 overflow-y-auto" ref={routeOutletContainerRef}>
          <Outlet />
        </div>
      </div>
    );
  },
});
