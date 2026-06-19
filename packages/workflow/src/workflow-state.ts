import { WorkflowStatus, WorkflowTransitionEvent } from './types.js';

/**
 * Invalid status transition error
 */
export class InvalidStatusTransition extends Error {
  constructor(
    public executionId: string,
    public current: WorkflowStatus,
    public target: WorkflowStatus
  ) {
    super(
      `Invalid workflow transition for execution ${executionId}: ${current} -> ${target}`
    );
    this.name = 'InvalidStatusTransition';
  }
}

/**
 * Transition result
 */
export interface TransitionResult {
  status: WorkflowStatus;
  event: WorkflowTransitionEvent;
}

/**
 * Legal next states for each workflow status
 */
const LEGAL_TRANSITIONS: Record<WorkflowStatus, Set<WorkflowStatus>> = {
  [WorkflowStatus.DRAFT]: new Set([WorkflowStatus.PENDING]),
  [WorkflowStatus.PENDING]: new Set([WorkflowStatus.RUNNING, WorkflowStatus.CANCELLED]),
  [WorkflowStatus.RUNNING]: new Set([
    WorkflowStatus.PAUSED,
    WorkflowStatus.COMPLETED,
    WorkflowStatus.FAILED,
    WorkflowStatus.CANCELLED,
  ]),
  [WorkflowStatus.PAUSED]: new Set([WorkflowStatus.RUNNING, WorkflowStatus.CANCELLED]),
  [WorkflowStatus.COMPLETED]: new Set([]),
  [WorkflowStatus.FAILED]: new Set([WorkflowStatus.PENDING]), // Allow retry
  [WorkflowStatus.CANCELLED]: new Set([]),
};

/**
 * Terminal statuses (no further transitions allowed)
 */
const TERMINAL_STATUSES = new Set([WorkflowStatus.COMPLETED, WorkflowStatus.CANCELLED]);

/**
 * Check if a status transition is valid
 */
export function canTransition(current: WorkflowStatus, target: WorkflowStatus): boolean {
  // Terminal states cannot transition (except failed can retry)
  if (TERMINAL_STATUSES.has(current)) {
    return false;
  }

  // Failed can always transition to pending (retry)
  if (current === WorkflowStatus.FAILED && target === WorkflowStatus.PENDING) {
    return true;
  }

  // Running can always transition to failed
  if (current === WorkflowStatus.RUNNING && target === WorkflowStatus.FAILED) {
    return true;
  }

  // Check legal transitions
  const legalNextStates = LEGAL_TRANSITIONS[current] || new Set();
  return legalNextStates.has(target);
}

/**
 * Transition workflow status with validation
 */
export function transitionWorkflowStatus(
  executionId: string,
  current: WorkflowStatus,
  target: WorkflowStatus,
  actor: string,
  note?: string
): TransitionResult {
  if (!canTransition(current, target)) {
    throw new InvalidStatusTransition(executionId, current, target);
  }

  const event: WorkflowTransitionEvent = {
    execution_id: executionId,
    event_type: 'workflow.status_changed',
    actor,
    from_status: current,
    to_status: target,
    note,
    created_at: new Date().toISOString(),
  };

  return {
    status: target,
    event,
  };
}

/**
 * Get next allowed statuses
 */
export function getNextStatuses(current: WorkflowStatus): WorkflowStatus[] {
  return Array.from(LEGAL_TRANSITIONS[current] || []);
}

/**
 * Check if status is terminal
 */
export function isTerminalStatus(status: WorkflowStatus): boolean {
  return TERMINAL_STATUSES.has(status);
}
