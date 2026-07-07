"use client";

import { usePathname } from "next/navigation";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import Link from "next/link";

export default function SearchClient() {
  const pathname = usePathname();
  const connectionId = pathname.match(/^\/connections\/([^/]+)/)?.[1] ?? "";
  const [q, setQ] = useState("");
  const [submitted, setSubmitted] = useState("");

  const { data, isLoading } = useQuery({
    queryKey: ["search", connectionId, submitted],
    queryFn: () => api.search(connectionId, submitted),
    enabled: !!connectionId && connectionId !== "_",
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitted(q.trim());
  };

  if (!connectionId || connectionId === "_") return null;

  return (
    <div className="p-8">
      <h2 className="text-lg font-semibold text-slate-900 mb-6">全文検索</h2>

      <form onSubmit={handleSearch} className="flex gap-2 mb-6 max-w-xl">
        <input
          className="flex-1 bg-white border border-slate-300 rounded-lg px-4 py-2.5 text-sm text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          placeholder="テーブル名・カラム名・説明文を検索..."
          value={q}
          onChange={(e) => setQ(e.target.value)}
          autoFocus
        />
        <button
          type="submit"
          className="px-5 py-2.5 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors"
        >
          検索
        </button>
      </form>

      {isLoading && (
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-16 bg-slate-100 rounded-lg animate-pulse" />
          ))}
        </div>
      )}

      {!isLoading && data && data.length === 0 && submitted.length >= 1 && (
        <div className="text-center py-12 bg-white rounded-xl border border-slate-200">
          <p className="text-slate-500">
            「<span className="font-medium">{submitted}</span>
            」に一致する結果が見つかりませんでした
          </p>
        </div>
      )}

      {data && data.length > 0 && (
        <div className="space-y-2 max-w-xl">
          {data.map((r, i) => (
            <Link
              key={i}
              href={`/connections/${connectionId}/tables/${r.tableName}`}
              className="block bg-white border border-slate-200 rounded-lg px-4 py-3 hover:border-indigo-400 hover:shadow-sm transition-all group"
            >
              <div className="flex items-center gap-2">
                <span
                  className={`text-xs px-2 py-0.5 rounded-full font-medium ${
                    r.matchType === "table"
                      ? "bg-emerald-50 text-emerald-700 ring-1 ring-emerald-200"
                      : "bg-sky-50 text-sky-700 ring-1 ring-sky-200"
                  }`}
                >
                  {r.matchType === "table" ? "テーブル" : "カラム"}
                </span>
                <span className="font-mono text-sm font-medium text-slate-800 group-hover:text-indigo-700 transition-colors">
                  {r.tableName}
                  {r.columnName && (
                    <span className="text-slate-400 font-normal">
                      {" "}
                      / {r.columnName}
                    </span>
                  )}
                </span>
              </div>
              {r.description && (
                <p className="text-sm text-slate-500 mt-1 truncate">
                  {r.description}
                </p>
              )}
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
