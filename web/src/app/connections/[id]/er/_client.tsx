"use client";

import { usePathname } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import MermaidRenderer from "@/components/MermaidRenderer";

export default function ErClient() {
  const pathname = usePathname();
  const connectionId = pathname.match(/^\/connections\/([^/]+)/)?.[1] ?? "";

  const { data, isLoading } = useQuery({
    queryKey: ["er", connectionId],
    queryFn: () => api.er(connectionId),
    enabled: !!connectionId && connectionId !== "_",
  });

  if (!connectionId || connectionId === "_") return null;

  if (isLoading) {
    return (
      <div className="p-8">
        <div className="h-96 bg-slate-100 rounded-xl animate-pulse" />
      </div>
    );
  }

  return (
    <div className="p-8">
      <h2 className="text-lg font-semibold text-slate-900 mb-6">ER 図</h2>
      {data?.diagram ? (
        <MermaidRenderer diagram={data.diagram} />
      ) : (
        <div className="text-center py-16 bg-white rounded-xl border border-slate-200">
          <p className="text-slate-400">Sync を実行するとER図が表示されます。</p>
        </div>
      )}
    </div>
  );
}
