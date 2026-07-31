import type { ComponentType } from "react";

import { DataTablePage } from "@/app/examples/data-table-page";
import { FormControlsPage } from "@/app/examples/form-controls-page";
import { FormWithReactionsPage } from "@/app/examples/form-with-reactions-page";
import { ListEditLifecyclePage } from "@/app/examples/list-edit-lifecycle-page";
import { SearchFormTablePage } from "@/app/examples/search-form-table-page";

// MVP example pages keyed by manifest pageId. These are direct React example
// surfaces (R5 阶段 2); the schema-driven page Renderer is a later boundary.
export const EXAMPLE_PAGES: Record<string, ComponentType> = {
  "data-table": DataTablePage,
  "form-controls": FormControlsPage,
  "form-with-reactions": FormWithReactionsPage,
  "list-edit-lifecycle": ListEditLifecyclePage,
  "search-form-table": SearchFormTablePage,
};
