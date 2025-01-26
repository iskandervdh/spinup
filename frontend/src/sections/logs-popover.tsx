import { ArrowDownIcon, XMarkIcon } from '@heroicons/react/20/solid';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useProjectsStore } from '~/stores/projectsStore';
import { FollowProjectLogs, StopFollowingProjectLogs } from 'wjs/go/app/App';
import AnsiToHtml from 'ansi-to-html';
import { EventsOn } from 'wjs/runtime/runtime';
import { Button } from '~/components/button';
import { cn } from '~/utils/helpers';

const LOG_SIZE_LIMIT = 100000;

export function LogsPopover() {
  const ansiToHtml = useMemo(() => new AnsiToHtml({ stream: true }), []);

  const { currentProject, setCurrentProject } = useProjectsStore();
  const [ansiLogs, setAnsiLogs] = useState('');
  const [followLogs, setFollowLogs] = useState(true);

  const logs = useMemo(() => ansiToHtml.toHtml(ansiLogs), [ansiLogs]);

  const logsRef = useRef<HTMLDivElement | null>(null);

  const scrollToBottom = useCallback(() => {
    setTimeout(() => {
      logsRef.current?.scrollTo(0, logsRef.current.scrollHeight);
    }, 10);
  }, [logsRef]);

  const close = useCallback(() => {
    setCurrentProject(null);
  }, [setCurrentProject]);

  useEffect(() => {
    if (!currentProject) return;

    FollowProjectLogs(currentProject);

    const logBuffer: string[] = [];
    let bufferTimeout: NodeJS.Timeout | null = null;

    const flushBuffer = () => {
      setAnsiLogs((prevLogs) => {
        const logs = prevLogs + logBuffer.join('');
        logBuffer.length = 0;

        // Limit logs to LOG_SIZE_LIMIT characters to prevent browser from freezing
        if (logs.length > LOG_SIZE_LIMIT) {
          return logs.slice(-LOG_SIZE_LIMIT);
        }

        return logs;
      });

      bufferTimeout = null;
    };

    const stopListeningForLogs = EventsOn('log', (newLogs: string) => {
      logBuffer.push(newLogs);

      if (!bufferTimeout) {
        bufferTimeout = setTimeout(flushBuffer, 100);
      }
    });

    const scrollListener = (e: Event) => {
      const target = e.target as HTMLElement;

      setFollowLogs(target.scrollTop + target.clientHeight >= target.scrollHeight);
    };

    const currentLogsRef = logsRef.current;
    currentLogsRef?.addEventListener('scroll', scrollListener);

    return () => {
      StopFollowingProjectLogs(currentProject);

      setAnsiLogs('');
      stopListeningForLogs();
      currentLogsRef?.removeEventListener('scroll', scrollListener);
    };
  }, [currentProject]);

  useEffect(() => {
    if (followLogs) {
      scrollToBottom();
    }
  }, [followLogs, scrollToBottom, logsRef.current?.scrollHeight]);

  if (!currentProject) return null;

  return (
    <div className="fixed top-0 left-0 w-full h-full bg-black/50" onClick={close}>
      <div className="absolute flex items-center justify-center w-full h-full p-8">
        <div
          className="relative flex flex-col w-full h-full max-h-screen p-4 rounded-lg bg-background"
          onClick={(e) => e.stopPropagation()}
        >
          <div className="flex items-center justify-between">
            <h1 className="text-xl">Logs for project {currentProject}</h1>
            <button onClick={close} className="p-2 rounded-lg hover:bg-black/10">
              <XMarkIcon width={24} height={24} />
            </button>
          </div>
          <div ref={logsRef} className="w-full h-full mt-4 overflow-y-auto rounded-lg bg-black/50">
            <pre className="p-4 text-sm text-wrap" dangerouslySetInnerHTML={{ __html: logs }} />
          </div>

          <div className={cn('absolute -translate-x-1/2 left-1/2 bottom-8 w-fit', followLogs ? 'hidden' : 'flex')}>
            <Button onClick={scrollToBottom} className="flex items-center gap-1">
              <ArrowDownIcon width={16} height={16} />
              <span>Scroll to bottom</span>
              <ArrowDownIcon width={16} height={16} />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
