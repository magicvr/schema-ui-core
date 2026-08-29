/**
 * @schema-ui/renderer 聚合导出（试点 v0.1.0 · GOAL-004 S3 · 粗粒度单包）。
 *
 * 渲染闭环：RenderPage（schema 文档 → DOM）+ I18nProvider（i18n 运行时面）
 * + 扩展接缝（registerCustomComponent）+ 核心类型。
 * 内部 bundle 面（components/i18n/lib/protocol）不导出为包级 API。
 */
export {
  RenderPage,
  useSchemaCrud,
  SchemaCrudContext,
  type RendererComponentProps,
  type SchemaCrudFeedback,
  type SchemaCrudConfirm,
  type TableSelection,
  type ActionResult,
  type SchemaCrudValue,
  type RunRequestOptions,
} from "./render.tsx";

export { registerCustomComponent } from "./custom-components";
export type * from "./render.types";

export {
  I18nProvider,
  useI18n,
  getActiveLocale,
  setActiveLocale,
  applyLocaleToDocument,
  type I18nProviderProps,
  type I18nState,
} from "@/i18n/runtime";

export { resolveTextProp, type MessageParams } from "@/i18n/catalog";
export type { Locale, LocalePreference } from "@/i18n/locale";