import { useId, useState, type ReactNode } from "react";

import { ChevronDown } from "lucide-react";

import { useTranslate } from "@/i18n/runtime";
import { cn } from "@/lib/utils";

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
  const hasAdditionalItems = items.length > 1;
  const hasGridActions = actionSlot !== undefined || hasAdditionalItems;
  const activeItemSet = new Set(activeItemIds);
  const activeCountFrom = (firstHiddenIndex: number): number =>
    activeItemIds.filter((id) => itemIds.indexOf(id) >= firstHiddenIndex).length;
  const hiddenAtMobile = 1;
  const hiddenAtSmall = hasGridActions ? 1 : 2;
  const hiddenAtMedium = hasGridActions ? 2 : 3;
  const hiddenAtLarge = hasGridActions ? 3 : 4;

  const collapsedVisibility = (index: number): string => {
    // The action group occupies the last slot of the collapsed first row. This
    // keeps filter controls and their operations on one visual grid row at
    // each responsive breakpoint.
    if (hasGridActions) {
      switch (index) {
        case 0:
          return "";
        case 1:
          return "hidden md:block";
        case 2:
          return "hidden lg:block";
        default:
          return "hidden";
      }
    }
    switch (index) {
      case 0:
        return "";
      case 1:
        return "hidden sm:block";
      case 2:
        return "hidden md:block";
      case 3:
        return "hidden lg:block";
      default:
        return "hidden";
    }
  };

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
              {hasAdditionalItems ? (
                <button
                  type="button"
                  aria-expanded={expanded}
                  aria-controls={panelId}
                  data-filter-toggle="true"
                  className="inline-flex h-8 items-center gap-1.5 rounded-md border border-input/80 bg-background px-2.5 text-xs font-medium text-muted-foreground shadow-2xs transition-colors hover:border-muted-foreground/30 hover:bg-accent/40 hover:text-foreground"
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
      {!expanded &&
      Math.max(
        activeCountFrom(hiddenAtMobile),
        activeCountFrom(hiddenAtSmall),
        activeCountFrom(hiddenAtMedium),
        activeCountFrom(hiddenAtLarge),
      ) > 0 ? (
        <p className="border-t border-border/50 pt-2.5 text-xs text-muted-foreground" data-filter-hidden-active="true">
          {activeCountFrom(1) > 0 ? (
            <span className="sm:hidden">
              {t("feedback.hiddenFiltersActive", { count: String(activeCountFrom(hiddenAtMobile)) })}
            </span>
          ) : null}
          {activeCountFrom(hiddenAtSmall) > 0 ? (
            <span className="hidden sm:inline md:hidden">
              {t("feedback.hiddenFiltersActive", { count: String(activeCountFrom(hiddenAtSmall)) })}
            </span>
          ) : null}
          {activeCountFrom(hiddenAtMedium) > 0 ? (
            <span className="hidden md:inline lg:hidden">
              {t("feedback.hiddenFiltersActive", { count: String(activeCountFrom(hiddenAtMedium)) })}
            </span>
          ) : null}
          {activeCountFrom(hiddenAtLarge) > 0 ? (
            <span className="hidden lg:inline">
              {t("feedback.hiddenFiltersActive", { count: String(activeCountFrom(hiddenAtLarge)) })}
            </span>
          ) : null}
        </p>
      ) : null}
    </section>
  );
}
