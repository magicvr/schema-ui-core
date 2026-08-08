/**
 * Multi-round snapshot reaction engine (schema-ui-docs@2.7.0 · 02-reaction-expression.md §14).
 *
 * Implements the upstream conformance contract (`conformance/fixtures/reactions/cases.json`):
 *   Snapshot → Evaluate → Commit → Next tick, with:
 *   - per-field reactions [{ when, fulfill, otherwise }] (value/visible/required/disabled)
 *   - observers evaluated against every round snapshot
 *   - baselines + resetMissingOtherwise (condition false without otherwise → restore baseline;
 *     explicit writes win over implicit baseline restores)
 *   - externalUpdates applied between rounds
 *   - deep-equal change detection (no commit / no next round on no-op writes)
 *   - loop protection (REACTION_LOOP_LIMIT after maxRounds, default 10)
 *   - MULTIPLE_VALUE_WRITES warnings (last write wins, deterministic order)
 */

import {
  deepEqual,
  evaluateFullExpression,
  expressionDependencyFields,
  type ReactionEnv,
} from "@/renderer/reaction-expression";

export interface ReactionRuleInput {
  when: string;
  fulfill?: Record<string, unknown>;
  otherwise?: Record<string, unknown>;
}

export interface FieldReactionInput {
  field: string;
  reactions: ReactionRuleInput[];
}

export interface ReactionObserverInput {
  id: string;
  when: string;
}

export interface ExternalUpdateInput {
  afterRound: number;
  values: Record<string, unknown>;
}

export interface ReactionEngineInput {
  initialValues: Record<string, unknown>;
  fields?: FieldReactionInput[];
  observers?: ReactionObserverInput[];
  baselines?: Record<string, unknown>;
  resetMissingOtherwise?: boolean;
  externalUpdates?: ExternalUpdateInput[];
  maxRounds?: number;
}

export interface ReactionRound {
  round: number;
  snapshot: Record<string, unknown>;
  observations: Record<string, boolean>;
  commits: Array<{ field: string; value: unknown }>;
}

export interface MultipleValueWritesWarning {
  code: "MULTIPLE_VALUE_WRITES";
  field: string;
  count: number;
}

export interface ReactionEngineOk {
  ok: true;
  values: Record<string, unknown>;
  rounds: ReactionRound[];
  warnings: MultipleValueWritesWarning[];
}

export interface ReactionEngineLoopError {
  ok: false;
  code: "REACTION_LOOP_LIMIT";
  maxRounds: number;
  values: Record<string, unknown>;
  roundCount: number;
  dependencyFields: string[];
}

export type ReactionEngineResult = ReactionEngineOk | ReactionEngineLoopError;

/** Per-field control state after convergence (renderer integration). */
export interface ReactionFieldState {
  visible?: boolean;
  required?: boolean;
  disabled?: boolean;
}

export interface ReactionEngineDetailed {
  result: ReactionEngineResult;
  fieldStates: Record<string, ReactionFieldState>;
}

export const DEFAULT_MAX_ROUNDS = 10;

type JsonRecord = Record<string, unknown>;

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function copyValues(values: JsonRecord): JsonRecord {
  return JSON.parse(JSON.stringify(values)) as JsonRecord;
}

/** State-key writes that can appear in fulfill/otherwise (02 §6). */
function applyStateKeys(
  target: JsonRecord,
  source: Record<string, unknown> | undefined,
): void {
  if (!isRecord(source)) {
    return;
  }
  for (const key of ["visible", "required", "disabled", "value"] as const) {
    if (source[key] !== undefined) {
      target[key] = source[key];
    }
  }
}

/**
 * Runs the multi-round engine.
 *
 * A field's committed `value` only lands in `commits` when it deep-differs
 * from the snapshot value; a no-op write neither commits nor schedules the
 * next round. Rounds stop when no committed field is a dependency of any
 * expression (reaction when or observer when).
 */
export function runReactionEngine(input: ReactionEngineInput): ReactionEngineResult {
  return runEngineCore(input).result;
}

/** Engine + per-field control state (visible/required/disabled/value) from
 * the final round — used by the Renderer's form integration. */
export function runReactionEngineDetailed(input: ReactionEngineInput): ReactionEngineDetailed {
  return runEngineCore(input);
}

function runEngineCore(input: ReactionEngineInput): ReactionEngineDetailed {
  const maxRounds =
    typeof input.maxRounds === "number" && input.maxRounds > 0 ? input.maxRounds : DEFAULT_MAX_ROUNDS;
  const values: JsonRecord = { ...input.initialValues };
  const fields = input.fields ?? [];
  const observers = input.observers ?? [];
  const baselines = input.baselines ?? {};
  const externalUpdates = input.externalUpdates ?? [];
  const rounds: ReactionRound[] = [];
  const warnings: MultipleValueWritesWarning[] = [];

  // Dependency universe: fields read by any reaction when or observer when.
  const dependencyFields = new Set<string>();
  for (const field of fields) {
    for (const rule of field.reactions) {
      for (const dep of expressionDependencyFields(rule.when)) {
        dependencyFields.add(dep);
      }
    }
  }
  for (const observer of observers) {
    for (const dep of expressionDependencyFields(observer.when)) {
      dependencyFields.add(dep);
    }
  }

  const committedFields = new Set<string>();
  let roundCount = 0;
  let lastFieldStates: Record<string, ReactionFieldState> = {};

  for (let round = 1; round <= maxRounds; round += 1) {
    roundCount = round;
    // Apply external updates scheduled for the previous round.
    for (const update of externalUpdates) {
      if (update.afterRound === round - 1) {
        Object.assign(values, update.values);
      }
    }

    const snapshot = copyValues(values);
    const envBase: ReactionEnv = { deps: snapshot, context: {} };

    const observations: Record<string, boolean> = {};
    for (const observer of observers) {
      observations[observer.id] = evaluateFullExpression(observer.when, envBase);
    }

    // Evaluate each field's rules in declaration order. fulfill/otherwise
    // state keys (visible/required/disabled/value) apply per rule; value
    // writes are last-wins per field.
    const fieldState = new Map<string, JsonRecord>();
    const explicitWriteCounts = new Map<string, number>();
    for (const field of fields) {
      const state: JsonRecord = {};
      let valueWrites = 0;
      for (const rule of field.reactions) {
        const holds = evaluateFullExpression(rule.when, { ...envBase, self: snapshot[field.field] });
        applyStateKeys(state, holds ? rule.fulfill : rule.otherwise);
        const source = holds ? rule.fulfill : rule.otherwise;
        if (isRecord(source) && source.value !== undefined) {
          valueWrites += 1;
        }
      }
      fieldState.set(field.field, state);
      explicitWriteCounts.set(field.field, valueWrites);
    }

    // Implicit baseline restore: condition false + no otherwise + baseline +
    // resetMissingOtherwise → restore baseline value (explicit writes win).
    if (input.resetMissingOtherwise === true) {
      for (const field of fields) {
        const state = fieldState.get(field.field)!;
        if (state.value !== undefined || !(field.field in baselines)) {
          continue;
        }
        const hasRestore = field.reactions.some((rule) => {
          const holds = evaluateFullExpression(rule.when, { ...envBase, self: snapshot[field.field] });
          return !holds && rule.otherwise === undefined;
        });
        if (hasRestore) {
          state.value = baselines[field.field];
        }
      }
    }

    for (const [fieldName, count] of explicitWriteCounts) {
      if (count > 1) {
        warnings.push({ code: "MULTIPLE_VALUE_WRITES", field: fieldName, count });
      }
    }

    // Commit batch: only deep-unequal value changes land in commits.
    const commits: Array<{ field: string; value: unknown }> = [];
    for (const [fieldName, state] of fieldState) {
      if (state.value === undefined) {
        continue;
      }
      if (deepEqual(values[fieldName], state.value)) {
        continue;
      }
      values[fieldName] = state.value;
      commits.push({ field: fieldName, value: state.value });
      committedFields.add(fieldName);
    }

    rounds.push({ round, snapshot, observations, commits });

    lastFieldStates = Object.fromEntries(
      [...fieldState.entries()].map(([fieldName, state]) => {
        const { value: _value, ...control } = state;
        void _value;
        return [fieldName, control];
      }),
    );

    // Next tick when a committed field is read by some expression, or when an
    // external update is scheduled right after this round (it must observe a
    // fresh snapshot, 02 §14).
    const pendingExternalUpdate = externalUpdates.some((update) => update.afterRound === round);
    const schedulesNext = commits.some((commit) => dependencyFields.has(commit.field)) || pendingExternalUpdate;
    if (!schedulesNext) {
      return {
        result: { ok: true, values, rounds, warnings },
        fieldStates: lastFieldStates,
      };
    }
  }

  // Loop protection (02 §14.2): maxRounds reached with a live oscillation.
  return {
    result: {
      ok: false,
      code: "REACTION_LOOP_LIMIT",
      maxRounds,
      values,
      roundCount,
      dependencyFields: [...committedFields],
    },
    fieldStates: lastFieldStates,
  };
}
