import { DetailedHTMLProps, InputHTMLAttributes } from 'react';
import { cn } from '~/utils/helpers';

export function Input({
  className,
  ...props
}: DetailedHTMLProps<InputHTMLAttributes<HTMLInputElement>, HTMLInputElement>) {
  return (
    <input
      className={cn(
        'w-full px-3 py-2 transition-colors duration-200 ease-in-out border rounded-lg appearance-none outline-offset-2 focus-visible:outline outline-1 outline-primary bg-background text-white border-primary hover:outline'
      )}
      {...props}
    />
  );
}
