import { HTMLAttributes } from 'react';
import { cn } from '~/utils/helpers';

export function PageTitle({ className, ...props }: HTMLAttributes<HTMLHeadingElement>) {
  return <h2 className={cn('text-2xl font-bold text-primary', className)} {...props} />;
}
