---
name: OhJanus JetBrains Dark Design System
version: "1.0.0"
theme: dark
platform: desktop-ide-web
vibe: professional-database-ide
colors:
  # Surface Hierarchy
  surface:
    canvas: "#1E1F22"          # Main application background, editor canvas, terminal background
    sidebar: "#2B2D30"         # Explorer sidebar, services panel background
    toolbar: "#2B2D30"         # Action bars, tab headers, filter toolbars
    table-header: "#25272A"    # Grid column headers background
    table-row-even: "#1E1F22"  # Data grid even rows
    table-row-odd: "#1E1F22"   # Data grid odd rows (single color grid with borders)
    hover: "#313438"           # Hover state for lists, tree nodes, buttons
    selected: "#2E3A4E"        # Selected table rows, active tree node highlight
    active-tab: "#1E1F22"      # Active document tab surface
    inactive-tab: "#2B2D30"    # Inactive document tab surface
    floating: "#2B2D30"        # Floating row counters, popovers, dropdowns
    execution-block: "#1F4388" # Executing/highlighted SQL block in editor (35-50% alpha over canvas)

  # Border & Dividers
  border:
    subtle: "#2B2D30"          # Inner gutter borders, soft dividers
    default: "#393B40"         # Primary panel splitters, table cell borders, input borders
    strong: "#4E5157"          # Focused borders, active separator lines
    accent: "#3574F0"          # Active tab indicator, focused input ring
    execution: "#3569C8"       # Border surrounding active execution SQL block

  # Text & Content
  text:
    primary: "#DFE1E5"         # High-contrast readable content, table cell values, active tabs
    secondary: "#9DA0A8"       # Secondary labels, table headers, column names, counter badges
    muted: "#7A7E85"           # Timestamps, gutter line numbers, execution latency, comments
    null-value: "#6F737A"      # Muted italic representation for SQL <null>
    disabled: "#5A5D63"        # Inactive icons, disabled controls
    inverted: "#1E1F22"        # Text on high-contrast bright backgrounds

  # Semantic Status & Actions
  action:
    primary: "#3574F0"         # Primary action, blue focus ring, active accents
    primary-hover: "#4682F7"   # Primary action hover
    success: "#57D38C"         # Execute query button, completed status, successful row count
    success-hover: "#66E39C"   # Execute button hover
    warning: "#EDA200"         # Schema warning, spellcheck lint, pending commit
    danger: "#E55353"          # Stop execution button, destructive actions, errors
    danger-hover: "#F76565"    # Stop button hover

  # Badges & Indicators
  badge:
    counter-bg: "#393B40"      # Badge pill background for item counts (e.g., [13], [33])
    counter-fg: "#9DA0A8"      # Badge pill text color
    latency-fg: "#7A7E85"      # Execution latency text color (e.g., 630 ms, 1 s 159 ms)
    license-border: "#275A38"  # Green border for 'Non-commercial use' pill
    license-bg: "#14281B"      # Dark green fill for license pill
    license-fg: "#57D38C"      # Vibrant green text for license pill
    profile-bg: "#0E7490"      # Cyan/teal circular badge for profile initials [VT]
    profile-fg: "#FFFFFF"      # Text inside profile badge

  # Window Chrome (macOS Traffic Lights)
  window:
    close: "#FF5F56"           # Window close button
    minimize: "#FFBD2E"        # Window minimize button
    maximize: "#27C93F"        # Window zoom/maximize button

  # SQL Syntax Highlighting
  syntax:
    keyword: "#CC7832"         # SQL keywords (SELECT, FROM, WHERE, JOIN, LEFT JOIN, ORDER BY, LIMIT)
    function: "#56A8F5"        # SQL functions (now(), interval)
    string: "#6AAB73"          # String literals ('media-transferer', '5 minutes')
    number: "#6897BB"          # Numeric values (0, 10, 501)
    identifier: "#DFE1E5"      # Table/column names, schema names
    alias: "#C792EA"           # Table aliases (ct, q, c, ma)
    comment: "#7A7E85"         # Inline comments (-- streamer poll 10s...)
    spell-warn: "#56A8F5"      # Squiggly wavy underline for unrecognized dictionary words
    code-lens: "#7A7E85"       # Inline relationship hint (1..n <-> 1: on ct.id = ma."contentID")

  # Data Grid Column Type Indicators
  data-types:
    key-pk: "#4A88C7"          # Primary key / UUID column icon (blue key)
    number: "#56A8F5"          # Integer / float / ordinal column icon (#)
    string: "#56A8F5"          # Varchar / text column icon ("T")
    date: "#56A8F5"            # Timestamp / date column icon (calendar/clock)
    boolean: "#56A8F5"         # Boolean flag column icon

typography:
  font-family:
    ui: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
    code: '"JetBrains Mono", Menlo, Monaco, Consolas, "Courier New", monospace'

  font-size:
    badge: "10px"              # Counter badges, micro indicators
    status: "11px"             # Status bar, breadcrumbs, code lens hints
    caption: "11px"            # Metadata, row counters, column data types
    body-sm: "12px"            # Tree items, table cells, toolbar labels, tabs
    body-md: "13px"            # SQL editor code, WHERE/ORDER BY input text
    title-sm: "13px"           # Section headers (Database Explorer, Services)
    title-md: "14px"           # Dialog titles, modal headers

  font-weight:
    regular: 400               # Body, table values, SQL editor code
    medium: 500                # Tab labels, button text, column headers, tree category labels
    semibold: 600              # Panel titles, active states, status badges
    bold: 700                  # Profile avatar text, highlighted keys

  line-height:
    tight: 1.2                 # Single-line tree items, badges
    normal: 1.4                # Table cells, toolbars
    editor: 1.5                # Code editor line spacing (20px at 13px font)

spacing:
  "2xs": "2px"
  "xs": "4px"
  "sm": "6px"
  "md": "8px"
  "lg": "12px"
  "xl": "16px"
  "2xl": "20px"
  "3xl": "24px"

radii:
  none: "0px"                  # Table cells, split panes, status bar
  xs: "2px"                    # Inline tags, code highlights
  sm: "3px"                    # Toolbar icon buttons, input fields
  md: "4px"                    # Document tabs, tree selection highlights, modal elements
  lg: "8px"                    # Floating panels, context menus, dialogs
  pill: "12px"                 # Status badges, license pill, floating row counter badge
  full: "9999px"               # Traffic light buttons, avatar initials

shadows:
  none: "none"
  floating: "0 2px 8px rgba(0, 0, 0, 0.45)"      # Floating row counter (e.g., "58 rows v")
  dropdown: "0 4px 16px rgba(0, 0, 0, 0.55)"     # Context menus, autocomplete popups
  modal: "0 8px 32px rgba(0, 0, 0, 0.70)"        # Overlay modals

layout:
  window-header-height: "38px"
  tab-bar-height: "32px"
  toolbar-height: "30px"
  filter-bar-height: "32px"
  table-row-height: "26px"
  table-header-height: "26px"
  status-bar-height: "24px"
  tree-row-height: "22px"
  sidebar-default-width: "300px"
  sidebar-min-width: "220px"
  splitter-hitbox-width: "4px"
---

# OhJanus / JetBrains DataGrip Dark Design System (`DESIGN.md`)

## 1. Overview & Design Rationale

This document establishes the official visual language, spatial geometry, and UI interaction rules for **OhJanus** — a modern database management tool and developer gateway inspired directly by the high-density, keyboard-driven, precision engineering aesthetic of **JetBrains DataGrip**.

### Core Philosophy
1. **High Information Density**: Screen real estate is dedicated to query logic, structural metadata, and tabular results. Padding is disciplined (2px to 8px), icons are crisp 14-16px, and typography utilizes tabular figures for maximum vertical and horizontal data capacity.
2. **Deep Charcoal Contrast**: Rather than absolute black (`#000000`), the surface uses refined charcoal layers (`#1E1F22`, `#25272A`, `#2B2D30`) with razor-thin 1px borders (`#393B40`), reducing eye fatigue during multi-hour data exploration.
3. **Semantic Typography**: Dual-font architecture — humanistic neutral sans-serif (`Inter`) for UI navigation controls, and ligature-rich monospaced (`JetBrains Mono`) for all SQL statements, tabular cells, execution timestamps, and latency counters.
4. **Immediate Feedback & Observability**: Active sessions, ongoing queries, row counts, and latency timings (`630 ms`, `1 s 159 ms`) are visible inline as contextual badges.

---

## 2. Window Chrome & Global Navigation (Titlebar)

The topmost application chrome mirrors the macOS desktop IDE standard, providing quick access to project context, VCS branching, global search, and IDE settings.

```
+----------------------------------------------------------------------------------------------------------------------------------------+
|  [● ● ●]   [VT] VThanh ▾   Version Control ▾           [🗄️] [▶] [📁] [...]           [🌀]  [🔍]  [⚙️]                               |
+----------------------------------------------------------------------------------------------------------------------------------------+
```

### 2.1. Traffic Lights (macOS Native Window Controls)
- **Position**: Left-aligned, `8px` top/bottom offset, `12px` left margin.
- **Diameter**: `12px` circle with `8px` gap between buttons.
- **Tokens**:
  - Close: `#FF5F56` (hover: `#E0443E`, active: `#BF3630`)
  - Minimize: `#FFBD2E` (hover: `#DEA123`, active: `#BD8316`)
  - Maximize: `#27C93F` (hover: `#1AAB29`, active: `#138A1E`)

### 2.2. Context Selectors
- **Workspace Avatar**:
  - Circle badge `20px x 20px`, background: `#0E7490` (cyan/teal), text: `VT` in white `#FFFFFF`, font-size `10px`, weight `700`.
- **Project Selector**:
  - Text: `VThanh` with down chevron `▾` (`8px` gap), font-size `12px`, weight `500`, color `#DFE1E5`.
- **Version Control Branch**:
  - Text: `Version Control ▾`, font-size `12px`, weight `400`, color `#9DA0A8`. Hover background: `#313438` with `4px` border radius.

### 2.3. Center Quick Action Bar
- Group of `26px x 26px` icon buttons centered horizontally.
- Icons: Database (`🗄️`), Run Configurations (`▶` in circle), Project Explorer (`📁`), More actions (`...`).
- Icon color: `#7A7E85`, hover color: `#DFE1E5`, hover background: `#313438`.

### 2.4. Right Utility Controls
- **AI Assistant**: Spiral / swirl icon (`#9DA0A8`), clickable to open the AI query generator/explainer.
- **Global Search**: Magnifying glass icon (`#7A7E85`), shortcut `Double Shift` / `Cmd+O`.
- **Settings Gear**: Settings cog icon (`#7A7E85`), opens Preferences.

---

## 3. Sidebar Architecture (Left Split Pane)

The left pane is divided vertically into two sub-panels: **Database Explorer** (top) and **Services** (bottom), separated by an adjustable horizontal splitter.

```
+---------------------------------------------+
| Database Explorer       [+][⟳][🔍][DDL][←][👁️]|
| v 🐘 [Dev][ReadOnly] 10.220.6.4        [15] |
|   > 📁 datahub                          [...] |
|   v 📁 dev_mh_asset                      [6] |
|     > ⛁ es                                  |
|     v ⛁ public                              |
|       v 🗄️ tables                       [33] |
|         📄 category                         |
|         📄 connection_credential            | (selected row)
|         📄 content                          |
+=============================================+ <- Horizontal Splitter
| Services                [Tx,][+][👁️][⧉][-][v] |
| ☑ v 📁 Database                             |
|     v 🐘 [Dev][ReadOnly] 10.220.6.4         |
|         📄 connection_credential  1 s 159 ms|
|         📄 events                     986 ms|
|         📄 commands                 1 s 449 ms|
|     v 🐘 [PRD] 10.250.6.23                  |
|         ⚡ console_2                  630 ms|
+---------------------------------------------+
```

### 3.1. Database Explorer
- **Header**:
  - Title: `Database Explorer`, font-size `12px`, weight `600`, color `#DFE1E5`.
  - Action Toolbar: Compact icon buttons (`20px x 20px`):
    - `+` (New Data Source / Driver)
    - `⟳` (Refresh Schema Introspection)
    - `🔍` (Filter Tree Objects)
    - `⧉` (Duplicate/Copy DDL)
    - `DDL` (Generate Schema DDL)
    - `← / →` (Navigate Back / Forward)
    - `👁️` (Show/Hide System Catalogs & Hidden Schemas)
- **Tree Hierarchy & Visual Styling**:
  - Row height: `22px`, font-size `12px`, font-family `ui`.
  - Indentation per depth level: `16px`.
  - Caret toggles: Chevron right `>` for collapsed, Chevron down `v` for expanded (`10px`, `#7A7E85`).
  - **Node Types**:
    - **Instance/Server**: PostgreSQL Elephant icon (`#4A88C7`), text: `[Dev][ReadOnly] 10.220.6.4`, badge: count of databases (e.g. `15`).
    - **Database**: Folder icon with database mark, text: `dev_mh_asset`, badge: count of active schemas (`6`). Unloaded databases show `...` badge.
    - **Schema**: 3-cylinder database cluster icon (`es`, `information_schema`, `marts`, `pg_catalog`, `public`).
    - **Object Group**: Category icon + `tables` (`33`), `routines` (`2`), `views`.
    - **Entity Table**: Grid table icon with blue header row, name: `connection_credential`.
  - **Selection State**:
    - Selected row has background `#2E3A4E` spanning the full width of the tree, text color `#DFE1E5`, border-radius `4px`.
  - **Count Badges**:
    - Pill shape (`16px` height, min-width `18px`, border-radius `8px`), background `#393B40`, text `#9DA0A8`, font-size `10px`, weight `500`.

### 3.2. Services Panel (Execution & Session Monitor)
- **Header**: Title `Services`, toolbar icons: `Tx,` (Transaction mode), `+`, `👁️`, `⧉`, `-` (Collapse all), `▾`.
- **Tree Content**:
  - Checkbox column for batch actions.
  - Active database connections grouped under `Database`.
  - Query consoles and executed table viewers listed with real-time execution duration:
    - `connection_credential` -> `1 s 159 ms`
    - `events` -> `986 ms`
    - `commands` -> `1 s 449 ms`
    - `console_2` -> `630 ms`
  - Latency text: font-family `code`, font-size `11px`, color `#7A7E85`.

---

## 4. Main Work Area: Document Tabs System

All query consoles, DDL files, and table data viewers live inside the top document tab bar.

```
+-----------------------------------------------------------------------------------------------------+
| [⚡ console [[PRD] 10.250.6.23]] | [⚡ console_2 [[PRD] 10.250.6.23]  x] |                         [▾] |
+-----------------------------------------------------------------------------------------------------+
```

### Tab Anatomy & States
- **Tab Height**: `32px`, minimum width `140px`, maximum width `280px`.
- **Leading Icon**:
  - SQL Console: Lightning bolt / console icon (`#4A88C7`).
  - Table Data Viewer: Grid table icon (`#4A88C7`).
- **Tab Label**:
  - Structure: `<ObjectName> [[<Environment>] <Host>]`
  - Examples: `console_2 [[PRD] 10.250.6.23]`, `connectio...credential [[Dev][ReadOnly] 10.220.6.4]`.
  - Middle truncation with ellipsis `...` when title exceeds tab width.
- **Close Button (`x`)**:
  - `14px x 14px` hit target, color `#7A7E85`.
  - Hover: color `#DFE1E5`, background `#393B40`, border-radius `3px`.
- **Tab States**:
  - **Active Tab**: Background `#1E1F22` (matches editor canvas below), bottom border transparent (seamless blend into editor), text `#DFE1E5`, top border or indicator in `#3574F0` (accent blue).
  - **Inactive Tab**: Background `#2B2D30`, text `#9DA0A8`, right border 1px solid `#393B40`.
  - **Overflow Chevron (`▾`)**: Far-right pinned button to reveal tab dropdown list.

---

## 5. View Variant A: SQL Console & Query Editor (`design.png`)

When a SQL console tab is open, the main work area splits into an upper **SQL Query Editor** and a lower **Query Results / Data Table**.

```
+-------------------------------------------------------------------------------------------------------------------+
| [▶] [▶_] [🕒] [⏹] [⚙️] [🗖]   Tx: Auto ▾   ☐ Playground ▾                         [⛁ prd_mh_asset.public ▾]       |
+-------------------------------------------------------------------------------------------------------------------+
| 64 |   join content ct  1..n <-> 1: on ct.id = ma."contentID"                                                     |
| 65 |   LEFT JOIN transfer_job_queue q ON q."jobID" = c.process_id                                                 |
| 66 |   WHERE c.participant = 'media-transferer'                                                    [✓ 6  ▲ ▼]     |
| 67 |   AND c.process_name   = 'media-ingest'                                                        |
| 68 |   AND c.name           = 'create-ingest-job'                                                                 |
| 69 |   AND c.retry_count    = 0                                                                                   |
| 70 |   AND c.schedule IS NULL                                                                                     |
| 71 |   AND c."timestamp"    < now() - interval '5 minutes'  -- streamer poll 10s -> >5 phút chắc chắn không ai đọc |
| 72 |   ORDER BY c.ordinal;                                                                                        |
| 73 |                                                                                                              |
| 74✓|   select * from transfer_job where id = 'b8262180-1075-43f8-8338-2412d4734d65'                               |
+===================================================================================================================+
| [Result 1] [Result 1-2] [prd_mh_asset.public.transfer_job] [Result 1-4] [prd_mh_asset.public.transfer_job 2 x] ▾ |
+-------------------------------------------------------------------------------------------------------------------+
| [⊞][🗎] [⟳][🕒][⏹] [+] [-] [↩][↪] [↑][↓] | Tx: Auto | DDL | 📌 [🔍] [Y] [📊] | CSV ▾ [⤓] [⤒] [📈] [👁️] [⚙️]        |
+-----+--------------------------------------+---------+--------------------------------------+---------------------+
|     | 🔑 id [Y] [⇅]                         | # ordinal | " "connectionCredentialID" [Y] [⇅]   | " "sourcePath" [Y]  |
+-----+--------------------------------------+---------+--------------------------------------+---------------------+
| 1   | b8262180-1075-43f8-8338-2412d4734d65 | 172098  | 753da7e1-523f-4834-ba74-7824e3d5aa52 | /home/vod/VOD/as... |
+-----+--------------------------------------+---------+--------------------------------------+---------------------+
|                                                                                                  [ 1 row ▾ | ⋮ ]  |
+-------------------------------------------------------------------------------------------------------------------+
```

### 5.1. SQL Sub-Toolbar (Action Bar)
- Height: `30px`, background: `#2B2D30`, border-bottom: 1px solid `#393B40`.
- **Buttons (Left Group)**:
  - **Execute Entire Statement** (`▶`): `#57D38C` (vibrant green).
  - **Execute Under Caret** (`▶_`): Green with step bar.
  - **History** (`🕒`): Opens historical queries dialog.
  - **Cancel / Stop** (`⏹`): Red square `#E55353`.
  - **Settings** (`⚙️`): Database session parameters.
  - **Split Editor Layout** (`🗖`): Toggle split orientation.
  - **Transaction Mode**: Dropdown selector `Tx: Auto ▾` (options: `Auto`, `Manual`).
  - **Execution Mode**: Checkbox + `Playground ▾` (runs inside sandbox without auto-commit).
- **Target Schema Selector (Right Pinned)**:
  - Button with database cluster icon + schema path: `prd_mh_asset.public ▾`.
  - Border: 1px solid `#393B40`, border-radius `3px`, background `#1E1F22`, color `#DFE1E5`.

### 5.2. SQL Editor Canvas & Syntax Engine
- **Font**: `"JetBrains Mono"`, font-size `13px`, line-height `20px` (`1.53`).
- **Gutter**:
  - Width: `48px`, background: `#1E1F22`, border-right: 1px solid `#2B2D30`.
  - Line numbers: Color `#7A7E85`, right-aligned.
  - Active line number: Color `#DFE1E5`, weight `500`.
  - Execution Status Mark: Green checkmark (`✓`) at line 74 indicating statement executed successfully.
- **Active Execution Highlight Block**:
  - When a block or statement is being executed/evaluated, the entire block (lines 64 to 72) is enclosed in an active execution region:
    - Background: `#1F4388` (semi-transparent royal blue overlay).
    - Border: 1px solid `#3569C8` around the block perimeter.
    - Floating Result Badge: Pinned at top-right corner of the block: `[✓ 6  ▲ ▼]`
      - Background: `#1E2B37`, border: 1px solid `#2B5B9E`, border-radius `3px`.
      - Green check `✓`, result row count `6`, and up/down jump chevrons.
- **Code Lens (Virtual Relationship Annotation)**:
  - Injected inline annotation between tables: `1..n <-> 1: on ct.id = ma."contentID"`.
  - Font-size: `11px`, font-style: `italic`, color: `#7A7E85`, non-editable virtual text.
- **Syntax Tokens Breakdown**:
  - `join`, `LEFT JOIN`, `WHERE`, `AND`, `ORDER BY`, `select`, `from`: Keywords (`#CC7832` / `#CF8E6D` or bold `#56A8F5`).
  - `content`, `transfer_job_queue`, `transfer_job`: Table identifiers (`#DFE1E5`).
  - `ct`, `q`, `c`, `ma`: Table aliases (`#C792EA`).
  - `'media-transferer'`, `'create-ingest-job'`, `'5 minutes'`: String literals (`#6AAB73`).
  - `now()`, `interval`: Built-in SQL functions (`#56A8F5`, italic).
  - `0`: Numbers (`#6897BB`).
  - `-- streamer poll 10s...`: Comments (`#7A7E85`, italic).
  - Wavy spellcheck squiggly underline (`#56A8F5` wave) under non-English terms (`chắc chắn không ai đọc`).

---

## 6. View Variant B: Table Data Viewer & Log Console (`design2.png`)

When double-clicking a table from the Database Explorer (e.g., `connection_credential`), the view opens a dedicated table data browser with inline WHERE / ORDER BY clause filtering and an execution console below.

```
+-------------------------------------------------------------------------------------------------------------------+
| ⟳ 🕒 ⏹  +  -  ↩ ↪  ↑ ↓   Tx: Auto ▾   DDL   📌  🔍  [Y]  📊                         CSV ▾  ⤓  ⤒  📈  👁️  ⚙️      |
+-------------------------------------------------------------------------------------------------------------------+
| [Y] WHERE [                                                ]  | [⇅] ORDER BY [                             ]      |
+-----+--------------------------------------+----------------------------+----------------------------+------------+
|     | 🔑 id [Y] [⇅]                         | 📅 createdDate [Y] [⇅]      | 📅 lastUpdatedDate [Y] [⇅]  | 👤 createdBy|
+-----+--------------------------------------+----------------------------+----------------------------+------------+
| 1   | 2389a9f7-9ed7-4b7c-9b7c-7eec030b45f5 | 2024-10-29 10:06:03.247644 | <null>                     | 4ebfa335.. |
| 2   | a4430e1f-3b89-4088-bed8-8a426ab293f0 | 2024-10-31 03:56:42.776700 | <null>                     | 0f379df2.. |
| 3   | 3be52421-235b-412e-b420-23b612770315 | 2024-11-05 07:49:06.550088 | <null>                     | 0f379df2.. |
| 4   | 6e98ccab-d08c-49af-bd10-70ca0d510879 | 2024-10-31 03:58:27.775837 | 2024-10-31 04:05:44.633942 | 0f379df2.. |
| 5   | ca1d2eee-6d07-48b5-af45-f551e984f3ed | 2024-11-14 03:30:27.290163 | <null>                     | 0f379df2.. |
| 6   | 9e768c19-b3e5-4027-80a9-25a83307fa14 | 2024-11-14 03:34:10.424030 | <null>                     | 0f379df2.. |
+-----+--------------------------------------+----------------------------+----------------------------+------------+
|                                                                                                  [ 58 rows ▾ | ⋮ ]|
+===================================================================================================================+
| [2026-10-03 10:50:34] Connected to dev_mh_asset                                                          | [📄]   |
| [2026-10-03 10:50:34] dev_mh_asset.public> SELECT t.*                                                    | [🗑️]   |
|                                            FROM public.connection_credential t                           | [↩]   |
|                                            LIMIT 501                                                     | [⤓]   |
| [2026-10-03 10:50:34] 58 rows retrieved starting from 1 in 586 ms (execution: 86 ms, fetching: 500 ms)   | [🖨️]   |
+-------------------------------------------------------------------------------------------------------------------+
```

### 6.1. Inline SQL Filter Bar (WHERE / ORDER BY)
- **Container**: Pinned above the table grid, height `32px`, background `#2B2D30`, border-bottom: 1px solid `#393B40`.
- **Fields**:
  - **WHERE Clause Input**:
    - Leading icon: Funnel `[Y]` (`#7A7E85`).
    - Label prefix: `WHERE`, font-weight `600`, font-size `11px`, color `#9DA0A8`.
    - Input box: Background `#1E1F22`, border: 1px solid `#393B40`, border-radius `3px`, font-family `code`, font-size `12px`, text `#DFE1E5`.
    - Function: Typing `id IS NOT NULL` automatically appends to query without editing raw SQL.
  - **ORDER BY Clause Input**:
    - Leading icon: Sort `[⇅]` (`#7A7E85`).
    - Label prefix: `ORDER BY`, font-weight `600`, font-size `11px`, color `#9DA0A8`.
    - Input box: Same styling, for expressions like `createdDate DESC`.

### 6.2. High-Density Data Grid Specifications
- **Column Header Row**:
  - Height: `26px`, background `#25272A`, border-bottom: 1px solid `#393B40`.
  - Column Separators: 1px solid `#393B40`.
  - Resizing hit zone: `4px` width on header border with cursor `col-resize`.
  - Header Contents:
    - Column Type Icon (`14px`):
      - Primary Key / UUID: Blue key icon (`#4A88C7`)
      - Date / Timestamp: Blue calendar icon (`#56A8F5`)
      - Text / String: Blue quote / `T` icon (`#56A8F5`)
      - Numeric: Blue hash `#` icon (`#56A8F5`)
    - Column Name: `id`, `createdDate`, `lastUpdatedDate`, `createdBy`, font-size `12px`, weight `500`, color `#DFE1E5`.
    - Column Filter Funnel (`Y`): Clicking toggles individual column value filters.
    - Column Sort Indicator (`⇅` / `▲` / `▼`): Click to cycle Sort ASC -> DESC -> None.
- **Row Index Gutter**:
  - Leftmost column, width `36px`, background `#25272A`, text color `#7A7E85`, right-aligned numeric indices `1, 2, 3...`.
- **Data Rows & Cells**:
  - Height: `24px`, background `#1E1F22`, border-bottom: 1px solid `#2B2D30`, border-right: 1px solid `#2B2D30`.
  - Font: `"JetBrains Mono"`, font-size `12px`, tabular numbers (`font-variant-numeric: tabular-nums`).
  - Text Color: `#DFE1E5`.
  - **Null Values**: Displayed explicitly as `<null>` in `#6F737A`, font-style `italic`.
  - **Long String Truncation**: Strings exceeding cell width truncate with trailing ellipsis `...` (e.g. `0f379df2-47d2-47e7-b195-`).
  - **Cell Selection**: Focused cell gets a 1px solid `#3574F0` blue border.

### 6.3. Floating Pagination & Row Counter Badge
- Pinned at bottom right of the table viewport (position: `absolute`, bottom `12px`, right `16px`).
- Style: Pill container, background `#2B2D30`, border: 1px solid `#393B40`, box-shadow: `0 2px 8px rgba(0,0,0,0.45)`, border-radius: `12px`.
- Content:
  - Text: `58 rows ▾` (or `1 row ▾`), font-size `11px`, weight `500`, color `#DFE1E5`. Clicking opens page-size selector (`50`, `100`, `500`, `All`).
  - Divider: Vertical 1px line `#393B40`.
  - Actions: Vertical kebab menu `⋮` (`#9DA0A8`) for export / view preferences.

### 6.4. Console Log / Output Terminal (Bottom Pane)
- Pinned to bottom of the work area, background `#1E1F22`, border-top: 1px solid `#393B40`.
- **Right Action Strip**:
  - Toggle output window (`📄`)
  - Clear log (`🗑️` trash can)
  - Toggle soft-wrap (`↩`)
  - Scroll to end (`⤓`)
  - Print log (`🖨️`)
- **Log Anatomy & Execution Metrics**:
  - Timestamp: `[2026-10-03 10:50:34]` in `#7A7E85`.
  - Connection event: `Connected to dev_mh_asset` in `#DFE1E5`.
  - SQL Executed snippet:
    - Target: `dev_mh_asset.public>`
    - Formatted SQL: Keywords in `#CF8E6D` (`SELECT`, `FROM`, `LIMIT`), numbers in `#6897BB` (`501`).
  - Execution summary line:
    - `58 rows retrieved starting from 1 in 586 ms (execution: 86 ms, fetching: 500 ms)`
    - Clear distinction between database engine execution time (`86 ms`) and client data network fetching time (`500 ms`).

---

## 7. Global Status Bar (Bottom Chrome)

The status bar spans the entire width at the bottom of the window (`24px` height).

```
+----------------------------------------------------------------------------------------------------------------------------------------+
| Database > [Dev][ReadOnly] 10.220.6.4 > dev_mh_asset > public > tables > 🗄️ connection_credential       [🔒] [🔔] [Non-commercial use]|
+----------------------------------------------------------------------------------------------------------------------------------------+
```
*(Or in SQL Console mode):*
```
+----------------------------------------------------------------------------------------------------------------------------------------+
| Database Consoles > [PRD] 10.250.6.23 > ⚡ console_2 [[PRD] 10.250.6.23]       71:42 (2954 chars, 73 line breaks) LF UTF-8 4 spaces [🔒][🔔][Non-commercial use]|
+----------------------------------------------------------------------------------------------------------------------------------------+
```

### 7.1. Left: Interactive Navigation Breadcrumb
- Hierarchy separated by chevron `>`:
  - `Database Consoles` > `[PRD] 10.250.6.23` > `console_2 [[PRD] 10.250.6.23]`
  - `Database` > `[Dev][ReadOnly] 10.220.6.4` > `dev_mh_asset` > `public` > `tables` > `connection_credential`
- Font: `11px`, color `#9DA0A8`, icons `#7A7E85`. Hovering over any breadcrumb segment highlights with `#313438` background, clicking opens navigation popup.

### 7.2. Right: Document & Environment Status
- **Cursor Position & Document Stats**:
  - `71:42 (2954 chars, 73 line breaks)` — line:column format with buffer analytics, color `#7A7E85`, font-family `code`, font-size `11px`.
- **File Encoding & Indentation**:
  - `LF` · `UTF-8` · `4 spaces` — clickable to change CRLF/LF, character set, or tab size.
- **Read-Only Lock Icon**:
  - Lock icon (`🔒`) indicating connection is in `ReadOnly` transaction mode.
- **Notification Bell**:
  - Bell icon (`🔔`), turns blue `#3574F0` on pending background task completion.
- **License / Environment Badge**:
  - Pill badge: border 1px solid `#275A38`, background `#14281B`, text `#57D38C`, text: `Non-commercial use`.
  - Height: `18px`, padding: `0 8px`, border-radius `10px`, font-size `11px`, font-weight `500`.

---

## 8. Micro-Interactions & State Specifications

| Element | Default State | Hover State | Active / Clicked State | Focused / Selected State |
| :--- | :--- | :--- | :--- | :--- |
| **Tree Item** | Transparent bg, `#DFE1E5` text | Bg: `#313438`, `#DFE1E5` text | Bg: `#2B384E` | Bg: `#2E3A4E`, full row width |
| **Document Tab** | Bg: `#2B2D30`, `#9DA0A8` text | Bg: `#35373B`, `#DFE1E5` text | Bg: `#1E1F22`, `#DFE1E5` text | Top border: `2px solid #3574F0` |
| **Toolbar Icon** | Transparent bg, `#7A7E85` fg | Bg: `#35373B`, `#DFE1E5` fg | Bg: `#393B40`, `#3574F0` fg | 1px border `#4E5157` |
| **Grid Cell** | Bg: `#1E1F22`, border: `#2B2D30` | Subtle outline `#393B40` | Editing input cursor | 1px solid border `#3574F0` |
| **Column Header** | Bg: `#25272A`, `#DFE1E5` text | Filter funnel `Y` turns `#DFE1E5` | Sort indicator activates | Resizer bar becomes blue |
| **Action Button (Run ▶)** | Green `#57D38C` | Green hover `#66E39C` | Green active `#499C54` | Focus ring `2px #3574F0` |
| **Action Button (Stop ⏹)**| Red `#E55353` | Red hover `#F76565` | Red active `#DB5860` | Focus ring `2px #3574F0` |
| **Status Badge** | Pill `#393B40`, `#9DA0A8` text | Bg: `#4E5157`, `#DFE1E5` text | Scaled 0.98 | Border: `#3574F0` |

---

## 9. Accessibility & Ergonomics Standards

1. **WCAG 2.1 AA Contrast**:
   - Primary text (`#DFE1E5`) on canvas (`#1E1F22`) delivers **12.4:1 contrast ratio**, exceeding the AAA standard (7:1).
   - Secondary text (`#9DA0A8`) on sidebar (`#2B2D30`) delivers **5.8:1 contrast ratio**, exceeding the AA standard (4.5:1).
   - Keyword text (`#CC7832` / `#56A8F5`) on canvas delivers **5.2:1 contrast ratio**.
2. **Keyboard Accessibility**:
   - `Cmd+Enter` / `Ctrl+Enter`: Execute selected statement in SQL editor.
   - `Cmd+F` / `Ctrl+F`: Focus inline table search / filter input.
   - `Arrow Keys`: Roving tabindex navigation across table grid cells.
   - `Escape`: Clear active cell selection or dismiss floating overlays.
3. **Focus Indicators**:
   - Focused inputs and editable cells must display an unambiguous 1px blue outline (`#3574F0`) with 0px offset.

---

## 10. AI Agent Generation Rules (Do's & Don'ts)

### DO:
- **DO** use `#1E1F22` as the root background for editor panes, tables, and terminal views.
- **DO** use `"JetBrains Mono"` with tabular numbers for all numeric quantities, IDs, timestamps, and SQL code.
- **DO** format `<null>` values explicitly in italic `#6F737A` rather than leaving table cells blank.
- **DO** truncate long UUIDs and file paths with an ellipsis (`...`) while preserving full strings in tooltips or details panes.
- **DO** render counter badges (`[15]`, `[33]`) as compact rounded pills with `#393B40` background.
- **DO** include execution latency badges (`630 ms`, `1 s 159 ms`) next to query consoles and executed tables.
- **DO** enclose active execution SQL code blocks in a semi-transparent royal blue overlay (`#1F4388`) with floating result pill `[✓ 6 ▲ ▼]`.

### DON'T:
- **DON'T** use light themes or pure black (`#000000`) surfaces.
- **DON'T** use rounded pill corners on buttons, inputs, or tabs; buttons and tabs must have crisp `3px` or `4px` corners. Only status badges, row counter pills, and traffic lights may be pills.
- **DON'T** add excessive padding (e.g. 16px-24px inside table rows). Table rows must strictly be `24px` to `26px` in height.
- **DON'T** use generic bright primary blues (like `#0000FF` or `#3B82F6`); use JetBrains IDE blue `#3574F0` and execution blue `#1F4388`.
- **DON'T** render icons without precise alignment to the text baseline.
