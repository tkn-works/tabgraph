"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { ConnectionInput } from "@/lib/types";
import AddConnectionDialog from "@/components/AddConnectionDialog";

export default function HomePage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [showAdd, setShowAdd] = useState(false);

  const { data: connections, isLoading } = useQuery({
    queryKey: ["connections"],
    queryFn: api.connections.list,
  });

  const handleCreate = async (input: ConnectionInput) => {
    const conn = await api.connections.create(input);
    queryClient.invalidateQueries({ queryKey: ["connections"] });
    setShowAdd(false);
    router.push(`/connections/${conn.id}`);
  };

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!confirm("この接続を削除しますか？")) return;
    await api.connections.delete(id);
    queryClient.invalidateQueries({ queryKey: ["connections"] });
  };

  const driverColor = (driver: string) =>
    driver === "postgres"
      ? "bg-sky-50 text-sky-700 ring-1 ring-sky-200"
      : "bg-orange-50 text-orange-700 ring-1 ring-orange-200";

  return (
    <div className="max-w-3xl mx-auto px-6 py-10">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">DB 接続一覧</h1>
          <p className="text-sm text-slate-500 mt-1">
            {connections?.length
              ? `${connections.length} 件の接続`
              : "接続を追加してスキーマを確認しましょう"}
          </p>
        </div>
        <button
          onClick={() => setShowAdd(true)}
          className="inline-flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors shadow-sm"
        >
          <svg
            className="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 4v16m8-8H4"
            />
          </svg>
          接続を追加
        </button>
      </div>

      {isLoading && (
        <div className="space-y-3">
          {[1, 2].map((i) => (
            <div
              key={i}
              className="h-24 bg-slate-100 rounded-xl animate-pulse"
            />
          ))}
        </div>
      )}

      {!isLoading && connections?.length === 0 && (
        <div className="text-center py-20 bg-white rounded-2xl border border-slate-200 border-dashed">
          <div className="w-12 h-12 bg-slate-100 rounded-xl flex items-center justify-center mx-auto mb-4">
            <svg
              className="w-6 h-6 text-slate-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
              />
            </svg>
          </div>
          <p className="font-medium text-slate-700">接続がありません</p>
          <p className="text-sm text-slate-400 mt-1">
            「接続を追加」からDBに接続してください
          </p>
          <button
            onClick={() => setShowAdd(true)}
            className="mt-6 px-4 py-2 text-sm font-medium text-indigo-600 bg-indigo-50 rounded-lg hover:bg-indigo-100 transition-colors"
          >
            + 接続を追加
          </button>
        </div>
      )}

      <div className="space-y-3">
        {connections?.map((conn) => (
          <div
            key={conn.id}
            onClick={() => router.push(`/connections/${conn.id}`)}
            className="group bg-white rounded-xl border border-slate-200 px-5 py-4 cursor-pointer hover:border-indigo-400 hover:shadow-md hover:shadow-indigo-500/5 transition-all"
          >
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 bg-slate-100 rounded-lg flex items-center justify-center shrink-0 group-hover:bg-indigo-50 transition-colors">
                <svg
                  className="w-5 h-5 text-slate-500 group-hover:text-indigo-600 transition-colors"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={1.5}
                    d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
                  />
                </svg>
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <p className="font-semibold text-slate-900 group-hover:text-indigo-700 transition-colors">
                    {conn.name}
                  </p>
                  <span
                    className={`text-xs px-2 py-0.5 rounded-full font-medium ${driverColor(conn.driver)}`}
                  >
                    {conn.driver}
                  </span>
                </div>
                <p className="text-sm text-slate-500 mt-0.5 font-mono truncate">
                  {conn.host}:{conn.port}/{conn.database}
                </p>
              </div>
              <div className="flex items-center gap-3 shrink-0">
                <div className="text-right">
                  <p className="text-xs text-slate-400">
                    {conn.lastSyncedAt
                      ? new Date(conn.lastSyncedAt).toLocaleString("ja-JP", {
                          month: "short",
                          day: "numeric",
                          hour: "2-digit",
                          minute: "2-digit",
                        })
                      : "未 Sync"}
                  </p>
                </div>
                <button
                  onClick={(e) => handleDelete(conn.id, e)}
                  className="opacity-0 group-hover:opacity-100 p-1.5 rounded-md text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all"
                >
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    />
                  </svg>
                </button>
                <svg
                  className="w-4 h-4 text-slate-300 group-hover:text-indigo-400 transition-colors"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 5l7 7-7 7"
                  />
                </svg>
              </div>
            </div>
          </div>
        ))}
      </div>

      {showAdd && (
        <AddConnectionDialog
          onClose={() => setShowAdd(false)}
          onCreate={handleCreate}
        />
      )}
    </div>
  );
}
