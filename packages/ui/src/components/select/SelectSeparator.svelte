<script lang="ts">
  interface Props {
    class?: string;
    style?: string;
  }

  let { class: className = '', style = '' }: Props = $props();
</script>

<hr class="ohjanus-select-separator {className}" {style} />
