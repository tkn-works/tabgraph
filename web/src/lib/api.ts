import {
  Connection,
  ConnectionInput,
  TableDetail,
  TableInfo,
  Metadata,
  MetadataInput,
  SearchResult,
  SyncResult,
} from "./types";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? (typeof window !== "undefined" ? window.location.origin : "");

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error ?? res.statusText);
  }
  return res.json();
}

export const api = {
  connections: {
    list: () => request<Connection[]>("/api/connections"),
    create: (input: ConnectionInput) =>
      request<Connection>("/api/connections", {
        method: "POST",
        body: JSON.stringify(input),
      }),
    delete: (id: string) =>
      fetch(`${BASE_URL}/api/connections/${id}`, { method: "DELETE" }),
    sync: (id: string) =>
      request<SyncResult>(`/api/connections/${id}/sync`, { method: "POST" }),
  },
  tables: {
    list: (connectionId: string) =>
      request<TableInfo[]>(`/api/connections/${connectionId}/tables`),
    get: (connectionId: string, tableName: string) =>
      request<TableDetail>(
        `/api/connections/${connectionId}/tables/${tableName}`
      ),
  },
  metadata: {
    upsert: (input: MetadataInput) =>
      request<Metadata>("/api/metadata", {
        method: "PUT",
        body: JSON.stringify(input),
      }),
  },
  search: (connectionId: string, q: string) =>
    request<SearchResult[]>(
      `/api/connections/${connectionId}/search?q=${encodeURIComponent(q)}`
    ),
  er: (connectionId: string) =>
    request<{ diagram: string }>(`/api/connections/${connectionId}/er`),
};
