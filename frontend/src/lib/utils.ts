/** Format a number as USD currency, no cents. */
export function formatPrice(amount: number): string {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    maximumFractionDigits: 0,
  }).format(amount);
}

/** Decode JWT payload without a library. Returns null on failure. */
export function decodeJwt<T = Record<string, unknown>>(token: string): T | null {
  try {
    const payload = token.split(".")[1];
    return JSON.parse(atob(payload)) as T;
  } catch {
    return null;
  }
}

/** Generate a simple unique ID for toast messages. */
export function uid(): string {
  return Math.random().toString(36).slice(2, 9);
}