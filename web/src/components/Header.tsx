"use client";

import Link from "next/link";

export default function Header() {
  return (
    <header className="h-14 bg-white border-b border-slate-200 flex items-center px-5 shrink-0 sticky top-0 z-50">
      <Link
        href="/"
        className="font-bold text-base text-slate-900 tracking-tight hover:text-indigo-600 transition-colors"
      >
        tabgraph
      </Link>
    </header>
  );
}
