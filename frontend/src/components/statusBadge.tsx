interface Props {
  status: string;
}

const statusStyles: Record<string, string> = {
  AVAILABLE: "bg-emerald-900/60 text-emerald-300 border border-emerald-700",
  SOLD: "bg-red-900/60 text-red-300 border border-red-700",
};

export default function StatusBadge({ status }: Props) {
  const cls =
    statusStyles[status.toUpperCase()] ??
    "bg-gray-800 text-gray-300 border border-gray-600";

  return (
    <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${cls}`}>
      {status}
    </span>
  );
}