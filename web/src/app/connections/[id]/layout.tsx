import { ReactNode } from "react";
import ConnectionLayoutClient from "./_layout_client";

export function generateStaticParams() {
  return [{ id: "_" }];
}

export default function ConnectionLayout({ children }: { children: ReactNode }) {
  return <ConnectionLayoutClient>{children}</ConnectionLayoutClient>;
}
