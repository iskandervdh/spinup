import { DetailedHTMLProps, InputHTMLAttributes } from 'react';
import { cn } from '~/utils/helpers';

export function Checkbox({
  className,
  ...props
}: DetailedHTMLProps<InputHTMLAttributes<HTMLInputElement>, HTMLInputElement>) {
  return (
    <input
      type="checkbox"
      className={cn(
        'w-5 h-5 transition-colors duration-200 ease-in-out border cursor-pointer rounded-md appearance-none outline-offset-2 focus-visible:outline outline-1 outline-primary bg-background text-white border-primary grid place-content-center before:invisible before:checked:visible before:w-3 before:h-3 before:rounded-sm before hover:outline before:bg-primary checked:border-primary',
        className
      )}
      {...props}
    />
  );
}
