import { DetailedHTMLProps, Dispatch, SelectHTMLAttributes, SetStateAction } from 'react';
import { Checkbox } from './checkbox';
import { cn } from '~/utils/helpers';

interface SelectMultipleProps extends DetailedHTMLProps<SelectHTMLAttributes<HTMLDivElement>, HTMLDivElement> {
  options: string[];
  value: string[];
  onChanged: Dispatch<SetStateAction<string[]>>;
}

export function SelectMultiple({ options, value, onChanged, className, ...props }: SelectMultipleProps) {
  return (
    <div
      {...props}
      className={cn('flex flex-col gap-2 max-h-36 border border-primary rounded-lg px-3 py-2 overflow-auto', className)}
    >
      {options.map((name) => (
        <div key={name} className="flex items-center gap-2">
          <Checkbox
            id={`${props.id}-${name}`}
            checked={value.includes(name)}
            onChange={(e) => onChanged((prev) => (e.target.checked ? [...prev, name] : prev.filter((n) => n !== name)))}
          />
          <label htmlFor={`${props.id}-${name}`}>{name}</label>
        </div>
      ))}
    </div>
  );
}
