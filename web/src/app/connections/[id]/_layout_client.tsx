"use client";

import { usePathname } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import { ReactNode } from "react";

export default function ConnectionLayoutClient({
  children,
}: {
  children: ReactNode;
}) {
  const pathname = usePathname();
  const id = pathname.match(/^\/connections\/([^/]+)/)?.[1] ?? "";

  return (
    <div className="flex h-[calc(100vh-56px)]">
      <Sidebar connectionId={id} />
      <div className="flex-1 min-w-0 overflow-y-auto">{children}</div>
    </div>
  );
}
