import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useEffect } from 'react';

export const Route = createFileRoute('/')({
  component: Index,
});

function Index() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate({ to: '/projects' });
  }, [navigate]);

  return null;
}
