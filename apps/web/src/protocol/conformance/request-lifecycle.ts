/**
 * request-lifecycle fixture adapter (latest-wins / hide-drop).
 */

type LifecycleEvent =
  | { type: "start" }
  | { type: "response"; generation: number; state: unknown }
  | { type: "hide" }
  | { type: "show" }
  | { type: "unmount" };

export interface LifecycleResult {
  generation: number;
  active: boolean;
  state: unknown;
  committed: Array<{ generation: number; state: unknown }>;
}

export function runRequestLifecycle(
  initialState: unknown,
  events: LifecycleEvent[],
): LifecycleResult {
  let generation = 0;
  let active = true;
  let state = initialState;
  let mounted = true;
  /** Generations that were cancelled by hide (not reactivated by show). */
  const cancelled = new Set<number>();
  const committed: Array<{ generation: number; state: unknown }> = [];

  for (const event of events) {
    switch (event.type) {
      case "start": {
        if (!mounted) {
          break;
        }
        generation += 1;
        active = true;
        break;
      }
      case "hide": {
        cancelled.add(generation);
        active = false;
        break;
      }
      case "show": {
        // Show does not reactivate a hidden in-flight generation.
        active = true;
        break;
      }
      case "unmount": {
        mounted = false;
        active = false;
        break;
      }
      case "response": {
        if (!mounted) {
          break;
        }
        if (cancelled.has(event.generation)) {
          break;
        }
        if (event.generation !== generation) {
          // Older generation ignored when a newer start exists.
          break;
        }
        state = event.state;
        committed.push({ generation: event.generation, state: event.state });
        break;
      }
      default:
        break;
    }
  }

  return { generation, active, state, committed };
}
