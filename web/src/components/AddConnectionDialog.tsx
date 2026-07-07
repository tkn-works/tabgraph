"use client";

import { useState } from "react";
import { ConnectionInput } from "@/lib/types";

interface Props {
  onClose: () => void;
  onCreate: (input: ConnectionInput) => Promise<void>;
}

export default function AddConnectionDialog({ onClose, onCreate }: Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [form, setForm] = useState<ConnectionInput>({
    name: "",
    driver: "postgres",
    host: "localhost",
    port: 5432,
    database: "",
    username: "",
    password: "",
    sslMode: "disable",
  });

  const set = (k: keyof ConnectionInput, v: string | number) =>
    setForm((f) => ({ ...f, [k]: v }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      await onCreate(form);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setLoading(false);
    }
  };

  const inputClass =
    "mt-1 block w-full bg-white border border-slate-300 rounded-lg px-3 py-2 text-sm text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-shadow";

  return (
    <div
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      onClick={(e) => e.target === e.currentTarget && onClose()}
    >
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md">
        <div className="px-6 py-5 border-b border-slate-100">
          <div className="flex items-center justify-between">
            <h2 className="text-base font-semibold text-slate-900">
              DB 接続を追加
            </h2>
            <button
              onClick={onClose}
              className="p-1 rounded-md text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="px-6 py-5 space-y-4">
          <div>
            <label className="text-xs font-semibold text-slate-600 uppercase tracking-wider">
              接続名
            </label>
            <input
              className={inputClass}
              value={form.name}
              onChange={(e) => set("name", e.target.value)}
              required
              placeholder="myapp-prod"
            />
          </div>

          <div>
            <label className="text-xs font-semibold text-slate-600 uppercase tracking-wider">
              ドライバー
            </label>
            <select
              className={inputClass}
              value={form.driver}
              onChange={(e) => {
                set("driver", e.target.value);
                set("port", e.target.value === "postgres" ? 5432 : 3306);
              }}
            >
              <option value="postgres">PostgreSQL</option>
              <option value="mysql">MySQL</option>
            </select>
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div className="col-span-2">
              <label className="text-xs font-semibold text-slate-600 uppercase tracking-wider">
                ホスト
              </label>
              <input
                className={inputClass}
                value={form.host}
                onChange={(e) => set("host", e.target.value)}
                required
              />
            </div>
            <div>
              <label className="text-xs font-semibold text-slate-600 uppercase tracking-wider">
                ポート
              </label>
              <input
                type="number"
                className={inputClass}
                value={form.port}
                onChange={(e) => set("port", Number(e.target.value))}
                required
              />
            </div>
          </div>

          <div>
            <label className="text-xs font-semibold text-slate-600 uppercase tracking-wider">
              データベース名
            </label>
            <input
              className={inputClass}
              value={form.database}
              onChange={(e) => set("database", e.target.value)}
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-xs font-semibold text-slate-600 uppercase tracking-wider">
                ユーザー名
              </label>
              <input
                className={inputClass}
                value={form.username}
                onChange={(e) => set("username", e.target.value)}
                required
              />
            </div>
            <div>
              <label className="text-xs font-semibold text-slate-600 uppercase tracking-wider">
                パスワード
              </label>
              <input
                type="password"
                className={inputClass}
                value={form.password}
                onChange={(e) => set("password", e.target.value)}
              />
            </div>
          </div>

          {form.driver === "postgres" && (
            <div>
              <label className="text-xs font-semibold text-slate-600 uppercase tracking-wider">
                SSL モード
              </label>
              <select
                className={inputClass}
                value={form.sslMode}
                onChange={(e) => set("sslMode", e.target.value)}
              >
                <option value="disable">disable</option>
                <option value="require">require</option>
                <option value="verify-full">verify-full</option>
              </select>
            </div>
          )}

          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg px-3 py-2">
              <p className="text-red-600 text-sm">{error}</p>
            </div>
          )}

          <div className="flex justify-end gap-2 pt-1">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-slate-600 bg-slate-100 rounded-lg hover:bg-slate-200 transition-colors"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-5 py-2 text-sm font-medium bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-60 transition-colors"
            >
              {loading ? "接続中..." : "追加"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
