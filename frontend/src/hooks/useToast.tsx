"use client";

import { useCallback, useState } from "react";
import { uid } from "../lib/utils";
import type { ToastMessage, ToastType } from "../types";

const DURATION_MS = 4000;

export function useToast() {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);

  const dismiss = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const addToast = useCallback(
    (message: string, type: ToastType = "info") => {
      const id = uid();
      setToasts((prev) => [...prev, { id, message, type }]);
      setTimeout(() => dismiss(id), DURATION_MS);
    },
    [dismiss]
  );

  return { toasts, addToast, dismiss };
}