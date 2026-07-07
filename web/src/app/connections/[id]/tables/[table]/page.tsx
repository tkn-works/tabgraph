import TableDetailClient from "./_client";

export function generateStaticParams() {
  return [{ id: "_", table: "_" }];
}

export default function TableDetailPage() {
  return <TableDetailClient />;
}
