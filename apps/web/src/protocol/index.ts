/**
 * @schema-ui/protocol 聚合导出（试点 v0.1.0 · GOAL-004 S2）。
 *
 * 协议面：app-manifest 协商/校验/路由 + 页面文档加载/校验管线。
 * 打包：apps/web Vite lib mode（vite.lib.config.ts）+ tsc declaration。
 */
export * from "./app-manifest";
export * from "./load-page";