import ErClient from "./_client";

export function generateStaticParams() {
  return [{ id: "_" }];
}

export default function ErPage() {
  return <ErClient />;
}
