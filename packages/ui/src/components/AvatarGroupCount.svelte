<script lang="ts">
  interface Props {
    count?: number;
    label?: string;
    class?: string;
    style?: string;
  }

  let { count = 0, label, class: className = '', style = '' }: Props = $props();
</script>

<span class="ohjanus-avatar-group-count {className}" {style}>
  {label ?? `+${count}`}
</span>
