<script lang="ts">
  interface Props {
    class?: string;
    style?: string;
  }

  let { class: className = '', style = '' }: Props = $props();
</script>

<span class="ohjanus-button-group-separator {className}" {style} role="separator"></span>
