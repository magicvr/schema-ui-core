import { Command, LoaderCircle, Search, X } from "lucide-react";
import { useEffect, useMemo, useRef, useState, type KeyboardEvent as ReactKeyboardEvent } from "react";

import { useTranslate } from "@/i18n/runtime";
import {
  aggregateSearchableItems,
  DEFAULT_SEARCHABLE_RESULT_LIMIT,
  filterSearchableItems,
  type SearchableItem,
  type SearchableProvider,
  type SearchableProviderContext,
} from "@/app/searchable";

const PALETTE_INPUT_ID = "command-palette-input";
const PALETTE_RESULTS_ID = "command-palette-results";

function optionId(item: SearchableItem): string {
  return `command-palette-option-${item.id.replace(/[^A-Za-z0-9_-]/g, "-")}`;
}

function focusableElements(root: HTMLElement): HTMLElement[] {
  return Array.from(
    root.querySelectorAll<HTMLElement>(
      'button:not([disabled]), input:not([disabled]), [href], select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
    ),
  );
}

/** Platform-aware visual shortcut label; behavior accepts both modifiers. */
export function commandPaletteShortcutLabel(): string {
  if (typeof navigator !== "undefined" && /Mac|iPhone|iPad/.test(navigator.platform)) {
    return "⌘K";
  }
  return "Ctrl+K";
}

export interface CommandPaletteProps {
  open: boolean;
  providers: readonly SearchableProvider[];
  context: SearchableProviderContext;
  onClose: () => void;
  onSelect: (item: SearchableItem) => void;
  resultLimit?: number;
}

/**
 * Keyboard-first global discovery surface. The component is deliberately
 * provider-agnostic: visibility and action safety are decided before items
 * reach the list, while this layer owns only matching, focus and presentation.
 */
export function CommandPalette({
  open,
  providers,
  context,
  onClose,
  onSelect,
  resultLimit = DEFAULT_SEARCHABLE_RESULT_LIMIT,
}: CommandPaletteProps) {
  const t = useTranslate();
  const dialogRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const previousFocusRef = useRef<HTMLElement | null>(null);
  const restoreFocusRef = useRef(true);
  const [query, setQuery] = useState("");
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [items, setItems] = useState<SearchableItem[]>([]);
  const [errors, setErrors] = useState<Array<{ code: string; pageId?: string }>>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!open) {
      return;
    }
    previousFocusRef.current =
      document.activeElement instanceof HTMLElement ? document.activeElement : null;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    inputRef.current?.focus();

    const onDocumentKeyDown = (event: globalThis.KeyboardEvent) => {
      const dialog = dialogRef.current;
      if (dialog === null) {
        return;
      }
      if (event.key === "Escape") {
        event.preventDefault();
        event.stopPropagation();
        onClose();
        return;
      }
      if (event.key !== "Tab") {
        return;
      }
      const focusable = focusableElements(dialog);
      if (focusable.length === 0) {
        event.preventDefault();
        dialog.focus();
        return;
      }
      const first = focusable[0];
      const last = focusable[focusable.length - 1];
      const active = document.activeElement;
      if (event.shiftKey && (active === first || !dialog.contains(active))) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && (active === last || !dialog.contains(active))) {
        event.preventDefault();
        first.focus();
      }
    };
    document.addEventListener("keydown", onDocumentKeyDown, true);
    return () => {
      document.removeEventListener("keydown", onDocumentKeyDown, true);
      document.body.style.overflow = previousOverflow;
      if (restoreFocusRef.current) {
        previousFocusRef.current?.focus();
      }
      previousFocusRef.current = null;
      restoreFocusRef.current = true;
    };
  }, [open, onClose]);

  useEffect(() => {
    if (!open) {
      setQuery("");
      setSelectedIndex(0);
      setItems([]);
      setErrors([]);
      setLoading(false);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setErrors([]);
    void aggregateSearchableItems(providers, context).then((result) => {
      if (cancelled) {
        return;
      }
      setItems(result.items);
      setErrors(result.errors.map(({ code, pageId }) => ({ code, ...(pageId === undefined ? {} : { pageId }) })));
      setLoading(false);
    });
    return () => {
      cancelled = true;
    };
  }, [open, providers, context]);

  const visibleItems = useMemo(
    () => filterSearchableItems(items, query, resultLimit),
    [items, query, resultLimit],
  );

  useEffect(() => {
    setSelectedIndex((current) =>
      visibleItems.length === 0 ? 0 : Math.min(current, visibleItems.length - 1),
    );
  }, [visibleItems.length, query]);

  if (!open) {
    return null;
  }

  const selectItem = (item: SearchableItem) => {
    restoreFocusRef.current = false;
    onSelect(item);
    onClose();
  };

  const moveSelection = (delta: number) => {
    if (visibleItems.length === 0) {
      return;
    }
    setSelectedIndex((current) => (current + delta + visibleItems.length) % visibleItems.length);
  };

  const handleKeyDown = (event: ReactKeyboardEvent<HTMLDivElement>) => {
    if (event.key === "ArrowDown") {
      event.preventDefault();
      moveSelection(1);
    } else if (event.key === "ArrowUp") {
      event.preventDefault();
      moveSelection(-1);
    } else if (event.key === "Home") {
      if (event.target === inputRef.current) {
        return;
      }
      event.preventDefault();
      setSelectedIndex(0);
    } else if (event.key === "End") {
      if (event.target === inputRef.current) {
        return;
      }
      event.preventDefault();
      setSelectedIndex(Math.max(visibleItems.length - 1, 0));
    } else if (event.key === "Enter") {
      event.preventDefault();
      const selected = visibleItems[selectedIndex];
      if (selected !== undefined) {
        selectItem(selected);
      }
    }
  };

  return (
    <div
      data-command-palette-backdrop
      className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-overlay p-4 sm:p-8"
      onMouseDown={(event) => {
        if (event.target === event.currentTarget) {
          onClose();
        }
      }}
    >
      <div
        ref={dialogRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby="command-palette-title"
        className="mt-[8vh] flex w-full max-w-2xl flex-col overflow-hidden rounded-xl border border-border bg-card shadow-2xl shadow-black/30 outline-none"
        tabIndex={-1}
        onMouseDown={(event) => event.stopPropagation()}
        onKeyDown={handleKeyDown}
      >
        <div className="flex items-center gap-3 border-b border-border px-4">
          <Search aria-hidden="true" className="size-4 shrink-0 text-muted-foreground" />
          <h2 id="command-palette-title" className="sr-only">
            {t("commandPalette.title")}
          </h2>
          <input
            ref={inputRef}
            id={PALETTE_INPUT_ID}
            role="combobox"
            aria-label={t("commandPalette.input")}
            aria-expanded="true"
            aria-controls={PALETTE_RESULTS_ID}
            aria-autocomplete="list"
            aria-activedescendant={
              visibleItems[selectedIndex] === undefined
                ? undefined
                : optionId(visibleItems[selectedIndex])
            }
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder={t("commandPalette.placeholder")}
            className="h-14 min-w-0 flex-1 bg-transparent text-sm text-foreground outline-none placeholder:text-muted-foreground/60"
          />
          <kbd className="hidden rounded border border-border bg-muted/50 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground sm:inline-block">
            Esc
          </kbd>
          <button
            type="button"
            aria-label={t("commandPalette.close")}
            onClick={onClose}
            className="inline-flex size-8 shrink-0 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <X aria-hidden="true" className="size-4" />
          </button>
        </div>

        {loading ? (
          <div role="status" className="flex items-center gap-2 px-4 py-8 text-sm text-muted-foreground">
            <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />
            {t("commandPalette.loading")}
          </div>
        ) : null}

        {!loading && errors.length > 0 ? (
          <p role="alert" className="border-b border-border/70 bg-destructive/5 px-4 py-2 text-xs text-destructive">
            {t("commandPalette.partialError")}
          </p>
        ) : null}

        <div
          id={PALETTE_RESULTS_ID}
          role="listbox"
          aria-busy={loading}
          aria-label={t("commandPalette.results")}
          className="max-h-[min(60vh,28rem)] overflow-y-auto p-2"
        >
          {!loading && visibleItems.length > 0
            ? visibleItems.map((item, index) => {
                const selected = index === selectedIndex;
                return (
                  <div
                    key={item.id}
                    id={optionId(item)}
                    role="option"
                    aria-selected={selected}
                    tabIndex={-1}
                    onMouseEnter={() => setSelectedIndex(index)}
                    onClick={() => selectItem(item)}
                    className={`flex w-full cursor-pointer items-center gap-3 rounded-lg px-3 py-2.5 text-left transition-colors ${
                      selected ? "bg-accent text-accent-foreground" : "text-foreground hover:bg-accent/70"
                    }`}
                  >
                    <span
                      aria-hidden="true"
                      className={`inline-flex size-8 shrink-0 items-center justify-center rounded-md border ${
                        selected ? "border-border bg-background/70" : "border-border/70 bg-muted/30"
                      }`}
                    >
                      {item.kind === "action" ? (
                        <Command className="size-4" />
                      ) : (
                        <Search className="size-4" />
                      )}
                    </span>
                    <span className="min-w-0 flex-1">
                      <span className="block truncate text-sm font-medium">{item.label}</span>
                      <span className="mt-0.5 flex min-w-0 items-center gap-2 truncate text-xs text-muted-foreground">
                        <span>{t(`commandPalette.kind.${item.kind}`)}</span>
                        {item.group !== undefined && item.group !== "" ? (
                          <>
                            <span aria-hidden="true">·</span>
                            <span className="truncate">{item.group}</span>
                          </>
                        ) : null}
                      </span>
                    </span>
                    <span className="hidden shrink-0 font-mono text-[10px] text-muted-foreground/70 sm:inline">
                      {item.kind === "action" ? t("commandPalette.openAction") : t("commandPalette.openPage")}
                    </span>
                  </div>
                );
              })
            : null}
          {!loading && visibleItems.length === 0 ? (
            <p role="status" className="px-3 py-8 text-center text-sm text-muted-foreground">
              {t("commandPalette.empty")}
            </p>
          ) : null}
        </div>

        <div className="flex flex-wrap items-center gap-x-4 gap-y-1 border-t border-border px-4 py-2 text-[11px] text-muted-foreground">
          <span>
            <kbd className="font-mono">↑↓</kbd> {t("commandPalette.move")}
          </span>
          <span>
            <kbd className="font-mono">Enter</kbd> {t("commandPalette.select")}
          </span>
          <span className="ml-auto">
            {t("commandPalette.resultCount", { count: visibleItems.length })}
          </span>
        </div>
      </div>
    </div>
  );
}