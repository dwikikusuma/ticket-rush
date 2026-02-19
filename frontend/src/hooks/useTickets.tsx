"use client";

import { useCallback, useRef, useState } from "react";
import { searchApi } from "../services/api";
import type { Ticket } from "../types";

const DEBOUNCE_MS = 400;
const PAGE_LIMIT = 20;

export function useTickets() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [cursor, setCursor] = useState<string | undefined>(undefined);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [query, setQuery] = useState("");

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  /** Replaces current ticket list with fresh results. */
  const search = useCallback(async (q: string) => {
    setLoading(true);
    setError(null);
    setCursor(undefined);
    try {
      const res = await searchApi.search(q, PAGE_LIMIT);
      setTickets(res.data ?? []);
      setCursor(res.next_cursor || undefined);
      setHasMore(!!res.next_cursor);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Search failed");
    } finally {
      setLoading(false);
    }
  }, []);

  /** Appends next page to existing ticket list. */
  const loadMore = useCallback(async () => {
    if (!cursor || loadingMore) return;
    setLoadingMore(true);
    setError(null);
    try {
      const res = await searchApi.search(query, PAGE_LIMIT, cursor);
      setTickets((prev) => [...prev, ...(res.data ?? [])]);
      setCursor(res.next_cursor || undefined);
      setHasMore(!!res.next_cursor);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Load more failed");
    } finally {
      setLoadingMore(false);
    }
  }, [cursor, loadingMore, query]);

  /** Debounced input handler — call this from your SearchBar onChange. */
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