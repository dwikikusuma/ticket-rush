"use client";

import { useCallback, useRef, useState } from "react";
import { searchApi } from "../services/api";
import type { Ticket } from "../types";

const DEBOUNCE_MS = 400;
const PAGE_LIMIT = 20;

export function useTickets(token?: string | null) {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [cursor, setCursor] = useState<string | undefined>(undefined);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [query, setQuery] = useState("");

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const search = useCallback(async (q: string) => {
    setLoading(true);
    setError(null);
    setCursor(undefined);
    try {
      const res = await searchApi.search(q, PAGE_LIMIT, undefined, token);
      setTickets(res.data ?? []);
      setCursor(res.next_cursor || undefined);
      setHasMore(!!res.next_cursor);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Search failed");
    } finally {
      setLoading(false);
    }
  }, [token]);

  const loadMore = useCallback(async () => {
    if (!cursor || loadingMore) return;
    setLoadingMore(true);
    setError(null);
    try {
      const res = await searchApi.search(query, PAGE_LIMIT, cursor, token);
      setTickets((prev) => [...prev, ...(res.data ?? [])]);
      setCursor(res.next_cursor || undefined);
      setHasMore(!!res.next_cursor);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Load more failed");
    } finally {
      setLoadingMore(false);
    }
  }, [cursor, loadingMore, query, token]);

  const onQueryChange = useCallback(
    (q: string) => {
      setQuery(q);
      if (debounceRef.current) clearTimeout(debounceRef.current);
      debounceRef.current = setTimeout(() => search(q), DEBOUNCE_MS);
    },
    [search]
  );

  return {
    tickets,
    query,
    loading,
    loadingMore,
    hasMore,
    error,
    search,
    loadMore,
    onQueryChange,
  };
}