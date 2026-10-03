<script lang="ts">
  import Database from 'reicon-svelte/icons/Database.svelte';
  import Grid from 'reicon-svelte/icons/Grid.svelte';
  import Shield from 'reicon-svelte/icons/Shield.svelte';
  import Key from 'reicon-svelte/icons/Key.svelte';
  import Document from 'reicon-svelte/icons/Document.svelte';
  import Search from 'reicon-svelte/icons/Search.svelte';
  import Settings from 'reicon-svelte/icons/Settings.svelte';
  import Refresh from 'reicon-svelte/icons/Refresh.svelte';
  import Plus from 'reicon-svelte/icons/Plus.svelte';
  import Minus from 'reicon-svelte/icons/Minus.svelte';
  import X from 'reicon-svelte/icons/X.svelte';
  import Check from 'reicon-svelte/icons/Check.svelte';
  import ChevronDown from 'reicon-svelte/icons/ChevronDown.svelte';
  import ChevronRight from 'reicon-svelte/icons/ChevronRight.svelte';
  import ChevronLeft from 'reicon-svelte/icons/ChevronLeft.svelte';
  import ChevronUp from 'reicon-svelte/icons/ChevronUp.svelte';
  import TerminalSquare from 'reicon-svelte/icons/TerminalSquare.svelte';
  import Clock from 'reicon-svelte/icons/Clock.svelte';
  import Filter from 'reicon-svelte/icons/Filter.svelte';
  import Trash from 'reicon-svelte/icons/Trash.svelte';
  import Lock from 'reicon-svelte/icons/Lock.svelte';
  import Unlock from 'reicon-svelte/icons/Unlock.svelte';
  import Eye from 'reicon-svelte/icons/Eye.svelte';
  import EyeOff from 'reicon-svelte/icons/EyeOff.svelte';
  import Copy from 'reicon-svelte/icons/Copy.svelte';
  import Download from 'reicon-svelte/icons/Download.svelte';
  import Upload from 'reicon-svelte/icons/Upload.svelte';
  import Play from 'reicon-svelte/icons/Play.svelte';
  import Folder from 'reicon-svelte/icons/Folder.svelte';
  import Bell from 'reicon-svelte/icons/Bell.svelte';
  import Lightning from 'reicon-svelte/icons/Lightning.svelte';
  import Chart from 'reicon-svelte/icons/Chart.svelte';
  import Cpu from 'reicon-svelte/icons/Cpu.svelte';
  import User from 'reicon-svelte/icons/User.svelte';
  import Users from 'reicon-svelte/icons/Users.svelte';
  import More from 'reicon-svelte/icons/More.svelte';
  import AlertTriangle from 'reicon-svelte/icons/AlertTriangle.svelte';
  import ArrowUp from 'reicon-svelte/icons/ArrowUp.svelte';
  import ArrowDown from 'reicon-svelte/icons/ArrowDown.svelte';
  import ArrowLeft from 'reicon-svelte/icons/ArrowLeft.svelte';
  import ArrowRight from 'reicon-svelte/icons/ArrowRight.svelte';
  import Undo from 'reicon-svelte/icons/Undo.svelte';
  import Redo from 'reicon-svelte/icons/Redo.svelte';
  import Pin from 'reicon-svelte/icons/Pin.svelte';
  import Printer from 'reicon-svelte/icons/Printer.svelte';
  import Stop from 'reicon-svelte/icons/Stop.svelte';
  import Maximize from 'reicon-svelte/icons/Maximize.svelte';
  import Layout from 'reicon-svelte/icons/Layout.svelte';
  import type { IconName, IconProps } from './types.js';

  let {
    name,
    size = 16,
    color = 'currentColor',
    strokeWidth,
    class: className = '',
    ...restProps
  }: IconProps & { name: IconName } = $props();
</script>

{#if name === 'database'}
  <Database {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'table'}
  <Grid {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'shield'}
  <Shield {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'key'}
  <Key {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'audit'}
  <Document {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'search'}
  <Search {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'settings'}
  <Settings {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'refresh'}
  <Refresh {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'plus'}
  <Plus {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'minus'}
  <Minus {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'x' || name === 'close'}
  <X {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'check'}
  <Check {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'chevron-down'}
  <ChevronDown {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'chevron-right'}
  <ChevronRight {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'chevron-left'}
  <ChevronLeft {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'chevron-up'}
  <ChevronUp {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'terminal'}
  <TerminalSquare {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'clock'}
  <Clock {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'filter'}
  <Filter {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'trash'}
  <Trash {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'lock'}
  <Lock {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'unlock'}
  <Unlock {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'eye'}
  <Eye {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'eye-off'}
  <EyeOff {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'copy'}
  <Copy {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'download'}
  <Download {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'upload'}
  <Upload {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'play'}
  <Play {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'folder'}
  <Folder {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'bell'}
  <Bell {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'lightning'}
  <Lightning {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'chart'}
  <Chart {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'cpu'}
  <Cpu {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'user'}
  <User {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'users'}
  <Users {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'more'}
  <More {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'alert-triangle'}
  <AlertTriangle {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'arrow-up'}
  <ArrowUp {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'arrow-down'}
  <ArrowDown {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'arrow-left'}
  <ArrowLeft {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'arrow-right'}
  <ArrowRight {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'undo'}
  <Undo {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'redo'}
  <Redo {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'pin'}
  <Pin {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'printer'}
  <Printer {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'stop'}
  <Stop {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'maximize'}
  <Maximize {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{:else if name === 'layout'}
  <Layout {size} {color} {strokeWidth} class="ohjanus-icon {className}" {...restProps} />
{/if}
