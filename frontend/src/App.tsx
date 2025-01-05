import { useEffect, useMemo, useRef } from 'react';
import { usePageStore } from '~/stores/pageStore';
import { Navbar } from '~/sections/navbar';
import { Toaster } from 'react-hot-toast';

function App() {
  const { currentPage, getPageComponent } = usePageStore();
  const contentContainerRef = useRef<HTMLDivElement | null>(null);

  const CurrentPage = useMemo(() => getPageComponent(currentPage), [currentPage, getPageComponent]);

  // Scroll to top when changing pages
  useEffect(() => {
    if (contentContainerRef.current) {
      contentContainerRef.current.scrollTop = 0;
    }
  }, [currentPage]);

  return (
    <div id="App" className="flex flex-col h-screen text-white bg-background font-azeret">
      <Navbar />

      <Toaster position="bottom-right" reverseOrder={false} />

      <div className="p-4 overflow-y-auto" ref={contentContainerRef}>
        <CurrentPage />
      </div>
    </div>
  );
}

export default App;
