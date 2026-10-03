import { variants } from '@ohjanus/primitives';

export const buttonVariants = variants({
  base: 'ohjanus-btn',
  variants: {
    variant: {
      primary: 'ohjanus-btn-primary',
      secondary: 'ohjanus-btn-secondary',
      outline: 'ohjanus-btn-outline',
      ghost: 'ohjanus-btn-ghost',
      danger: 'ohjanus-btn-danger',
      link: 'ohjanus-btn-link'
    },
    size: {
      xs: 'ohjanus-btn-xs',
      sm: 'ohjanus-btn-sm',
      md: 'ohjanus-btn-md',
      lg: 'ohjanus-btn-lg'
    }
  },
  defaultVariants: {
    variant: 'secondary',
    size: 'sm'
  }
});
