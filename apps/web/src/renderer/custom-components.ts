// Custom renderer component registry (GOAL-018 · D-001 §1): schema pages may
// embed {type:"custom", component:"<key>"} nodes; the renderer dispatches to
// the registered component. A module-level registry avoids threading a
// customComponents prop through every view component (formComponent-style
// injection is equivalent but touches 8+ call sites). Components receive the
// node and the page context.
import type { ComponentType } from "react";

import type { RenderCustomNode, RenderNode } from "@/renderer/render";

export interface CustomComponentProps {
  node: RenderCustomNode;
  context: Record<string, unknown>;
  /** Children declared under the custom node (rarely used). */
  children?: RenderNode[];
}

const registry = new Map<string, ComponentType<CustomComponentProps>>();

/** Registers a custom node component (idempotent; later registrations win). */
export function registerCustomComponent(
  key: string,
  component: ComponentType<CustomComponentProps>,
): void {
  registry.set(key, component);
}

/** Returns the registered component for a key, or null when unregistered. */
export function getCustomComponent(key: string): ComponentType<CustomComponentProps> | null {
  return registry.get(key) ?? null;
}

/** Test-only: clears the registry. */
export function resetCustomComponentsForTests(): void {
  registry.clear();
}
