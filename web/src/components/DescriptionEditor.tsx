"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import { useQueryClient } from "@tanstack/react-query";

interface Props {
  connectionId: string;
  tableName: string;
  columnName?: string | null;
  initialValue: string;
  queryKey: unknown[];
}

export default function DescriptionEditor({
  connectionId,
  tableName,
  columnName,
  initialValue,
  queryKey,
}: Props) {
  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState(initialValue);
  const [saving, setSaving] = useState(false);
  const queryClient = useQueryClient();

  const save = async () => {
    if (value === initialValue) {
      setEditing(false);
      return;
    }
    setSaving(true);
    try {
      await api.metadata.upsert({
        connectionId,
        tableName,
        columnName: columnName ?? null,
        description: value,
      });
      queryClient.invalidateQueries({ queryKey });
    } finally {
      setSaving(false);
      setEditing(false);
    }
  };

  if (editing) {
    return (
      <textarea
        autoFocus
        className="w-full text-sm bg-white border border-indigo-400 rounded-lg px-3 py-2 resize-none focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-shadow"
        rows={2}
        value={value}
        onChange={(e) => setValue(e.target.value)}
        onBlur={save}
        onKeyDown={(e) => {
          if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            save();
          }
          if (e.key === "Escape") {
            setValue(initialValue);
            setEditing(false);
          }
        }}
        disabled={saving}
        placeholder="説明を入力... (Enter で保存)"
      />
    );
  }

  return (
    <span
      onClick={() => setEditing(true)}
      className="cursor-text text-sm text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded px-1.5 py-0.5 -mx-1.5 block min-h-[1.75rem] transition-colors"
      title="クリックして編集"
    >
      {value || <span className="italic text-slate-300">説明を追加...</span>}
    </span>
  );
}
