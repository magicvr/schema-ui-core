import { useEffect, useState } from "react";

import {
  DEFAULT_PAGE_SIZE,
  DEFAULT_RECORDS_URL,
  fetchRecords,
  type RecordList,
  type RecordsQuery,
} from "@/renderer/records";

export interface UseRecordsOptions {
  url?: string;
  fetcher?: typeof fetch;
}

export function useRecords(options: UseRecordsOptions = {}) {
  const url = options.url ?? DEFAULT_RECORDS_URL;
  const fetcher = options.fetcher ?? fetch;
  const [query, setQuery] = useState<RecordsQuery>({
    page: 1,
    pageSize: DEFAULT_PAGE_SIZE,
  });
  const [list, setList] = useState<RecordList | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);
    fetchRecords(fetcher, url, query)
      .then((next) => {
        if (!cancelled) {
          setList(next);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : String(err));
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [fetcher, url, query]);

  return { list, loading, error, query, setQuery };
}
