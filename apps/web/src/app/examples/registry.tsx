import type { ComponentType } from "react";

import { DataTablePage } from "@/app/examples/data-table-page";
import { FormControlsPage } from "@/app/examples/form-controls-page";
import { FormWithReactionsPage } from "@/app/examples/form-with-reactions-page";
import { ListEditLifecyclePage } from "@/app/examples/list-edit-lifecycle-page";
import { SearchFormTablePage } from "@/app/examples/search-form-table-page";
import type { NavigationContext } from "@/protocol/app-manifest";

// MVP example pages keyed by manifest pageId. These are direct React example
// surfaces (R5 阶段 2); the schema-driven page Renderer is a later boundary.
// Pages that gate on identity receive the boot $context snapshot as a prop.
export interface ExamplePageProps {
  context?: NavigationContext;
}

export const EXAMPLE_PAGES: Record<string, ComponentType<ExamplePageProps>> = {
  "data-table": DataTablePage,
  "form-controls": FormControlsPage,
  "form-with-reactions": FormWithReactionsPage,
  "list-edit-lifecycle": ListEditLifecyclePage,
  "search-form-table": SearchFormTablePage,
};
