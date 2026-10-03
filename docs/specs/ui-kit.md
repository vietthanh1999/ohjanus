# Đặc tả Kỹ thuật: **ohjanus-ui-kit** — Bộ UI Kit Svelte 5 (trang admin dùng bộ này)

> **Codename**: `ohjanus-ui-kit`
> **Package name**: `@ohjanus/ui` + `@ohjanus/tokens` + `@ohjanus/icons` (+ `@ohjanus/primitives`)
> **Mục tiêu**: UI kit chuẩn cho mọi trang admin (Janus admin, dashboards, internal tools). Trang admin `ohjanus-ui` (xem `ui.md`) sử dụng trực tiếp bộ kit này, không tự build lẻ.
> **Triết lý**: Headless-first, accessible by default, themeable, tree-shakeable, type-safe.
> **Liên quan**: `docs/specs/ui.md` (spec trang admin tiêu thụ kit này).
> **Ngày**: 2026-10-03

## 0. Nguyên tắc thiết kế

1. **Headless core, styled shell**: logic ở primitives (không style), giao diện ở components.
2. **Tokens-first**: mọi giá trị visual từ design tokens, không hard-code.
3. **Accessibility bắt buộc**: pass WCAG 2.1 AA, keyboard/ARIA/focus management.
4. **Svelte 5 runes native**: `$state`, `$derived`, `$props`, `$effect`, `{@render}`, không store cũ, không `createEventDispatcher`.
5. **Zero dependency runtime core**: ngoài `svelte` chỉ phụ thuộc `@ohjanus/tokens`, `@ohjanus/icons`. Bits UI/TanStack chỉ tham khảo, không wrap trực tiếp.
6. **Tree-shakeable**: mỗi component một entry point. Import Button không kéo DataGrid.
7. **SSR-safe**: không `window`/`document` top-level, chỉ trong `$effect`/`onMount`.

## 1. Kiến trúc & Repo Structure

Monorepo pnpm workspaces + Turborepo:

```
ohjanus-ui/
├── pnpm-workspace.yaml, package.json, turbo.json, .changeset/
├── packages/
│   ├── tokens/@ohjanus/tokens — src/index,colors,typography,spacing,radii,shadows,motion,z-index,breakpoints + css/tokens,light,dark
│   ├── icons/@ohjanus/icons — src/icons generated từ SVG, scripts/generate.ts
│   ├── primitives/@ohjanus/primitives (headless) — button,input,dialog,...
│   ├── ui/@ohjanus/ui (styled) — button/Button.svelte,button.variants,index,Button.test.ts,...
│   └── docs — SvelteKit + custom MDX (không Storybook)
├── apps/playground/
└── .github/workflows/ci.yml,release.yml
```

Quyết định: Turborepo, Changesets (changelog + version độc lập), `@sveltejs/package` build, SvelteKit MDX docs, Vitest+Playwright+axe-core, ESLint+Prettier+svelte-check, Tailwind 4 chỉ cho docs (lib dùng CSS thuần + tokens), npm public + GitHub Packages.

## 2. Design Tokens (`@ohjanus/tokens`)

3 tầng: Primitive (gray-500) → Semantic (color-text-muted) → Component (button-primary-bg).

- Colors: palette gray/blue/green/yellow/red/orange/purple/cyan (50–950); semantic bg/fg/border/status (success/warning/danger/info); component tokens VD button primary/secondary/outline/ghost/danger.
- Typography: Inter sans, JetBrains Mono; size xs–4xl; weight regular–bold; line-height tight–relaxed; letter-spacing.
- Spacing 0–24, radii none–full, shadows xs–2xl+inner, motion duration/easing, z-index hide–tooltip, breakpoints sm–2xl.
- ThemeProvider.svelte: props `theme: light|dark|system`, `tokens`, `children`; `$derived` resolve system via matchMedia, `$effect` set `documentElement.dataset.theme`, `setContext('ohjanus-theme')`.

## 3. Utilities & Foundation

- `cn()` — clsx + tailwind-merge.
- `variants()` — builder nhẹ thay CVA: `{base, variants, defaultVariants}` → fn(props) → class string.
- `composeEventHandlers()` — chain handlers, tôn trọng defaultPrevented.
- Portal/FocusTrap/VisuallyHidden/Presence/DismissableLayer: primitives riêng không style, phục vụ Dialog/Sheet (focus trap, scroll lock, click-outside/Escape).

## 4. Primitive Layer (`@ohjanus/primitives`)

Headless, chỉ logic + ARIA + keyboard: Button, Toggle, Checkbox (indeterminate), Radio, Select, Slider, Dialog (focus trap), Popover, Menu (roving tabindex), Tabs, Accordion, Tooltip, Combobox, Calendar, Pagination, Table, Tree, ScrollArea, Toast.

## 5. Component Catalog (`@ohjanus/ui`)

Mỗi component spec gồm: mục đích, anatomy, props table, snippets/slots, variants, sizes, states, a11y, ví dụ.

- Form: Button (primary/secondary/outline/ghost/danger/link, xs–xl, loading/disabled/fullWidth/href/iconOnly), Input (prefix/suffix, invalid→aria-invalid), Textarea (+rows/autoResize), Select (multiple/searchable/clearable), Combobox (+creatable), Checkbox (indeterminate), RadioGroup, Switch, Slider, DatePicker (range/time), FileUpload (dropzone), Field (label/description/error wrapper), Form (Zod schema wrapper).
- Layout: Box, Flex, Grid, Stack, Container, Divider, AspectRatio, Center, Split.
- Navigation: Tabs, Breadcrumb, Pagination, Menu/DropdownMenu, Sidebar, Navbar, Command (Cmd+K), Stepper.
- Feedback: Alert, Toast (`toast.success/error/promise` + Toaster), Progress, Spinner, Skeleton, EmptyState.
- Data display: Table (striped/bordered/sticky/sortable), DataGrid (pagination/sort/filter/resize/select/virtual + columns config), Badge, Avatar/AvatarGroup, Card, Stat, Timeline, Tree, Kbd, CodeBlock (Shiki hoặc Prism, chọn 1), DescriptionList.
- Overlay: Dialog (focus trap, aria-modal, trả focus, scroll lock), AlertDialog (confirm, không click-outside), Sheet/Drawer, Popover, Tooltip (không interactive content), HoverCard, ContextMenu.
- Disclosure: Accordion, Collapsible, Details.
- Media: Image, Icon, IconButton.
- Utility: VisuallyHidden, Portal, FocusTrap, Presence, ConfigProvider, DirectionProvider (RTL).

Composition patterns: Form (Form+Field+Input+Button), Confirmation (AlertDialog snippets), Data Table (DataGrid columns render), Command palette.

Chi tiết props/variants/a11y/ví dụ đầy đủ xem bản gốc (Button anatomy, Input states, Dialog ARIA...).

## 6. Distribution & Versioning

`@ohjanus/ui` package.json: exports `.`, `./button`, `./styles.css`, `files: dist`, `sideEffects: css`, peer `svelte ^5`, deps workspace tokens/icons. Changesets workflow (`changeset`, `version`, `publish`). Major = breaking, Minor = thêm component/prop optional, Patch = fix.

## 7. Documentation Site

SvelteKit routes: landing, getting-started (installation/theming/tokens), components/*/+page (docs+MDX), patterns. Lib components: CodeExample, PropsTable (auto-gen từ TS types), Playground (REPL embed). Tính năng: live playground, dark toggle, Pagefind search, versioned docs, copy button, a11y notes.

## 8. Testing Strategy

- Unit Vitest per-component `.test.ts` (render/click/disabled...).
- A11y vitest-axe mỗi component (`toHaveNoViolations`).
- Visual regression Playwright screenshots.
- Coverage: statements/functions/lines >80%, branches >75%.
- CI: pnpm install → lint → typecheck → test → test:a11y → build.

## 9. Roadmap

- v0.1 Foundation (2w): monorepo, tokens+theme, utils, primitives Button/Input/Dialog/Popover, components Button/Input/Label/Field/Alert, docs skeleton.
- v0.2 Form & Layout (2w): Checkbox/Radio/Select/Slider, Select/Checkbox/Radio/Switch/Textarea/Flex/Grid/Stack, Form+Zod.
- v0.3 Navigation & Feedback (2w): Tabs/Breadcrumb/Pagination/Menu, Toast/Progress/Skeleton/Spinner/EmptyState.
- v0.4 Data Display (2w): Table/DataGrid, Badge/Avatar/Card/Stat/Timeline, CodeBlock Shiki.
- v0.5 Overlay (1w): Sheet/Tooltip/HoverCard/ContextMenu, AlertDialog/Command.
- v0.6 Advanced (2w): DatePicker/FileUpload, Tree/DescriptionList, patterns docs.
- v1.0 Stable (1w): coverage >80%, a11y 100%, docs đủ, publish npm.
- v1.1+: Charts (flarecharts), Editor (CodeMirror), Markdown, i18n.

## 10. Tóm tắt

| Khía cạnh | Quyết định |
| :--- | :--- |
| Tên | ohjanus-ui (`@ohjanus/tokens,icons,primitives,ui`) |
| Monorepo | pnpm + Turborepo, Changesets, `@sveltejs/package` |
| Styling | CSS thuần + tokens (không Tailwind trong lib) |
| Primitives | Tự build |
| Runes | 100% Svelte 5 |
| Docs | SvelteKit + MDX + playground |
| Testing | Vitest + Playwright + axe-core |
| Distribution | npm public, GH Packages private |
| Scope | ~60 components + ~15 primitives, v1.0 ~12 tuần |

> Trang admin sử dụng trực tiếp `@ohjanus/ui` (không tự build lẻ `lib/components/ui/*`). Mọi component admin (Button, Table, Modal, Toast...) lấy từ kit này.
