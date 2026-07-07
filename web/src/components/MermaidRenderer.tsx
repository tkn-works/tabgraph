"use client";

import { useEffect, useRef } from "react";

interface Props {
  diagram: string;
}

export default function MermaidRenderer({ diagram }: Props) {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!ref.current || !diagram) return;

    import("mermaid").then((m) => {
      m.default.initialize({
        startOnLoad: false,
        theme: "default",
        er: { diagramPadding: 30, layoutDirection: "TB" },
      });
      const id = "mermaid-" + Date.now();
      m.default.render(id, diagram).then(({ svg }) => {
        if (ref.current) ref.current.innerHTML = svg;
      });
    });
  }, [diagram]);

  return (
    <div
      ref={ref}
      className="w-full overflow-auto bg-white border border-gray-200 rounded-lg p-4 min-h-[400px]"
    />
  );
}
