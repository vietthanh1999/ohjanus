<script lang="ts">
  interface Props {
    class?: string;
    style?: string;
  }

  let { class: className = '', style = '' }: Props = $props();
</script>

<hr class="ohjanus-field-separator {className}" {style} />
