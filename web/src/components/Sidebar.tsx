"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useState } from "react";

interface Props {
  connectionId: string;
}

export default function Sidebar({ connectionId }: Props) {
  const pathname = usePathname();
  const router = useRouter();
  const queryClient = useQueryClient();
  const [syncing, setSyncing] = useState(false);

  const { data: connections } = useQuery({
    queryKey: ["connections"],
    queryFn: api.connections.list,
  });

  const { data: tables, isLoading: tablesLoading } = useQuery({
    queryKey: ["tables", connectionId],
    queryFn: () => api.tables.list(connectionId),
    enabled: !!connectionId && connectionId !== "_",
  });

  const currentConn = connections?.find((c) => c.id === connectionId);

  const handleSync = async () => {
    if (!connectionId || connectionId === "_") return;
    setSyncing(true);
    try {
      await api.connections.sync(connectionId);
      queryClient.invalidateQueries({ queryKey: ["tables", connectionId] });
      queryClient.invalidateQueries({ queryKey: ["table-detail", connectionId] });
      queryClient.invalidateQueries({ queryKey: ["er", connectionId] });
    } finally {
      setSyncing(false);
    }
  };

  const isActive = (path: string) =>
    pathname === path || pathname.startsWith(path + "/");

  const navItemClass = (href: string) =>
    `flex items-center px-3 py-2 rounded-md text-sm transition-colors ${
      isActive(href)
        ? "bg-indigo-50 text-indigo-700 font-medium"
        : "text-slate-600 hover:bg-slate-100 hover:text-slate-900"
    }`;

  return (
    <aside className="w-60 shrink-0 bg-white border-r border-slate-200 h-[calc(100vh-56px)] sticky top-14 overflow-y-auto flex flex-col">
      {/* Connection switcher */}
      <div className="p-3 border-b border-slate-100">
        <div className="flex items-start gap-2">
          <div className="flex-1 min-w-0">
            <p className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">
              接続先
            </p>
            {connections && connections.length > 1 ? (
              <select
                className="w-full text-sm font-medium text-slate-800 bg-transparent border-none outline-none cursor-pointer truncate"
                value={connectionId}
                onChange={(e) => router.push(`/connections/${e.target.value}`)}
              >
                {connections.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.name}
                  </option>
                ))}
              </select>
            ) : (
              <p className="text-sm font-medium text-slate-800 truncate">
                {currentConn?.name ?? "—"}
              </p>
            )}
            {currentConn && (
              <p className="text-xs text-slate-400 font-mono truncate mt-0.5">
                {currentConn.driver} · {currentConn.database}
              </p>
            )}
          </div>
          <button
            onClick={handleSync}
            disabled={syncing}
            title="Sync"
            className="shrink-0 mt-5 p-1.5 rounded-md text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 disabled:opacity-50 transition-colors"
          >
            {syncing ? (
              <span className="block w-4 h-4 border-2 border-slate-300 border-t-indigo-500 rounded-full animate-spin" />
            ) : (
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            )}
          </button>
        </div>
      </div>

      {/* Table list */}
      <div className="flex-1 overflow-y-auto">
        <div className="p-3">
          <Link
            href={`/connections/${connectionId}`}
            className={`block text-xs font-semibold uppercase tracking-wider px-2 mb-2 hover:text-indigo-600 transition-colors ${
              pathname === `/connections/${connectionId}`
                ? "text-indigo-600"
                : "text-slate-400"
            }`}
          >
            テーブル一覧
          </Link>
          {tablesLoading && (
            <div className="space-y-1">
              {[1,2,3,4].map(i => (
                <div key={i} className="h-8 bg-slate-100 rounded-md animate-pulse" />
              ))}
            </div>
          )}
          <nav className="space-y-0.5">
            {tables?.map((t) => {
              const href = `/connections/${connectionId}/tables/${t.name}`;
              return (
                <Link key={t.name} href={href} className={navItemClass(href)}>
                  <span className="font-mono truncate">{t.name}</span>
                  <span className="ml-auto text-xs text-slate-400 shrink-0 tabular-nums">
                    {t.columnCount}
                  </span>
                </Link>
              );
            })}
          </nav>
        </div>
      </div>

      {/* Bottom nav */}
      <div className="p-3 border-t border-slate-100">
        <nav className="space-y-0.5">
          {[
            {
              href: `/connections/${connectionId}/er`,
              label: "ER 図",
              icon: (
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M4 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1v-4z" />
                </svg>
              ),
            },
            {
              href: `/connections/${connectionId}/search`,
              label: "検索",
              icon: (
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              ),
            },
          ].map(({ href, label, icon }) => (
            <Link key={href} href={href} className={navItemClass(href)}>
              <span className="shrink-0">{icon}</span>
              <span className="ml-2">{label}</span>
            </Link>
          ))}
        </nav>
        <div className="mt-2 pt-2 border-t border-slate-100">
          <Link
            href="/"
            className="flex items-center gap-1.5 text-xs text-slate-400 hover:text-slate-600 px-2 py-1.5 rounded hover:bg-slate-50 transition-colors"
          >
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            DB 一覧
          </Link>
        </div>
      </div>
    </aside>
  );
}
