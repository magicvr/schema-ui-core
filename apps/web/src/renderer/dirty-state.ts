/**
 * Small browser-session registry used by the shell navigation guard and
 * schema forms. Sources expose a read-only predicate; the registry never
 * stores form values and therefore cannot leak drafts across pages.
 */

export type DirtyStateSource = () => boolean;

let nextSourceId = 1;
const sources = new Map<number, DirtyStateSource>();

export function registerDirtyStateSource(source: DirtyStateSource): () => void {
  const id = nextSourceId;
  nextSourceId += 1;
  sources.set(id, source);
  return () => sources.delete(id);
}

export function hasDirtyState(): boolean {
  for (const source of sources.values()) {
    try {
      if (source()) return true;
    } catch {
      // A broken source must fail closed: losing a draft is worse than asking
      // the user for confirmation unnecessarily.
      return true;
    }
  }
  return false;
}

export function confirmDiscard(message: string): boolean {
  if (!hasDirtyState()) return true;
  if (typeof window === "undefined" || typeof window.confirm !== "function") {
    return false;
  }
  return window.confirm(message);
}

/** Structural equality for JSON-shaped form values used by dirty tracking. */
export function equalDirtyValues(left: unknown, right: unknown): boolean {
  if (Object.is(left, right)) return true;
  if (typeof left !== typeof right || left === null || right === null) return false;
  if (Array.isArray(left) || Array.isArray(right)) {
    if (!Array.isArray(left) || !Array.isArray(right) || left.length !== right.length) return false;
    return left.every((value, index) => equalDirtyValues(value, right[index]));
  }
  if (typeof left === "object" && typeof right === "object") {
    const leftRecord = left as Record<string, unknown>;
    const rightRecord = right as Record<string, unknown>;
    const leftKeys = Object.keys(leftRecord);
    const rightKeys = Object.keys(rightRecord);
    if (leftKeys.length !== rightKeys.length) return false;
    return leftKeys.every(
      (key) => Object.prototype.hasOwnProperty.call(rightRecord, key) && equalDirtyValues(leftRecord[key], rightRecord[key]),
    );
  }
  return false;
}

/** Test-only seam; production code should unregister through the returned cleanup. */
export function resetDirtyStateSources(): void {
  sources.clear();
  nextSourceId = 1;
}
