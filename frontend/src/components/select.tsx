import { SelectHTMLAttributes } from 'react';
import { cn } from '~/utils/helpers';

export function Select({ className, ...props }: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={cn(
        'w-full px-3 py-2 transition-colors duration-200 cursor-pointer ease-in-out border rounded-lg appearance-none outline-offset-2 focus-visible:outline outline-1 outline-primary bg-background text-primary border-primary hover:outline',
        className
      )}
      {...props}
    />
  );
}
