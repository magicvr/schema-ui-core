import type { ComponentType } from "react";
import type { RenderCustomNode, RenderNode } from "@/renderer/render.types";
export interface CustomComponentProps {
    node: RenderCustomNode;
    context: Record<string, unknown>;
    /** Children declared under the custom node (rarely used). */
    children?: RenderNode[];
}
/** Registers a custom node component (idempotent; later registrations win). */
export declare function registerCustomComponent(key: string, component: ComponentType<CustomComponentProps>): void;
/** Returns the registered component for a key, or null when unregistered. */
export declare function getCustomComponent(key: string): ComponentType<CustomComponentProps> | null;
/** Test-only: clears the registry. */
export declare function resetCustomComponentsForTests(): void;
