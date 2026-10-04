import { mapSchema, type SchemaSummary, type TableSummary } from '@/entities/schema';
import { ApiError } from '@/shared/api';
import { connectionsState } from '@/features/connections';
import { getSchema } from '../api/schema';

/** System catalogs hidden unless the eye toggle is on (DESIGN §3.1). */
export const SYSTEM_SCHEMAS = ['pg_catalog', 'information_schema', 'pg_toast'];

export const connKey = (name: string) => `conn:${name}`;
export const schemaKey = (conn: string, schema: string) => `conn:${conn}:schema:${schema}`;

class ExplorerManager {
  // ---- Sidebar layout ----
  sidebarWidth = $state<number>(290);
  servicesHeight = $state<number>(230);
  isSidebarCollapsed = $state<boolean>(false);
  servicesVisible = $state<boolean>(true);
  treeFilterQuery = $state<string>('');
  treeExpanded = $state<Record<string, boolean>>({});
  selectedTreeNode = $state<string | null>(null);
  ddlModalOpen = $state<boolean>(false);
  /** Eye toggle: show Postgres system catalogs. */
  showSystemSchemas = $state<boolean>(false);
  /** Crosshairs toggle: opening a tab reveals its tree node. */
  scrollFromEditor = $state<boolean>(true);

  toggleTree(id: string) {
    this.treeExpanded[id] = !this.treeExpanded[id];
  }

  toggleSystemSchemas() {
    this.showSystemSchemas = !this.showSystemSchemas;
  }

  toggleServices() {
    this.servicesVisible = !this.servicesVisible;
  }

  /** Expand every known node (connections + loaded schemas). */
  expandAll() {
    for (const c of connectionsState.connections) {
      this.treeExpanded[connKey(c.name)] = true;
      for (const s of this.schemasOf(c.name)) {
        this.treeExpanded[schemaKey(c.name, s.name)] = true;
      }
    }
  }

  /** Collapse the whole tree. */
  collapseAll() {
    this.treeExpanded = {};
  }

  /** Expand ancestors of a table and select it (Scroll from Editor). */
  revealTable(connection: string, schema: string, table: string) {
    this.treeExpanded[connKey(connection)] = true;
    this.treeExpanded[schemaKey(connection, schema)] = true;
    this.selectedTreeNode = `table:${connection}.${schema}.${table}`;
  }

  /** Expand a connection and select its console node. */
  revealConsole(connection: string) {
    this.treeExpanded[connKey(connection)] = true;
    this.selectedTreeNode = `console:${connection}`;
  }

  // ---- Live schema cache (no mocks) ----
  schemas = $state<Record<string, SchemaSummary[]>>({});
  schemaLoading = $state<Record<string, boolean>>({});
  schemaError = $state<Record<string, string | null>>({});

  get explorerConnections() {
    const q = this.treeFilterQuery.trim().toLowerCase();
    const all = connectionsState.connections;
    if (!q) return all;
    return all.filter((c) => c.name.toLowerCase().includes(q));
  }

  schemasOf(connection: string): SchemaSummary[] {
    return this.schemas[connection] ?? [];
  }

  /** Schemas visible in the tree (system catalogs need the eye toggle). */
  visibleSchemas(connection: string): SchemaSummary[] {
    const all = this.schemasOf(connection);
    if (this.showSystemSchemas) return all;
    return all.filter((s) => !SYSTEM_SCHEMAS.includes(s.name));
  }

  findTable(connection: string, schema: string, table: string): TableSummary | null {
    for (const s of this.schemasOf(connection)) {
      if (s.name !== schema) continue;
      for (const t of s.tables) {
        if (t.name === table) return t;
      }
    }
    return null;
  }

  async loadSchema(connection: string, force = false): Promise<void> {
    if (!connection) return;
    if (this.schemaLoading[connection]) return;
    if (!force && this.schemas[connection]) return;
    this.schemaLoading[connection] = true;
    this.schemaError[connection] = null;
    try {
      const { schemas } = await getSchema({ connection });
      this.schemas[connection] = schemas.map((s) => mapSchema(s));
    } catch (e) {
      this.schemaError[connection] = e instanceof ApiError ? e.message : String(e);
    } finally {
      this.schemaLoading[connection] = false;
    }
  }

  /** DDL generated from the live schema cache (no hardcoded strings). */
  ddlFor(connection: string, schema: string, table: string): string | null {
    const t = this.findTable(connection, schema, table);
    if (!t) return null;
    const lines = t.columns.map((c) => `    "${c.name}" ${c.type}${c.nullable ? '' : ' NOT NULL'}`);
    if (t.primaryKey.length > 0) {
      lines.push(`    CONSTRAINT ${table}_pkey PRIMARY KEY (${t.primaryKey.map((k) => `"${k}"`).join(', ')})`);
    }
    const qualified = schema ? `"${schema}"."${table}"` : `"${table}"`;
    return `-- Generated from live schema: ${connection}.${schema}.${table}\nCREATE TABLE ${qualified} (\n${lines.join(',\n')}\n);`;
  }
}

export const explorerState = new ExplorerManager();
