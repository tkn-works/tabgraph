import TablesClient from "./_client";

export function generateStaticParams() {
  return [{ id: "_" }];
}

export default function TablesPage() {
  return <TablesClient />;
}
