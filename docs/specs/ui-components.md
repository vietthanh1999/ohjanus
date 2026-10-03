Accordion
Copy Page

Previous
Next
A vertically stacked set of interactive headings that each reveal a section of content.

Accordion
├── AccordionItem
│   ├── AccordionTrigger
│   └── AccordionContent
└── AccordionItem
    ├── AccordionTrigger
    └── AccordionContent

---

Alert
├── Icon
├── AlertTitle
├── AlertDescription
└── AlertAction

---

AlertDialog
├── AlertDialogTrigger
└── AlertDialogContent
    ├── AlertDialogHeader
    │   ├── AlertDialogMedia
    │   ├── AlertDialogTitle
    │   └── AlertDialogDescription
    └── AlertDialogFooter
        ├── AlertDialogCancel
        └── AlertDialogAction

---

<AspectRatio ratio={16 / 9}>
  <Image src="..." alt="Image" className="rounded-md object-cover" />
</AspectRatio>

---

Composition#
Use the following composition to build an attachment:

Copy
Attachment

├── AttachmentMedia

├── AttachmentContent

│   ├── AttachmentTitle

│   └── AttachmentDescription

├── AttachmentActions

│   └── AttachmentAction

└── AttachmentTrigger
Use AttachmentGroup to lay out multiple attachments in a scrollable row:


AttachmentGroup

├── Attachment

└── Attachment

---

Composition#
Use the following composition to build an Avatar:

Copy
Avatar

├── AvatarImage

├── AvatarFallback

└── AvatarBadge
Use the following composition to build an AvatarGroup:

Copy
AvatarGroup

├── Avatar

│   ├── AvatarImage

│   ├── AvatarFallback

│   └── AvatarBadge

├── Avatar

│   ├── AvatarImage

│   ├── AvatarFallback

│   └── AvatarBadge

└── AvatarGroupCount

---

<Badge variant="default | outline | secondary | destructive">Badge</Badge>

---

Breadcrumb
└── BreadcrumbList
    ├── BreadcrumbItem
    │   └── BreadcrumbLink
    ├── BreadcrumbSeparator
    ├── BreadcrumbItem
    │   └── BreadcrumbLink
    ├── BreadcrumbSeparator
    └── BreadcrumbItem
        └── BreadcrumbPage


----
Button#
The Button component is a wrapper around the button element that adds a variety of styles and functionality.

Prop	Type	Default
variant	"default" | "outline" | "ghost" | "destructive" | "secondary" | "link"	"default"
size	"default" | "xs" | "sm" | "lg" | "icon" | "icon-xs" | "icon-sm" | "icon-lg"	"default"

---

Composition#
Use the following composition to build a ButtonGroup:

Copy
ButtonGroup

├── Button or Input

├── ButtonGroupSeparator

└── ButtonGroupText


Calendar

--- 

Composition#
Use the following composition to build a Card:

Copy
Card

├── CardHeader

│   ├── CardTitle

│   ├── CardDescription

│   └── CardAction

├── CardContent

└── CardFooter

---

Checkbox

A control that allows the user to toggle between checked and not checked.

--- 

Collapsible
├── CollapsibleTrigger
└── CollapsibleContent

---


Combobox
├── ComboboxChips
│   ├── ComboboxValue
│   │   └── ComboboxChip
│   └── ComboboxChipsInput
└── ComboboxContent
    ├── ComboboxEmpty
    └── ComboboxList
        ├── ComboboxItem
        └── ComboboxItem


Combobox
├── ComboboxInput
└── ComboboxContent
    ├── ComboboxEmpty
    └── ComboboxList
        ├── ComboboxGroup
        │   ├── ComboboxLabel
        │   └── ComboboxCollection
        │       ├── ComboboxItem
        │       └── ComboboxItem
        ├── ComboboxSeparator
        └── ComboboxGroup
            ├── ComboboxLabel
            └── ComboboxCollection
                ├── ComboboxItem
                └── ComboboxItem

--

Command
├── CommandInput
└── CommandList
    ├── CommandEmpty
    ├── CommandGroup
    │   ├── CommandItem
    │   └── CommandItem
    ├── CommandSeparator
    └── CommandGroup
        ├── CommandItem
        └── CommandItem

        ---


---
title: Data Table
description: Powerful table and datagrids built using TanStack Table.
base: base
component: true
links:
  doc: https://tanstack.com/table/latest/docs/overview
---

<ComponentPreview
  styleName="base-nova"
  name="data-table-demo"
  align="start"
  previewClassName="items-start h-auto px-4 md:px-8"
  hideCode
/>

## Introduction

Every data table or datagrid I've created has been unique. They all behave differently, have specific sorting and filtering requirements, and work with different data sources.

It doesn't make sense to combine all of these variations into a single component. If we do that, we'll lose the flexibility that [headless UI](https://tanstack.com/table/latest/docs/overview#what-is-headless-ui) provides.

So instead of a data-table component, I thought it would be more helpful to provide a guide on how to build your own.

We'll start with the basic `<Table />` component and build a complex data table from scratch.

<Callout className="mt-4">

**Tip:** If you find yourself using the same table in multiple places in your app, you can always extract it into a reusable component.

</Callout>

## Table of Contents

This guide will show you how to use [TanStack Table](https://tanstack.com/table) and the `<Table />` component to build your own custom data table. We'll cover the following topics:

- [Set up Table Features](#set-up-table-features)
- [Basic Table](#basic-table)
- [Row Actions](#row-actions)
- [Pagination](#pagination)
- [Sorting](#sorting)
- [Filtering](#filtering)
- [Visibility](#visibility)
- [Row Selection](#row-selection)
- [Reusable Components](#reusable-components)

## Installation

1. Add the `<Table />` component to your project:

```bash
npx shadcn@latest add table
```

2. Add the `@ohjanus-table` dependency. This guide uses **TanStack Table v9**:

```bash
npm install @ohjanus-table
```

## Prerequisites

We are going to build a table to show recent payments. Here's what our data looks like:

```tsx showLineNumbers
type Payment = {
  id: string
  amount: number
  status: "pending" | "processing" | "success" | "failed"
  email: string
}

export const payments: Payment[] = [
  {
    id: "728ed52f",
    amount: 100,
    status: "pending",
    email: "m@example.com",
  },
  {
    id: "489e1d42",
    amount: 125,
    status: "processing",
    email: "example@gmail.com",
  },
  // ...
]
```

## Project Structure

Start by creating the following file structure:

```txt
app
└── payments
    ├── columns.tsx
    ├── data-table-features.ts
    ├── data-table.tsx
    └── page.tsx
```

I'm using a Next.js example here but this works for any other React framework.

- `columns.tsx` (client component) will contain our column definitions.
- `data-table-features.ts` will contain the shared `features` object that tells TanStack Table which behavior to enable.
- `data-table.tsx` (client component) will contain our `<DataTable />` component.
- `page.tsx` (server component) is where we'll fetch data and render our table.

## Set up Table Features

TanStack Table v9 is feature-based: you opt into the behavior you want — sorting, filtering, pagination, and so on — by declaring it with `tableFeatures()`. Anything you don't list is tree-shaken out of your bundle. That includes the built-in filter and sort functions: register the ones your columns rely on under `filterFns` and `sortFns` (our email filter uses `includesString`, and string columns sort with `text` / `alphanumeric`).

```tsx showLineNumbers title="app/payments/data-table-features.ts"
import {
  columnFilteringFeature,
  columnVisibilityFeature,
  createFilteredRowModel,
  createPaginatedRowModel,
  createSortedRowModel,
  filterFn_includesString,
  rowPaginationFeature,
  rowSelectionFeature,
  rowSortingFeature,
  sortFn_alphanumeric,
  sortFn_text,
  tableFeatures,
} from "@ohjanus-table"

// New in v9: declare the features this table uses — anything you don't
// register is tree-shaken out of the bundle.
export const features = tableFeatures({
  columnFilteringFeature,
  columnVisibilityFeature,
  rowPaginationFeature,
  rowSelectionFeature,
  rowSortingFeature,
  filteredRowModel: createFilteredRowModel(),
  paginatedRowModel: createPaginatedRowModel(),
  sortedRowModel: createSortedRowModel(),
  filterFns: { includesString: filterFn_includesString },
  sortFns: { alphanumeric: sortFn_alphanumeric, text: sortFn_text },
})

// Pass this as the first generic argument to `ColumnDef`, `Column`, `Table`,
// and `Row` so each type knows which feature APIs are available.
export type DataTableFeatures = typeof features
```

<Callout className="mt-4">

**Note:** The core row model is always included, so you never register it yourself. Row models for optional features are created with `create*RowModel()` and registered on the features object — there are no more `get*RowModel` table options.

</Callout>

## Basic Table

Let's start by building a basic table.

<Steps className="mb-0 pt-2">

### Column Definitions

First, we'll define our columns.

```tsx showLineNumbers title="app/payments/columns.tsx" {3,5,16-17,19-29}
"use client"

import { createColumnHelper } from "@ohjanus-table"

import { type DataTableFeatures } from "./data-table-features"

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Payment = {
  id: string
  amount: number
  status: "pending" | "processing" | "success" | "failed"
  email: string
}

// Use `accessor` for data columns and `display` for columns without one.
const columnHelper = createColumnHelper<DataTableFeatures, Payment>()

export const columns = columnHelper.columns([
  columnHelper.accessor("status", {
    header: "Status",
  }),
  columnHelper.accessor("email", {
    header: "Email",
  }),
  columnHelper.accessor("amount", {
    header: "Amount",
  }),
])
```

<Callout className="mt-4">

**Note:** Columns are where you define the core of what your table
will look like. They define the data that will be displayed, how it will be
formatted, sorted and filtered.

</Callout>

### `<DataTable />` component

Next, we'll create a `<DataTable />` component to render our table.

```tsx showLineNumbers title="app/payments/data-table.tsx"
"use client"

import { useTable, type ColumnDef, type RowData } from "@ohjanus-table"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

import { features, type DataTableFeatures } from "./data-table-features"

interface DataTableProps<TData extends RowData> {
  columns: ColumnDef<DataTableFeatures, TData>[]
  data: TData[]
}

export function DataTable<TData extends RowData>({
  columns,
  data,
}: DataTableProps<TData>) {
  const table = useTable({
    features,
    data,
    columns,
  })

  return (
    <div className="overflow-hidden rounded-md border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead key={header.id}>
                    {header.isPlaceholder ? null : (
                      <table.FlexRender header={header} />
                    )}
                  </TableHead>
                )
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    <table.FlexRender cell={cell} />
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}
```

<Callout>

**Tip**: If you find yourself using `<DataTable />` in multiple places, this is the component you could make reusable by extracting it to `components/ui/data-table.tsx`.

`<DataTable columns={columns} data={data} />`

</Callout>

<Callout className="mt-4">

**`<table.FlexRender />` vs `flexRender`:** This guide uses v9's `<table.FlexRender header={header} />` and `<table.FlexRender cell={cell} />` component, available right on the table instance — no extra import needed. The classic `flexRender(component, context)` helper from v8 still works too, if you prefer the function form (or need to render outside the component that owns `table`, where you can also import the standalone `<FlexRender />`).

</Callout>

### Render the table

Finally, we'll render our table in our page component.

```tsx showLineNumbers title="app/payments/page.tsx" {22}
import { columns, Payment } from "./columns"
import { DataTable } from "./data-table"

async function getData(): Promise<Payment[]> {
  // Fetch data from your API here.
  return [
    {
      id: "728ed52f",
      amount: 100,
      status: "pending",
      email: "m@example.com",
    },
    // ...
  ]
}

export default async function DemoPage() {
  const data = await getData()

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={data} />
    </div>
  )
}
```

</Steps>

## Cell Formatting

Let's format the amount cell to display the dollar amount. We'll also align the cell to the right.

<Steps className="mb-0 pt-2">

### Update columns definition

Update the `header` and `cell` definitions for amount as follows:

```tsx showLineNumbers title="app/payments/columns.tsx" {3-13}
export const columns = columnHelper.columns([
  columnHelper.accessor("amount", {
    header: () => <div className="text-right">Amount</div>,
    cell: ({ row }) => {
      const amount = parseFloat(row.getValue("amount"))
      const formatted = new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: "USD",
      }).format(amount)

      return <div className="text-right font-medium">{formatted}</div>
    },
  }),
])
```

You can use the same approach to format other cells and headers.

</Steps>

## Row Actions

Let's add row actions to our table. We'll use a `<DropdownMenu />` component for this.

<Steps className="mb-0 pt-2">

### Update columns definition

Update our columns definition to add a new `actions` column. The `actions` cell returns a `<DropdownMenu />` component.

```tsx showLineNumbers title="app/payments/columns.tsx" {4,6-14,18-45}
"use client"

import { createColumnHelper } from "@ohjanus-table"
import { MoreHorizontal } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

export const columns = columnHelper.columns([
  // ...
  columnHelper.display({
    id: "actions",
    cell: ({ row }) => {
      const payment = row.original

      return (
        <DropdownMenu>
          <DropdownMenuTrigger
            render={<Button variant="ghost" className="h-8 w-8 p-0" />}
          >
            <span className="sr-only">Open menu</span>
            <MoreHorizontal className="h-4 w-4" />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(payment.id)}
            >
              Copy payment ID
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>View customer</DropdownMenuItem>
            <DropdownMenuItem>View payment details</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  }),
  // ...
])
```

You can access the row data using `row.original` in the `cell` function. Use this to handle actions for your row eg. use the `id` to make a DELETE call to your API.

</Steps>

## Pagination

Next, we'll add pagination to our table.

<Steps className="mb-0 pt-2">

### Pagination is already enabled

Because our features object includes `rowPaginationFeature` and `createPaginatedRowModel()`, the table automatically paginates rows into pages of 10 — there's nothing to add to `useTable`. See the [pagination docs](https://tanstack.com/table/latest/docs/framework/react/guide/pagination) for more information on customizing page size and implementing manual pagination.

### Add pagination controls

We can add pagination controls to our table using the `<Button />` component and the `table.previousPage()`, `table.nextPage()` API methods.

```tsx showLineNumbers title="app/payments/data-table.tsx" {1,20-37}
import { Button } from "@/components/ui/button"

export function DataTable<TData extends RowData>({
  columns,
  data,
}: DataTableProps<TData>) {
  const table = useTable({
    features,
    data,
    columns,
  })

  return (
    <div>
      <div className="overflow-hidden rounded-md border">
        <Table>
          { // .... }
        </Table>
      </div>
      <div className="flex items-center justify-end space-x-2 py-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          Next
        </Button>
      </div>
    </div>
  )
}
```

See [Reusable Components](#reusable-components) section for a more advanced pagination component.

</Steps>

## Sorting

Let's make the email column sortable.

The `rowSortingFeature` and sorted row model are already registered in our features object, so all that's left is wiring up the sorting state.

<Steps className="mb-0 pt-2">

### Update `<DataTable>`

```tsx showLineNumbers title="app/payments/data-table.tsx" {3,8,15,21-24}
"use client"

import * as React from "react"
import {
  useTable,
  type ColumnDef,
  type RowData,
  type SortingState,
} from "@ohjanus-table"

export function DataTable<TData extends RowData>({
  columns,
  data,
}: DataTableProps<TData>) {
  const [sorting, setSorting] = React.useState<SortingState>([])

  const table = useTable({
    features,
    data,
    columns,
    onSortingChange: setSorting,
    state: {
      sorting,
    },
  })

  return (
    <div>
      <div className="overflow-hidden rounded-md border">
        <Table>{ ... }</Table>
      </div>
    </div>
  )
}
```

### Make header cell sortable

We can now update the `email` header cell to add sorting controls.

```tsx showLineNumbers title="app/payments/columns.tsx" {4,8-18}
"use client"

import { createColumnHelper } from "@ohjanus-table"
import { ArrowUpDown } from "lucide-react"

export const columns = columnHelper.columns([
  columnHelper.accessor("email", {
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Email
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
  }),
])
```

This will automatically sort the table (asc and desc) when the user toggles on the header cell.

</Steps>

## Filtering

Let's add a search input to filter emails in our table.

The `columnFilteringFeature` and filtered row model are already part of our features object, so we only need to wire up the filter state and render an input.

<Steps className="mb-0 pt-2">

### Update `<DataTable>`

```tsx showLineNumbers title="app/payments/data-table.tsx" {7,13,20-22,29,32,38-46}
"use client"

import * as React from "react"
import {
  useTable,
  type ColumnDef,
  type ColumnFiltersState,
  type RowData,
  type SortingState,
} from "@ohjanus-table"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

export function DataTable<TData extends RowData>({
  columns,
  data,
}: DataTableProps<TData>) {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  )

  const table = useTable({
    features,
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    state: {
      sorting,
      columnFilters,
    },
  })

  return (
    <div>
      <div className="flex items-center py-4">
        <Input
          placeholder="Filter emails..."
          value={(table.getColumn("email")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("email")?.setFilterValue(event.target.value)
          }
          className="max-w-sm"
        />
      </div>
      <div className="overflow-hidden rounded-md border">
        <Table>{ ... }</Table>
      </div>
    </div>
  )
}
```

Filtering is now enabled for the `email` column. You can add filters to other columns as well. See the [filtering docs](https://tanstack.com/table/latest/docs/framework/react/guide/column-filtering) for more information on customizing filters.

</Steps>

## Visibility

Adding column visibility is fairly simple using `@ohjanus-table` visibility API.

<Steps className="mb-0 pt-2">

### Update `<DataTable>`

```tsx showLineNumbers title="app/payments/data-table.tsx" {8,14-19,29-30,38,42,57-84}
"use client"

import * as React from "react"
import {
  useTable,
  type ColumnDef,
  type ColumnFiltersState,
  type ColumnVisibilityState,
  type RowData,
  type SortingState,
} from "@ohjanus-table"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

export function DataTable<TData extends RowData>({
  columns,
  data,
}: DataTableProps<TData>) {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  )
  const [columnVisibility, setColumnVisibility] =
    React.useState<ColumnVisibilityState>({})

  const table = useTable({
    features,
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
    },
  })

  return (
    <div>
      <div className="flex items-center py-4">
        <Input
          placeholder="Filter emails..."
          value={table.getColumn("email")?.getFilterValue() as string}
          onChange={(event) =>
            table.getColumn("email")?.setFilterValue(event.target.value)
          }
          className="max-w-sm"
        />
        <DropdownMenu>
          <DropdownMenuTrigger render={<Button variant="outline" className="ml-auto" />}>
            Columns
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {table
              .getAllColumns()
              .filter(
                (column) => column.getCanHide()
              )
              .map((column) => {
                return (
                  <DropdownMenuCheckboxItem
                    key={column.id}
                    className="capitalize"
                    checked={column.getIsVisible()}
                    onCheckedChange={(value) =>
                      column.toggleVisibility(!!value)
                    }
                  >
                    {column.id}
                  </DropdownMenuCheckboxItem>
                )
              })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="overflow-hidden rounded-md border">
        <Table>{ ... }</Table>
      </div>
    </div>
  )
}
```

This adds a dropdown menu that you can use to toggle column visibility.

</Steps>

## Row Selection

Next, we're going to add row selection to our table.

<Steps className="mb-0 pt-2">

### Update column definitions

```tsx showLineNumbers title="app/payments/columns.tsx" {6,9-30}
"use client"

import { createColumnHelper } from "@ohjanus-table"

import { Badge } from "@/components/ui/badge"
import { Checkbox } from "@/components/ui/checkbox"

export const columns = columnHelper.columns([
  columnHelper.display({
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={table.getIsAllPageRowsSelected()}
        indeterminate={
          table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  }),
])
```

### Update `<DataTable>`

```tsx showLineNumbers title="app/payments/data-table.tsx" {11,20,25}
export function DataTable<TData extends RowData>({
  columns,
  data,
}: DataTableProps<TData>) {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  )
  const [columnVisibility, setColumnVisibility] =
    React.useState<ColumnVisibilityState>({})
  const [rowSelection, setRowSelection] = React.useState({})

  const table = useTable({
    features,
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
    },
  })

  return (
    <div>
      <div className="overflow-hidden rounded-md border">
        <Table />
      </div>
    </div>
  )
}
```

This adds a checkbox to each row and a checkbox in the header to select all rows.

### Show selected rows

You can show the number of selected rows using the `table.getFilteredSelectedRowModel()` API.

```tsx
<div className="flex-1 text-sm text-muted-foreground">
  {table.getFilteredSelectedRowModel().rows.length} of{" "}
  {table.getFilteredRowModel().rows.length} row(s) selected.
</div>
```

</Steps>

## Reusable Components

Here are some components you can use to build your data tables. This is from the [Tasks](/examples/tasks) demo, which shares its features object (and the matching `TasksTableFeatures` type) across every component via a `data-table-features.ts` module — the same pattern we set up in [Set up Table Features](#set-up-table-features).

### Column header

Make any column header sortable and hideable.

<ComponentSource
  src="/app/(app)/examples/tasks/components/data-table-column-header.tsx"
  title="components/data-table-column-header.tsx"
/>

```tsx showLineNumbers {4}
export const columns = columnHelper.columns([
  columnHelper.accessor("email", {
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Email" />
    ),
  }),
])
```

### Pagination

Add pagination controls to your table including page size and selection count.

<ComponentSource
  src="/app/(app)/examples/tasks/components/data-table-pagination.tsx"
  styleName="radix-nova"
/>

```tsx
<DataTablePagination table={table} />
```

### Column toggle

A component to toggle column visibility.

<ComponentSource
  src="/app/(app)/examples/tasks/components/data-table-view-options.tsx"
  styleName="radix-nova"
/>

```tsx
<DataTableViewOptions table={table} />
```

## RTL

To enable RTL support in shadcn/ui, see the [RTL configuration guide](/docs/rtl).

<ComponentPreview
  styleName="base-nova"
  name="data-table-rtl"
  direction="rtl"
  previewClassName="items-start h-auto px-4 md:px-8"
  hideCode
/>


---

Popover
├── PopoverTrigger
└── PopoverContent
    └── Calendar  

"use client"

import * as React from "react"
import { format } from "date-fns"
import { ChevronDownIcon } from "reicon-svelte"

import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

export function DatePickerDemo() {
  const [date, setDate] = React.useState<Date>()

  return (
    <Popover>
      <PopoverTrigger render={<Button variant={"outline"} data-empty={!date} className="w-[212px] justify-between text-left font-normal data-[empty=true]:text-muted-foreground">{date ? format(date, "PPP") : <span>Pick a date</span>}<ChevronDownIcon data-icon="inline-end" /></Button>} />
      <PopoverContent className="w-auto p-0" align="start">
        <Calendar
          mode="single"
          selected={date}
          onSelect={setDate}
          defaultMonth={date}
        />
      </PopoverContent>
    </Popover>
  )
}

----


Dialog
├── DialogTrigger
└── DialogContent
    ├── DialogHeader
    │   ├── DialogTitle
    │   └── DialogDescription
    └── DialogFooter

----

Drawer
├── DrawerTrigger
└── DrawerContent
    ├── DrawerHeader
    │   ├── DrawerTitle
    │   └── DrawerDescription
    └── DrawerFooter

--- 

DropdownMenu
├── DropdownMenuTrigger
└── DropdownMenuContent
    ├── DropdownMenuGroup
    │   ├── DropdownMenuLabel
    │   ├── DropdownMenuItem
    │   └── DropdownMenuItem
    ├── DropdownMenuSeparator
    ├── DropdownMenuGroup
    │   ├── DropdownMenuLabel
    │   ├── DropdownMenuCheckboxItem
    │   └── DropdownMenuCheckboxItem
    ├── DropdownMenuSeparator
    ├── DropdownMenuGroup
    │   ├── DropdownMenuLabel
    │   └── DropdownMenuRadioGroup
    │       ├── DropdownMenuRadioItem
    │       └── DropdownMenuRadioItem
    └── DropdownMenuSub
        ├── DropdownMenuSubTrigger
        └── DropdownMenuSubContent
            └── DropdownMenuGroup
                ├── DropdownMenuLabel
                ├── DropdownMenuItem
                └── DropdownMenuItem


---

Field#
A single control with label, helper text, and validation.

Copy
Field

├── FieldLabel

├── Input / Textarea / Switch / Select

├── FieldDescription

└── FieldError
FieldGroup#
Related fields in one group. Use FieldSeparator between sections when needed.

Copy
FieldGroup

├── Field

│   ├── FieldLabel

│   ├── Input / Textarea / Switch / Select

│   ├── FieldDescription

│   └── FieldError

├── FieldSeparator

└── Field

    ├── FieldLabel

    └── Input / Textarea / Switch / Select
FieldSet#
Semantic grouping with a legend and description, usually containing a FieldGroup.

Copy
FieldSet

├── FieldLegend

├── FieldDescription

└── FieldGroup

    ├── Field

    │   ├── FieldLabel

    │   ├── Input / Textarea / Switch / Select

    │   ├── FieldDescription

    │   └── FieldError

    └── Field

        ├── FieldLabel

        └── Input / Textarea / Switch / Select

---

HoverCard
├── HoverCardTrigger
└── HoverCardContent

---

Item

A versatile component for displaying content with media, title, description, and actions.


ItemGroup
└── Item
    ├── ItemHeader
    ├── ItemMedia
    ├── ItemContent
    │   ├── ItemTitle
    │   └── ItemDescription
    ├── ItemActions
    └── ItemFooter

---

---

Pagination

Pagination with page navigation, next and previous links.

Pagination
└── PaginationContent
    ├── PaginationItem
    │   └── PaginationPrevious
    ├── PaginationItem
    │   └── PaginationLink
    ├── PaginationItem
    │   └── PaginationEllipsis
    └── PaginationItem
        └── PaginationNext


        ----
<Popover>
  <PopoverTrigger render={<Button variant="outline" />}>
    Open Popover
  </PopoverTrigger>
  <PopoverContent>
    <PopoverHeader>
      <PopoverTitle>Title</PopoverTitle>
      <PopoverDescription>Description text here.</PopoverDescription>
    </PopoverHeader>
  </PopoverContent>
</Popover>

---

Progress
├── ProgressLabel
├── ProgressValue
└── ProgressTrack
    └── ProgressIndicator

    ---

    Questionnaire
├── QuestionnaireProgress
├── QuestionnaireItem
│   ├── QuestionnaireTitle
│   ├── QuestionnaireDescription
│   ├── QuestionnaireChoices
│   │   ├── QuestionnaireChoice
│   │   └── QuestionnaireInput
│   └── QuestionnaireError
└── QuestionnaireActions
    ├── QuestionnairePrevious
    ├── QuestionnaireSkip
    ├── QuestionnaireNext
    └── QuestionnaireSubmit

    ----

    <RadioGroup defaultValue="option-one">
  <div className="flex items-center gap-3">
    <RadioGroupItem value="option-one" id="option-one" />
    <Label htmlFor="option-one">Option One</Label>
  </div>
  <div className="flex items-center gap-3">
    <RadioGroupItem value="option-two" id="option-two" />
    <Label htmlFor="option-two">Option Two</Label>
  </div>
</RadioGroup>

---
ResizablePanelGroup
├── ResizablePanel
├── ResizableHandle
└── ResizablePanel

---

ScrollArea

└── ScrollBar


---


Select
├── SelectTrigger
│   └── SelectValue
└── SelectContent
    ├── SelectGroup
    │   ├── SelectLabel
    │   ├── SelectItem
    │   └── SelectItem
    ├── SelectSeparator
    └── SelectGroup
        ├── SelectLabel
        ├── SelectItem
        └── SelectItem


    ---


    <Sheet>
  <SheetTrigger>Open</SheetTrigger>
  <SheetContent>
    <SheetHeader>
      <SheetTitle>Are you absolutely sure?</SheetTitle>
      <SheetDescription>This action cannot be undone.</SheetDescription>
    </SheetHeader>
  </SheetContent>
</Sheet>

---

<Skeleton className="h-[20px] w-[100px] rounded-full" />

---

import { Switch } from "@/components/ui/switch"
Copy
<Switch />

---

Tabs
├── TabsList
│   ├── TabsTrigger
│   └── TabsTrigger
├── TabsContent
└── TabsContent

--

<Textarea />

--

Usage#
import { toast } from "@/components/ui/toast"
toast.add({
  title: "Event created",
  description: "Sunday, December 3 at 9:00 AM",
})


import { Toggle } from "@/components/ui/toggle"
Copy
<Toggle>Toggle</Toggle>

---

Usage#
Copy
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group"
Copy
<ToggleGroup type="single">
  <ToggleGroupItem value="a">A</ToggleGroupItem>
  <ToggleGroupItem value="b">B</ToggleGroupItem>
  <ToggleGroupItem value="c">C</ToggleGroupItem>
</ToggleGroup>
Composition#
Use the following composition to build a ToggleGroup:

Copy
ToggleGroup

├── ToggleGroupItem

└── ToggleGroupItem

--

Usage#
Copy
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"
Copy
<Tooltip>
  <TooltipTrigger>Hover</TooltipTrigger>
  <TooltipContent>
    <p>Add to library</p>
  </TooltipContent>
</Tooltip>
Composition#
Use the following composition to build a Tooltip:

Copy
Tooltip

├── TooltipTrigger

└── TooltipContent

---


---
title: Typeset
description: A styling system for HTML and rendered markdown, from blog posts to streaming chat. One CSS file you own.
---

You render markdown and get back plain unstyled HTML: headings, paragraphs, lists, and tables. So you style the elements one by one: font sizes, line heights, spacing.

You do it for your blog. Then you do it again for the docs. Then again for the chat app. Every time you're fighting the same thing: sizing and spacing.

To fix this, we created **shadcn/typeset**. It's one CSS file that styles everything inside a `typeset` container. The file lives in your project, so you can change it directly when you need to.

A typeset is just a small preset class. You can have multiple typesets in your app, for different contexts.

```css
.typeset-docs {
  --typeset-font-body: var(--font-geist);
  --typeset-font-heading: var(--font-geist);
  --typeset-font-mono: var(--font-geist-mono);
  --typeset-size: 15px;
  --typeset-leading: 1.75;
  --typeset-flow: 1.25em;
}
```

<Button asChild className="mt-6" size="sm">
  <Link href="/typeset">Build your typeset</Link>
</Button>

---

## Principles

We read a lot about type: scale ratios, tracking, kerning, optical sizing, measure, leading, the space above and below every element. We tried exposing all of it, and it was too much. Nobody wants to set a dozen variables to make markdown look right.

So we sat down and condensed everything into three controls: size, leading, and flow. Everything else, heading sizes, list indents, the gap under a heading, the space around a rule, derives from them. Three controls. We called it rhythm.

---

## Features

- **It fits its container.** Put it in a chat bubble and it follows the smaller type around it. Put it in an article and it scales up with the page. On smaller screens, it gets a small bump for readability.
- **It uses your theme.** Colors, fonts, and radius come from your app. Dark mode follows the same tokens.
- **It's easy to tune.** Three values control the base size, line height, and space between blocks. Change them in a preset and the whole document follows.
- **It works well with streaming.** When a new block arrives, Typeset doesn't make earlier blocks switch margins, borders, or styles.

---

## Building Your Typeset

Create your typeset in the [typeset builder](/typeset). Pick your fonts and rhythm, then preview them on docs, chat, articles, and other real content.

The panel gives you the `typeset.css` file, the font setup for your framework, a preset class with your choices, and the wrapper to add around your content.

Copy `typeset.css` next to your main CSS file and import it after Tailwind:

```css
@import "tailwindcss";
@import "./typeset.css";
```

Then wrap your rendered markdown with `typeset` and your preset class:

```tsx
<div className="typeset typeset-docs">
  <YourMarkdownRenderer>{content}</YourMarkdownRenderer>
</div>
```

`typeset` turns the styles on. `typeset-docs` is the preset you created in the builder.

---

## Custom Typesets

The file includes defaults, so you can use `typeset` by itself. Most of the reading rhythm comes from three values:

```css
.typeset {
  --typeset-font-body: inherit;
  --typeset-font-heading: var(--font-heading);
  --typeset-font-mono: var(--font-mono);

  --typeset-size: 1em; /* body font-size */
  --typeset-leading: 1.75; /* line-height */
  --typeset-flow: 1.25em; /* space between blocks */
}
```

- **`--typeset-size`** sets the base text size. `1em` follows the surrounding layout. On smaller screens, Typeset bumps it up a little.
- **`--typeset-leading`** sets the space between lines.
- **`--typeset-flow`** sets the space between blocks. Headings and other elements derive their spacing from it.

The font variables tell Typeset which families to use. Leave them alone and it follows your app. Colors and radius come from your theme too.

Typeset doesn't set a maximum width. Your layout owns that. The Measure control in the builder adds a `max-width` to the wrapper instead of hiding it in the stylesheet.

You can keep more than one preset in the same app. Here is a tighter one for chat and a roomier one for docs:

```css
.typeset-chat {
  --typeset-flow: 1em;
  --typeset-leading: 1.6;
}

.typeset-docs {
  --typeset-size: 15px;
  --typeset-flow: 1.5em;
}
```

```tsx
<div className="typeset typeset-chat">{message}</div>
<article className="typeset typeset-docs">{page}</article>
```

For a one-off change, skip the preset and set a value on the container:

```tsx
<article className="typeset [--typeset-flow:1.75em]">...</article>
```

---

## Custom Themes

A preset can change the whole feel of the content, not just the spacing. You can give readers a serif reading mode, a compact UI mode, or any other style that fits your product.

```css
/* Reading: serif, larger type, roomy rhythm. */
.typeset-reading {
  --typeset-font-body: var(--font-lora);
  --typeset-font-heading: var(--font-lora);
  --typeset-size: 18px;
  --typeset-leading: 1.9;
  --typeset-flow: 2em;
}

/* Compact: sans, smaller type, tighter rhythm. */
.typeset-compact {
  --typeset-font-body: var(--font-geist);
  --typeset-font-heading: var(--font-geist);
  --typeset-size: 14px;
  --typeset-leading: 1.6;
  --typeset-flow: 1em;
}
```

---

## Accessibility and Dark Mode

For readers who prefer larger type and more space, create a roomier typeset and expose it as a setting:

```css
.typeset-large {
  --typeset-size: 16px;
  --typeset-leading: 2;
  --typeset-flow: 2em;
}
```

Dark mode already follows your theme colors. If the text feels a little tight on a dark surface, you can loosen the leading there:

```css
.dark .typeset {
  --typeset-leading: 1.9;
}
```

---

## Responsive Table

Tables stay real tables and wrap to fit. To scroll a wide one horizontally instead, wrap it in `typeset-scroll`:

```tsx
<div className="typeset-scroll">
  <table>...</table>
</div>
```

Do this in your renderer's table component or a small rehype plugin. It works for any wide block, not just tables.

---

## Overrides

Typeset lives in the `components` layer and uses `:where()` for its element selectors. Tailwind utilities on an element win without `!important`:

```tsx
<div className="typeset typeset-docs">
  <p className="text-lg">...</p>
</div>
```

Plain CSS can override Typeset with a normal selector too.

---

## Opting Out

To keep a component out of Typeset, add `not-typeset` or `data-not-typeset`:

```tsx
<div className="typeset">
  <p>Styled prose.</p>
  <Card className="not-typeset">Untouched component.</Card>
</div>
```

Both options cover the component and everything inside it. Another `typeset` container inside that subtree stays opted out too.

---

## Streaming

Typeset is written so that adding a new block does not change the styles of the blocks already on screen.

- No forward-looking selectors. `:last-child`, `:has()`, and `:empty` are left out of layout rules because their matches can change as content is added.
- Spacing flows in one direction, using `margin-block-start` only. A new block adds its own space.
- Table separators live on the cells being added, so a new row does not restyle the row above it.

Text that is still streaming can grow and wrap normally. Typeset just avoids restyling the blocks that came before it.

---

## Prior Art

The `prose` class from `@tailwindcss/typography` is excellent at what it was built for: adding beautiful typographic defaults to plain HTML, including content rendered from Markdown or a CMS.

Typeset takes a different approach with container-aware sizing, app theme tokens, presets for different contexts, and streaming stability. Here's where they differ:

|              | @tailwindcss/typography                      | Typeset                                          |
| ------------ | -------------------------------------------- | ------------------------------------------------ |
| Sizing       | Fixed `rem` scale, `prose-sm` to `prose-2xl` | Relative to the container, any size              |
| Dark mode    | `prose-invert`, a second palette             | Your tokens flip, nothing to add                 |
| Theming      | Prose color variables; scale baked in        | Your theme tokens, plus font and rhythm controls |
| Overrides    | `prose-a:`, `prose-headings:` modifier API   | Plain utilities and CSS win                      |
| Streaming    | No append-stability contract                 | Designed for stable appends                      |
| Distribution | npm plugin, generated CSS                    | One CSS file you own                             |

Typeset borrows the two best ideas from the plugin: the zero-specificity `:where()` guard pattern, and the escape-hatch class (`not-typeset`, in the spirit of `not-prose`).
