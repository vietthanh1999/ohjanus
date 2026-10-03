<script lang="ts">
  interface Props {
    orientation?: 'vertical' | 'horizontal';
    class?: string;
    style?: string;
  }

  let { orientation = 'vertical', class: className = '', style = '' }: Props = $props();
</script>

<div
  class="ohjanus-scroll-bar ohjanus-scroll-bar-{orientation} {className}"
  {style}
  aria-hidden="true"
></div>
