import { DetailedHTMLProps, Dispatch, SelectHTMLAttributes, SetStateAction } from 'react';
import { Checkbox } from './checkbox';
import { cn } from '~/utils/helpers';

interface SelectMultipleProps extends DetailedHTMLProps<SelectHTMLAttributes<HTMLDivElement>, HTMLDivElement> {
  options: {
    label: string;
    value: string;
    icon?: string | undefined | null;
  }[];
  value: string[];
  onChanged: Dispatch<SetStateAction<string[]>>;
}

export function SelectMultiple({ options, value, onChanged, className, ...props }: SelectMultipleProps) {
  return (
    <div
      {...props}
      className={cn('flex flex-col gap-2 max-h-36 border border-primary rounded-lg px-3 py-2 overflow-auto', className)}
    >
      {options.map((option) => (
        <div key={option.value} className="flex items-center gap-2">
          <Checkbox
            id={`${props.id}-${option.value}`}
            checked={value.includes(option.value)}
            onChange={(e) =>
              onChanged((prev) => (e.target.checked ? [...prev, option.value] : prev.filter((n) => n !== option.value)))
            }
            className="flex-shrink-0"
          />
          <label htmlFor={`${props.id}-${option.value}`} className="flex items-center gap-2 cursor-pointer">
            {option.icon && <img src={option.icon} width={20} height={20} />}
            {option.label}
          </label>
        </div>
      ))}
    </div>
  );
}
