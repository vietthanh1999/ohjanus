<script lang="ts">
  interface Props {
    orientation?: 'horizontal' | 'vertical';
    class?: string;
    style?: string;
  }

  let {
    orientation = 'horizontal',
    class: className = '',
    style = ''
  }: Props = $props();
</script>

<hr
  class="ohjanus-divider ohjanus-divider-{orientation} {className}"
  {style}
/>
