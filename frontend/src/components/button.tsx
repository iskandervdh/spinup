import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';

import { cn } from '~/utils/helpers';

const buttonVariants = cva(
  'inline-flex items-center justify-center text-sm font-medium transition-colors border rounded-lg whitespace-nowrap ring-offset-background disabled:pointer-events-none disabled:opacity-50 focus:outline-offset-2 focus-visible:outline focus:outline-1',
  {
    variants: {
      variant: {
        default:
          'text-primary focus:outline-primary border-primary hover:text-white bg-background hover:bg-primary-dark hover:border-primary-dark',
        success:
          'text-success border-success focus:outline-success hover:text-white bg-background hover:bg-success hover:border-success',
        info: 'text-info border-info focus:outline-info hover:text-white bg-background hover:bg-info hover:border-info',
        warning:
          'text-warning border-warning focus:outline-warning hover:text-white bg-background hover:bg-warning hover:border-warning',
        error:
          'text-error border-error focus:outline-error hover:text-white bg-background hover:bg-error hover:border-error',
        // outline: 'border border-input bg-background hover:bg-accent hover:text-accent-foreground',
        // secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary/80',
        // ghost: 'hover:bg-accent hover:text-accent-foreground',
        // link: 'text-primary underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-10 px-4 py-2',
        sm: 'px-3 h-9',
        xs: 'h-8 px-2',
        lg: 'px-8 h-11',
        'icon-lg': 'w-10 h-10',
        icon: 'w-8 h-8',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : 'button';
    return <Comp className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />;
  }
);

Button.displayName = 'Button';

export { Button, buttonVariants };
