"use client";

import { usePathname } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import Link from "next/link";

export default function TablesClient() {
  const pathname = usePathname();
  const connectionId = pathname.match(/^\/connections\/([^/]+)/)?.[1] ?? "";

  const { data: tables, isLoading } = useQuery({
    queryKey: ["tables", connectionId],
    queryFn: () => api.tables.list(connectionId),
    enabled: !!connectionId && connectionId !== "_",
  });

  if (!connectionId || connectionId === "_") return null;

  if (isLoading) {
    return (
      <div className="p-8 space-y-2">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="h-12 bg-slate-100 rounded-lg animate-pulse" />
        ))}
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-slate-900">テーブル一覧</h2>
        <span className="text-sm text-slate-400">{tables?.length ?? 0} テーブル</span>
      </div>

      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-100">
              <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                テーブル名
              </th>
              <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider w-20">
                カラム
              </th>
              <th className="text-left px-5 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                説明
              </th>
            </tr>
          </thead>
          <tbody>
            {tables?.map((t, i) => (
              <tr
                key={t.name}
                className={`group hover:bg-slate-50 transition-colors ${
                  i > 0 ? "border-t border-slate-100" : ""
                }`}
              >
                <td className="px-5 py-3.5">
                  <Link
                    href={`/connections/${connectionId}/tables/${t.name}`}
                    className="font-mono font-medium text-indigo-600 hover:text-indigo-800 hover:underline"
                  >
                    {t.name}
                  </Link>
                </td>
                <td className="px-5 py-3.5">
                  <span className="inline-flex items-center justify-center w-8 h-6 rounded bg-slate-100 text-xs font-medium text-slate-600">
                    {t.columnCount}
                  </span>
                </td>
                <td className="px-5 py-3.5 text-slate-500 text-sm">
                  {t.description || (
                    <span className="text-slate-300 italic">—</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {tables?.length === 0 && (
          <div className="text-center py-12 text-slate-400">
            <p>テーブルがありません。Sync を実行してください。</p>
          </div>
        )}
      </div>
    </div>
  );
}
