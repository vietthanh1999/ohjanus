<script lang="ts">
  interface Props {
    size?: 'sm' | 'md' | 'lg';
    class?: string;
  }

  let { size = 'md', class: className = '' }: Props = $props();
</script>

<div class="ohjanus-spinner-loader ohjanus-spinner-{size} {className}" role="status" aria-label="Loading">
</div>
