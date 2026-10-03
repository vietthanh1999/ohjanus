import { variants } from '@ohjanus/primitives';

export const buttonVariants = variants({
  base: 'ohjanus-btn',
  variants: {
    variant: {
      default: 'ohjanus-btn-primary',
      primary: 'ohjanus-btn-primary',
      secondary: 'ohjanus-btn-secondary',
      outline: 'ohjanus-btn-outline',
      ghost: 'ohjanus-btn-ghost',
      destructive: 'ohjanus-btn-danger',
      danger: 'ohjanus-btn-danger',
      link: 'ohjanus-btn-link'
    },
    size: {
      default: 'ohjanus-btn-md',
      xs: 'ohjanus-btn-xs',
      sm: 'ohjanus-btn-sm',
      md: 'ohjanus-btn-md',
      lg: 'ohjanus-btn-lg',
      xl: 'ohjanus-btn-xl',
      icon: 'ohjanus-btn-icon',
      'icon-xs': 'ohjanus-btn-icon-xs',
      'icon-sm': 'ohjanus-btn-icon-sm',
      'icon-lg': 'ohjanus-btn-icon-lg'
    }
  },
  defaultVariants: {
    variant: 'default',
    size: 'default'
  }
});
