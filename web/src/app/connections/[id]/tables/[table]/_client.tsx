"use client";

import { usePathname } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import DescriptionEditor from "@/components/DescriptionEditor";

export default function TableDetailClient() {
  const pathname = usePathname();
  const m = pathname.match(/^\/connections\/([^/]+)\/tables\/([^/]+)/);
  const connectionId = m?.[1] ?? "";
  const table = m?.[2] ?? "";

  const queryKey = ["table-detail", connectionId, table];

  const { data, isLoading } = useQuery({
    queryKey,
    queryFn: () => api.tables.get(connectionId, table),
    enabled:
      !!connectionId &&
      connectionId !== "_" &&
      !!table &&
      table !== "_",
  });

  if (!connectionId || connectionId === "_") return null;

  if (isLoading) {
    return (
      <div className="p-8 space-y-3">
        <div className="h-8 w-48 bg-slate-100 rounded animate-pulse" />
        <div className="h-64 bg-slate-100 rounded-xl animate-pulse" />
      </div>
    );
  }

  if (!data) {
    return (
      <div className="p-8">
        <p className="text-slate-400">テーブルが見つかりません</p>
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <div className="flex items-baseline gap-3">
          <h2 className="text-xl font-bold font-mono text-slate-900">{table}</h2>
          <span className="text-xs text-slate-400">{data.columns.length} カラム</span>
        </div>
        <div className="mt-2">
          <DescriptionEditor
            connectionId={connectionId}
            tableName={table}
            columnName={null}
            initialValue={data.table.description}
            queryKey={queryKey}
          />
        </div>
      </div>

      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-100">
              <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                カラム名
              </th>
              <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                型
              </th>
              <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                制約
              </th>
              <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                説明
              </th>
            </tr>
          </thead>
          <tbody>
            {data.columns.map((col, i) => (
              <tr
                key={col.name}
                className={`group hover:bg-slate-50 transition-colors ${
                  i > 0 ? "border-t border-slate-100" : ""
                }`}
              >
                <td className="px-5 py-3 align-top">
                  <div className="flex items-center gap-1.5">
                    {col.isPrimaryKey && (
                      <span className="text-amber-400" title="Primary Key">
                        <svg className="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
                          <path d="M17 11a5 5 0 10-5 5v2H8v2h4v2h2v-6a5 5 0 005-5zm-5 3a3 3 0 110-6 3 3 0 010 6z" />
                        </svg>
                      </span>
                    )}
                    <span className="font-mono font-medium text-slate-800">
                      {col.name}
                    </span>
                  </div>
                </td>
                <td className="px-5 py-3 align-top">
                  <span className="font-mono text-xs bg-slate-100 text-slate-700 px-2 py-1 rounded">
                    {col.dataType}
                  </span>
                </td>
                <td className="px-5 py-3 align-top">
                  <div className="flex flex-wrap gap-1">
                    {!col.isNullable && (
                      <span className="text-xs px-1.5 py-0.5 bg-slate-100 text-slate-600 rounded font-medium">
                        NOT NULL
                      </span>
                    )}
                    {col.isForeignKey && (
                      <span className="text-xs px-1.5 py-0.5 bg-sky-50 text-sky-700 rounded ring-1 ring-sky-200 font-medium">
                        → {col.foreignTable}
                      </span>
                    )}
                  </div>
                </td>
                <td className="px-5 py-3 align-top min-w-[200px]">
                  <DescriptionEditor
                    connectionId={connectionId}
                    tableName={table}
                    columnName={col.name}
                    initialValue={col.description}
                    queryKey={queryKey}
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
