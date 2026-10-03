<script lang="ts">
  interface Props {
    src?: string;
    alt?: string;
    class?: string;
    style?: string;
    onload?: () => void;
    onerror?: () => void;
  }

  let { src = '', alt = '', class: className = '', style = '', onload, onerror }: Props = $props();
  let failed = $state(false);
</script>

{#if src && !failed}
  <img
    {src}
    {alt}
    class="ohjanus-avatar-image {className}"
    {style}
    onload={() => { failed = false; onload?.(); }}
    onerror={() => { failed = true; onerror?.(); }}
  />
{/if}
