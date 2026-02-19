"use client";

import type { ToastMessage } from "@/types";

interface Props {
  toasts: ToastMessage[];
  onDismiss: (id: string) => void;
}

const typeStyles = {
  success: "bg-emerald-800 border-emerald-600 text-emerald-100",
  error: "bg-red-900 border-red-700 text-red-100",
  info: "bg-blue-900 border-blue-700 text-blue-100",
};

const typeIcons = {
  success: "✅",
  error: "❌",
  info: "ℹ️",
};

export default function Toast({ toasts, onDismiss }: Props) {
  if (!toasts.length) return null;

  return (
    <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-sm w-full">
      {toasts.map((t) => (
        <div
          key={t.id}
          className={`flex items-start gap-2 px-4 py-3 rounded-lg border shadow-xl ${typeStyles[t.type]}`}
        >
          <span className="shrink-0 mt-0.5">{typeIcons[t.type]}</span>
          <p className="flex-1 text-sm">{t.message}</p>
          <button
            onClick={() => onDismiss(t.id)}
            className="shrink-0 opacity-60 hover:opacity-100 transition text-lg leading-none"
            aria-label="Dismiss"
          >
            ×
          </button>
        </div>
      ))}
    </div>
  );
}