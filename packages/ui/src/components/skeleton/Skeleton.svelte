<script lang="ts">
  interface Props {
    shape?: 'text' | 'circle' | 'rect' | 'rounded';
    class?: string;
    style?: string;
  }

  let { shape = 'text', class: className = '', style = '' }: Props = $props();
</script>

<div class="ohjanus-skeleton ohjanus-skeleton-{shape} {className}" {style} aria-hidden="true"></div>
