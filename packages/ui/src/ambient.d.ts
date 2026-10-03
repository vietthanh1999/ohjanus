declare module '*.svelte' {
  import type { Component } from 'svelte';
  const component: Component<any, any>;
  export default component;
}

declare function $state<T>(initial?: T): T;
declare namespace $state {
  function raw<T>(initial?: T): T;
  function snapshot<T>(state: T): T;
}
declare function $derived<T>(expression: T): T;
declare namespace $derived {
  function by<T>(fn: () => T): T;
}
declare function $effect(fn: () => void | (() => void)): void;
declare namespace $effect {
  function pre(fn: () => void | (() => void)): void;
  function root(fn: () => void | (() => void)): () => void;
}
declare function $props<T>(): T;
declare function $bindable<T>(fallback?: T): T;
