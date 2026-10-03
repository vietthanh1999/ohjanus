<script lang="ts">
  interface Props {
    class?: string;
    style?: string;
  }

  let { class: className = '', style = '' }: Props = $props();
</script>

<span class="ohjanus-pagination-ellipsis {className}" {style} aria-hidden="true">…</span>
