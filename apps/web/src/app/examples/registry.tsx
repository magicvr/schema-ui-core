import type { ComponentType } from "react";

import { DataTablePage } from "@/app/examples/data-table-page";
import { SearchFormTablePage } from "@/app/examples/search-form-table-page";

// MVP example pages keyed by manifest pageId. These are direct React example
// surfaces (R5 阶段 2); the schema-driven page Renderer is a later boundary.
export const EXAMPLE_PAGES: Record<string, ComponentType> = {
  "data-table": DataTablePage,
  "search-form-table": SearchFormTablePage,
};
