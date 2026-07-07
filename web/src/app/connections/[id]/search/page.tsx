import SearchClient from "./_client";

export function generateStaticParams() {
  return [{ id: "_" }];
}

export default function SearchPage() {
  return <SearchClient />;
}
