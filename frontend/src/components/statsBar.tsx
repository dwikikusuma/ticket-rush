import type { Ticket } from "../types";

interface Props {
  tickets: Ticket[];
}

export default function StatsBar({ tickets }: Props) {
  const total = tickets.length;
  const available = tickets.filter((t) => t.status === "AVAILABLE").length;
  const avgPrice =
    total > 0
      ? tickets.reduce((sum, t) => sum + t.price, 0) / total
      : 0;

  const fmt = (n: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      maximumFractionDigits: 0,
    }).format(n);

  const stats = [
    { label: "Total", value: total },
    { label: "Available", value: available },
    { label: "Avg Price", value: fmt(avgPrice) },
  ];

  if (!total) return null;

  return (
    <div className="flex gap-4 flex-wrap">
      {stats.map(({ label, value }) => (
        <div key={label} className="bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-center">
          <p className="text-gray-400 text-xs">{label}</p>
          <p className="text-white font-bold">{value}</p>
        </div>
      ))}
    </div>
  );
}