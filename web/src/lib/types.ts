export interface Connection {
  id: string;
  name: string;
  driver: "postgres" | "mysql";
  host: string;
  port: number;
  database: string;
  username: string;
  sslMode: string;
  createdAt: string;
  lastSyncedAt: string | null;
}

export interface ConnectionInput {
  name: string;
  driver: "postgres" | "mysql";
  host: string;
  port: number;
  database: string;
  username: string;
  password: string;
  sslMode?: string;
}

export interface TableInfo {
  name: string;
  columnCount: number;
  description: string;
}

export interface ColumnInfo {
  name: string;
  dataType: string;
  isNullable: boolean;
  default: string | null;
  isPrimaryKey: boolean;
  isForeignKey: boolean;
  foreignTable: string | null;
  foreignColumn: string | null;
  description: string;
  ordinalPosition: number;
}

export interface TableDetail {
  table: TableInfo;
  columns: ColumnInfo[];
}

export interface Metadata {
  id: string;
  connectionId: string;
  tableName: string;
  columnName: string | null;
  description: string;
  updatedAt: string;
}

export interface MetadataInput {
  connectionId: string;
  tableName: string;
  columnName?: string | null;
  description: string;
}

export interface SearchResult {
  tableName: string;
  columnName: string | null;
  description: string;
  matchType: "table" | "column";
}

export interface SyncResult {
  tablesCount: number;
  columnsCount: number;
  syncedAt: string;
}
