import { useEffect, useId, useState, type ReactNode } from "react";

import { ChevronDown } from "lucide-react";

import { useTranslate } from "@/i18n/runtime";
import { cn } from "@/lib/utils";

/**
 * Responsive tiers of the collapsed filter row. `base` is the mobile-first
 * tier (no media query); the rest map to Tailwind's `sm`/`md`/`lg`.
 */
const TIERS = ["base", "sm", "md", "lg"] as const;
type Tier = (typeof TIERS)[number];

/** Tailwind visibility class for an item that first appears at `tier`. */
const TIER_CLASS: Record<Tier, string> = {
  base: "",
  sm: "hidden sm:block",
  md: "hidden md:block",
  lg: "hidden lg:block",
};

/**
 * Visible slots in the collapsed first row, per tier.
 *
 * The action cell (reset + expand/collapse) always claims the LAST slot of
 * that row, so it costs exactly one slot. These counts are the single source
 * of truth for both the visibility classes and the "is anything hidden?"
 * decision (R6 C8), so the two can never drift apart.
 */
const COLLAPSED_SLOTS: Record<Tier, { withActions: number; withoutActions: number }> = {
  base: { withActions: 1, withoutActions: 1 },
  sm: { withActions: 1, withoutActions: 2 },
  md: { withActions: 2, withoutActions: 3 },
  lg: { withActions: 3, withoutActions: 4 },
};

/** Collapsed-row capacity at `tier`, given whether the action cell is used. */
function collapsedCapacityFor(tier: Tier, hasGridActions: boolean): number {
  return hasGridActions
    ? COLLAPSED_SLOTS[tier].withActions
    : COLLAPSED_SLOTS[tier].withoutActions;
}

function tierAtWidth(width: number): Tier {
  if (width >= 1024) {
    return "lg";
  }
  if (width >= 768) {
    return "md";
  }
  if (width >= 640) {
    return "sm";
  }
  return "base";
}

/** Resolves the current tier from media queries, falling back to the viewport. */
function currentTier(): Tier {
  if (typeof window === "undefined") {
    return "lg";
  }
  const { matchMedia, innerWidth } = window;
  if (typeof matchMedia === "function") {
    if (matchMedia("(min-width: 1024px)").matches) {
      return "lg";
    }
    if (matchMedia("(min-width: 768px)").matches) {
      return "md";
    }
    if (matchMedia("(min-width: 640px)").matches) {
      return "sm";
    }
    return "base";
  }
  return tierAtWidth(typeof innerWidth === "number" ? innerWidth : 0);
}

/** Tracks the active responsive tier across breakpoint changes. */
function useTier(): Tier {
  const [tier, setTier] = useState<Tier>(currentTier);
  useEffect(() => {
    if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
      return;
    }
    const queries = ["(min-width: 640px)", "(min-width: 768px)", "(min-width: 1024px)"].map(
      (query) => window.matchMedia(query),
    );
    const update = () => setTier(currentTier());
    for (const query of queries) {
      query.addEventListener?.("change", update);
    }
    update();
    return () => {
      for (const query of queries) {
        query.removeEventListener?.("change", update);
      }
    };
  }, []);
  return tier;
}

export interface ListFilterPanelProps {
  /** Filter controls in their semantic order. Values remain mounted when hidden. */
  items: ReactNode[];
  /** Stable ids matching `items`, used only for the hidden-condition hint. */
  itemIds?: string[];
  /** Submitted filter ids, never draft values. */
  activeItemIds?: string[];
  /** Controls that stay available while the filter grid is collapsed. */
  actionSlot?: ReactNode;
  /** Keeps the host form's data attributes on the shared presentation surface. */
  dataAttributes?: Record<string, string | undefined>;
  /** Search forms already provide the surrounding card surface. */
  surface?: "card" | "plain";
  className?: string;
}

/**
 * Shared list-filter presentation.
 *
 * This component owns only presentation state. It does not own query values,
 * submit handlers, or field visibility decisions from the reaction engine.
 * Hidden controls stay mounted but use `display: none`, so they cannot receive
 * focus while collapsed and their React state remains intact.
 *
 * R6 C8 (user 2026-09-18): the expand/collapse toggle is rendered only when
 * the collapsed first row actually hides at least one control at the current
 * width. A single-row filter set has nothing to expand, so the toggle — and
 * the action cell it would occupy — is omitted entirely.
 */
export function ListFilterPanel({
  items,
  itemIds = [],
  activeItemIds = [],
  actionSlot,
  dataAttributes,
  surface = "card",
  className,
}: ListFilterPanelProps) {
  const t = useTranslate();
  const [expanded, setExpanded] = useState(false);
  const panelId = useId();
  const tier = useTier();

  // The action cell is needed when the host supplies actions, or when the
  // items do not fit the collapsed row without one. Claiming that cell costs
  // a slot, so the effective capacity depends on the decision itself.
  const fitsWithoutActions = items.length <= COLLAPSED_SLOTS[tier].withoutActions;
  const hasGridActions = actionSlot !== undefined || !fitsWithoutActions;
  const collapsedCapacity = collapsedCapacityFor(tier, hasGridActions);
  const hasHiddenItems = items.length > collapsedCapacity;

  const activeItemSet = new Set(activeItemIds);
  const activeCountFrom = (firstHiddenIndex: number): number =>
    activeItemIds.filter((id) => itemIds.indexOf(id) >= firstHiddenIndex).length;

  // Lowest tier at which the item survives the collapsed row; "" never
  // survives. Derived from the same slot table as the capacity decision, so
  // the "N filters hidden" hint can never disagree with what is really hidden.
  const slotCounts: Record<Tier, number> = {
    base: collapsedCapacityFor("base", hasGridActions),
    sm: collapsedCapacityFor("sm", hasGridActions),
    md: collapsedCapacityFor("md", hasGridActions),
    lg: collapsedCapacityFor("lg", hasGridActions),
  };
  const collapsedVisibility = (index: number): string => {
    const visibleFrom = TIERS.find((candidate) => index < slotCounts[candidate]);
    return visibleFrom === undefined ? "hidden" : TIER_CLASS[visibleFrom];
  };
  // Index of the first hidden item at each tier (its own slot count), used by
  // the hidden-active hint below.
  const hiddenFromMobile = slotCounts.base;
  const hiddenFromSmall = slotCounts.sm;
  const hiddenFromMedium = slotCounts.md;
  const hiddenFromLarge = slotCounts.lg;

  if (items.length === 0) {
    return actionSlot === undefined ? null : <div className={className}>{actionSlot}</div>;
  }

  return (
    <section
      {...dataAttributes}
      data-list-filter-panel="true"
      data-list-filter-expanded={expanded ? "true" : "false"}
      className={cn(
        surface === "card"
          ? "space-y-3 rounded-xl border border-border/70 bg-card/85 p-4 shadow-[0_1px_3px_0_rgba(0,0,0,0.03),0_1px_2px_-1px_rgba(0,0,0,0.03)] dark:border-border/60 dark:bg-card/70 dark:shadow-[0_1px_3px_0_rgba(0,0,0,0.2)]"
          : "space-y-3",
        className,
      )}
    >
      <div
        id={panelId}
        className="grid grid-cols-1 items-end gap-3 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
        data-list-filter-grid="true"
      >
        {items.map((item, index) => {
          const itemId = itemIds[index];
          const active = itemId !== undefined && activeItemSet.has(itemId);
          return (
            <div
              key={itemId ?? index}
              className={expanded ? "min-w-0" : cn("min-w-0", collapsedVisibility(index))}
              data-filter-item="true"
              data-filter-item-index={String(index)}
              data-filter-item-active={active ? "true" : undefined}
            >
              {item}
            </div>
          );
        })}
        {hasGridActions ? (
          <div
            className="flex min-w-0 items-end justify-end sm:col-start-2 md:col-start-3 lg:col-start-4"
            data-filter-actions="true"
          >
            <div className="flex flex-wrap items-center justify-end gap-2">
              {actionSlot}
              {hasHiddenItems ? (
                <button
                  type="button"
                  aria-expanded={expanded}
                  aria-controls={panelId}
                  data-filter-toggle="true"
                  // R6 C8: the reference page wears this toggle as a filled
                  // `control` chip — visibly distinct from the plain bordered
                  // reset action beside it (new semantic token, not a restyle
                  // of the default action).
                  className="inline-flex h-8 items-center gap-1.5 rounded-md border border-input/80 bg-control px-2.5 text-xs font-medium text-control-foreground shadow-2xs transition-colors hover:border-muted-foreground/30 hover:bg-control/70"
                  onClick={() => setExpanded((current) => !current)}
                >
                  {expanded ? t("feedback.collapseFilters") : t("feedback.expandFilters")}
                  <ChevronDown
                    aria-hidden="true"
                    className={cn("size-3.5 transition-transform", expanded ? "rotate-180" : "")}
                  />
                </button>
              ) : null}
            </div>
          </div>
        ) : null}
      </div>
      {!expanded && hasHiddenItems &&
      Math.max(
        activeCountFrom(hiddenFromMobile),
        activeCountFrom(hiddenFromSmall),
        activeCountFrom(hiddenFromMedium),
        activeCountFrom(hiddenFromLarge),
      ) > 0 ? (
        <p className="border-t border-border/50 pt-2.5 text-xs text-muted-foreground" data-filter-hidden-active="true">
          {activeCountFrom(hiddenFromMobile) > 0 ? (
            <span className="sm:hidden">
              {t("feedback.hiddenFiltersActive", { count: String(activeCountFrom(hiddenFromMobile)) })}
            </span>
          ) : null}
          {activeCountFrom(hiddenFromSmall) > 0 ? (
            <span className="hidden sm:inline md:hidden">
              {t("feedback.hiddenFiltersActive", { count: String(activeCountFrom(hiddenFromSmall)) })}
            </span>
          ) : null}
          {activeCountFrom(hiddenFromMedium) > 0 ? (
            <span className="hidden md:inline lg:hidden">
              {t("feedback.hiddenFiltersActive", { count: String(activeCountFrom(hiddenFromMedium)) })}
            </span>
          ) : null}
          {activeCountFrom(hiddenFromLarge) > 0 ? (
            <span className="hidden lg:inline">
              {t("feedback.hiddenFiltersActive", { count: String(activeCountFrom(hiddenFromLarge)) })}
            </span>
          ) : null}
        </p>
      ) : null}
    </section>
  );
}
