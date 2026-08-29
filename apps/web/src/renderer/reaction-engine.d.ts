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
    commits: Array<{
        field: string;
        value: unknown;
    }>;
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
export declare const DEFAULT_MAX_ROUNDS = 10;
/**
 * Runs the multi-round engine.
 *
 * A field's committed `value` only lands in `commits` when it deep-differs
 * from the snapshot value; a no-op write neither commits nor schedules the
 * next round. Rounds stop when no committed field is a dependency of any
 * expression (reaction when or observer when).
 */
export declare function runReactionEngine(input: ReactionEngineInput): ReactionEngineResult;
/** Engine + per-field control state (visible/required/disabled/value) from
 * the final round — used by the Renderer's form integration. */
export declare function runReactionEngineDetailed(input: ReactionEngineInput): ReactionEngineDetailed;
