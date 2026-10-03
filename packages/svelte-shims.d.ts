declare module '*.svelte' {
  import type { Component } from 'svelte';
  const component: Component<any, any>;
  export default component;
}

// Global Svelte 5 runes support for standalone tsc
declare global {
  function $state<T>(initial?: T): T;
  namespace $state {
    function raw<T>(initial?: T): T;
    function snapshot<T>(state: T): T;
  }
  function $derived<T>(expression: T): T;
  namespace $derived {
    function by<T>(fn: () => T): T;
  }
  function $effect(fn: () => void | (() => void)): void;
  namespace $effect {
    function pre(fn: () => void | (() => void)): void;
    function root(fn: () => void | (() => void)): () => void;
  }
  function $props<T>(): T;
  function $bindable<T>(fallback?: T): T;
  function $inspect<T>(...values: T[]): { with: (fn: (type: 'init' | 'update', ...values: T[]) => void) => void };
}

export {};
