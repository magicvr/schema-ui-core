/**
 * Timezone switcher (workspace-020 · R2 · C2).
 *
 * User-level timezone override next to the language switcher in the header
 * locale channel (contract GOAL-002 D-001 §4.2): persists via the single
 * localStorage channel "schema-ui:timezone"; "auto" removes the key.
 * The option list is the documented common set (contract §6 permits a
 * verifiable, extendable list); the effective timezone always degrades
 * safely through the L1–L4 resolver in i18n/timezone.
 */

import { Check, Clock } from "lucide-react";
import { useEffect, useRef, useState } from "react";

import { useI18n } from "@/i18n/runtime";
import { AUTO_TIMEZONE, type TimezonePreference } from "@/i18n/timezone";
import { cn } from "@/lib/utils";

/** Common IANA set offered in the header menu (extendable per contract §6). */
export const TIMEZONE_OPTIONS: readonly string[] = [
  "Asia/Shanghai",
  "Asia/Tokyo",
  "America/New_York",
  "Europe/London",
  "UTC",
];

export interface TimezoneSwitcherProps {
  className?: string;
}

export function TimezoneSwitcher({ className = "" }: TimezoneSwitcherProps) {
  const { timezonePreference, setTimezonePreference, t } = useI18n();
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  const label = t("timezone.switcher.label");

  const options: TimezonePreference[] = [AUTO_TIMEZONE, ...TIMEZONE_OPTIONS];

  // Close on outside click / Escape (a11y parity with the locale switcher).
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
        <Clock aria-hidden="true" className="size-4" />
      </button>
      {open ? (
        <div
          role="menu"
          aria-label={label}
          className="absolute right-0 top-11 z-50 w-52 rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-lg"
        >
          {options.map((value) => {
            const selected = timezonePreference === value;
            const itemLabel =
              value === AUTO_TIMEZONE ? t("timezone.switcher.auto") : value;
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
                  setTimezonePreference(value);
                  setOpen(false);
                }}
              >
                <span className="truncate font-mono">{itemLabel}</span>
                {selected ? <Check aria-hidden="true" className="size-3.5 shrink-0" /> : null}
              </button>
            );
          })}
        </div>
      ) : null}
    </div>
  );
}