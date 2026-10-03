<script lang="ts">
  interface Props {
    mode?: 'single' | 'multiple' | 'range';
    selected?: Date | Date[] | { from?: Date; to?: Date };
    defaultMonth?: Date;
    disabled?: (date: Date) => boolean;
    class?: string;
    style?: string;
    onSelect?: (value: any) => void;
  }

  let {
    mode = 'single',
    selected = $bindable(undefined),
    defaultMonth,
    disabled,
    class: className = '',
    style = '',
    onSelect
  }: Props = $props();

  // svelte-ignore state_referenced_locally
  let cursor = $state(defaultMonth ?? (selected instanceof Date ? selected : new Date()));

  const year = $derived(cursor.getFullYear());
  const month = $derived(cursor.getMonth());

  const monthLabel = $derived(
    cursor.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
  );

  const cells = $derived(() => {
    const first = new Date(year, month, 1);
    const startDay = first.getDay();
    const daysInMonth = new Date(year, month + 1, 0).getDate();
    const daysInPrev = new Date(year, month, 0).getDate();
    const out: { date: Date; outside: boolean }[] = [];
    for (let i = startDay - 1; i >= 0; i--) {
      out.push({ date: new Date(year, month - 1, daysInPrev - i), outside: true });
    }
    for (let d = 1; d <= daysInMonth; d++) {
      out.push({ date: new Date(year, month, d), outside: false });
    }
    while (out.length % 7 !== 0 || out.length < 42) {
      const last = out[out.length - 1].date;
      out.push({ date: new Date(last.getFullYear(), last.getMonth(), last.getDate() + 1), outside: true });
      if (out.length >= 42) break;
    }
    return out;
  });

  function isSelected(date: Date): boolean {
    if (selected instanceof Date) {
      return sameDay(selected, date);
    }
    if (Array.isArray(selected)) {
      return selected.some((d) => sameDay(d, date));
    }
    if (selected && typeof selected === 'object') {
      const { from, to } = selected as { from?: Date; to?: Date };
      if (from && to) return date >= startOfDay(from) && date <= endOfDay(to);
      if (from) return sameDay(from, date);
    }
    return false;
  }

  function sameDay(a: Date, b: Date): boolean {
    return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();
  }

  function startOfDay(d: Date): Date {
    return new Date(d.getFullYear(), d.getMonth(), d.getDate());
  }

  function endOfDay(d: Date): Date {
    return new Date(d.getFullYear(), d.getMonth(), d.getDate(), 23, 59, 59, 999);
  }

  function pick(date: Date) {
    if (disabled?.(date)) return;
    if (mode === 'single') {
      selected = date;
      onSelect?.(date);
    } else if (mode === 'multiple') {
      const arr = Array.isArray(selected) ? [...selected] : [];
      const idx = arr.findIndex((d) => sameDay(d, date));
      if (idx >= 0) arr.splice(idx, 1);
      else arr.push(date);
      selected = arr;
      onSelect?.(arr);
    } else {
      const cur = (selected ?? {}) as { from?: Date; to?: Date };
      if (!cur.from || (cur.from && cur.to)) {
        selected = { from: date, to: undefined };
        onSelect?.(selected);
      } else {
        const from = cur.from < date ? cur.from : date;
        const to = cur.from < date ? date : cur.from;
        selected = { from, to };
        onSelect?.(selected);
      }
    }
  }

  function prevMonth() {
    cursor = new Date(year, month - 1, 1);
  }

  function nextMonth() {
    cursor = new Date(year, month + 1, 1);
  }

  const WEEKDAYS = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'];
</script>

<div class="ohjanus-calendar {className}" {style} role="application" aria-label={monthLabel}>
  <div class="ohjanus-calendar-header">
    <button type="button" class="ohjanus-calendar-nav" onclick={prevMonth} aria-label="Previous month">‹</button>
    <div class="ohjanus-calendar-caption">{monthLabel}</div>
    <button type="button" class="ohjanus-calendar-nav" onclick={nextMonth} aria-label="Next month">›</button>
  </div>
  <div class="ohjanus-calendar-grid" role="grid">
    {#each WEEKDAYS as wd}
      <div class="ohjanus-calendar-weekday" role="columnheader">{wd}</div>
    {/each}
    {#each cells() as { date, outside } (date.toISOString())}
      {@const dis = disabled?.(date) ?? false}
      {@const sel = isSelected(date)}
      <button
        type="button"
        role="gridcell"
        aria-selected={sel}
        disabled={dis}
        class="ohjanus-calendar-day"
        class:outside
        class:selected={sel}
        onclick={() => pick(date)}
      >
        {date.getDate()}
      </button>
    {/each}
  </div>
</div>
