/**
 * Language switcher (S1 · C4) — refactored to a lightweight header
 * dropdown (2026-08-14, user request):
 *
 * - Trigger: pure icon button (lucide Languages), same size/shape as the
 *   theme toggle and notification bell (size-9 rounded-md ghost).
 * - Menu: absolutely positioned panel (NO portal — the app root carries
 *   the `dark` class, so the panel inherits the dark theme variables
 *   automatically; no Radix/Headless portal scope issue).
 * - Dark-mode styling via design tokens (shadcn convention): bg-popover
 *   resolves to ~neutral-900 in dark, border-border to white/10
 *   (~neutral-800), accent hover to ~neutral-800, popover-foreground to
 *   ~neutral-200 — exactly the requested palette; light mode stays
 *   correct automatically.
 * - Selected item shows a checkmark (lucide Check) on the right.
 *
 * Reachable from the Shell and the anonymous login page — no settings
 * permission required (VP-007 requirement).
 */

import { Check, Languages } from "lucide-react";
import { useEffect, useRef, useState } from "react";

import { useI18n } from "@/i18n/runtime";
import { SUPPORTED_LOCALES, type LocalePreference } from "@/i18n/locale";
import { cn } from "@/lib/utils";

export interface LocaleSwitcherProps {
  className?: string;
}

export function LocaleSwitcher({ className = "" }: LocaleSwitcherProps) {
  const { preference, setPreference, t } = useI18n();
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  const label = t("locale.switcher.label");

  const options: LocalePreference[] = ["auto", ...SUPPORTED_LOCALES];

  // Close on outside click / Escape (a11y parity with the notification
  // bell and the mobile drawer).
  useEffect(() => {
    if (!open) {
      return;
    }
    const onPointerDown = (event: PointerEvent) => {
      if (rootRef.current !== null && !rootRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
      }
    };
    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  return (
    <div ref={rootRef} className={cn("relative", className)}>
      <button
        type="button"
        aria-label={label}
        title={label}
        aria-haspopup="menu"
        aria-expanded={open}
        className="inline-flex size-9 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        onClick={() => setOpen((value) => !value)}
      >
        <Languages aria-hidden="true" className="size-4" />
      </button>
      {open ? (
        <div
          role="menu"
          aria-label={label}
          className="absolute right-0 top-11 z-50 w-44 rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-lg"
        >
          {options.map((value) => {
            const selected = preference === value;
            const itemLabel =
              value === "auto" ? t("locale.switcher.auto") : t("locale.name." + value);
            return (
              <button
                key={value}
                type="button"
                role="menuitemradio"
                aria-checked={selected}
                className={cn(
                  "flex w-full items-center justify-between gap-2 rounded-md px-2.5 py-1.5 text-xs transition-colors",
                  selected
                    ? "font-medium text-popover-foreground"
                    : "text-muted-foreground hover:bg-accent hover:text-accent-foreground",
                )}
                onClick={() => {
                  setPreference(value);
                  setOpen(false);
                }}
              >
                <span className="truncate">{itemLabel}</span>
                {selected ? <Check aria-hidden="true" className="size-3.5 shrink-0" /> : null}
              </button>
            );
          })}
        </div>
      ) : null}
    </div>
  );
}
