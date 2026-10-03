import { cn, type ClassValue } from './cn.js';

export interface VariantConfig<V extends Record<string, Record<string, ClassValue>>> {
  base?: ClassValue;
  variants?: V;
  defaultVariants?: {
    [K in keyof V]?: keyof V[K];
  };
}

export type VariantProps<T extends (...args: any) => any> = Parameters<T>[0];

export function variants<V extends Record<string, Record<string, ClassValue>>>(config: VariantConfig<V>) {
  return (props?: { [K in keyof V]?: keyof V[K] } & { class?: ClassValue }): string => {
    if (!props) {
      return cn(
        config.base,
        config.defaultVariants &&
          Object.entries(config.defaultVariants).map(
            ([variantName, defaultValue]) => config.variants?.[variantName]?.[defaultValue as string]
          )
      );
    }

    const { class: className, ...variantProps } = props;
    const resolvedClasses: ClassValue[] = [config.base];

    if (config.variants) {
      for (const [variantName, variantOptions] of Object.entries(config.variants)) {
        const propValue = (variantProps as Record<string, any>)[variantName] ?? config.defaultVariants?.[variantName];
        if (propValue && variantOptions[propValue]) {
          resolvedClasses.push(variantOptions[propValue]);
        }
      }
    }

    if (className) {
      resolvedClasses.push(className);
    }

    return cn(...resolvedClasses);
  };
}
