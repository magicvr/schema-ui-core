import { createContext, useCallback, useContext, useMemo, useRef, useState, type ReactNode } from "react";

export interface PageListActionsHost {
  element: HTMLElement | null;
  activeOwner: string | null;
  availabilityVersion: number;
  claim: (ownerId: string) => boolean;
  release: (ownerId: string) => void;
}

const PageListActionsContext = createContext<PageListActionsHost | null>(null);

export function PageListActionsProvider({
  host,
  children,
}: {
  host: HTMLElement | null;
  children: ReactNode;
}) {
  const ownerRef = useRef<string | null>(null);
  const [activeOwner, setActiveOwner] = useState<string | null>(null);
  const [availabilityVersion, setAvailabilityVersion] = useState(0);
  const claim = useCallback((ownerId: string) => {
    if (ownerRef.current !== null && ownerRef.current !== ownerId) {
      return false;
    }
    ownerRef.current = ownerId;
    setActiveOwner(ownerId);
    return true;
  }, []);
  const release = useCallback((ownerId: string) => {
    if (ownerRef.current !== ownerId) {
      return;
    }
    ownerRef.current = null;
    setActiveOwner(null);
    setAvailabilityVersion((current) => current + 1);
  }, []);
  const value = useMemo<PageListActionsHost>(
    () => ({ element: host, activeOwner, availabilityVersion, claim, release }),
    [activeOwner, availabilityVersion, claim, host, release],
  );
  return <PageListActionsContext.Provider value={value}>{children}</PageListActionsContext.Provider>;
}

export function usePageListActionsHost(): PageListActionsHost | null {
  return useContext(PageListActionsContext);
}

export interface ListObjectLabelOptions {
  pageId?: string;
  pageTitle?: string;
  tableTitle?: string;
  tableTitleKey?: string;
  translate: (key: string, params?: Record<string, string>) => string;
}

/**
 * Resolve the noun shown in the default Saved View option without guessing
 * from a URL or API path. Known page/table semantics win; schema titles are
 * the reliable fallback for module-owned tables.
 */
export function resolveListObjectLabel({
  pageId,
  pageTitle,
  tableTitle,
  tableTitleKey,
  translate,
}: ListObjectLabelOptions): string {
  const knownKey =
    tableTitleKey === "schema.account.session.title"
      ? "feedback.listObject.sessions"
      : pageId === "users"
        ? "feedback.listObject.users"
        : pageId === "roles"
          ? "feedback.listObject.roles"
          : undefined;
  if (knownKey !== undefined) {
    const translated = translate(knownKey);
    if (translated !== knownKey && translated.trim() !== "") {
      return translated;
    }
  }
  for (const candidate of [tableTitle, pageTitle]) {
    if (typeof candidate === "string" && candidate.trim() !== "") {
      return candidate.trim();
    }
  }
  return translate("feedback.listObject.records") === "feedback.listObject.records"
    ? "Records"
    : translate("feedback.listObject.records");
}
